package router

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"blog/internal/model"
	"blog/pkg/reply"

	"github.com/labstack/echo/v4"
)

// doMultipart 发送 multipart 请求
// 与 do 的区别只在于请求体构造，仍然复用同一套 env/handler
func (env *auditEnv) doMultipart(t *testing.T, path, filename string, content []byte, uid string) (reply.Reply, echo.Context) {
	t.Helper()
	h, ok := env.handler[http.MethodPost+" "+path]
	if !ok {
		t.Fatalf("未注册的路由: POST %s", path)
	}
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		t.Fatalf("构造 multipart 失败: %v", err)
	}
	if _, err := part.Write(content); err != nil {
		t.Fatalf("写入 multipart 失败: %v", err)
	}
	if err := writer.Close(); err != nil {
		t.Fatalf("关闭 multipart 失败: %v", err)
	}
	req := httptest.NewRequest(http.MethodPost, path, body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	ec := env.engine.NewContext(req, rec)
	ec.SetPath(path)
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
	return out, ec
}

// TestAuditMultipartUsesHandlerSummary 上传请求的审计内容取自处理器摘要
//   - 中间件不能预先读取 multipart 请求体，否则 http.MaxBytesReader 失效
//   - 入库内容应包含原始文件名、大小、存储路径，而不是「非 JSON 请求体」
func TestAuditMultipartUsesHandlerSummary(t *testing.T) {
	env := newAuditEnv(t)
	env.handleWith(http.MethodPost, "/api/files", "file:create|700|system|系统|文件上传", func(ctx echo.Context) error {
		// 处理器仍能读取上传内容，说明中间件没有消费请求体
		if err := ctx.Request().ParseMultipartForm(1 << 20); err != nil {
			return reply.Fail(ctx, "解析失败", err)
		}
		if _, _, err := ctx.Request().FormFile("file"); err != nil {
			return reply.Fail(ctx, "取文件失败", err)
		}
		ctx.Set(model.CtxAuditData, "原始文件名=a.txt; 扩展名=.txt; 大小=5B; 存储路径=static/upload/x/a.txt")
		return reply.Succ(ctx, "上传成功", "http://x/static/upload/x/a.txt")
	})
	env.buildIndex()

	out, _ := env.doMultipart(t, "/api/files", "a.txt", []byte("hello"), "")
	if out.Code != reply.CodeSucc {
		t.Fatalf("期望成功, got code=%d msg=%s", out.Code, out.Msg)
	}
	mods := logsOf(t)
	if len(mods) != 1 {
		t.Fatalf("期望 1 条日志, got %d", len(mods))
	}
	mod := mods[0]
	if mod.Perm != "file:create" {
		t.Fatalf("权限标识异常: %s", mod.Perm)
	}
	if mod.Module != "系统" || mod.Action != "文件上传" {
		t.Fatalf("分组/动作异常: %s/%s", mod.Module, mod.Action)
	}
	for _, want := range []string{"原始文件名=a.txt", "大小=5B", "存储路径=static/upload/"} {
		if !strings.Contains(mod.Data, want) {
			t.Fatalf("审计内容缺少 %q: %s", want, mod.Data)
		}
	}
	if strings.Contains(mod.Data, "非 JSON") {
		t.Fatalf("上传请求不应记录成非 JSON 请求体: %s", mod.Data)
	}
}

// TestAuditMultipartNoSummary 处理器没给摘要时不应记录成二进制长度
func TestAuditMultipartNoSummary(t *testing.T) {
	env := newAuditEnv(t)
	env.handleWith(http.MethodPost, "/api/files", "file:create|700|system|系统|文件上传", func(ctx echo.Context) error {
		return reply.Succ(ctx, "上传成功")
	})
	env.buildIndex()

	if out, _ := env.doMultipart(t, "/api/files", "a.txt", []byte("hello"), ""); out.Code != reply.CodeSucc {
		t.Fatalf("期望成功, got code=%d msg=%s", out.Code, out.Msg)
	}
	mods := logsOf(t)
	if len(mods) != 1 {
		t.Fatalf("期望 1 条日志, got %d", len(mods))
	}
	if mods[0].Data != "" {
		t.Fatalf("无摘要时应为空, got %q", mods[0].Data)
	}
}
