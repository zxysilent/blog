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

// CateList doc
// @Tags cate
// @Summary 栏目列表
// @Security KeyAuth
// @Param query query model.CateListArgs true "请求数据"
// @Success 200 {object} reply.Reply{data=[]model.Cate} "返回数据"
// @Router /api/cates [get]
func CateList(ctx echo.Context) error {
	filter := &model.CateListArgs{}
	if err := ctx.Bind(filter); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	rctx := ctx.Request().Context()
	mods, count, err := repo.CateList(rctx, filter)
	if err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.CateList")
		return reply.Fail(ctx, "查询失败")
	}
	return reply.List(ctx, "ok", mods, count)
}

// CateTree doc
// @Tags cate
// @Summary 栏目树
// @Security KeyAuth
// @Param query query model.CateTreeArgs true "请求数据"
// @Success 200 {object} reply.Reply{data=[]model.Cate} "返回数据"
// @Router /api/cates/tree [get]
func CateTree(ctx echo.Context) error {
	args := &model.CateTreeArgs{}
	if err := ctx.Bind(args); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	// 直接从内存快照取完整树，再按条件过滤
	mods := repo.CateTree()
	mods = cateTreeFilter(mods, args)
	if args.Flat {
		flat := make([]model.Cate, 0, len(mods))
		flat = cateTreeFlatten(flat, mods)
		mods = flat
	}
	if args.Root {
		mods = append([]model.Cate{{Name: "全部栏目", Use: true, Show: true, Inner: true}}, mods...)
	}
	return reply.Succ(ctx, "ok", mods)
}

// cateTreeFilter 递归过滤栏目树：子级全部被过滤时父级也去掉
func cateTreeFilter(mods []model.Cate, args *model.CateTreeArgs) []model.Cate {
	out := make([]model.Cate, 0, len(mods))
	for idx := range mods {
		it := mods[idx]
		if args.Kind != nil && *args.Kind > 0 && it.Kind != *args.Kind {
			continue
		}
		if args.Use != nil && it.Use != *args.Use {
			continue
		}
		if args.Show != nil && it.Show != *args.Show {
			continue
		}
		if args.Pid != nil && it.Pid != *args.Pid {
			continue
		}
		it.Children = cateTreeFilter(it.Children, args)
		out = append(out, it)
	}
	return out
}

// cateTreeFlatten 展平栏目树，保留 level 便于前端下拉缩进
func cateTreeFlatten(dst []model.Cate, mods []model.Cate) []model.Cate {
	for idx := range mods {
		it := mods[idx]
		children := it.Children
		it.Children = nil
		dst = append(dst, it)
		dst = cateTreeFlatten(dst, children)
	}
	return dst
}

// CateGet doc
// @Tags cate
// @Summary 栏目详情
// @Security KeyAuth
// @Param id path string true "栏目id"
// @Success 200 {object} reply.Reply{data=model.Cate} "返回数据"
// @Router /api/cates/{id} [get]
func CateGet(ctx echo.Context) error {
	rctx := ctx.Request().Context()
	mod, err := repo.CateGet(rctx, ctx.Param("id"))
	if err != nil {
		logs.Ctx(rctx).Err(err).Str("id", ctx.Param("id")).Error("repo.CateGet")
		return reply.Fail(ctx, "栏目不存在")
	}
	mod.Children = repo.CateChildren(mod.Id)
	mod.Count = repo.CatePostCount(rctx, repo.CateSubIds(mod.Id))
	return reply.Succ(ctx, "ok", mod)
}

// cateLevelOf 计算目标父级下的层级
func cateLevelOf(pid string) (int, error) {
	if pid == "" {
		return 1, nil
	}
	parent, ok := repo.CateCached(pid)
	if !ok {
		return 0, fmt.Errorf("父级栏目不存在")
	}
	if parent.Level >= model.CateMaxLevel {
		return 0, fmt.Errorf("栏目最多 %d 级，不能在「%s」下继续创建", model.CateMaxLevel, parent.Name)
	}
	return parent.Level + 1, nil
}

// cateSelfDepth 以 id 为根的子树高度，自身计 1
func cateSelfDepth(id string) int {
	children := repo.CateChildren(id)
	depth := 1
	for idx := range children {
		if d := cateSelfDepth(children[idx].Id) + 1; d > depth {
			depth = d
		}
	}
	return depth
}

