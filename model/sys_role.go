package model

import "time"

// SysRole 角色表
type SysRole struct {
	Id    int       `xorm:"pk autoincr INT(11)"`   //主键
	Name  string    `xorm:"VARCHAR(255)"`          //名称
	Intro string    `xorm:"VARCHAR(255)"`          //介绍
	Ctime time.Time `xorm:"DATETIME" json:"ctime"` //时间
}

// SysRoleGet 单条角色信息
func SysRoleGet(id int) (*SysRole, bool) {
	mod := &SysRole{}
	has, _ := Db.ID(id).Get(mod)
	return mod, has
}

// SysRoleAll 所有角色信息
func SysRoleAll() ([]SysRole, error) {
	mods := make([]SysRole, 0, 8)
	err := Db.Find(&mods)
	return mods, err
}

// SysRolePage 角色分页信息
func SysRolePage(pi int, ps int, cols ...string) ([]SysRole, error) {
	mods := make([]SysRole, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Desc("Utime").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// SysRoleCount 角色分页信息总数
func SysRoleCount() int {
	mod := &SysRole{}
	sess := Db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// SysRoleAdd 添加角色信息
func SysRoleAdd(mod *SysRole) error {
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

// SysRoleEdit 编辑角色信息
func SysRoleEdit(mod *SysRole, cols ...string) error {
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

// SysRoleIds 返回角色信息-ids
func SysRoleIds(ids []int) map[int]*SysRole {
	mods := make([]SysRole, 0, len(ids))
	Db.In("id", ids).Find(&mods)
	if len(mods) > 0 {
		mapMods := make(map[int]*SysRole, len(mods))
		for idx := range mods {
			mapMods[mods[idx].Id] = &mods[idx]
		}
		return mapMods
	}
	return nil
}

// SysRoleDrop 删除单条角色信息
func SysRoleDrop(id int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&SysRole{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}
