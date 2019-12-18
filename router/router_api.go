package router

import (
	"blog/control"

	"github.com/labstack/echo"
)

// apiRouter 通用访问
func apiRouter(api *echo.Group) {
	// 判断账号是否存在
	api.GET(`/user/exist/:num`, control.UserExist)
	// 登陆
	api.POST(`/login`, control.UserLogin)
	// 注销
	api.POST(`/logout`, control.UserLogout)
	// 分类列表
	api.GET(`/cate/all`, control.CateAll)
	// 通过分类查询文章
	api.GET(`/post/tag/get/:id`, control.PostTagGet)
	// 文章
	api.GET(`/post/get/:id`, control.PostGet)
	// 通过分类查询文章
	api.GET(`/cate/post/:cid`, control.CatePost)
	// 获取配置项
	api.GET(`/opts/:key`, control.OptsGet)
	// 页面
	api.GET(`/page/all`, control.PostPageAll)
	// 标签列表
	api.GET(`/tag/all`, control.TagAll)
	// 配置
	api.GET(`/opts/base`, control.OptsBase)
}