// cateContains 判断 id 是否位于 rootId 的子树中（含自身）
func cateContains(rootId, id string) bool {
	if rootId == id {
		return true
	}
	for _, it := range repo.CateChildren(rootId) {
		if cateContains(it.Id, id) {
			return true
		}
	}
	return false
}

// catePathCheck 校验公开路径：格式合法且全站唯一
func catePathCheck(ctx echo.Context, path, excludeId string) (string, error) {
	path = model.CateNormalizePath(path)
	if !model.CatePathValid(path) {
		return "", fmt.Errorf("公开路径只能包含小写字母、数字、下划线和短横线，且不能以符号开头")
	}
	if repo.CateExistPath(ctx.Request().Context(), path, excludeId) {
		return "", fmt.Errorf("公开路径已被占用")
	}
	return path, nil
}

// CateAdd doc
// @Tags cate
// @Summary 栏目新增
// @Security KeyAuth
// @Param body body model.CateAddArgs true "请求数据"
// @Success 200 {object} reply.Reply{data=string} "返回数据"
// @Router /api/cates [post]
func CateAdd(ctx echo.Context) error {
	ipt := &model.CateAddArgs{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	ipt.Name = strings.TrimSpace(ipt.Name)
	if ipt.Name == "" || len(ipt.Name) > 255 {
		return reply.Fail(ctx, "栏目名称不能为空且不超过 255 字符")
	}
	level, err := cateLevelOf(ipt.Pid)
	if err != nil {
		return reply.Fail(ctx, err.Error())
	}
	rctx := ctx.Request().Context()
	if repo.CateExistName(rctx, ipt.Pid, ipt.Name, "") {
		return reply.Fail(ctx, "同级下栏目名称已存在")
	}
	if ipt.Kind < 1 {
		ipt.Kind = model.KindPost
	}
	mod := &model.Cate{
		Pid:   ipt.Pid,
		Kind:  ipt.Kind,
		Name:  ipt.Name,
		Url:   strings.TrimSpace(ipt.Url),
		Intro: strings.TrimSpace(ipt.Intro),
		Icon:  strings.TrimSpace(ipt.Icon),
		Color: strings.TrimSpace(ipt.Color),
		Level: level,
		Itpl:  strings.TrimSpace(ipt.Itpl),
		Ltpl:  strings.TrimSpace(ipt.Ltpl),
		Dtpl:  strings.TrimSpace(ipt.Dtpl),
		Use:   true,
		Show:  true,
		Sort:  1000,
	}
	if ipt.Sort > 0 {
		mod.Sort = ipt.Sort
	}
	if ipt.Use != nil {
		mod.Use = *ipt.Use
	}
	if ipt.Show != nil {
		mod.Show = *ipt.Show
	}
	if ipt.Auth != nil {
		mod.Auth = *ipt.Auth
	}
	// 主键先生成，公开路径留空时用主键小写作为默认值（保证合法且唯一）
	mod.BeforeInsert()
	path := strings.TrimSpace(ipt.Path)
	if path == "" {
		mod.Path = model.CateNormalizePath(mod.Id)
	} else {
		if mod.Path, err = catePathCheck(ctx, path, ""); err != nil {
			return reply.Fail(ctx, err.Error())
		}
	}
	milli := time.Now().UnixMilli()
	mod.Updated = milli
	mod.Created = milli
	if err := repo.CateAdd(rctx, mod); err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.CateAdd")
		return reply.Fail(ctx, "添加失败")
	}
	return reply.Succ(ctx, "添加成功", mod.Id)
}

