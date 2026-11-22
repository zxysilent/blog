package router

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"

	"blog/docs"

	"github.com/labstack/echo/v4"
)

// swaggerPaths 取出 swag 生成的文档中声明的全部「路径 -> 方法」组合
// docs.SwaggerInfo.ReadDoc() 返回的是原始模板，需先去掉模板占位符再解析
func swaggerPaths(t *testing.T) map[string]map[string]bool {
	t.Helper()
	raw := docs.SwaggerInfo.ReadDoc()
	// 模板占位符出现在 schemes/host/basePath 等处，先清空避免 JSON 解析失败
	raw = strings.NewReplacer("{{ marshal .Schemes }}", "[]", "{{.Host}}", "", "{{.BasePath}}", "",
		"{{escape .Description}}", "", "{{.Title}}", "", "{{.Version}}", "").Replace(raw)
	doc := struct {
		Paths map[string]map[string]json.RawMessage `json:"paths"`
	}{}
	if err := json.Unmarshal([]byte(raw), &doc); err != nil {
		t.Fatalf("解析 swagger 文档失败: %v", err)
	}
	out := make(map[string]map[string]bool, len(doc.Paths))
	for path, ops := range doc.Paths {
		methods := make(map[string]bool, len(ops))
		for method := range ops {
			methods[strings.ToUpper(method)] = true
		}
		out[path] = methods
	}
	return out
}

// toSwaggerPath 把 echo 路由模式转成 swagger 的路径写法
// /api/users/:id → /api/users/{id}
func toSwaggerPath(path string) string {
	parts := strings.Split(path, "/")
	for i := range parts {
		if strings.HasPrefix(parts[i], ":") {
			parts[i] = "{" + parts[i][1:] + "}"
		}
		if parts[i] == "*" {
			parts[i] = "{path}"
		}
	}
	return strings.Join(parts, "/")
}

// apiRoutes 取出引擎中真正的业务路由
// 排除两类登记项：
//   - 非标准方法的登记项（分组 404 兜底路由的 Method 是 RouteNotFound）
//   - 挂在公开分组上的 NotFoundHandler 兜底（/api 与 /api/*）
func apiRoutes(engine *echo.Echo) []*echo.Route {
	out := make([]*echo.Route, 0, 64)
	for _, r := range engine.Routes() {
		if !strings.HasPrefix(r.Path, "/api") {
			continue
		}
		if r.Path == "/api" || r.Path == "/api/*" {
			continue
		}
		switch r.Method {
		case http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete,
			http.MethodPatch, http.MethodHead, http.MethodOptions:
		default:
			continue
		}
		out = append(out, r)
	}
	return out
}

// TestSwaggerMatchesRoutes swagger 注释里的 @Router 必须与实际注册的路由完全一致
// 两者靠人工维护，很容易在改路径时只改一边，导致文档与接口对不上。
func TestSwaggerMatchesRoutes(t *testing.T) {
	engine := buildEngine(t)
	declared := swaggerPaths(t)

	actual := make(map[string]map[string]bool, 64)
	for _, r := range apiRoutes(engine) {
		path := toSwaggerPath(r.Path)
		if actual[path] == nil {
			actual[path] = make(map[string]bool, 2)
		}
		actual[path][r.Method] = true
	}

	// 文档中声明但代码里不存在的路由（最典型的改路径漏改注释）
	for path, methods := range declared {
		if _, ok := actual[path]; !ok {
			t.Errorf("swagger 声明了不存在的路径: %s", path)
			continue
		}
		for method := range methods {
			if !actual[path][method] {
				t.Errorf("swagger 声明了不存在的方法: %s %s", method, path)
			}
		}
	}
	// 代码里存在但文档没写的路由
	for path, methods := range actual {
		if _, ok := declared[path]; !ok {
			t.Errorf("路由缺少 swagger 注释: %s", path)
			continue
		}
		for method := range methods {
			if !declared[path][method] {
				t.Errorf("路由缺少 swagger 注释: %s %s", method, path)
			}
		}
	}
}

// TestSwaggerPathParams swagger 的 {id} 必须与路由的 :id 同名
func TestSwaggerPathParams(t *testing.T) {
	engine := buildEngine(t)
	declared := swaggerPaths(t)
	for _, r := range apiRoutes(engine) {
		path := toSwaggerPath(r.Path)
		if _, ok := declared[path]; !ok {
			// 缺注释的情况由上一个用例报告
			continue
		}
		// 路由里的参数名逐个出现在 swagger 路径中
		for _, seg := range strings.Split(r.Path, "/") {
			if !strings.HasPrefix(seg, ":") {
				continue
			}
			if !strings.Contains(path, "{"+seg[1:]+"}") {
				t.Errorf("%s 的参数名 %q 未体现在 swagger 路径 %q", r.Path, seg[1:], path)
			}
		}
	}
}
