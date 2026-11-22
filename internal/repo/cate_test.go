package repo

import (
	"testing"

	"blog/internal/model"
)

// buildCateTree 构造一棵三层的测试栏目树，返回各节点副本
//
//	技术(tech, 显示)  → Go(go, 显示) → Gin(gin, 显示)
//	生活(life, 不显示)
//
// 直接调 repo 层时布尔字段不会走服务层默认值，这里显式给定。
func buildCateTree(t *testing.T) (tech, goCate, gin, life model.Cate) {
	t.Helper()
	ctx := setupCateDB(t)
	tech = *mustCateAdd(t, ctx, model.Cate{Name: "技术", Path: "tech", Sort: 10, Use: true, Show: true})
	goCate = *mustCateAdd(t, ctx, model.Cate{Name: "Go", Path: "go", Pid: tech.Id, Sort: 20, Use: true, Show: true})
	gin = *mustCateAdd(t, ctx, model.Cate{Name: "Gin", Path: "gin", Pid: goCate.Id, Sort: 30, Use: true, Show: true})
	life = *mustCateAdd(t, ctx, model.Cate{Name: "生活", Path: "life", Sort: 5, Use: true})
	return tech, goCate, gin, life
}

// TestCateLevel 层级由父链自动推导
func TestCateLevel(t *testing.T) {
	tech, goCate, gin, life := buildCateTree(t)
	tests := []struct {
		name  string
		id    string
		level int
	}{
		{"顶层", tech.Id, 1},
		{"二级", goCate.Id, 2},
		{"三级", gin.Id, 3},
		{"另一顶层", life.Id, 1},
	}
	for _, it := range tests {
		got, ok := CateCached(it.id)
		if !ok {
			t.Fatalf("%s: 缓存未命中", it.name)
		}
		if got.Level != it.level {
			t.Errorf("%s level = %d, want %d", it.name, got.Level, it.level)
		}
	}
}

// TestCateIndex 按 id / 路径 / 名称三种方式都能命中
func TestCateIndex(t *testing.T) {
	tech, _, _, _ := buildCateTree(t)
	if _, ok := CateCached(tech.Id); !ok {
		t.Error("按 id 未命中")
	}
	if _, ok := CateByPath("tech"); !ok {
		t.Error("按路径未命中")
	}
	// 名称仅用于兼容旧的前台链接
	if _, ok := CateByName("技术"); !ok {
		t.Error("按名称未命中")
	}
	if _, ok := CateByPath("nope"); ok {
		t.Error("不存在的路径竟然命中")
	}
}

// TestCateChildren 子级列表与树形结构
func TestCateChildren(t *testing.T) {
	tech, goCate, gin, _ := buildCateTree(t)
	if n := len(CateChildren("")); n != 2 {
		t.Errorf("顶层子级数量 = %d, want 2", n)
	}
	if n := len(CateChildren(tech.Id)); n != 1 {
		t.Errorf("技术子级数量 = %d, want 1", n)
	}
	children := CateChildren(goCate.Id)
	if len(children) != 1 || children[0].Id != gin.Id {
		t.Fatalf("Go 的子级不为 Gin: %+v", children)
	}
	if n := len(CateChildren(gin.Id)); n != 0 {
		t.Errorf("Gin 不应有子级, got %d", n)
	}
}

// TestCateTreeOrder 同级按 Sort 升序，父一定在子之前
func TestCateTreeOrder(t *testing.T) {
	tech, goCate, gin, life := buildCateTree(t)
	tree := CateTree()
	if len(tree) != 2 {
		t.Fatalf("顶层节点数 = %d, want 2", len(tree))
	}
	// life 的 Sort=5 < tech 的 Sort=10
	if tree[0].Id != life.Id || tree[1].Id != tech.Id {
		t.Errorf("顶层顺序错误: %s, %s", tree[0].Name, tree[1].Name)
	}
	techNode := tree[1]
	if len(techNode.Children) != 1 || techNode.Children[0].Id != goCate.Id {
		t.Fatalf("技术的子级错误: %+v", techNode.Children)
	}
	if len(techNode.Children[0].Children) != 1 || techNode.Children[0].Children[0].Id != gin.Id {
		t.Fatalf("Go 的子级错误: %+v", techNode.Children[0].Children)
	}
}

