package repo

import (
	"blog/internal/model"
	"context"
	"sync"
)

// RolePermIds 查询角色已授权的权限id集合
func RolePermIds(ctx context.Context, roleId string) []string {
	ids := make([]string, 0, 16)
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.Table("role_permission").Cols("perm_id").Where("role_id = ?", roleId).Find(&ids)
	return ids
}

// RolePermAll 查询角色已授权的权限详情
// 根角色直接返回全部权限
func RolePermAll(ctx context.Context, roleId string) ([]model.Perm, error) {
	if roleId == model.RoleRootId {
		return PermAll(ctx)
	}
	ids := RolePermIds(ctx, roleId)
	if len(ids) < 1 {
		return []model.Perm{}, nil
	}
	mods := make([]model.Perm, 0, len(ids))
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	err := sess.In("id", ids).Asc("`group`").Asc("sort").Asc("id").Find(&mods)
	return mods, err
}

// rolePermPrune 清理指向已不存在权限点的角色授权
// 权限点主键由 perm 稳定派生，正常情况下不会产生孤儿行，
// 仅在从旧版本（随机主键）升级时可能出现，这里顺带清一次。
func rolePermPrune(ctx context.Context) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	total, err := sess.Count(&model.Perm{})
	if err != nil || total < 1 {
		return err
	}
	_, err = sess.Exec("DELETE FROM role_permission WHERE perm_id NOT IN (SELECT id FROM permission)")
	return err
}

// RolePermSync 全量替换角色授权
func RolePermSync(ctx context.Context, roleId string, permIds []string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if err := sess.Begin(); err != nil {
		return err
	}
	// 只删除当前角色的绑定，避免误删其它角色
	if _, err := sess.Where("role_id = ?", roleId).Delete(&model.RolePerm{}); err != nil {
		sess.Rollback()
		return err
	}
	if len(permIds) > 0 {
		seen := make(map[string]struct{}, len(permIds))
		mods := make([]model.RolePerm, 0, len(permIds))
		for _, gid := range permIds {
			if gid == "" {
				continue
			}
			if _, has := seen[gid]; has {
				continue
			}
			seen[gid] = struct{}{}
			mods = append(mods, model.RolePerm{RoleId: roleId, PermId: gid})
		}
		for idx := range mods {
			mods[idx].BeforeInsert()
		}
		if _, err := sess.Insert(&mods); err != nil {
			sess.Rollback()
			return err
		}
	}
	if err := sess.Commit(); err != nil {
		return err
	}
	SyncRolePerms()
	return nil
}

var (
	_rolePermMu sync.RWMutex
	// _rolePerms 角色id -> 权限标识集合（perm 字符串）
	_rolePerms map[string]map[string]struct{}
)

// SyncRolePerms 重建角色权限缓存
// 在角色授权变更、权限点变更、启动初始化后调用
func SyncRolePerms() {
	all, err := PermAll(context.Background())
	if err != nil {
		return
	}
	// 权限id -> 权限标识
	permById := make(map[string]string, len(all))
	for idx := range all {
		if all[idx].Perm == "" {
			continue
		}
		permById[all[idx].Id] = all[idx].Perm
	}
	// 查出所有角色授权关系
	type rolePermRow struct {
		RoleId string `xorm:"role_id"`
		PermId string `xorm:"perm_id"`
	}
	rows := make([]rolePermRow, 0, 64)
	db.Table("role_permission").Cols("role_id", "perm_id").Find(&rows)

	set := make(map[string]map[string]struct{}, 16)
	for _, row := range rows {
		perm, ok := permById[row.PermId]
		if !ok {
			continue
		}
		perms, ok := set[row.RoleId]
		if !ok {
			perms = make(map[string]struct{}, 16)
			set[row.RoleId] = perms
		}
		perms[perm] = struct{}{}
	}
	_rolePermMu.Lock()
	_rolePerms = set
	_rolePermMu.Unlock()
}

// RolePermAllow 判断角色是否拥有指定权限标识
// 根角色直接放行；perm 为空表示该接口无需授权
func RolePermAllow(ctx context.Context, roleId string, perm string) bool {
	if perm == "" || roleId == model.RoleRootId {
		return true
	}
	_rolePermMu.RLock()
	perms := _rolePerms[roleId]
	_rolePermMu.RUnlock()
	if len(perms) < 1 {
		return false
	}
	_, ok := perms[perm]
	return ok
}
