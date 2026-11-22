package repo

import (
	"blog/internal/model"
	"context"
)

// RoleExist 角色是否已存在
func RoleExist(ctx context.Context, mod *model.Role) bool {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	exist, _ := sess.Exist(mod)
	return exist
}

// RoleGet 角色详情
func RoleGet(ctx context.Context, id string) (*model.Role, error) {
	mod := &model.Role{}
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	has, err := sess.ID(id).NoAutoCondition().Get(mod)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotExist
	}
	return mod, nil
}

// RoleList 角色分页列表
// actorLevel / selfId 由服务层按当前登录角色注入：
//   - 默认只能看到层级低于自己的角色（数值 > actorLevel）
//   - all 为真时（根角色）不做层级限制
//   - 无论哪种情况都不返回自己所在的角色：自己拥有什么权限，不应在管理列表里回显
func RoleList(ctx context.Context, args *model.RoleListArgs, actorLevel int, selfId string, all bool, cols ...string) ([]model.Role, int, error) {
	mods := make([]model.Role, 0, args.GetCap())
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	// 层级数值越小权限越高，因此「比自己低」即 level 严格大于自身
	if !all {
		sess.And("level > ?", actorLevel)
	}
	if selfId != "" {
		sess.And("id <> ?", selfId)
	}
	if args.Mult != "" {
		sess.And("(name like ? or intro like ?)", args.SafeMult(), args.SafeMult())
	}
	sess.Asc("level").Asc("id")
	if args.IsZero() {
		err := sess.Find(&mods)
		fillRoleCounts(ctx, mods)
		return mods, len(mods), err
	}
	count, err := sess.Limit(args.GetLimit(), args.GetOffset()).FindAndCount(&mods)
	fillRoleCounts(ctx, mods)
	return mods, int(count), err
}

// fillRoleCounts 填充角色的关联用户数与已授权权限数
// 用分组聚合一次查出全部角色的计数，避免逐个角色查库
func fillRoleCounts(ctx context.Context, mods []model.Role) {
	if len(mods) < 1 {
		return
	}
	ids := make([]string, 0, len(mods))
	hasRoot := false
	for _, it := range mods {
		ids = append(ids, it.Id)
		if it.Id == model.RoleRootId {
			hasRoot = true
		}
	}
	userCount := make(map[string]int, len(mods))
	permCount := make(map[string]int, len(mods))

	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)

	rows := make([]struct {
		RoleId string `xorm:"role_id"`
		Total  int    `xorm:"total"`
	}, 0, len(mods))
	sess.Table("user").Select("role_id, count(*) as total").In("role_id", ids).GroupBy("role_id").Find(&rows)
	for _, it := range rows {
		userCount[it.RoleId] = it.Total
	}

	rows = rows[:0]
	sess.Table("role_permission").Select("role_id, count(*) as total").In("role_id", ids).GroupBy("role_id").Find(&rows)
	for _, it := range rows {
		permCount[it.RoleId] = it.Total
	}

	// 根角色天然拥有全部权限，授权表里没有它的行，单独按权限总数兜底
	rootTotal := 0
	if hasRoot {
		total, _ := sess.Count(&model.Perm{})
		rootTotal = int(total)
	}
	for idx := range mods {
		mods[idx].UserCount = userCount[mods[idx].Id]
		if mods[idx].Id == model.RoleRootId {
			mods[idx].PermCount = int(rootTotal)
			continue
		}
		mods[idx].PermCount = permCount[mods[idx].Id]
	}
}

// RoleMaxLevel 返回层级低于 atLevel（数值更大）的角色中的最大层级
// 没有更低层级的角色时返回 0，调用方据此取 atLevel+1
// 用于给新角色自动分配层级：追加到可见范围的最深处
func RoleMaxLevel(ctx context.Context, atLevel int) int {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	row := struct {
		Level int `xorm:"level"`
	}{}
	has, err := sess.Table("role").Where("level > ?", atLevel).
		Select("MAX(level) AS level").Get(&row)
	if err != nil || !has {
		return 0
	}
	return row.Level
}

// RoleAll 全部角色
func RoleAll(ctx context.Context) ([]model.Role, error) {
	mods := make([]model.Role, 0, 8)
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	err := sess.Asc("level").Asc("id").Find(&mods)
	return mods, err
}

// RoleAdd 新增角色
func RoleAdd(ctx context.Context, mod *model.Role) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if _, err := sess.InsertOne(mod); err != nil {
		return err
	}
	return nil
}

// RoleEdit 修改角色
func RoleEdit(ctx context.Context, mod *model.Role, cols ...string) error {
	return EditAny(ctx, mod, cols...)
}

// RoleDrop 删除角色（同时清理角色授权）
func RoleDrop(ctx context.Context, id string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if err := sess.Begin(); err != nil {
		return err
	}
	if _, err := sess.ID(id).NoAutoCondition().Delete(&model.Role{}); err != nil {
		sess.Rollback()
		return err
	}
	if _, err := sess.Where("role_id = ?", id).Delete(&model.RolePerm{}); err != nil {
		sess.Rollback()
		return err
	}
	return sess.Commit()
}

// RoleMapIds 角色id集合查询
func RoleMapIds(ctx context.Context, ids []string) map[string]*model.Role {
	set := make(map[string]*model.Role, len(ids))
	if len(ids) < 1 {
		return set
	}
	mods := make([]model.Role, 0, len(ids))
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.In("id", ids).Find(&mods)
	for idx := range mods {
		set[mods[idx].Id] = &mods[idx]
	}
	return set
}
