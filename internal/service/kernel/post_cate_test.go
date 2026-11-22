package kernel

import (
	"net/http"
	"strings"
	"testing"
)

// postTitles 取文章列表的标题集合，便于断言过滤结果
func postTitles(t *testing.T, target string) []string {
	t.Helper()
	items, _ := callList(t, PostList, http.MethodGet, target)
	titles := make([]string, 0, len(items))
	for _, it := range items {
		mod, ok := it.(map[string]any)
		if !ok {
			t.Fatalf("文章结构异常: %v", it)
		}
		title, _ := mod["title"].(string)
		titles = append(titles, title)
	}
	return titles
}

// TestPostFilterByCate 按栏目过滤：默认只看本栏目，fuzzy=true 时模糊匹配含子孙栏目
func TestPostFilterByCate(t *testing.T) {
	setupKernelDB(t)
	tech := mustAddCate(t, `{"name":"技术","path":"tech"}`)
	goCate := mustAddCate(t, `{"name":"Go","path":"go","pid":"`+tech+`"}`)
	gin := mustAddCate(t, `{"name":"Gin","path":"gin","pid":"`+goCate+`"}`)
	life := mustAddCate(t, `{"name":"生活","path":"life"}`)

	for _, it := range []struct {
		title  string
		path   string
		cateId string
	}{
		{"Gin 入门", "p-gin", gin},
		{"Go 泛型", "p-go", goCate},
		{"生活随笔", "p-life", life},
		{"无栏目", "p-none", ""},
	} {
		callOk(t, PostAdd, http.MethodPost, "/api/posts",
			`{"title":"`+it.title+`","path":"`+it.path+`","cate_id":"`+it.cateId+`"}`, nil)
	}

	// 只看 Gin 本身
	if got := postTitles(t, "/api/posts?cate_id="+gin); len(got) != 1 || got[0] != "Gin 入门" {
		t.Errorf("按 Gin 过滤 = %v", got)
	}
	// 技术本身没有直属内容
	if got := postTitles(t, "/api/posts?cate_id="+tech); len(got) != 0 {
		t.Errorf("技术直属内容 = %v, want 空", got)
	}
	// 技术含子孙：Gin + Go
	got := postTitles(t, "/api/posts?cate_id="+tech+"&fuzzy=true")
	if len(got) != 2 {
		t.Fatalf("技术子树内容 = %v, want 2 条", got)
	}
	want := map[string]bool{"Gin 入门": true, "Go 泛型": true}
	for _, title := range got {
		if !want[title] {
			t.Errorf("意外的文章 %q", title)
		}
	}
	// 生活只有自己
	if got := postTitles(t, "/api/posts?cate_id="+life+"&fuzzy=true"); len(got) != 1 {
		t.Errorf("生活子树内容 = %v", got)
	}
}

// TestPostCateUsable 挂载内容时校验栏目存在且可用
func TestPostCateUsable(t *testing.T) {
	setupKernelDB(t)
	cateId := mustAddCate(t, `{"name":"技术","path":"tech"}`)

	// 栏目不存在
	msg := callFail(t, PostAdd, http.MethodPost, "/api/posts",
		`{"title":"X","path":"p1","cate_id":"not-exist"}`, nil)
	if msg != "栏目不存在" {
		t.Errorf("错误信息 = %q", msg)
	}

	// 栏目被禁用
	callOk(t, CateEdit, http.MethodPut, "/api/cates/"+cateId, `{"use":false}`, map[string]string{"id": cateId})
	msg = callFail(t, PostAdd, http.MethodPost, "/api/posts",
		`{"title":"X","path":"p1","cate_id":"`+cateId+`"}`, nil)
	if !strings.Contains(msg, "已禁用") {
		t.Errorf("错误信息 = %q", msg)
	}

	// 修改时同样校验
	callOk(t, CateEdit, http.MethodPut, "/api/cates/"+cateId, `{"use":true}`, map[string]string{"id": cateId})
	postId := callOk(t, PostAdd, http.MethodPost, "/api/posts",
		`{"title":"X","path":"p1","cate_id":"`+cateId+`"}`, nil).(string)
	callOk(t, CateEdit, http.MethodPut, "/api/cates/"+cateId, `{"use":false}`, map[string]string{"id": cateId})
	msg = callFail(t, PostEdit, http.MethodPut, "/api/posts/"+postId,
		`{"title":"X","path":"p1","cate_id":"`+cateId+`"}`, map[string]string{"id": postId})
	if !strings.Contains(msg, "已禁用") {
		t.Errorf("修改时错误信息 = %q", msg)
	}
}

// TestPostGetFillsCate 详情自动带出栏目信息
func TestPostGetFillsCate(t *testing.T) {
	setupKernelDB(t)
	cateId := mustAddCate(t, `{"name":"技术","path":"tech"}`)
	postId := callOk(t, PostAdd, http.MethodPost, "/api/posts",
		`{"title":"文章","path":"p1","cate_id":"`+cateId+`"}`, nil).(string)

	data := callOk(t, PostGet, http.MethodGet, "/api/posts/"+postId, "", map[string]string{"id": postId})
	mod, ok := data.(map[string]any)
	if !ok {
		t.Fatalf("详情结构异常: %v", data)
	}
	cate, ok := mod["cate"].(map[string]any)
	if !ok {
		t.Fatalf("未带出栏目: %v", mod["cate"])
	}
	if cate["name"] != "技术" || cate["path"] != "tech" {
		t.Errorf("栏目信息 = %v", cate)
	}
}
