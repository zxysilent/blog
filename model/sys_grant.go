package model

import (
	"time"
)

// Grant 授权
type Grant struct {
	Id       int       `xorm:"INT(11) PK AUTOINCR comment('主键')" json:"id"`
	Guid     string    `xorm:"VARCHAR(255) comment('标识')" json:"guid"` //全局唯一标识符（GUID，Globally Unique Identifier）
	Name     string    `xorm:"VARCHAR(255) comment('名称')" json:"name"`
	Pid      int       `xorm:"INT(11) DEFAULT 0 comment('父id')" json:"pid"`
	Kind     int       `xorm:"INT(11) DEFAULT 0 comment('种类')" json:"kind"`
	Title    string    `xorm:"VARCHAR(255) comment('菜单')" json:"title"`
	Path     string    `xorm:"VARCHAR(255) comment('路径')" json:"path"`
	Use      bool      `xorm:"TINYINT(4) DEFAULT 1 comment('是否使用')" json:"use"`
	Icon     string    `xorm:"VARCHAR(255) comment('菜单图标')" json:"icon"`
	Show     bool      `xorm:"TINYINT(4) DEFAULT 1 comment('导航显示')" json:"show"`
	Comp     string    `xorm:"VARCHAR(255) comment('vue文件路径Component')" json:"comp"`
	Sort     int       `xorm:"INT(11) DEFAULT 1000 comment('排序')" json:"sort"`
	Inner    bool      `xorm:"TINYINT(4) DEFAULT 0 comment('内部禁止删除')" json:"inner"`
	Ctime    time.Time `xorm:"DATETIME comment('时间')" json:"ctime"`
	Children []Menu    `xorm:"-" swaggerignore:"true" json:"children"` //忽略文档生成
}

// 根角色
const RootRoleId = 1

func (Grant) TableName() string {
	return "sys_grant"
}

// GrantTree 授权导航树
func GrantTree(roleId int) (map[string][]Grant, error) {
	mods := make([]Grant, 0, 8)
	sess := Db.NewSession()
	defer sess.Close()
	if roleId != RootRoleId {
		grants := getGrantsByRole(roleId)
		sess.In("id", grants)
	}
	err := sess.Asc("sort", "Id").Find(&mods)
	if err != nil {
		return nil, err
	}
	modMap := make(map[string][]Grant, 8)
	// for idx := range mods {
	// 	itm := mods[idx]
	// 	if _, ok := modMap[itm.Group]; !ok {
	// 		modMap[itm.Group] = make([]Grant, 0, 4)
	// 	}
	// 	modMap[itm.Group] = append(modMap[itm.Group], itm)
	// }
	// menus := make([]Menu, 0, 10)
	// iters := modMap[0]
	// for _, menu := range iters {
	// 	mod := menu
	// 	if children, has := modMap[mod.Group]; has {
	// 		mod.Children = children
	// 	}
	// 	menus = append(menus, mod)
	// }
	// return menus, err
	return modMap, nil
}

// GrantGet 单条授权导航信息
func GrantGet(id int) (*Grant, bool) {
	mod := &Grant{}
	has, _ := Db.ID(id).Get(mod)
	return mod, has
}

// GrantAll 所有授权导航信息
func GrantAll() ([]Grant, error) {
	mods := make([]Grant, 0, 8)
	err := Db.Asc("sort", "Id").Find(&mods)
	return mods, err
}

// GrantPage 授权导航分页信息
func GrantPage(pi int, ps int, cols ...string) ([]Grant, error) {
	mods := make([]Grant, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Desc("Ctime").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// GrantCount 授权导航分页信息总数
func GrantCount() int {
	mod := &Grant{}
	sess := Db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// GrantAdd 添加授权导航信息
func GrantAdd(mod *Grant) error {
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

// GrantEdit 编辑授权导航信息
func GrantEdit(mod *Grant, cols ...string) error {
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

// GrantMapIds 通过id集合返回授权导航信息
func GrantMapIds(ids []int) map[int]*Grant {
	mods := make([]Grant, 0, len(ids))
	Db.In("id", ids).Find(&mods)
	mapSet := make(map[int]*Grant, len(mods))
	for idx := range mods {
		mapSet[mods[idx].Id] = &mods[idx]
	}
	return mapSet
}

// GrantDrop 删除单条授权导航信息
func GrantDrop(id int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&Grant{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// ------------------------------------------------------ 角色授权 ------------------------------------------------------

// RoleGrantAll 通过RoleId查询所有授权信息
func RoleGrantAll(roleId int) ([]Grant, error) {
	mods := make([]Grant, 0, 8)
	// 内置根角色,拥有所有权限
	if roleId == RootRoleId {
		err := Db.Find(&mods)
		return mods, err
	}
	grants := getGrantsByRole(roleId)
	err := Db.In("id", grants).Find(&mods)
	return mods, err
}
