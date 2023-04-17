package repo

import "blog/internal/model"

// RoleGet 单条角色信息
func RoleGet(id int) (*model.Role, bool) {
	mod := &model.Role{}
	has, _ := db.ID(id).Get(mod)
	return mod, has
}

// RoleAll 所有角色信息
func RoleAll(rid int) ([]model.Role, error) {
	mods := make([]model.Role, 0, 8)
	sess := db.NewSession()
	defer sess.Close()
	if rid > 0 {
		sess.Where("id > ?", rid)
	}
	err := sess.Find(&mods)
	return mods, err
}

// RolePage 角色分页信息
func RolePage(rid int, pi int, ps int, cols ...string) ([]model.Role, error) {
	mods := make([]model.Role, 0, ps)
	sess := db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	if rid > 0 {
		sess.Where("id > ?", rid)
	}
	err := sess.Asc("id").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// RoleCount 角色分页信息总数
func RoleCount(rid int) int {
	mod := &model.Role{}
	sess := db.NewSession()
	defer sess.Close()
	if rid > 0 {
		sess.Where("id > ?", rid)
	}
	count, _ := sess.Count(mod)
	return int(count)
}

// RoleAdd 添加角色信息
func RoleAdd(mod *model.Role) error {
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

// RoleEdit 编辑角色信息
func RoleEdit(mod *model.Role, cols ...string) error {
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

// RoleIds 通过id集合返回角色信息
func RoleIds(ids []int) map[int]*model.Role {
	mods := make([]model.Role, 0, len(ids))
	db.In("id", ids).Find(&mods)
	mapSet := make(map[int]*model.Role, len(mods))
	for idx := range mods {
		mapSet[mods[idx].Id] = &mods[idx]
	}
	return mapSet
}

// RoleDrop 删除单条角色信息
func RoleDrop(id int) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&model.Role{}); err != nil {
		sess.Rollback()
		return err
	}
	// 删除角色授权
	if _, err := sess.Where("role_id = ?", id).Delete(&model.RoleGrant{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}
