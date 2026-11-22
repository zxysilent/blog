package model

import (
	"strings"
	"testing"
)

// TestAuthDefaultValid 默认值必须落在合法区间内
func TestAuthDefaultValid(t *testing.T) {
	mod := AuthDefault()
	if err := mod.Valid(); err != nil {
		t.Fatalf("默认认证配置不合法: %v", err)
	}
	if mod.TokenExp != AuthTokenExpDefault {
		t.Errorf("TokenExp = %d, want %d", mod.TokenExp, AuthTokenExpDefault)
	}
}

// TestAuthValidRange 过期时间与私钥长度的边界
func TestAuthValidRange(t *testing.T) {
	secret := strings.Repeat("a", AuthSecretMin)
	cases := []struct {
		name string
		mod  Auth
		ok   bool
	}{
		{"最小过期时间", Auth{TokenExp: AuthTokenExpMin, TokenSecret: secret}, true},
		{"最大过期时间", Auth{TokenExp: AuthTokenExpMax, TokenSecret: secret}, true},
		{"过期时间为 0", Auth{TokenExp: 0, TokenSecret: secret}, false},
		{"过期时间超上限", Auth{TokenExp: AuthTokenExpMax + 1, TokenSecret: secret}, false},
		{"私钥过短", Auth{TokenExp: 24, TokenSecret: strings.Repeat("a", AuthSecretMin-1)}, false},
		{"私钥过长", Auth{TokenExp: 24, TokenSecret: strings.Repeat("a", AuthSecretMax+1)}, false},
		{"私钥含空格", Auth{TokenExp: 24, TokenSecret: "aaaa bbbb cccc dddd"}, false},
	}
	for _, c := range cases {
		mod := c.mod
		mod.Normalize()
		err := mod.Valid()
		if c.ok && err != nil {
			t.Errorf("%s: 应合法，实际 %v", c.name, err)
		}
		if !c.ok && err == nil {
			t.Errorf("%s: 应拒绝，实际通过", c.name)
		}
	}
}

// TestConfigValidAuth 配置层面校验认证配置
// 直接改库或通过接口提交非法值都必须被拦下
func TestConfigValidAuth(t *testing.T) {
	good := `{"token_exp":24,"token_secret":"` + strings.Repeat("x", AuthSecretMin) + `"}`
	cases := []struct {
		name  string
		value string
		ok    bool
	}{
		{"合法", good, true},
		{"过期时间为负", `{"token_exp":-1,"token_secret":"` + strings.Repeat("x", AuthSecretMin) + `"}`, false},
		{"过期时间越界", `{"token_exp":9999,"token_secret":"` + strings.Repeat("x", AuthSecretMin) + `"}`, false},
		{"私钥太短", `{"token_exp":24,"token_secret":"short"}`, false},
		{"非 JSON", `not-json`, false},
	}
	for _, c := range cases {
		mod := &Config{Id: ConfigIdAuth, Type: CfgTypeJson, Value: c.value}
		err := mod.Valid()
		if c.ok && err != nil {
			t.Errorf("%s: 应合法，实际 %v", c.name, err)
		}
		if !c.ok && err == nil {
			t.Errorf("%s: 应拒绝，实际通过", c.name)
		}
	}
}

// TestConfigToAuthFallback 库中配置残缺时回落默认值，避免签名不可用
// 只回落出问题的那一项，合法的另一项保持库中原值
func TestConfigToAuthFallback(t *testing.T) {
	okSecret := strings.Repeat("a", AuthSecretMin)
	cases := []struct {
		name       string
		value      string
		wantExp    int
		wantSecret string
	}{
		{"空值", "", AuthTokenExpDefault, AuthSecretDefault},
		{"空对象", "{}", AuthTokenExpDefault, AuthSecretDefault},
		{"私钥过短", `{"token_exp":24,"token_secret":"x"}`, 24, AuthSecretDefault},
		{"过期时间越界", `{"token_exp":0,"token_secret":"` + okSecret + `"}`, AuthTokenExpDefault, okSecret},
		{"非法 JSON", "not-json", AuthTokenExpDefault, AuthSecretDefault},
	}
	for _, c := range cases {
		mod := (&Config{Id: ConfigIdAuth, Type: CfgTypeJson, Value: c.value}).ToAuth()
		if err := mod.Valid(); err != nil {
			t.Errorf("%s: 回落后的配置仍不合法: %v", c.name, err)
		}
		if mod.TokenExp != c.wantExp {
			t.Errorf("%s: exp = %d, want %d", c.name, mod.TokenExp, c.wantExp)
		}
		if mod.TokenSecret != c.wantSecret {
			t.Errorf("%s: secret = %q, want %q", c.name, mod.TokenSecret, c.wantSecret)
		}
	}
}
