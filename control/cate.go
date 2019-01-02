package control

import (
	"blog/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/zxysilent/util"
)

// CatePostView 标签文章列表
func CatePostView(ctx echo.Context) error {
	cate := ctx.Param("cate")
	if cate == "" {
		return ctx.Redirect(302, "/")
	}
	mod, has := model.CateName(cate)
	if !has {
		return ctx.Redirect(302, "/")
	}
	pi, _ := strconv.Atoi(ctx.FormValue("page"))
	if pi == 0 {
		pi = 1
	}
	ps := util.Atoi(model.MapOpts.MustGet("page_size"), 6)
	mods, err := model.CatePostList(mod.Id, pi, ps, true)
	if err != nil || len(mods) < 1 {
		return ctx.Redirect(302, "/")
	}
	total := model.CatePostCount(mod.Id, true)
	naver := model.Naver{}
	if pi > 1 {
		naver.Prev = "/cate/" + mod.Name + "?page=1"
	}
	if total > (pi * ps) {
		naver.Next = "/cate/" + mod.Name + "?page=" + strconv.Itoa(pi+1)
	}
	return ctx.Render(http.StatusOK, "cate-post.html", map[string]interface{}{
		"Cate":      mod,
		"CatePosts": mods,
		"Naver":     naver,
	})
}

// CateAll 所有分类
func CateAll(ctx echo.Context) error {
	mods, err := model.CateAll()
	if err != nil {
		return ctx.JSON(util.NewErrOpt(`未查询到分类信息`, err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(util.NewErrOpt(`未查询到分类信息`, "len"))
	}
	return ctx.JSON(util.NewSucc(`分类信息`, mods))
}

// CatePost 分类文章列表
// @Param cid path int true "分类id" default(1)
// @Param pi query int true "分页页数pi" default(1)
// @Param ps query int true "分页大小ps" default(6)
// @Success 200 {object} util.Result "成功数据"
func CatePost(ctx echo.Context) error {
	cid, err := strconv.Atoi(ctx.Param("cid"))
	if err != nil {
		return ctx.JSON(util.NewErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	ipt := &model.Page{}
	err = ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(util.NewErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	count := model.CatePostCount(cid, false)
	if count == 0 {
		return ctx.JSON(util.NewErrOpt(`未查询到文章信息,请重试`))
	}
	mods, err := model.CatePostList(cid, ipt.Pi, ipt.Ps, false)
	if err != nil {
		return ctx.JSON(util.NewErrOpt(`未查询到文章信息,请重试`, err.Error()))
	}
	return ctx.JSON(util.NewPage(`文章信息`, mods, count))
}

// CateAdd 添加分类
func CateAdd(ctx echo.Context) error {
	ipt := &model.Cate{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(util.NewErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	if !model.CateAdd(ipt) {
		return ctx.JSON(util.NewFail(`添加分类失败,请重试`))
	}
	return ctx.JSON(util.NewSucc(`添加分类成功`))
}

// CateEdit 修改分类
func CateEdit(ctx echo.Context) error {
	ipt := &model.Cate{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(util.NewErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	if !model.CateEdit(ipt) {
		return ctx.JSON(util.NewFail(`分类修改失败`))
	}
	return ctx.JSON(util.NewSucc(`分类修改成功`))
}

// CateDel  删除分类
func CateDel(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(util.NewErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	if !model.CateDel(id) {
		return ctx.JSON(util.NewFail(`分类删除失败,请重试`))
	}
	return ctx.JSON(util.NewSucc(`分类删除成功`))
}
