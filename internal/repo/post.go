package repo

import (
	"blog/internal/model"
	"context"
	"reflect"
)

// PostExist 博文是否存在
func PostExist(ctx context.Context, mod *model.Post) bool {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	exist, _ := sess.Exist(mod)
	return exist
}

// PostGet 博文单条数据
func PostGet(ctx context.Context, mod *model.Post) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if !reflect.ValueOf(mod.Id).IsZero() {
		sess.ID(mod.Id)
		sess.NoAutoCondition()
	}
	if has, _ := sess.Get(mod); !has {
		return ErrNotExist
	}
	return nil
}

// PostList 博文列表数据
func PostList(ctx context.Context, filter *model.PostListArgs, cols ...string) ([]model.Post, int, error) {
	mods := make([]model.Post, 0, filter.GetCap())
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	sess.Asc("created")
	if filter.IsZero() {
		err := sess.Find(&mods)
		return mods, len(mods), err
	}
	sess.Limit(filter.GetLimit(), filter.GetOffset())
	count, err := sess.FindAndCount(&mods)
	return mods, int(count), err
}

// PostCount 博文数据总数
func PostCount(ctx context.Context, filter *model.PostListArgs) int {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	// if filter.Mult != "" {
	// 	sess.And("field like concat('%',?,'%')", filter.Mult)
	// }
	count, _ := sess.Count(&model.Post{})
	return int(count)
}

