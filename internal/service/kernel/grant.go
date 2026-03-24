package kernel

import (
	"blog/internal/model"
	"blog/internal/repo"
	"blog/internal/utils"
	"blog/pkg/reply"
	"time"

	"github.com/labstack/echo/v4"
)

// GrantGet doc
// @Auth mgmt
// @Tags kernel,grant
// @Summary 通过id获取单条授权信息
// @Param id query int true "id"
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=model.Grant} "返回数据"
// @Router /api/grant/{id} [get]
func GrantGet(ctx echo.Context) error {
	rctx := utils.WrapCtx(ctx)
	mod, err := repo.GetAny[model.Grant](rctx, ctx.Param("id"))
	if err != nil {
		return reply.Fail(ctx, "未查询到数据", err)
	}
	return reply.Succ(ctx, "succ", mod)
}

// GrantPage doc
// @Auth mgmt
// @Tags kernel,grant
// @Summary 获取授权分页信息
// @Security KeyAuth
// @Param query query model.List true "请求数据"
// @Success 200 {object} reply.Reply{data=[]model.Grant} "返回数据"
// @Security KeyAuth
// @Router /api/grant/page [get]
func GrantPage(ctx echo.Context) error {
	// ipt := &model.Page{}
	// err := ctx.Bind(ipt)
	// if err != nil {
	// 	return reply.Fail(ctx, "输入有误", err)
	// }
	// if err = ipt.Stat(); err != nil {
	// 	return reply.Fail(ctx, "输入有误", err)
	// }
	// count := repo.GrantCount()
	// if count < 1 {
	// 	return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	// }
	// mods, err := repo.GrantPage(ipt.Pi, ipt.Ps)
	// if err != nil {
	// 	return reply.Fail(ctx, "查询数据错误", err)
	// }
	// if len(mods) < 1 {
	// 	return reply.Fail(ctx, "未查询到数据")
	// }
	// return reply.List(ctx, "succ", mods, count)
	return nil
}

// GrantAdd doc
// @Auth mgmt
// @Tags kernel,grant
// @Summary 添加授权信息
// @Security KeyAuth
// @Param body body model.Grant true "请求数据"
// @Success 200 {object} reply.Reply{data=string} "返回数据"
// @Router /api/grant/add [post]
func GrantAdd(ctx echo.Context) error {
	ipt := &model.Grant{}
	err := ctx.Bind(ipt)
	if err != nil {
		return reply.Fail(ctx, "输入有误", err)
	}
	ipt.Created = time.Now().UnixMilli()
	ipt.Updated = ipt.Created
	// err = repo.GrantAdd(ipt)
	// if err != nil {
	// 	return reply.Fail(ctx, "添加失败", err)
	// }
	return reply.Succ(ctx, "succ")
}

// GrantEdit doc
// @Auth mgmt
// @Tags kernel,grant
// @Summary 修改授权信息
// @Security KeyAuth
// @Param body body model.Grant true "请求数据"
// @Success 200 {object} reply.Reply{data=string} "返回数据"
// @Router /api/grant/edit [post]
func GrantEdit(ctx echo.Context) error {
	ipt := &model.Grant{}
	err := ctx.Bind(ipt)
	if err != nil {
		return reply.Fail(ctx, "输入有误", err)
	}
	ipt.Updated = time.Now().UnixMilli()
	// err = repo.GrantEdit(ipt)
	// if err != nil {
	// 	return reply.Fail(ctx, "修改失败", err)
	// }
	return reply.Succ(ctx, "succ")
}

// GrantDrop doc
// @Auth mgmt
// @Tags kernel,grant
// @Summary 通过id删除单条授权信息
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=string} "返回数据"
// @Router /api/grant/{id} [delete]
func GrantDrop(ctx echo.Context) error {
	rctx := utils.WrapCtx(ctx)
	mod, err := repo.GetAny[model.Grant](rctx, ctx.Param("id"))
	if err != nil {
		return reply.Fail(ctx, "未查询到数据", err)
	}
	if mod.Inner {
		return reply.Fail(ctx, "内置数据无法删除")
	}
	// err = repo.GrantDrop(mod.Id)
	// if err != nil {
	// 	return reply.Fail(ctx, "删除失败", err)
	// }
	return reply.Succ(ctx, "succ")
}
