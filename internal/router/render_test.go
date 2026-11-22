package router

import (
	"bytes"
	"html/template"
	"strings"
	"testing"

	"blog/internal/utils"
)

// loadViews 用与生产一致的 funcMap 加载视图目录
func loadViews(root string) (*template.Template, error) {
	return utils.LoadTmpl(root, funcMap)
}

// TestLoadViewsTemplates 项目自带的前台模板必须能被按名字加载
// 这层校验能提前发现模板命名规则回归：早期实现把相对路径截错，
// 导致 ExecuteTemplate("index.html") 找不到模板，前台全站白屏。
func TestLoadViewsTemplates(t *testing.T) {
	tmpl, err := loadViews("../../views")
	if err != nil {
		t.Fatalf("加载模板失败: %v", err)
	}
	names := make(map[string]bool, 16)
	for _, item := range tmpl.Templates() {
		names[item.Name()] = true
	}
	// 前台入口与被引用的公共片段
	for _, name := range []string{
		"index.html",
		"post.html",
		"page.html",
		"post-cate.html",
		"post-tag.html",
		"tags.html",
		"search.html",
		"archive.html",
		"note.html",
		"error.html",
		"tpl-naver.html",
		"tpl-linker.html",
		"tpl-footer.html",
		"tpl-script.html",
	} {
		if !names[name] {
			t.Errorf("模板 %q 未被加载", name)
		}
	}
	// 模板名不应带目录前缀
	for name := range names {
		if strings.Contains(name, "/views/") || strings.HasPrefix(name, "views") {
			t.Errorf("模板名带了多余前缀: %q", name)
		}
	}
}

// TestRenderIndex 首页模板能带着最小数据渲染出完整 HTML
// 覆盖 {{template}} 片段引用与 funcMap 两条链路
func TestRenderIndex(t *testing.T) {
	tmpl, err := loadViews("../../views")
	if err != nil {
		t.Fatalf("加载模板失败: %v", err)
	}
	data := map[string]any{
		"basic": map[string]any{
			"Title":   "测试站点",
			"LogoUrl": "/static/logo.png",
			"SiteUrl": "https://example.com",
			"Author":  "测",
		},
		"blog": map[string]any{
			"PageSize":  10,
			"Analytic":  "",
			"Comment":   "",
			"GithubUrl": "",
			"WeiboUrl":  "",
		},
		"appjs":  "/static/js/app.js",
		"appcss": "/static/css/app.css",
		"Posts":  []any{},
		"Naver":  struct{ Prev, Next string }{},
		"cates":  []any{},
	}
	w := bytes.NewBuffer(nil)
	if err := tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		t.Fatalf("渲染首页失败: %v", err)
	}
	html := w.String()
	if !strings.Contains(html, "<!DOCTYPE html>") {
		t.Error("输出缺少文档声明")
	}
	// 侧边导航与页头片段应已被内联
	if !strings.Contains(html, `id="sidebar"`) {
		t.Error("未包含 sidebar 片段")
	}
	if !strings.Contains(html, `id="header"`) {
		t.Error("未包含 header 片段")
	}
	if !strings.Contains(html, "测试站点") {
		t.Error("未渲染站点标题")
	}
}
