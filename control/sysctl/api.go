package sysctl

import (
	"blog/model"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// ApiGet doc
// @Tags sysauth
// @Summary 通过id获取单条接口信息
// @Param id path int true "pk id" default(1)
// @Router /sys/auth/get/{id} [get]
func ApiGet(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	mod, has := model.ApiGet(id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到接口信息"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// ApiAll doc
// @Tags sysauth
// @Summary 获取所有接口信息
// @Router /sys/auth/all [get]
func ApiAll(ctx echo.Context) error {
	mods, err := model.ApiAll()
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到接口信息", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// ApiPage doc
// @Tags sysauth
// @Summary 获取接口分页信息
// @Param cid path int true "分类id" default(1)
// @Param pi query int true "分页数" default(1)
// @Param ps query int true "每页条数[5,20]" default(5)
// @Router /sys/auth/page/{cid} [get]
func ApiPage(ctx echo.Context) error {
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
	count := model.ApiCount()
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := model.ApiPage(ipt.Pi, ipt.Ps)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, int(count)))
}

// ApiAdd doc
// @Tags sysauth
// @Summary 添加接口信息
// @Param token query string true "hmt" default(token)
// @Router /sys/auth/add [post]
func ApiAdd(ctx echo.Context) error {
	ipt := &model.Api{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// ipt.Utime = time.Now()
	err = model.ApiAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// ApiEdit doc
// @Tags sysauth
// @Summary 修改接口信息
// @Param token query string true "hmt" default(token)
// @Router /sys/auth/edit [post]
func ApiEdit(ctx echo.Context) error {
	ipt := &model.Api{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// ipt.Utime = time.Now()
	err = model.ApiEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// ApiDrop doc
// @Tags sysauth
// @Summary 通过id删除单条接口信息
// @Param id path int true "pk id" default(1)
// @Param token query string true "hmt" default(token)
// @Router /sys/auth/drop/{id} [get]
func ApiDrop(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	err = model.ApiDrop(id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
