package kernel

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"blog/internal/model"
	"blog/internal/repo"
)

// TestCateAddDefaults 新增时未传的字段按服务端默认值落库
func TestCateAddDefaults(t *testing.T) {
	setupKernelDB(t)
	id := mustAddCate(t, `{"name":"技术"}`)
	mod := cateOf(t, id)
	if mod["name"] != "技术" {
		t.Errorf("name = %v", mod["name"])
	}
	// 公开路径留空时用主键小写自动生成
	if want := strings.ToLower(id); mod["path"] != want {
		t.Errorf("path = %v, want %v", mod["path"], want)
	}
	if levelOf(t, mod) != 1 {
		t.Errorf("顶层栏目 level = %d, want 1", levelOf(t, mod))
	}
	if mod["use"] != true || mod["show"] != true {
		t.Errorf("默认应可用且显示: use=%v show=%v", mod["use"], mod["show"])
	}
	if mod["kind"] != float64(1) {
		t.Errorf("默认归属类型 = %v, want 1", mod["kind"])
	}
	if mod["sort"] != float64(1000) {
		t.Errorf("默认排序 = %v, want 1000", mod["sort"])
	}
}

// TestCateAddLevelChain 层级按父链自动推导
func TestCateAddLevelChain(t *testing.T) {
	setupKernelDB(t)
	top := mustAddCate(t, `{"name":"技术","path":"tech"}`)
	lv2 := mustAddCate(t, `{"name":"Go","path":"go","pid":"`+top+`"}`)
	lv3 := mustAddCate(t, `{"name":"Gin","path":"gin","pid":"`+lv2+`"}`)

	for _, it := range []struct {
		id   string
		want int
	}{
		{top, 1},
		{lv2, 2},
		{lv3, 3},
	} {
		if got := levelOf(t, cateOf(t, it.id)); got != it.want {
			t.Errorf("level = %d, want %d", got, it.want)
		}
	}
}

// TestCateAddMaxLevel 超过层级上限时拒绝
func TestCateAddMaxLevel(t *testing.T) {
	setupKernelDB(t)
	lv1 := mustAddCate(t, `{"name":"L1"}`)
	lv2 := mustAddCate(t, `{"name":"L2","pid":"`+lv1+`"}`)
	lv3 := mustAddCate(t, `{"name":"L3","pid":"`+lv2+`"}`)
	msg := callFail(t, CateAdd, http.MethodPost, "/api/cates",
		`{"name":"L4","pid":"`+lv3+`"}`, nil)
	if !strings.Contains(msg, "最多") {
		t.Errorf("错误信息 = %q, 期望提示层级上限", msg)
	}
}

// TestCateAddParentNotExist 父级不存在时拒绝
func TestCateAddParentNotExist(t *testing.T) {
	setupKernelDB(t)
	msg := callFail(t, CateAdd, http.MethodPost, "/api/cates",
		`{"name":"X","pid":"not-exist"}`, nil)
	if msg != "父级栏目不存在" {
		t.Errorf("错误信息 = %q", msg)
	}
}

// TestCateAddNameRequired 名称必填且长度受限
func TestCateAddNameRequired(t *testing.T) {
	setupKernelDB(t)
	if msg := callFail(t, CateAdd, http.MethodPost, "/api/cates", `{"name":"  "}`, nil); msg == "" {
		t.Error("空名称应被拒绝")
	}
	long := strings.Repeat("名", 256)
	if msg := callFail(t, CateAdd, http.MethodPost, "/api/cates",
		`{"name":"`+long+`"}`, nil); msg == "" {
		t.Error("超长名称应被拒绝")
	}
}

// TestCateAddPathUnique 公开路径全站唯一
func TestCateAddPathUnique(t *testing.T) {
	setupKernelDB(t)
	mustAddCate(t, `{"name":"技术","path":"tech"}`)
	msg := callFail(t, CateAdd, http.MethodPost, "/api/cates", `{"name":"另一个","path":"tech"}`, nil)
	if msg != "公开路径已被占用" {
		t.Errorf("错误信息 = %q", msg)
	}
	// 不同路径可以共存
	mustAddCate(t, `{"name":"另一个","path":"other"}`)
}

