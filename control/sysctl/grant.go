package sysctl

import (
	"blog/model"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// GrantAll doc
// @Tags grant
// @Summary 获取所有授权树
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=[]model.Grant} "返回数据"
// @Router /api/grant/tree [get]
func GrantTree(ctx echo.Context) error {
	rid, _ := ctx.Get("rid").(int)
	// 授权给别人的权限,不能高于自己的权限
	mods, err := model.GrantTree(rid)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到授权树", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// GrantGet doc
// @Tags grant
// @Summary 通过id获取单条授权信息
// @Param id query int true "id"
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=model.Grant} "返回数据"
// @Router /adm/grant/get [get]
func GrantGet(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod, has := model.GrantGet(ipt.Id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到授权信息"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// GrantAll doc
// @Tags grant
// @Summary 获取所有授权信息
// @Param token query string true "token"
// @Param body body object{slt=bool,show=bool} true "json"
// @Success 200 {object} model.Reply{data=[]model.Grant} "返回数据"
// @Router /adm/grant/all [get]
func GrantAll(ctx echo.Context) error {
	mods, err := model.GrantAll()
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到授权信息", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// GrantPage doc
// @Tags grant
// @Summary 获取授权分页信息
// @Param pi query int true "分页数"
// @Param ps query int true "每页条数[5,30]"
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=[]model.Grant} "返回数据"
// @Param token query string true "token"
// @Router /adm/grant/page [get]
func GrantPage(ctx echo.Context) error {
	ipt := &model.Page{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if ipt.Ps > 20 || ipt.Ps < 5 {
		return ctx.JSON(utils.ErrIpt("分页大小输入错误", ipt.Ps))
	}
	count := model.GrantCount()
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := model.GrantPage(ipt.Pi, ipt.Ps)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, int(count)))
}

// GrantAdd doc
// @Tags grant
// @Summary 添加授权信息
// @Param token query string true "token"
// @Param body body model.Grant true "请求数据"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/grant/add [post]
func GrantAdd(ctx echo.Context) error {
	ipt := &model.Grant{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	ipt.Ctime = time.Now()
	err = model.GrantAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// GrantEdit doc
// @Tags grant
// @Summary 修改授权信息
// @Param token query string true "token"
// @Param body body model.Grant true "请求数据"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/grant/edit [post]
func GrantEdit(ctx echo.Context) error {
	ipt := &model.Grant{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	err = model.GrantEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// GrantDrop doc
// @Tags grant
// @Summary 通过id删除单条授权信息
// @Param token query string true "token"
// @Param body body model.IptId true "json"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/grant/drop [post]
func GrantDrop(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod, has := model.GrantGet(ipt.Id)
	if !has {
		return ctx.JSON(utils.ErrOpt("删除失败", "不存在当前数据"))
	}
	if mod.Inner {
		return ctx.JSON(utils.ErrOpt("内置数据无法删除"))
	}
	err = model.GrantDrop(ipt.Id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
