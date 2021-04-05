package sysctl

import (
	"blog/model"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// RoleMenuAll doc
// @Tags role,menu
// @Summary 通过角色id获取所有角色菜单信息
// @Param id query int true "id"
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=model.Role} "成功数据"
// @Router /adm/role/menu/get [get]
func RoleMenuAll(ctx echo.Context) error {
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

// RoleMenuEdit doc
// @Tags role,menu
// @Summary 修改菜单导航信息菜单导航
// @Param id query int true "id"
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=string} "成功数据"
// @Router /adm/role/menu/edit [post]
func RoleMenuEdit(ctx echo.Context) error {
	ipt := &model.RoleMenu{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	err = model.RoleMenuEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
