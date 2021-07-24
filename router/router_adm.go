package router

import (
	"blog/control/appctl"
	"blog/control/sysctl"

	"github.com/labstack/echo/v4"
)

// admRouter 登录访问
func admRouter(adm *echo.Group) {
	adm.GET("/auth/get", sysctl.AuthGet)             // 获取当前信息
	adm.POST("/auth/edit", sysctl.AuthEdit)          // 修改自己信息
	adm.POST("/auth/passwd", sysctl.AuthPasswd)      // 修改自己密码
	adm.GET("/status/goinfo", sysctl.StatusGoinfo)   // 服务器信息
	adm.GET("/status/appinfo", sysctl.StatusAppinfo) // 统计信息
	adm.POST("/upload/file", sysctl.UploadFile)      // 文件上传
	adm.POST("/upload/image", sysctl.UploadImage)    // 图片上传
	adm.POST("/global/edit", sysctl.GlobalEdit)      // 配置修改
	adm.GET("/cate/drop/:id", appctl.CateDrop)       // 删除分类
	adm.POST("/cate/add", appctl.CateAdd)            // 添加分类
	adm.POST("/cate/edit", appctl.CateEdit)          // 编辑分类
	adm.POST("/post/opts", appctl.PostOpts)          // 文章/页面-编辑/添加
	adm.GET("/post/drop/:id", appctl.PostDrop)       // 删除文章/页面
	adm.GET("/tag/drop/:id", appctl.TagDrop)         // 删除标签
	adm.POST("/tag/add", appctl.TagAdd)              // 添加标签
	adm.POST("/tag/edit", appctl.TagEdit)            // 编辑标签
}
