//go:build !prod
// +build !prod

package router

import (
	_ "blog/docs"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/logs"
)

const AppJsUrl = "/static/js/app.js"
const AppCssUrl = "/static/css/app.css"

func init() {
	logs.SetLevel(logs.LDEBUG)
	logs.SetCaller(true)
	logs.SetSep("/blog")
	logs.SetFile("./logs/app.log")
	logs.SetCons(true)
}

// RegDocs 注册文档
// dev[开发] 模式需要文档
func RegDocs(engine *echo.Echo) {
	swagger := EchoWrapHandler(DocExpansion("none"))
	engine.GET("/swagger/*", swagger)
}
