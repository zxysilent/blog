package sysctl

import (
	"blog/model"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// SysRoleGet doc
// @Tags sysrole
// @Summary 通过id获取单条角色信息
// @Param id path int true "pk id" default(1)
// @Router /sys/role/get/{id} [get]
func SysRoleGet(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	mod, has := model.SysRoleGet(id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到角色信息"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// SysRoleAll doc
// @Tags sysrole
// @Summary 获取所有角色信息
// @Router /sys/role/all [get]
func SysRoleAll(ctx echo.Context) error {
	mods, err := model.SysRoleAll()
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到角色信息", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// SysRolePage doc
// @Tags sysrole
// @Summary 获取角色分页信息
// @Param cid path int true "分类id" default(1)
// @Param pi query int true "分页数" default(1)
// @Param ps query int true "每页条数[5,20]" default(5)
// @Router /sys/role/page/{cid} [get]
func SysRolePage(ctx echo.Context) error {
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
	count := model.SysRoleCount()
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := model.SysRolePage(ipt.Pi, ipt.Ps)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, int(count)))
}

// SysRoleAdd doc
// @Tags sysrole
// @Summary 添加角色信息
// @Param token query string true "hmt" default(token)
// @Router /sys/role/add [post]
func SysRoleAdd(ctx echo.Context) error {
	ipt := &model.SysRole{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// ipt.Utime = time.Now()
	err = model.SysRoleAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// SysRoleEdit doc
// @Tags sysrole
// @Summary 修改角色信息
// @Param token query string true "hmt" default(token)
// @Router /sys/role/edit [post]
func SysRoleEdit(ctx echo.Context) error {
	ipt := &model.SysRole{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// ipt.Utime = time.Now()
	err = model.SysRoleEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// SysRoleDrop doc
// @Tags sysrole
// @Summary 通过id删除单条角色信息
// @Param id path int true "pk id" default(1)
// @Param token query string true "hmt" default(token)
// @Router /sys/role/drop/{id} [get]
func SysRoleDrop(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	err = model.SysRoleDrop(id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
