package repo

import (
	"blog/internal/model"
	"context"
	"strings"

	"xorm.io/xorm"
)

// UserNum 按账号查询用户
func UserNum(ctx context.Context, num string) (*model.User, error) {
	mod := &model.User{}
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	has, err := sess.Where("num = ?", num).Get(mod)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotExist
	}
	return mod, nil
}

// UserGet 用户详情（附带角色）
func UserGet(ctx context.Context, id string) (*model.User, error) {
	mod, err := GetAny[model.User](ctx, id)
	if err != nil {
		return nil, err
	}
	fillUserRole(ctx, &mod)
	return &mod, nil
}

// UserExistNum 账号是否已存在
func UserExistNum(ctx context.Context, num string) bool {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	exist, _ := sess.Where("num = ?", num).Exist(&model.User{})
	return exist
}

// UserOpenid 按 openid 查用户
// 供公众号/小程序登录时把微信身份映射到本系统账号
func UserOpenid(ctx context.Context, openid string) (*model.User, error) {
	openid = model.UserNormalizeOpenid(openid)
	if openid == "" {
		return nil, ErrNotExist
	}
	mod := &model.User{}
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	has, err := sess.Where("openid = ?", openid).Get(mod)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotExist
	}
	fillUserRole(ctx, mod)
	return mod, nil
}

// UserExistOpenid openid 是否已被占用
// 空串表示未绑定，不参与唯一性判定；excludeId 用于编辑时排除自身
func UserExistOpenid(ctx context.Context, openid, excludeId string) bool {
	openid = model.UserNormalizeOpenid(openid)
	if openid == "" {
		return false
	}
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.Where("openid = ?", openid)
	if excludeId != "" {
		sess.And("id <> ?", excludeId)
	}
	exist, _ := sess.Exist(&model.User{})
	return exist
}

// UserCount 用户计数
func UserCount(ctx context.Context, args *model.UserListArgs) int {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	userCondition(sess, args)
	count, _ := sess.Count(&model.User{})
	return int(count)
}

// UserList 用户分页列表（附带角色）
// actorLevel / selfId 由服务层注入：
//   - 默认只能看到角色层级低于自己（数值 > actorLevel）的用户
//   - all 为真时（根角色）不做层级限制
//   - 无论哪种情况都不返回自己：自己的资料走 /api/auth，不在用户管理里回显
func UserList(ctx context.Context, args *model.UserListArgs, actorLevel int, selfId string, all bool, cols ...string) ([]model.User, int, error) {
	mods := make([]model.User, 0, args.GetCap())
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	userCondition(sess, args)
	if !all {
		// 层级数值越小权限越高，因此「比自己低」即 level 严格大于自身
		sess.And("role_id IN (SELECT id FROM role WHERE level > ?)", actorLevel)
	}
	if selfId != "" {
		sess.And("id <> ?", selfId)
	}
	sess.Desc("id")
	if args.IsZero() {
		err := sess.Find(&mods)
		fillUserRoles(ctx, mods)
		return mods, len(mods), err
	}
	count, err := sess.Limit(args.GetLimit(), args.GetOffset()).FindAndCount(&mods)
	fillUserRoles(ctx, mods)
	return mods, int(count), err
}

// UserAdd 新增用户
func UserAdd(ctx context.Context, mod *model.User) error {
	return AddAny(ctx, mod)
}

// UserEdit 修改用户
func UserEdit(ctx context.Context, mod *model.User, cols ...string) error {
	return EditAny(ctx, mod, cols...)
}

// UserDrop 删除用户
func UserDrop(ctx context.Context, id string) error {
	return DelAny[model.User](ctx, id)
}

// UserMapIds 用户id集合查询
func UserMapIds(ctx context.Context, ids []string) map[string]*model.User {
	set := make(map[string]*model.User, len(ids))
	if len(ids) < 1 {
		return set
	}
	mods := make([]model.User, 0, len(ids))
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.In("id", ids).Find(&mods)
	for idx := range mods {
		set[mods[idx].Id] = &mods[idx]
	}
	return set
}

// userCondition 组装用户查询条件
func userCondition(sess *xorm.Session, args *model.UserListArgs) {
	if args.Mult != "" {
		sess.And("(num like ? or name like ? or phone like ? or email like ?)",
			args.SafeMult(), args.SafeMult(), args.SafeMult(), args.SafeMult())
	}
	if args.Num != "" {
		sess.And("num = ?", strings.TrimSpace(args.Num))
	}
	if args.Openid != "" {
		sess.And("openid = ?", model.UserNormalizeOpenid(args.Openid))
	}
	if args.RoleId != "" {
		sess.And("role_id = ?", args.RoleId)
	}
	if args.Status > 0 {
		sess.And("status = ?", args.Status)
	}
}

// fillUserRole 填充单个用户的角色
func fillUserRole(ctx context.Context, mod *model.User) {
	if mod == nil || mod.RoleId == "" {
		return
	}
	role, err := RoleGet(ctx, mod.RoleId)
	if err == nil {
		mod.Role = role
	}
}

// fillUserRoles 批量填充用户角色，避免 N+1
func fillUserRoles(ctx context.Context, mods []model.User) {
	if len(mods) < 1 {
		return
	}
	rids := make([]string, 0, len(mods))
	for idx := range mods {
		if mods[idx].RoleId != "" {
			rids = append(rids, mods[idx].RoleId)
		}
	}
	if len(rids) < 1 {
		return
	}
	rmp := RoleMapIds(ctx, rids)
	for idx := range mods {
		if role, ok := rmp[mods[idx].RoleId]; ok {
			mods[idx].Role = role
		}
	}
}
