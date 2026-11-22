package kernel

import (
	"blog/internal/model"
	"blog/internal/repo"
	"blog/pkg/reply"
	"encoding/json"
	"path"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/logs"
)

// configProtected 内置配置主键集合，禁止删除
var configProtected = func() map[string]struct{} {
	mp := make(map[string]struct{}, len(model.BuiltinConfigIds))
	for _, id := range model.BuiltinConfigIds {
		mp[id] = struct{}{}
	}
	return mp
}()

// configIdOf 解析当前请求要操作的配置主键
// 内置配置走独立静态路径（如 /api/configs/basic），没有 :id 路径参数；
// 自建配置走 /api/configs/:id。两种入口统一在这里取到主键。
func configIdOf(ctx echo.Context) string {
	if id := strings.TrimSpace(ctx.Param("id")); id != "" {
		return id
	}
	return path.Base(ctx.Path())
}

// checkConfigPerm 内置配置的写权限校验
// 每个内置配置都有独立权限点（config:basic 等），只拥有 config:update 的账号
// 不能借通用接口 /configs/:id 绕过去；根角色由 hasPerm 直接放行
func checkConfigPerm(ctx echo.Context, id string) bool {
	if _, ok := configProtected[id]; !ok {
		return true
	}
	return hasPerm(ctx, model.ConfigPerm(id))
}

// ConfigList doc
// @Tags config
// @Summary 配置列表
// @Security KeyAuth
// @Param query query model.ConfigListArgs true "请求数据"
// @Success 200 {object} reply.Reply{data=[]model.Config} "返回数据"
// @Router /api/configs [get]
func ConfigList(ctx echo.Context) error {
	args := &model.ConfigListArgs{}
	if err := ctx.Bind(args); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	rctx := ctx.Request().Context()
	mods, count, err := repo.ConfigList(rctx, args)
	if err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.ConfigList")
		return reply.Fail(ctx, "查询失败")
	}
	return reply.List(ctx, "ok", mods, count)
}

// ConfigGroups doc
// @Tags config
// @Summary 配置分组列表（含数量，供后台折叠展示）
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=[]model.ConfigGroup} "返回数据"
// @Router /api/configs/groups [get]
func ConfigGroups(ctx echo.Context) error {
	rctx := ctx.Request().Context()
	mods, err := repo.ConfigGroups(rctx)
	if err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.ConfigGroups")
		return reply.Fail(ctx, "查询失败")
	}
	return reply.Succ(ctx, "ok", mods)
}

// ConfigGet doc
// @Tags config
// @Summary 配置详情
// @Security KeyAuth
// @Param id path string true "配置主键"
// @Success 200 {object} reply.Reply{data=model.Config} "返回数据"
// @Router /api/configs/{id} [get]
func ConfigGet(ctx echo.Context) error {
	rctx := ctx.Request().Context()
	mod, err := repo.ConfigGet(rctx, configIdOf(ctx))
	if err != nil {
		logs.Ctx(rctx).Err(err).Str("id", configIdOf(ctx)).Error("repo.ConfigGet")
		return reply.Fail(ctx, "配置不存在")
	}
	return reply.Succ(ctx, "ok", mod)
}

// ConfigAdd doc
// @Tags config
// @Summary 配置新增（id 可自定义，留空自动生成）
// @Security KeyAuth
// @Param body body model.ConfigAddArgs true "请求数据"
// @Success 200 {object} reply.Reply{data=string} "返回数据"
// @Router /api/configs [post]
func ConfigAdd(ctx echo.Context) error {
	ipt := &model.ConfigAddArgs{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	milli := time.Now().UnixMilli()
	mod := &model.Config{
		Id:       strings.TrimSpace(ipt.Id),
		Group:    ipt.Group,
		Name:     ipt.Name,
		Value:    ipt.Value,
		Min:      ipt.Min,
		Max:      ipt.Max,
		Type:     ipt.Type,
		Describe: ipt.Describe,
		Editable: true,
		Sort:     ipt.Sort,
		Updated:  milli,
		Created:  milli,
	}
	// 未指定主键时生成 ULID，再校验格式
	mod.BeforeInsert()
	if err := mod.Valid(); err != nil {
		return reply.Fail(ctx, err.Error())
	}
	if err := mod.NormalizeValue(); err != nil {
		return reply.Fail(ctx, err.Error())
	}
	rctx := ctx.Request().Context()
	if repo.ConfigExist(rctx, mod.Id) {
		return reply.Fail(ctx, "配置已存在")
	}
	if err := repo.ConfigAdd(rctx, mod); err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.ConfigAdd")
		return reply.Fail(ctx, "添加失败")
	}
	if err := repo.SyncCfg(rctx); err != nil {
		logs.Ctx(rctx).Err(err).Warn("repo.SyncCfg")
	}
	return reply.Succ(ctx, "添加成功", mod.Id)
}

