package router

import (
	"blog/internal/model"
	"blog/internal/repo"
	"blog/pkg/reply"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
)

// auditEnv 构造一个挂了 Audit 的最小引擎
// 用假 handler 代替真实业务，专注验证「记了什么」而不是「做了什么」。
// 这里不走 engine.ServeHTTP，而是自己建 context 再手��串中间件：
// 只有这样才能在调用前注入 uid（真实环境由 AuthToken 中间件写入）。
type auditEnv struct {
	engine  *echo.Echo
	handler map[string]echo.HandlerFunc //method + " " + path -> handler
}

func newAuditEnv(t *testing.T) *auditEnv {
	t.Helper()
	setupRouterDB(t)
	return &auditEnv{
		engine:  echo.New(),
		handler: make(map[string]echo.HandlerFunc, 8),
	}
}

// handle 注册一个返回固定业务码的路由
func (env *auditEnv) handle(method, path, name string, code int) {
	env.handleWith(method, path, name, func(ctx echo.Context) error {
		return reply.Succ(ctx, "ok", code)
	})
}

// handleWith 注册自定义 handler
func (env *auditEnv) handleWith(method, path, name string, h echo.HandlerFunc) {
	env.handler[method+" "+path] = h
	// v4 的 Add 返回 *Route，名字可以注册后再改
	env.engine.Add(method, path, h, Audit).Name = name
}

// buildIndex 构建权限与名称索引（Audit 依赖）
func (env *auditEnv) buildIndex() {
	BuildRouteNames(env.engine)
}

// do 发一次请求，返回统一响应结构
func (env *auditEnv) do(t *testing.T, method, path, body, uid string, params map[string]string) reply.Reply {
	t.Helper()
	h, ok := env.handler[method+" "+path]
	if !ok {
		t.Fatalf("未注册的路由: %s %s", method, path)
	}
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	rec := httptest.NewRecorder()
	ec := env.engine.NewContext(req, rec)
	ec.SetPath(path)
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
	// 模拟 AuthToken 中间件写入的登录态
	if uid != "" {
		ec.Set("uid", uid)
		ec.Set("rid", model.RoleRootId)
	}
	if err := Audit(h)(ec); err != nil {
		t.Fatalf("请求失败: %v", err)
	}
	out := reply.Reply{}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("解析响应失败: %v, body=%s", err, rec.Body.String())
	}
	return out
}

// logsOf 取全部日志（先等异步队列排空）
func logsOf(t *testing.T) []model.Log {
	t.Helper()
	repo.LogFlush(2 * time.Second)
	mods, _, err := repo.LogList(context.Background(), &model.LogListArgs{})
	if err != nil {
		t.Fatalf("LogList: %v", err)
	}
	return mods
}

// TestAuditIsWriteMethod 写方法判定
func TestAuditIsWriteMethod(t *testing.T) {
	for _, it := range []string{http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete} {
		if !isWriteMethod(it) {
			t.Errorf("isWriteMethod(%q) = false, want true", it)
		}
	}
	for _, it := range []string{http.MethodGet, http.MethodHead, http.MethodOptions} {
		if isWriteMethod(it) {
			t.Errorf("isWriteMethod(%q) = true, want false", it)
		}
	}
}

// TestAuditEnabledByDefault 审计默认开启
func TestAuditEnabledByDefault(t *testing.T) {
	setupRouterDB(t)
	if !auditEnabled() {
		t.Error("默认应开启审计")
	}
	if !model.AuthDefault().Audit {
		t.Error("AuthDefault().Audit 应为 true")
	}
}

// TestAuditDisabledByConfig 关闭配置后中间件完全跳过
func TestAuditDisabledByConfig(t *testing.T) {
	env := newAuditEnv(t)
	// 关闭审计
	cfg := repo.GetCfg().Auth
	prev := cfg.Audit
	cfg.Audit = false
	t.Cleanup(func() { cfg.Audit = prev })

	env.handle(http.MethodDelete, "/api/users/:id", "user:delete|100|user|用户管理|用户删除", 200)
	env.buildIndex()
	env.do(t, http.MethodDelete, "/api/users/:id", `{"id":"u1"}`, "uid-1", map[string]string{"id": "u1"})

	if mods := logsOf(t); len(mods) != 0 {
		t.Errorf("关闭审计后仍记录了日志: %+v", mods)
	}
}

