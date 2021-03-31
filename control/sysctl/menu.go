package sysctl

import (
	"blog/model"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// MenuAll doc
// @Tags sysmenu
// @Summary 获取所有菜单导航菜单导航树
// @Success 200 {object} model.Reply{data=model.Menu} "成功数据"
// @Router /api/menu/tree [get]
func MenuTree(ctx echo.Context) error {
	mods, err := model.MenuTree()
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到菜单导航树", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// MenuGet doc
// @Tags sysmenu
// @Summary 通过id获取单条菜单导航信息
// @Param id path int true "pk id" default(1)
// @Router /sys/menu/get/{id} [get]
func MenuGet(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	mod, has := model.MenuGet(id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到菜单导航信息"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// MenuAll doc
// @Tags sysmenu
// @Summary 获取所有菜单导航信息
// @Router /sys/menu/all [get]
func MenuAll(ctx echo.Context) error {
	mods, err := model.MenuAll()
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到菜单导航信息", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// MenuPage doc
// @Tags sysmenu
// @Summary 获取菜单导航分页信息
// @Param cid path int true "分类id" default(1)
// @Param pi query int true "分页数" default(1)
// @Param ps query int true "每页条数[5,20]" default(5)
// @Router /sys/menu/page/{cid} [get]
func MenuPage(ctx echo.Context) error {
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
	count := model.MenuCount()
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := model.MenuPage(ipt.Pi, ipt.Ps)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, int(count)))
}

// MenuAdd doc
// @Tags sysmenu
// @Summary 添加菜单导航信息
// @Param token query string true "hmt" default(token)
// @Router /sys/menu/add [post]
func MenuAdd(ctx echo.Context) error {
	ipt := &model.Menu{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// ipt.Utime = time.Now()
	err = model.MenuAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// MenuEdit doc
// @Tags sysmenu
// @Summary 修改菜单导航信息
// @Param token query string true "hmt" default(token)
// @Router /sys/menu/edit [post]
func MenuEdit(ctx echo.Context) error {
	ipt := &model.Menu{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// ipt.Utime = time.Now()
	err = model.MenuEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// MenuDrop doc
// @Tags sysmenu
// @Summary 通过id删除单条菜单导航信息
// @Param id path int true "pk id" default(1)
// @Param token query string true "hmt" default(token)
// @Router /sys/menu/drop/{id} [get]
func MenuDrop(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	err = model.MenuDrop(id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
