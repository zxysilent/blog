package sysctl

import (
	"blog/model"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// SysAuthGet doc
// @Tags sysauth
// @Summary 通过id获取单条认证信息
// @Param id path int true "pk id" default(1)
// @Router /api/sysauth/get/{id} [get]
func SysAuthGet(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	mod, has := model.SysAuthGet(id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到认证信息"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// SysAuthAll doc
// @Tags sysauth
// @Summary 获取所有认证信息
// @Router /api/sysauth/all [get]
func SysAuthAll(ctx echo.Context) error {
	mods, err := model.SysAuthAll()
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到认证信息", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// SysAuthPage doc
// @Tags sysauth
// @Summary 获取认证分页信息
// @Param cid path int true "分类id" default(1)
// @Param pi query int true "分页数" default(1)
// @Param ps query int true "每页条数[5,20]" default(5)
// @Router /api/sysauth/page/{cid} [get]
func SysAuthPage(ctx echo.Context) error {
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
	count := model.SysAuthCount()
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := model.SysAuthPage(ipt.Pi, ipt.Ps)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, int(count)))
}

// SysAuthAdd doc
// @Tags sysauth
// @Summary 添加认证信息
// @Param token query string true "hmt" default(token)
// @Router /adm/sysauth/add [post]
func SysAuthAdd(ctx echo.Context) error {
	ipt := &model.SysAuth{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// ipt.Utime = time.Now()
	err = model.SysAuthAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// SysAuthEdit doc
// @Tags sysauth
// @Summary 修改认证信息
// @Param token query string true "hmt" default(token)
// @Router /adm/sysauth/edit [post]
func SysAuthEdit(ctx echo.Context) error {
	ipt := &model.SysAuth{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// ipt.Utime = time.Now()
	err = model.SysAuthEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// SysAuthDrop doc
// @Tags sysauth
// @Summary 通过id删除单条认证信息
// @Param id path int true "pk id" default(1)
// @Param token query string true "hmt" default(token)
// @Router /adm/sysauth/drop/{id} [get]
func SysAuthDrop(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	err = model.SysAuthDrop(id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
