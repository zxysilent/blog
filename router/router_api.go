package router

import (
	"blog/control/appctl"
	"blog/control/sysctl"

	"github.com/labstack/echo/v4"
)

// apiRouter 通用访问
func apiRouter(api *echo.Group) {
	api.GET("/auth/vcode", sysctl.AuthVcode)        // 验证码
	api.GET("/auth/exist", appctl.UserExist)        // 判断账号是否存在
	api.POST("/auth/login", sysctl.AuthLogin)       // 登陆
	api.POST("/auth/logout", sysctl.UserLogout)     // 注销
	api.GET("/cate/all", appctl.CateAll)            // 分类列表
	api.GET("/post/tag/get/:id", appctl.PostTagGet) // 通过分类查询文章
	api.GET("/post/get/:id", appctl.PostGet)        // 文章
	api.GET("/cate/post/:cid", appctl.CatePost)     // 通过分类查询文章
	api.GET("/opts/:key", appctl.OptsGet)           // 获取配置项
	api.GET("/page/all", appctl.PostPageAll)        // 页面
	api.GET("/tag/all", appctl.TagAll)              // 标签列表
	api.GET("/opts/base", appctl.OptsBase)          // 配置

}
