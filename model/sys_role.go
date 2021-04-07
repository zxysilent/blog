package model

import "time"

// Role 角色表
type Role struct {
	Id    int       `xorm:"INT(11) PK AUTOINCR comment('主键')" json:"id"`
	Name  string    `xorm:"VARCHAR(255) comment('角色名称')" json:"name"`
	Intro string    `xorm:"VARCHAR(255) comment('角色描述')" json:"intro"`
	Inner bool      `xorm:"TINYINT(4) DEFAULT 0 comment('内部禁止删除')" json:"show"`
	Ctime time.Time `xorm:"DATETIME comment('时间')" json:"ctime"`
}

func (Role) TableName() string {
	return "sys_role"
}

// RoleGet 单条角色信息
func RoleGet(id int) (*Role, bool) {
	mod := &Role{}
	has, _ := Db.ID(id).Get(mod)
	return mod, has
}

// RoleAll 所有角色信息
func RoleAll() ([]Role, error) {
	mods := make([]Role, 0, 8)
	err := Db.Find(&mods)
	return mods, err
}

// RolePage 角色分页信息
func RolePage(pi int, ps int, cols ...string) ([]Role, error) {
	mods := make([]Role, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Desc("Ctime").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// RoleCount 角色分页信息总数
func RoleCount() int {
	mod := &Role{}
	sess := Db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// RoleAdd 添加角色信息
func RoleAdd(mod *Role) error {
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

// RoleEdit 编辑角色信息
func RoleEdit(mod *Role, cols ...string) error {
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

// RoleMapIds 通过id集合返回角色信息
func RoleMapIds(ids []int) map[int]*Role {
	mods := make([]Role, 0, len(ids))
	Db.In("id", ids).Find(&mods)
	mapSet := make(map[int]*Role, len(mods))
	for idx := range mods {
		mapSet[mods[idx].Id] = &mods[idx]
	}
	return mapSet
}

// RoleDrop 删除单条角色信息
func RoleDrop(id int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&Role{}); err != nil {
		sess.Rollback()
		return err
	}
	// 删除角色菜单
	sess.Exec("DELETE FROM sys_role_menu WHERE role_id = ?", id)
	// 删除角色接口
	// sess.Exec("DELETE FROM sys_role_api WHERE role_id = ?", id)
	sess.Commit()
	return nil
}
