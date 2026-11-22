package router

import (
	"context"
	"fmt"
	"io"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"blog/conf"
	"blog/internal/model"
	"blog/internal/repo"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/logs"
)

// init 在任何测试（以及它们启动的 goroutine）之前导走日志输出
// 放到 setup 里会与 repo 日志协程并发读写 logs 的全局配置，触发数据竞争
func init() {
	logs.SetOutput(io.Discard)
}

// setupRouterDB 初始化隔离的内存库，供权限同步使用
func setupRouterDB(t *testing.T) {
	t.Helper()
	name := strings.NewReplacer("/", "_", " ", "_", "\\", "_").Replace(t.Name())
	dsn := fmt.Sprintf("file:%s?mode=memory&cache=shared", name)
	path := filepath.Join(t.TempDir(), "conf.toml")
	body := fmt.Sprintf(`title = "test"
mode = "test"
port = 0
host = "127.0.0.1"
token_exp = 1
token_secret = "test-secret"

[db]
type = "sqlite"
sqlite = %q
`, dsn)
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("write conf: %v", err)
	}
	conf.Init(context.Background(), path)
	repo.Init(context.Background())
	t.Cleanup(func() { _ = repo.GetDb().Close() })
}

// buildEngine 注册 API 路由并构建权限索引
func buildEngine(t *testing.T) *echo.Echo {
	t.Helper()
	setupRouterDB(t)
	engine := echo.New()
	apiRouter(engine)
	BuildRouteNames(engine)
	return engine
}

// TestSplitPermName 路由 Name 必须是 5 段
func TestSplitPermName(t *testing.T) {
	ok := []string{
		"user:list|100|user|用户管理|用户列表",
		"post:read|600|content|内容管理|内容详情",
	}
	for _, it := range ok {
		if _, valid := splitPermName(it); !valid {
			t.Errorf("splitPermName(%q) 应通过", it)
		}
	}
	// 4 段（旧格式）、echo 自动生成的 handler 名、空串都应被拒绝
	bad := []string{
		"user:list|100|用户管理|用户列表",
		"blog/internal/service/kernel.AuthMe",
		"",
		"a|1|b|c|",
	}
	for _, it := range bad {
		if _, valid := splitPermName(it); valid {
			t.Errorf("splitPermName(%q) 应被拒绝", it)
		}
	}
}

// TestParsePermNameRejectsBadSort 段数够但排序不是数字时也不解析
func TestParsePermNameRejectsBadSort(t *testing.T) {
	for _, it := range []string{"a|x|b|c|d", "a|1|b|c|d|e"} {
		if _, ok := parsePermName(it); ok {
			t.Errorf("parsePermName(%q) 应被拒绝", it)
		}
	}
}

// TestParsePermName 解析出的字段顺序正确
func TestParsePermName(t *testing.T) {
	def, ok := parsePermName("post:list|600|content|内容管理|内容列表")
	if !ok {
		t.Fatal("解析失败")
	}
	want := model.PermInitArgs{
		Perm:     "post:list",
		Sort:     600,
		GroupKey: "content",
		Group:    "内容管理",
		Name:     "内容列表",
	}
	if def != want {
		t.Errorf("解析结果 = %+v, want %+v", def, want)
	}
}

// TestCollectPermsAddsGroupGrant 每个分组自动派生一个授权权限
func TestCollectPermsAddsGroupGrant(t *testing.T) {
	defs := []model.PermInitArgs{
		{Perm: "user:list", Name: "用户列表", GroupKey: "user", Group: "用户管理", Sort: 100},
		{Perm: "user:create", Name: "用户新增", GroupKey: "user", Group: "用户管理", Sort: 100},
		{Perm: "post:list", Name: "内容列表", GroupKey: "content", Group: "内容管理", Sort: 600},
	}
	out := collectPerms(defs)
	seen := make(map[string]bool, len(out))
	for _, it := range out {
		seen[it.Perm] = true
		if it.GroupKey == "" || it.Group == "" {
			t.Errorf("%q 缺少分组信息", it.Perm)
		}
	}
	for _, want := range []string{"user:grant", "content:grant"} {
		if !seen[want] {
			t.Errorf("缺少分组授权权限 %q，got %v", want, seen)
		}
	}
	// 同组只派生一次
	count := 0
	for _, it := range out {
		if it.Perm == "user:grant" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("user:grant 派生次数 = %d, want 1", count)
	}
}

// TestApiRouteNamesAllParsable 所有声明了权限的路由都必须能被解析
// 这类错误会让权限点静默丢失，只能在运行时表现为「没有权限」
func TestApiRouteNamesAllParsable(t *testing.T) {
	engine := buildEngine(t)
	total := 0
	for _, r := range engine.Routes() {
		// 未显式设置 Name 的路由由 echo 自动填入 handler 函数名，天然不含分隔符
		if !strings.Contains(r.Name, permSep) {
			continue
		}
		def, ok := parsePermName(r.Name)
		if !ok {
			t.Errorf("路由 %s %s 的 Name 无法解析: %q", r.Method, r.Path, r.Name)
			continue
		}
		total++
		if def.Perm == "" || def.GroupKey == "" || def.Group == "" || def.Name == "" {
			t.Errorf("路由 %s %s 解析结果不完整: %+v", r.Method, r.Path, def)
		}
	}
	if total < 30 {
		t.Errorf("声明权限的路由数 = %d, 期望 >= 30", total)
	}
}

// newTestContext 构造带真实请求的上下文，便于设置方法与方法/路径匹配结果
func newTestContext(engine *echo.Echo, method, path string) echo.Context {
	req := httptest.NewRequest(method, path, nil)
	ec := engine.NewContext(req, httptest.NewRecorder())
	ec.SetPath(path)
	return ec
}

