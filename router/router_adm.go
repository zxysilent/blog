package router

import (
	"blog/control/appctl"

	"github.com/labstack/echo/v4"
)

// admRouter 登录访问
func admRouter(adm *echo.Group) {
	adm.GET(`/sys`, appctl.Sys)                      // 服务器信息
	adm.GET(`/collect`, appctl.Collect)              // 统计信息
	adm.GET(`/auth`, appctl.UserAuth)                // 获取当前登陆信息
	adm.POST(`/upload`, appctl.Upload)               // 图片上传
	adm.POST(`/user/edit/self`, appctl.UserEditSelf) // 修改自身信息
	adm.POST(`/user/pass`, appctl.UserPass)          // 修改密码
	adm.GET(`/cate/drop/:id`, appctl.CateDrop)       // 删除分类
	adm.POST(`/cate/add`, appctl.CateAdd)            // 添加分类
	adm.POST(`/cate/edit`, appctl.CateEdit)          // 编辑分类
	adm.POST(`/post/opts`, appctl.PostOpts)          // 文章/页面-编辑/添加
	adm.GET(`/post/drop/:id`, appctl.PostDrop)       // 删除文章/页面
	adm.GET(`/tag/drop/:id`, appctl.TagDrop)         // 删除标签
	adm.POST(`/tag/add`, appctl.TagAdd)              // 添加标签
	adm.POST(`/tag/edit`, appctl.TagEdit)            // 编辑标签
	adm.POST(`/opts/edit`, appctl.OptsEdit)          // 编辑配置项
	//sysctl

}
