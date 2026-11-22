package model

import (
	"strings"
	"testing"
)

// TestLogNeedMask 敏感字段识别，大小写与空格都应该容错
func TestLogNeedMask(t *testing.T) {
	sensitive := []string{
		"passwd", "PASSWD", " Passwd ", "password",
		"opasswd", "npasswd", "old_password", "new_password",
		"secret", "token", "vcode", "vreal",
		"private_key", "access_key",
	}
	for _, it := range sensitive {
		if !LogNeedMask(it) {
			t.Errorf("LogNeedMask(%q) = false, want true", it)
		}
	}
	// 普通字段不脱敏
	normal := []string{
		"", "num", "name", "email", "title", "path", "content",
		"passwd2", "mytoken", "token_type",
	}
	for _, it := range normal {
		if LogNeedMask(it) {
			t.Errorf("LogNeedMask(%q) = true, want false", it)
		}
	}
}

// TestLogSensitiveFieldsCoverage 常见密码字段必须都在脱敏名单里
func TestLogSensitiveFieldsCoverage(t *testing.T) {
	// 与 model/user.go 的密码参数保持一致
	for _, it := range []string{"passwd", "opasswd", "npasswd"} {
		found := false
		for _, f := range LogSensitiveFields {
			if strings.EqualFold(f, it) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("脱敏名单缺少 %q", it)
		}
	}
}

func TestLogBeforeInsert(t *testing.T) {
	mod := &Log{Perm: "user:delete"}
	mod.BeforeInsert()
	if len(mod.Id) != 26 {
		t.Fatalf("id = %q, len = %d, want 26", mod.Id, len(mod.Id))
	}
	fixed := &Log{Id: "fixed"}
	fixed.BeforeInsert()
	if fixed.Id != "fixed" {
		t.Errorf("指定主键被覆盖: %q", fixed.Id)
	}
}

func TestLogPkHelpers(t *testing.T) {
	if !(&Log{}).IsPkZero() {
		t.Error("IsPkZero() = false for empty id")
	}
	if (&Log{Id: "x"}).IsPkZero() {
		t.Error("IsPkZero() = true for non-empty id")
	}
	if (Log{Id: "abc"}).GetId() != "abc" {
		t.Error("GetId() 不一致")
	}
	if (Log{}).TableName() != "log" {
		t.Errorf("TableName() = %q, want log", (Log{}).TableName())
	}
}
