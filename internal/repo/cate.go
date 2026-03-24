package repo

import (
	"blog/internal/model"
	"context"
	"reflect"
)

// CateExist 分类是否存在
func CateExist(ctx context.Context, mod *model.Cate) bool {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	exist, _ := sess.Exist(mod)
	return exist
}

// CateGet 分类单条数据
func CateGet(ctx context.Context, mod *model.Cate) error {
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

// CateList 分类列表数据
func CateList(ctx context.Context, filter *model.CateListArgs, cols ...string) ([]model.Cate, int, error) {
	mods := make([]model.Cate, 0, filter.GetCap())
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

// CateCount 分类数据总数
func CateCount(ctx context.Context, filter *model.CateListArgs) int {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	// if filter.Mult != "" {
	// 	sess.And("field like concat('%',?,'%')", filter.Mult)
	// }
	count, _ := sess.Count(&model.Cate{})
	return int(count)
}

// CateAdd 分类添加数据
func CateAdd(ctx context.Context, mod *model.Cate) error {
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

// CateEdit 分类修改数据
func CateEdit(ctx context.Context, mod *model.Cate, cols ...string) error {
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

// CateDrop 分类删除数据
func CateDrop(ctx context.Context, mod *model.Cate) error {
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

// CateIds 分类id集合查询
func CateIds(ctx context.Context, ids []string) []model.Cate {
	mods := make([]model.Cate, 0, len(ids))
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.In("id", ids).Find(&mods)
	return mods
}

// CateMapIds 分类id集合查询
func CateMapIds(ctx context.Context, ids []int) map[string]*model.Cate {
	mods := make([]model.Cate, 0, len(ids))
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.In("id", ids).Find(&mods)
	set := make(map[string]*model.Cate, len(mods))
	for idx := range mods {
		set[mods[idx].Id] = &mods[idx]
	}
	return set
}
