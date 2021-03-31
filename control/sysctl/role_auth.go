package sysctl

import (
	"blog/model"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// RoleAuthGet doc
// @Tags sysroleauth
// @Summary 通过id获取单条角色认证信息
// @Param id path int true "pk id" default(1)
// @Router /sys/roleauth/get/{id} [get]
func RoleAuthGet(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	mod, has := model.RoleAuthGet(id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到角色认证信息"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// RoleAuthAll doc
// @Tags sysroleauth
// @Summary 获取所有角色认证信息
// @Router /sys/roleauth/all [get]
func RoleAuthAll(ctx echo.Context) error {
	mods, err := model.RoleAuthAll()
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到角色认证信息", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// RoleAuthPage doc
// @Tags sysroleauth
// @Summary 获取角色认证分页信息
// @Param cid path int true "分类id" default(1)
// @Param pi query int true "分页数" default(1)
// @Param ps query int true "每页条数[5,20]" default(5)
// @Router /sys/roleauth/page/{cid} [get]
func RoleAuthPage(ctx echo.Context) error {
	// cid, err := strconv.Atoi(ctx.Param("cid"))
	// if err != nil {
	//  return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	// }
	ipt := &model.Page{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if ipt.Ps > 20 || ipt.Ps < 5 {
		return ctx.JSON(utils.ErrIpt("分页大小输入错误", ipt.Ps))
	}
	count := model.RoleAuthCount()
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := model.RoleAuthPage(ipt.Pi, ipt.Ps)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, int(count)))
}

// RoleAuthAdd doc
// @Tags sysroleauth
// @Summary 添加角色认证信息
// @Param token query string true "hmt" default(token)
// @Router /sys/roleauth/add [post]
func RoleAuthAdd(ctx echo.Context) error {
	ipt := &model.RoleAuth{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// ipt.Utime = time.Now()
	err = model.RoleAuthAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// RoleAuthEdit doc
// @Tags sysroleauth
// @Summary 修改角色认证信息
// @Param token query string true "hmt" default(token)
// @Router /sys/roleauth/edit [post]
func RoleAuthEdit(ctx echo.Context) error {
	ipt := &model.RoleAuth{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// ipt.Utime = time.Now()
	err = model.RoleAuthEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// RoleAuthDrop doc
// @Tags sysroleauth
// @Summary 通过id删除单条角色认证信息
// @Param id path int true "pk id" default(1)
// @Param token query string true "hmt" default(token)
// @Router /sys/roleauth/drop/{id} [get]
func RoleAuthDrop(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	err = model.RoleAuthDrop(id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
