package model

import (
	"encoding/json"
	"os"
	"testing"
)

func TestJson(t *testing.T) {
	defer os.RemoveAll("logs")
	mod := Blog{}
	b, _ := json.Marshal(mod)
	t.Log(string(b))
}

// TestListNormalize 分页默认值统一由 Normalize 填充
func TestListNormalize(t *testing.T) {
	cases := []struct {
		name   string
		in     List
		ps     int
		wantPi int
		wantPs int
	}{
		{"全空取默认", List{}, 0, DefaultPage, DefaultPageSize},
		{"指定默认条数", List{}, 10, DefaultPage, 10},
		{"负数回落", List{Pi: -1, Ps: -5}, 0, DefaultPage, DefaultPageSize},
		{"合法值保留", List{Pi: 3, Ps: 20}, 0, 3, 20},
		{"只缺页码", List{Ps: 50}, 10, DefaultPage, 50},
	}
	for _, c := range cases {
		got := c.in
		got.Normalize(c.ps)
		if got.Pi != c.wantPi || got.Ps != c.wantPs {
			t.Errorf("%s: Normalize(%d) = (%d,%d), want (%d,%d)",
				c.name, c.ps, got.Pi, got.Ps, c.wantPi, c.wantPs)
		}
		if got.IsZero() {
			t.Errorf("%s: 归一化后 IsZero 应为 false", c.name)
		}
	}
}

// TestListLimitOffset 归一化前的取值也不应越界
func TestListLimitOffset(t *testing.T) {
	if got := (&List{}).GetLimit(); got != DefaultPageSize {
		t.Errorf("空 Ps 的 GetLimit = %d, want %d", got, DefaultPageSize)
	}
	if got := (&List{}).GetOffset(); got != 0 {
		t.Errorf("空参数的 GetOffset = %d, want 0", got)
	}
	if got := (&List{Pi: 3, Ps: 10}).GetOffset(); got != 20 {
		t.Errorf("第 3 页每页 10 条的 GetOffset = %d, want 20", got)
	}
	if got := (&List{}).GetCap(); got != DefaultPageSize {
		t.Errorf("空 Ps 的 GetCap = %d, want %d", got, DefaultPageSize)
	}
	if got := (&List{Ps: 1}).GetCap(); got != 8 {
		t.Errorf("小页的 GetCap = %d, want 8", got)
	}
}
