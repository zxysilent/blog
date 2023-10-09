package applet

import (
	"blog/conf"
	"blog/internal/model"
	"blog/internal/repo"
	"blog/internal/utils"
	"blog/pkg/token"
	"time"

	"github.com/labstack/echo/v4"
)

// NoteSave doc
// @Auth
// @Tags note
// @Summary 笔记保存
// @Param token query string true "token"
// @Param body body model.Post true "请求数据"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/note/save [post]
func NoteSave(ctx echo.Context) error {
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

// NoteShare doc
// @Auth
// @Tags note
// @Summary 笔记分享
// @Param token query string true "token"
// @Param body body model.NoteShareArgs true "请求数据"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/note/share [post]
func NoteShare(ctx echo.Context) error {
	ipt := &model.NoteShareArgs{}
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

// NoteList doc
// @Tags note
// @Summary 笔记目录&列表数据
// @Param token query string true "token"
// @Param query query model.NoteFilterList true "请求数据"
// @Success 200 {object} utils.Reply{data=[]model.NoteReply} "返回数据"
// @Router /api/note/list [get]
func NoteList(ctx echo.Context) error {
	ipt := &model.NoteFilterList{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod := model.NoteReply{}
	if !ipt.Newest {
		cateFilter := &model.CateFilterList{}
		cateFilter.Kind = utils.Any2Ptr[int](model.KindNote)
		cateFilter.Mult = ipt.Mult
		cateFilter.Pid = ipt.CateId
		mod.Cates, _ = repo.CateList(cateFilter)
	}
	postFilter := &model.PostFilterList{}
	if ipt.Newest {
		postFilter.Limit = 20
	}
	mod.Posts, _ = repo.NoteList(ipt)
	mod.Empty = len(mod.Cates) == 0 && len(mod.Posts) == 0
	return ctx.JSON(utils.Succ("succ", mod))
}
