package router

import (
	"blog/control/appctl"
	"blog/control/sysctl"

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
	// adm.GET("/role/page", sysctl.RolePage)              //角色分页
	// adm.GET("/role/all", sysctl.RoleAll)                //所有角色
	// adm.GET("/role/drop/:id", sysctl.RoleDrop)          //角色分页
	// adm.POST("/role/add", sysctl.RoleAdd)               //添加角色
	// adm.POST("/role/edit", sysctl.RoleEdit)             //编辑角色
	adm.GET("/api/page", sysctl.ApiPage)                //接口分页
	adm.GET("/api/all", sysctl.ApiAll)                  //所有接口
	adm.GET("/api/drop/:id", sysctl.ApiDrop)            //接口分页
	adm.POST("/api/add", sysctl.ApiAdd)                 //添加接口
	adm.POST("/api/edit", sysctl.ApiEdit)               //编辑接口
	adm.GET("/menu/all", sysctl.MenuAll)                //所有菜单导航
	adm.POST("/menu/drop", sysctl.MenuDrop)             //菜单导航分页
	adm.POST("/menu/add", sysctl.MenuAdd)               //添加菜单导航
	adm.POST("/menu/edit", sysctl.MenuEdit)             //编辑菜单导航
	adm.POST("/menu/edit/show", sysctl.MenuEditShow)    //编辑菜单导航显隐
	adm.GET("/role/api/page", sysctl.RoleApiPage)       //角色接口分页
	adm.GET("/role/api/all", sysctl.RoleApiAll)         //所有角色接口
	adm.GET("/role/api/drop/:id", sysctl.RoleApiDrop)   //角色接口分页
	adm.POST("/role/api/add", sysctl.RoleApiAdd)        //添加角色接口
	adm.POST("/role/api/edit", sysctl.RoleApiEdit)      //编辑角色接口
	adm.GET("/role/menu/page", sysctl.RoleMenuPage)     //角色菜单导航分页
	adm.GET("/role/menu/all", sysctl.RoleMenuAll)       //所有角色菜单导航
	adm.GET("/role/menu/drop/:id", sysctl.RoleMenuDrop) //角色菜单导航分页
	adm.POST("/role/menu/add", sysctl.RoleMenuAdd)      //添加角色菜单导航
	adm.POST("/role/menu/edit", sysctl.RoleMenuEdit)    //编辑角色菜单导航
}