// TestRoutePermIndex 请求期按「方法 + 路由模式」取权限标识
func TestRoutePermIndex(t *testing.T) {
	engine := buildEngine(t)
	tests := []struct {
		method string
		path   string
		want   string
	}{
		{"GET", "/api/posts", "post:list"},
		{"POST", "/api/posts", "post:create"},
		{"GET", "/api/posts/tree", "post:list"},
		{"GET", "/api/posts/:id", "post:read"},
		{"PUT", "/api/posts/:id", "post:update"},
		{"PUT", "/api/posts/:id/top", "post:update"},
		{"DELETE", "/api/posts/:id", "post:delete"},
		{"GET", "/api/cates/tree", "cate:tree"},
		{"GET", "/api/cates/:id", "cate:read"},
		{"PUT", "/api/roles/:id/permissions", "role:assign"},
		{"GET", "/api/status", "status:read"},
	}
	for _, it := range tests {
		ec := newTestContext(engine, it.method, it.path)
		if got := routePerm(ec); got != it.want {
			t.Errorf("%s %s 权限 = %q, want %q", it.method, it.path, got, it.want)
		}
	}
}

// TestRoutePermEmptyForPublic 公开接口与登录态接口不参与权限校验
func TestRoutePermEmptyForPublic(t *testing.T) {
	engine := buildEngine(t)
	for _, it := range []struct{ method, path string }{
		{"GET", "/api/captcha"},
		{"POST", "/api/tokens"},
		{"GET", "/api/site"},
		{"GET", "/api/auth"},
		{"PUT", "/api/auth"},
		{"GET", "/api/auth/permissions"},
		{"GET", "/api/dicts/:id/entries"},
		{"GET", "/api/configs/:id/val"},
	} {
		ec := newTestContext(engine, it.method, it.path)
		if got := routePerm(ec); got != "" {
			t.Errorf("%s %s 不应有权限标识, got %q", it.method, it.path, got)
		}
	}
}

// TestPermSyncedToDB 启动同步后权限点入库，且含分组授权权限
func TestPermSyncedToDB(t *testing.T) {
	buildEngine(t)
	ctx := context.Background()
	all, err := repo.PermAll(ctx)
	if err != nil {
		t.Fatalf("PermAll: %v", err)
	}
	if len(all) < 30 {
		t.Fatalf("权限点数量 = %d, 期望 >= 30", len(all))
	}
	perms := make(map[string]model.Perm, len(all))
	for _, it := range all {
		perms[it.Perm] = it
	}
	for _, want := range []string{
		"user:list", "user:grant",
		"role:assign", "role:grant",
		"post:list", "post:update", "content:grant",
		"cate:tree", "dict:list", "config:update", "status:read",
	} {
		if _, ok := perms[want]; !ok {
			t.Errorf("权限点 %q 未入库", want)
		}
	}
	// 分组键与中文名都要带上，权限树依赖它们聚合
	for cate, it := range perms {
		if it.GroupKey == "" {
			t.Errorf("%q 缺少分组键", cate)
		}
		if it.Group == "" || it.Name == "" {
			t.Errorf("%q 缺少中文名，got group=%q name=%q", cate, it.Group, it.Name)
		}
	}
	// 根角色不落授权表，但根用户仍应放行全部权限
	if !repo.RolePermAllow(ctx, model.RoleRootId, "post:create") {
		t.Error("根角色应放行任意权限")
	}
}

// TestRoutePermNamesUniquePerPath 同一路径的不同方法可以有不同权限
func TestRoutePermNamesUniquePerPath(t *testing.T) {
	engine := buildEngine(t)
	got := make(map[string]string, 8)
	for _, r := range engine.Routes() {
		if r.Path != "/api/posts/:id" {
			continue
		}
		if def, ok := parsePermName(r.Name); ok {
			got[r.Method] = def.Perm
		}
	}
	want := map[string]string{
		"GET":    "post:read",
		"PUT":    "post:update",
		"DELETE": "post:delete",
	}
	for method, perm := range want {
		if got[method] != perm {
			t.Errorf("%s /api/posts/:id 权限 = %q, want %q", method, got[method], perm)
		}
	}
}

// TestPublicPathsUnchanged 内部标识缩写为 perm，对外路径与权限标识保持不变
// 前端与权限声明都按这些字符串工作，改名不能波及它们。
func TestPublicPathsUnchanged(t *testing.T) {
	engine := buildEngine(t)
	paths := make(map[string]bool, 64)
	perms := make(map[string]bool, 64)
	for _, r := range engine.Routes() {
		paths[r.Method+" "+r.Path] = true
		if def, ok := parsePermName(r.Name); ok {
			perms[def.Perm] = true
		}
	}
	for _, want := range []string{
		"GET /api/permissions",
		"GET /api/permissions/tree",
		"GET /api/permissions/:id",
		"GET /api/auth/permissions",
		"GET /api/roles/:id/permissions",
		"PUT /api/roles/:id/permissions",
	} {
		if !paths[want] {
			t.Errorf("缺失路由 %q", want)
		}
	}
	// 权限标识里的资源名同样是对外契约
	for _, want := range []string{
		"permission:list", "permission:tree", "permission:read",
		"role:assign",
	} {
		if !perms[want] {
			t.Errorf("缺失权限标识 %q", want)
		}
	}
	// 分组授权权限由分组键派生，不在路由 Name 中
	defs := make([]model.PermInitArgs, 0, len(perms))
	for perm := range perms {
		defs = append(defs, model.PermInitArgs{Perm: perm, GroupKey: "permission", Group: "权限管理"})
	}
	found := false
	for _, def := range collectPerms(defs) {
		if def.Perm == "permission:grant" {
			found = true
		}
	}
	if !found {
		t.Error("未派生出 permission:grant")
	}
}
