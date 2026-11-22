package repo

import (
	"blog/internal/model"
	"context"
	"sync"
)

// ConfigSet 进程内配置缓存
type ConfigSet struct {
	Auth   *model.Auth              //认证配置
	Basic  *model.Basic             //站点基础配置
	Blog   *model.Blog              //博客运行配置
	Login  *model.Login             //登录页配置
	Layout *model.Layout            //布局配置（管理端导航、主题、多页签）
	Sql    *model.Sql               //SQL 日志配置
	Upload *model.Upload            //上传配置
	mp     map[string]*model.Config //其余配置
}

// Get 取任意配置，不存在时返回空配置
func (cs *ConfigSet) Get(id string) *model.Config {
	if mod, ok := cs.mp[id]; ok {
		return mod
	}
	return &model.Config{Id: id}
}

var (
	_cfgMu sync.RWMutex
	_cfg   = newConfigSet()
)

func newConfigSet() *ConfigSet {
	return &ConfigSet{
		Auth:   model.AuthDefault(),
		Basic:  &model.Basic{},
		Blog:   &model.Blog{},
		Login:  model.LoginDefault(),
		Layout: model.LayoutDefault(),
		Sql:    model.SqlDefault(),
		Upload: model.UploadDefault(),
		mp:     make(map[string]*model.Config, 16),
	}
}

// GetCfg 取配置集合（只读）
func GetCfg() *ConfigSet {
	_cfgMu.RLock()
	defer _cfgMu.RUnlock()
	return _cfg
}

// GetCfgx 按主键取配置
func GetCfgx(id string) *model.Config {
	_cfgMu.RLock()
	defer _cfgMu.RUnlock()
	return _cfg.Get(id)
}

// SyncCfg 重建配置缓存
func SyncCfg(ctx context.Context) error {
	mods, err := ConfigAll(ctx)
	if err != nil {
		return err
	}
	// 先建空集合再填充，避免删除后旧值残留
	cs := newConfigSet()
	for idx := range mods {
		mod := mods[idx]
		switch mod.Id {
		case model.ConfigIdAuth:
			cs.Auth = mod.ToAuth()
		case model.ConfigIdBasic:
			cs.Basic = mod.ToBasic()
		case model.ConfigIdBlog:
			cs.Blog = mod.ToBlog()
		case model.ConfigIdLogin:
			cs.Login = mod.ToLogin()
		case model.ConfigIdLayout:
			cs.Layout = mod.ToLayout()
		case model.ConfigIdUpload:
			cs.Upload = mod.ToUpload()
		case model.ConfigIdSql:
			cs.Sql = mod.ToSql()
		}
		cs.mp[mod.Id] = &mod
	}
	_cfgMu.Lock()
	_cfg = cs
	_cfgMu.Unlock()
	return nil
}

// ConfigExist 配置主键是否已存在
func ConfigExist(ctx context.Context, id string) bool {
	if id == "" {
		return false
	}
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	exist, _ := sess.ID(id).NoAutoCondition().Exist(&model.Config{})
	return exist
}

// ConfigGet 按主键查询配置
func ConfigGet(ctx context.Context, id string) (*model.Config, error) {
	mod, err := GetAny[model.Config](ctx, id)
	if err != nil {
		return nil, err
	}
	return &mod, nil
}

// ConfigList 配置分页列表
func ConfigList(ctx context.Context, args *model.ConfigListArgs, cols ...string) ([]model.Config, int, error) {
	mods := make([]model.Config, 0, args.GetCap())
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	if args.Mult != "" {
		sess.And("(id like ? or name like ?)", args.SafeMult(), args.SafeMult())
	}
	if args.Group != "" {
		sess.And("`group` = ?", args.Group)
	}
	if args.Name != "" {
		sess.And("name like ?", "%"+args.Name+"%")
	}
	if args.Type != "" {
		sess.And("type = ?", args.Type)
	}
	sess.Asc("`group`").Asc("sort").Asc("id")
	if args.IsZero() {
		err := sess.Find(&mods)
		return mods, len(mods), err
	}
	count, err := sess.Limit(args.GetLimit(), args.GetOffset()).FindAndCount(&mods)
	return mods, int(count), err
}

// ConfigGroups 配置分组及数量，供后台按分组折叠展示
func ConfigGroups(ctx context.Context) ([]model.ConfigGroup, error) {
	groups := make([]model.ConfigGroup, 0, 8)
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	err := sess.Table("config").Select("`group` AS `group`, COUNT(*) AS `count`").
		GroupBy("`group`").Asc("`group`").Find(&groups)
	return groups, err
}

// ConfigAll 全部配置
func ConfigAll(ctx context.Context) ([]model.Config, error) {
	mods := make([]model.Config, 0, 16)
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	err := sess.Asc("`group`").Asc("sort").Asc("id").Find(&mods)
	return mods, err
}

// ConfigAdd 新增配置
func ConfigAdd(ctx context.Context, mod *model.Config) error {
	return AddAny(ctx, mod)
}

// ConfigEdit 修改配置
func ConfigEdit(ctx context.Context, mod *model.Config, cols ...string) error {
	return EditAny(ctx, mod, cols...)
}

// ConfigDrop 删除配置
func ConfigDrop(ctx context.Context, id string) error {
	return DelAny[model.Config](ctx, id)
}
