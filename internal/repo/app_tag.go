package repo

import (
	"blog/internal/model"
)

// TagExist 标签是否存在
func TagExist(mod *model.Tag) bool {
	sess := db.NewSession()
	defer sess.Close()
	exist, _ := db.Exist(mod)
	return exist
}

// TagGet 标签单条数据
func TagGet(mod *model.Tag) error {
	sess := db.NewSession()
	defer sess.Close()
	if mod.Id > 0 {
		sess.ID(mod.Id)
	}
	if has, _ := sess.Get(mod); !has {
		return ErrNotFound
	}
	return nil
}

// TagList 标签列表数据
func TagList(filter *model.TagFilterList) ([]model.Tag, error) {
	mods := make([]model.Tag, 0, 8)
	sess := db.NewSession()
	defer sess.Close()
	err := sess.Find(&mods)
	return mods, err
}

// TagPage 标签分页数据
func TagPage(filter *model.TagFilterPage, cols ...string) ([]model.Tag, int, error) {
	mods := make([]model.Tag, 0, filter.Ps)
	sess := db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	// if filter.Mult != "" {
	// 	sess.And("field like concat('%',?,'%')", filter.Mult)
	// }
	sess.Limit(filter.Ps, (filter.Pi-1)*filter.Ps)
	count, err := sess.Desc("id").FindAndCount(&mods)
	return mods, int(count), err
}

// TagCount 标签数据总数
func TagCount() int {
	sess := db.NewSession()
	defer sess.Close()
	// if filter.Mult != "" {
	// 	sess.And("field like concat('%',?,'%')", filter.Mult)
	// }
	count, _ := sess.Count(&model.Tag{})
	return int(count)
}

// TagAdd 标签添加数据
func TagAdd(mod *model.Tag) error {
	sess := db.NewSession()
	defer sess.Close()
	if _, err := sess.InsertOne(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// TagEdit 标签修改数据
func TagEdit(mod *model.Tag, cols ...string) error {
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

// TagDrop 标签删除数据
func TagDrop(mod *model.Tag) error {
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

// TagIds 标签id集合查询
func TagIds(ids []string) []model.Tag {
	mods := make([]model.Tag, 0, len(ids))
	sess := db.NewSession()
	defer sess.Close()
	sess.In("id", ids).Find(&mods)
	return mods
}

// TagMapIds 标签id集合查询
func TagMapIds(ids []int) map[int]*model.Tag {
	mods := make([]model.Tag, 0, len(ids))
	db.In("id", ids).Find(&mods)
	mapSet := make(map[int]*model.Tag, len(mods))
	for idx := range mods {
		mapSet[mods[idx].Id] = &mods[idx]
	}
	return mapSet
}

// TagStatus 所有标签统计 当前标签下有文章才显示
func TagStatus() ([]model.TagState, error) {
	mods := make([]model.TagState, 0, 8)
	err := db.SQL("SELECT `name`,intro,count(tag_id) as count FROM post_tag ,tag WHERE tag.id=tag_id GROUP BY tag_id HAVING count>0").Find(&mods)
	return mods, err
}