// TestCateNaver 导航只包含 Show=true 的栏目
func TestCateNaver(t *testing.T) {
	_, _, _, life := buildCateTree(t)
	naver := CateNaver()
	for idx := range naver {
		if naver[idx].Id == life.Id {
			t.Fatal("Show=false 的栏目不应出现在导航中")
		}
	}
	if len(naver) != 1 {
		t.Errorf("导航顶层数量 = %d, want 1", len(naver))
	}
	// 子级同样受 Show 过滤
	if len(naver[0].Children) != 1 {
		t.Errorf("导航子级数量 = %d, want 1", len(naver[0].Children))
	}
}

// TestCateSubIds 含自身的全部子孙 id
func TestCateSubIds(t *testing.T) {
	tech, goCate, gin, life := buildCateTree(t)
	ids := CateSubIds(tech.Id)
	want := map[string]bool{tech.Id: true, goCate.Id: true, gin.Id: true}
	if len(ids) != len(want) {
		t.Fatalf("CateSubIds = %v, want %d 个", ids, len(want))
	}
	for _, id := range ids {
		if !want[id] {
			t.Errorf("CateSubIds 含意外 id %s", id)
		}
	}
	if sub := CateSubIds(life.Id); len(sub) != 1 || sub[0] != life.Id {
		t.Errorf("叶子节点的 CateSubIds = %v, want 仅自身", sub)
	}
}

// TestCateExistPath 路径唯一性校验，编辑时排除自身
func TestCateExistPath(t *testing.T) {
	ctx := setupCateDB(t)
	tech := mustCateAdd(t, ctx, model.Cate{Name: "技术", Path: "tech"})
	if !CateExistPath(ctx, "tech", "") {
		t.Error("已存在的路径应判为占用")
	}
	if CateExistPath(ctx, "tech", tech.Id) {
		t.Error("编辑自身时应排除自身")
	}
	if CateExistPath(ctx, "other", "") {
		t.Error("不存在的路径不应判为占用")
	}
	if CateExistPath(ctx, "", "") {
		t.Error("空路径不应判为占用")
	}
}

// TestCateExistName 同级重名校验，跨级允许同名
func TestCateExistName(t *testing.T) {
	ctx := setupCateDB(t)
	tech := mustCateAdd(t, ctx, model.Cate{Name: "技术", Path: "tech"})
	goCate := mustCateAdd(t, ctx, model.Cate{Name: "Go", Path: "go", Pid: tech.Id})
	if !CateExistName(ctx, "", "技术", "") {
		t.Error("顶层同名应判为重复")
	}
	if CateExistName(ctx, "", "技术", tech.Id) {
		t.Error("编辑自身时应排除自身")
	}
	// 不同父级下允许同名
	if CateExistName(ctx, tech.Id, "技术", "") {
		t.Error("不同父级下的同名不应判为重复")
	}
	if !CateExistName(ctx, goCate.Pid, "Go", "") {
		t.Error("同父级下同名应判为重复")
	}
}

// TestCateHasChild 删除前的子级校验
func TestCateHasChild(t *testing.T) {
	_, goCate, gin, life := buildCateTree(t)
	if !CateHasChild(goCate.Id) {
		t.Error("Go 有子级，应为 true")
	}
	if CateHasChild(gin.Id) {
		t.Error("Gin 无子级，应为 false")
	}
	if CateHasChild(life.Id) {
		t.Error("生活 无子级，应为 false")
	}
}

// TestCateListFilter 列表过滤与分页
func TestCateListFilter(t *testing.T) {
	ctx := setupCateDB(t)
	tech := mustCateAdd(t, ctx, model.Cate{Name: "技术", Path: "tech", Kind: model.KindPost})
	mustCateAdd(t, ctx, model.Cate{Name: "笔记", Path: "note", Kind: model.KindNote})

	kind := model.KindNote
	mods, count, err := CateList(ctx, &model.CateListArgs{CateFilter: model.CateFilter{Kind: &kind}})
	if err != nil {
		t.Fatalf("CateList: %v", err)
	}
	if count != 1 || len(mods) != 1 || mods[0].Name != "笔记" {
		t.Errorf("按 Kind 过滤结果 = %+v, count=%d", mods, count)
	}

	// 关键字命中名称或路径
	mods, _, err = CateList(ctx, &model.CateListArgs{List: model.List{Mult: "tec"}})
	if err != nil {
		t.Fatalf("CateList: %v", err)
	}
	if len(mods) != 1 || mods[0].Id != tech.Id {
		t.Errorf("按关键字过滤结果 = %+v", mods)
	}

	// 按父级过滤
	pid := tech.Id
	mods, _, err = CateList(ctx, &model.CateListArgs{CateFilter: model.CateFilter{Pid: &pid}})
	if err != nil {
		t.Fatalf("CateList: %v", err)
	}
	if len(mods) != 0 {
		t.Errorf("技术下暂无子级, got %+v", mods)
	}
}

