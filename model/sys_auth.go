package model

import "time"

// Auth 认证
type Auth struct {
	Id    int       `xorm:"pk autoincr INT(11)"`   //主键
	Path  string    `xorm:"VARCHAR(255)"`          //路径
	Intro string    `xorm:"VARCHAR(255)"`          //介绍
	Ctime time.Time `xorm:"DATETIME" json:"ctime"` //时间
}

// AuthGet 单条认证信息
func AuthGet(id int) (*Auth, bool) {
	mod := &Auth{}
	has, _ := Db.ID(id).Get(mod)
	return mod, has
}

// AuthAll 所有认证信息
func AuthAll() ([]Auth, error) {
	mods := make([]Auth, 0, 8)
	err := Db.Find(&mods)
	return mods, err
}

// AuthPage 认证分页信息
func AuthPage(pi int, ps int, cols ...string) ([]Auth, error) {
	mods := make([]Auth, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// AuthCount 认证分页信息总数
func AuthCount() int {
	mod := &Auth{}
	sess := Db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// AuthAdd 添加认证信息
func AuthAdd(mod *Auth) error {
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

// AuthEdit 编辑认证信息
func AuthEdit(mod *Auth, cols ...string) error {
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

// AuthIds 返回认证信息-ids
func AuthIds(ids []int) map[int]*Auth {
	mods := make([]Auth, 0, len(ids))
	Db.In("id", ids).Find(&mods)
	if len(mods) > 0 {
		mapMods := make(map[int]*Auth, len(mods))
		for idx := range mods {
			mapMods[mods[idx].Id] = &mods[idx]
		}
		return mapMods
	}
	return nil
}

// AuthDrop 删除单条认证信息
func AuthDrop(id int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&Auth{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}
