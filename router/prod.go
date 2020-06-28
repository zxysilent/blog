// +build prod

package router

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
)

const AppJsUrl = "/static/js/app.min.js"
const AppCssUrl = "/static/css/app.min.css"

func init() {
	log.SetFlags(log.Lshortfile | log.Ldate | log.Ltime)
	log.SetPrefix("[app] ")
	fd, err := os.OpenFile("app.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalln("init error:", err)
	}
	log.SetOutput(fd)
}

// RegDocs 注册文档
// prod[正式] 模式不需要文档
func RegDocs(engine *echo.Echo) {

}

// midLogger 中间件-日志记录
func midLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return next
}

/*  正式模式 编译 取消文档
 *  生成文档 swag init
 *  go build -tags=prod -o blog.exe .\main.go
 *
 *  开发模式 编译 添加文档
 *  go build -o blogdev.exe .\main.go
 */
