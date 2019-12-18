package router

import (
	"blog/conf"
	"blog/control"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// RunApp 入口
func RunApp() {
	engine := echo.New()
	// 初始渲染引擎
	engine.Renderer = initRender()
	// 恢复 日志记录
	engine.Use(midRecover, midLogger)
	// 跨域设置
	engine.Use(middleware.CORSWithConfig(crosConfig))
	// 不显示横幅
	engine.HideBanner = true
	// 自定义错误处理
	engine.HTTPErrorHandler = HTTPErrorHandler
	engine.Debug = conf.App.IsDev()
	// 注册文档
	RegDocs(engine)
	// 静态目录
	engine.Static(`/static`, "static")
	// ico
	engine.File(`/favicon.ico`, "favicon.ico")
	// 前后端分离页面
	engine.File("/dashboard/*", "static/dist/index.html")
	//--- 页面 -- start
	// 首页
	engine.GET(`/`, control.IndexView)
	// 归档
	engine.GET(`/archives`, control.ArchivesView)
	// 标签
	engine.GET(`/tags`, control.TagsView)
	// 具体某个标签
	engine.GET(`/tag/:tag`, control.TagPostView)
	// 分类
	engine.GET(`/cate/:cate`, control.CatePostView)
	// 关于
	engine.GET(`/about`, control.AboutView)
	// 友链
	engine.GET(`/links`, control.LinksView)
	// 具体某个文章
	engine.GET(`/post/*`, control.PostView)
	// 具体某个页面
	engine.GET(`/page/*`, control.PageView)
	//--- 页面 -- end

	api := engine.Group("/api")
	apiRouter(api)
	// 需要登陆才能访问
	adm := engine.Group("/adm", midJwt)
	admRouter(adm)
	engine.Start(conf.App.Addr)
}
