package model

import "strconv"

// Tag 标签
type Tag struct {
	Id    int    `xorm:"INT(11) PK AUTOINCR comment('主键')" json:"id"`    //主键
	Name  string `xorm:"UNIQUE VARCHAR(255) comment('标签名')" json:"name"` //标签名
	Intro string `xorm:"VARCHAR(255) comment('描述')" json:"intro"`        //描述
}

// TagGet 单条标签
// int	==>	id
// str	==>	name
func TagGet(id interface{}) (*Tag, bool) {
	mod := &Tag{}
	switch val := id.(type) {
	case int:
		has, _ := db.ID(val).Get(mod)
		return mod, has
	case string:
		has, _ := db.Where("name = ?", val).Get(mod)
		return mod, has
	default:
		return mod, false
	}
}

// TagAll 所有标签信息
func TagAll() ([]Tag, error) {
	mods := make([]Tag, 0, 8)
	err := db.Find(&mods)
	return mods, err
}

// TagPage 标签分页
func TagPage(pi int, ps int, cols ...string) ([]Tag, error) {
	mods := make([]Tag, 0, ps)
	sess := db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Desc("Id").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// TagCount 标签分页总数
func TagCount() int {
	mod := &Tag{}
	sess := db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// TagIds 通过id集合返回标签
func TagIds(ids []int) map[int]*Tag {
	mods := make([]Tag, 0, len(ids))
	db.In("id", ids).Find(&mods)
	mapSet := make(map[int]*Tag, len(mods))
	for idx := range mods {
		mapSet[mods[idx].Id] = &mods[idx]
	}
	return mapSet
}

// TagAdd 添加标签
func TagAdd(mod *Tag) error {
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

// TagEdit 编辑标签
func TagEdit(mod *Tag, cols ...string) error {
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

// TagDrop 删除单条标签
func TagDrop(id int) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&Tag{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	db.ClearCacheBean(&Tag{}, strconv.Itoa(id))
	return nil
}

// ------------------------------------------------------ 前台使用 ------------------------------------------------------

// TagState 统计
type TagState struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Intro string `json:"intro"`
}

// TagStateAll 所有标签统计 当前标签下有文章才显示
func TagStateAll() ([]TagState, error) {
	mods := make([]TagState, 0, 8)
	err := db.SQL("SELECT `name`,intro,count(tag_id) as count FROM post_tag ,tag WHERE tag.id=tag_id GROUP BY tag_id HAVING count>0").Find(&mods)
	return mods, err
}