// TestCateAddPathInvalid 公开路径格式校验
func TestCateAddPathInvalid(t *testing.T) {
	setupKernelDB(t)
	for _, path := range []string{"Bad Path", "-abc", "a/b", "中文"} {
		msg := callFail(t, CateAdd, http.MethodPost, "/api/cates",
			`{"name":"X","path":"`+path+`"}`, nil)
		if !strings.Contains(msg, "公开路径") {
			t.Errorf("path=%q 错误信息 = %q", path, msg)
		}
	}
}

// TestCateAddPathNormalized 路径自动转小写并去空格
func TestCateAddPathNormalized(t *testing.T) {
	setupKernelDB(t)
	id := mustAddCate(t, `{"name":"技术","path":"  TECH  "}`)
	if got := cateOf(t, id)["path"]; got != "tech" {
		t.Errorf("path = %v, want tech", got)
	}
}

// TestCateAddNameUniqueInSameParent 同级重名拒绝，跨级允许
func TestCateAddNameUniqueInSameParent(t *testing.T) {
	setupKernelDB(t)
	top := mustAddCate(t, `{"name":"技术","path":"tech"}`)
	msg := callFail(t, CateAdd, http.MethodPost, "/api/cates", `{"name":"技术"}`, nil)
	if msg != "同级下栏目名称已存在" {
		t.Errorf("顶层重名错误信息 = %q", msg)
	}
	// 不同父级下允许同名
	mustAddCate(t, `{"name":"技术","path":"tech2","pid":"`+top+`"}`)
	// 但同父级下仍不允许
	msg = callFail(t, CateAdd, http.MethodPost, "/api/cates",
		`{"name":"技术","pid":"`+top+`"}`, nil)
	if msg != "同级下栏目名称已存在" {
		t.Errorf("子级重名错误信息 = %q", msg)
	}
}

// TestCateEditPartial 局部更新：未提交的字段保持不变
func TestCateEditPartial(t *testing.T) {
	setupKernelDB(t)
	id := mustAddCate(t, `{"name":"技术","path":"tech","intro":"原描述","icon":"icon-a"}`)
	callOk(t, CateEdit, http.MethodPut, "/api/cates/"+id,
		`{"name":"技术2"}`, map[string]string{"id": id})
	mod := cateOf(t, id)
	if mod["name"] != "技术2" {
		t.Errorf("name = %v", mod["name"])
	}
	if mod["path"] != "tech" || mod["intro"] != "原描述" || mod["icon"] != "icon-a" {
		t.Errorf("未提交字段被改动: %v", mod)
	}
}

// TestCateEditIdempotent 重复提交相同值不报错
func TestCateEditIdempotent(t *testing.T) {
	setupKernelDB(t)
	id := mustAddCate(t, `{"name":"技术","path":"tech"}`)
	callOk(t, CateEdit, http.MethodPut, "/api/cates/"+id, `{"path":"tech"}`, map[string]string{"id": id})
	callOk(t, CateEdit, http.MethodPut, "/api/cates/"+id, `{}`, map[string]string{"id": id})
}

// TestCateEditPathUnique 路径唯一性校验需排除自身
func TestCateEditPathUnique(t *testing.T) {
	setupKernelDB(t)
	a := mustAddCate(t, `{"name":"A","path":"a"}`)
	mustAddCate(t, `{"name":"B","path":"b"}`)
	// 改成自己的路径允许
	callOk(t, CateEdit, http.MethodPut, "/api/cates/"+a, `{"path":"a"}`, map[string]string{"id": a})
	// 占用别人的路径拒绝
	msg := callFail(t, CateEdit, http.MethodPut, "/api/cates/"+a, `{"path":"b"}`, map[string]string{"id": a})
	if msg != "公开路径已被占用" {
		t.Errorf("错误信息 = %q", msg)
	}
}

