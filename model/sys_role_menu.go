package model

import "time"

// RoleMenu 角色菜单导航
type RoleMenu struct {
	Id     int       `xorm:"INT(11) PK AUTOINCR comment('主键')" json:"id"`
	RoleId int       `xorm:"INT(11) DEFAULT 0 comment('角色id')" json:"role_id"`
	MenuId int       `xorm:"INT(11) DEFAULT 0 comment('菜单id')" json:"menu_id"`
	Ctime  time.Time `xorm:"DATETIME comment('时间')" json:"ctime"`
}

func (RoleMenu) TableName() string {
	return "sys_role_menu"
}

// RoleMenuAddMulti 批量添加角色菜单导航信息
func RoleMenuAddMulti(mods []RoleMenu) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.InsertMulti(&mods); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// RoleMenuDropMulti 批量删除角色菜单导航信息
func RoleMenuDropMulti(menuIds []int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.In("menu_id", menuIds).Delete(&RoleMenu{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// RoleMenuEdit 编辑角色菜单导航
// param	mods	"待添加集合"
// param	menuIds	"待删除集合"
func RoleMenuEdit(mods []RoleMenu, menuIds []int) error {
	sess := Db.NewSession()
	defer sess.Close()
	if _, err := sess.InsertMulti(&mods); err != nil {
		sess.Rollback()
		return err
	}
	sess.Begin()
	if _, err := sess.In("menu_id", menuIds).Delete(&RoleMenu{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}
