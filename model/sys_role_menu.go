package model

import "time"

// SysRoleMenu 角色菜单导航
type SysRoleMenu struct {
	Id       int           `xorm:"pk autoincr INT(11)"` //主键
	RoleId   int           `xorm:"default 0 INT(11)"`   //RoleId
	MenuId   int           `xorm:"default 0 INT(11)"`   //MenuId
	Add      bool          `xorm:"BOOL default false"`  //添加
	Edit     bool          `xorm:"BOOL default false"`  //编辑
	Drop     bool          `xorm:"BOOL default false"`  //删除
	Fetch    bool          `xorm:"BOOL default false"`  //查询
	Export   bool          `xorm:"BOOL default false"`  //导出
	Import   bool          `xorm:"BOOL default false"`  //导入
	Menu     SysMenu       `xorm:"-"`
	Children []SysRoleMenu `xorm:"-"` //子菜单导航
	Ctime    time.Time     `xorm:"DATETIME" json:"ctime"`
}

// SysRoleMenuGet 单条角色菜单导航信息
func SysRoleMenuGet(id int) (*SysRoleMenu, bool) {
	mod := &SysRoleMenu{}
	has, _ := Db.ID(id).Get(mod)
	return mod, has
}

// SysRoleMenuAll 所有角色菜单导航信息
func SysRoleMenuAll() ([]SysRoleMenu, error) {
	mods := make([]SysRoleMenu, 0, 8)
	err := Db.Find(&mods)
	return mods, err
}

// SysRoleMenuPage 角色菜单导航分页信息
func SysRoleMenuPage(pi int, ps int, cols ...string) ([]SysRoleMenu, error) {
	mods := make([]SysRoleMenu, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Desc("Utime").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// SysRoleMenuCount 角色菜单导航分页信息总数
func SysRoleMenuCount() int {
	mod := &SysRoleMenu{}
	sess := Db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// SysRoleMenuAdd 添加角色菜单导航信息
func SysRoleMenuAdd(mod *SysRoleMenu) error {
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

// SysRoleMenuEdit 编辑角色菜单导航信息
func SysRoleMenuEdit(mod *SysRoleMenu, cols ...string) error {
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

// SysRoleMenuIds 返回角色菜单导航信息-ids
func SysRoleMenuIds(ids []int) map[int]*SysRoleMenu {
	mods := make([]SysRoleMenu, 0, len(ids))
	Db.In("id", ids).Find(&mods)
	if len(mods) > 0 {
		mapMods := make(map[int]*SysRoleMenu, len(mods))
		for idx := range mods {
			mapMods[mods[idx].Id] = &mods[idx]
		}
		return mapMods
	}
	return nil
}

// SysRoleMenuDrop 删除单条角色菜单导航信息
func SysRoleMenuDrop(id int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&SysRoleMenu{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}
