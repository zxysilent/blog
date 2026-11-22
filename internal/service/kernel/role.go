package kernel

import (
	"blog/internal/model"
	"blog/internal/repo"
	"blog/pkg/reply"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/logs"
)

// errInvalidPerm 权限id无效
var errInvalidPerm = errors.New("存在无效的权限id")

// RoleList doc
// @Tags role
// @Summary 角色列表
// @Security KeyAuth
// @Param query query model.RoleListArgs true "请求数据"
// @Success 200 {object} reply.Reply{data=[]model.Role} "返回数据"
// @Router /api/roles [get]
func RoleList(ctx echo.Context) error {
	actor, err := currentRole(ctx)
	if err != nil {
		return reply.ErrToken(ctx, "请重新登陆", "角色不存在")
	}
	args := &model.RoleListArgs{}
	if err := ctx.Bind(args); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	rctx := ctx.Request().Context()
	// 只能看到层级低于自己（数值更大）的角色；根角色可见全部，但同样看不到自己
	mods, count, err := repo.RoleList(rctx, args, actor.Level, actor.Id, actor.Id == model.RoleRootId)
	if err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.RoleList")
		return reply.Fail(ctx, "查询失败")
	}
	return reply.List(ctx, "ok", mods, count)
}

// RoleGet doc
// @Tags role
// @Summary 角色详情
// @Security KeyAuth
// @Param id path string true "角色id"
// @Success 200 {object} reply.Reply{data=model.Role} "返回数据"
// @Router /api/roles/{id} [get]
func RoleGet(ctx echo.Context) error {
	actor, err := currentRole(ctx)
	if err != nil {
		return reply.ErrToken(ctx, "请重新登陆", "角色不存在")
	}
	rctx := ctx.Request().Context()
	mod, err := repo.RoleGet(rctx, ctx.Param("id"))
	if err != nil {
		logs.Ctx(rctx).Err(err).Str("id", ctx.Param("id")).Error("repo.RoleGet")
		return reply.Fail(ctx, "角色不存在")
	}
	// 同级不可查看，只能查看层级低于自己的角色
	if !canView(actor, mod) {
		return reply.Fail(ctx, "无权查看该角色")
	}
	return reply.Succ(ctx, "ok", mod)
}

// RoleAdd doc
// @Tags role
// @Summary 角色新增
// @Security KeyAuth
// @Param body body model.RoleAddArgs true "请求数据"
// @Success 200 {object} reply.Reply{data=string} "返回数据"
// @Router /api/roles [post]
func RoleAdd(ctx echo.Context) error {
	actor, err := currentRole(ctx)
	if err != nil {
		return reply.ErrToken(ctx, "请重新登陆", "角色不存在")
	}
	ipt := &model.RoleAddArgs{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	ipt.Name = strings.TrimSpace(ipt.Name)
	ipt.Intro = strings.TrimSpace(ipt.Intro)
	if ipt.Name == "" || len(ipt.Name) > 255 {
		return reply.Fail(ctx, "角色名称不能为空且不超过 255 字符")
	}
	// 层级属于权限模型的内部细节，前端不传时由服务端决定：
	// 追加到「我可见的最深一级」之下，即权限最低的位置。
	// 这样从根角色一路新建下来会自然形成 L1、L2、L3 的层级；
	// 若统一取 actor+1，根角色建的每个角色都会落在 L1 变成同级，彼此不可见。
	rctx := ctx.Request().Context()
	level := actor.Level + 1
	if ipt.Level != nil {
		level = *ipt.Level
	} else if max := repo.RoleMaxLevel(rctx, actor.Level); max >= level {
		level = max + 1
	}
	// 层级数值越小权限越高，新角色必须比创建者权限更低，即数值更大
	if level <= actor.Level {
		return reply.Fail(ctx, "角色层级必须大于当前角色层级（数值越大权限越低）")
	}
	if repo.RoleExist(rctx, &model.Role{Name: ipt.Name}) {
		return reply.Fail(ctx, "角色名称已存在")
	}
	milli := time.Now().UnixMilli()
	mod := &model.Role{
		Name:    ipt.Name,
		Intro:   ipt.Intro,
		Level:   level,
		Updated: milli,
		Created: milli,
	}
	if err := repo.RoleAdd(rctx, mod); err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.RoleAdd")
		return reply.Fail(ctx, "添加失败")
	}
	return reply.Succ(ctx, "添加成功", mod.Id)
}