// TestCateEditMove 移动节点后整棵子树层级重算
func TestCateEditMove(t *testing.T) {
	setupKernelDB(t)
	a := mustAddCate(t, `{"name":"A"}`)
	b := mustAddCate(t, `{"name":"B","pid":"`+a+`"}`)
	c := mustAddCate(t, `{"name":"C","pid":"`+b+`"}`)

	// 把 B 提到顶层，C 跟着上移一级
	callOk(t, CateEdit, http.MethodPut, "/api/cates/"+b, `{"pid":""}`, map[string]string{"id": b})
	for _, it := range []struct {
		id   string
		want int
	}{
		{a, 1},
		{b, 1},
		{c, 2},
	} {
		if got := levelOf(t, cateOf(t, it.id)); got != it.want {
			t.Errorf("移动后 level = %d, want %d", got, it.want)
		}
	}
	// 再放回 A 下，层级恢复
	callOk(t, CateEdit, http.MethodPut, "/api/cates/"+b, `{"pid":"`+a+`"}`, map[string]string{"id": b})
	if got := levelOf(t, cateOf(t, c)); got != 3 {
		t.Errorf("移回后 C level = %d, want 3", got)
	}
}

// TestCateEditMoveCycle 不能移动到自身或自己的子孙下
func TestCateEditMoveCycle(t *testing.T) {
	setupKernelDB(t)
	a := mustAddCate(t, `{"name":"A"}`)
	b := mustAddCate(t, `{"name":"B","pid":"`+a+`"}`)
	b2 := mustAddCate(t, `{"name":"B2","pid":"`+b+`"}`)

	for _, pid := range []string{a, b2} {
		msg := callFail(t, CateEdit, http.MethodPut, "/api/cates/"+a,
			`{"pid":"`+pid+`"}`, map[string]string{"id": a})
		if msg != "不能移动到自身或其子级下" {
			t.Errorf("pid=%s 错误信息 = %q", pid, msg)
		}
	}
}

// TestCateEditMoveExceedLevel 移动后子树超过层级上限时拒绝
func TestCateEditMoveExceedLevel(t *testing.T) {
	setupKernelDB(t)
	a := mustAddCate(t, `{"name":"A"}`)
	b := mustAddCate(t, `{"name":"B","pid":"`+a+`"}`)
	mustAddCate(t, `{"name":"C","pid":"`+b+`"}`) // 挂在 A 下的子树高度 3
	top := mustAddCate(t, `{"name":"TOP"}`)

	// A 移到 TOP 下会变成 2..4 层，超过 3 级
	msg := callFail(t, CateEdit, http.MethodPut, "/api/cates/"+a,
		`{"pid":"`+top+`"}`, map[string]string{"id": a})
	if !strings.Contains(msg, "层级上限") {
		t.Errorf("错误信息 = %q", msg)
	}
}

// TestCateEditNameUnique 改名时同级判重
func TestCateEditNameUnique(t *testing.T) {
	setupKernelDB(t)
	a := mustAddCate(t, `{"name":"A"}`)
	mustAddCate(t, `{"name":"B"}`)
	msg := callFail(t, CateEdit, http.MethodPut, "/api/cates/"+a, `{"name":"B"}`, map[string]string{"id": a})
	if msg != "同级下栏目名称已存在" {
		t.Errorf("错误信息 = %q", msg)
	}
	// 改成自己的名字允许
	callOk(t, CateEdit, http.MethodPut, "/api/cates/"+a, `{"name":"A"}`, map[string]string{"id": a})
}

// TestCateEditNotFound 修改不存在的栏目
func TestCateEditNotFound(t *testing.T) {
	setupKernelDB(t)
	msg := callFail(t, CateEdit, http.MethodPut, "/api/cates/nope", `{"name":"X"}`, map[string]string{"id": "nope"})
	if msg != "栏目不存在" {
		t.Errorf("错误信息 = %q", msg)
	}
}

// TestCateDropWithChildren 有子级时禁止删除
func TestCateDropWithChildren(t *testing.T) {
	setupKernelDB(t)
	a := mustAddCate(t, `{"name":"A"}`)
	mustAddCate(t, `{"name":"B","pid":"`+a+`"}`)
	msg := callFail(t, CateDrop, http.MethodDelete, "/api/cates/"+a, "", map[string]string{"id": a})
	if msg != "请先删除子级栏目" {
		t.Errorf("错误信息 = %q", msg)
	}
}

