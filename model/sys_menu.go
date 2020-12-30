package model

import "time"

// SysMenu 菜单导航
type SysMenu struct {
	Id    int       `xorm:"pk autoincr INT(11)"`   //主键
	Pid   int       `xorm:"default 0 INT(11)"`     //父id
	Path  string    `xorm:"VARCHAR(255)"`          //路径
	Intro string    `xorm:"VARCHAR(255)"`          //介绍
	Sort  int       `xorm:"default 999 INT(11)"`   //排序
	Ctime time.Time `xorm:"DATETIME" json:"ctime"` //时间
}

// SysMenuGet 单条菜单导航信息
func SysMenuGet(id int) (*SysMenu, bool) {
	mod := &SysMenu{}
	has, _ := Db.ID(id).Get(mod)
	return mod, has
}

// SysMenuAll 所有菜单导航信息
func SysMenuAll() ([]SysMenu, error) {
	mods := make([]SysMenu, 0, 8)
	err := Db.Find(&mods)
	return mods, err
}

// SysMenuPage 菜单导航分页信息
func SysMenuPage(pi int, ps int, cols ...string) ([]SysMenu, error) {
	mods := make([]SysMenu, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Desc("Utime").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// SysMenuCount 菜单导航分页信息总数
func SysMenuCount() int {
	mod := &SysMenu{}
	sess := Db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// SysMenuAdd 添加菜单导航信息
func SysMenuAdd(mod *SysMenu) error {
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

// SysMenuEdit 编辑菜单导航信息
func SysMenuEdit(mod *SysMenu, cols ...string) error {
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

// SysMenuIds 返回菜单导航信息-ids
func SysMenuIds(ids []int) map[int]*SysMenu {
	mods := make([]SysMenu, 0, len(ids))
	Db.In("id", ids).Find(&mods)
	if len(mods) > 0 {
		mapMods := make(map[int]*SysMenu, len(mods))
		for idx := range mods {
			mapMods[mods[idx].Id] = &mods[idx]
		}
		return mapMods
	}
	return nil
}

// SysMenuDrop 删除单条菜单导航信息
func SysMenuDrop(id int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&SysMenu{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}
