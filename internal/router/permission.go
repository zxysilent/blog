package router

import (
	"blog/conf"
	"blog/internal/model"
	"blog/internal/repo"
	"context"
	"sort"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/logs"
)

// 权限元信息直接编码在路由 Name 中，注册路由处即为唯一来源，平铺书写：
//
//	auth.GET("/users", kernel.UserList).Name = "user:list|100|user|用户管理|用户列表"
//
// Name 由 5 段组成：权限标识|分组排序|分组键|分组中文名|权限中文名。
// 未设置 Name 或切分后不足 5 段的路由只校验登录态，不参与权限同步。
// echo 会为未显式设置 Name 的路由自动填入 handler 函数名（含包路径），
// 其切分后不足 5 段，天然被过滤掉。
//
// 每个分组会自动派生一个「分组授权权限」<分组键>:grant，
// 例如 user -> user:grant，修改该分组下任意权限都需要先拥有它。
const permSep = "|"

// permParts Name 切分后的段数
const permParts = 5

// splitPermName 按分隔符切分路由 Name，段数不符时返回 false
func splitPermName(name string) ([permParts]string, bool) {
	var out [permParts]string
	parts := strings.Split(name, permSep)
	if len(parts) != permParts {
		return out, false
	}
	for i := 0; i < permParts; i++ {
		out[i] = strings.TrimSpace(parts[i])
		if out[i] == "" {
			return out, false
		}
	}
	return out, true
}

// parsePermName 把路由 Name 解析为入库定义
func parsePermName(name string) (model.PermInitArgs, bool) {
	parts, ok := splitPermName(name)
	if !ok {
		return model.PermInitArgs{}, false
	}
	sortNum, err := strconv.Atoi(parts[1])
	if err != nil {
		return model.PermInitArgs{}, false
	}
	return model.PermInitArgs{
		Perm:     parts[0],
		Sort:     sortNum,
		GroupKey: parts[2],
		Group:    parts[3],
		Name:     parts[4],
	}, true
}

// collectPerms 补齐分组授权权限，并按分组排序 + 权限标识排序
// 分组授权权限由分组自动派生，无需在每个路由上重复声明
func collectPerms(defs []model.PermInitArgs) []model.PermInitArgs {
	seen := make(map[string]struct{}, len(defs)+8)
	groups := make(map[string]model.PermInitArgs, 8)
	for _, item := range defs {
		seen[item.Perm] = struct{}{}
		group, ok := groups[item.GroupKey]
		// 取分组内最小排序作为分组授权权限的排序，保证它排在同组最前
		if !ok || item.Sort < group.Sort {
			groups[item.GroupKey] = item
		}
	}
	for key, group := range groups {
		perm := model.GroupGrantPerm(key)
		if _, ok := seen[perm]; ok {
			continue // 已显式声明，不重复派生
		}
		seen[perm] = struct{}{}
		defs = append(defs, model.PermInitArgs{
			Perm:     perm,
			Name:     group.Group + "授权",
			GroupKey: key,
			Group:    group.Group,
			Sort:     group.Sort,
		})
	}
	sort.Slice(defs, func(i, j int) bool {
		if defs[i].Sort != defs[j].Sort {
			return defs[i].Sort < defs[j].Sort
		}
		if defs[i].Group != defs[j].Group {
			return defs[i].Group < defs[j].Group
		}
		return defs[i].Perm < defs[j].Perm
	})
	return defs
}

// routePerms 路由索引：path -> method -> 权限标识
// routeNames 路由索引：path -> method -> 原始 Name（审计日志取分组与中文名）
// 启动注册完路由后由 BuildRouteNames 一次性构建，之后只读，故无需加锁。
// 请求期查询走两次 map 查找，无字符串拼接、无分配。
var (
	routePerms map[string]map[string]string
	routeNames map[string]map[string]string
)

// BuildRouteNames 构建路由索引，并把权限点同步入库
// 必须在 apiRouter 注册完成后调用
func BuildRouteNames(engine *echo.Echo) {
	routes := engine.Routes()
	// Routes() 内部是 map 遍历，顺序不稳定；
	// 先按 方法+路径 排序，保证同一权限的去重结果在多次启动间一致
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Path != routes[j].Path {
			return routes[i].Path < routes[j].Path
		}
		return routes[i].Method < routes[j].Method
	})
	routePerms = make(map[string]map[string]string, len(routes))
	routeNames = make(map[string]map[string]string, len(routes))
	seen := make(map[string]struct{}, 64)
	defs := make([]model.PermInitArgs, 0, 64)
	for _, r := range routes {
		// 原始 Name 全部索引下来，审计日志的非权限路由（如登录）也能取到
		names := routeNames[r.Path]
		if names == nil {
			names = make(map[string]string, 2)
			routeNames[r.Path] = names
		}
		names[r.Method] = r.Name
		def, ok := parsePermName(r.Name)
		if !ok {
			continue // 未声明权限的路由，只校验登录态
		}
		// 索引只保留剥离后的权限标识，供请求期直接相等判断
		methods := routePerms[r.Path]
		if methods == nil {
			methods = make(map[string]string, 2)
			routePerms[r.Path] = methods
		}
		methods[r.Method] = def.Perm
		// 同一权限可控制多个路由，按权限标识去重后入库
		if _, ok := seen[def.Perm]; ok {
			continue
		}
		seen[def.Perm] = struct{}{}
		defs = append(defs, def)
	}
	ctx := logs.TraceCtx(context.Background(), logs.TraceId())
	if err := repo.PermSync(ctx, collectPerms(defs)); err != nil {
		logs.Ctx(ctx).Err(err).Int("count", len(defs)).Error("repo.PermSync")
		return
	}
	// 权限点可能新增/下线，重建角色权限缓存
	repo.SyncRolePerms()
	logs.Ctx(ctx).Info("router.permissions.synced")
}

// routePerm 取当前请求匹配到的路由权限标识
// ctx.Path() 返回 echo 匹配到的路由模式（如 /api/users/:id），未声明权限时返回空串
func routePerm(ctx echo.Context) string {
	methods := routePerms[ctx.Path()]
	if methods == nil {
		return ""
	}
	return methods[ctx.Request().Method]
}

// extractToken 从 header / query / form / cookie 中提取 token
func extractToken(ctx echo.Context) string {
	if raw := ctx.Request().Header.Get(echo.HeaderAuthorization); raw != "" {
		return trimBearer(raw)
	}
	// 管理端（mgmt）按惯例把令牌放在 token 头里，这里一并支持
	if raw := ctx.Request().Header.Get(conf.TokenKey); raw != "" {
		return trimBearer(raw)
	}
	if raw := ctx.FormValue(conf.TokenKey); raw != "" {
		return trimBearer(raw)
	}
	if cookie, err := ctx.Cookie(conf.TokenKey); err == nil && cookie != nil {
		return trimBearer(cookie.Value)
	}
	return ""
}

// trimBearer 去掉可选的 Bearer 前缀
// 有的客户端习惯写 Authorization: Bearer <token>，令牌本身不含空格，
// 因此这里宽松处理，两种写法都能用
func trimBearer(raw string) string {
	raw = strings.TrimSpace(raw)
	if len(raw) > 7 && strings.EqualFold(raw[:7], "Bearer ") {
		return strings.TrimSpace(raw[7:])
	}
	return raw
}
