//go:build prod
// +build prod

package router

import (
	"github.com/labstack/echo/v4"
	"github.com/zxysilent/logs"
)

const AppJsUrl = "/static/js/app.min.js"
const AppCssUrl = "/static/css/app.min.css"

func init() {
	logs.SetLevel(logs.LDEBUG)
	logs.SetCaller(true)
	logs.SetSep("/blog")
	logs.SetFile("./logs/app.log")
}

// RegDocs 注册文档
// prod[正式] 模式不需要文档
func RegDocs(engine *echo.Echo) {
}

/*  正式模式 编译 取消文档
 *  go build -tags=prod -ldflags "-s -w" -o prod.exe .\main.go
 *
 *  开发模式 编译 添加文档
 *  go build -ldflags "-s -w" -o dev.exe .\main.go
 */
