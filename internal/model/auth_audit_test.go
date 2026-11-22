package model

import (
	"encoding/json"
	"strings"
	"testing"
)

// TestAuthAuditDefault 审计开关默认开启
func TestAuthAuditDefault(t *testing.T) {
	if !AuthDefault().Audit {
		t.Error("AuthDefault().Audit 应为 true，审计是安全基线")
	}
}

// TestAuthAuditJSON 开关随配置 JSON 一起序列化/反序列化
// 后台修改配置走的是同一份 JSON，字段名必须是 audit
func TestAuthAuditJSON(t *testing.T) {
	raw, err := json.Marshal(AuthDefault())
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if !strings.Contains(string(raw), `"audit":true`) {
		t.Errorf("默认配置 JSON 缺少 audit=true: %s", raw)
	}

	// 关闭状态能正确往返
	secret := strings.Repeat("a", AuthSecretMin)
	in := `{"token_exp":24,"token_secret":"` + secret + `","audit":false}`
	mod := &Auth{}
	if err := json.Unmarshal([]byte(in), mod); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if mod.Audit {
		t.Error("audit=false 未生效")
	}
	if err := mod.Valid(); err != nil {
		t.Errorf("合法配置被拒: %v", err)
	}

	// 老配置没有 audit 字段时按 false 处理，保持向后兼容
	legacy := &Auth{}
	if err := json.Unmarshal([]byte(`{"token_exp":24,"token_secret":"`+secret+`"}`), legacy); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if legacy.Audit {
		t.Error("缺失 audit 字段时不应默认开启")
	}
}

// TestAuthAuditDoesNotAffectValid 开关不影响其余字段校验
func TestAuthAuditDoesNotAffectValid(t *testing.T) {
	mod := AuthDefault()
	mod.Audit = false
	if err := mod.Valid(); err != nil {
		t.Errorf("关闭审计后校验失败: %v", err)
	}
	// 非法过期时间仍应被拒
	bad := AuthDefault()
	bad.Audit = false
	bad.TokenExp = 0
	if err := bad.Valid(); err == nil {
		t.Error("非法过期时间应被拒绝")
	}
}

// TestConfigToAuthAudit 从 Config 反序列化时带上开关
func TestConfigToAuthAudit(t *testing.T) {
	secret := strings.Repeat("a", AuthSecretMin)
	for _, want := range []bool{true, false} {
		raw := `{"token_exp":12,"token_secret":"` + secret + `","audit":` + map[bool]string{true: "true", false: "false"}[want] + `}`
		auth := (&Config{Id: ConfigIdAuth, Type: CfgTypeJson, Value: raw}).ToAuth()
		if auth.Audit != want {
			t.Errorf("Audit = %v, want %v", auth.Audit, want)
		}
	}
}
