package router

import (
	"blog/conf"
	"blog/control"
	"blog/control/sysctl"
	"log"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// RunApp 入口
func RunApp() {
	engine := echo.New()
	engine.Renderer = initRender()                    // 初始渲染引擎
	engine.Use(midRecover, midLogger)                 // 恢复 日志记录
	engine.Use(middleware.CORSWithConfig(crosConfig)) // 跨域设置
	engine.HideBanner = true                          // 不显示横幅
	engine.HTTPErrorHandler = HTTPErrorHandler        // 自定义错误处理
	engine.Debug = conf.App.IsDev()                   // 运行模式 - echo框架好像没怎么使用这个
	RegDocs(engine)                                   // 注册文档
	engine.Static(`/dist`, "dist")                    // 静态目录 - 后端专用
	engine.Static(`/static`, "static")                // 静态目录
	engine.File(`/favicon.ico`, "favicon.ico")        // ico
	engine.File("/dashboard*", "dist/index.html")     // 前后端分离页面
	//--- 页面 -- start
	engine.GET(`/`, control.IndexView)                 // 首页
	engine.GET(`/archives`, control.ArchivesView)      // 归档
	engine.GET(`/archives.json`, control.ArchivesJson) // 归档 json
	engine.GET(`/tags`, control.TagsView)              // 标签
	engine.GET(`/tags.json`, control.TagsJson)         // 标签 json
	engine.GET(`/tag/:tag`, control.TagPostView)       // 具体某个标签
	engine.GET(`/cate/:cate`, control.CatePostView)    // 分类
	engine.GET(`/about`, control.AboutView)            // 关于
	engine.GET(`/links`, control.LinksView)            // 友链
	engine.GET(`/post/*`, control.PostView)            // 具体某个文章
	engine.GET(`/page/*`, control.PageView)            // 具体某个页面
	//--- 页面 -- end
	api := engine.Group("/api")          // api/
	apiRouter(api)                       // 注册分组路由
	adm := engine.Group("/adm", midAuth) // adm/ 需要登陆才能访问
	admRouter(adm)                       // 注册分组路由
	sys := engine.Group("/sys", midAuth) // sys/ 需要登陆才能访问
	sysRouter(sys)                       // 注册分组路由
	err := engine.Start(conf.App.Addr)
	if err != nil {
		log.Fatalln("run error :", err)
	}
}

// sysRouter RBAC权限
func sysRouter(sys *echo.Group) {
	sys.GET("/role/page", sysctl.SysRolePage)              //"角色分页"
	sys.GET("/role/all", sysctl.SysRoleAll)                //"所有角色"
	sys.GET("/role/drop/:id", sysctl.SysRoleDrop)          //"角色分页"
	sys.POST("/role/add", sysctl.SysRoleAdd)               //"添加角色"
	sys.POST("/role/edit", sysctl.SysRoleEdit)             //"编辑角色"
	sys.GET("/auth/page", sysctl.SysAuthPage)              //"认证分页"
	sys.GET("/auth/all", sysctl.SysAuthAll)                //"所有认证"
	sys.GET("/auth/drop/:id", sysctl.SysAuthDrop)          //"认证分页"
	sys.POST("/auth/add", sysctl.SysAuthAdd)               //"添加认证"
	sys.POST("/auth/edit", sysctl.SysAuthEdit)             //"编辑认证"
	sys.GET("/menu/all", sysctl.SysMenuAll)                //"所有菜单导航"
	sys.GET("/menu/drop/:id", sysctl.SysMenuDrop)          //"菜单导航分页"
	sys.POST("/menu/add", sysctl.SysMenuAdd)               //"添加菜单导航"
	sys.POST("/menu/edit", sysctl.SysMenuEdit)             //"编辑菜单导航"
	sys.GET("/role/auth/page", sysctl.SysRoleAuthPage)     //"角色认证分页"
	sys.GET("/role/auth/all", sysctl.SysRoleAuthAll)       //"所有角色认证"
	sys.GET("/role/auth/drop/:id", sysctl.SysRoleAuthDrop) //"角色认证分页"
	sys.POST("/role/auth/add", sysctl.SysRoleAuthAdd)      //"添加角色认证"
	sys.POST("/role/auth/edit", sysctl.SysRoleAuthEdit)    //"编辑角色认证"
	sys.GET("/role/menu/page", sysctl.SysRoleMenuPage)     //"角色菜单导航分页"
	sys.GET("/role/menu/all", sysctl.SysRoleMenuAll)       //"所有角色菜单导航"
	sys.GET("/role/menu/drop/:id", sysctl.SysRoleMenuDrop) //"角色菜单导航分页"
	sys.POST("/role/menu/add", sysctl.SysRoleMenuAdd)      //"添加角色菜单导航"
	sys.POST("/role/menu/edit", sysctl.SysRoleMenuEdit)    //"编辑角色菜单导航"
}
