package sysctl

import (
	"blog/model"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// SysUserGet doc
// @Tags sysuser
// @Summary 通过id获取单条用户信息
// @Param id path int true "pk id" default(1)
// @Router /api/sysuser/get/{id} [get]
func SysUserGet(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	mod, has := model.SysUserGet(id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到用户信息"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// SysUserAll doc
// @Tags sysuser
// @Summary 获取所有用户信息
// @Router /api/sysuser/all [get]
func SysUserAll(ctx echo.Context) error {
	mods, err := model.SysUserAll()
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到用户信息", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// SysUserPage doc
// @Tags sysuser
// @Summary 获取用户分页信息
// @Param cid path int true "分类id" default(1)
// @Param pi query int true "分页数" default(1)
// @Param ps query int true "每页条数[5,20]" default(5)
// @Router /api/sysuser/page/{cid} [get]
func SysUserPage(ctx echo.Context) error {
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
	count := model.SysUserCount()
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := model.SysUserPage(ipt.Pi, ipt.Ps)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, int(count)))
}

// SysUserAdd doc
// @Tags sysuser
// @Summary 添加用户信息
// @Param token query string true "hmt" default(token)
// @Router /adm/sysuser/add [post]
func SysUserAdd(ctx echo.Context) error {
	ipt := &model.SysUser{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// ipt.Utime = time.Now()
	err = model.SysUserAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// SysUserEdit doc
// @Tags sysuser
// @Summary 修改用户信息
// @Param token query string true "hmt" default(token)
// @Router /adm/sysuser/edit [post]
func SysUserEdit(ctx echo.Context) error {
	ipt := &model.SysUser{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// ipt.Utime = time.Now()
	err = model.SysUserEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// SysUserDrop doc
// @Tags sysuser
// @Summary 通过id删除单条用户信息
// @Param id path int true "pk id" default(1)
// @Param token query string true "hmt" default(token)
// @Router /adm/sysuser/drop/{id} [get]
func SysUserDrop(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	err = model.SysUserDrop(id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
