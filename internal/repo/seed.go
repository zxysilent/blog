package repo

import (
	"blog/internal/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// 默认根账号
const (
	defaultRootNum    = "root"
	defaultRootPasswd = "admin888"
)

// ConfigIdLayoutLegacy 布局配置的历史主键
// 该配置最初叫「项目配置」、主键为 project，现改名为 layout；
// 启动时若发现旧行则把值搬到新行并删除旧行，避免升级后布局被重置
const ConfigIdLayoutLegacy = "project"

// seed 幂等初始化数据
// 根角色 / 根用户 / 配置 / 字典
// 权限点由路由注册后统一同步（见 router.BuildRouteNames）
func seed(ctx context.Context) error {
	if err := seedRole(ctx); err != nil {
		return err
	}
	if err := seedUser(ctx); err != nil {
		return err
	}
	if err := seedConfig(ctx); err != nil {
		return err
	}
	return seedDict(ctx)
}

// seedRole 内置根角色，主键与层级固定
func seedRole(ctx context.Context) error {
	mod, err := RoleGet(ctx, model.RoleRootId)
	if err == nil {
		// 修复历史数据：根角色始终是内置且层级最高（数值最小）
		if mod.Level == model.RoleRootLevel && mod.Inner {
			return nil
		}
		mod.Level = model.RoleRootLevel
		mod.Inner = true
		return RoleEdit(ctx, mod, "level", "inner")
	}
	if !errors.Is(err, ErrNotExist) {
		return err
	}
	milli := time.Now().UnixMilli()
	mod = &model.Role{
		Id:      model.RoleRootId,
		Name:    "根角色",
		Intro:   "内置根角色-拥有所有权限且无法修改",
		Level:   model.RoleRootLevel,
		Inner:   true,
		Updated: milli,
		Created: milli,
	}
	return RoleAdd(ctx, mod)
}

// seedUser 默认根用户
func seedUser(ctx context.Context) error {
	if UserExistNum(ctx, defaultRootNum) {
		return nil
	}
	milli := time.Now().UnixMilli()
	mod := &model.User{
		Num:     defaultRootNum,
		Name:    "根管理员",
		RoleId:  model.RoleRootId,
		Status:  model.UserStatusNormal,
		Updated: milli,
		Created: milli,
	}
	if err := mod.SetPasswd(defaultRootPasswd); err != nil {
		return err
	}
	return UserAdd(ctx, mod)
}

// migrateLoginStyle 把登录页风格从 login 配置迁到 layout 配置
//
// 旧版本 style 存在 login.value 里，现在它与导航模式、导航风格同属布局配置，
// 在「布局配置」抽屉里一起切换。这里只迁一次：
// layout 已经有 login_style 时不动，避免覆盖用户后来在抽屉里选的风格。
func migrateLoginStyle(ctx context.Context) error {
	layoutMod, err := ConfigGet(ctx, model.ConfigIdLayout)
	if err != nil {
		return nil // 布局配置尚未创建，后续插入流程会用默认值
	}
	raw := map[string]any{}
	if err := json.Unmarshal([]byte(layoutMod.Value), &raw); err != nil {
		return nil
	}
	if _, ok := raw["login_style"]; ok {
		return nil // 已迁移过
	}
	loginMod, err := ConfigGet(ctx, model.ConfigIdLogin)
	if err != nil {
		return nil
	}
	old := struct {
		Style string `json:"style"`
	}{}
	if err := json.Unmarshal([]byte(loginMod.Value), &old); err != nil || old.Style == "" {
		return nil // 历史数据没有 style，保持默认
	}
	layout := layoutMod.ToLayout()
	layout.LoginStyle = old.Style
	layout.Normalize()
	buf, err := json.Marshal(layout)
	if err != nil {
		return err
	}
	milli := time.Now().UnixMilli()
	layoutMod.Value = string(buf)
	layoutMod.Updated = milli
	if err := ConfigEdit(ctx, layoutMod, "value", "updated"); err != nil {
		return err
	}
	// 顺手把 login 里的 style 去掉，避免同一个值存在两处
	login := loginMod.ToLogin()
	buf, err = json.Marshal(login)
	if err != nil {
		return err
	}
	loginMod.Value = string(buf)
	loginMod.Updated = milli
	return ConfigEdit(ctx, loginMod, "value", "updated")
}

// migrateAuthDocs 给历史认证配置补上「接口文档」开关
//
// 老版本的 auth.value 只有 token_exp/token_secret/audit，反序列化时缺失的 docs
// 会取 Go 零值 false，等于升级后文档被静默关闭。这里只在键不存在时补默认值，
// 用户后来自己在后台关掉的状态不会被覆盖。
func migrateAuthDocs(ctx context.Context) error {
	mod, err := ConfigGet(ctx, model.ConfigIdAuth)
	if err != nil {
		return nil // 认证配置尚未创建，后续插入流程会用默认值
	}
	raw := map[string]any{}
	if err := json.Unmarshal([]byte(mod.Value), &raw); err != nil {
		return nil
	}
	if _, ok := raw["docs"]; ok {
		return nil // 已有该字段，尊重当前设置
	}
	auth := mod.ToAuth()
	auth.Docs = model.AuthDefault().Docs
	buf, err := json.Marshal(auth)
	if err != nil {
		return err
	}
	mod.Value = string(buf)
	mod.Updated = time.Now().UnixMilli()
	return ConfigEdit(ctx, mod, "value", "updated")
}

// seedConfig 内置配置
func seedConfig(ctx context.Context) error {
	// 先把历史留在 login 里的页面风格搬到 layout
	if err := migrateLoginStyle(ctx); err != nil {
		return err
	}
	// 再给历史认证配置补上接口文档开关
	if err := migrateAuthDocs(ctx); err != nil {
		return err
	}
	milli := time.Now().UnixMilli()
	// 认证配置：默认值来自 model，避免两处硬编码默认值漂移
	authValue, err := json.Marshal(model.AuthDefault())
	if err != nil {
		return err
	}
	uploadValue, err := json.Marshal(model.UploadDefault())
	if err != nil {
		return err
	}
	sqlValue, err := json.Marshal(model.SqlDefault())
	if err != nil {
		return err
	}
	// 布局配置：默认值同样来自 model，避免前后端默认值各写一份
	layoutValue, err := json.Marshal(model.LayoutDefault())
	if err != nil {
		return err
	}
	items := []model.Config{
		{
			Id:       model.ConfigIdBasic,
			Group:    model.ConfigGroupSite,
			Name:     "站点基础配置",
			Type:     model.CfgTypeJson,
			Value:    `{"title":"blog","site_url":"","logo_url":"","author":"","copyright":"","keywords":"","description":"","beian_miit":"","beian_nism":""}`,
			Describe: "站点标题、Logo、备案等基础信息",
			Editable: true,
			Sort:     1,
		},
		{
			Id:       model.ConfigIdBlog,
			Group:    model.ConfigGroupSite,
			Name:     "博客运行配置",
			Type:     model.CfgTypeJson,
			Value:    `{"page_size":10,"analytic":"","comment":"","github_url":"","weibo_url":""}`,
			Describe: "分页大小、统计与评论脚本、社交链接",
			Editable: true,
			Sort:     2,
		},
		{
			Id:       model.ConfigIdLogin,
			Group:    model.ConfigGroupAuth,
			Name:     "登录页配置",
			Type:     model.CfgTypeJson,
			Value:    `{"title":"欢迎登录","desc":"便捷、高效、可靠"}`,
			Describe: "登录页标题与描述（页面风格在布局配置里切换）",
			Editable: true,
			Sort:     3,
		},
		{
			Id:       model.ConfigIdLayout,
			Group:    model.ConfigGroupUI,
			Name:     "布局配置",
			Type:     model.CfgTypeJson,
			Value:    string(layoutValue),
			Describe: "导航模式、导航风格、登录页风格、多页签、面包屑与路由动画",
			Editable: true,
			Sort:     7,
		},
		{
			Id:    model.ConfigIdAuth,
			Group: model.ConfigGroupAuth,
			Name:  "认证配置",
			Type:  model.CfgTypeJson,
			Value: string(authValue),
			Describe: fmt.Sprintf("令牌过期时间 %d~%d 小时，加密私钥 %d~%d 字符，接口文档开关",
				model.AuthTokenExpMin, model.AuthTokenExpMax, model.AuthSecretMin, model.AuthSecretMax),
			Editable: true,
			Sort:     4,
		},
		{
			Id:    model.ConfigIdUpload,
			Group: model.ConfigGroupSystem,
			Name:  "上传配置",
			Type:  model.CfgTypeJson,
			Value: string(uploadValue),
			Describe: fmt.Sprintf("文件 %d~%d MB，图片 %d~%d MB，边长 %d~%d px",
				model.UploadFileMaxSizeMin, model.UploadFileMaxSizeMax,
				model.UploadImageMaxSizeMin, model.UploadImageMaxSizeMax,
				model.UploadImagePxMin, model.UploadImagePxMax),
			Editable: true,
			Sort:     5,
		},
		{
			Id:    model.ConfigIdSql,
			Group: model.ConfigGroupSystem,
			Name:  "SQL 日志配置",
			Type:  model.CfgTypeJson,
			Value: string(sqlValue),
			Describe: fmt.Sprintf("慢 SQL 阈值 %d~%d 毫秒（0 关闭），可选是否打印全部 SQL",
				model.SqlSlowMsMin, model.SqlSlowMsMax),
			Editable: true,
			Sort:     6,
		},
	}
	for idx := range items {
		it := items[idx]
		if ConfigExist(ctx, it.Id) {
			// 历史数据：分组为空时补上，其余字段保持用户已改的值
			if old, err := ConfigGet(ctx, it.Id); err == nil && old.Group == "" {
				old.Group = it.Group
				if err := ConfigEdit(ctx, old, "group"); err != nil {
					return err
				}
			}
			continue
		} // 布局配置的主键由 project 改为 layout：把旧行的值迁过来再删旧行，
		// 保证升级后用户此前调过的布局不丢
		if it.Id == model.ConfigIdLayout {
			if old, err := ConfigGet(ctx, ConfigIdLayoutLegacy); err == nil {
				it.Value = old.Value
				it.Created = old.Created
				it.Updated = old.Updated
				if err := ConfigAdd(ctx, &it); err != nil {
					return err
				}
				if err := ConfigDrop(ctx, ConfigIdLayoutLegacy); err != nil {
					return err
				}
				continue
			}
		}
		it.Updated = milli
		it.Created = milli
		if err := ConfigAdd(ctx, &it); err != nil {
			return err
		}
	}
	return nil
}

// seedDict 内置字典
// 已存在时跳过，不做覆盖，避免冲掉用户改动
func seedDict(ctx context.Context) error {
	items := []struct {
		mod  *model.Dict
		rows []model.DictEntry
	}{
		{
			mod: &model.Dict{
				Id:     model.DictIdGender,
				Name:   "性别",
				Type:   model.DictTypeInt,
				Status: true,
				Inner:  true,
				Remark: "0-未知 1-男 2-女",
			},
			rows: []model.DictEntry{
				{Field: "unknown", Label: "未知", Value: "0", Sort: 1, Status: true},
				{Field: "male", Label: "男", Value: "1", Sort: 2, Status: true},
				{Field: "female", Label: "女", Value: "2", Sort: 3, Status: true},
			},
		},
		{
			mod: &model.Dict{
				Id:     model.DictIdPageSize,
				Name:   "分页大小",
				Type:   model.DictTypeInt,
				Status: true,
				Inner:  true,
				Remark: "列表接口可选的分页条数",
			},
			rows: []model.DictEntry{
				{Field: "size_15", Label: "15 条/页", Value: "15", Sort: 1, Status: true},
				{Field: "size_20", Label: "20 条/页", Value: "20", Sort: 2, Status: true},
				{Field: "size_30", Label: "30 条/页", Value: "30", Sort: 3, Status: true},
				{Field: "size_50", Label: "50 条/页", Value: "50", Sort: 4, Status: true},
				{Field: "size_100", Label: "100 条/页", Value: "100", Sort: 5, Status: true},
			},
		},
	}
	milli := time.Now().UnixMilli()
	for _, item := range items {
		if DictExist(ctx, item.mod.Id) {
			continue
		}
		mod := item.mod
		mod.Updated = milli
		mod.Created = milli
		if err := mod.Valid(); err != nil {
			return err
		}
		// 先分配主键，选项入库时才有 dict_id 可用
		mod.BeforeInsert()
		rows, err := model.FlatDictEntries(mod.Id, mod.Type, item.rows)
		if err != nil {
			return err
		}
		for idx := range rows {
			rows[idx].Created = milli
			rows[idx].Updated = milli
		}
		if err := DictAdd(ctx, mod, rows); err != nil {
			return err
		}
	}
	return nil
}
