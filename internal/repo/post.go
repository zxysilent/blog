package repo

import (
	"blog/internal/model"
	"context"

	"xorm.io/xorm"
)

// PostNeighbor 相邻内容，供管理端预览使用
// 只看同栏目、已发布、前台可见的内容，与前台上一篇/下一篇的选取口径一致
// 找不到时返回零值结构，由调用方按 Id 是否为空判断
func PostNeighbor(ctx context.Context, cateId string, created int64) *model.PostNeighbor {
	naver := &model.PostNeighbor{Prev: &model.PostPart{}, Next: &model.PostPart{}}
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)

	prevSess := sess.Where("kind = ? AND status = ? AND `show` = ? AND created > ?",
		model.KindPost, model.StatusFinish, true, created)
	if cateId != "" {
		prevSess.And("cate_id = ?", cateId)
	}
	if has, _ := prevSess.Asc("created").Get(naver.Prev); !has {
		naver.Prev = &model.PostPart{}
	}

	nextSess := sess.Where("kind = ? AND status = ? AND `show` = ? AND created < ?",
		model.KindPost, model.StatusFinish, true, created)
	if cateId != "" {
		nextSess.And("cate_id = ?", cateId)
	}
	if has, _ := nextSess.Desc("created").Get(naver.Next); !has {
		naver.Next = &model.PostPart{}
	}
	return naver
}

// PostGet 内容详情（含栏目）
func PostGet(ctx context.Context, id string) (*model.Post, error) {
	mod, err := GetAny[model.Post](ctx, id)
	if err != nil {
		return nil, err
	}
	PostFillCate(&mod)
	return &mod, nil
}

// PostFillCate 从栏目缓存填充栏目信息，不额外查库
func PostFillCate(mod *model.Post) {
	if mod.CateId == "" {
		return
	}
	cate, ok := CateCached(mod.CateId)
	if !ok {
		return
	}
	copied := *cate
	copied.Children = nil
	mod.Cate = &copied
}

// PostGetPath 按公开路径读取已发布且可见的内容
// 前台详情页使用，未发布或已隐藏一律视为不存在
func PostGetPath(ctx context.Context, path string) (*model.Post, error) {
	mod := &model.Post{}
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	has, err := sess.Where("path = ? AND status = ? AND `show` = ?", path, model.StatusFinish, true).Get(mod)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotExist
	}
	PostFillCate(mod)
	return mod, nil
}

// PostHits 浏览量自增
// 用 SQL 自增而不是先读后写，并发下不会丢计数
func PostHits(ctx context.Context, id string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	_, err := sess.Exec("UPDATE post SET hits = hits + 1 WHERE id = ?", id)
	return err
}

// PostExistPath 访问路径是否已被占用，excludeId 用于编辑时排除自身
func PostExistPath(ctx context.Context, path, excludeId string) bool {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.Where("path = ?", path)
	if excludeId != "" {
		sess.And("id <> ?", excludeId)
	}
	exist, _ := sess.Exist(&model.Post{})
	return exist
}

// postCondition 组装内容查询条件，管理端与前台共用同一套语义
func postCondition(sess *xorm.Session, filter *model.PostFilter) {
	if filter.Kind != nil && *filter.Kind > 0 {
		sess.And("kind = ?", *filter.Kind)
	}
	if filter.Status != nil && *filter.Status > 0 {
		sess.And("status = ?", *filter.Status)
	}
	if filter.Show != nil {
		sess.And("`show` = ?", *filter.Show)
	}
	if filter.Istop != nil {
		sess.And("istop = ?", *filter.Istop)
	}
	if filter.Uid != nil && *filter.Uid != "" {
		sess.And("uid = ?", *filter.Uid)
	}
	if filter.Pid != nil && *filter.Pid != "" {
		sess.And("pid = ?", *filter.Pid)
	}
	// CateIds 由服务层按「栏目 + 子树」展开，优先于单个栏目条件
	if len(filter.CateIds) > 0 {
		sess.In("cate_id", filter.CateIds)
	} else if filter.CateId != nil && *filter.CateId != "" {
		sess.And("cate_id = ?", *filter.CateId)
	}
}

// PostList 管理端内容列表
// 默认不返回 markdown / richtext 两个大字段，需要时通过 cols 显式指定
func PostList(ctx context.Context, filter *model.PostListArgs, cols ...string) ([]model.Post, int, error) {
	mods := make([]model.Post, 0, filter.GetCap())
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if len(cols) > 0 {
		sess.Cols(cols...)
	} else {
		sess.Omit("markdown", "richtext")
	}
	postCondition(sess, &filter.PostFilter)
	if filter.Mult != "" {
		sess.And("(title like ? or path like ? or summary like ?)",
			filter.SafeMult(), filter.SafeMult(), filter.SafeMult())
	}
	sess.Desc("istop").Desc("created").Desc("id")
	if filter.IsZero() {
		err := sess.Find(&mods)
		return mods, len(mods), err
	}
	count, err := sess.Limit(filter.GetLimit(), filter.GetOffset()).FindAndCount(&mods)
	return mods, int(count), err
}

