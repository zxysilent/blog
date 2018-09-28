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

	engine.GET(`/`, control.Index)
	engine.GET(`/archives`, control.Archives)
	engine.GET(`/tags`, control.Tags)
	engine.GET(`/tag/:tag`, control.TagPost)
	engine.GET(`/about`, control.About)
	engine.GET(`/links`, control.Links)
	engine.GET(`/post/*`, control.Post)
	engine.GET(`/page/*`, control.Page)
	engine.File(`/favicon.ico`, "favicon.ico")

	engine.GET(`/user/exist/:num`, control.UserExist)
	engine.POST(`/login`, control.UserLogin)
	engine.POST(`/logout`, control.UserLogout)

	app := engine.Group(``, midJwt)
	app.GET(`/auth`, control.UserAuth)
	app.POST(`/user/edit/self`, control.UserEditSelf)
	app.POST(`/user/pass`, control.UserPass)

	adm := app.Group(`/adm`)
	adm.GET(`/cate/all`, control.CateAll)
	adm.GET(`/cate/del/:id`, control.CateDel)
	adm.POST(`/cate/add`, control.CateAdd)
	adm.POST(`/cate/edit`, control.CateEdit)

	adm.GET(`/post/get/:id`, control.PostGet)

	adm.GET(`/cate/post/:cid`, control.CatePost)
	adm.GET(`/page/all`, control.PostPageAll)

	adm.GET(`/tag/all`, control.TagAll)
	adm.GET(`/tag/del/:id`, control.TagDel)
	adm.POST(`/tag/add`, control.TagAdd)
	adm.POST(`/tag/edit`, control.TagEdit)

	adm.GET(`/opts/base`, control.OptsBase)
	adm.GET(`/opts/:key`, control.OptsGet)
	adm.POST(`/opts/edit`, control.OptsEdit)
	engine.Start(addr)
}
