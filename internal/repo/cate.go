package repo

import (
	"blog/internal/model"
	"context"
	"sort"
	"sync/atomic"

	"github.com/zxysilent/logs"
	"xorm.io/xorm"
)

// cateCache 栏目内存快照
// 栏目数据量小、读远多于写，适合全量常驻内存：
//   - 写操作重建整份快照后原子替换，读路径完全无锁、无分配
//   - 树形渲染、导航、按路径访问都不再查库
//   - 层级（Level）在重建时按父链重算，保证与结构一致
type cateCache struct {
	all    []model.Cate            //前序遍历，父一定在子之前
	byId   map[string]*model.Cate  //id 索引
	byPath map[string]*model.Cate  //公开路径索引
	byName map[string]*model.Cate  //名称索引，用于兼容旧的前台链接
	childs map[string][]model.Cate //pid -> 子级，已按 sort 排序
}

var _cateCache atomic.Pointer[cateCache]

func emptyCateCache() *cateCache {
	return &cateCache{
		all:    make([]model.Cate, 0, 8),
		byId:   make(map[string]*model.Cate, 8),
		byPath: make(map[string]*model.Cate, 8),
		byName: make(map[string]*model.Cate, 8),
		childs: make(map[string][]model.Cate, 8),
	}
}

// cateSnapshot 取当前快照，未初始化时返回空快照（读路径永不 panic）
func cateSnapshot() *cateCache {
	if cc := _cateCache.Load(); cc != nil {
		return cc
	}
	return emptyCateCache()
}

// CateSync 全量重建栏目缓存
// 任何写操作（新增/修改/删除）之后都必须调用
func CateSync() error {
	rows := make([]model.Cate, 0, 32)
	sess := db.NewSession()
	defer sess.Close()
	if err := sess.Asc("sort").Asc("created").Asc("id").Find(&rows); err != nil {
		return err
	}
	cc := emptyCateCache()
	byId := make(map[string]*model.Cate, len(rows))
	for idx := range rows {
		rows[idx].Children = nil
		byId[rows[idx].Id] = &rows[idx]
	}
	// 按父级分组；父级不存在或自引用的一律挂到顶层，避免脏数据整棵消失
	grouped := make(map[string][]model.Cate, len(rows)+1)
	for idx := range rows {
		pid := cateParentOf(&rows[idx], byId)
		grouped[pid] = append(grouped[pid], rows[idx])
	}
	cateSort(grouped)
	// 前序遍历得到稳定顺序与层级；visited 防御环路
	visited := make(map[string]struct{}, len(rows))
	var walk func(pid string, level int)
	walk = func(pid string, level int) {
		items := grouped[pid]
		for idx := range items {
			it := items[idx]
			if _, ok := visited[it.Id]; ok {
				continue
			}
			visited[it.Id] = struct{}{}
			it.Level = level
			it.Children = nil
			cc.all = append(cc.all, it)
			walk(it.Id, level+1)
		}
	}
	walk("", 1)
	// 环路中的节点无法从顶层到达，按顶层兜底追加，保证不丢数据
	for idx := range rows {
		if _, ok := visited[rows[idx].Id]; ok {
			continue
		}
		it := rows[idx]
		it.Level = 1
		it.Children = nil
		cc.all = append(cc.all, it)
	}
	// 用重算过层级的副本重建索引与分组，保证三处数据一致
	cc.childs = make(map[string][]model.Cate, len(grouped))
	for idx := range cc.all {
		it := cc.all[idx]
		cc.byId[it.Id] = &cc.all[idx]
		if it.Path != "" {
			cc.byPath[it.Path] = &cc.all[idx]
		}
		// 名称不保证唯一，保留前序最靠前的一个，仅供旧链接兼容
		if _, ok := cc.byName[it.Name]; !ok {
			cc.byName[it.Name] = &cc.all[idx]
		}
		pid := cateParentOf(&it, cc.byId)
		cc.childs[pid] = append(cc.childs[pid], it)
	}
	_cateCache.Store(cc)
	return nil
}

// cateParentOf 取有效父级id：自引用或父级不存在时视为顶层
func cateParentOf(mod *model.Cate, byId map[string]*model.Cate) string {
	if mod.Pid == "" || mod.Pid == mod.Id {
		return ""
	}
	if _, ok := byId[mod.Pid]; !ok {
		return ""
	}
	return mod.Pid
}

// cateSort 同级排序：sort 升序，其次创建时间、主键，保证多次启动结果一致
func cateSort(grouped map[string][]model.Cate) {
	for pid := range grouped {
		items := grouped[pid]
		sort.SliceStable(items, func(i, j int) bool {
			if items[i].Sort != items[j].Sort {
				return items[i].Sort < items[j].Sort
			}
			if items[i].Created != items[j].Created {
				return items[i].Created < items[j].Created
			}
			return items[i].Id < items[j].Id
		})
		grouped[pid] = items
	}
}

