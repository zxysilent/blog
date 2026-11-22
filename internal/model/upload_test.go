package model

import (
	"strings"
	"testing"
)

// TestUploadExt 扩展名提取必须大小写无关且只取最后一段
func TestUploadExt(t *testing.T) {
	cases := map[string]string{
		"a.txt":     ".txt",
		"A.TXT":     ".txt",
		"a.tar.gz":  ".gz",
		"  a.PnG  ": ".png",
		"noext":     "",
		"a.":        ".",
		"dir/a.md":  ".md",
	}
	for in, want := range cases {
		if got := UploadExt(in); got != want {
			t.Fatalf("UploadExt(%q) = %q, want %q", in, got, want)
		}
	}
}

// TestUploadDangerRejected 所有显式标记的危险类型都不能通过任何校验
// 这些类型一旦落到 /static/ 下就是同源可执行内容（XSS / RCE）
func TestUploadDangerRejected(t *testing.T) {
	if len(UploadDangerExt) == 0 {
		t.Fatal("危险类型列表为空")
	}
	for ext := range UploadDangerExt {
		if UploadExtValid(ext) {
			t.Fatalf("%s 不应出现在通用白名单中", ext)
		}
		if UploadImageExtValid(ext) {
			t.Fatalf("%s 不应出现在图片白名单中", ext)
		}
		if _, err := UploadCheck("a" + ext); err == nil {
			t.Fatalf("%s 应被 UploadCheck 拒绝", ext)
		}
	}
}

// TestUploadImageExtSubset 图片白名单必须是通用白名单的子集
func TestUploadImageExtSubset(t *testing.T) {
	for ext := range UploadImageExt {
		if !UploadExtValid(ext) {
			t.Fatalf("图片白名单中的 %s 不在通用白名单内", ext)
		}
	}
}

// TestUploadCheck 通用上传校验
func TestUploadCheck(t *testing.T) {
	ok := []string{"a.txt", "a.PDF", "a.zip", "a.mp4", "a.md"}
	for _, name := range ok {
		ext, err := UploadCheck(name)
		if err != nil {
			t.Fatalf("%s 应该被允许: %v", name, err)
		}
		if ext != UploadExt(name) {
			t.Fatalf("%s 返回扩展名不一致: %s", name, ext)
		}
	}
	bad := []string{"a.html", "a.svg", "a.js", "a.css", "a.php", "a.jsp", "a.exe", "a.sh", "noext"}
	for _, name := range bad {
		if _, err := UploadCheck(name); err == nil {
			t.Fatalf("%s 不应被允许", name)
		}
	}
}

// TestUploadImageCheck 图片上传校验
func TestUploadImageCheck(t *testing.T) {
	for _, name := range []string{"a.jpg", "b.JPEG", "c.png", "d.gif", "e.webp", "f.bmp"} {
		if _, err := UploadImageCheck(name); err != nil {
			t.Fatalf("%s 应该被允许: %v", name, err)
		}
	}
	// 文档类合法但非图片，图片接口必须拒绝
	for _, name := range []string{"a.pdf", "a.txt", "a.zip", "a.gif.html"} {
		if _, err := UploadImageCheck(name); err == nil {
			t.Fatalf("%s 不应被图片接口接受", name)
		}
	}
}

// TestUploadExtError 错误提示要能区分「危险类型」与「不支持」
func TestUploadExtError(t *testing.T) {
	if msg := UploadExtError(".html").Error(); !strings.Contains(msg, "不允许上传") {
		t.Fatalf("危险类型提示不明确: %s", msg)
	}
	if msg := UploadExtError(".xyz").Error(); !strings.Contains(msg, "不支持") {
		t.Fatalf("未知类型提示不明确: %s", msg)
	}
	if msg := UploadExtError("").Error(); !strings.Contains(msg, "缺少扩展名") {
		t.Fatalf("空扩展名提示不明确: %s", msg)
	}
}

// TestUploadSizeText 体积文本
func TestUploadSizeText(t *testing.T) {
	cases := map[int64]string{
		0:        "0B",
		512:      "512B",
		1024:     "1.0KB",
		32 << 20: "32.0MB",
	}
	for in, want := range cases {
		if got := UploadSizeText(in); got != want {
			t.Fatalf("UploadSizeText(%d) = %q, want %q", in, got, want)
		}
	}
}

// TestUploadMaxSize 上限必须是正数，否则 MaxBytesReader 会拒绝一切请求
func TestUploadMaxSize(t *testing.T) {
	if UploadMaxSize <= 0 {
		t.Fatalf("UploadMaxSize 必须为正数: %d", UploadMaxSize)
	}
}

// TestUploadFileNameText 原始文件名入库前必须清洗
func TestUploadFileNameText(t *testing.T) {
	cases := map[string]string{
		"a.txt":       "a.txt",
		"  a.txt  ":   "a.txt",
		"a\nb\tc.txt": "abc.txt",
		"a\r\n.txt":   "a.txt",
		"":            "(未命名)",
		"   ":         "(未命名)",
	}
	for in, want := range cases {
		if got := UploadFileNameText(in); got != want {
			t.Fatalf("UploadFileNameText(%q) = %q, want %q", in, got, want)
		}
	}
	long := strings.Repeat("a", UploadNameTextLen+10)
	got := UploadFileNameText(long)
	if !strings.HasSuffix(got, "...") {
		t.Fatalf("超长文件名应截断: %s", got)
	}
	if len([]rune(strings.TrimSuffix(got, "..."))) != UploadNameTextLen {
		t.Fatalf("截断长度异常: %d", len([]rune(got)))
	}
}
