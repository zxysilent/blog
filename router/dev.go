// +build !prod

package router

import (

	// docs
	_ "blog/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

const AppJsUrl = "/static/js/app.js"
const AppCssUrl = "/static/css/app.css"

// RegDocs 注册文档
// dev[开发] 模式需要文档
func RegDocs(engine *echo.Echo) {
	docUrl := echoSwagger.URL("/swagger/doc.json")
	engine.GET("/swagger/*", echoSwagger.EchoWrapHandler(docUrl))
}
