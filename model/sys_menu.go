package model

import (
	"time"
)

// Menu 菜单导航
type Menu struct {
	Id       int       `xorm:"INT(11) PK AUTOINCR comment('主键')" json:"id"`
	Pid      int       `xorm:"INT(11) DEFAULT 0 comment('父id')" json:"pid"`
	Title    string    `xorm:"VARCHAR(255) comment('菜单')" json:"title"`
	Name     string    `xorm:"VARCHAR(255) comment('名称')" json:"name"`
	Path     string    `xorm:"VARCHAR(255) comment('路径')" json:"path"`
	Intro    string    `xorm:"VARCHAR(255) comment('介绍')" json:"intro"`
	Icon     string    `xorm:"VARCHAR(255) comment('菜单图标')" json:"icon"`
	Show     bool      `xorm:"TINYINT(4) DEFAULT 1 comment('导航显示')" json:"show"`
	Comp     string    `xorm:"VARCHAR(255) comment('vue文件路径Component')" json:"comp"`
	Sort     int       `xorm:"INT(11) DEFAULT 1000 comment('排序')" json:"sort"`
	Inner    bool      `xorm:"TINYINT(4) DEFAULT 0 comment('内部禁止删除')" json:"inner"`
	Ctime    time.Time `xorm:"DATETIME comment('时间')" json:"ctime"`
	Children []Menu    `xorm:"-" json:"children"`
}

func (Menu) TableName() string {
	return "sys_menu"
}

type Meta struct {
	Title string `json:"title"`
	Show  bool   `json:"show"`
	Icon  string `json:"icon"`
}

// interface RouteConfig = {
// 	path: string,
// 	component?: Component,
// 	name?: string, // 命名路由
// 	components?: { [name: string]: Component }, // 命名视图组件
// 	redirect?: string | Location | Function,
// 	props?: boolean | Object | Function,
// 	alias?: string | Array<string>,
// 	children?: Array<RouteConfig>, // 嵌套路由
// 	beforeEnter?: (to: Route, from: Route, next: Function) => void,
// 	meta?: any,

// 	// 2.6.0+
// 	caseSensitive?: boolean, // 匹配规则是否大小写敏感？(默认值：false)
// 	pathToRegexpOptions?: Object // 编译正则的选项
//}

// MenuTree 菜单导航树
func MenuTree() ([]Menu, error) {
	mods := make([]Menu, 0, 8)
	err := Db.Asc("Pid", "Sort", "Id").Find(&mods)
	if err != nil {
		return nil, err
	}
	modMap := make(map[int][]Menu, 8)
	for idx := range mods {
		itm := mods[idx]
		if _, ok := modMap[itm.Pid]; !ok {
			modMap[itm.Pid] = make([]Menu, 0, 4)
		}
		modMap[itm.Pid] = append(modMap[itm.Pid], itm)
	}
	menus := make([]Menu, 0, 10)
	iters := modMap[0]
	for _, menu := range iters {
		mod := menu
		if children, has := modMap[mod.Id]; has {
			mod.Children = children
		}
		menus = append(menus, mod)
	}
	return menus, err
}

// MenuGet 单条菜单导航信息
func MenuGet(id int) (*Menu, bool) {
	mod := &Menu{}
	has, _ := Db.ID(id).Get(mod)
	return mod, has
}

// MenuAll 所有菜单导航信息
func MenuAll() ([]Menu, error) {
	mods := make([]Menu, 0, 8)
	err := Db.Find(&mods)
	return mods, err
}

// MenuPage 菜单导航分页信息
func MenuPage(pi int, ps int, cols ...string) ([]Menu, error) {
	mods := make([]Menu, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Desc("Utime").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// MenuCount 菜单导航分页信息总数
func MenuCount() int {
	mod := &Menu{}
	sess := Db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// MenuAdd 添加菜单导航信息
func MenuAdd(mod *Menu) error {
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

// MenuEdit 编辑菜单导航信息
func MenuEdit(mod *Menu, cols ...string) error {
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

// MenuIds 返回菜单导航信息-ids
func MenuIds(ids []int) map[int]*Menu {
	mods := make([]Menu, 0, len(ids))
	Db.In("id", ids).Find(&mods)
	if len(mods) > 0 {
		mapMods := make(map[int]*Menu, len(mods))
		for idx := range mods {
			mapMods[mods[idx].Id] = &mods[idx]
		}
		return mapMods
	}
	return nil
}

// MenuDrop 删除单条菜单导航信息
func MenuDrop(id int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&Menu{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}
