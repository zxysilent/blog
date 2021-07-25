package router

import (
	"blog/conf"
	"blog/control/appctl"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zxysilent/logs"
)

// RunApp 入口
func RunApp() {
	engine := echo.New()                              // 实例化echo
	engine.Renderer = initRender()                    // 初始渲染引擎
	engine.Use(midRecover, midLogger)                 // 恢复 日志记录
	engine.Use(middleware.CORSWithConfig(crosConfig)) // 跨域设置
	engine.HideBanner = true                          // 不显示横幅
	engine.HTTPErrorHandler = HTTPErrorHandler        // 自定义错误处理
	engine.Debug = conf.App.IsDev()                   // 运行模式 - echo框架好像没怎么使用这个
	RegDocs(engine)                                   // 注册文档
	engine.Static("/dist", "dist")                    // 静态目录 - 后端专用
	engine.Static("/static", "static")                // 静态目录
	engine.File("/favicon.ico", "favicon.ico")        // ico
	engine.File("/dashboard*", "dist/index.html")     // 前后端分离页面
	engine.GET("/dashboard", func(ctx echo.Context) error {
		// 301 永久重定向
		// 302 临时重定向
		return ctx.Redirect(302, "/dashboard/")
	})
	//--- 页面 -- start
	engine.GET("/", appctl.IndexView)              // 首页
	engine.GET("/archives", appctl.ViewArchives)   // 归档
	engine.GET("/tags", appctl.ViewTags)           // 标签
	engine.GET("/tag/:tag", appctl.ViewTagPost)    // 具体某个标签
	engine.GET("/cate/:cate", appctl.ViewCatePost) // 分类
	engine.GET("/about", appctl.ViewAbout)         // 关于
	engine.GET("/links", appctl.ViewLinks)         // 友链
	engine.GET("/page/*", appctl.ViewPage)         // 具体某个页面
	engine.GET("/post/*", appctl.ViewPost)         // 具体某个文章
	//--- 页面 -- end
	api := engine.Group("/api")          // api/
	apiRouter(api)                       // 注册分组路由
	adm := engine.Group("/adm", midAuth) // adm/ 需要登陆才能访问
	admRouter(adm)                       // 注册分组路由
	err := engine.Start(conf.App.Addr)
	if err != nil {
		logs.Fatal("run error :", err.Error())
	}
}
