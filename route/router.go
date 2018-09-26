package route

import (
	"blog/control"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Run 入口
func Run(addr string) {
	engine := echo.New()
	engine.Renderer = initRender()
	engine.Use(middleware.Recover())
	engine.Use(middleware.LoggerWithConfig(logConfig))
	engine.Use(middleware.CORSWithConfig(crosConfig)) //允许跨域
	engine.HideBanner = true
	engine.HTTPErrorHandler = HTTPErrorHandler
	engine.Static(`/res`, "res")
	engine.Static(`/static`, "static")
	engine.File(`/favicon.ico`, "favicon.ico")
	engine.GET(`/`, control.Index)
	engine.GET(`/archives`, control.Archives)
	engine.GET(`/tags`, control.Tags)
	engine.GET(`/about`, control.About)
	engine.GET(`/links`, control.Links)
	engine.GET(`/post/*`, control.Post)
	engine.GET(`/page/*`, control.Page)

	admin := engine.Group(`/adm`)
	admin.GET(`/opts/base`, control.OptsBase)
	admin.GET(`/opts/:key`, control.OptsGet)
	admin.POST(`/opts/edit`, control.OptsEdit)
	engine.Start(addr)
}
