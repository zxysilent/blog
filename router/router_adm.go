package router

import (
	"blog/control"

	"github.com/labstack/echo"
)

// admRouter 登录访问
func admRouter(adm *echo.Group) {
	// 服务器信息
	adm.GET(`/sys`, control.Sys)
	// 统计信息
	adm.GET(`/collect`, control.Collect)
	// 获取当前登陆信息
	adm.GET(`/auth`, control.UserAuth)
	// 图片上传
	adm.POST(`/upload`, control.Upload)
	// 修改自身信息
	adm.POST(`/user/edit/self`, control.UserEditSelf)
	// 修改密码
	adm.POST(`/user/pass`, control.UserPass)
	// 删除分类
	adm.GET(`/cate/drop/:id`, control.CateDrop)
	// 添加分类
	adm.POST(`/cate/add`, control.CateAdd)
	// 编辑分类
	adm.POST(`/cate/edit`, control.CateEdit)
	// 文章/页面-编辑/添加
	adm.POST(`/post/opts`, control.PostOpts)
	// 删除文章/页面
	adm.GET(`/post/drop/:id`, control.PostDrop)
	// 删除标签
	adm.GET(`/tag/drop/:id`, control.TagDrop)
	// 添加标签
	adm.POST(`/tag/add`, control.TagAdd)
	// 编辑标签
	adm.POST(`/tag/edit`, control.TagEdit)
	// 编辑配置项
	adm.POST(`/opts/edit`, control.OptsEdit)
}
