package sysctl

import (
	"blog/model"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// SysRoleMenuGet doc
// @Tags sysrolemenu
// @Summary 通过id获取单条角色菜单导航信息
// @Param id path int true "pk id" default(1)
// @Router /sys/rolemenu/get/{id} [get]
func SysRoleMenuGet(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	mod, has := model.SysRoleMenuGet(id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到角色菜单导航信息"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// SysRoleMenuAll doc
// @Tags sysrolemenu
// @Summary 获取所有角色菜单导航信息
// @Router /sys/rolemenu/all [get]
func SysRoleMenuAll(ctx echo.Context) error {
	mods, err := model.SysRoleMenuAll()
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到角色菜单导航信息", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// SysRoleMenuPage doc
// @Tags sysrolemenu
// @Summary 获取角色菜单导航分页信息
// @Param cid path int true "分类id" default(1)
// @Param pi query int true "分页数" default(1)
// @Param ps query int true "每页条数[5,20]" default(5)
// @Router /sys/rolemenu/page/{cid} [get]
func SysRoleMenuPage(ctx echo.Context) error {
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
	count := model.SysRoleMenuCount()
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := model.SysRoleMenuPage(ipt.Pi, ipt.Ps)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, int(count)))
}

// SysRoleMenuAdd doc
// @Tags sysrolemenu
// @Summary 添加角色菜单导航信息
// @Param token query string true "hmt" default(token)
// @Router /sys/rolemenu/add [post]
func SysRoleMenuAdd(ctx echo.Context) error {
	ipt := &model.SysRoleMenu{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// ipt.Utime = time.Now()
	err = model.SysRoleMenuAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// SysRoleMenuEdit doc
// @Tags sysrolemenu
// @Summary 修改角色菜单导航信息
// @Param token query string true "hmt" default(token)
// @Router /sys/rolemenu/edit [post]
func SysRoleMenuEdit(ctx echo.Context) error {
	ipt := &model.SysRoleMenu{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// ipt.Utime = time.Now()
	err = model.SysRoleMenuEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// SysRoleMenuDrop doc
// @Tags sysrolemenu
// @Summary 通过id删除单条角色菜单导航信息
// @Param id path int true "pk id" default(1)
// @Param token query string true "hmt" default(token)
// @Router /sys/rolemenu/drop/{id} [get]
func SysRoleMenuDrop(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	err = model.SysRoleMenuDrop(id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
