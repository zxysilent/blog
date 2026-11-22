package kernel

import (
	"blog/internal/model"
	"blog/internal/repo"
	"blog/pkg/reply"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/logs"
)

// PermList doc
// @Tags permission
// @Summary 权限列表
// @Security KeyAuth
// @Param query query model.PermListArgs true "请求数据"
// @Success 200 {object} reply.Reply{data=[]model.Perm} "返回数据"
// @Router /api/permissions [get]
func PermList(ctx echo.Context) error {
	args := &model.PermListArgs{}
	if err := ctx.Bind(args); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	rctx := ctx.Request().Context()
	mods, count, err := repo.PermList(rctx, args)
	if err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.PermList")
		return reply.Fail(ctx, "查询失败")
	}
	return reply.List(ctx, "ok", mods, count)
}

// PermTree doc
// @Tags permission
// @Summary 权限树-按分组聚合，用于角色授权勾选
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=[]model.PermNode} "返回数据"
// @Router /api/permissions/tree [get]
func PermTree(ctx echo.Context) error {
	rctx := ctx.Request().Context()
	mods, err := repo.PermTree(rctx)
	if err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.PermTree")
		return reply.Fail(ctx, "查询失败")
	}
	return reply.Succ(ctx, "ok", mods)
}

// PermGet doc
// @Tags permission
// @Summary 权限详情
// @Security KeyAuth
// @Param id path string true "权限id"
// @Success 200 {object} reply.Reply{data=model.Perm} "返回数据"
// @Router /api/permissions/{id} [get]
func PermGet(ctx echo.Context) error {
	rctx := ctx.Request().Context()
	mod, err := repo.PermGet(rctx, ctx.Param("id"))
	if err != nil {
		logs.Ctx(rctx).Err(err).Str("id", ctx.Param("id")).Error("repo.PermGet")
		return reply.Fail(ctx, "权限不存在")
	}
	return reply.Succ(ctx, "ok", mod)
}
