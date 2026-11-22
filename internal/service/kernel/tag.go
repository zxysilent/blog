package kernel

import (
	"blog/internal/model"
	"blog/internal/repo"
	"blog/pkg/reply"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/logs"
)

// TagList doc
// @Tags tag
// @Summary 标签列表
// @Security KeyAuth
// @Param query query model.TagListArgs true "请求数据"
// @Param stats query int false "为 1 时附带标签文章计数"
// @Success 200 {object} reply.Reply{data=[]model.Tag} "返回数据"
// @Router /api/tags [get]
func TagList(ctx echo.Context) error {
	args := &model.TagListArgs{}
	if err := ctx.Bind(args); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	rctx := ctx.Request().Context()
	// 统计模式返回 name/count 列表
	if ctx.QueryParam("stats") == "1" {
		mods, err := repo.TagStatus()
		if err != nil {
			logs.Ctx(rctx).Err(err).Error("repo.TagStatus")
			return reply.Fail(ctx, "查询失败")
		}
		return reply.List(ctx, "ok", mods, len(mods))
	}
	// 分页参数缺省值统一由 model.List 处理
	args.Normalize(0)
	mods, count, err := repo.TagList(rctx, args)
	if err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.TagList")
		return reply.Fail(ctx, "查询失败")
	}
	return reply.List(ctx, "ok", mods, count)
}

// TagGet doc
// @Tags tag
// @Summary 标签详情
// @Security KeyAuth
// @Param id path string true "标签id"
// @Success 200 {object} reply.Reply{data=model.Tag} "返回数据"
// @Router /api/tags/{id} [get]
func TagGet(ctx echo.Context) error {
	rctx := ctx.Request().Context()
	mod, err := repo.GetAny[model.Tag](rctx, ctx.Param("id"))
	if err != nil {
		logs.Ctx(rctx).Err(err).Str("id", ctx.Param("id")).Error("repo.GetAny")
		return reply.Fail(ctx, "标签不存在")
	}
	return reply.Succ(ctx, "ok", mod)
}

// TagAdd doc
// @Tags tag
// @Summary 标签新增
// @Security KeyAuth
// @Param body body model.Tag true "请求数据"
// @Success 200 {object} reply.Reply{data=string} "返回数据"
// @Router /api/tags [post]
func TagAdd(ctx echo.Context) error {
	ipt := &model.Tag{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	ipt.Name = strings.TrimSpace(ipt.Name)
	if ipt.Name == "" {
		return reply.Fail(ctx, "标签名不能为空")
	}
	ipt.Id = ""
	milli := time.Now().UnixMilli()
	ipt.Updated = milli
	ipt.Created = milli
	rctx := ctx.Request().Context()
	if repo.TagExist(rctx, &model.Tag{Name: ipt.Name}) {
		return reply.Fail(ctx, "标签名已存在")
	}
	if err := repo.TagAdd(rctx, ipt); err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.TagAdd")
		return reply.Fail(ctx, "添加失败")
	}
	return reply.Succ(ctx, "添加成功", ipt.Id)
}

// TagEdit doc
// @Tags tag
// @Summary 标签修改
// @Security KeyAuth
// @Param id path string true "标签id"
// @Param body body model.Tag true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/tags/{id} [put]
func TagEdit(ctx echo.Context) error {
	ipt := &model.Tag{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	id := ctx.Param("id")
	rctx := ctx.Request().Context()
	old, err := repo.GetAny[model.Tag](rctx, id)
	if err != nil {
		return reply.Fail(ctx, "标签不存在")
	}
	ipt.Id = id
	ipt.Name = strings.TrimSpace(ipt.Name)
	if ipt.Name == "" {
		return reply.Fail(ctx, "标签名不能为空")
	}
	if ipt.Name != old.Name && repo.TagExist(rctx, &model.Tag{Name: ipt.Name}) {
		return reply.Fail(ctx, "标签名已存在")
	}
	ipt.Created = old.Created
	ipt.Updated = time.Now().UnixMilli()
	if err := repo.TagEdit(rctx, ipt, "name", "intro", "color", "updated"); err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.TagEdit")
		return reply.Fail(ctx, "修改失败")
	}
	return reply.Succ(ctx, "修改成功")
}

// TagDrop doc
// @Tags tag
// @Summary 标签删除
// @Security KeyAuth
// @Param id path string true "标签id"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/tags/{id} [delete]
func TagDrop(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return reply.Fail(ctx, "参数有误")
	}
	rctx := ctx.Request().Context()
	if err := repo.TagPostDrop(id); err != nil {
		logs.Ctx(rctx).Err(err).Str("id", id).Error("repo.TagPostDrop")
	}
	if err := repo.TagDrop(rctx, &model.Tag{Id: id}); err != nil {
		logs.Ctx(rctx).Err(err).Str("id", id).Error("repo.TagDrop")
		return reply.Fail(ctx, "删除失败")
	}
	return reply.Succ(ctx, "删除成功")
}
