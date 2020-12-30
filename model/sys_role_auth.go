package model

import "time"

// SysRoleMenu 角色认证
type SysRoleAuth struct {
	Id     int       `xorm:"pk autoincr INT(11)"` //主键
	RoleId int       `xorm:"default 0 INT(11)"`   //RoleId
	AuthId int       `xorm:"default 0 INT(11)"`   //AuthId
	Ctime  time.Time `xorm:"DATETIME" json:"ctime"`
}

// SysRoleAuthGet 单条角色认证信息
func SysRoleAuthGet(id int) (*SysRoleAuth, bool) {
	mod := &SysRoleAuth{}
	has, _ := Db.ID(id).Get(mod)
	return mod, has
}

// SysRoleAuthAll 所有角色认证信息
func SysRoleAuthAll() ([]SysRoleAuth, error) {
	mods := make([]SysRoleAuth, 0, 8)
	err := Db.Find(&mods)
	return mods, err
}

// SysRoleAuthPage 角色认证分页信息
func SysRoleAuthPage(pi int, ps int, cols ...string) ([]SysRoleAuth, error) {
	mods := make([]SysRoleAuth, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Desc("Utime").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// SysRoleAuthCount 角色认证分页信息总数
func SysRoleAuthCount() int {
	mod := &SysRoleAuth{}
	sess := Db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// SysRoleAuthAdd 添加角色认证信息
func SysRoleAuthAdd(mod *SysRoleAuth) error {
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

// SysRoleAuthEdit 编辑角色认证信息
func SysRoleAuthEdit(mod *SysRoleAuth, cols ...string) error {
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

// SysRoleAuthIds 返回角色认证信息-ids
func SysRoleAuthIds(ids []int) map[int]*SysRoleAuth {
	mods := make([]SysRoleAuth, 0, len(ids))
	Db.In("id", ids).Find(&mods)
	if len(mods) > 0 {
		mapMods := make(map[int]*SysRoleAuth, len(mods))
		for idx := range mods {
			mapMods[mods[idx].Id] = &mods[idx]
		}
		return mapMods
	}
	return nil
}

// SysRoleAuthDrop 删除单条角色认证信息
func SysRoleAuthDrop(id int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&SysRoleAuth{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}
