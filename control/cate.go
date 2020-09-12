package control

import (
	"blog/model"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// CateAll doc
// @Tags 分类
// @Summary 所有分类
// @Param token query string true "凭证jwt" default(jwt)
// @Router /api/cate/all [get]
func CateAll(ctx echo.Context) error {
	mods, err := model.CateAll()
	if err != nil {
		return ctx.JSON(utils.ErrOpt(`未查询到分类信息`, err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt(`未查询到分类信息`, "len"))
	}
	return ctx.JSON(utils.Succ(`分类信息`, mods))
}

// CatePost doc
// @Tags 分类
// @Summary 分类文章列表
// @Param cid path int true "分类id" default(1)
// @Param pi query int true "分页页数pi" default(1)
// @Param ps query int true "分页大小ps" default(6)
// @Param token query string true "凭证jwt" default(jwt)
// @Router /api/cate/post/{cid} [get]
func CatePost(ctx echo.Context) error {
	cid, err := strconv.Atoi(ctx.Param("cid"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	ipt := &model.Page{}
	err = ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	count := model.CatePostCount(cid, false)
	if count == 0 {
		return ctx.JSON(utils.ErrOpt(`未查询到文章信息,请重试`))
	}
	mods, err := model.CatePostList(cid, ipt.Pi, ipt.Ps, false)
	if err != nil {
		return ctx.JSON(utils.ErrOpt(`未查询到文章信息,请重试`, err.Error()))
	}
	return ctx.JSON(utils.Page(`文章信息`, mods, count))
}

// CateAdd doc
// @Tags 分类
// @Summary 添加分类
// @Param body body model.Cate true "分类 struct"
// @Param token query string true "凭证jwt" default(jwt)
// @Router /api/cate/add [post]
func CateAdd(ctx echo.Context) error {
	ipt := &model.Cate{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	if !model.CateAdd(ipt) {
		return ctx.JSON(utils.Fail(`添加分类失败,请重试`))
	}
	return ctx.JSON(utils.Succ(`添加分类成功`))
}

// CateEdit doc
// @Tags 分类
// @Summary 修改分类
// @Param body body model.Cate true "分类 struct"
// @Param token query string true "凭证jwt" default(jwt)
// @Router /api/cate/edit [post]
func CateEdit(ctx echo.Context) error {
	ipt := &model.Cate{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	if !model.CateEdit(ipt) {
		return ctx.JSON(utils.Fail(`分类修改失败`))
	}
	return ctx.JSON(utils.Succ(`分类修改成功`))
}

// CateDrop doc
// @Tags 分类
// @Summary 删除分类
// @Param id path int true "id-分类" default(0)
// @Param token query string true "凭证jwt" default(jwt)
// @Router /api/cate/drop/{id} [get]
func CateDrop(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	if !model.CateDrop(id) {
		return ctx.JSON(utils.Fail(`分类删除失败,请重试`))
	}
	return ctx.JSON(utils.Succ(`分类删除成功`))
}
