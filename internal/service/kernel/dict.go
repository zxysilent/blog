package kernel

import (
	"blog/internal/model"
	"blog/internal/repo"
	"blog/internal/utils"
	"blog/pkg/reply"
	"encoding/json"
	"time"

	"github.com/labstack/echo/v4"
)

// DictVal doc
// @Tags kernel,dict
// @Alias DictVal
// @Summary 通过id获取单条字典
// @Param key path string true "key"
// @Success 200 {object} reply.Reply{data=model.Dict} "返回数据"
// @Router /api/dict/{key} [get]
func DictVal(ctx echo.Context) error {
	rctx := utils.WrapCtx(ctx)
	key := ctx.Param("key")
	mod, err := repo.GetAny[model.Dict](rctx, key)
	if err != nil {
		return reply.Fail(ctx, "不存在", err)
	}
	return reply.Succ(ctx, "succ", json.RawMessage(mod.Value))
}

// DictGet doc
// @Tags kernel,dict
// @Alias DictGet
// @Summary 通过id获取单条字典
// @Param key path string true "key"
// @Param type query string false "raw,"
// @Success 200 {object} reply.Reply{data=model.Dict} "返回数据"
// @Router /api/dict/get/{key} [get]
func DictGet(ctx echo.Context) error {
	rctx := utils.WrapCtx(ctx)
	key := ctx.Param("key")
	mod, err := repo.GetAny[model.Dict](rctx, key)
	if err != nil {
		return reply.Fail(ctx, "不存在", err)
	}
	if ctx.QueryParam("type") == "raw" {
		return reply.Succ(ctx, "succ", json.RawMessage(mod.Value))
	}
	return reply.Succ(ctx, "succ", mod)
}

// DictList doc
// @Auth mgmt
// @Tags kernel,dict
// @Summary 获取字典分页
// @Alias DictList
// @Param query query model.DictListArgs true "请求数据"
// @Success 200 {object} reply.Reply{data=[]model.Dict} "返回数据"
// @Router /api/dict/list [get]
func DictList(ctx echo.Context) error {
	rctx := utils.WrapCtx(ctx)
	args := &model.DictListArgs{}
	if err := ctx.Bind(args); err != nil {
		return reply.Fail(ctx, "输入有误", err)
	}
	mods, count, err := repo.DictList(rctx, args)
	if err != nil {
		return reply.Fail(ctx, "查询数据错误", err)
	}
	if len(mods) < 1 {
		return reply.Fail(ctx, "未查询到数据")
	}
	return reply.List(ctx, "succ", mods, count)
}

// DictAdd doc
// @Auth mgmt
// @Tags kernel,dict
// @Summary 添加字典
// @Security KeyAuth
// @Param body body model.Dict true "请求数据"
// @Success 200 {object} reply.Reply{data=string} "返回数据"
// @Router /api/dict/add [post]
func DictAdd(ctx echo.Context) error {
	rctx := utils.WrapCtx(ctx)
	ipt := &model.Dict{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err)
	}
	// if err = ipt.Check(); err != nil {
	// 	return reply.Fail(ctx, "输入有误", err)
	// }
	ipt.Created = time.Now().UnixMilli()
	ipt.Updated = ipt.Created
	err := repo.AddAny(rctx, ipt)
	if err != nil {
		return reply.Fail(ctx, "添加失败", err)
	}
	return reply.Succ(ctx, "succ")
}

// DictEdit doc
// @Auth mgmt
// @Tags kernel,dict
// @Summary 修改字典
// @Security KeyAuth
// @Param body body model.Dict true "请求数据"
// @Success 200 {object} reply.Reply{data=string} "返回数据"
// @Router /api/dict/edit [put]
func DictEdit(ctx echo.Context) error {
	rctx := utils.WrapCtx(ctx)
	ipt := &model.Dict{}
	err := ctx.Bind(ipt)
	if err != nil {
		return reply.Fail(ctx, "输入有误", err)
	}
	if err = ipt.Check(); err != nil {
		return reply.Fail(ctx, "输入有误", err)
	}
	ipt.Updated = time.Now().UnixMilli()
	err = repo.DictEdit(rctx, ipt)
	if err != nil {
		return reply.Fail(ctx, "修改失败", err)
	}
	return reply.Succ(ctx, "succ")
}

// DictDrop doc
// @Auth mgmt
// @Tags kernel,dict
// @Summary 通过key删除单条字典
// @Param key path string true "key"
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=string} "返回数据"
// @Router /api/dict/{key} [delete]
func DictDrop(ctx echo.Context) error {
	rctx := utils.WrapCtx(ctx)
	key := ctx.Param("key")
	mod, err := repo.GetAny[model.Dict](rctx, key)
	if err != nil {
		return reply.Fail(ctx, "不存在", err)
	}
	if mod.Inner {
		return reply.Fail(ctx, "内置数据无法删除")
	}
	err = repo.DelAny[model.Dict](rctx, mod.Key)
	if err != nil {
		return reply.Fail(ctx, "删除失败", err)
	}
	return reply.Succ(ctx, "succ")
}
