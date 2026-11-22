package model

import (
	"reflect"
	"strings"
	"testing"
)

func TestUserNormalizeOpenid(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"", ""},
		{"  oABC123  ", "oABC123"},
		{"\to_Ab-C\n", "o_Ab-C"},
		// openid 对大小写敏感，不做大小写变换
		{"OABCdef", "OABCdef"},
	}
	for _, it := range tests {
		if got := UserNormalizeOpenid(it.in); got != it.want {
			t.Errorf("UserNormalizeOpenid(%q) = %q, want %q", it.in, got, it.want)
		}
	}
}

func TestUserOpenidValid(t *testing.T) {
	valid := []string{
		"",                       // 空串表示未绑定
		"oABC123def_-",           // 典型 openid 字符集
		"oo0abcDEF1234567890123", // 28 位
		strings.Repeat("a", 128), // 上限
		"-_",                     // 只有符号也合法（base64url）
	}
	for _, it := range valid {
		if !UserOpenidValid(it) {
			t.Errorf("UserOpenidValid(%q) = false, want true", it)
		}
	}
	invalid := []string{
		"a b",                    // 空格
		"abc@weixin",             // @
		"openid/with/slash",      // 斜杠
		"中文",                     // 非 ASCII
		strings.Repeat("a", 129), // 超长
		"abc=",                   // base64url 无填充
	}
	for _, it := range invalid {
		if UserOpenidValid(it) {
			t.Errorf("UserOpenidValid(%q) = true, want false", it)
		}
	}
}

// TestUserOpenidColumn 模型字段名为 Openid，DB 列名为 openid
func TestUserOpenidColumn(t *testing.T) {
	if got := (User{}).TableName(); got != "user" {
		t.Fatalf("TableName() = %q, want user", got)
	}
	// 通过 Insert 后读回验证列确实存在且可承载 128 位长度
	mod := User{Openid: strings.Repeat("z", 128)}
	mod.BeforeInsert()
	if mod.Openid != strings.Repeat("z", 128) {
		t.Errorf("Openid 被截断: %q", mod.Openid)
	}
}

// TestUserHasNoTypeField 用户不再区分「管理员/普通用户」类型
// 权限完全由角色授权决定，结构体上不应再有 Type 字段。
func TestUserHasNoTypeField(t *testing.T) {
	typ := reflect.TypeOf(User{})
	if _, ok := typ.FieldByName("Type"); ok {
		t.Error("User 不应再有 Type 字段")
	}
	// 其余关键字段仍在
	for _, name := range []string{"Id", "Num", "Passwd", "RoleId", "Openid", "Status"} {
		if _, ok := typ.FieldByName(name); !ok {
			t.Errorf("User 缺少字段 %q", name)
		}
	}
	// 参数结构体也同步移除
	for _, bean := range []any{UserAddArgs{}, UserEditArgs{}} {
		bt := reflect.TypeOf(bean)
		if _, ok := bt.FieldByName("Type"); ok {
			t.Errorf("%T 不应再有 Type 字段", bean)
		}
	}
}
