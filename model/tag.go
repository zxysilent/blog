package model

// Tag 标签
type Tag struct {
	Id    int    `xorm:"pk autoincr INT(11)" json:"id" form:"id"`
	Name  string `xorm:"unique VARCHAR(64)" json:"name" form:"name"`
	Intro string `xorm:"VARCHAR(64)" json:"intro" form:"intro"`
}

// TagState 统计
type TagState struct {
	Name  string `xorm:"VARCHAR(64)" json:"name" form:"name"`
	Count int    `xorm:"INT(11)" json:"count" form:"count"`
}

// tagIds 通过id返回标签集合
func tagIds(ids []int) map[int]*Tag {
	mods := make([]Tag, 0, 6)
	Db.In("id", ids).Find(&mods)
	if len(mods) > 0 {
		mapSet := make(map[int]*Tag, len(mods))
		for idx := range mods {
			mapSet[mods[idx].Id] = &mods[idx]
		}
		return mapSet
	}
	return nil
}

// TagName 通过name 查询标签
func TagName(nam string) (*Tag, bool) {
	mod := &Tag{
		Name: nam,
	}
	has, _ := Db.Get(mod)
	return mod, has
}

// TagAll 所有标签
func TagAll() ([]Tag, error) {
	mods := make([]Tag, 0, 8)
	err := Db.Asc("id").Find(&mods)
	return mods, err
}

// TagStateAll 所有标签统计 当前标签下有文章才显示
func TagStateAll() ([]TagState, error) {
	mods := make([]TagState, 0, 8)
	err := Db.SQL("SELECT `name`,count(tag_id) as count FROM post_tag ,tag WHERE tag.id=tag_id GROUP BY tag_id HAVING count>0").Find(&mods)
	return mods, err
}

// TagAdd 添加标签
func TagAdd(mod *Tag) bool {
	sess := Db.NewSession()
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

// TagEdit 修改标签
func TagEdit(mod *Tag) bool {
	sess := Db.NewSession()
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

// TagDel 删除标签
func TagDel(id int) bool {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if affect, err := sess.ID(id).Delete(&Tag{}); affect > 0 && err == nil {
		sess.Commit()
		Db.ClearCacheBean(&Tag{}, string(id))
		return true
	}
	sess.Rollback()
	return false
}
