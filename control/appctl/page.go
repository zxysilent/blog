package appctl

import (
	"blog/model"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// PageGet doc
// @Tags page-页面
// @Summary 通过id获取单条页面
// @Param id query int true "id"
// @Success 200 {object} model.Reply{data=model.Page} "返回数据"
// @Router /api/page/get [get]
func PageGet(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod, has := model.PostGet(ipt.Id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到页面"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// PostPage doc
// @Tags post-文章页面
// @Summary 获取文章分页
// @Param pi query int true "分页数" default(1)
// @Param ps query int true "每页条数" default(5)
// @Success 200 {object} model.Reply{data=[]model.Post} "返回数据"
// @Router /api/post/page [get]
func PagePage(ctx echo.Context) error {
	ipt := &model.Page{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if err = ipt.Stat(); err != nil {
		return ctx.JSON(utils.ErrIpt("分页大小输入错误", err.Error()))
	}
	count := model.PostCount(-1, model.PostKindPage)
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := model.PostPage(-1, model.PostKindPage, ipt.Pi, ipt.Ps, "id", "title", "path", "created", "summary", "updated", "status")
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, int(count)))
}

// PageAdd doc
// @Tags page-页面
// @Summary 添加页面
// @Param token query string true "token"
// @Param body body model.Page true "请求数据"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/page/add [post]
func PageAdd(ctx echo.Context) error {
	ipt := &model.Post{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if model.PostExist(ipt.Path) {
		return ctx.JSON(utils.ErrIpt("当前访问路径已经存在,请重新输入"))
	}
	if strings.Contains(ipt.Richtext, "<!--more-->") {
		ipt.Summary = strings.Split(ipt.Richtext, "<!--more-->")[0]
	}
	ipt.Updated = ipt.Created
	err = model.PostAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// PageEdit doc
// @Tags page-页面
// @Summary 修改页面
// @Param token query string true "token"
// @Param body body model.Page true "请求数据"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/page/edit [post]
func PageEdit(ctx echo.Context) error {
	ipt := &model.Post{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if strings.Contains(ipt.Richtext, "<!--more-->") {
		ipt.Summary = strings.Split(ipt.Richtext, "<!--more-->")[0]
	}
	ipt.Updated = time.Now()
	err = model.PostEdit(ipt, "cate_id", "kind", "status", "title", "path", "summary", "markdown", "richtext", "allow", "created", "updated")
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// PageDrop doc
// @Tags page-页面
// @Summary 通过id删除单条页面
// @Param id query int true "id"
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/page/drop [post]
func PageDrop(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	err = model.PostDrop(ipt.Id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
