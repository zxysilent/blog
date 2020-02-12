package control

import (
	"blog/model"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// TagAll 所有标签
func TagAll(ctx echo.Context) error {
	mods, err := model.TagAll()
	if err != nil {
		return ctx.JSON(utils.ErrOpt(`未查询到标签信息`, err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt(`未查询到标签信息`, "len"))
	}
	return ctx.JSON(utils.Succ(`分类信息`, mods))
}

// TagAdd 添加标签
func TagAdd(ctx echo.Context) error {
	ipt := &model.Tag{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	if !model.TagAdd(ipt) {
		return ctx.JSON(utils.Fail(`添加标签失败,请重试`))
	}
	return ctx.JSON(utils.Succ(`添加标签成功`))
}

// TagEdit 修改标签
func TagEdit(ctx echo.Context) error {
	ipt := &model.Tag{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	if !model.TagEdit(ipt) {
		return ctx.JSON(utils.Fail(`标签修改失败`))
	}
	return ctx.JSON(utils.Succ(`标签修改成功`))
}

// TagDrop  删除标签
func TagDrop(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	if !model.TagDrop(id) {
		return ctx.JSON(utils.Fail(`标签删除失败,请重试`))
	}
	// 删除标签相关联的数据
	model.TagPostDrop(id)
	return ctx.JSON(utils.Succ(`标签删除成功`))
}