// CateEdit doc
// @Tags cate
// @Summary 栏目修改（支持局部字段）
// @Security KeyAuth
// @Param id path string true "栏目id"
// @Param body body model.CateEditArgs true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/cates/{id} [put]
func CateEdit(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return reply.Fail(ctx, "参数有误")
	}
	ipt := &model.CateEditArgs{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	rctx := ctx.Request().Context()
	mod, err := repo.CateGet(rctx, id)
	if err != nil {
		return reply.Fail(ctx, "栏目不存在")
	}
	cols := make([]string, 0, 16)
	pid := mod.Pid
	if ipt.Pid != nil {
		pid = strings.TrimSpace(*ipt.Pid)
		if pid != mod.Pid {
			if pid != "" {
				// 不能移动到自身或自己的子孙下，否则形成环
				if cateContains(id, pid) {
					return reply.Fail(ctx, "不能移动到自身或其子级下")
				}
				if _, ok := repo.CateCached(pid); !ok {
					return reply.Fail(ctx, "父级栏目不存在")
				}
			}
			level, err := cateLevelOf(pid)
			if err != nil {
				return reply.Fail(ctx, err.Error())
			}
			// 整棵子树一起移动，必须保证子树最低层也不超上限
			if level+cateSelfDepth(id)-1 > model.CateMaxLevel {
				return reply.Fail(ctx, fmt.Sprintf("移动后将超过 %d 级层级上限", model.CateMaxLevel))
			}
			mod.Pid = pid
			cols = append(cols, "pid")
		}
	}
	if ipt.Name != nil {
		name := strings.TrimSpace(*ipt.Name)
		if name == "" || len(name) > 255 {
			return reply.Fail(ctx, "栏目名称不能为空且不超过 255 字符")
		}
		if name != mod.Name && repo.CateExistName(rctx, pid, name, id) {
			return reply.Fail(ctx, "同级下栏目名称已存在")
		}
		mod.Name = name
		cols = append(cols, "name")
	}
	if ipt.Path != nil {
		path, err := catePathCheck(ctx, *ipt.Path, id)
		if err != nil {
			return reply.Fail(ctx, err.Error())
		}
		mod.Path = path
		cols = append(cols, "path")
	}
	if ipt.Kind != nil {
		mod.Kind = *ipt.Kind
		cols = append(cols, "kind")
	}
	if ipt.Url != nil {
		mod.Url = strings.TrimSpace(*ipt.Url)
		cols = append(cols, "url")
	}
	if ipt.Intro != nil {
		mod.Intro = strings.TrimSpace(*ipt.Intro)
		cols = append(cols, "intro")
	}
	if ipt.Icon != nil {
		mod.Icon = strings.TrimSpace(*ipt.Icon)
		cols = append(cols, "icon")
	}
	if ipt.Color != nil {
		mod.Color = strings.TrimSpace(*ipt.Color)
		cols = append(cols, "color")
	}
	if ipt.Sort != nil {
		mod.Sort = *ipt.Sort
		cols = append(cols, "sort")
	}
	if ipt.Use != nil {
		mod.Use = *ipt.Use
		cols = append(cols, "use")
	}
	if ipt.Show != nil {
		mod.Show = *ipt.Show
		cols = append(cols, "show")
	}
	if ipt.Auth != nil {
		mod.Auth = *ipt.Auth
		cols = append(cols, "auth")
	}
	if ipt.Itpl != nil {
		mod.Itpl = strings.TrimSpace(*ipt.Itpl)
		cols = append(cols, "itpl")
	}
	if ipt.Ltpl != nil {
		mod.Ltpl = strings.TrimSpace(*ipt.Ltpl)
		cols = append(cols, "ltpl")
	}
	if ipt.Dtpl != nil {
		mod.Dtpl = strings.TrimSpace(*ipt.Dtpl)
		cols = append(cols, "dtpl")
	}
	if len(cols) < 1 {
		// PUT 语义幂等：提交内容与现状一致时不报错
		return reply.Succ(ctx, "修改成功")
	}
	mod.Updated = time.Now().UnixMilli()
	cols = append(cols, "updated")
	if err := repo.CateEdit(rctx, mod, cols...); err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.CateEdit")
		return reply.Fail(ctx, "修改失败")
	}
	return reply.Succ(ctx, "修改成功")
}

// CateDrop doc
// @Tags cate
// @Summary 栏目删除
// @Security KeyAuth
// @Param id path string true "栏目id"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/cates/{id} [delete]
func CateDrop(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return reply.Fail(ctx, "参数有误")
	}
	rctx := ctx.Request().Context()
	mod, err := repo.CateGet(rctx, id)
	if err != nil {
		return reply.Fail(ctx, "栏目不存在")
	}
	if mod.Inner {
		return reply.Fail(ctx, "内置栏目不可删除")
	}
	if repo.CateHasChild(id) {
		return reply.Fail(ctx, "请先删除子级栏目")
	}
	if repo.CatePostCount(rctx, repo.CateSubIds(id)) > 0 {
		return reply.Fail(ctx, "栏目下还有内容，无法删除")
	}
	if err := repo.CateDrop(rctx, id); err != nil {
		logs.Ctx(rctx).Err(err).Str("id", id).Error("repo.CateDrop")
		return reply.Fail(ctx, "删除失败")
	}
	return reply.Succ(ctx, "删除成功")
}
