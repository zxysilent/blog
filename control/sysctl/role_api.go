package sysctl

import (
	"blog/model"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// RoleApiGet doc
// @Tags sysroleauth
// @Summary 通过id获取单条角色接口信息
// @Param id path int true "pk id" default(1)
// @Router /sys/roleauth/get/{id} [get]
func RoleApiGet(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	mod, has := model.RoleApiGet(id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到角色接口信息"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// RoleApiAll doc
// @Tags sysroleauth
// @Summary 获取所有角色接口信息
// @Router /sys/roleauth/all [get]
func RoleApiAll(ctx echo.Context) error {
	mods, err := model.RoleApiAll()
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到角色接口信息", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// RoleApiPage doc
// @Tags sysroleauth
// @Summary 获取角色接口分页信息
// @Param cid path int true "分类id" default(1)
// @Param pi query int true "分页数" default(1)
// @Param ps query int true "每页条数[5,20]" default(5)
// @Router /sys/roleauth/page/{cid} [get]
func RoleApiPage(ctx echo.Context) error {
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
	count := model.RoleApiCount()
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := model.RoleApiPage(ipt.Pi, ipt.Ps)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, int(count)))
}

// RoleApiAdd doc
// @Tags sysroleauth
// @Summary 添加角色接口信息
// @Param token query string true "hmt" default(token)
// @Router /sys/roleauth/add [post]
func RoleApiAdd(ctx echo.Context) error {
	ipt := &model.RoleApi{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// ipt.Utime = time.Now()
	err = model.RoleApiAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// RoleApiEdit doc
// @Tags sysroleauth
// @Summary 修改角色接口信息
// @Param token query string true "hmt" default(token)
// @Router /sys/roleauth/edit [post]
func RoleApiEdit(ctx echo.Context) error {
	ipt := &model.RoleApi{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// ipt.Utime = time.Now()
	err = model.RoleApiEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// RoleApiDrop doc
// @Tags sysroleauth
// @Summary 通过id删除单条角色接口信息
// @Param id path int true "pk id" default(1)
// @Param token query string true "hmt" default(token)
// @Router /sys/roleauth/drop/{id} [get]
func RoleApiDrop(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	err = model.RoleApiDrop(id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
