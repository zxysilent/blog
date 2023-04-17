package repo

import "blog/internal/model"

// DictGet 单条字典
func DictGet(key string) (*model.Dict, bool) {
	mod := &model.Dict{}
	has, _ := db.ID(key).Get(mod)
	return mod, has
}

// DictExist 是否存在
func DictExist(key string) bool {
	has, _ := db.Exist(&model.Dict{
		Key: key,
	})
	return has
}

// DictPage 字典分页
func DictPage(pi int, ps int, cols ...string) ([]model.Dict, error) {
	mods := make([]model.Dict, 0, ps)
	sess := db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Desc("created").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// DictCount 字典分页总数
func DictCount() int {
	mod := &model.Dict{}
	sess := db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// DictAdd 添加字典
func DictAdd(mod *model.Dict) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.InsertOne(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// DictEdit 编辑字典
func DictEdit(mod *model.Dict, cols ...string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	sess.Omit("key")
	if _, err := sess.ID(mod.Key).Cols(cols...).Update(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	if mod.Key == model.DictKeyBasic {
		initBasic()
	}
	if mod.Key == model.DictKeyBlog {
		initBlog()
	}
	return nil
}

// DictDrop 删除单条字典
func DictDrop(key string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(key).Delete(&model.Dict{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}