// CateCached 按id取栏目（内存命中）
func CateCached(id string) (*model.Cate, bool) {
	mod, ok := cateSnapshot().byId[id]
	return mod, ok
}

// CateByPath 按公开路径取栏目（内存命中）
func CateByPath(path string) (*model.Cate, bool) {
	mod, ok := cateSnapshot().byPath[path]
	return mod, ok
}

// CateByName 按名称取栏目（内存命中）
// 名称不保证唯一，仅用于兼容旧的前台链接 /cates/<名称>
func CateByName(name string) (*model.Cate, bool) {
	mod, ok := cateSnapshot().byName[name]
	return mod, ok
}

// CateChildren 子级列表
func CateChildren(pid string) []model.Cate {
	items := cateSnapshot().childs[pid]
	out := make([]model.Cate, len(items))
	copy(out, items)
	return out
}

// CateNaver 前台导航：可显示栏目，子级一并带上
func CateNaver() []model.Cate {
	return cateNavOf(cateSnapshot(), "")
}

func cateNavOf(cc *cateCache, pid string) []model.Cate {
	items := cc.childs[pid]
	out := make([]model.Cate, 0, len(items))
	for idx := range items {
		it := items[idx]
		if !it.Show {
			continue
		}
		it.Children = cateNavOf(cc, it.Id)
		out = append(out, it)
	}
	return out
}

// CateTree 完整树形结构
func CateTree() []model.Cate {
	return cateTreeOf(cateSnapshot(), "")
}

func cateTreeOf(cc *cateCache, pid string) []model.Cate {
	items := cc.childs[pid]
	out := make([]model.Cate, 0, len(items))
	for idx := range items {
		it := items[idx]
		it.Children = cateTreeOf(cc, it.Id)
		out = append(out, it)
	}
	return out
}

// CateSubIds 含自身的所有子孙id，用于按栏目树筛选内容
func CateSubIds(id string) []string {
	cc := cateSnapshot()
	ids := make([]string, 0, 8)
	visited := make(map[string]struct{}, 8)
	var walk func(cur string)
	walk = func(cur string) {
		if _, ok := visited[cur]; ok {
			return
		}
		visited[cur] = struct{}{}
		ids = append(ids, cur)
		for _, it := range cc.childs[cur] {
			walk(it.Id)
		}
	}
	walk(id)
	return ids
}

// CateHasChild 是否存在子级，删除前校验
func CateHasChild(id string) bool {
	return len(cateSnapshot().childs[id]) > 0
}

// CateExistPath 公开路径是否已被占用，excludeId 用于编辑时排除自身
func CateExistPath(ctx context.Context, path string, excludeId string) bool {
	if path == "" {
		return false
	}
	if mod, ok := CateByPath(path); ok && mod.Id != excludeId {
		return true
	}
	return false
}

// CateExistName 同级下名称是否重复，excludeId 用于编辑时排除自身
func CateExistName(ctx context.Context, pid, name, excludeId string) bool {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.Where("pid = ? AND name = ?", pid, name)
	if excludeId != "" {
		sess.And("id <> ?", excludeId)
	}
	exist, _ := sess.Exist(&model.Cate{})
	return exist
}

// CateGet 栏目详情
func CateGet(ctx context.Context, id string) (*model.Cate, error) {
	mod := &model.Cate{}
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	has, err := sess.ID(id).NoAutoCondition().Get(mod)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotExist
	}
	return mod, nil
}

// cateCondition 组装栏目查询条件
func cateCondition(sess *xorm.Session, filter *model.CateFilter) {
	if filter.Pid != nil {
		sess.And("pid = ?", *filter.Pid)
	}
	if filter.Kind != nil && *filter.Kind > 0 {
		sess.And("kind = ?", *filter.Kind)
	}
	if filter.Use != nil {
		sess.And("`use` = ?", *filter.Use)
	}
	if filter.Show != nil {
		sess.And("`show` = ?", *filter.Show)
	}
}

// CateList 栏目分页列表
func CateList(ctx context.Context, filter *model.CateListArgs, cols ...string) ([]model.Cate, int, error) {
	mods := make([]model.Cate, 0, filter.GetCap())
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	cateCondition(sess, &filter.CateFilter)
	if filter.Mult != "" {
		sess.And("(name like ? or path like ? or intro like ?)", filter.SafeMult(), filter.SafeMult(), filter.SafeMult())
	}
	sess.Asc("level").Asc("sort").Asc("created").Asc("id")
	if filter.IsZero() {
		err := sess.Find(&mods)
		cateCountFill(ctx, mods)
		return mods, len(mods), err
	}
	count, err := sess.Limit(filter.GetLimit(), filter.GetOffset()).FindAndCount(&mods)
	cateCountFill(ctx, mods)
	return mods, int(count), err
}