// TestCateDropWithPosts 有内容时禁止删除
func TestCateDropWithPosts(t *testing.T) {
	setupKernelDB(t)
	a := mustAddCate(t, `{"name":"A"}`)
	b := mustAddCate(t, `{"name":"B","pid":"`+a+`"}`)
	// 内容挂在子级，删除父级同样应被内容保护拦住
	callOk(t, PostAdd, http.MethodPost, "/api/posts",
		`{"title":"文章","path":"p1","cate_id":"`+b+`"}`, nil)
	msg := callFail(t, CateDrop, http.MethodDelete, "/api/cates/"+b, "", map[string]string{"id": b})
	if !strings.Contains(msg, "还有内容") {
		t.Errorf("错误信息 = %q", msg)
	}
	if got := callFail(t, CateDrop, http.MethodDelete, "/api/cates/"+a, "", map[string]string{"id": a}); got != "请先删除子级栏目" {
		t.Errorf("父级错误信息 = %q", got)
	}
}

// TestCateDropSuccess 删除叶子栏目成功且缓存同步
func TestCateDropSuccess(t *testing.T) {
	setupKernelDB(t)
	a := mustAddCate(t, `{"name":"A","path":"a"}`)
	callOk(t, CateDrop, http.MethodDelete, "/api/cates/"+a, "", map[string]string{"id": a})
	msg := callFail(t, CateGet, http.MethodGet, "/api/cates/"+a, "", map[string]string{"id": a})
	if msg != "栏目不存在" {
		t.Errorf("删除后仍可查询: %q", msg)
	}
	// 路径被释放，可以重新使用
	mustAddCate(t, `{"name":"A2","path":"a"}`)
}

// TestCateGetDetail 详情带子级与内容数量
func TestCateGetDetail(t *testing.T) {
	setupKernelDB(t)
	a := mustAddCate(t, `{"name":"A"}`)
	b := mustAddCate(t, `{"name":"B","pid":"`+a+`"}`)
	callOk(t, PostAdd, http.MethodPost, "/api/posts", `{"title":"文章","path":"p1","cate_id":"`+b+`"}`, nil)

	data := callOk(t, CateGet, http.MethodGet, "/api/cates/"+a, "", map[string]string{"id": a})
	mod, ok := data.(map[string]any)
	if !ok {
		t.Fatalf("详情结构异常: %v", data)
	}
	children, ok := mod["children"].([]any)
	if !ok || len(children) != 1 {
		t.Fatalf("子级 = %v", mod["children"])
	}
	// 数量含子树
	if n := atoi(mod["count"]); n != 1 {
		t.Errorf("count = %d, want 1", n)
	}
}

// TestCateTreeNested 树形返回按层级嵌套
func TestCateTreeNested(t *testing.T) {
	setupKernelDB(t)
	a := mustAddCate(t, `{"name":"A"}`)
	b := mustAddCate(t, `{"name":"B","pid":"`+a+`"}`)
	c := mustAddCate(t, `{"name":"C","pid":"`+b+`"}`)

	data := callOk(t, CateTree, http.MethodGet, "/api/cates/tree", "", nil)
	nodes := cateNodes(t, data)
	if len(nodes) != 1 || nodes[0].Id != a {
		t.Fatalf("顶层 = %+v", nodes)
	}
	if len(nodes[0].Children) != 1 || nodes[0].Children[0].Id != b {
		t.Fatalf("二级 = %+v", nodes[0].Children)
	}
	if len(nodes[0].Children[0].Children) != 1 || nodes[0].Children[0].Children[0].Id != c {
		t.Fatalf("三级 = %+v", nodes[0].Children[0].Children)
	}
}

// TestCateTreeFlat 扁平返回带 level，便于下拉缩进
func TestCateTreeFlat(t *testing.T) {
	setupKernelDB(t)
	a := mustAddCate(t, `{"name":"A"}`)
	b := mustAddCate(t, `{"name":"B","pid":"`+a+`"}`)

	data := callOk(t, CateTree, http.MethodGet, "/api/cates/tree?flat=true", "", nil)
	nodes := cateNodes(t, data)
	if len(nodes) != 2 {
		t.Fatalf("扁平节点数 = %d, want 2", len(nodes))
	}
	if nodes[0].Id != a || nodes[0].Level != 1 {
		t.Errorf("第 1 个节点 = %+v", nodes[0])
	}
	if nodes[1].Id != b || nodes[1].Level != 2 {
		t.Errorf("第 2 个节点 = %+v", nodes[1])
	}
	if len(nodes[0].Children) != 0 {
		t.Error("扁平模式不应带子级")
	}
}

