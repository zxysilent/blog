package targz

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"
)

// writeFile 在 dir 下写一个文件，必要时创建父目录
func writeFile(t *testing.T, dir, name, body string) string {
	t.Helper()
	path := filepath.Join(dir, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll: %v", err)
	}
	if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	return path
}

// TestTargzUnTargzRoundTrip 目录压缩再解压，内容与层级保持一致
func TestTargzUnTargzRoundTrip(t *testing.T) {
	src := t.TempDir()
	writeFile(t, src, "a.txt", "hello")
	writeFile(t, src, "sub/b.txt", "world")
	writeFile(t, src, "sub/deep/c.txt", "!!")

	work := t.TempDir()
	archive := filepath.Join(work, "out.tar.gz")
	if err := Targz(archive, src); err != nil {
		t.Fatalf("Targz: %v", err)
	}
	if info, err := os.Stat(archive); err != nil || info.Size() == 0 {
		t.Fatalf("归档文件异常: %v", err)
	}

	dst := filepath.Join(work, "out")
	if err := UnTargz(dst, archive); err != nil {
		t.Fatalf("UnTargz: %v", err)
	}
	// 源目录名会作为归档内的顶层目录
	base := filepath.Base(src)
	for _, it := range []struct{ name, body string }{
		{"a.txt", "hello"},
		{"sub/b.txt", "world"},
		{"sub/deep/c.txt", "!!"},
	} {
		got, err := os.ReadFile(filepath.Join(dst, base, filepath.FromSlash(it.name)))
		if err != nil {
			t.Fatalf("读取 %s: %v", it.name, err)
		}
		if string(got) != it.body {
			t.Errorf("%s = %q, want %q", it.name, got, it.body)
		}
	}
	// 解压不应破坏归档本身
	if info, err := os.Stat(archive); err != nil || info.Size() == 0 {
		t.Errorf("解压后归档文件被损坏: %v", err)
	}
}

// TestTargzSuffix 目标文件缺少后缀时自动补全
func TestTargzSuffix(t *testing.T) {
	src := writeFile(t, t.TempDir(), "one.txt", "x")
	work := t.TempDir()
	dst := filepath.Join(work, "pkg")
	if err := Targz(dst, src); err != nil {
		t.Fatalf("Targz: %v", err)
	}
	if _, err := os.Stat(dst + ".tar.gz"); err != nil {
		t.Errorf("未生成 .tar.gz: %v", err)
	}
}

// TestTargzSourceMissing 源路径不存在时报错且不留半成品
func TestTargzSourceMissing(t *testing.T) {
	work := t.TempDir()
	archive := filepath.Join(work, "out.tar.gz")
	if err := Targz(archive, filepath.Join(work, "not-exist")); err == nil {
		t.Fatal("期望报错")
	}
	if _, err := os.Stat(archive); !os.IsNotExist(err) {
		t.Error("失败时应清理归档文件")
	}
}

// TestUnTargzRejectsEscape 含 ../ 的归档必须被拒绝
func TestUnTargzRejectsEscape(t *testing.T) {
	work := t.TempDir()
	archive := filepath.Join(work, "evil.tar.gz")

	f, err := os.Create(archive)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	gw := gzip.NewWriter(f)
	tw := tar.NewWriter(gw)
	body := []byte("pwned")
	if err := tw.WriteHeader(&tar.Header{Name: "../escape.txt", Mode: 0o644, Size: int64(len(body))}); err != nil {
		t.Fatalf("WriteHeader: %v", err)
	}
	if _, err := tw.Write(body); err != nil {
		t.Fatalf("Write: %v", err)
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("tar close: %v", err)
	}
	if err := gw.Close(); err != nil {
		t.Fatalf("gzip close: %v", err)
	}
	if err := f.Close(); err != nil {
		t.Fatalf("file close: %v", err)
	}

	dst := filepath.Join(work, "out")
	if err := UnTargz(dst, archive); err == nil {
		t.Fatal("越界路径应被拒绝")
	}
	if _, err := os.Stat(filepath.Join(work, "escape.txt")); !os.IsNotExist(err) {
		t.Error("越界文件不应被写入")
	}
}

// TestUnTargzEmptyDst 目标目录为空时直接报错
func TestUnTargzEmptyDst(t *testing.T) {
	if err := UnTargz("", "any.tar.gz"); err == nil {
		t.Fatal("空目标目录应被拒绝")
	}
}

// TestUnTargzBadGzip 非 gzip 内容返回错误
func TestUnTargzBadGzip(t *testing.T) {
	work := t.TempDir()
	bad := filepath.Join(work, "bad.tar.gz")
	if err := os.WriteFile(bad, []byte("not gzip"), 0o644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	if err := UnTargz(filepath.Join(work, "out"), bad); err == nil {
		t.Fatal("非 gzip 内容应报错")
	}
}

// TestTargzKeepsNestedPath 嵌套目录的相对路径不应带上临时目录前缀
func TestTargzKeepsNestedPath(t *testing.T) {
	src := t.TempDir()
	writeFile(t, src, "sub/keep.txt", "keep")
	work := t.TempDir()
	archive := filepath.Join(work, "out.tar.gz")
	if err := Targz(archive, src); err != nil {
		t.Fatalf("Targz: %v", err)
	}
	raw, err := os.ReadFile(archive)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	gr, err := gzip.NewReader(bytes.NewReader(raw))
	if err != nil {
		t.Fatalf("gzip: %v", err)
	}
	tr := tar.NewReader(gr)
	found := false
	for {
		hdr, err := tr.Next()
		if err != nil {
			break
		}
		if hdr.Name == "sub/keep.txt" || hdr.Name == filepath.Base(src)+"/sub/keep.txt" {
			found = true
		}
		if filepath.IsAbs(hdr.Name) {
			t.Errorf("归档内出现绝对路径: %s", hdr.Name)
		}
	}
	if !found {
		t.Error("未找到 sub/keep.txt 的相对路径条目")
	}
}
