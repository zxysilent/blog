package applet

import (
	"blog/conf"
	"blog/internal/model"
	"blog/internal/repo"
	"blog/internal/utils"
	"blog/pkg/token"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// PostGet doc
// @Auth
// @Tags post
// @Summary 博文单条数据
// @Param id query int true "id"
// @Success 200 {object} utils.Reply{data=model.Post} "返回数据"
// @Router /api/post/get [get]
func PostGet(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod := &model.Post{Id: ipt.Id}
	err = repo.PostGet(mod)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到数据", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// PostList doc
// @Auth
// @Tags post
// @Summary 博文列表数据
// @Param query query model.PostFilterList true "请求数据"
// @Success 200 {object} utils.Reply{data=[]model.Post} "返回数据"
// @Router /api/post/list [get]
func PostList(ctx echo.Context) error {
	ipt := &model.PostFilterList{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mods, err := repo.PostList(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到数据", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// PostPage doc
// @Auth
// @Tags post
// @Summary 获取博文分页
// @Param query query model.PostFilterPage true "请求数据"
// @Success 200 {object} utils.Reply{data=[]model.Post} "返回数据"
// @Router /api/post/page [get]
func PostPage(ctx echo.Context) error {
	// cid, err := strconv.Atoi(ctx.Param("cid"))
	// if err != nil {
	//  return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	// }
	ipt := &model.PostFilterPage{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if err = ipt.Stat(); err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mods, count, err := repo.PostPage(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if count == 0 || len(mods) == 0 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "empty"))
	}
	return ctx.JSON(utils.Page("succ", mods, count))
}

// PostAdd doc
// @Auth
// @Tags post
// @Summary 博文添加数据
// @Param token query string true "token"
// @Param body body model.Post true "请求数据"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/post/add [post]
func PostAdd(ctx echo.Context) error {
	ipt := &model.Post{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if ipt.Kind == model.KindNote {
		ipt.Path = utils.UUID()
	}
	if repo.PostExist(&model.Post{Path: ipt.Path}) {
		return ctx.JSON(utils.ErrIpt("当前访问路径已经存在,请重新输入"))
	}
	if strings.Contains(ipt.Richtext, "<!--more-->") {
		ipt.Summary = strings.Split(ipt.Richtext, "<!--more-->")[0]
	}
	ipt.Richtext = utils.GenToc(ipt.Richtext)
	ipt.Updated = time.Now().UnixMilli()
	ipt.Created = ipt.Updated
	err = repo.PostAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// PostAdd doc
// @Auth
// @Tags post
// @Summary 博文笔记类型保存
// @Param token query string true "token"
// @Param body body model.Post true "请求数据"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/post/save [post]
func PostSave(ctx echo.Context) error {
	ipt := &struct {
		model.Post
		Tags []int `json:"tags"`
	}{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if ipt.Kind == model.KindNote && ipt.Path == "" {
		ipt.Path = utils.UUID()
	}
	if ipt.Id == 0 {
		ipt.Updated = time.Now().UnixMilli()
		ipt.Created = ipt.Updated
		err = repo.PostAdd(&ipt.Post)
		if err != nil {
			return ctx.JSON(utils.Fail("添加失败", err.Error()))
		}
		return ctx.JSON(utils.Succ("succ"))
	}
	ipt.Updated = time.Now().UnixMilli()
	err = repo.PostEdit(&ipt.Post)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// PostShare doc
// @Auth
// @Tags post
// @Summary 博文分享
// @Param token query string true "token"
// @Param body body model.PostShareArgs true "请求数据"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/post/share [post]
func PostShare(ctx echo.Context) error {
	ipt := &model.PostShareArgs{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	s := token.Token[token.Share]{Claim: token.Share{PostId: ipt.PostId}}
	if ipt.Day > 0 {
		s.ExpAt = time.Now().AddDate(0, 0, ipt.Day).Unix()
	}
	return ctx.JSON(utils.Succ("succ", s.Encode(conf.App.TokenSecret)))
}

// PostEdit doc
// @Auth
// @Tags post
// @Summary 博文修改数据
// @Param token query string true "token"
// @Param body body model.Post true "请求数据"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/post/edit [post]
func PostEdit(ctx echo.Context) error {
	ipt := &model.Post{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if strings.Contains(ipt.Richtext, "<!--more-->") {
		ipt.Summary = strings.Split(ipt.Richtext, "<!--more-->")[0]
	}
	ipt.Richtext = utils.GenToc(ipt.Richtext)
	ipt.Updated = time.Now().UnixMilli()
	err = repo.PostEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// PostDrop doc
// @Auth
// @Tags post
// @Summary 博文删除数据
// @Param token query string true "token"
// @Param query query model.IptId true "请求数据"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/post/drop [post]
func PostDrop(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	err = repo.PostDrop(&model.Post{Id: ipt.Id})
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
