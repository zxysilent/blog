package router

import (
	"blog/control"

	"github.com/labstack/echo/v4"
)

// apiRouter 通用访问
func apiRouter(api *echo.Group) {
	api.GET(`/user/exist/:num`, control.UserExist)   // 判断账号是否存在
	api.POST(`/login`, control.UserLogin)            // 登陆
	api.GET("/vcode", control.Vcode)                 // 验证码
	api.POST(`/logout`, control.UserLogout)          // 注销
	api.GET(`/cate/all`, control.CateAll)            // 分类列表
	api.GET(`/post/tag/get/:id`, control.PostTagGet) // 通过分类查询文章
	api.GET(`/post/get/:id`, control.PostGet)        // 文章
	api.GET(`/cate/post/:cid`, control.CatePost)     // 通过分类查询文章
	api.GET(`/opts/:key`, control.OptsGet)           // 获取配置项
	api.GET(`/page/all`, control.PostPageAll)        // 页面
	api.GET(`/tag/all`, control.TagAll)              // 标签列表
	api.GET(`/opts/base`, control.OptsBase)          // 配置
}
