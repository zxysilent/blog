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

// DictList doc
// @Tags dict
// @Summary 字典列表
// @Security KeyAuth
// @Param query query model.DictListArgs true "请求数据"
// @Success 200 {object} reply.Reply{data=[]model.Dict} "返回数据"
// @Router /api/dicts [get]
func DictList(ctx echo.Context) error {
	args := &model.DictListArgs{}
	if err := ctx.Bind(args); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	rctx := ctx.Request().Context()
	mods, count, err := repo.DictList(rctx, args)
	if err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.DictList")
		return reply.Fail(ctx, "查询失败")
	}
	return reply.List(ctx, "ok", mods, count)
}

// DictGet doc
// @Tags dict
// @Summary 字典详情（含选项树）
// @Security KeyAuth
// @Param id path string true "字典主键"
// @Success 200 {object} reply.Reply{data=model.Dict} "返回数据"
// @Router /api/dicts/{id} [get]
func DictGet(ctx echo.Context) error {
	rctx := ctx.Request().Context()
	mod, err := repo.DictGet(rctx, ctx.Param("id"))
	if err != nil {
		logs.Ctx(rctx).Err(err).Str("id", ctx.Param("id")).Error("repo.DictGet")
		return reply.Fail(ctx, "字典不存在")
	}
	return reply.Succ(ctx, "ok", mod)
}

// DictEntries doc
// @Tags dict
// @Summary 字典选项（匿名可读，已过滤禁用项）
// @Param id path string true "字典主键"
// @Success 200 {object} reply.Reply{data=[]model.DictEntry} "返回数据"
// @Router /api/dicts/{id}/entries [get]
func DictEntries(ctx echo.Context) error {
	rctx := ctx.Request().Context()
	mod, err := repo.DictGet(rctx, ctx.Param("id"))
	if err != nil {
		// 字典不存在时返回空数组，避免前端额外处理
		return reply.Succ(ctx, "ok", []model.DictEntry{})
	}
	return reply.Succ(ctx, "ok", mod.Flatten())
}

// DictAdd doc
// @Tags dict
// @Summary 字典新增（含选项，id 可自定义，留空自动生成）
// @Security KeyAuth
// @Param body body model.DictAddArgs true "请求数据"
// @Success 200 {object} reply.Reply{data=string} "返回数据"
// @Router /api/dicts [post]
func DictAdd(ctx echo.Context) error {
	ipt := &model.DictAddArgs{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	milli := time.Now().UnixMilli()
	mod := &model.Dict{
		Id:      strings.TrimSpace(ipt.Id),
		Name:    ipt.Name,
		Type:    ipt.Type,
		Status:  true,
		Remark:  ipt.Remark,
		Updated: milli,
		Created: milli,
	}
	if ipt.Status != nil {
		mod.Status = *ipt.Status
	}
	// 未指定主键时生成 ULID，再校验格式
	mod.BeforeInsert()
	if err := mod.Valid(); err != nil {
		return reply.Fail(ctx, err.Error())
	}
	rctx := ctx.Request().Context()
	if repo.DictExist(rctx, mod.Id) {
		return reply.Fail(ctx, "字典已存在")
	}
	rows, err := model.FlatDictEntries(mod.Id, mod.Type, ipt.Entries)
	if err != nil {
		return reply.Fail(ctx, err.Error())
	}
	for idx := range rows {
		rows[idx].Created = milli
		rows[idx].Updated = milli
	}
	if err := repo.DictAdd(rctx, mod, rows); err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.DictAdd")
		return reply.Fail(ctx, "添加失败")
	}
	return reply.Succ(ctx, "添加成功", mod.Id)
}

// DictEdit doc
// @Tags dict
// @Summary 字典修改（提交 entries 时全量替换选项）
// @Security KeyAuth
// @Param id path string true "字典主键"
// @Param body body model.DictEditArgs true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/dicts/{id} [put]
func DictEdit(ctx echo.Context) error {
	ipt := &model.DictEditArgs{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	rctx := ctx.Request().Context()
	mod, err := repo.DictGet(rctx, ctx.Param("id"))
	if err != nil {
		return reply.Fail(ctx, "字典不存在")
	}
	cols := make([]string, 0, 6)
	if ipt.Name != nil {
		mod.Name = *ipt.Name
		cols = append(cols, "name")
	}
	if ipt.Type != nil {
		mod.Type = *ipt.Type
		cols = append(cols, "type")
	}
	if ipt.Status != nil {
		mod.Status = *ipt.Status
		cols = append(cols, "status")
	}
	if ipt.Remark != nil {
		mod.Remark = *ipt.Remark
		cols = append(cols, "remark")
	}
	if err := mod.Valid(); err != nil {
		return reply.Fail(ctx, err.Error())
	}
	// rows 为 nil 表示不动选项；非 nil（含空切片）表示全量替换
	var rows []model.DictEntry
	if ipt.Entries != nil {
		rows, err = model.FlatDictEntries(mod.Id, mod.Type, *ipt.Entries)
		if err != nil {
			return reply.Fail(ctx, err.Error())
		}
		milli := time.Now().UnixMilli()
		for idx := range rows {
			rows[idx].Created = milli
			rows[idx].Updated = milli
		}
	} else if ipt.Type != nil {
		// 仅改类型时，需确保已有选项仍符合新类型
		if _, err := model.FlatDictEntries(mod.Id, mod.Type, mod.Entries); err != nil {
			return reply.Fail(ctx, "已有选项与新类型不兼容："+err.Error())
		}
	}
	if len(cols) < 1 && rows == nil {
		return reply.Fail(ctx, "没有需要修改的数据")
	}
	mod.Updated = time.Now().UnixMilli()
	cols = append(cols, "updated")
	if err := repo.DictEdit(rctx, mod, rows, cols...); err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.DictEdit")
		return reply.Fail(ctx, "修改失败")
	}
	return reply.Succ(ctx, "修改成功")
}

// DictDrop doc
// @Tags dict
// @Summary 字典删除（级联删除选项）
// @Security KeyAuth
// @Param id path string true "字典主键"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/dicts/{id} [delete]
func DictDrop(ctx echo.Context) error {
	id := strings.TrimSpace(ctx.Param("id"))
	if id == "" {
		return reply.Fail(ctx, "参数有误")
	}
	rctx := ctx.Request().Context()
	mod, err := repo.DictGet(rctx, id)
	if err != nil {
		return reply.Fail(ctx, "字典不存在")
	}
	if mod.Inner {
		return reply.Fail(ctx, "内置字典不可删除")
	}
	if err := repo.DictDrop(rctx, mod.Id); err != nil {
		logs.Ctx(rctx).Err(err).Str("id", id).Error("repo.DictDrop")
		return reply.Fail(ctx, "删除失败")
	}
	return reply.Succ(ctx, "删除成功")
}
