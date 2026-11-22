package kernel

import (
	"bytes"
	"encoding/json"
	"image"
	"image/png"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"blog/internal/model"
	"blog/pkg/reply"

	"github.com/labstack/echo/v4"
)

// uploadCall 构造真实的 multipart 请求并执行 upload handler
// 返回业务响应与上下文，便于校验审计摘要等副作用
func uploadCall(t *testing.T, handler echo.HandlerFunc, filename string, content []byte) (reply.Reply, echo.Context) {
	t.Helper()
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
	req := httptest.NewRequest(http.MethodPost, "/api/files", body)
	req.Header.Set(echo.HeaderContentType, writer.FormDataContentType())
	rec := httptest.NewRecorder()
	engine := echo.New()
	ec := engine.NewContext(req, rec)
	if err := handler(ec); err != nil {
		t.Fatalf("handler 返回错误: %v", err)
	}
	out := reply.Reply{}
	if err := json.Unmarshal(rec.Body.Bytes(), &out); err != nil {
		t.Fatalf("解析响应失败: %v, body=%s", err, rec.Body.String())
	}
	return out, ec
}

// uploadEnv 隔离的文件系统环境：切到临时目录，避免污染仓库 static/
func uploadEnv(t *testing.T) {
	t.Helper()
	setupKernelDB(t)
	t.Chdir(t.TempDir())
}

// pngBytes 生成一张合法的 1x1 PNG
func pngBytes(t *testing.T) []byte {
	t.Helper()
	buf := &bytes.Buffer{}
	if err := png.Encode(buf, image.NewRGBA(image.Rect(0, 0, 1, 1))); err != nil {
		t.Fatalf("生成 PNG 失败: %v", err)
	}
	return buf.Bytes()
}

func TestUploadFileOK(t *testing.T) {
	uploadEnv(t)
	out, ec := uploadCall(t, UploadFile, "说明.txt", []byte("hello 上传"))
	if out.Code != reply.CodeSucc {
		t.Fatalf("期望上传成功, got code=%d msg=%s", out.Code, out.Msg)
	}
	url, ok := out.Data.(string)
	if !ok || url == "" {
		t.Fatalf("返回值异常: %v", out.Data)
	}
	// 返回地址是「endpoint/相对路径」，落盘位置应存在
	rel := strings.TrimPrefix(url[strings.Index(url, "static/"):], "")
	if _, err := os.Stat(filepath.FromSlash(rel)); err != nil {
		t.Fatalf("文件未落盘: %v", err)
	}
	// 文件名必须是服务端生成的 ULID，不能回显原始文件名
	base := filepath.Base(rel)
	if base == "说明.txt" || strings.Contains(base, "说明") {
		t.Fatalf("文件名不应包含客户端提供的名称: %s", base)
	}
	if !strings.HasSuffix(base, ".txt") {
		t.Fatalf("扩展名丢失: %s", base)
	}
	if len(strings.TrimSuffix(base, ".txt")) != model.UploadNameLen {
		t.Fatalf("文件名长度应为 ULID 的 %d 位: %s", model.UploadNameLen, base)
	}
	// 审计摘要应带上原始文件名与大小
	data, _ := ec.Get(model.CtxAuditData).(string)
	if !strings.Contains(data, "说明.txt") || !strings.Contains(data, "存储路径=") {
		t.Fatalf("审计摘要异常: %q", data)
	}
}

func TestUploadFileRejectDanger(t *testing.T) {
	uploadEnv(t)
	for _, name := range []string{"x.html", "x.svg", "x.js", "x.php", "x.sh", "X.HTML"} {
		out, _ := uploadCall(t, UploadFile, name, []byte("<script>alert(1)</script>"))
		if out.Code == reply.CodeSucc {
			t.Fatalf("%s 危险类型不应上传成功", name)
		}
		if !strings.Contains(out.Msg, "不允许上传") {
			t.Fatalf("%s 提示不明确: %s", name, out.Msg)
		}
	}
}

func TestUploadFileRejectUnknown(t *testing.T) {
	uploadEnv(t)
	for _, name := range []string{"a.xyz", "noext", "a.exe"} {
		out, _ := uploadCall(t, UploadFile, name, []byte("x"))
		if out.Code == reply.CodeSucc {
			t.Fatalf("%s 不在白名单内不应成功", name)
		}
	}
}

func TestUploadImageOK(t *testing.T) {
	uploadEnv(t)
	out, _ := uploadCall(t, UploadImage, "a.png", pngBytes(t))
	if out.Code != reply.CodeSucc {
		t.Fatalf("期望图片上传成功, got code=%d msg=%s", out.Code, out.Msg)
	}
}

func TestUploadImageRejectSpoof(t *testing.T) {
	uploadEnv(t)
	// 扩展名合法但内容不是图片：客户端声明的 Content-Type 不可信，必须解码校验
	out, _ := uploadCall(t, UploadImage, "a.png", []byte("<?php echo 1; ?>"))
	if out.Code == reply.CodeSucc {
		t.Fatal("伪装成 png 的非图片内容不应上传成功")
	}
	if !strings.Contains(out.Msg, "图片") {
		t.Fatalf("提示不明确: %s", out.Msg)
	}
}

func TestUploadImageRejectNonImageExt(t *testing.T) {
	uploadEnv(t)
	out, _ := uploadCall(t, UploadImage, "a.pdf", []byte("%PDF-1.4"))
	if out.Code == reply.CodeSucc {
		t.Fatal("图片接口不应接受非图片扩展名")
	}
}

func TestUploadFileRejectTooLarge(t *testing.T) {
	uploadEnv(t)
	big := bytes.Repeat([]byte("a"), int(model.UploadMaxSize)+1024)
	out, _ := uploadCall(t, UploadFile, "big.txt", big)
	if out.Code == reply.CodeSucc {
		t.Fatal("超过上限的文件不应上传成功")
	}
	if !strings.Contains(out.Msg, "文件过大") {
		t.Fatalf("提示不明确: %s", out.Msg)
	}
}
