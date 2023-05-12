package router

import (
	"github.com/labstack/echo/v4"

	"blog/internal/service/applet"
	"blog/internal/service/kernel"
)

// apiRouter
func apiRouter(engine *echo.Echo) {
	auth := engine.Group("/api")
	mgmt := engine.Group("/api")
	api := engine.Group("/api")
	auth.Use(midToken)
	mgmt.Use(midAdmin)
	api.Any("", echo.NotFoundHandler)
	api.Any("/*", echo.NotFoundHandler)
	api.GET("/auth/vcode", kernel.AuthVcode)              //验证码
	api.GET("/dict/:key", kernel.DictVal)                 //通过id获取单条字典
	api.GET("/dict/get/:key", kernel.DictGet)             //通过id获取单条字典
	api.POST("/auth/login", kernel.AuthLogin)             //登陆
	auth.GET("/auth/grant", kernel.AuthGrant)             //获取当前用户的授权
	auth.GET("/cate/get", applet.CateGet)                 //分类单条数据
	auth.GET("/cate/list", applet.CateList)               //分类列表数据
	auth.GET("/cate/page", applet.CatePage)               //获取分类分页
	auth.GET("/post/get", applet.PostGet)                 //博文单条数据
	auth.GET("/post/list", applet.PostList)               //博文列表数据
	auth.GET("/post/page", applet.PostPage)               //获取博文分页
	auth.GET("/status/app", kernel.StatusGo)              //获取服务器go信息
	auth.GET("/status/go", kernel.StatusGo)               //获取服务器go信息
	auth.GET("/tag/get", applet.TagGet)                   //标签单条数据
	auth.GET("/tag/list", applet.TagList)                 //标签列表数据
	auth.GET("/tag/page", applet.TagPage)                 //获取标签分页
	auth.POST("/cate/add", applet.CateAdd)                //分类添加数据
	auth.POST("/cate/drop", applet.CateDrop)              //分类删除数据
	auth.POST("/cate/edit", applet.CateEdit)              //分类修改数据
	auth.POST("/post/add", applet.PostAdd)                //博文添加数据
	auth.POST("/post/drop", applet.PostDrop)              //博文删除数据
	auth.POST("/post/edit", applet.PostEdit)              //博文修改数据
	auth.POST("/post/save", applet.PostSave)              //博文笔记类型保存
	auth.POST("/post/share", applet.PostShare)            //博文分享
	auth.POST("/tag/add", applet.TagAdd)                  //标签添加数据
	auth.POST("/tag/drop", applet.TagDrop)                //标签删除数据
	auth.POST("/tag/edit", applet.TagEdit)                //标签修改数据
	auth.POST("/upload/file", kernel.UploadFile)          //上传文件
	auth.POST("/upload/image", kernel.UploadImage)        //上传图片并裁剪
	mgmt.GET("/admin/all", kernel.AdminAll)               //获取admin分页数据
	mgmt.GET("/admin/exist", kernel.AdminExist)           //获取某个用户信息
	mgmt.GET("/admin/get", kernel.AdminGet)               //通过id获取admin信息
	mgmt.GET("/admin/page", kernel.AdminPage)             //获取admin分页数据
	mgmt.GET("/auth/get", kernel.AuthGet)                 //当前登录信息
	mgmt.GET("/dict/page", kernel.DictPage)               //获取字典分页
	mgmt.GET("/grant/get", kernel.GrantGet)               //通过id获取单条授权信息
	mgmt.GET("/grant/page", kernel.GrantPage)             //获取授权分页信息
	mgmt.GET("/grant/tree", kernel.GrantTree)             //获取所有授权树
	mgmt.GET("/log/page", kernel.LogPage)                 //获取日志分页
	mgmt.GET("/role/all", kernel.RoleAll)                 //获取所有角色信息
	mgmt.GET("/role/get", kernel.RoleGet)                 //通过id获取单条角色信息
	mgmt.GET("/role/grant/all", kernel.RoleGrantAll)      //通过角色id获取所有角色授权信息
	mgmt.GET("/role/page", kernel.RolePage)               //获取角色分页信息
	mgmt.POST("/admin/add", kernel.AdminAdd)              //添加admin信息
	mgmt.POST("/admin/drop", kernel.AdminDrop)            //删除admin信息
	mgmt.POST("/admin/edit", kernel.AdminEdit)            //修改admin信息
	mgmt.POST("/admin/edit/lock", kernel.AdminEditLock)   //修改用户锁定状态
	mgmt.POST("/auth/edit", kernel.AuthEdit)              //修改自己的信息
	mgmt.POST("/auth/edit/passwd", kernel.AuthEditPasswd) //修改自己的密码
	mgmt.POST("/dict/add", kernel.DictAdd)                //添加字典
	mgmt.POST("/dict/drop", kernel.DictDrop)              //通过key删除单条字典
	mgmt.POST("/dict/edit", kernel.DictEdit)              //修改字典
	mgmt.POST("/grant/add", kernel.GrantAdd)              //添加授权信息
	mgmt.POST("/grant/drop", kernel.GrantDrop)            //通过id删除单条授权信息
	mgmt.POST("/grant/edit", kernel.GrantEdit)            //修改授权信息
	mgmt.POST("/role/add", kernel.RoleAdd)                //添加角色信息
	mgmt.POST("/role/drop", kernel.RoleDrop)              //通过id删除单条角色信息
	mgmt.POST("/role/edit", kernel.RoleEdit)              //修改角色信息
	mgmt.POST("/role/grant/edit", kernel.RoleGrantEdit)   //修改角色授权信息
}