// CateAdd 栏目添加
// 写入后重建缓存，并把推导出的层级回写库，保证库与树始终一致
func CateAdd(ctx context.Context, mod *model.Cate) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if _, err := sess.InsertOne(mod); err != nil {
		return err
	}
	return cateRefresh(ctx)
}

// CateEdit 栏目修改
// 父级变化时整棵子树的层级会跟着重算并回写
func CateEdit(ctx context.Context, mod *model.Cate, cols ...string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if _, err := sess.ID(mod.Id).NoAutoCondition().Cols(cols...).Update(mod); err != nil {
		return err
	}
	return cateRefresh(ctx)
}

// CateDrop 栏目删除
// 删除后原本的子级会变成孤儿挂到顶层，层级同样需要重算
func CateDrop(ctx context.Context, id string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if _, err := sess.ID(id).NoAutoCondition().Delete(&model.Cate{}); err != nil {
		return err
	}
	return cateRefresh(ctx)
}

// cateRefresh 重建缓存并回写层级
// 写操作统一入口，调用方不必关心层级的持久化细节
func cateRefresh(ctx context.Context) error {
	if err := CateSync(); err != nil {
		return err
	}
	_, err := CateResyncLevel(ctx)
	return err
}

// CatePostCount 栏目（含子树）下的内容数量
func CatePostCount(ctx context.Context, ids []string) int {
	if len(ids) < 1 {
		return 0
	}
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	count, _ := sess.In("cate_id", ids).Count(&model.Post{})
	return int(count)
}

// cateCountFill 为栏目列表填充内容数量（含子孙栏目）
// 只发一条 GROUP BY 查询，再在内存里自底向上累加，避免 N+1
func cateCountFill(ctx context.Context, mods []model.Cate) {
	if len(mods) < 1 {
		return
	}
	type row struct {
		CateId string `xorm:"cate_id"`
		Total  int    `xorm:"total"`
	}
	rows := make([]row, 0, 32)
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if err := sess.Table("post").Select("cate_id, count(*) as total").
		GroupBy("cate_id").Find(&rows); err != nil {
		return
	}
	own := make(map[string]int, len(rows))
	for idx := range rows {
		own[rows[idx].CateId] = rows[idx].Total
	}
	memo := make(map[string]int, len(mods))
	var sum func(id string) int
	sum = func(id string) int {
		if v, ok := memo[id]; ok {
			return v
		}
		total := own[id]
		for _, it := range CateChildren(id) {
			total += sum(it.Id)
		}
		memo[id] = total
		return total
	}
	for idx := range mods {
		mods[idx].Count = sum(mods[idx].Id)
	}
}

// CateInit 栏目初始化：补齐历史数据缺失的公开路径，重算层级，重建缓存
// 全部操作幂等，可安全地在每次启动时执行
func CateInit(ctx context.Context) error {
	fixed, err := catePathFix(ctx)
	if err != nil {
		return err
	}
	if fixed > 0 {
		logs.Ctx(ctx).Int("count", fixed).Warn("cate.path auto filled with id, please edit them in admin")
	}
	if err := CateSync(); err != nil {
		return err
	}
	levelFixed, err := CateResyncLevel(ctx)
	if err != nil {
		return err
	}
	if levelFixed > 0 {
		logs.Ctx(ctx).Int("count", levelFixed).Info("cate.level recalculated")
	}
	return nil
}

// catePathFix 为缺失公开路径的历史栏目补一个合法且唯一的路径
// 用主键小写作为路径：天然满足 [a-z0-9] 格式，且与 id 一一对应不会撞车
func catePathFix(ctx context.Context) (int, error) {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	mods := make([]model.Cate, 0, 16)
	if err := sess.Where("path = ? OR path IS NULL", "").Find(&mods); err != nil {
		return 0, err
	}
	for idx := range mods {
		it := mods[idx]
		if _, err := sess.ID(it.Id).NoAutoCondition().Cols("path").
			Update(&model.Cate{Path: model.CateNormalizePath(it.Id)}); err != nil {
			return 0, err
		}
	}
	return len(mods), nil
}

// CateResyncLevel 把内存中重算好的层级回写库，只写有变化的行
func CateResyncLevel(ctx context.Context) (int, error) {
	type row struct {
		Id    string `xorm:"id"`
		Level int    `xorm:"level"`
	}
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	rows := make([]row, 0, 32)
	if err := sess.Table("cate").Cols("id", "level").Find(&rows); err != nil {
		return 0, err
	}
	cc := cateSnapshot()
	changed := make([]model.Cate, 0, 8)
	for _, it := range rows {
		want, ok := cc.byId[it.Id]
		if !ok || want.Level == it.Level {
			continue
		}
		changed = append(changed, model.Cate{Id: it.Id, Level: want.Level})
	}
	if len(changed) < 1 {
		return 0, nil
	}
	for idx := range changed {
		if _, err := sess.ID(changed[idx].Id).NoAutoCondition().Cols("level").Update(&changed[idx]); err != nil {
			return 0, err
		}
	}
	return len(changed), nil
}
