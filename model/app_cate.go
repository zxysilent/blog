package model

import (
	"strconv"
)

// Cate 分类
type Cate struct {
	Id    int    `xorm:"INT(11) PK AUTOINCR comment('主键')" json:"id"`    //主键
	Name  string `xorm:"UNIQUE VARCHAR(255) comment('分类名')" json:"name"` //分类名
	Intro string `xorm:"VARCHAR(255) comment('描述')" json:"intro"`        //描述
}

// CateGet 单条分类
// int	==>	id
// str	==>	name
func CateGet(id interface{}) (*Cate, bool) {
	mod := &Cate{}
	switch val := id.(type) {
	case int:
		has, _ := Db.ID(val).Get(mod)
		return mod, has
	case string:
		has, _ := Db.Where("Name = ?", val).Get(mod)
		return mod, has
	default:
		return mod, false
	}
}

// CateAll 所有分类
func CateAll() ([]Cate, error) {
	mods := make([]Cate, 0, 8)
	err := Db.Find(&mods)
	return mods, err
}

// CatePage 分类分页
func CatePage(pi int, ps int, cols ...string) ([]Cate, error) {
	mods := make([]Cate, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Desc("Id").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// CateCount 分类分页总数
func CateCount() int {
	mod := &Cate{}
	sess := Db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// CateIds 通过id集合返回分类
func CateIds(ids []int) map[int]*Cate {
	mods := make([]Cate, 0, len(ids))
	Db.In("id", ids).Find(&mods)
	mapSet := make(map[int]*Cate, len(mods))
	for idx := range mods {
		mapSet[mods[idx].Id] = &mods[idx]
	}
	return mapSet
}

// CateAdd 添加分类
func CateAdd(mod *Cate) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.InsertOne(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// CateEdit 编辑分类
func CateEdit(mod *Cate, cols ...string) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(mod.Id).Cols(cols...).Update(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// CateDrop 删除单条分类
func CateDrop(id int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&Cate{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	Db.ClearCacheBean(&Cate{}, strconv.Itoa(id))
	return nil
}

// ------------------------------------------------------ 前台使用 ------------------------------------------------------

// CatePostCount 通过标签查询文章分页总数
// lmt 是否前台限制
func CatePostCount(cid int, lmt bool) int {
	sess := Db.NewSession()
	defer sess.Close()
	if cid > 0 {
		sess.Where("Cate_id = ?  ", cid)
	}
	if lmt {
		sess.Where("Is_Public = 1 and Status = 3 ")
	}
	sess.Where("Type = 0")
	mod := &Post{}
	count, _ := sess.Count(mod)
	return int(count)
}

// CatePostList 通过分类查询文章分页
// lmt 是否前台限制
func CatePostList(cid, pi, ps int, lmt bool) ([]Post, error) {
	mods := make([]Post, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if cid > 0 {
		sess.Where("Cate_id = ? ", cid)
	}
	if lmt {
		sess.Where("Is_Public = 1 and Status = 3 ")
	}
	sess.Where("Type = 0")
	err := sess.Cols("id", "title", "path", "create_time", "summary", "comment_num", "options", "Status", "Is_Public").Desc("create_time").Limit(ps, (pi-1)*ps).Find(&mods)
	if len(mods) > 0 {
		for idx := range mods {
			if !lmt {
				mods[idx].Summary = ""
			}
			mods[idx].Richtext = ""
			mods[idx].Markdown = ""
		}
	}
	return mods, err
}
