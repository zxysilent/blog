package repo

import (
	"blog/internal/model"
	"context"
	"reflect"
)

// RoleExist 角色是否存在
func RoleExist(ctx context.Context, mod *model.Role) bool {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	exist, _ := sess.Exist(mod)
	return exist
}

// RoleGet 角色单条数据
func RoleGet(ctx context.Context, mod *model.Role) error {
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

// RoleList 角色列表数据
func RoleList(ctx context.Context, filter *model.RoleListArgs, cols ...string) ([]model.Role, int, error) {
	mods := make([]model.Role, 0, filter.GetCap())
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

// RoleCount 角色数据总数
func RoleCount(ctx context.Context, filter *model.RoleListArgs) int {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	// if filter.Mult != "" {
	// 	sess.And("field like concat('%',?,'%')", filter.Mult)
	// }
	count, _ := sess.Count(&model.Role{})
	return int(count)
}

// RoleAdd 角色添加数据
func RoleAdd(ctx context.Context, mod *model.Role) error {
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

// RoleEdit 角色修改数据
func RoleEdit(ctx context.Context, mod *model.Role, cols ...string) error {
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

// RoleDrop 角色删除数据
func RoleDrop(ctx context.Context, mod *model.Role) error {
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

// RoleIds 角色id集合查询
func RoleIds(ctx context.Context, ids []string) []model.Role {
	mods := make([]model.Role, 0, len(ids))
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.In("id", ids).Find(&mods)
	return mods
}

// RoleMapIds 角色id集合查询
func RoleMapIds(ctx context.Context, ids []int) map[string]*model.Role {
	mods := make([]model.Role, 0, len(ids))
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.In("id", ids).Find(&mods)
	set := make(map[string]*model.Role, len(mods))
	for idx := range mods {
		set[mods[idx].Id] = &mods[idx]
	}
	return set
}
