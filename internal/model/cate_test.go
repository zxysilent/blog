package model

import (
	"strings"
	"testing"
)

func TestCateNormalizePath(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"  Tech  ", "tech"},
		{"A-B_c", "a-b_c"},
		{"中文", "中文"},
		{"\tMulti Word ", "multi word"},
	}
	for _, it := range tests {
		if got := CateNormalizePath(it.in); got != it.want {
			t.Errorf("CateNormalizePath(%q) = %q, want %q", it.in, got, it.want)
		}
	}
}

func TestCatePathValid(t *testing.T) {
	valid := []string{
		"a", "0", "tech", "go-lang", "v2_api", "abc123",
		strings.Repeat("a", 64),
	}
	for _, it := range valid {
		if !CatePathValid(it) {
			t.Errorf("CatePathValid(%q) = false, want true", it)
		}
	}
	invalid := []string{
		"",                      // 空
		" Tech",                 // 大写+空格
		"-abc",                  // 符号开头
		"_abc",                  // 符号开头
		"a b",                   // 空格
		"a/b",                   // 斜杠
		"a.b",                   // 点
		"中文",                    // 非 ASCII
		strings.Repeat("a", 65), // 超长
	}
	for _, it := range invalid {
		if CatePathValid(it) {
			t.Errorf("CatePathValid(%q) = true, want false", it)
		}
	}
}

func TestCateBeforeInsert(t *testing.T) {
	// 未指定主键时自动生成 26 位 ULID
	mod := &Cate{Name: "技术"}
	mod.BeforeInsert()
	if len(mod.Id) != 26 {
		t.Fatalf("BeforeInsert id = %q, len = %d, want 26", mod.Id, len(mod.Id))
	}
	// 已指定主键时不覆盖
	mod2 := &Cate{Id: "fixed-id"}
	mod2.BeforeInsert()
	if mod2.Id != "fixed-id" {
		t.Errorf("BeforeInsert overwrote id: %q", mod2.Id)
	}
}

func TestCatePkHelpers(t *testing.T) {
	if !(&Cate{}).IsPkZero() {
		t.Error("IsPkZero() = false for empty id")
	}
	if (&Cate{Id: "x"}).IsPkZero() {
		t.Error("IsPkZero() = true for non-empty id")
	}
	mod := Cate{Id: "abc"}
	if mod.GetId() != "abc" {
		t.Errorf("GetId() = %q, want abc", mod.GetId())
	}
}

func TestCateMaxLevel(t *testing.T) {
	// 层级上限变更会影响服务层校验，锁住值便于发现意外改动
	if CateMaxLevel != 3 {
		t.Errorf("CateMaxLevel = %d, want 3", CateMaxLevel)
	}
}
