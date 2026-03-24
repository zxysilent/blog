package repo

import (
	"blog/internal/model"
	"context"
	"reflect"
)

// GrantTree 授权导航树
func GrantTree(roleId string) (map[string][]model.Grant, error) {
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

// GrantExist 授权是否存在
func GrantExist(ctx context.Context, mod *model.Grant) bool {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	exist, _ := sess.Exist(mod)
	return exist
}

// GrantGet 授权单条数据
func GrantGet(ctx context.Context, mod *model.Grant) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if !reflect.ValueOf(mod.Id).IsZero() {
		sess.ID(mod.Id)
		sess.NoAutoCondition()
	}
	if has, _ := sess.Get(mod); !has {
		return ErrNotExist
	}
	return nil
}

// GrantList 授权列表数据
func GrantList(ctx context.Context, filter *model.GrantListArgs, cols ...string) ([]model.Grant, int, error) {
	mods := make([]model.Grant, 0, filter.GetCap())
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	sess.Asc("created")
	if filter.IsZero() {
		err := sess.Find(&mods)
		return mods, len(mods), err
	}
	sess.Limit(filter.GetLimit(), filter.GetOffset())
	count, err := sess.FindAndCount(&mods)
	return mods, int(count), err
}

// GrantCount 授权数据总数
func GrantCount(ctx context.Context, filter *model.GrantListArgs) int {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	// if filter.Mult != "" {
	// 	sess.And("field like concat('%',?,'%')", filter.Mult)
	// }
	count, _ := sess.Count(&model.Grant{})
	return int(count)
}

// GrantAdd 授权添加数据
func GrantAdd(ctx context.Context, mod *model.Grant) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if _, err := sess.InsertOne(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// GrantEdit 授权修改数据
func GrantEdit(ctx context.Context, mod *model.Grant, cols ...string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.Begin()
	if _, err := sess.ID(mod.Id).Cols(cols...).Update(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// GrantDrop 授权删除数据
func GrantDrop(ctx context.Context, mod *model.Grant) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.Begin()
	if !reflect.ValueOf(mod.Id).IsZero() {
		sess.ID(mod.Id)
		sess.NoAutoCondition()
	}
	if _, err := sess.Delete(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// GrantIds 授权id集合查询
func GrantIds(ctx context.Context, ids []string) []model.Grant {
	mods := make([]model.Grant, 0, len(ids))
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.In("id", ids).Find(&mods)
	return mods
}

// GrantMapIds 授权id集合查询
func GrantMapIds(ctx context.Context, ids []int) map[string]*model.Grant {
	mods := make([]model.Grant, 0, len(ids))
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.In("id", ids).Find(&mods)
	set := make(map[string]*model.Grant, len(mods))
	for idx := range mods {
		set[mods[idx].Id] = &mods[idx]
	}
	return set
}

// ------------------------------------------------------ 角色授权 ------------------------------------------------------

// RoleGrantAll 通过RoleId查询所有授权信息
func RoleGrantAll(ctx context.Context, roleId string) ([]model.Grant, error) {
	mods := make([]model.Grant, 0, 8)
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	// 内置根角色,拥有所有权限
	if roleId == model.RootRoleId {
		err := sess.Find(&mods)
		return mods, err
	}
	grants := getGrantsByRole(roleId)
	if len(grants) < 1 {
		return mods, ErrNotExist
	}
	err := sess.In("id", grants).Find(&mods)
	return mods, err
}
