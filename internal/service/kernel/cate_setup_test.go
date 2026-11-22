package kernel

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"blog/conf"
	"blog/internal/repo"
	"blog/pkg/reply"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/logs"
)

// init 在任何测试（以及它们启动的 goroutine）之前导走日志输出
// 放到 setup 里会与 repo 日志协程并发读写 logs 的全局配置，触发数据竞争
func init() {
	logs.SetOutput(io.Discard)
}

// setupKernelDB 初始化一套隔离的内存库环境
//   - 每个测试独立命名的内存库，互不干扰
//   - 走真实的 conf.Init + repo.Init 流程，表结构与种子数据都会就绪
//   - 结束后关闭引擎，内存库销毁，不留落盘文件
func setupKernelDB(t *testing.T) {
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

// call 以 HTTP handler 的方式调用接口，返回统一响应结构
// 这是对 echo.Context 的真实绑定，因此参数解析、校验、响应格式都会被覆盖到
func call(t *testing.T, handler echo.HandlerFunc, method, target string, body string, params map[string]string) reply.Reply {
	t.Helper()
	return callCtx(t, handler, method, target, body, params, nil)
}

// callAs 以指定登录用户身份调用，手动注入 AuthToken 中间件原本会写入的上下文值
func callAs(t *testing.T, uid string, handler echo.HandlerFunc, method, target string, body string, params map[string]string) reply.Reply {
	t.Helper()
	return callCtx(t, handler, method, target, body, params, func(ec echo.Context) {
		ec.Set("uid", uid)
		if mod, err := repo.UserGet(context.Background(), uid); err == nil {
			ec.Set("rid", mod.RoleId)
		}
	})
}

// callCtx 构造真实 echo.Context 并执行 handler
func callCtx(t *testing.T, handler echo.HandlerFunc, method, target string, body string, params map[string]string, prepare func(echo.Context)) reply.Reply {
	t.Helper()
	var reader io.Reader
	if body != "" {
		reader = strings.NewReader(body)
	}
	req := httptest.NewRequest(method, target, reader)
	if body != "" {
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	}
	rec := httptest.NewRecorder()
	engine := echo.New()
	ec := engine.NewContext(req, rec)
	if len(params) > 0 {
		names := make([]string, 0, len(params))
		values := make([]string, 0, len(params))
		for k, v := range params {
			names = append(names, k)
			values = append(values, v)
		}
		ec.SetParamNames(names...)
		ec.SetParamValues(values...)
	}
	if prepare != nil {
		prepare(ec)
	}
	if err := handler(ec); err != nil {
		t.Fatalf("handler 返回错误: %v", err)
	}
	out := reply.Reply{}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("解析响应失败: %v, body=%s", err, rec.Body.String())
	}
	return out
}

// callOk 断言成功并返回原始 data
func callOk(t *testing.T, handler echo.HandlerFunc, method, target, body string, params map[string]string) any {
	t.Helper()
	out := call(t, handler, method, target, body, params)
	if out.Code != reply.CodeSucc {
		t.Fatalf("%s %s 期望成功, got code=%d msg=%s", method, target, out.Code, out.Msg)
	}
	return out.Data
}

// callFail 断言失败并返回错误信息
func callFail(t *testing.T, handler echo.HandlerFunc, method, target, body string, params map[string]string) string {
	t.Helper()
	out := call(t, handler, method, target, body, params)
	if out.Code == reply.CodeSucc {
		t.Fatalf("%s %s 期望失败, got code=200 data=%v", method, target, out.Data)
	}
	return out.Msg
}

// callList 断言成功并返回分页结果
func callList(t *testing.T, handler echo.HandlerFunc, method, target string) ([]any, int) {
	t.Helper()
	return callListCtx(t, handler, method, target, nil)
}

// callListAs 以指定登录用户身份调用分页接口
func callListAs(t *testing.T, uid string, handler echo.HandlerFunc, method, target string) ([]any, int) {
	t.Helper()
	return callListCtx(t, handler, method, target, func(ec echo.Context) {
		ec.Set("uid", uid)
		if mod, err := repo.UserGet(context.Background(), uid); err == nil {
			ec.Set("rid", mod.RoleId)
		}
	})
}

// callListCtx 构造真实上下文并解析分页响应
func callListCtx(t *testing.T, handler echo.HandlerFunc, method, target string, prepare func(echo.Context)) ([]any, int) {
	t.Helper()
	req := httptest.NewRequest(method, target, nil)
	rec := httptest.NewRecorder()
	engine := echo.New()
	ec := engine.NewContext(req, rec)
	if prepare != nil {
		prepare(ec)
	}
	if err := handler(ec); err != nil {
		t.Fatalf("handler 返回错误: %v", err)
	}
	out := struct {
		reply.Reply
		Data []any `json:"data"`
	}{}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("解析响应失败: %v, body=%s", err, rec.Body.String())
	}
	if out.Code != reply.CodeSucc {
		t.Fatalf("%s %s 期望成功, got code=%d msg=%s", method, target, out.Code, out.Msg)
	}
	return out.Data, out.Count
}

// mustAddCate 走接口新增栏目并返回主键
func mustAddCate(t *testing.T, body string) string {
	t.Helper()
	data := callOk(t, CateAdd, http.MethodPost, "/api/cates", body, nil)
	id, ok := data.(string)
	if !ok || id == "" {
		t.Fatalf("新增栏目未返回主键: %v", data)
	}
	return id
}

// cateOf 通过详情接口取栏目
func cateOf(t *testing.T, id string) map[string]any {
	t.Helper()
	data := callOk(t, CateGet, http.MethodGet, "/api/cates/"+id, "", map[string]string{"id": id})
	mod, ok := data.(map[string]any)
	if !ok {
		t.Fatalf("详情返回结构异常: %v", data)
	}
	return mod
}

// levelOf 取栏目层级
func levelOf(t *testing.T, mod map[string]any) int {
	t.Helper()
	v, ok := mod["level"].(float64)
	if !ok {
		t.Fatalf("level 字段缺失: %v", mod)
	}
	return int(v)
}

// cateNodes 把树响应转换为统一结构
func cateNodes(t *testing.T, data any) []cateNode {
	t.Helper()
	raw, err := json.Marshal(data)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}
	nodes := make([]cateNode, 0, 8)
	if err := json.Unmarshal(raw, &nodes); err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}
	return nodes
}

type cateNode struct {
	Id       string     `json:"id"`
	Name     string     `json:"name"`
	Path     string     `json:"path"`
	Level    int        `json:"level"`
	Pid      string     `json:"pid"`
	Use      bool       `json:"use"`
	Show     bool       `json:"show"`
	Inner    bool       `json:"inner"`
	Children []cateNode `json:"children"`
}

// atoi 便于断言 code / count 之类的数值
func atoi(v any) int {
	f, _ := v.(float64)
	return int(f)
}
