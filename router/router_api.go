package router

import (
	"blog/control/appctl"
	"blog/control/sysctl"

	"github.com/labstack/echo/v4"
)

// apiRouter 通用访问
func apiRouter(api *echo.Group) {
	api.GET("/auth/vcode", sysctl.AuthVcode)    // 验证码
	api.GET("/global/get", sysctl.GlobalGet)    // 全局配置
	api.POST("/auth/login", sysctl.AuthLogin)   // 登陆
	api.POST("/auth/logout", sysctl.UserLogout) // 注销
	api.GET("/cate/get", appctl.CateGet)        // 单个分类
	api.GET("/cate/all", appctl.CateAll)        // 所有分类
	api.GET("/cate/page", appctl.CatePage)      // 分类分页
	api.GET("/post/tag/get", appctl.PostTagGet) // 通过分类查询文章
	api.GET("/post/get", appctl.PostGet)        // 文章
	api.GET("/post/page", appctl.PostPage)      // 文章分页
	api.GET("/page/all", appctl.PostPageAll)    // 页面
	api.GET("/tag/all", appctl.TagAll)          // 所有标签
	api.GET("/tag/page", appctl.TagPage)        // 标签分页
}
