package repo

import (
	"blog/internal/model"
	"errors"
)

// CateExist 分类是否存在
func CateExist(mod *model.Cate) bool {
	sess := db.NewSession()
	defer sess.Close()
	exist, _ := db.Exist(mod)
	return exist
}

// CateGet 分类单条数据
func CateGet(mod *model.Cate) error {
	sess := db.NewSession()
	defer sess.Close()
	if mod.Id > 0 {
		sess.ID(mod.Id)
	}
	if has, _ := sess.Get(mod); !has {
		return errors.New("No records found")
	}
	return nil
}

// CateList 分类列表数据
func CateList(filter *model.CateFilterList) ([]model.Cate, error) {
	mods := make([]model.Cate, 0, 8)
	sess := db.NewSession()
	defer sess.Close()
	if filter.Kind != nil && *filter.Kind > 0 {
		sess.And("kind = ?", *filter.Kind)
	}
	err := sess.Omit("updated", "created").Find(&mods)
	return mods, err
}

// CateList 分类列表数据
func CateTree(filter *model.CateFilterTree) ([]model.Cate, error) {
	mods := make([]model.Cate, 0, 8)
	sess := db.NewSession()
	defer sess.Close()
	if filter.Kind != nil && *filter.Kind > 0 {
		sess.And("kind = ?", *filter.Kind)
	}
	err := sess.Omit("updated", "created").Find(&mods)
	return mods, err
}

// CatePage 分类分页数据
func CatePage(filter *model.CateFilterPage, cols ...string) ([]model.Cate, int, error) {
	mods := make([]model.Cate, 0, filter.Ps)
	sess := db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	// if filter.Mult != "" {
	// 	sess.And("field like concat('%',?,'%')", filter.Mult)
	// }
	if filter.Kind != nil && *filter.Kind > 0 {
		sess.And("kind = ?", *filter.Kind)
	}
	sess.Limit(filter.Ps, (filter.Pi-1)*filter.Ps)
	count, err := sess.Desc("id").FindAndCount(&mods)
	return mods, int(count), err
}

// CateCount 分类数据总数
func CateCount() int {
	sess := db.NewSession()
	defer sess.Close()
	// if filter.Mult != "" {
	// 	sess.And("field like concat('%',?,'%')", filter.Mult)
	// }
	count, _ := sess.Count(&model.Cate{})
	return int(count)
}

// CateAdd 分类添加数据
func CateAdd(mod *model.Cate) error {
	sess := db.NewSession()
	defer sess.Close()
	if _, err := sess.InsertOne(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// CateEdit 分类修改数据
func CateEdit(mod *model.Cate, cols ...string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(mod.Id).Cols(cols...).Update(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// CateDrop 分类删除数据
func CateDrop(mod *model.Cate) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	if mod.Id > 0 {
		sess.ID(mod.Id)
	}
	if _, err := sess.Delete(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// CateIds 分类id集合查询
func CateIds(ids []string) []model.Cate {
	mods := make([]model.Cate, 0, len(ids))
	sess := db.NewSession()
	defer sess.Close()
	sess.In("id", ids).Find(&mods)
	return mods
}

// CateMapIds 分类id集合查询
func CateMapIds(ids []int) map[int]*model.Cate {
	mods := make([]model.Cate, 0, len(ids))
	db.In("id", ids).Find(&mods)
	mapSet := make(map[int]*model.Cate, len(mods))
	for idx := range mods {
		mapSet[mods[idx].Id] = &mods[idx]
	}
	return mapSet
}
