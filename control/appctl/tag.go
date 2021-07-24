package appctl

import (
	"blog/model"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// TagGet doc
// @Tags tag-标签
// @Summary 通过id获取单条标签
// @Param id query int true "id"
// @Success 200 {object} model.Reply{data=model.Tag} "返回数据"
// @Router /api/tag/get [get]
func TagGet(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod, has := model.TagGet(ipt.Id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到标签"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// TagAll doc
// @Tags tag-标签
// @Summary 获取所有标签
// @Success 200 {object} model.Reply{data=[]model.Tag} "返回数据"
// @Router /api/tag/all [get]
func TagAll(ctx echo.Context) error {
	mods, err := model.TagAll()
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到标签", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// TagPage doc
// @Tags tag-标签
// @Summary 获取标签分页
// @Param cid path int true "分类id" default(1)
// @Param pi query int true "分页数" default(1)
// @Param ps query int true "每页条数[5,30]" default(5)
// @Success 200 {object} model.Reply{data=[]model.Tag} "返回数据"
// @Router /api/tag/page [get]
func TagPage(ctx echo.Context) error {
	ipt := &model.Page{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if ipt.Ps > 30 || ipt.Ps < 1 {
		return ctx.JSON(utils.ErrIpt("分页大小输入错误", ipt.Ps))
	}
	count := model.TagCount()
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := model.TagPage(ipt.Pi, ipt.Ps)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, int(count)))
}

// TagAdd doc
// @Tags tag-标签
// @Summary 添加标签
// @Param token query string true "token"
// @Param body body model.Tag true "请求数据"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/tag/add [post]
func TagAdd(ctx echo.Context) error {
	ipt := &model.Tag{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	err = model.TagAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// TagEdit doc
// @Tags tag-标签
// @Summary 修改标签
// @Param token query string true "token"
// @Param body body model.Tag true "请求数据"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/tag/edit [post]
func TagEdit(ctx echo.Context) error {
	ipt := &model.Tag{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	err = model.TagEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// TagDrop doc
// @Tags tag-标签
// @Summary 通过id删除单条标签
// @Param id query int true "id"
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/tag/drop [post]
func TagDrop(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	err = model.TagDrop(ipt.Id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