// TestAuditMask 脱敏递归处理嵌套结构
func TestAuditMask(t *testing.T) {
	in := map[string]any{
		"num":    "root",
		"passwd": "admin888",
		"nested": map[string]any{
			"token": "abc",
			"name":  "keep",
		},
		"list": []any{
			map[string]any{"secret": "s1", "title": "t1"},
		},
	}
	out := auditMask(in).(map[string]any)
	if out["num"] != "root" {
		t.Errorf("普通字段被改动: %v", out["num"])
	}
	if out["passwd"] != model.LogMaskedValue {
		t.Errorf("passwd 未脱敏: %v", out["passwd"])
	}
	nested := out["nested"].(map[string]any)
	if nested["token"] != model.LogMaskedValue || nested["name"] != "keep" {
		t.Errorf("嵌套脱敏异常: %v", nested)
	}
	list := out["list"].([]any)
	item := list[0].(map[string]any)
	if item["secret"] != model.LogMaskedValue || item["title"] != "t1" {
		t.Errorf("数组脱敏异常: %v", item)
	}
}

// TestAuditBody 请求体快照的边界处理
func TestAuditBody(t *testing.T) {
	if got := auditBody(nil); got != "" {
		t.Errorf("空体 = %q, want 空串", got)
	}
	if got := auditBody([]byte("  ")); got != "" {
		t.Errorf("空白体 = %q, want 空串", got)
	}
	if got := auditBody([]byte("not json")); !strings.Contains(got, "非 JSON") {
		t.Errorf("非 JSON 处理 = %q", got)
	}
	got := auditBody([]byte(`{"passwd":"x","num":"root"}`))
	if strings.Contains(got, `"x"`) {
		t.Errorf("密码泄漏: %s", got)
	}
	if !strings.Contains(got, model.LogMaskedValue) || !strings.Contains(got, "root") {
		t.Errorf("脱敏结果 = %s", got)
	}
}

// TestAuditRecordsWrite 写操作被记录，模块与动作取自路由 Name
func TestAuditRecordsWrite(t *testing.T) {
	env := newAuditEnv(t)
	env.handle(http.MethodDelete, "/api/users/:id", "user:delete|100|user|用户管理|用户删除", 200)
	env.buildIndex()

	env.do(t, http.MethodDelete, "/api/users/:id", `{"id":"u1"}`, "uid-1", map[string]string{"id": "u1"})

	mods := logsOf(t)
	if len(mods) != 1 {
		t.Fatalf("日志条数 = %d, want 1", len(mods))
	}
	mod := mods[0]
	if mod.Perm != "user:delete" {
		t.Errorf("Perm = %q", mod.Perm)
	}
	if mod.Module != "用户管理" || mod.Action != "用户删除" {
		t.Errorf("模块/动作 = %q/%q", mod.Module, mod.Action)
	}
	if mod.Method != http.MethodDelete || mod.Path != "/api/users/:id" {
		t.Errorf("方法/路径 = %q/%q", mod.Method, mod.Path)
	}
	if mod.UserId != "uid-1" {
		t.Errorf("UserId = %q", mod.UserId)
	}
	if mod.Status != reply.CodeSucc {
		t.Errorf("Status = %d, want %d", mod.Status, reply.CodeSucc)
	}
}

// TestAuditSkipsRead 只读接口不记录
func TestAuditSkipsRead(t *testing.T) {
	env := newAuditEnv(t)
	env.handle(http.MethodGet, "/api/users", "user:list|100|user|用户管理|用户列表", 200)
	env.buildIndex()

	env.do(t, http.MethodGet, "/api/users", "", "uid-1", nil)
	if mods := logsOf(t); len(mods) != 0 {
		t.Errorf("只读接口被记录: %+v", mods)
	}
}

