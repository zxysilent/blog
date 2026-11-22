package router

import (
	"blog/internal/model"
	"blog/internal/service/kernel"

	"github.com/labstack/echo/v4"
)

// apiRouter 注册 RESTful API 路由
// 只有两组：api 公开、auth 登录态
//
// 权限声明：路由 Name 平铺写完整信息
//
// 权限标识|分组排序|分组键|分组中文名|权限中文名
//
// 启动时 BuildRouteNames 通过 engine.Routes() 汇总入库，无需在别处维护。
// 每个分组会自动派生 <分组键>:grant 授权权限，用于控制该分组下的授权行为。
// 不设置 Name 的路由只校验登录态，不参与权限同步。
func apiRouter(engine *echo.Echo) {
	api := engine.Group("/api")
	auth := engine.Group("/api")
	auth.Use(AuthToken)
	// Group.Use 会为分组注册 404 兜底路由，未匹配的 /api/* 会先被 auth 分组
	// 的 AuthToken 拦成「请重新登陆」。在公开分组挂 NotFoundHandler 挡在前面，
	// 让未匹配的接口直接落到全局错误处理，统一返回「接口不存在」。
	api.Any("", echo.NotFoundHandler)
	api.Any("/*", echo.NotFoundHandler)

	// ---------------- 公开接口（无需登录） ----------------
	api.GET("/captcha", kernel.AuthCaptcha)           //图形验证码
	api.POST("/tokens", kernel.AuthLogin)             //登录-创建令牌
	api.GET("/site", kernel.Site)                     //站点公开配置
	api.GET("/dicts/:id/entries", kernel.DictEntries) //字典选项-匿名可读
	api.GET("/configs/:id/val", kernel.ConfigVal)     //公开配置项-仅限 basic/blog/login

	// ---------------- 当前登录账号（无需权限） ----------------
	auth.GET("/auth", kernel.AuthMe)                   //当前登录用户
	auth.PUT("/auth", kernel.AuthEditMe)               //修改自己资料
	auth.PUT("/auth/password", kernel.AuthEditPasswd)  //修改自己密码
	auth.POST("/auth/verify", kernel.AuthVerifyPasswd) //校验自己密码-锁屏解锁
	auth.GET("/auth/permissions", kernel.AuthPerms)    //我的权限标识

	// ---------------- 用户 ----------------
	auth.GET("/users", kernel.UserList).Name = "user:list|100|user|用户管理|用户列表"
	auth.POST("/users", kernel.UserAdd).Name = "user:create|100|user|用户管理|用户新增"
	auth.GET("/users/:id", kernel.UserGet).Name = "user:read|100|user|用户管理|用户详情"
	auth.PUT("/users/:id", kernel.UserEdit).Name = "user:update|100|user|用户管理|用户修改"
	auth.DELETE("/users/:id", kernel.UserDrop).Name = "user:delete|100|user|用户管理|用户删除"

	// ---------------- 角色 ----------------
	auth.GET("/roles", kernel.RoleList).Name = "role:list|200|role|角色管理|角色列表"
	auth.POST("/roles", kernel.RoleAdd).Name = "role:create|200|role|角色管理|角色新增"
	auth.GET("/roles/:id", kernel.RoleGet).Name = "role:read|200|role|角色管理|角色详情"
	auth.PUT("/roles/:id", kernel.RoleEdit).Name = "role:update|200|role|角色管理|角色修改"
	auth.DELETE("/roles/:id", kernel.RoleDrop).Name = "role:delete|200|role|角色管理|角色删除"
	auth.GET("/roles/:id/permissions", kernel.RolePermGet).Name = "role:assign|200|role|角色管理|角色授权"
	auth.PUT("/roles/:id/permissions", kernel.RolePermEdit).Name = "role:assign|200|role|角色管理|角色授权"

	// ---------------- 权限点（由路由自动同步，只读） ----------------
	auth.GET("/permissions", kernel.PermList).Name = "permission:list|300|permission|权限管理|权限列表"
	auth.GET("/permissions/tree", kernel.PermTree).Name = "permission:tree|300|permission|权限管理|权限树"
	auth.GET("/permissions/:id", kernel.PermGet).Name = "permission:read|300|permission|权限管理|权限详情"

	// ---------------- 字典 ----------------
	auth.GET("/dicts", kernel.DictList).Name = "dict:list|400|dict|字典管理|字典列表"
	auth.POST("/dicts", kernel.DictAdd).Name = "dict:create|400|dict|字典管理|字典新增"
	auth.GET("/dicts/:id", kernel.DictGet).Name = "dict:read|400|dict|字典管理|字典详情"
	auth.PUT("/dicts/:id", kernel.DictEdit).Name = "dict:update|400|dict|字典管理|字典修改"
	auth.DELETE("/dicts/:id", kernel.DictDrop).Name = "dict:delete|400|dict|字典管理|字典删除"

	// ---------------- 配置 ----------------
	auth.GET("/configs", kernel.ConfigList).Name = "config:list|500|config|配置管理|配置列表"
	auth.POST("/configs", kernel.ConfigAdd).Name = "config:create|500|config|配置管理|配置新增"
	auth.GET("/configs/groups", kernel.ConfigGroups).Name = "config:list|500|config|配置管理|配置列表"
	// 内置配置的独立入口：一项配置一个权限点，便于按配置分权
	// 只给运营角色「站点基础配置」权限，它就改不了认证私钥与上传上限。
	// 静态路径在 echo 中优先于 /configs/:id 匹配，两条路由可共存。
	auth.GET("/configs/basic", kernel.ConfigBasicGet).Name = "config:basic|500|config|配置管理|基础配置"
	auth.PUT("/configs/basic", kernel.ConfigBasicEdit).Name = "config:basic|500|config|配置管理|基础配置"
	auth.GET("/configs/login", kernel.ConfigLoginGet).Name = "config:login|500|config|配置管理|登录页配置"
	auth.PUT("/configs/login", kernel.ConfigLoginEdit).Name = "config:login|500|config|配置管理|登录页配置"
	auth.GET("/configs/blog", kernel.ConfigBlogGet).Name = "config:blog|500|config|配置管理|博文配置"
	auth.PUT("/configs/blog", kernel.ConfigBlogEdit).Name = "config:blog|500|config|配置管理|博文配置"
	auth.GET("/configs/auth", kernel.ConfigAuthGet).Name = "config:auth|500|config|配置管理|认证配置"
	auth.PUT("/configs/auth", kernel.ConfigAuthEdit).Name = "config:auth|500|config|配置管理|认证配置"
	auth.GET("/configs/upload", kernel.ConfigUploadGet).Name = "config:upload|500|config|配置管理|上传配置"
	auth.PUT("/configs/upload", kernel.ConfigUploadEdit).Name = "config:upload|500|config|配置管理|上传配置"
	auth.GET("/configs/layout", kernel.ConfigLayoutGet) // 管理端外壳布局，登录可读（不单独授权，避免每次进后台都报无权限）
	// 布局配置含登录页风格（login_style）：改该字段还需 login:style 权限，
	// 由 ConfigEdit 内部校验，这里不额外拆路由
	auth.PUT("/configs/layout", kernel.ConfigLayoutEdit).Name = "config:layout|500|config|配置管理|布局配置"
	auth.GET("/configs/sql", kernel.ConfigSqlGet).Name = "config:sql|500|config|配置管理|SQL 日志配置"
	auth.PUT("/configs/sql", kernel.ConfigSqlEdit).Name = "config:sql|500|config|配置管理|SQL 日志配置"
	// 登录页风格（存于 layout 配置的 login_style）：站点门面，单独授权。
	// 单独成路由是为了让 login:style 这个权限点存在且可授予；
	// 通用 /configs/layout 改动该字段会被 ConfigEdit 拦下，只有这里能写。
	auth.PUT("/configs/login/style", kernel.ConfigLoginStyle).Name = model.PermLoginStyle + "|500|config|配置管理|登录页风格"
	// 通用入口：自建配置只能用这条；内置配置同样受 config:<id> 保护
	auth.GET("/configs/:id", kernel.ConfigGet).Name = "config:read|500|config|配置管理|配置详情"
	auth.PUT("/configs/:id", kernel.ConfigEdit).Name = "config:update|500|config|配置管理|配置修改"
	auth.DELETE("/configs/:id", kernel.ConfigDrop).Name = "config:delete|500|config|配置管理|配置删除"

	// ---------------- 内容 ----------------
	auth.GET("/posts", kernel.PostList).Name = "post:list|600|content|内容管理|内容列表"
	auth.POST("/posts", kernel.PostAdd).Name = "post:create|600|content|内容管理|内容新增"
	auth.GET("/posts/tree", kernel.PostTree).Name = "post:list|600|content|内容管理|内容列表"
	auth.GET("/posts/:id", kernel.PostGet).Name = "post:read|600|content|内容管理|内容详情"
	auth.PUT("/posts/:id", kernel.PostEdit).Name = "post:update|600|content|内容管理|内容修改"
	auth.PUT("/posts/:id/top", kernel.PostEditTop).Name = "post:update|600|content|内容管理|内容修改"
	auth.DELETE("/posts/:id", kernel.PostDrop).Name = "post:delete|600|content|内容管理|内容删除"
	auth.GET("/posts/:id/tags", kernel.PostTags).Name = "post:read|600|content|内容管理|内容详情"
	auth.PUT("/posts/:id/tags", kernel.PostTagEdit).Name = "post:update|600|content|内容管理|内容修改"
	auth.GET("/posts/:id/neighbors", kernel.PostNeighbors).Name = "post:read|600|content|内容管理|内容详情"

	// ---------------- 栏目 ----------------
	auth.GET("/cates", kernel.CateList).Name = "cate:list|600|content|内容管理|栏目列表"
	auth.POST("/cates", kernel.CateAdd).Name = "cate:create|600|content|内容管理|栏目新增"
	auth.GET("/cates/tree", kernel.CateTree).Name = "cate:tree|600|content|内容管理|栏目树"
	auth.GET("/cates/:id", kernel.CateGet).Name = "cate:read|600|content|内容管理|栏目详情"
	auth.PUT("/cates/:id", kernel.CateEdit).Name = "cate:update|600|content|内容管理|栏目修改"
	auth.DELETE("/cates/:id", kernel.CateDrop).Name = "cate:delete|600|content|内容管理|栏目删除"

	// ---------------- 标签 ----------------
	auth.GET("/tags", kernel.TagList).Name = "tag:list|600|content|内容管理|标签列表"
	auth.POST("/tags", kernel.TagAdd).Name = "tag:create|600|content|内容管理|标签新增"
	auth.GET("/tags/:id", kernel.TagGet).Name = "tag:read|600|content|内容管理|标签详情"
	auth.PUT("/tags/:id", kernel.TagEdit).Name = "tag:update|600|content|内容管理|标签修改"
	auth.DELETE("/tags/:id", kernel.TagDrop).Name = "tag:delete|600|content|内容管理|标签删除"

	// ---------------- 系统 ----------------
	auth.GET("/status", kernel.Status).Name = "status:read|700|system|系统|服务器状态"
	auth.POST("/files", kernel.UploadFile).Name = "file:create|700|system|系统|文件上传"
	auth.POST("/images", kernel.UploadImage).Name = "file:create|700|system|系统|图片上传"

	// ---------------- 操作日志 ----------------
	// 日志由审计中间件自动写入，这里只提供查询与清理
	auth.GET("/logs", kernel.LogList).Name = "log:list|800|log|日志管理|操作日志"
	auth.DELETE("/logs", kernel.LogClear).Name = "log:delete|800|log|日志管理|日志清理"
	auth.GET("/logs/:id", kernel.LogGet).Name = "log:read|800|log|日志管理|日志详情"
	auth.DELETE("/logs/:id", kernel.LogDrop).Name = "log:delete|800|log|日志管理|日志删除"
}
