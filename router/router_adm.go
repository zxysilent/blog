package router

import (
	"blog/control"

	"github.com/labstack/echo/v4"
)

// admRouter 登录访问
func admRouter(adm *echo.Group) {
	adm.GET(`/sys`, control.Sys)                      // 服务器信息
	adm.GET(`/collect`, control.Collect)              // 统计信息
	adm.GET(`/auth`, control.UserAuth)                // 获取当前登陆信息
	adm.POST(`/upload`, control.Upload)               // 图片上传
	adm.POST(`/user/edit/self`, control.UserEditSelf) // 修改自身信息
	adm.POST(`/user/pass`, control.UserPass)          // 修改密码
	adm.GET(`/cate/drop/:id`, control.CateDrop)       // 删除分类
	adm.POST(`/cate/add`, control.CateAdd)            // 添加分类
	adm.POST(`/cate/edit`, control.CateEdit)          // 编辑分类
	adm.POST(`/post/opts`, control.PostOpts)          // 文章/页面-编辑/添加
	adm.GET(`/post/drop/:id`, control.PostDrop)       // 删除文章/页面
	adm.GET(`/tag/drop/:id`, control.TagDrop)         // 删除标签
	adm.POST(`/tag/add`, control.TagAdd)              // 添加标签
	adm.POST(`/tag/edit`, control.TagEdit)            // 编辑标签
	adm.POST(`/opts/edit`, control.OptsEdit)          // 编辑配置项
}
