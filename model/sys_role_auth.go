package model

import "time"

// RoleMenu 角色认证
type RoleAuth struct {
	Id     int       `xorm:"pk autoincr INT(11)"` //主键
	RoleId int       `xorm:"default 0 INT(11)"`   //RoleId
	AuthId int       `xorm:"default 0 INT(11)"`   //AuthId
	Ctime  time.Time `xorm:"DATETIME" json:"ctime"`
}

// RoleAuthGet 单条角色认证信息
func RoleAuthGet(id int) (*RoleAuth, bool) {
	mod := &RoleAuth{}
	has, _ := Db.ID(id).Get(mod)
	return mod, has
}

// RoleAuthAll 所有角色认证信息
func RoleAuthAll() ([]RoleAuth, error) {
	mods := make([]RoleAuth, 0, 8)
	err := Db.Find(&mods)
	return mods, err
}

// RoleAuthPage 角色认证分页信息
func RoleAuthPage(pi int, ps int, cols ...string) ([]RoleAuth, error) {
	mods := make([]RoleAuth, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Desc("Utime").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// RoleAuthCount 角色认证分页信息总数
func RoleAuthCount() int {
	mod := &RoleAuth{}
	sess := Db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// RoleAuthAdd 添加角色认证信息
func RoleAuthAdd(mod *RoleAuth) error {
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

// RoleAuthEdit 编辑角色认证信息
func RoleAuthEdit(mod *RoleAuth, cols ...string) error {
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

// RoleAuthIds 返回角色认证信息-ids
func RoleAuthIds(ids []int) map[int]*RoleAuth {
	mods := make([]RoleAuth, 0, len(ids))
	Db.In("id", ids).Find(&mods)
	if len(mods) > 0 {
		mapMods := make(map[int]*RoleAuth, len(mods))
		for idx := range mods {
			mapMods[mods[idx].Id] = &mods[idx]
		}
		return mapMods
	}
	return nil
}

// RoleAuthDrop 删除单条角色认证信息
func RoleAuthDrop(id int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&RoleAuth{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}