// TestCateTreeRoot 顶部附加「全部」节点
func TestCateTreeRoot(t *testing.T) {
	setupKernelDB(t)
	mustAddCate(t, `{"name":"A"}`)
	data := callOk(t, CateTree, http.MethodGet, "/api/cates/tree?root=true", "", nil)
	nodes := cateNodes(t, data)
	if len(nodes) != 2 {
		t.Fatalf("节点数 = %d, want 2", len(nodes))
	}
	if nodes[0].Name != "全部栏目" || nodes[0].Level != 0 || !nodes[0].Inner {
		t.Errorf("顶部节点 = %+v", nodes[0])
	}
	if nodes[1].Name != "A" {
		t.Errorf("第二个节点 = %+v", nodes[1])
	}
}

// TestCateTreeFilter 按可用/显示过滤，子级同步过滤
func TestCateTreeFilter(t *testing.T) {
	setupKernelDB(t)
	a := mustAddCate(t, `{"name":"A"}`)
	mustAddCate(t, `{"name":"B","pid":"`+a+`"}`)
	hidden := mustAddCate(t, `{"name":"隐藏"}`)
	callOk(t, CateEdit, http.MethodPut, "/api/cates/"+hidden, `{"show":false}`, map[string]string{"id": hidden})

	data := callOk(t, CateTree, http.MethodGet, "/api/cates/tree?show=true", "", nil)
	nodes := cateNodes(t, data)
	if len(nodes) != 1 || nodes[0].Id != a {
		t.Fatalf("按 show 过滤结果 = %+v", nodes)
	}

	// 禁用其中一个子级后，其自身从结果中消失
	b := nodes[0].Children[0].Id
	callOk(t, CateEdit, http.MethodPut, "/api/cates/"+b, `{"use":false}`, map[string]string{"id": b})
	data = callOk(t, CateTree, http.MethodGet, "/api/cates/tree?use=true", "", nil)
	nodes = cateNodes(t, data)
	if len(nodes) != 2 { // A 与 隐藏（隐藏只是 show=false，use 仍为 true）
		t.Fatalf("按 use 过滤结果 = %+v", nodes)
	}
	if len(nodes[0].Children) != 0 {
		t.Errorf("禁用的子级应被过滤: %+v", nodes[0].Children)
	}
}

// TestCateAddAllFields 一次性提交全部字段并逐项回读
func TestCateAddAllFields(t *testing.T) {
	setupKernelDB(t)
	body := `{
		"name":"技术",
		"path":"tech",
		"kind":3,
		"url":"https://example.com",
		"intro":"描述",
		"icon":"icon-a",
		"color":"#fff",
		"sort":10,
		"use":false,
		"show":false,
		"auth":true,
		"itpl":"i.html",
		"ltpl":"l.html",
		"dtpl":"d.html"
	}`
	id := mustAddCate(t, body)
	mod := cateOf(t, id)
	want := map[string]any{
		"name":  "技术",
		"path":  "tech",
		"kind":  float64(3),
		"url":   "https://example.com",
		"intro": "描述",
		"icon":  "icon-a",
		"color": "#fff",
		"sort":  float64(10),
		"use":   false,
		"show":  false,
		"auth":  true,
		"itpl":  "i.html",
		"ltpl":  "l.html",
		"dtpl":  "d.html",
	}
	for key, expect := range want {
		if mod[key] != expect {
			t.Errorf("%s = %v, want %v", key, mod[key], expect)
		}
	}
}

