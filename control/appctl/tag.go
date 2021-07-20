package appctl

import (
	"blog/model"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// TagAll doc
// @Tags tag-标签
// @Summary 获取所有标签信息
// @Success 200 {object} model.Reply{data=[]model.Tag} "返回数据"
// @Router /api/tag/all [get]
func TagAll(ctx echo.Context) error {
	mods, err := model.TagAll()
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到标签信息", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到标签信息", "len"))
	}
	return ctx.JSON(utils.Succ("分类信息", mods))
}

// TagAdd doc
// @Tags tag-标签
// @Summary 添加标签信息
// @Param token query string true "token"
// @Param body body model.Tag true "请求数据"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/tag/add [post]
func TagAdd(ctx echo.Context) error {
	ipt := &model.Tag{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误,请重试", err.Error()))
	}
	err = model.TagAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("添加标签成功"))
}

// TagEdit doc
// @Tags tag-标签
// @Summary 修改标签信息
// @Param token query string true "token"
// @Param body body model.Tag true "请求数据"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/tag/edit [post]
func TagEdit(ctx echo.Context) error {
	ipt := &model.Tag{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误,请重试", err.Error()))
	}
	err = model.TagEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("标签修改成功"))
}

// TagDrop doc
// @Tags tag-标签
// @Summary 通过id删除单条标签信息
// @Param token query string true "token"
// @Param body body model.IptId true "json"
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
	// 删除标签相关联的数据
	model.TagPostDrop(ipt.Id)
	return ctx.JSON(utils.Succ("标签删除成功"))
}
