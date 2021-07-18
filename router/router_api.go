package router

import (
	"blog/control/appctl"
	"blog/control/sysctl"

	"github.com/labstack/echo/v4"
)

// apiRouter 通用访问
func apiRouter(api *echo.Group) {
	api.GET("/auth/vcode", sysctl.AuthVcode)        // 验证码
	api.POST("/auth/login", sysctl.AuthLogin)       // 登陆
	api.POST("/auth/logout", sysctl.UserLogout)     // 注销
	api.GET("/user/exist", sysctl.UserExist)        // 判断账号是否存在
	api.GET("/cate/all", appctl.CateAll)            // 分类列表
	api.GET("/post/tag/get/:id", appctl.PostTagGet) // 通过分类查询文章
	api.GET("/post/get/:id", appctl.PostGet)        // 文章
	api.GET("/cate/post/:cid", appctl.CatePost)     // 通过分类查询文章
	api.GET("/page/all", appctl.PostPageAll)        // 页面
	api.GET("/tag/all", appctl.TagAll)              // 标签列表
	api.GET("/global/get", sysctl.GlobalGet)        // 全局配置
}
