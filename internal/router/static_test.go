package router

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/labstack/echo/v4"
)

// TestStaticSafeNosniff 静态资源必须带 nosniff
// 上传目录与站点同源，缺少该头时浏览器可能对内容做 MIME 嗅探，
// 从而把「内容为 HTML 的 .txt」当页面执行，扩展名白名单就白做了
func TestStaticSafeNosniff(t *testing.T) {
	engine := echo.New()
	engine.Group("", staticSafe).Static("/static", t.TempDir())

	req := httptest.NewRequest(http.MethodGet, "/static/a.txt", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if got := rec.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("X-Content-Type-Options = %q, want nosniff", got)
	}
}

// TestStaticNoTraversal 静态目录不能越权读取上级文件
func TestStaticNoTraversal(t *testing.T) {
	root := t.TempDir()
	// 在 root 的上一级放一个「机密」文件，尝试用 ../ 访问
	secret := filepath.Join(filepath.Dir(root), "secret.txt")
	if err := os.WriteFile(secret, []byte("top-secret"), 0o644); err != nil {
		t.Fatalf("准备密文失败: %v", err)
	}
	t.Cleanup(func() { os.Remove(secret) })

	engine := echo.New()
	engine.Group("", staticSafe).Static("/static", root)

	req := httptest.NewRequest(http.MethodGet, "/static/../secret.txt", nil)
	rec := httptest.NewRecorder()
	engine.ServeHTTP(rec, req)
	if rec.Code == http.StatusOK {
		t.Fatalf("越权读取成功, body=%s", rec.Body.String())
	}
}
