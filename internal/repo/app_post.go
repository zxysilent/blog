package repo

import (
	"blog/internal/model"
	"errors"
)

// PostExist 博文是否存在
func PostExist(mod *model.Post) bool {
	sess := db.NewSession()
	defer sess.Close()
	exist, _ := db.Exist(mod)
	return exist
}

// PostGet 博文单条数据
func PostGet(mod *model.Post) error {
	sess := db.NewSession()
	defer sess.Close()
	if mod.Id > 0 {
		sess.ID(mod.Id)
		sess.NoAutoCondition()
	}
	if has, _ := sess.Get(mod); !has {
		return errors.New("No records found")
	}
	return nil
}

// PostList 博文列表数据
func PostList(filter *model.PostFilterList, cols ...string) ([]model.PostPart, error) {
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
	if filter.CateId != nil && *filter.CateId > 0 {
		sess.And("cate_id = ?", *filter.CateId)
	}
	if filter.Limit > 0 {
		sess.Limit(filter.Limit)
	}
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
	if filter.CateId != nil && *filter.CateId > 0 {
		sess.And("cate_id = ?", *filter.CateId)
	}
	sess.Limit(filter.Ps, (filter.Pi-1)*filter.Ps)
	count, err := sess.Desc("id").FindAndCount(&mods)
	return mods, int(count), err
}

// PostCount 博文数据总数
func PostCount() int {
	sess := db.NewSession()
	defer sess.Close()
	// if filter.Mult != "" {
	// 	sess.And("field like concat('%',?,'%')", filter.Mult)
	// }
	count, _ := sess.Count(&model.Post{})
	return int(count)
}

// PostAdd 博文添加数据
func PostAdd(mod *model.Post) error {
	sess := db.NewSession()
	defer sess.Close()
	if _, err := sess.InsertOne(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// PostEdit 博文修改数据
func PostEdit(mod *model.Post, cols ...string) error {
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

// PostDrop 博文删除数据
func PostDrop(mod *model.Post) error {
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

// PostIds 博文id集合查询
func PostIds(ids []string) []model.Post {
	mods := make([]model.Post, 0, len(ids))
	sess := db.NewSession()
	defer sess.Close()
	sess.In("id", ids).Find(&mods)
	return mods
}

// PostMapIds 博文id集合查询
func PostMapIds(ids []int) map[int]*model.Post {
	mods := make([]model.Post, 0, len(ids))
	db.In("id", ids).Find(&mods)
	mapSet := make(map[int]*model.Post, len(mods))
	for idx := range mods {
		mapSet[mods[idx].Id] = &mods[idx]
	}
	return mapSet
}