// RoleEdit doc
// @Tags role
// @Summary 角色修改
// @Security KeyAuth
// @Param id path string true "角色id"
// @Param body body model.RoleEditArgs true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/roles/{id} [put]
func RoleEdit(ctx echo.Context) error {
	actor, err := currentRole(ctx)
	if err != nil {
		return reply.ErrToken(ctx, "请重新登陆", "角色不存在")
	}
	ipt := &model.RoleEditArgs{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	id := ctx.Param("id")
	if id == model.RoleRootId {
		return reply.Fail(ctx, "根角色不可修改")
	}
	rctx := ctx.Request().Context()
	mod, err := repo.RoleGet(rctx, id)
	if err != nil {
		return reply.Fail(ctx, "角色不存在")
	}
	// 只能操作层级严格低于自己的角色，同级不可操作
	if !canManage(actor, mod) {
		return reply.Fail(ctx, "只能修改层级低于自己（数值更大）的角色")
	}
	oldName := mod.Name
	mod.Name = strings.TrimSpace(ipt.Name)
	mod.Intro = strings.TrimSpace(ipt.Intro)
	if mod.Name == "" {
		return reply.Fail(ctx, "角色名称不能为空")
	}
	// 名称变更时校验重名
	if mod.Name != oldName && repo.RoleExist(rctx, &model.Role{Name: mod.Name}) {
		return reply.Fail(ctx, "角色名称已存在")
	}
	cols := []string{"name", "intro", "updated"}
	if ipt.Level != nil {
		// 层级数值越小权限越高，被修改角色必须仍比操作者权限更低
		if *ipt.Level <= actor.Level {
			return reply.Fail(ctx, "角色层级必须大于当前角色层级（数值越大权限越低）")
		}
		mod.Level = *ipt.Level
		cols = append(cols, "level")
	}
	mod.Updated = time.Now().UnixMilli()
	if err := repo.RoleEdit(rctx, mod, cols...); err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.RoleEdit")
		return reply.Fail(ctx, "修改失败")
	}
	return reply.Succ(ctx, "修改成功")
}

// RoleDrop doc
// @Tags role
// @Summary 角色删除
// @Security KeyAuth
// @Param id path string true "角色id"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/roles/{id} [delete]
func RoleDrop(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return reply.Fail(ctx, "参数有误")
	}
	if id == model.RoleRootId {
		return reply.Fail(ctx, "根角色不可删除")
	}
	if rid, ok := currentRid(ctx); ok && rid == id {
		return reply.Fail(ctx, "不能删除自己所属角色")
	}
	actor, err := currentRole(ctx)
	if err != nil {
		return reply.ErrToken(ctx, "请重新登陆", "角色不存在")
	}
	rctx := ctx.Request().Context()
	mod, err := repo.RoleGet(rctx, id)
	if err != nil {
		return reply.Fail(ctx, "角色不存在")
	}
	if !canManage(actor, mod) {
		return reply.Fail(ctx, "只能删除层级低于自己（数值更大）的角色")
	}
	if mod.Inner {
		return reply.Fail(ctx, "内置角色不可删除")
	}
	if err := repo.RoleDrop(rctx, id); err != nil {
		logs.Ctx(rctx).Err(err).Str("id", id).Error("repo.RoleDrop")
		return reply.Fail(ctx, "删除失败")
	}
	repo.SyncRolePerms()
	return reply.Succ(ctx, "删除成功")
}

// RolePermGet doc
// @Tags role
// @Summary 获取角色已授权限
// @Security KeyAuth
// @Param id path string true "角色id"
// @Success 200 {object} reply.Reply{data=[]model.Perm} "返回数据"
// @Router /api/roles/{id}/permissions [get]
func RolePermGet(ctx echo.Context) error {
	actor, err := currentRole(ctx)
	if err != nil {
		return reply.ErrToken(ctx, "请重新登陆", "角色不存在")
	}
	rctx := ctx.Request().Context()
	id := ctx.Param("id")
	if id != model.RoleRootId {
		mod, err := repo.RoleGet(rctx, id)
		if err != nil {
			return reply.Fail(ctx, "角色不存在")
		}
		if !canView(actor, mod) {
			return reply.Fail(ctx, "无权查看该角色的权限")
		}
	}
	mods, err := repo.RolePermAll(rctx, id)
	if err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.RolePermAll")
		return reply.Fail(ctx, "查询失败")
	}
	return reply.Succ(ctx, "ok", mods)
}

