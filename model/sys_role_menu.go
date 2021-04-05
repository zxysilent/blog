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

// RoleMenuGet 单条角色菜单导航信息
func RoleMenuGet(id int) (*RoleMenu, bool) {
	mod := &RoleMenu{}
	has, _ := Db.ID(id).Get(mod)
	return mod, has
}

// RoleMenuAdd 添加角色菜单导航信息
func RoleMenuAdd(mod *RoleMenu) error {
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

// RoleMenuEdit 编辑角色菜单导航信息
func RoleMenuEdit(mod *RoleMenu, cols ...string) error {
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

// RoleMenuDrop 删除单条角色菜单导航信息
func RoleMenuDrop(id int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&RoleMenu{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}
