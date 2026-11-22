package kernel

import (
	"blog/internal/model"
	"blog/internal/repo"
	"context"
	"strings"

	"github.com/labstack/echo/v4"
)

// 角色层级约定
//
//	层级数值越小权限越高，只能操作「严格大于」自己层级的角色，
//	同级角色互相不可见操作，避免平级越权。
//	根角色（model.RoleRootId）层级固定为 model.RoleRootLevel（0）且全部放行。

// currentRole 当前登录角色
// 根角色不查库：主键与层级固定，避免库中数据被改写后失效
func currentRole(ctx echo.Context) (*model.Role, error) {
	rid, ok := currentRid(ctx)
	if !ok {
		return nil, repo.ErrNotExist
	}
	if rid == model.RoleRootId {
		return &model.Role{
			Id:    model.RoleRootId,
			Name:  "根角色",
			Level: model.RoleRootLevel,
			Inner: true,
		}, nil
	}
	return repo.RoleGet(ctx.Request().Context(), rid)
}

// canView 能否查看目标角色
// 严格低于自己（数值更大）才可查看，同级与更高层级都不可见；根角色可查看全部
func canView(actor, target *model.Role) bool {
	if actor == nil || target == nil {
		return false
	}
	if actor.Id == model.RoleRootId {
		return true
	}
	return actor.Level < target.Level
}

// canManage 能否操作目标角色：层级必须严格高于目标（数值更小）
func canManage(actor, target *model.Role) bool {
	if actor == nil || target == nil {
		return false
	}
	if actor.Id == model.RoleRootId {
		return true
	}
	return actor.Level < target.Level
}

// permGroupKey 权限所属分组键，兼容历史空值
func permGroupKey(mod *model.Perm) string {
	if mod.GroupKey != "" {
		return mod.GroupKey
	}
	key, _, _ := strings.Cut(mod.Perm, ":")
	return key
}

// rolePermSet 角色拥有的权限标识集合
func rolePermSet(ctx context.Context, roleId string) map[string]struct{} {
	mods, err := repo.RolePermAll(ctx, roleId)
	if err != nil {
		return map[string]struct{}{}
	}
	set := make(map[string]struct{}, len(mods))
	for idx := range mods {
		set[mods[idx].Perm] = struct{}{}
	}
	return set
}

// hasPerm 当前登录角色是否拥有指定权限标识
// 走内存缓存，不查库；根角色由 repo.RolePermAllow 直接放行
// 用于「路由级权限之外」的细粒度判断，例如同一接口内按字段区分权限
func hasPerm(ctx echo.Context, perm string) bool {
	rid, ok := currentRid(ctx)
	if !ok {
		return false
	}
	return repo.RolePermAllow(ctx.Request().Context(), rid, perm)
}

// userRole 用户所属角色，角色缺失时返回兜底角色（层级取最大值即权限最低）
// 关键：不能用零值，因为 0 是根角色层级（权限最高）
func userRole(ctx echo.Context, roleId string) *model.Role {
	if roleId == "" {
		return &model.Role{Level: model.RoleNilLevel}
	}
	mod, err := repo.RoleGet(ctx.Request().Context(), roleId)
	if err != nil {
		return &model.Role{Id: roleId, Level: model.RoleNilLevel}
	}
	return mod
}
