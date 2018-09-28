package model

// Cate 分类
type Cate struct {
	Id    int    `xorm:"not null pk autoincr INT(11)" json:"id" form:"id"`
	Name  string `xorm:"not null unique VARCHAR(64)" json:"name" form:"name"`
	Intro string `xorm:"default 'NULL' VARCHAR(64)" json:"intro" form:"intro"`
}

// CateIds 通过id返回新闻类别信息集合
func cateIds(ids []int) map[int]*Cate {
	mods := make([]Cate, 0, 6)
	DB.In("id", ids).Find(&mods)
	if len(mods) > 0 {
		mapSet := make(map[int]*Cate, len(mods))
		for idx := range mods {
			mapSet[mods[idx].Id] = &mods[idx]
		}
		return mapSet
	}
	return nil
}

// CateAll 所有分类
func CateAll() ([]Cate, error) {
	mods := make([]Cate, 0, 4)
	err := DB.Asc("id").Find(&mods)
	return mods, err
}

// CateAdd 添加分类
func CateAdd(mod *Cate) bool {
	sess := DB.NewSession()
	defer sess.Close()
	sess.Begin()
	affect, _ := sess.InsertOne(mod)
	if affect != 1 {
		sess.Rollback()
		return false
	}
	sess.Commit()
	return true
}

// CateEdit 修改分类
func CateEdit(mod *Cate) bool {
	sess := DB.NewSession()
	defer sess.Close()
	sess.Begin()
	affect, err := sess.ID(mod.Id).Cols("Name", "Intro").Update(mod)
	if affect >= 0 && err == nil {
		sess.Commit()
		return true
	}
	sess.Rollback()
	return false
}

// CateDel 删除分类
func CateDel(id int) bool {
	sess := DB.NewSession()
	defer sess.Close()
	sess.Begin()
	if affect, err := sess.ID(id).Delete(&Cate{}); affect > 0 && err == nil {
		sess.Commit()
		DB.ClearCacheBean(&Cate{}, string(id))
		return true
	}
	sess.Rollback()
	return false
}