// PostCount 内容计数
func PostCount(ctx context.Context, filter *model.PostFilter) int {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	postCondition(sess, filter)
	count, _ := sess.Count(&model.Post{})
	return int(count)
}

// PostAdd 内容添加
func PostAdd(ctx context.Context, mod *model.Post) error {
	return AddAny(ctx, mod)
}

// PostEdit 内容修改
func PostEdit(ctx context.Context, mod *model.Post, cols ...string) error {
	return EditAny(ctx, mod, cols...)
}

// PostEditTop 修改置顶
func PostEditTop(ctx context.Context, id string, istop bool) error {
	return PostEdit(ctx, &model.Post{Id: id, Istop: istop}, "istop")
}

// PostDrop 内容删除
// 与标签关联的清理放在同一个事务里，避免出现「文章已删、关联还在」的中间态
func PostDrop(ctx context.Context, id string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if err := sess.Begin(); err != nil {
		return err
	}
	if _, err := sess.Where("post_id = ?", id).Delete(&model.PostTag{}); err != nil {
		sess.Rollback()
		return err
	}
	if _, err := sess.ID(id).NoAutoCondition().Delete(&model.Post{}); err != nil {
		sess.Rollback()
		return err
	}
	return sess.Commit()
}

// PostPage 摘要列表，分页
func PostPage(filter *model.PostFilterPage, cols ...string) ([]model.PostPart, int, error) {
	mods := make([]model.PostPart, 0, filter.GetCap())
	sess := db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	postCondition(sess, &filter.PostFilter)
	if filter.Mult != "" {
		sess.And("title like ?", filter.SafeMult())
	}
	count, err := sess.Limit(filter.GetLimit(), filter.GetOffset()).
		Desc("istop").Desc("created").Desc("id").FindAndCount(&mods)
	return mods, int(count), err
}

// PostTree 摘要树形列表，按 Pid 自关联组装
// 笔记类内容依赖它做目录树；非树形内容平铺返回即可
func PostTree(filter *model.PostFilterPage, cols ...string) ([]model.PostPart, int, error) {
	mods := make([]model.PostPart, 0, filter.GetCap())
	sess := db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	postCondition(sess, &filter.PostFilter)
	if filter.Mult != "" {
		sess.And("title like ?", filter.SafeMult())
	}
	count, err := sess.Desc("istop").Desc("created").Desc("id").FindAndCount(&mods)
	if err != nil {
		return nil, 0, err
	}
	return postBuildTree(mods), int(count), nil
}

// postBuildTree 把平铺结果组装成树
// 父级不在结果集中的节点提升为顶层，避免数据被藏起来
func postBuildTree(mods []model.PostPart) []model.PostPart {
	byId := make(map[string]*model.PostPart, len(mods))
	for idx := range mods {
		mods[idx].Children = nil
		byId[mods[idx].Id] = &mods[idx]
	}
	tree := make([]model.PostPart, 0, len(mods))
	for idx := range mods {
		mod := mods[idx]
		if mod.Pid != "" {
			if parent, ok := byId[mod.Pid]; ok && parent.Id != mod.Id {
				parent.Children = append(parent.Children, mod)
				continue
			}
		}
		tree = append(tree, mod)
	}
	// 上一步 append 的是值拷贝，这里按 id 回填子级
	var fill func(nodes []model.PostPart, level int) []model.PostPart
	fill = func(nodes []model.PostPart, level int) []model.PostPart {
		out := make([]model.PostPart, 0, len(nodes))
		seen := make(map[string]struct{}, len(nodes))
		for _, node := range nodes {
			if _, ok := seen[node.Id]; ok {
				continue
			}
			seen[node.Id] = struct{}{}
			node.Level = level
			if parent, ok := byId[node.Id]; ok {
				node.Children = fill(parent.Children, level+1)
			} else {
				node.Children = nil
			}
			out = append(out, node)
		}
		return out
	}
	return fill(tree, 1)
}

// PostNaver 上下篇
// cateId 非空时只在同一栏目内找，避免跨栏目跳转
func PostNaver(cateId string, created int64) *model.Naver {
	naver := &model.Naver{}
	prev := &model.Post{}
	prevSess := db.Where("kind = ? AND status = ? AND `show` = ? AND created > ?",
		model.KindPost, model.StatusFinish, true, created)
	if cateId != "" {
		prevSess.And("cate_id = ?", cateId)
	}
	if has, _ := prevSess.Asc("created").Get(prev); has {
		naver.Prev = `<a href="/posts/` + prev.Path + `.html" class="prev">&laquo; ` + prev.Title + `</a>`
	}
	next := &model.Post{}
	nextSess := db.Where("kind = ? AND status = ? AND `show` = ? AND created < ?",
		model.KindPost, model.StatusFinish, true, created)
	if cateId != "" {
		nextSess.And("cate_id = ?", cateId)
	}
	if has, _ := nextSess.Desc("created").Get(next); has {
		naver.Next = `<a href="/posts/` + next.Path + `.html" class="next"> ` + next.Title + ` &raquo;</a>`
	}
	return naver
}
