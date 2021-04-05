package sysctl

import (
	"blog/model"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// RoleApiAll doc
// @Tags role,api
// @Summary 通过角色id获取所有api接口信息
// @Param id query int true "id"
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=model.Role} "成功数据"
// @Router /adm/role/api/get [get]
func RoleApiAll(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mods, has := model.RoleMenuAll(ipt.Id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到角色菜单导航信息"))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// RoleApiEdit doc
// @Tags role,api
// @Summary 修改角色接口信息
// @Param id query int true "id"
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=string} "成功数据"
// @Router /adm/role/api/edit [post]
func RoleApiEdit(ctx echo.Context) error {
	ipt := &model.RoleApi{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	err = model.RoleApiEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