// TestAuditSkipsNoPerm 没有权限标识的写接口不记录（登录除外）
func TestAuditSkipsNoPerm(t *testing.T) {
	env := newAuditEnv(t)
	env.handle(http.MethodPut, "/api/auth", "blog/internal/service/kernel.AuthEditMe", 200)
	env.buildIndex()

	env.do(t, http.MethodPut, "/api/auth", `{"name":"x"}`, "uid-1", nil)
	if mods := logsOf(t); len(mods) != 0 {
		t.Errorf("无权限标识的写接口被记录: %+v", mods)
	}
}

// TestAuditRecordsLogin 登录留痕，密码脱敏
func TestAuditRecordsLogin(t *testing.T) {
	env := newAuditEnv(t)
	env.handle(http.MethodPost, "/api/tokens", "", 200)
	env.buildIndex()

	env.do(t, http.MethodPost, "/api/tokens", `{"num":"root","passwd":"admin888"}`, "", nil)

	mods := logsOf(t)
	if len(mods) != 1 {
		t.Fatalf("登录日志条数 = %d, want 1", len(mods))
	}
	mod := mods[0]
	if mod.Module != "账号" || mod.Action != "登录" {
		t.Errorf("模块/动作 = %q/%q", mod.Module, mod.Action)
	}
	if mod.Perm != "" {
		t.Errorf("登录不应有 perm, got %q", mod.Perm)
	}
	if strings.Contains(mod.Data, "admin888") {
		t.Errorf("密码泄漏: %s", mod.Data)
	}
	if !strings.Contains(mod.Data, model.LogMaskedValue) {
		t.Errorf("密码未脱敏: %s", mod.Data)
	}
}

// TestAuditKeepsBody 审计读取请求体后，下游仍能拿到完整内容
func TestAuditKeepsBody(t *testing.T) {
	env := newAuditEnv(t)
	payload := `{"num":"root","passwd":"admin888","name":"测试"}`
	env.handleWith(http.MethodPost, "/api/tokens", "", func(ctx echo.Context) error {
		ipt := &struct {
			Num    string `json:"num"`
			Passwd string `json:"passwd"`
			Name   string `json:"name"`
		}{}
		if err := ctx.Bind(ipt); err != nil {
			t.Errorf("下游 Bind 失败: %v", err)
		}
		if ipt.Num != "root" || ipt.Passwd != "admin888" || ipt.Name != "测试" {
			t.Errorf("下游未拿到完整请求体: %+v", ipt)
		}
		return reply.Succ(ctx, "ok")
	})
	env.buildIndex()

	env.do(t, http.MethodPost, "/api/tokens", payload, "", nil)
}

// TestAuditCaptchaSkipped 验证码接口永不记录
func TestAuditCaptchaSkipped(t *testing.T) {
	env := newAuditEnv(t)
	env.handle(http.MethodPost, "/api/captcha", "", 200)
	env.buildIndex()

	env.do(t, http.MethodPost, "/api/captcha", `{}`, "", nil)
	if mods := logsOf(t); len(mods) != 0 {
		t.Errorf("验证码被记录: %+v", mods)
	}
}

// TestAuditRecordsFailure 失败的写操作同样留痕，记录失败业务码
func TestAuditRecordsFailure(t *testing.T) {
	env := newAuditEnv(t)
	env.handleWith(http.MethodPost, "/api/users", "user:create|100|user|用户管理|用户新增",
		func(ctx echo.Context) error {
			return reply.Fail(ctx, "账号已存在")
		})
	env.buildIndex()

	env.do(t, http.MethodPost, "/api/users", `{"num":"dup"}`, "uid-1", nil)

	mods := logsOf(t)
	if len(mods) != 1 {
		t.Fatalf("日志条数 = %d, want 1", len(mods))
	}
	if mods[0].Status != reply.CodeFail {
		t.Errorf("Status = %d, want %d", mods[0].Status, reply.CodeFail)
	}
}
