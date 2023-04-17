package repo

import (
	"blog/internal/model"
	"errors"
)

// GrantTree 授权导航树
func GrantTree(roleId int) (map[string][]model.Grant, error) {
	mods := make([]model.Grant, 0, 8)
	sess := db.NewSession()
	defer sess.Close()
	if roleId != model.RootRoleId {
		grants := getGrantsByRole(roleId)
		sess.In("id", grants)
	}
	err := sess.Asc("group", "sort", "Id").Find(&mods)
	if err != nil {
		return nil, err
	}
	modMap := make(map[string][]model.Grant, 8)
	for idx := range mods {
		itm := mods[idx]
		if _, ok := modMap[itm.Group]; !ok {
			modMap[itm.Group] = make([]model.Grant, 0, 4)
		}
		modMap[itm.Group] = append(modMap[itm.Group], itm)
	}
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
func GrantGet(id int) (*model.Grant, bool) {
	mod := &model.Grant{}
	has, _ := db.ID(id).Get(mod)
	return mod, has
}

// GrantAll 所有授权导航信息
func GrantAll() ([]model.Grant, error) {
	mods := make([]model.Grant, 0, 8)
	err := db.Asc("group", "sort", "Id").Find(&mods)
	return mods, err
}

// GrantPage 授权导航分页信息
func GrantPage(pi int, ps int, cols ...string) ([]model.Grant, error) {
	mods := make([]model.Grant, 0, ps)
	sess := db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Asc("id").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// GrantCount 授权导航分页信息总数
func GrantCount() int {
	mod := &model.Grant{}
	sess := db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// GrantAdd 添加授权导航信息
func GrantAdd(mod *model.Grant) error {
	sess := db.NewSession()
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
func GrantEdit(mod *model.Grant, cols ...string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(mod.Id).Cols(cols...).Update(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// GrantIds 通过id集合返回授权导航信息
func GrantIds(ids []int) map[int]*model.Grant {
	mods := make([]model.Grant, 0, len(ids))
	db.In("id", ids).Find(&mods)
	mapSet := make(map[int]*model.Grant, len(mods))
	for idx := range mods {
		mapSet[mods[idx].Id] = &mods[idx]
	}
	return mapSet
}

// GrantDrop 删除单条授权导航信息
func GrantDrop(id int) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&model.Grant{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// ------------------------------------------------------ 角色授权 ------------------------------------------------------

// RoleGrantAll 通过RoleId查询所有授权信息
func RoleGrantAll(roleId int) ([]model.Grant, error) {
	mods := make([]model.Grant, 0, 8)
	// 内置根角色,拥有所有权限
	if roleId == model.RootRoleId {
		err := db.Find(&mods)
		return mods, err
	}
	grants := getGrantsByRole(roleId)
	if len(grants) < 1 {
		return mods, errors.New("No records found")
	}
	err := db.In("id", grants).Find(&mods)
	return mods, err
}
