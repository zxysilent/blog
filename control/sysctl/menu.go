package sysctl

import (
	"blog/model"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// SysMenuGet doc
// @Tags sysmenu
// @Summary 通过id获取单条菜单导航信息
// @Param id path int true "pk id" default(1)
// @Router /api/sysmenu/get/{id} [get]
func SysMenuGet(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	mod, has := model.SysMenuGet(id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到菜单导航信息"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// SysMenuAll doc
// @Tags sysmenu
// @Summary 获取所有菜单导航信息
// @Router /api/sysmenu/all [get]
func SysMenuAll(ctx echo.Context) error {
	mods, err := model.SysMenuAll()
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到菜单导航信息", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// SysMenuPage doc
// @Tags sysmenu
// @Summary 获取菜单导航分页信息
// @Param cid path int true "分类id" default(1)
// @Param pi query int true "分页数" default(1)
// @Param ps query int true "每页条数[5,20]" default(5)
// @Router /api/sysmenu/page/{cid} [get]
func SysMenuPage(ctx echo.Context) error {
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
	count := model.SysMenuCount()
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := model.SysMenuPage(ipt.Pi, ipt.Ps)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, int(count)))
}

// SysMenuAdd doc
// @Tags sysmenu
// @Summary 添加菜单导航信息
// @Param token query string true "hmt" default(token)
// @Router /adm/sysmenu/add [post]
func SysMenuAdd(ctx echo.Context) error {
	ipt := &model.SysMenu{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// ipt.Utime = time.Now()
	err = model.SysMenuAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// SysMenuEdit doc
// @Tags sysmenu
// @Summary 修改菜单导航信息
// @Param token query string true "hmt" default(token)
// @Router /adm/sysmenu/edit [post]
func SysMenuEdit(ctx echo.Context) error {
	ipt := &model.SysMenu{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// ipt.Utime = time.Now()
	err = model.SysMenuEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// SysMenuDrop doc
// @Tags sysmenu
// @Summary 通过id删除单条菜单导航信息
// @Param id path int true "pk id" default(1)
// @Param token query string true "hmt" default(token)
// @Router /adm/sysmenu/drop/{id} [get]
func SysMenuDrop(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	err = model.SysMenuDrop(id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
