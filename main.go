package main

import (
	"blog/route"
)

// @Title Blog’s Api文档
// @Version 1.0
// @Description token传递方式包括 [get]token、[post] token [header] Authorization="Bearer xxxx"
// @Description 数据传递方式包括 json、formData 文档使用 formData 前后对接推荐使用 json
// @Host 127.0.0.1:88
// @BasePath /
func main() {
	route.Run()
}
