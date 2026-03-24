package repo

import (
	"blog/internal/model"
	"context"
	"reflect"
)

// UserExist 用户是否存在
func UserExist(ctx context.Context, mod *model.User) bool {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	exist, _ := sess.Exist(mod)
	return exist
}

// UserGet 用户单条数据
func UserGet(ctx context.Context, mod *model.User) error {
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

// UserList 用户列表数据
func UserList(ctx context.Context, filter *model.UserListArgs, cols ...string) ([]model.User, int, error) {
	mods := make([]model.User, 0, filter.GetCap())
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

// UserCount 用户数据总数
func UserCount(ctx context.Context, filter *model.UserListArgs) int {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	// if filter.Mult != "" {
	// 	sess.And("field like concat('%',?,'%')", filter.Mult)
	// }
	count, _ := sess.Count(&model.User{})
	return int(count)
}

// UserAdd 用户添加数据
func UserAdd(ctx context.Context, mod *model.User) error {
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

// UserEdit 用户修改数据
func UserEdit(ctx context.Context, mod *model.User, cols ...string) error {
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

// UserDrop 用户删除数据
func UserDrop(ctx context.Context, mod *model.User) error {
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

// UserIds 用户id集合查询
func UserIds(ctx context.Context, ids []string) []model.User {
	mods := make([]model.User, 0, len(ids))
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.In("id", ids).Find(&mods)
	return mods
}

// UserMapIds 用户id集合查询
func UserMapIds(ctx context.Context, ids []string) map[string]*model.User {
	mods := make([]model.User, 0, len(ids))
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.In("id", ids).Find(&mods)
	set := make(map[string]*model.User, len(mods))
	for idx := range mods {
		set[mods[idx].Id] = &mods[idx]
	}
	return set
}
