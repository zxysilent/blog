package sysctl

import (
	"blog/model"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// ------------------------------------------------------ 配置中心 ------------------------------------------------------

// GlobalGet doc
// @Tags global
// @Summary 获取单个global信息
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /api/global/get/{id} [get]
func GlobalGet(ctx echo.Context) error {
	mod, has := model.GlobalGet()
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到global信息"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// GlobalEdit doc
// @Tags global
// @Summary 编辑global信息
// @Param token query string true "token"
// @Param body body model.Global true "请求数据"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/global/edit [post]
func GlobalEdit(ctx echo.Context) error {
	ipt := &model.Global{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	err = model.GlobalEdit(ipt, "")
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
