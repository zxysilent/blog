package control

import (
	"blog/model"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/zxysilent/util"
)

// TagsView 标签列表
func TagsView(ctx echo.Context) error {
	mods, err := model.TagAll()
	if err != nil {
		return ctx.Redirect(302, "/")
	}
	return ctx.Render(http.StatusOK, "tags.html", map[string]interface{}{
		"Tags": mods,
	})
}

// TagPostView 标签文章列表
func TagPostView(ctx echo.Context) error {
	tag := ctx.Param("tag")
	if tag == "" {
		return ctx.Redirect(302, "/tags")
	}
	mod, has := model.TagName(tag)
	if !has {
		return ctx.Redirect(302, "/tags")
	}
	pi, _ := strconv.Atoi(ctx.FormValue("page"))
	if pi == 0 {
		pi = 1
	}
	ps := util.Atoi(model.MapOpts.MustGet("page_size"), 6)
	mods, err := model.TagPostList(mod.Id, pi, ps)
	if err != nil || len(mods) < 1 {
		return ctx.Redirect(302, "/tags")
	}
	total := model.TagPostCount(mod.Id)
	naver := model.Naver{}
	if pi > 1 {
		naver.Prev = "/tag/" + mod.Name + "?page=1"
	}
	if total > (pi * ps) {
		naver.Next = "/tag/" + mod.Name + "?page=" + strconv.Itoa(pi+1)
	}
	return ctx.Render(http.StatusOK, "tag-post.html", map[string]interface{}{
		"Tag":      mod,
		"TagPosts": mods,
		"Naver":    naver,
	})
}

// TagAll 所有标签
func TagAll(ctx echo.Context) error {
	mods, err := model.TagAll()
	if err != nil {
		return ctx.JSON(util.NewErrOpt(`未查询到标签信息`, err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(util.NewErrOpt(`未查询到标签信息`, "len"))
	}
	return ctx.JSON(util.NewSucc(`分类信息`, mods))
}

// TagAdd 添加标签
func TagAdd(ctx echo.Context) error {
	ipt := &model.Tag{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(util.NewErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	if !model.TagAdd(ipt) {
		return ctx.JSON(util.NewFail(`添加标签失败,请重试`))
	}
	return ctx.JSON(util.NewSucc(`添加标签成功`))
}

// TagEdit 修改标签
func TagEdit(ctx echo.Context) error {
	ipt := &model.Tag{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(util.NewErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	if !model.TagEdit(ipt) {
		return ctx.JSON(util.NewFail(`标签修改失败`))
	}
	return ctx.JSON(util.NewSucc(`标签修改成功`))
}

// TagDel  删除标签
func TagDel(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(util.NewErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	if !model.TagDel(id) {
		return ctx.JSON(util.NewFail(`标签删除失败,请重试`))
	}
	// 删除标签相关联的数据
	model.TagPostDels(id)
	return ctx.JSON(util.NewSucc(`标签删除成功`))
}