// TestCatePostCount 栏目内容数量含子树累加
func TestCatePostCount(t *testing.T) {
	ctx := setupCateDB(t)
	tech := mustCateAdd(t, ctx, model.Cate{Name: "技术", Path: "tech"})
	goCate := mustCateAdd(t, ctx, model.Cate{Name: "Go", Path: "go", Pid: tech.Id})
	gin := mustCateAdd(t, ctx, model.Cate{Name: "Gin", Path: "gin", Pid: goCate.Id})

	for _, it := range []struct {
		path   string
		cateId string
	}{
		{"gin-1", gin.Id},
		{"gin-2", gin.Id},
		{"go-1", goCate.Id},
	} {
		if err := PostAdd(ctx, &model.Post{Path: it.path, Title: it.path, CateId: it.cateId}); err != nil {
			t.Fatalf("PostAdd(%s): %v", it.path, err)
		}
	}
	if n := CatePostCount(ctx, CateSubIds(tech.Id)); n != 3 {
		t.Errorf("技术子树内容数 = %d, want 3", n)
	}
	if n := CatePostCount(ctx, CateSubIds(goCate.Id)); n != 3 {
		t.Errorf("Go 子树内容数 = %d, want 3", n)
	}
	if n := CatePostCount(ctx, CateSubIds(gin.Id)); n != 2 {
		t.Errorf("Gin 内容数 = %d, want 2", n)
	}
}

// TestCateListCount 列表返回的 Count 为该栏目子树的内容总数
func TestCateListCount(t *testing.T) {
	ctx := setupCateDB(t)
	tech := mustCateAdd(t, ctx, model.Cate{Name: "技术", Path: "tech"})
	goCate := mustCateAdd(t, ctx, model.Cate{Name: "Go", Path: "go", Pid: tech.Id})
	if err := PostAdd(ctx, &model.Post{Path: "p1", Title: "p1", CateId: goCate.Id}); err != nil {
		t.Fatalf("PostAdd: %v", err)
	}
	mods, _, err := CateList(ctx, &model.CateListArgs{})
	if err != nil {
		t.Fatalf("CateList: %v", err)
	}
	for _, it := range mods {
		switch it.Id {
		case tech.Id:
			if it.Count != 1 {
				t.Errorf("技术 Count = %d, want 1", it.Count)
			}
		case goCate.Id:
			if it.Count != 1 {
				t.Errorf("Go Count = %d, want 1", it.Count)
			}
		}
	}
}

// TestCateDropCascade 删除栏目后缓存同步，子级审批由服务层负责
func TestCateDrop(t *testing.T) {
	ctx := setupCateDB(t)
	tech := mustCateAdd(t, ctx, model.Cate{Name: "技术", Path: "tech"})
	if err := CateDrop(ctx, tech.Id); err != nil {
		t.Fatalf("CateDrop: %v", err)
	}
	if _, ok := CateCached(tech.Id); ok {
		t.Error("删除后仍能从缓存命中")
	}
	if _, ok := CateByPath("tech"); ok {
		t.Error("删除后路径索引未清理")
	}
	if _, err := CateGet(ctx, tech.Id); err == nil {
		t.Error("删除后仍能查到记录")
	}
}

// TestCateSyncRecoverOrphan 父级不存在的脏数据挂到顶层，不会整棵消失
func TestCateSyncRecoverOrphan(t *testing.T) {
	ctx := setupCateDB(t)
	// 直接构造一个父级不存在的行
	orphan := model.Cate{Name: "孤儿", Path: "orphan", Pid: "not-exist-id"}
	orphan.BeforeInsert()
	if err := CateAdd(ctx, &orphan); err != nil {
		t.Fatalf("CateAdd: %v", err)
	}
	got, ok := CateCached(orphan.Id)
	if !ok {
		t.Fatal("孤儿节点丢失")
	}
	if got.Level != 1 {
		t.Errorf("孤儿节点 level = %d, want 1", got.Level)
	}
	// 应当出现在顶层
	found := false
	for _, it := range CateChildren("") {
		if it.Id == orphan.Id {
			found = true
			break
		}
	}
	if !found {
		t.Error("孤儿节点未挂到顶层")
	}
}

