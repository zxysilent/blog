package kernel

import (
	"blog/internal/model"
	"blog/internal/repo"
	"blog/internal/utils"
	"time"

	"github.com/labstack/echo/v4"
)

// GrantAll doc
// @Auth mgmt
// @Tags kernel
// @Summary 获取所有授权树
// @Security ApiKeyAuth
// @Success 200 {object} utils.Reply{data=[]model.Grant} "返回数据"
// @Router /api/grant/tree [get]
func GrantTree(ctx echo.Context) error {
	rid, _ := ctx.Get("rid").(int)
	// 授权给别人的权限,不能高于自己的权限
	mods, err := repo.GrantTree(rid)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到授权树", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// GrantGet doc
// @Auth mgmt
// @Tags kernel
// @Summary 通过id获取单条授权信息
// @Param id query int true "id"
// @Security ApiKeyAuth
// @Success 200 {object} utils.Reply{data=model.Grant} "返回数据"
// @Router /api/grant/get [get]
func GrantGet(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod, has := repo.GrantGet(ipt.Id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到授权信息"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// GrantPage doc
// @Auth mgmt
// @Tags kernel
// @Summary 获取授权分页信息
// @Security ApiKeyAuth
// @Param query query model.Page true "请求数据"
// @Success 200 {object} utils.Reply{data=[]model.Grant} "返回数据"
// @Security ApiKeyAuth
// @Router /api/grant/page [get]
func GrantPage(ctx echo.Context) error {
	ipt := &model.Page{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if err = ipt.Stat(); err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	count := repo.GrantCount()
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := repo.GrantPage(ipt.Pi, ipt.Ps)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, count))
}

// GrantAdd doc
// @Auth mgmt
// @Tags kernel
// @Summary 添加授权信息
// @Security ApiKeyAuth
// @Param body body model.Grant true "请求数据"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/grant/add [post]
func GrantAdd(ctx echo.Context) error {
	ipt := &model.Grant{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	ipt.Created = time.Now().UnixMilli()
	ipt.Updated = ipt.Created
	err = repo.GrantAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// GrantEdit doc
// @Auth mgmt
// @Tags kernel
// @Summary 修改授权信息
// @Security ApiKeyAuth
// @Param body body model.Grant true "请求数据"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/grant/edit [post]
func GrantEdit(ctx echo.Context) error {
	ipt := &model.Grant{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	ipt.Updated = time.Now().UnixMilli()
	err = repo.GrantEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// GrantDrop doc
// @Auth mgmt
// @Tags kernel
// @Summary 通过id删除单条授权信息
// @Security ApiKeyAuth
// @Param body body model.IptId true "json"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/grant/drop [post]
func GrantDrop(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod, has := repo.GrantGet(ipt.Id)
	if !has {
		return ctx.JSON(utils.ErrOpt("删除失败", "不存在当前数据"))
	}
	if mod.Inner {
		return ctx.JSON(utils.ErrOpt("内置数据无法删除"))
	}
	err = repo.GrantDrop(ipt.Id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
