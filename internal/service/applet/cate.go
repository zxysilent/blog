package applet

import (
	"blog/internal/model"
	"blog/internal/repo"
	"blog/internal/utils"
	"time"

	"github.com/labstack/echo/v4"
)

// CateGet doc
// Auth
// @Tags cate
// @Summary 分类单条数据
// @Param id query int true "id"
// @Success 200 {object} utils.Reply{data=model.Cate} "返回数据"
// @Router /api/cate/get [get]
func CateGet(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod := &model.Cate{Id: ipt.Id}
	err = repo.CateGet(mod)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到数据", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// CateList doc
// Auth
// @Tags cate
// @Summary 分类列表数据
// @Param query query model.CateFilterList true "请求数据"
// @Success 200 {object} utils.Reply{data=[]model.Cate} "返回数据"
// @Router /api/cate/list [get]
func CateList(ctx echo.Context) error {
	ipt := &model.CateFilterList{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mods, err := repo.CateList(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到数据", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// CatePage doc
// Auth
// @Tags cate
// @Summary 获取分类分页
// @Param query query model.CateFilterPage true "请求数据"
// @Success 200 {object} utils.Reply{data=[]model.Cate} "返回数据"
// @Router /api/cate/page [get]
func CatePage(ctx echo.Context) error {
	// cid, err := strconv.Atoi(ctx.Param("cid"))
	// if err != nil {
	//  return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	// }
	ipt := &model.CateFilterPage{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if err = ipt.Stat(); err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mods, count, err := repo.CatePage(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if count == 0 || len(mods) == 0 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "empty"))
	}
	return ctx.JSON(utils.Page("succ", mods, count))
}

// CateAdd doc
// Auth
// @Tags cate
// @Summary 分类添加数据
// @Param token query string true "token"
// @Param body body model.Cate true "请求数据"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/cate/add [post]
func CateAdd(ctx echo.Context) error {
	ipt := &model.Cate{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	ipt.Updated = time.Now().UnixMilli()
	ipt.Created = ipt.Updated
	err = repo.CateAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// CateEdit doc
// Auth
// @Tags cate
// @Summary 分类修改数据
// @Param token query string true "token"
// @Param body body model.Cate true "请求数据"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/cate/edit [post]
func CateEdit(ctx echo.Context) error {
	ipt := &model.Cate{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	ipt.Updated = time.Now().UnixMilli()
	err = repo.CateEdit(ipt, "pid", "name", "color", "intro", "updated")
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// CateDrop doc
// Auth
// @Tags cate
// @Summary 分类删除数据
// @Param token query string true "token"
// @Param query query model.IptId true "请求数据"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/cate/drop [post]
func CateDrop(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	err = repo.CateDrop(&model.Cate{Id: ipt.Id})
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
