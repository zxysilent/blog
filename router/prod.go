// +build prod

package router

import "github.com/labstack/echo"

// RegDocs 注册文档
// prod[正式] 模式不需要文档
func RegDocs(engine *echo.Echo) {

}

/*  正式模式 编译 取消文档
 *  生成文档 swag init
 *  go build -tags=prod -o blog.exe .\main.go
 *
 *  开发模式 编译 添加文档
 *  go build -o blogdev.exe .\main.go
 */