// TestCateEditAllFields 逐字段修改均可落库
func TestCateEditAllFields(t *testing.T) {
	setupKernelDB(t)
	id := mustAddCate(t, `{"name":"A","path":"a"}`)
	body := `{
		"name":"B",
		"path":"b",
		"kind":3,
		"url":"https://example.com",
		"intro":"描述",
		"icon":"icon-a",
		"color":"#fff",
		"sort":10,
		"use":false,
		"show":false,
		"auth":true,
		"itpl":"i.html",
		"ltpl":"l.html",
		"dtpl":"d.html"
	}`
	callOk(t, CateEdit, http.MethodPut, "/api/cates/"+id, body, map[string]string{"id": id})
	mod := cateOf(t, id)
	want := map[string]any{
		"name":  "B",
		"path":  "b",
		"kind":  float64(3),
		"url":   "https://example.com",
		"intro": "描述",
		"icon":  "icon-a",
		"color": "#fff",
		"sort":  float64(10),
		"use":   false,
		"show":  false,
		"auth":  true,
		"itpl":  "i.html",
		"ltpl":  "l.html",
		"dtpl":  "d.html",
	}
	for key, expect := range want {
		if mod[key] != expect {
			t.Errorf("%s = %v, want %v", key, mod[key], expect)
		}
	}
}

// TestCateEditInvalidName 修改时名称不能为空或超长
func TestCateEditInvalidName(t *testing.T) {
	setupKernelDB(t)
	id := mustAddCate(t, `{"name":"A"}`)
	if msg := callFail(t, CateEdit, http.MethodPut, "/api/cates/"+id,
		`{"name":"  "}`, map[string]string{"id": id}); msg == "" {
		t.Error("空名称应被拒绝")
	}
	long := strings.Repeat("名", 256)
	if msg := callFail(t, CateEdit, http.MethodPut, "/api/cates/"+id,
		`{"name":"`+long+`"}`, map[string]string{"id": id}); msg == "" {
		t.Error("超长名称应被拒绝")
	}
}

// TestCateDropInner 内置栏目禁止删除
func TestCateDropInner(t *testing.T) {
	setupKernelDB(t)
	id := mustAddCate(t, `{"name":"内置","path":"inner"}`)
	// 内置标记只由服务端/种子设置，这里直接改库模拟
	if err := repo.CateEdit(context.Background(), &model.Cate{Id: id, Inner: true}, "inner"); err != nil {
		t.Fatalf("CateEdit: %v", err)
	}
	msg := callFail(t, CateDrop, http.MethodDelete, "/api/cates/"+id, "", map[string]string{"id": id})
	if msg != "内置栏目不可删除" {
		t.Errorf("错误信息 = %q", msg)
	}
}

// TestCateDropEmptyId 缺少主键时拒绝
func TestCateDropEmptyId(t *testing.T) {
	setupKernelDB(t)
	msg := callFail(t, CateDrop, http.MethodDelete, "/api/cates/", "", map[string]string{"id": ""})
	if msg != "参数有误" {
		t.Errorf("错误信息 = %q", msg)
	}
}

// TestCateListFilters 列表过滤与分页
func TestCateListFilters(t *testing.T) {
	setupKernelDB(t)
	mustAddCate(t, `{"name":"技术","path":"tech"}`)
	mustAddCate(t, `{"name":"生活","path":"life"}`)
	mustAddCate(t, `{"name":"笔记","path":"note","kind":3}`)

	items, count := callList(t, CateList, http.MethodGet, "/api/cates?page=0&size=0")
	if count != 3 || len(items) != 3 {
		t.Fatalf("全部栏目 = %d, count=%d", len(items), count)
	}

	items, count = callList(t, CateList, http.MethodGet, "/api/cates?page=0&size=0&kind=3")
	if count != 1 || len(items) != 1 {
		t.Fatalf("按 kind 过滤 = %d, count=%d", len(items), count)
	}

	items, _ = callList(t, CateList, http.MethodGet, "/api/cates?page=0&size=0&q=tec")
	if len(items) != 1 {
		t.Fatalf("按关键字过滤 = %d", len(items))
	}

	// 分页
	items, count = callList(t, CateList, http.MethodGet, "/api/cates?page=1&size=2")
	if len(items) != 2 || count != 3 {
		t.Fatalf("分页 = %d, count=%d", len(items), count)
	}
}