// PostAdd 博文添加数据
func PostAdd(ctx context.Context, mod *model.Post) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if _, err := sess.InsertOne(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// PostEdit 博文修改数据
func PostEdit(ctx context.Context, mod *model.Post, cols ...string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.Begin()
	if _, err := sess.ID(mod.Id).Cols(cols...).Update(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// PostDrop 博文删除数据
func PostDrop(ctx context.Context, mod *model.Post) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.Begin()
	if !reflect.ValueOf(mod.Id).IsZero() {
		sess.ID(mod.Id)
		sess.NoAutoCondition()
	}
	if _, err := sess.Delete(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// PostIds 博文id集合查询
func PostIds(ctx context.Context, ids []string) []model.Post {
	mods := make([]model.Post, 0, len(ids))
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.In("id", ids).Find(&mods)
	return mods
}

// PostMapIds 博文id集合查询
func PostMapIds(ctx context.Context, ids []string) map[string]*model.Post {
	mods := make([]model.Post, 0, len(ids))
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.In("id", ids).Find(&mods)
	set := make(map[string]*model.Post, len(mods))
	for idx := range mods {
		set[mods[idx].Id] = &mods[idx]
	}
	return set
}

// PostList 博文列表数据
func PostList1(filter *model.PostFilterList, cols ...string) ([]model.PostPart, error) {
	mods := make([]model.PostPart, 0, 8)
	sess := db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	if filter.Mult != "" {
		sess.And("title like ?", filter.SafeMult())
	}
	if filter.Kind != nil && *filter.Kind > 0 {
		sess.And("kind = ?", *filter.Kind)
	}
	if filter.Status != nil && *filter.Status > 0 {
		sess.And("status = ?", *filter.Status)
	}
	if filter.CateId != nil && *filter.CateId != "" {
		sess.And("cate_id = ?", *filter.CateId)
	}
	// if filter.Limit > 0 {
	// 	sess.Limit(filter.Limit)
	// }
	err := sess.Desc("id").Find(&mods)
	return mods, err
}

// PostPage 博文分页数据
func PostPage(filter *model.PostFilterPage, cols ...string) ([]model.PostPart, int, error) {
	mods := make([]model.PostPart, 0, filter.Ps)
	sess := db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	if filter.Mult != "" {
		sess.And("title like ?", filter.SafeMult())
	}
	if filter.Kind != nil && *filter.Kind > 0 {
		sess.And("kind = ?", *filter.Kind)
	}
	if filter.Status != nil && *filter.Status > 0 {
		sess.And("status = ?", *filter.Status)
	}
	if filter.CateId != nil && *filter.CateId != "" {
		sess.And("cate_id = ?", *filter.CateId)
	}
	sess.Limit(filter.Ps, (filter.Pi-1)*filter.Ps)
	count, err := sess.Desc("id").FindAndCount(&mods)
	return mods, int(count), err
}

func PostTree(filter *model.PostFilterPage, cols ...string) ([]model.PostPart, int, error) {
	mods := make([]model.PostPart, 0, filter.Ps)
	sess := db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	if filter.Mult != "" {
		sess.And("title like ?", filter.SafeMult())
	}
	if filter.Kind != nil && *filter.Kind > 0 {
		sess.And("kind = ?", *filter.Kind)
	}
	if filter.Status != nil && *filter.Status > 0 {
		sess.And("status = ?", *filter.Status)
	}
	if filter.CateId != nil && *filter.CateId != "" {
		sess.And("cate_id = ?", *filter.CateId)
	}
	sess.Limit(filter.Ps, (filter.Pi-1)*filter.Ps)
	count, err := sess.Desc("id").FindAndCount(&mods)
	return mods, int(count), err
}

func PostNaver(cateId string, crated int64) *model.Naver {
	naver := &model.Naver{}
	prev := &model.Post{}
	if has, _ := db.Where("kind = ? and status = ? and created >?", model.KindPost, model.StatusFinish, crated).Asc("created").Get(prev); has {
		// <a href="{{.Naver.Prev}}" class="prev">&laquo; 上一页</a>
		naver.Prev = `<a href="/posts/` + prev.Path + `.html" class="prev">&laquo; ` + prev.Title + `</a>`
	}
	next := &model.Post{}
	if has, _ := db.Where("kind = ? and status = ? and created <?", model.KindPost, model.StatusFinish, crated).Desc("created").Get(next); has {
		//<a href="{{.Naver.Next}}" class="next">下一页 &raquo;</a>
		naver.Next = `<a href="/posts/` + next.Path + `.html" class="next"> ` + next.Title + ` &raquo;</a>`
	}
	return naver
}

// NoteList 笔记列表数据
func NoteList(filter *model.List, cols ...string) ([]model.PostPart, error) {
	arr := make([]model.PostPart, 0, 8)
	sess := db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	if filter.Mult != "" {
		sess.And("title like ?", filter.SafeMult())
	}
	sess.And("kind = ?", model.KindNote)
	sess.Desc("created").Find(&arr)
	mp := make(map[int][]model.PostPart, len(arr))
	for _, val := range arr {
		if _, ok := mp[val.Pid]; !ok {
			mp[val.Pid] = make([]model.PostPart, 0, 4)
		}
		mp[val.Pid] = append(mp[val.Pid], val)
	}
	level := 1
	var getChildren func(pid int, mp map[int][]model.PostPart) []model.PostPart
	getChildren = func(pid int, mp map[int][]model.PostPart) []model.PostPart {
		lv := level
		mods, ok := mp[pid]
		if !ok {
			return nil
		}
		level++
		for idx := range mods {
			mods[idx].Level = lv
			mods[idx].Children = getChildren(mods[idx].Id, mp)
		}
		return mods
	}
	mods := getChildren(0, mp)
	// for idx := range mods {
	// 	mods[idx].Children = getChildren(&mods[idx], mp)
	// }
	// for _, sub1 := range first {
	// 	second := mp[sub1.Id]
	// 	sub1.Children = make([]model.PostPart, 0, len(second))
	// 	for _, sub2 := range second {
	// 		third := mp[sub2.Id]
	// 		sub2.Children = make([]model.PostPart, 0, len(third))
	// 		for _, sub3 := range third {
	// 			sub2.Children = append(sub2.Children, sub3)
	// 		}
	// 		sub1.Children = append(sub1.Children, sub2)
	// 	}
	// 	mods = append(mods, sub1)
	// }
	return mods, nil
}
