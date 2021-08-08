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
	api.GET("/dict/get", sysctl.DictGet)        //单个字典
	api.GET("/dict/page", sysctl.DictPage)      //字典分页
	api.POST("/auth/login", sysctl.AuthLogin)   // 登陆
	api.POST("/auth/logout", sysctl.UserLogout) // 注销
	api.GET("/cate/get", appctl.CateGet)        // 单个分类
	api.GET("/cate/all", appctl.CateAll)        // 所有分类
	api.GET("/cate/page", appctl.CatePage)      // 分类分页
	api.GET("/tag/get", appctl.TagGet)          // 单个标签
	api.GET("/tag/all", appctl.TagAll)          // 所有标签
	api.GET("/tag/page", appctl.TagPage)        // 标签分页
	api.GET("/post/get", appctl.PostGet)        // 单个文章
	api.GET("/post/page", appctl.PostPage)      // 文章分页
	api.GET("/page/get", appctl.PageGet)        // 单个页面
	api.GET("/page/page", appctl.PagePage)      // 页面分页

}