// ConfigEdit doc
// @Tags config
// @Summary 配置修改
// @Security KeyAuth
// @Param id path string true "配置主键"
// @Param body body model.ConfigEditArgs true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/configs/{id} [put]
func ConfigEdit(ctx echo.Context) error {
	ipt := &model.ConfigEditArgs{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	rctx := ctx.Request().Context()
	mod, err := repo.ConfigGet(rctx, configIdOf(ctx))
	if err != nil {
		return reply.Fail(ctx, "配置不存在")
	}
	if !mod.Editable {
		return reply.Fail(ctx, "该配置为只读，不可修改")
	}
	// 内置配置独立授权：走通用接口同样要校验对应权限
	if !checkConfigPerm(ctx, mod.Id) {
		return reply.Fail(ctx, "无权修改该配置")
	}
	cols := make([]string, 0, 3)
	if ipt.Value != nil {
		value := *ipt.Value
		// 布局配置里混着登录页风格（需要 login:style），而布局抽屉保存时只提交
		// 自己那几个字段。这里把缺失的 login_style 从旧值补回来：
		// 既避免整份配置被写成空风格，也让「没带该字段」不等于「要清空该字段」。
		if mod.Id == model.ConfigIdLayout {
			value = mergeLoginStyle(mod.Value, value)
		}
		// 登录页风格有独立权限：它属于站点门面，只拥有 config:layout 不足以改动。
		// 比对旧值判断是否真的改了，未改则放行
		if mod.Id == model.ConfigIdLayout && !hasPerm(ctx, model.PermLoginStyle) {
			if layoutLoginStyleOf(value) != mod.ToLayout().LoginStyle {
				return reply.Fail(ctx, "无权修改登录页风格")
			}
		}
		mod.Value = value
		if err := mod.NormalizeValue(); err != nil {
			return reply.Fail(ctx, err.Error())
		}
		if err := mod.Valid(); err != nil {
			return reply.Fail(ctx, err.Error())
		}
		cols = append(cols, "value")
	}
	if ipt.Describe != nil {
		mod.Describe = strings.TrimSpace(*ipt.Describe)
		cols = append(cols, "describe")
	}
	if ipt.Group != nil {
		mod.Group = strings.TrimSpace(*ipt.Group)
		if mod.Group == "" {
			return reply.Fail(ctx, "配置分组不能为空")
		}
		if len(mod.Group) > 64 {
			return reply.Fail(ctx, "配置分组长度不能超过 64")
		}
		cols = append(cols, "group")
	}
	if ipt.Name != nil {
		mod.Name = strings.TrimSpace(*ipt.Name)
		if mod.Name == "" {
			return reply.Fail(ctx, "配置名称不能为空")
		}
		if len(mod.Name) > 255 {
			return reply.Fail(ctx, "配置名称长度不能超过 255")
		}
		cols = append(cols, "name")
	}
	if ipt.Min != nil {
		mod.Min = *ipt.Min
		cols = append(cols, "min")
	}
	if ipt.Max != nil {
		mod.Max = *ipt.Max
		cols = append(cols, "max")
	}
	if ipt.Sort != nil {
		mod.Sort = *ipt.Sort
		cols = append(cols, "sort")
	}
	// 区间校验：只有 min/max 参与时才有意义，0/0 表示不限制
	if mod.Max != 0 && mod.Min > mod.Max {
		return reply.Fail(ctx, "最小值不能大于最大值")
	}
	if len(cols) < 1 {
		return reply.Fail(ctx, "没有需要修改的数据")
	}
	mod.Updated = time.Now().UnixMilli()
	cols = append(cols, "updated")
	if err := repo.ConfigEdit(rctx, mod, cols...); err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.ConfigEdit")
		return reply.Fail(ctx, "修改失败")
	}
	if err := repo.SyncCfg(rctx); err != nil {
		logs.Ctx(rctx).Err(err).Warn("repo.SyncCfg")
	}
	return reply.Succ(ctx, "修改成功")
}

