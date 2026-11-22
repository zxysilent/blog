package kernel

import (
	"blog/internal/model"
	"blog/internal/repo"
	"blog/pkg/reply"
	"fmt"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/logs"
)

// PostList doc
// @Tags post
// @Summary 内容列表
// @Security KeyAuth
// @Param query query model.PostListArgs true "请求数据"
// @Success 200 {object} reply.Reply{data=[]model.Post} "返回数据"
// @Router /api/posts [get]
func PostList(ctx echo.Context) error {
	args := &model.PostListArgs{}
	if err := ctx.Bind(args); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	postFilterExpand(&args.PostFilter)
	rctx := ctx.Request().Context()
	mods, count, err := repo.PostList(rctx, args)
	if err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.PostList")
		return reply.Fail(ctx, "查询失败")
	}
	return reply.List(ctx, "ok", mods, count)
}

// PostTree doc
// @Tags post
// @Summary 内容树（笔记目录使用）
// @Security KeyAuth
// @Param query query model.PostFilterPage true "请求数据"
// @Success 200 {object} reply.Reply{data=[]model.PostPart} "返回数据"
// @Router /api/posts/tree [get]
func PostTree(ctx echo.Context) error {
	args := &model.PostFilterPage{}
	if err := ctx.Bind(args); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	postFilterExpand(&args.PostFilter)
	mods, count, err := repo.PostTree(args)
	if err != nil {
		logs.Ctx(ctx.Request().Context()).Err(err).Error("repo.PostTree")
		return reply.Fail(ctx, "查询失败")
	}
	return reply.List(ctx, "ok", mods, count)
}

// PostGet doc
// @Tags post
// @Summary 内容详情
// @Security KeyAuth
// @Param id path string true "内容id"
// @Success 200 {object} reply.Reply{data=model.Post} "返回数据"
// @Router /api/posts/{id} [get]
func PostGet(ctx echo.Context) error {
	id := ctx.Param("id")
	rctx := ctx.Request().Context()
	mod, err := repo.PostGet(rctx, id)
	if err != nil {
		logs.Ctx(rctx).Err(err).Str("id", id).Error("repo.PostGet")
		return reply.Fail(ctx, "内容不存在")
	}
	mod.Tags = repo.PostTags(id)
	return reply.Succ(ctx, "ok", mod)
}

// postFilterExpand 展开栏目子树过滤条件
// cate_id + fuzzy=true 时把过滤条件扩展为「当前栏目及其所有子孙栏目」，
// 相当于对栏目维度做模糊匹配；fuzzy=false 时精确匹配单个栏目
func postFilterExpand(filter *model.PostFilter) {
	if !filter.Fuzzy || filter.CateId == nil || *filter.CateId == "" {
		return
	}
	filter.CateIds = repo.CateSubIds(*filter.CateId)
}

// cateUsable 校验栏目可用于挂载内容
// 栏目必须存在且 Use 为真，否则不允许保存内容
func cateUsable(cateId string) error {
	if cateId == "" {
		return nil
	}
	mod, ok := repo.CateCached(cateId)
	if !ok {
		return fmt.Errorf("栏目不存在")
	}
	if !mod.Use {
		return fmt.Errorf("栏目「%s」已禁用，不能挂载内容", mod.Name)
	}
	return nil
}

// PostAdd doc
// @Tags post
// @Summary 内容新增
// @Security KeyAuth
// @Param body body model.PostAddArgs true "请求数据"
// @Success 200 {object} reply.Reply{data=string} "返回数据"
// @Router /api/posts [post]
func PostAdd(ctx echo.Context) error {
	ipt := &model.PostAddArgs{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	if err := cateUsable(ipt.CateId); err != nil {
		return reply.Fail(ctx, err.Error())
	}
	mod := &model.Post{
		Pid:      strings.TrimSpace(ipt.Pid),
		CateId:   strings.TrimSpace(ipt.CateId),
		Kind:     ipt.Kind,
		Status:   ipt.Status,
		Title:    strings.TrimSpace(ipt.Title),
		Path:     model.PostNormalizePath(ipt.Path),
		Summary:  ipt.Summary,
		Author:   strings.TrimSpace(ipt.Author),
		Origin:   strings.TrimSpace(ipt.Origin),
		Icon:     strings.TrimSpace(ipt.Icon),
		Markdown: ipt.Markdown,
		Richtext: ipt.Richtext,
	}
	// 默认值：类型按文章、状态按草稿、默认可见
	if mod.Kind == 0 {
		mod.Kind = model.KindPost
	}
	if mod.Status == 0 {
		mod.Status = model.StatusDraft
	}
	mod.Show = true
	if ipt.Show != nil {
		mod.Show = *ipt.Show
	}
	if ipt.Istop != nil {
		mod.Istop = *ipt.Istop
	}
	if ipt.Allow != nil {
		mod.Allow = *ipt.Allow
	}
	if err := model.PostValidate(mod.Kind, mod.Status, mod.Title, mod.Path, mod.Markdown, mod.Richtext); err != nil {
		return reply.Fail(ctx, err.Error())
	}
	// 笔记允许挂在另一篇笔记下，其他类型必须平铺
	if mod.Pid != "" && mod.Kind != model.KindNote {
		return reply.Fail(ctx, "只有笔记支持父子层级")
	}
	rctx := ctx.Request().Context()
	if repo.PostExistPath(rctx, mod.Path, "") {
		return reply.Fail(ctx, "访问路径已存在")
	}
	if uid, ok := currentUid(ctx); ok {
		mod.Uid = uid
	}
	milli := time.Now().UnixMilli()
	mod.Updated = milli
	mod.Created = milli
	if err := repo.PostAdd(rctx, mod); err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.PostAdd")
		return reply.Fail(ctx, "添加失败")
	}
	return reply.Succ(ctx, "添加成功", mod.Id)
}

