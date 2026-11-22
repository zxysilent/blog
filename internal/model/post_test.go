package model

import (
	"strings"
	"testing"
)

func TestPostNormalizePath(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"  vscode-ext  ", "vscode-ext"},
		{"KeepCase", "KeepCase"}, // 不改大小写，避免历史链接失效
		{"a_b-c", "a_b-c"},
		{"\tone\n", "one"},
	}
	for _, it := range tests {
		if got := PostNormalizePath(it.in); got != it.want {
			t.Errorf("PostNormalizePath(%q) = %q, want %q", it.in, got, it.want)
		}
	}
}

func TestPostPathValid(t *testing.T) {
	valid := []string{
		"one", "vscode-ext", "a_b", "ABC123", "0",
		strings.Repeat("a", 128),
	}
	for _, it := range valid {
		if !PostPathValid(it) {
			t.Errorf("PostPathValid(%q) = false, want true", it)
		}
	}
	invalid := []string{
		"",                       // 空
		"-abc",                   // 符号开头
		"_abc",                   // 符号开头
		"a b",                    // 空格
		"a/b",                    // 斜杠
		"a.html",                 // 点
		"中文",                     // 非 ASCII
		strings.Repeat("a", 129), // 超长
	}
	for _, it := range invalid {
		if PostPathValid(it) {
			t.Errorf("PostPathValid(%q) = true, want false", it)
		}
	}
}

func TestPostKindAndStatusValid(t *testing.T) {
	for _, kind := range []int{KindPost, KindPage, KindNote} {
		if !PostKindValid(kind) {
			t.Errorf("PostKindValid(%d) = false", kind)
		}
	}
	for _, kind := range []int{0, 4, -1, 99} {
		if PostKindValid(kind) {
			t.Errorf("PostKindValid(%d) = true", kind)
		}
	}
	for _, status := range []int{StatusDraft, StatusFinish} {
		if !PostStatusValid(status) {
			t.Errorf("PostStatusValid(%d) = false", status)
		}
	}
	for _, status := range []int{0, 3, -1} {
		if PostStatusValid(status) {
			t.Errorf("PostStatusValid(%d) = true", status)
		}
	}
}

func TestPostHasContent(t *testing.T) {
	tests := []struct {
		markdown string
		richtext string
		want     bool
	}{
		{"", "", false},
		{"  \n\t ", "   ", false},
		{"# 标题", "", true},
		{"", "<p>x</p>", true},
		{"a", "b", true},
	}
	for _, it := range tests {
		if got := PostHasContent(it.markdown, it.richtext); got != it.want {
			t.Errorf("PostHasContent(%q,%q) = %v, want %v", it.markdown, it.richtext, got, it.want)
		}
	}
}

func TestPostValidate(t *testing.T) {
	title := "标题"
	path := "some-path"
	body := "正文"
	tests := []struct {
		name     string
		kind     int
		status   int
		title    string
		path     string
		markdown string
		richtext string
		wantErr  string
	}{
		{"正常文章", KindPost, StatusFinish, title, path, body, "", ""},
		{"正常页面", KindPage, StatusFinish, title, path, "", body, ""},
		{"草稿可空正文", KindPost, StatusDraft, title, path, "", "", ""},
		{"类型无效", 0, StatusDraft, title, path, body, "", "内容类型无效"},
		{"状态无效", KindPost, 0, title, path, body, "", "状态无效"},
		{"标题为空", KindPost, StatusDraft, "  ", path, body, "", "标题不能为空"},
		{"标题超长", KindPost, StatusDraft, strings.Repeat("字", 256), path, body, "", "标题不能超过 255 字符"},
		{"路径非法", KindPost, StatusDraft, title, "a b", body, "", "访问路径只能包含"},
		{"已发布缺正文", KindPost, StatusFinish, title, path, "  ", "", "已发布的内容必须填写正文"},
	}
	for _, it := range tests {
		err := PostValidate(it.kind, it.status, it.title, it.path, it.markdown, it.richtext)
		if it.wantErr == "" {
			if err != nil {
				t.Errorf("%s: 非预期错误 %v", it.name, err)
			}
			continue
		}
		if err == nil {
			t.Errorf("%s: 期望错误包含 %q", it.name, it.wantErr)
			continue
		}
		if !strings.Contains(err.Error(), it.wantErr) {
			t.Errorf("%s: 错误 = %q, want 包含 %q", it.name, err.Error(), it.wantErr)
		}
	}
}

func TestPostBeforeInsert(t *testing.T) {
	mod := &Post{Title: "x"}
	mod.BeforeInsert()
	if len(mod.Id) != 26 {
		t.Fatalf("id = %q, len = %d, want 26", mod.Id, len(mod.Id))
	}
	fixed := &Post{Id: "fixed"}
	fixed.BeforeInsert()
	if fixed.Id != "fixed" {
		t.Errorf("指定主键被覆盖: %q", fixed.Id)
	}
}

func TestPostPkHelpers(t *testing.T) {
	if !(&Post{}).IsPkZero() {
		t.Error("IsPkZero() = false for empty id")
	}
	if (&Post{Id: "x"}).IsPkZero() {
		t.Error("IsPkZero() = true for non-empty id")
	}
	if (Post{Id: "abc"}).GetId() != "abc" {
		t.Error("GetId() 不一致")
	}
}
