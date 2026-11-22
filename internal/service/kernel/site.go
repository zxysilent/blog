package kernel

import (
	"blog/internal/model"
	"blog/internal/repo"
	"blog/pkg/reply"

	"github.com/labstack/echo/v4"
)

// SiteConfig 站点公开配置聚合
type SiteConfig struct {
	Basic any  `json:"basic"` //站点基础配置
	Blog  any  `json:"blog"`  //博客运行配置
	Login any  `json:"login"` //登录页配置
	Docs  bool `json:"docs"`  //是否开放接口文档
}

// Site doc
// @Tags site
// @Summary 获取站点公开配置（basic/blog/login/docs）
// @Success 200 {object} reply.Reply{data=kernel.SiteConfig} "返回数据"
// @Router /api/site [get]
func Site(ctx echo.Context) error {
	cfg := repo.GetCfg()
	// 登录页风格存在 layout 配置里，但登录页只认 login.style 一个字段，
	// 这里做一次投影，前端不用自己拼两份配置
	login := model.LoginSite{
		Title: cfg.Login.Title,
		Desc:  cfg.Login.Desc,
		Style: cfg.Layout.LoginStyle,
	}
	// docs 一并下发：管理端首页据此决定是否展示「接口文档」入口
	return reply.Succ(ctx, "ok", SiteConfig{
		Basic: cfg.Basic,
		Blog:  cfg.Blog,
		Login: login,
		Docs:  cfg.Auth.Docs,
	})
}

// ConfigVal doc
// @Tags site
// @Summary 获取公开配置项（仅限 basic/blog/login）
// @Param id path string true "配置主键"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/configs/{id}/val [get]
func ConfigVal(ctx echo.Context) error {
	id := ctx.Param("id")
	switch id {
	case model.ConfigIdBasic, model.ConfigIdBlog, model.ConfigIdLogin:
	default:
		return reply.Fail(ctx, "配置不可公开读取")
	}
	cfg := repo.GetCfg()
	return reply.Succ(ctx, "ok", cfg.Get(id))
}
