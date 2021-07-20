package sysctl

import (
	"blog/model"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// MenuAll doc
// @Tags menu
// @Summary 获取所有菜单导航菜单导航树
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=[]model.Menu} "返回数据"
// @Router /api/menu/tree [get]
func MenuTree(ctx echo.Context) error {
	mods, err := model.MenuTree()
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到菜单导航树", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// MenuGet doc
// @Tags menu
// @Summary 通过id获取单条菜单导航信息
// @Param id query int true "id"
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=model.Menu} "返回数据"
// @Router /adm/menu/get [get]
func MenuGet(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod, has := model.MenuGet(ipt.Id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到菜单导航信息"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// MenuAll doc
// @Tags menu
// @Summary 获取所有菜单导航信息
// @Param token query string true "token"
// @Param body body object{slt=bool,show=bool} true "json"
// @Success 200 {object} model.Reply{data=[]model.Menu} "返回数据"
// @Router /adm/menu/all [get]
func MenuAll(ctx echo.Context) error {
	ipt := &struct {
		Slt  bool `query:"slt" json:"slt"`
		Root bool `query:"root" json:"root"`
	}{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mods, err := model.MenuAll(ipt.Root, ipt.Slt)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到菜单导航信息", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// MenuPage doc
// @Tags menu
// @Summary 获取菜单导航分页信息
// @Param pi query int true "分页数"
// @Param ps query int true "每页条数[5,30]"
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=[]model.Menu} "返回数据"
// @Param token query string true "token"
// @Router /adm/menu/page [get]
func MenuPage(ctx echo.Context) error {
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
// @Tags menu
// @Summary 添加菜单导航信息
// @Param token query string true "token"
// @Param body body model.Menu true "请求数据"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/menu/add [post]
func MenuAdd(ctx echo.Context) error {
	ipt := &model.Menu{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	ipt.Ctime = time.Now()
	err = model.MenuAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// MenuEdit doc
// @Tags menu
// @Summary 修改菜单导航信息
// @Param token query string true "token"
// @Param body body model.Menu true "请求数据"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/menu/edit [post]
func MenuEdit(ctx echo.Context) error {
	ipt := &model.Menu{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	err = model.MenuEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// ClassChgShow doc
// @Tags menu
// @Summary 修改菜单导航显示信息
// @Param token query string true "token"
// @Param body body object{id=int,show=bool} true "json"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/menu/edit/show [post]
func MenuEditShow(ctx echo.Context) error {
	ipt := &struct {
		Id   int  `json:"id"`
		Show bool `json:"show"`
	}{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod := &model.Menu{Id: ipt.Id, Show: ipt.Show}
	err = model.MenuEdit(mod, "Show")
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// MenuDrop doc
// @Tags menu
// @Summary 通过id删除单条菜单导航信息
// @Param token query string true "token"
// @Param body body model.IptId true "json"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/menu/drop [post]
func MenuDrop(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod, has := model.MenuGet(ipt.Id)
	if !has {
		return ctx.JSON(utils.ErrOpt("删除失败", "不存在当前数据"))
	}
	if mod.Inner {
		return ctx.JSON(utils.ErrOpt("内置数据无法删除"))
	}
	err = model.MenuDrop(ipt.Id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