// TestCateSelfCycle 自引用的脏数据不会造成死循环
func TestCateSelfCycle(t *testing.T) {
	ctx := setupCateDB(t)
	self := model.Cate{Name: "自环", Path: "self"}
	self.BeforeInsert()
	self.Pid = self.Id
	if err := CateAdd(ctx, &self); err != nil {
		t.Fatalf("CateAdd: %v", err)
	}
	got, ok := CateCached(self.Id)
	if !ok {
		t.Fatal("自环节点丢失")
	}
	if got.Level != 1 {
		t.Errorf("自环节点 level = %d, want 1", got.Level)
	}
	if len(CateChildren(self.Id)) != 0 {
		t.Error("自环节点不应成为自己的子级")
	}
}

// TestCateEditLevelRecalc 移动节点后整棵子树层级重算并回写库
func TestCateEditLevelRecalc(t *testing.T) {
	ctx := setupCateDB(t)
	a := mustCateAdd(t, ctx, model.Cate{Name: "A", Path: "a"})
	b := mustCateAdd(t, ctx, model.Cate{Name: "B", Path: "b", Pid: a.Id})
	c := mustCateAdd(t, ctx, model.Cate{Name: "C", Path: "c", Pid: b.Id})

	// 把 B 提到顶层，B 与其子树 C 的层级都应减 1
	if err := CateEdit(ctx, &model.Cate{Id: b.Id, Pid: ""}, "pid"); err != nil {
		t.Fatalf("CateEdit: %v", err)
	}
	// CateEdit 内部已完成回写，再次重算应无变化
	if n, err := CateResyncLevel(ctx); err != nil {
		t.Fatalf("CateResyncLevel: %v", err)
	} else if n != 0 {
		t.Errorf("已回写后重算行数 = %d, want 0", n)
	}
	for _, it := range []struct {
		id    string
		level int
	}{
		{a.Id, 1},
		{b.Id, 1},
		{c.Id, 2},
	} {
		got, ok := CateCached(it.id)
		if !ok {
			t.Fatalf("缓存未命中: %s", it.id)
		}
		if got.Level != it.level {
			t.Errorf("缓存 level = %d, want %d", got.Level, it.level)
		}
		dbMod, err := CateGet(ctx, it.id)
		if err != nil {
			t.Fatalf("CateGet: %v", err)
		}
		if dbMod.Level != it.level {
			t.Errorf("库里 level = %d, want %d", dbMod.Level, it.level)
		}
	}
	// 父子关系也应已变更
	if n := len(CateChildren("")); n != 2 {
		t.Errorf("顶层数量 = %d, want 2", n)
	}
}

// TestCateResyncLevelIdempotent 写操作完成后层级已一致，重算不产生写操作
func TestCateResyncLevelIdempotent(t *testing.T) {
	ctx := setupCateDB(t)
	a := mustCateAdd(t, ctx, model.Cate{Name: "A", Path: "a"})
	mustCateAdd(t, ctx, model.Cate{Name: "B", Path: "b", Pid: a.Id})
	n, err := CateResyncLevel(ctx)
	if err != nil {
		t.Fatalf("CateResyncLevel: %v", err)
	}
	if n != 0 {
		t.Errorf("无变化时回写行数 = %d, want 0", n)
	}
}

// TestCateInitBackfillPath 历史数据缺路径时用主键补齐
func TestCateInitBackfillPath(t *testing.T) {
	ctx := setupCateDB(t)
	// 绕过 CateAdd 直接插入一条无路径的行，模拟历史数据
	mod := model.Cate{Name: "老栏目"}
	mod.BeforeInsert()
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if _, err := sess.InsertOne(&mod); err != nil {
		t.Fatalf("InsertOne: %v", err)
	}
	if err := CateInit(ctx); err != nil {
		t.Fatalf("CateInit: %v", err)
	}
	got, ok := CateCached(mod.Id)
	if !ok {
		t.Fatal("缓存未命中")
	}
	if got.Path != model.CateNormalizePath(mod.Id) {
		t.Errorf("补齐路径 = %q, want %q", got.Path, model.CateNormalizePath(mod.Id))
	}
}
