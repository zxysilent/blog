package sysctl

import (
	"blog/model"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// DictGet doc
// @Tags dict
// @Summary 通过id获取单条字典
// @Param key query string true "key"
// @Success 200 {object} model.Reply{data=model.Dict} "返回数据"
// @Router /api/dict/get [get]
func DictGet(ctx echo.Context) error {
	key := ctx.QueryParam("key")
	mod, has := model.DictGet(key)
	if !has {
		return ctx.JSON(utils.Fail("不存在"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// DictAllCache doc
// @Tags dict
// @Summary 缓存所有字典
// @Success 200 {object} model.Reply "成功数据"
// @Router /api/dict/all/cache [get]
func DictAllCache(ctx echo.Context) error {
	if err := model.DictCache(); err != nil {
		return ctx.JSON(utils.ErrOpt("缓存失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// DictAll doc
// @Tags dict
// @Summary 获取所有字典
// @Success 200 {object} model.Reply{data=[]model.Dict} "返回数据"
// @Router /api/dict/all [get]
func DictAll(ctx echo.Context) error {
	mods, err := model.DictAll()
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到字典", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// DictPage doc
// @Tags dict
// @Summary 获取字典分页
// @Param cid path int true "分类id" default(1)
// @Param pi query int true "分页数" default(1)
// @Param ps query int true "每页条数[5,30]" default(5)
// @Success 200 {object} model.Reply{data=[]model.Dict} "返回数据"
// @Router /api/dict/page [get]
func DictPage(ctx echo.Context) error {
	ipt := &model.Page{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if ipt.Ps > 30 || ipt.Ps < 1 {
		return ctx.JSON(utils.ErrIpt("分页大小输入错误", ipt.Ps))
	}
	count := model.DictCount()
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := model.DictPage(ipt.Pi, ipt.Ps)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, int(count)))
}

// DictAdd doc
// @Tags dict
// @Summary 添加字典
// @Param token query string true "token"
// @Param body body model.Dict true "请求数据"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/dict/add [post]
func DictAdd(ctx echo.Context) error {
	ipt := &model.Dict{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	ipt.Ctime = time.Now()
	err = model.DictAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// DictEdit doc
// @Tags dict
// @Summary 修改字典
// @Param token query string true "token"
// @Param body body model.Dict true "请求数据"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/dict/edit [post]
func DictEdit(ctx echo.Context) error {
	ipt := &model.Dict{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	ipt.Ctime = time.Now()
	err = model.DictEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// DictDrop doc
// @Tags dict
// @Summary 通过key删除单条字典
// @Param key query string true "key"
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/dict/drop [post]
func DictDrop(ctx echo.Context) error {
	key := ctx.QueryParam("key")
	mod, has := model.DictGet(key)
	if !has {
		return ctx.JSON(utils.Fail("不存在"))
	}
	if mod.Inner {
		return ctx.JSON(utils.ErrOpt("内置数据无法删除"))
	}
	err := model.DictDrop(key)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// ------------------------------------------------------ 配置中心 ------------------------------------------------------

// GlobalGet doc
// @Tags global
// @Summary 获取单个global信息
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /api/global/get/{id} [get]
func GlobalGet(ctx echo.Context) error {
	mod, has := model.GlobalGet()
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到global信息"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// GlobalEdit doc
// @Tags global
// @Summary 编辑global信息
// @Param token query string true "token"
// @Param body body model.Global true "请求数据"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/global/edit [post]
func GlobalEdit(ctx echo.Context) error {
	ipt := &model.Global{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	err = model.GlobalEdit(ipt, "")
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
