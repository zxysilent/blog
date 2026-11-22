package model

import (
	"testing"
	"time"
)

// TestSqlDefault 默认值必须合法
func TestSqlDefault(t *testing.T) {
	mod := SqlDefault()
	if err := mod.Valid(); err != nil {
		t.Fatalf("默认 SQL 配置不合法: %v", err)
	}
	if mod.Show {
		t.Error("默认不应打印全部 SQL")
	}
	if mod.SlowMs != SqlSlowMsDefault {
		t.Errorf("SlowMs = %d, want %d", mod.SlowMs, SqlSlowMsDefault)
	}
	if got := mod.SlowDuration(); got != time.Duration(SqlSlowMsDefault)*time.Millisecond {
		t.Errorf("SlowDuration = %v", got)
	}
}

// TestSqlValid 阈值边界与关闭开关
func TestSqlValid(t *testing.T) {
	cases := []struct {
		name string
		slow int64
		ok   bool
	}{
		{"最小阈值", SqlSlowMsMin, true},
		{"最大阈值", SqlSlowMsMax, true},
		{"阈值过小", SqlSlowMsMin - 1, false},
		{"阈值超上限", SqlSlowMsMax + 1, false},
		{"0 表示关闭", 0, true},
		{"负数也表示关闭", -1, true},
	}
	for _, c := range cases {
		mod := &Sql{SlowMs: c.slow}
		err := mod.Valid()
		if c.ok && err != nil {
			t.Errorf("%s: 应合法，实际 %v", c.name, err)
		}
		if !c.ok && err == nil {
			t.Errorf("%s: 应拒绝，实际通过", c.name)
		}
	}
}

// TestSqlSlowDuration 关闭时返回 0，调用方据此跳过慢日志
func TestSqlSlowDuration(t *testing.T) {
	if got := (&Sql{SlowMs: 0}).SlowDuration(); got != 0 {
		t.Errorf("关闭时 SlowDuration = %v, want 0", got)
	}
	if got := (&Sql{SlowMs: -5}).SlowDuration(); got != 0 {
		t.Errorf("负数时 SlowDuration = %v, want 0", got)
	}
	if got := (&Sql{SlowMs: 150}).SlowDuration(); got != 150*time.Millisecond {
		t.Errorf("SlowDuration = %v, want 150ms", got)
	}
}

// TestConfigValidSql 配置层面校验
func TestConfigValidSql(t *testing.T) {
	cases := []struct {
		name  string
		value string
		ok    bool
	}{
		{"合法", `{"show":true,"slow_ms":500}`, true},
		{"关闭慢日志", `{"show":false,"slow_ms":0}`, true},
		{"阈值过小", `{"show":false,"slow_ms":5}`, false},
		{"阈值过大", `{"show":false,"slow_ms":99999}`, false},
		{"非 JSON", `nope`, false},
	}
	for _, c := range cases {
		mod := &Config{Id: ConfigIdSql, Type: CfgTypeJson, Value: c.value}
		err := mod.Valid()
		if c.ok && err != nil {
			t.Errorf("%s: 应合法，实际 %v", c.name, err)
		}
		if !c.ok && err == nil {
			t.Errorf("%s: 应拒绝，实际通过", c.name)
		}
	}
}

// TestConfigToSqlFallback 脏值只重置出问题的那一项
func TestConfigToSqlFallback(t *testing.T) {
	cases := []struct {
		name     string
		value    string
		wantShow bool
		wantSlow int64
	}{
		{"空值", "", false, SqlSlowMsDefault},
		{"空对象", "{}", false, SqlSlowMsDefault},
		{"阈值越界", `{"show":true,"slow_ms":1}`, true, SqlSlowMsDefault},
		{"关闭慢日志", `{"show":false,"slow_ms":0}`, false, 0},
		{"合法值", `{"show":true,"slow_ms":800}`, true, 800},
	}
	for _, c := range cases {
		mod := (&Config{Id: ConfigIdSql, Type: CfgTypeJson, Value: c.value}).ToSql()
		if err := mod.Valid(); err != nil {
			t.Errorf("%s: 回落后的配置仍不合法: %v", c.name, err)
		}
		if mod.Show != c.wantShow || mod.SlowMs != c.wantSlow {
			t.Errorf("%s: 得到 show=%v slow=%d, want show=%v slow=%d",
				c.name, mod.Show, mod.SlowMs, c.wantShow, c.wantSlow)
		}
	}
}