// layoutLoginStyleOf 从布局配置 JSON 文本中取出登录页风格
// 解析失败或字段缺失时返回默认值，与 model.Config.ToLayout 的口径保持一致
func layoutLoginStyleOf(value string) string {
	mod := model.LayoutDefault()
	if err := json.Unmarshal([]byte(value), mod); err != nil {
		return model.LoginStyleDefault
	}
	mod.Normalize()
	return mod.LoginStyle
}

// mergeLoginStyle 把旧布局配置里的 login_style 补进新值
//
// 布局抽屉保存的是「它自己那几个字段」，不带 login_style；
// 若直接落库会让风格被重置为默认，甚至触发 Valid 的枚举校验失败。
// 这里只补字段，不动调用方显式传入的值。
func mergeLoginStyle(oldValue, newValue string) string {
	var oldRaw, newRaw map[string]any
	if err := json.Unmarshal([]byte(oldValue), &oldRaw); err != nil {
		return newValue // 旧值坏了就按新值走，交给 Valid 兜底
	}
	if _, ok := oldRaw["login_style"]; !ok {
		return newValue
	}
	if err := json.Unmarshal([]byte(newValue), &newRaw); err != nil {
		return newValue
	}
	if _, ok := newRaw["login_style"]; ok {
		return newValue // 调用方显式传了，尊重它（权限校验随后处理）
	}
	newRaw["login_style"] = oldRaw["login_style"]
	buf, err := json.Marshal(newRaw)
	if err != nil {
		return newValue
	}
	return string(buf)
}

// ConfigLoginStyle doc
// @Tags config
// @Summary 切换登录页风格（需 login:style 权限）
// @Security KeyAuth
// @Param body body model.ConfigLoginStyleArgs true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/configs/login/style [put]
func ConfigLoginStyle(ctx echo.Context) error {
	ipt := &model.ConfigLoginStyleArgs{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	style := strings.ToLower(strings.TrimSpace(ipt.Style))
	if style != model.LoginStyleMgmt && style != model.LoginStyleVue {
		return reply.Fail(ctx, "登录页风格仅支持 "+model.LoginStyleMgmt+" 或 "+model.LoginStyleVue)
	}
	rctx := ctx.Request().Context()
	// 风格存在布局配置里，与导航模式、导航风格同源
	mod, err := repo.ConfigGet(rctx, model.ConfigIdLayout)
	if err != nil {
		return reply.Fail(ctx, "布局配置不存在")
	}
	layout := mod.ToLayout()
	if layout.LoginStyle == style {
		return reply.Succ(ctx, "风格未变更")
	}
	layout.LoginStyle = style
	buf, err := json.Marshal(layout)
	if err != nil {
		return reply.Fail(ctx, "风格保存失败")
	}
	mod.Value = string(buf)
	if err := mod.NormalizeValue(); err != nil {
		return reply.Fail(ctx, err.Error())
	}
	mod.Updated = time.Now().UnixMilli()
	if err := repo.ConfigEdit(rctx, mod, "value", "updated"); err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.ConfigEdit")
		return reply.Fail(ctx, "修改失败")
	}
	if err := repo.SyncCfg(rctx); err != nil {
		logs.Ctx(rctx).Err(err).Warn("repo.SyncCfg")
	}
	return reply.Succ(ctx, "风格已切换为 "+style)
}

// ------------------------------------------------------ 内置配置独立入口 ------------------------------------------------------
//
// 每个内置配置都有独立的读写路径与独立权限点（config:basic 等）：
//   GET  /api/configs/basic   -> config:basic
//   PUT  /api/configs/basic   -> config:basic
//
// 好处是可以按配置分权：只给运营角色「站点基础配置」权限，它就改不了认证私钥
// 与上传上限。路径是静态段，echo 里优先于 /configs/:id 匹配。
// 处理函数直接复用 ConfigGet / ConfigEdit，权限由路由 Name 决定。

// ConfigBasicGet doc
// @Tags config
// @Summary 站点基础配置-读取
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=model.Config} "返回数据"
// @Router /api/configs/basic [get]
func ConfigBasicGet(ctx echo.Context) error { return ConfigGet(ctx) }

// ConfigBasicEdit doc
// @Tags config
// @Summary 站点基础配置-修改
// @Security KeyAuth
// @Param body body model.ConfigEditArgs true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/configs/basic [put]
func ConfigBasicEdit(ctx echo.Context) error { return ConfigEdit(ctx) }

