package router

import (
	"blog/control/appctl"
	"blog/control/sysctl"

	"github.com/labstack/echo/v4"
)

// admRouter 登录访问
func admRouter(adm *echo.Group) {
	// sys
	adm.GET("/auth/get", sysctl.AuthGet)             // 获取当前信息
	adm.GET("/auth/grant", sysctl.AuthGrant)         // 获取当前授权
	adm.POST("/auth/edit", sysctl.AuthEdit)          // 修改自己信息
	adm.POST("/auth/passwd", sysctl.AuthPasswd)      // 修改自己密码
	adm.GET("/status/goinfo", sysctl.StatusGoinfo)   // 服务器信息
	adm.GET("/status/appinfo", sysctl.StatusAppinfo) // 统计信息
	adm.POST("/upload/file", sysctl.UploadFile)      // 文件上传
	adm.POST("/upload/image", sysctl.UploadImage)    // 图片上传
	adm.POST("/global/edit", sysctl.GlobalEdit)
	// role 角色
	adm.GET("/role/get", sysctl.RoleGet)               // 单条角色
	adm.GET("/role/all", sysctl.RoleAll)               // 所有角色
	adm.POST("/role/drop", sysctl.RoleDrop)            // 角色分页
	adm.POST("/role/add", sysctl.RoleAdd)              // 添加角色
	adm.POST("/role/edit", sysctl.RoleEdit)            // 编辑角色
	adm.GET("/role/grant/all", sysctl.RoleGrantAll)    // 所有角色授权 // role-grant
	adm.POST("/role/grant/edit", sysctl.RoleGrantEdit) // 编辑角色授权
	// grant 授权
	adm.GET("/grant/get", sysctl.GrantGet)    // 单条授权
	adm.GET("/grant/all", sysctl.GrantAll)    // 所有授权
	adm.GET("/grant/tree", sysctl.GrantTree)  // 授权树形
	adm.POST("/grant/drop", sysctl.GrantDrop) // 删除授权
	adm.POST("/grant/add", sysctl.GrantAdd)   // 添加授权
	adm.POST("/grant/edit", sysctl.GrantEdit) // 编辑授权
	// user
	adm.GET("/user/get", sysctl.UserGet)               // 单条用户
	adm.GET("/user/page", sysctl.UserPage)             // 用户分页
	adm.POST("/user/add", sysctl.UserAdd)              // 添加用户
	adm.POST("/user/drop", sysctl.UserDrop)            // 删除用户
	adm.POST("/user/edit", sysctl.UserEdit)            // 编辑用户
	adm.POST("/user/edit/lock", sysctl.UserEditLock)   // 编辑用户启用禁用
	adm.POST("/user/edit/reset", sysctl.UserEditReset) // 编辑用户重置密码
	// role
	adm.GET("/role/get", sysctl.RoleGet)       // 单条角色
	adm.GET("/role/all", sysctl.RoleAll)       // 所有角色
	adm.POST("/role/drop", sysctl.RoleDrop)    // 角色分页
	adm.POST("/role/add", sysctl.RoleAdd)      // 添加角色
	adm.POST("/role/edit", sysctl.RoleEdit)    // 编辑角色
	adm.GET("/cate/drop/:id", appctl.CateDrop) // 删除分类
	adm.POST("/cate/add", appctl.CateAdd)      // 添加分类
	adm.POST("/cate/edit", appctl.CateEdit)    // 编辑分类
	adm.POST("/post/opts", appctl.PostOpts)    // 文章/页面-编辑/添加
	adm.GET("/post/drop/:id", appctl.PostDrop) // 删除文章/页面
	adm.GET("/tag/drop/:id", appctl.TagDrop)   // 删除标签
	adm.POST("/tag/add", appctl.TagAdd)        // 添加标签
	adm.POST("/tag/edit", appctl.TagEdit)      // 编辑标签
	// menu
	adm.GET("/menu/get", sysctl.MenuGet)             // 单条菜单导航
	adm.GET("/menu/all", sysctl.MenuAll)             // 所有菜单导航
	adm.GET("/menu/tree", sysctl.MenuTree)           // 菜单导航树形
	adm.POST("/menu/drop", sysctl.MenuDrop)          // 删除菜单导航
	adm.POST("/menu/add", sysctl.MenuAdd)            // 添加菜单导航
	adm.POST("/menu/edit", sysctl.MenuEdit)          // 编辑菜单导航
	adm.POST("/menu/edit/show", sysctl.MenuEditShow) // 编辑菜单导航显隐
}