// PostEdit doc
// @Tags post
// @Summary 内容修改（支持局部字段）
// @Security KeyAuth
// @Param id path string true "内容id"
// @Param body body model.PostEditArgs true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/posts/{id} [put]
func PostEdit(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return reply.Fail(ctx, "参数有误")
	}
	ipt := &model.PostEditArgs{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	rctx := ctx.Request().Context()
	mod, err := repo.PostGet(rctx, id)
	if err != nil {
		return reply.Fail(ctx, "内容不存在")
	}
	// 在旧值基础上合并，便于做完整的保存前校验
	kind := mod.Kind
	status := mod.Status
	title := mod.Title
	path := mod.Path
	markdown := mod.Markdown
	richtext := mod.Richtext
	cols := make([]string, 0, 16)

	if ipt.Pid != nil {
		mod.Pid = strings.TrimSpace(*ipt.Pid)
		cols = append(cols, "pid")
	}
	if ipt.CateId != nil {
		cateId := strings.TrimSpace(*ipt.CateId)
		if err := cateUsable(cateId); err != nil {
			return reply.Fail(ctx, err.Error())
		}
		mod.CateId = cateId
		cols = append(cols, "cate_id")
	}
	if ipt.Status != nil {
		status = *ipt.Status
		mod.Status = status
		cols = append(cols, "status")
	}
	if ipt.Title != nil {
		title = strings.TrimSpace(*ipt.Title)
		mod.Title = title
		cols = append(cols, "title")
	}
	if ipt.Path != nil {
		// 内置内容的路径由服务端固定，不接受客户端修改
		if mod.Inner {
			return reply.Fail(ctx, "内置内容不可修改访问路径")
		}
		path = model.PostNormalizePath(*ipt.Path)
		mod.Path = path
		cols = append(cols, "path")
	}
	if ipt.Summary != nil {
		mod.Summary = *ipt.Summary
		cols = append(cols, "summary")
	}
	if ipt.Author != nil {
		mod.Author = strings.TrimSpace(*ipt.Author)
		cols = append(cols, "author")
	}
	if ipt.Origin != nil {
		mod.Origin = strings.TrimSpace(*ipt.Origin)
		cols = append(cols, "origin")
	}
	if ipt.Icon != nil {
		mod.Icon = strings.TrimSpace(*ipt.Icon)
		cols = append(cols, "icon")
	}
	if ipt.Markdown != nil {
		markdown = *ipt.Markdown
		mod.Markdown = markdown
		cols = append(cols, "markdown")
	}
	if ipt.Richtext != nil {
		richtext = *ipt.Richtext
		mod.Richtext = richtext
		cols = append(cols, "richtext")
	}
	if ipt.Allow != nil {
		mod.Allow = *ipt.Allow
		cols = append(cols, "allow")
	}
	if ipt.Show != nil {
		mod.Show = *ipt.Show
		cols = append(cols, "show")
	}
	if ipt.Istop != nil {
		mod.Istop = *ipt.Istop
		cols = append(cols, "istop")
	}
	if len(cols) < 1 {
		// PUT 语义幂等：提交内容与现状一致时不报错
		return reply.Succ(ctx, "修改成功")
	}
	if err := model.PostValidate(kind, status, title, path, markdown, richtext); err != nil {
		return reply.Fail(ctx, err.Error())
	}
	if repo.PostExistPath(rctx, path, id) {
		return reply.Fail(ctx, "访问路径已存在")
	}
	if mod.Pid != "" && kind != model.KindNote {
		return reply.Fail(ctx, "只有笔记支持父子层级")
	}
	mod.Updated = time.Now().UnixMilli()
	cols = append(cols, "updated")
	if err := repo.PostEdit(rctx, mod, cols...); err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.PostEdit")
		return reply.Fail(ctx, "修改失败")
	}
	return reply.Succ(ctx, "修改成功")
}

