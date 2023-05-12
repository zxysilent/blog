package applet

import (
	"blog/internal/model"
	"blog/internal/repo"
	"blog/internal/utils"
	"time"

	"github.com/labstack/echo/v4"
)

// TagGet doc
// @Auth
// @Tags tag
// @Summary 标签单条数据
// @Param id query int true "id"
// @Success 200 {object} utils.Reply{data=model.Tag} "返回数据"
// @Router /api/tag/get [get]
func TagGet(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod := &model.Tag{Id: ipt.Id}
	err = repo.TagGet(mod)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到数据", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// TagList doc
// @Auth
// @Tags tag
// @Summary 标签列表数据
// @Param query query model.TagFilterList true "请求数据"
// @Success 200 {object} utils.Reply{data=[]model.Tag} "返回数据"
// @Router /api/tag/list [get]
func TagList(ctx echo.Context) error {
	ipt := &model.TagFilterList{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mods, err := repo.TagList(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到数据", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// TagPage doc
// @Auth
// @Tags tag
// @Summary 获取标签分页
// @Param query query model.TagFilterPage true "请求数据"
// @Success 200 {object} utils.Reply{data=[]model.Tag} "返回数据"
// @Router /api/tag/page [get]
func TagPage(ctx echo.Context) error {
	// cid, err := strconv.Atoi(ctx.Param("cid"))
	// if err != nil {
	//  return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	// }
	ipt := &model.TagFilterPage{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if err = ipt.Stat(); err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mods, count, err := repo.TagPage(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if count == 0 || len(mods) == 0 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "empty"))
	}
	return ctx.JSON(utils.Page("succ", mods, count))
}

// TagAdd doc
// @Auth
// @Tags tag
// @Summary 标签添加数据
// @Param token query string true "token"
// @Param body body model.Tag true "请求数据"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/tag/add [post]
func TagAdd(ctx echo.Context) error {
	ipt := &model.Tag{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	ipt.Updated = time.Now().UnixMilli()
	ipt.Created = ipt.Updated
	err = repo.TagAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// TagEdit doc
// @Auth
// @Tags tag
// @Summary 标签修改数据
// @Param token query string true "token"
// @Param body body model.Tag true "请求数据"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/tag/edit [post]
func TagEdit(ctx echo.Context) error {
	ipt := &model.Tag{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	ipt.Updated = time.Now().UnixMilli()
	err = repo.TagEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// TagDrop doc
// @Auth
// @Tags tag
// @Summary 标签删除数据
// @Param token query string true "token"
// @Param query query model.IptId true "请求数据"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/tag/drop [post]
func TagDrop(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	err = repo.TagDrop(&model.Tag{Id: ipt.Id})
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
