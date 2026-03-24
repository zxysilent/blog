package router

import (
	"github.com/labstack/echo/v4"

	"blog/internal/service/kernel"
)

// apiRouter
func apiRouter(engine *echo.Echo) {
	auth := engine.Group("/api")
	mgmt := engine.Group("/api")
	api := engine.Group("/api")
	auth.Use(AuthToken)
	mgmt.Use(MgmtToken)
	api.Any("", echo.NotFoundHandler)
	api.Any("/*", echo.NotFoundHandler)
	api.GET("/auth/vcode", kernel.AuthVcode)              //验证码
	api.GET("/dict/get/:key", kernel.DictGet)             //通过id获取单条字典
	api.POST("/auth/login", kernel.AuthLogin)             //登陆
	auth.GET("/auth/grant", kernel.AuthGrant)             //获取当前用户的授权
	auth.GET("/status/app", kernel.StatusGo)              //获取服务器go信息
	auth.GET("/status/go", kernel.StatusGo)               //获取服务器go信息
	auth.POST("/upload/file", kernel.UploadFile)          //上传文件
	auth.POST("/upload/image", kernel.UploadImage)        //上传图片并裁剪
	mgmt.DELETE("/dict/drop", kernel.DictDrop)            //通过key删除单条字典
	mgmt.DELETE("/grant/:id", kernel.GrantDrop)           //通过id删除单条授权信息
	mgmt.GET("/admin/all", kernel.AdminAll)               //获取admin分页数据
	mgmt.GET("/admin/page", kernel.AdminPage)             //获取admin分页数据
	mgmt.GET("/admin/:id", kernel.AdminGet)               //通过id获取admin信息
	mgmt.GET("/admin/:num/exist", kernel.AdminExist)      //获取某个用户信息
	mgmt.GET("/auth/get", kernel.AuthGet)                 //当前登录信息
	mgmt.GET("/dict/page", kernel.DictList)               //获取字典分页
	mgmt.GET("/grant/page", kernel.GrantPage)             //获取授权分页信息
	mgmt.GET("/grant/:id", kernel.GrantGet)               //通过id获取单条授权信息
	mgmt.GET("/role/all", kernel.RoleAll)                 //获取所有角色信息
	mgmt.GET("/role/grant/all", kernel.RoleGrantAll)      //通过角色id获取所有角色授权信息
	mgmt.GET("/role/page", kernel.RolePage)               //获取角色分页信息
	mgmt.GET("/role/:id", kernel.RoleGet)                 //通过id获取单条角色信息
	mgmt.POST("/admin/add", kernel.AdminAdd)              //添加admin信息
	mgmt.POST("/admin/drop", kernel.AdminDrop)            //删除admin信息
	mgmt.POST("/admin/edit", kernel.AdminEdit)            //修改admin信息
	mgmt.POST("/admin/edit/lock", kernel.AdminEditLock)   //修改用户锁定状态
	mgmt.POST("/auth/edit", kernel.AuthEdit)              //修改自己的信息
	mgmt.POST("/auth/edit/passwd", kernel.AuthEditPasswd) //修改自己的密码
	mgmt.POST("/dict/add", kernel.DictAdd)                //添加字典
	mgmt.POST("/grant/add", kernel.GrantAdd)              //添加授权信息
	mgmt.POST("/grant/edit", kernel.GrantEdit)            //修改授权信息
	mgmt.POST("/role/add", kernel.RoleAdd)                //添加角色信息
	mgmt.POST("/role/drop", kernel.RoleDrop)              //通过id删除单条角色信息
	mgmt.POST("/role/edit", kernel.RoleEdit)              //修改角色信息
	mgmt.POST("/role/grant/edit", kernel.RoleGrantEdit)   //修改角色授权信息
	mgmt.PUT("/dict/edit", kernel.DictEdit)               //修改字典
}
