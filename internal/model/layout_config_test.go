package model

import (
	"encoding/json"
	"testing"
)

// TestLayoutDefaultValid 默认值必须自身合法，否则新库种子数据就会写不进去
func TestLayoutDefaultValid(t *testing.T) {
	mod := LayoutDefault()
	if err := mod.Valid(); err != nil {
		t.Fatalf("默认布局配置不合法: %v", err)
	}
	buf, err := json.Marshal(mod)
	if err != nil {
		t.Fatalf("序列化失败: %v", err)
	}
	back := &Layout{}
	if err := json.Unmarshal(buf, back); err != nil {
		t.Fatalf("反序列化失败: %v", err)
	}
	if *back != *mod {
		t.Fatalf("默认值往返不一致:\n got %+v\nwant %+v", back, mod)
	}
}

// TestLayoutValidReject 枚举越界必须被拒绝，避免前端拿到无法渲染的值
func TestLayoutValidReject(t *testing.T) {
	cases := []struct {
		name string
		mod  Layout
	}{
		{"导航模式", Layout{NavMode: "left", NavTheme: NavThemeDark, MenuWidth: 200, MenuMinWidth: 64, PageAnimateType: PageAnimateZoomFade}},
		{"导航风格", Layout{NavMode: NavModeVertical, NavTheme: "blue", MenuWidth: 200, MenuMinWidth: 64, PageAnimateType: PageAnimateZoomFade}},
		{"动画类型", Layout{NavMode: NavModeVertical, NavTheme: NavThemeDark, MenuWidth: 200, MenuMinWidth: 64, PageAnimateType: "spin"}},
		{"菜单宽度", Layout{NavMode: NavModeVertical, NavTheme: NavThemeDark, MenuWidth: 10, MenuMinWidth: 64, PageAnimateType: PageAnimateZoomFade}},
		{"折叠宽度", Layout{NavMode: NavModeVertical, NavTheme: NavThemeDark, MenuWidth: 200, MenuMinWidth: 999, PageAnimateType: PageAnimateZoomFade}},
	}
	for _, c := range cases {
		if err := c.mod.Valid(); err == nil {
			t.Errorf("%s: 越界值应被拒绝", c.name)
		}
	}
}

// TestLayoutNormalize 读取历史数据时静默纠正，不因为一个脏值让整份配置不可用
func TestLayoutNormalize(t *testing.T) {
	mod := &Layout{
		NavMode:         "LEFT",
		NavTheme:        "",
		PageAnimateType: "spin",
		MenuWidth:       9999,
		MenuMinWidth:    1,
		MenuMix:         true,
	}
	mod.Normalize()
	if mod.NavMode != NavModeVertical {
		t.Errorf("导航模式未回落: %q", mod.NavMode)
	}
	if mod.NavTheme != NavThemeDark {
		t.Errorf("导航风格未回落: %q", mod.NavTheme)
	}
	if mod.PageAnimateType != PageAnimateZoomFade {
		t.Errorf("动画类型未回落: %q", mod.PageAnimateType)
	}
	if mod.MenuWidth != MenuWidthDefault || mod.MenuMinWidth != 64 {
		t.Errorf("菜单宽度未回落: %d/%d", mod.MenuWidth, mod.MenuMinWidth)
	}
	// 非混合模式下不允许出现「分割菜单」
	if mod.MenuMix {
		t.Error("非 horizontal-mix 模式下 MenuMix 应为 false")
	}
	if err := mod.Valid(); err != nil {
		t.Errorf("收敛后应合法: %v", err)
	}
}

// TestConfigPerm 权限标识必须与主键一一对应，供路由与业务校验共用
func TestConfigPerm(t *testing.T) {
	if got := ConfigPerm(ConfigIdBasic); got != "config:basic" {
		t.Errorf("ConfigPerm(basic) = %q", got)
	}
	seen := make(map[string]struct{}, len(BuiltinConfigIds))
	for _, id := range BuiltinConfigIds {
		perm := ConfigPerm(id)
		if _, ok := seen[perm]; ok {
			t.Errorf("权限标识重复: %s", perm)
		}
		seen[perm] = struct{}{}
	}
}