// RolePermEdit doc
// @Tags role
// @Summary 角色授权（全量替换）
// @Security KeyAuth
// @Param id path string true "角色id"
// @Param body body model.RolePermArgs true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/roles/{id}/permissions [put]
func RolePermEdit(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return reply.Fail(ctx, "参数有误")
	}
	if id == model.RoleRootId {
		return reply.Fail(ctx, "根角色拥有所有权限，无需授权")
	}
	actor, err := currentRole(ctx)
	if err != nil {
		return reply.ErrToken(ctx, "请重新登陆", "角色不存在")
	}
	ipt := &model.RolePermArgs{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	rctx := ctx.Request().Context()
	target, err := repo.RoleGet(rctx, id)
	if err != nil {
		return reply.Fail(ctx, "角色不存在")
	}
	// 只能给层级低于自己（数值更大）的角色授权
	if !canManage(actor, target) {
		return reply.Fail(ctx, "只能给层级低于自己的角色授权")
	}
	permIds := normalizePermIds(ipt.PermIds)
	if len(permIds) < 1 {
		// 清空授权：需逐项校验被移除的权限
		if err := roleGrantCheck(rctx, actor, id, nil); err != nil {
			return reply.Fail(ctx, err.Error())
		}
	} else {
		if exist := repo.PermMapIds(rctx, permIds); len(exist) != len(permIds) {
			return reply.Fail(ctx, errInvalidPerm.Error())
		}
		if err := roleGrantCheck(rctx, actor, id, permIds); err != nil {
			return reply.Fail(ctx, err.Error())
		}
	}
	if err := repo.RolePermSync(rctx, id, permIds); err != nil {
		logs.Ctx(rctx).Err(err).Str("roleId", id).Error("repo.RolePermSync")
		return reply.Fail(ctx, "授权失败")
	}
	return reply.Succ(ctx, "授权成功")
}

// normalizePermIds 权限id去重并剔除空值
func normalizePermIds(ids []string) []string {
	out := make([]string, 0, len(ids))
	seen := make(map[string]struct{}, len(ids))
	for _, id := range ids {
		if id == "" {
			continue
		}
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		out = append(out, id)
	}
	return out
}

// roleGrantCheck 授权变更越权校验
// 变更集合 = 新增 + 移除，逐项要求操作者：
//  1. 自己已拥有该权限（不能授予自己没有的权限，也不能剥夺自己没有的权限）
//  2. 拥有该权限所属分组的授权权限 <分组键>:grant
//
// 根角色由调用方提前放行
func roleGrantCheck(rctx context.Context, actor *model.Role, roleId string, newIds []string) error {
	if actor.Id == model.RoleRootId {
		return nil
	}
	oldIds := repo.RolePermIds(rctx, roleId)
	newSet := make(map[string]struct{}, len(newIds))
	for _, id := range newIds {
		newSet[id] = struct{}{}
	}
	oldSet := make(map[string]struct{}, len(oldIds))
	for _, id := range oldIds {
		oldSet[id] = struct{}{}
	}
	delta := make([]string, 0, len(oldIds)+len(newIds))
	for _, id := range newIds {
		if _, ok := oldSet[id]; !ok {
			delta = append(delta, id)
		}
	}
	for _, id := range oldIds {
		if _, ok := newSet[id]; !ok {
			delta = append(delta, id)
		}
	}
	if len(delta) < 1 {
		return nil
	}
	infos := repo.PermMapIds(rctx, delta)
	if len(infos) != len(delta) {
		return errInvalidPerm
	}
	owned := rolePermSet(rctx, actor.Id)
	for _, id := range delta {
		mod, ok := infos[id]
		if !ok {
			return errInvalidPerm
		}
		if _, ok := owned[mod.Perm]; !ok {
			return fmt.Errorf("不能授予自己没有的权限：%s", mod.Name)
		}
		if _, ok := owned[model.GroupGrantPerm(permGroupKey(mod))]; !ok {
			return fmt.Errorf("缺少「%s」分组的授权权限", mod.Group)
		}
	}
	return nil
}