// PostEditTop doc
// @Tags post
// @Summary 修改置顶
// @Security KeyAuth
// @Param id path string true "内容id"
// @Param body body model.PostEditTopArgs true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/posts/{id}/top [put]
func PostEditTop(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return reply.Fail(ctx, "参数有误")
	}
	ipt := &model.PostEditTopArgs{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	rctx := ctx.Request().Context()
	if _, err := repo.PostGet(rctx, id); err != nil {
		return reply.Fail(ctx, "内容不存在")
	}
	if err := repo.PostEditTop(rctx, id, ipt.Istop); err != nil {
		logs.Ctx(rctx).Err(err).Str("id", id).Error("repo.PostEditTop")
		return reply.Fail(ctx, "修改失败")
	}
	return reply.Succ(ctx, "修改成功")
}

// PostDrop doc
// @Tags post
// @Summary 内容删除
// @Security KeyAuth
// @Param id path string true "内容id"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/posts/{id} [delete]
func PostDrop(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return reply.Fail(ctx, "参数有误")
	}
	rctx := ctx.Request().Context()
	mod, err := repo.PostGet(rctx, id)
	if err != nil {
		return reply.Fail(ctx, "内容不存在")
	}
	if mod.Inner {
		return reply.Fail(ctx, "内置内容不可删除")
	}
	// 有子级时先要求处理子级，避免目录树被直接挖空
	if repo.PostCount(rctx, &model.PostFilter{Pid: &id}) > 0 {
		return reply.Fail(ctx, "请先删除子级内容")
	}
	if err := repo.PostDrop(rctx, id); err != nil {
		logs.Ctx(rctx).Err(err).Str("id", id).Error("repo.PostDrop")
		return reply.Fail(ctx, "删除失败")
	}
	return reply.Succ(ctx, "删除成功")
}

// PostTags doc
// @Tags post
// @Summary 获取内容标签
// @Security KeyAuth
// @Param id path string true "内容id"
// @Success 200 {object} reply.Reply{data=[]model.Tag} "返回数据"
// @Router /api/posts/{id}/tags [get]
func PostTags(ctx echo.Context) error {
	return reply.Succ(ctx, "ok", repo.PostTags(ctx.Param("id")))
}

// PostTagEdit doc
// @Tags post
// @Summary 全量替换内容标签
// @Security KeyAuth
// @Param id path string true "内容id"
// @Param body body object{tag_ids=[]string} true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/posts/{id}/tags [put]
func PostTagEdit(ctx echo.Context) error {
	ipt := &struct {
		TagIds []string `json:"tag_ids" form:"tag_ids"`
	}{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	id := ctx.Param("id")
	rctx := ctx.Request().Context()
	if _, err := repo.PostGet(rctx, id); err != nil {
		return reply.Fail(ctx, "内容不存在")
	}
	if len(ipt.TagIds) > 0 {
		if exist := repo.TagMapIds(rctx, ipt.TagIds); len(exist) != len(ipt.TagIds) {
			return reply.Fail(ctx, "存在无效的标签id")
		}
	}
	if err := repo.PostTagSync(id, ipt.TagIds); err != nil {
		logs.Ctx(rctx).Err(err).Str("id", id).Error("repo.PostTagSync")
		return reply.Fail(ctx, "修改失败")
	}
	return reply.Succ(ctx, "修改成功")
}

// PostNeighbors doc
// @Tags post
// @Summary 获取内容上下篇
// @Security KeyAuth
// @Param id path string true "内容id"
// @Success 200 {object} reply.Reply{data=model.PostNeighbor} "返回数据"
// @Router /api/posts/{id}/neighbors [get]
func PostNeighbors(ctx echo.Context) error {
	id := ctx.Param("id")
	rctx := ctx.Request().Context()
	mod, err := repo.PostGet(rctx, id)
	if err != nil {
		return reply.Fail(ctx, "内容不存在")
	}
	return reply.Succ(ctx, "ok", repo.PostNeighbor(rctx, mod.CateId, mod.Created))
}