// ConfigBlogGet doc
// @Tags config
// @Summary 博客运行配置-读取
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=model.Config} "返回数据"
// @Router /api/configs/blog [get]
func ConfigBlogGet(ctx echo.Context) error { return ConfigGet(ctx) }

// ConfigBlogEdit doc
// @Tags config
// @Summary 博客运行配置-修改
// @Security KeyAuth
// @Param body body model.ConfigEditArgs true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/configs/blog [put]
func ConfigBlogEdit(ctx echo.Context) error { return ConfigEdit(ctx) }

// ConfigLoginGet doc
// @Tags config
// @Summary 登录页配置-读取
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=model.Config} "返回数据"
// @Router /api/configs/login [get]
func ConfigLoginGet(ctx echo.Context) error { return ConfigGet(ctx) }

// ConfigLoginEdit doc
// @Tags config
// @Summary 登录页配置-修改（风格字段另需 login:style）
// @Security KeyAuth
// @Param body body model.ConfigEditArgs true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/configs/login [put]
func ConfigLoginEdit(ctx echo.Context) error { return ConfigEdit(ctx) }

// ConfigAuthGet doc
// @Tags config
// @Summary 认证配置-读取
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=model.Config} "返回数据"
// @Router /api/configs/auth [get]
func ConfigAuthGet(ctx echo.Context) error { return ConfigGet(ctx) }

// ConfigAuthEdit doc
// @Tags config
// @Summary 认证配置-修改
// @Security KeyAuth
// @Param body body model.ConfigEditArgs true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/configs/auth [put]
func ConfigAuthEdit(ctx echo.Context) error { return ConfigEdit(ctx) }

// ConfigUploadGet doc
// @Tags config
// @Summary 上传配置-读取
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=model.Config} "返回数据"
// @Router /api/configs/upload [get]
func ConfigUploadGet(ctx echo.Context) error { return ConfigGet(ctx) }

// ConfigUploadEdit doc
// @Tags config
// @Summary 上传配置-修改
// @Security KeyAuth
// @Param body body model.ConfigEditArgs true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/configs/upload [put]
func ConfigUploadEdit(ctx echo.Context) error { return ConfigEdit(ctx) }

// ConfigSqlGet doc
// @Tags config
// @Summary SQL 日志配置-读取
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=model.Config} "返回数据"
// @Router /api/configs/sql [get]
func ConfigSqlGet(ctx echo.Context) error { return ConfigGet(ctx) }

// ConfigSqlEdit doc
// @Tags config
// @Summary SQL 日志配置-修改
// @Security KeyAuth
// @Param body body model.ConfigEditArgs true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/configs/sql [put]
func ConfigSqlEdit(ctx echo.Context) error { return ConfigEdit(ctx) }

// ConfigLayoutGet doc
// @Tags config
// @Summary 布局配置-读取（管理端导航、主题、多页签）
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=model.Config} "返回数据"
// @Router /api/configs/layout [get]
func ConfigLayoutGet(ctx echo.Context) error { return ConfigGet(ctx) }

// ConfigLayoutEdit doc
// @Tags config
// @Summary 布局配置-修改（管理端导航、主题、多页签）
// @Security KeyAuth
// @Param body body model.ConfigEditArgs true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/configs/layout [put]
func ConfigLayoutEdit(ctx echo.Context) error { return ConfigEdit(ctx) }

// ConfigDrop doc
// @Tags config
// @Summary 配置删除
// @Security KeyAuth
// @Param id path string true "配置主键"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/configs/{id} [delete]
func ConfigDrop(ctx echo.Context) error {
	id := strings.TrimSpace(ctx.Param("id"))
	if id == "" {
		return reply.Fail(ctx, "参数有误")
	}
	if _, ok := configProtected[id]; ok {
		return reply.Fail(ctx, "内置配置不可删除")
	}
	rctx := ctx.Request().Context()
	if err := repo.ConfigDrop(rctx, id); err != nil {
		logs.Ctx(rctx).Err(err).Str("id", id).Error("repo.ConfigDrop")
		return reply.Fail(ctx, "删除失败")
	}
	if err := repo.SyncCfg(rctx); err != nil {
		logs.Ctx(rctx).Err(err).Warn("repo.SyncCfg")
	}
	return reply.Succ(ctx, "删除成功")
}
