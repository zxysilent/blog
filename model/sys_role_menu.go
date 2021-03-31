package model

import "time"

// RoleMenu 角色菜单导航
type RoleMenu struct {
	Id     int       `xorm:"pk autoincr INT(11)"` //主键
	RoleId int       `xorm:"default 0 INT(11)"`   //RoleId
	MenuId int       `xorm:"default 0 INT(11)"`   //MenuId
	Ctime  time.Time `xorm:"DATETIME" json:"ctime"`
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

// RoleMenuAll 所有角色菜单导航信息
func RoleMenuAll() ([]RoleMenu, error) {
	mods := make([]RoleMenu, 0, 8)
	err := Db.Find(&mods)
	return mods, err
}

// RoleMenuPage 角色菜单导航分页信息
func RoleMenuPage(pi int, ps int, cols ...string) ([]RoleMenu, error) {
	mods := make([]RoleMenu, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Desc("Utime").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// RoleMenuCount 角色菜单导航分页信息总数
func RoleMenuCount() int {
	mod := &RoleMenu{}
	sess := Db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
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

// RoleMenuIds 返回角色菜单导航信息-ids
func RoleMenuIds(ids []int) map[int]*RoleMenu {
	mods := make([]RoleMenu, 0, len(ids))
	Db.In("id", ids).Find(&mods)
	if len(mods) > 0 {
		mapMods := make(map[int]*RoleMenu, len(mods))
		for idx := range mods {
			mapMods[mods[idx].Id] = &mods[idx]
		}
		return mapMods
	}
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
