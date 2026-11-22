package kernel

import (
	"blog/internal/model"
	"blog/internal/repo"
	"blog/pkg/reply"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/logs"
)

// LogList doc
// @Tags log
// @Summary 操作日志列表
// @Security KeyAuth
// @Param query query model.LogListArgs true "请求数据"
// @Success 200 {object} reply.Reply{data=[]model.Log} "返回数据"
// @Router /api/logs [get]
func LogList(ctx echo.Context) error {
	args := &model.LogListArgs{}
	if err := ctx.Bind(args); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	rctx := ctx.Request().Context()
	mods, count, err := repo.LogList(rctx, args)
	if err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.LogList")
		return reply.Fail(ctx, "查询失败")
	}
	return reply.List(ctx, "ok", mods, count)
}

// LogGet doc
// @Tags log
// @Summary 操作日志详情
// @Security KeyAuth
// @Param id path string true "日志id"
// @Success 200 {object} reply.Reply{data=model.Log} "返回数据"
// @Router /api/logs/{id} [get]
func LogGet(ctx echo.Context) error {
	id := ctx.Param("id")
	rctx := ctx.Request().Context()
	mod, err := repo.LogGet(rctx, id)
	if err != nil {
		logs.Ctx(rctx).Err(err).Str("id", id).Error("repo.LogGet")
		return reply.Fail(ctx, "日志不存在")
	}
	return reply.Succ(ctx, "ok", mod)
}

// LogDrop doc
// @Tags log
// @Summary 删除单条操作日志
// @Security KeyAuth
// @Param id path string true "日志id"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/logs/{id} [delete]
func LogDrop(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return reply.Fail(ctx, "参数有误")
	}
	rctx := ctx.Request().Context()
	if _, err := repo.LogGet(rctx, id); err != nil {
		return reply.Fail(ctx, "日志不存在")
	}
	if err := repo.LogDrop(rctx, id); err != nil {
		logs.Ctx(rctx).Err(err).Str("id", id).Error("repo.LogDrop")
		return reply.Fail(ctx, "删除失败")
	}
	return reply.Succ(ctx, "删除成功")
}

// LogClearArgs 日志清理参数
type LogClearArgs struct {
	Before int64 `json:"before" form:"before"` //只删早于该时刻（毫秒）的日志，0 表示全部清空
}

// LogClear doc
// @Tags log
// @Summary 批量清理操作日志
// @Security KeyAuth
// @Param body body kernel.LogClearArgs true "请求数据"
// @Success 200 {object} reply.Reply{data=int64} "返回数据"
// @Router /api/logs [delete]
func LogClear(ctx echo.Context) error {
	ipt := &LogClearArgs{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	rctx := ctx.Request().Context()
	count, err := repo.LogClear(rctx, ipt.Before)
	if err != nil {
		logs.Ctx(rctx).Err(err).Int64("before", ipt.Before).Error("repo.LogClear")
		return reply.Fail(ctx, "清理失败")
	}
	return reply.Succ(ctx, "清理成功", count)
}
