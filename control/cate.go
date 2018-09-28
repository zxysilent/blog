package control

import (
	"blog/model"
	"blog/util"

	"github.com/labstack/echo"
)

// CateAll 所有分类
func CateAll(ctx echo.Context) error {
	mods, err := model.CateAll()
	if err != nil {
		return ctx.Res(util.NewErrOpt(`未查询到分类信息`, err.Error()))
	}
	if len(mods) < 1 {
		return ctx.Res(util.NewErrOpt(`未查询到分类信息`, "len"))
	}
	return ctx.Res(util.NewSucc(`分类信息`, mods))
}

// CatePost 分类文章列表
// @Param cid path int true "分类id" default(1)
// @Param pi query int true "分页页数pi" default(1)
// @Param ps query int true "分页大小ps" default(6)
// @Success 200 {object} util.Result "成功数据"
func CatePost(ctx echo.Context) error {
	cid := util.Atoi(ctx.Param("cid"))
	ipt := &model.Page{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.Res(util.NewErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	count := model.CatePostCount(cid)
	if count == 0 {
		return ctx.Res(util.NewErrOpt(`未查询到文章信息,请重试`))
	}
	mods, err := model.CatePostList(cid, ipt.Pi, ipt.Ps)
	if err != nil {
		return ctx.Res(util.NewErrOpt(`未查询到文章信息,请重试`, err.Error()))
	}
	return ctx.Res(util.NewPage(`文章信息`, mods, count))
}
