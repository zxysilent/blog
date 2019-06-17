// +build !prod

package route

import (
	// docs
	_ "blog/docs"

	"github.com/labstack/echo"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// RegDocs 注册文档
// dev[开发] 模式需要文档
func RegDocs(engine *echo.Echo) {
	engine.GET("/swagger/*", echoSwagger.WrapHandler)
}
