package kernel

import (
	"context"
	"net/http"
	"strings"
	"testing"

	"blog/internal/model"
	"blog/internal/repo"
	"blog/pkg/reply"
)

// postOf 通过详情接口取内容
func postOf(t *testing.T, id string) map[string]any {
	t.Helper()
	data := callOk(t, PostGet, http.MethodGet, "/api/posts/"+id, "", map[string]string{"id": id})
	mod, ok := data.(map[string]any)
	if !ok {
		t.Fatalf("详情结构异常: %v", data)
	}
	return mod
}

// addPost 走接口新增内容并返回主键
func addPost(t *testing.T, body string) string {
	t.Helper()
	id, ok := callOk(t, PostAdd, http.MethodPost, "/api/posts", body, nil).(string)
	if !ok || id == "" {
		t.Fatal("新增内容未返回主键")
	}
	return id
}

// TestPostAddDefaults 新增时未传的字段按服务端默认值落库
func TestPostAddDefaults(t *testing.T) {
	setupKernelDB(t)
	id := addPost(t, `{"title":"标题","path":"p1"}`)
	mod := postOf(t, id)
	if mod["kind"] != float64(1) {
		t.Errorf("默认类型 = %v, want 1", mod["kind"])
	}
	if mod["status"] != float64(1) {
		t.Errorf("默认状态 = %v, want 1（草稿）", mod["status"])
	}
	if mod["show"] != true {
		t.Errorf("默认应前台可见: %v", mod["show"])
	}
	if mod["istop"] != false {
		t.Errorf("默认不置顶: %v", mod["istop"])
	}
	if mod["hits"] != float64(0) {
		t.Errorf("默认浏览量 = %v, want 0", mod["hits"])
	}
	if mod["inner"] != false {
		t.Errorf("默认非内置: %v", mod["inner"])
	}
}

// TestPostAddKindStatus 类型与状态校验
func TestPostAddKindStatus(t *testing.T) {
	setupKernelDB(t)
	msg := callFail(t, PostAdd, http.MethodPost, "/api/posts",
		`{"title":"x","path":"p1","kind":9}`, nil)
	if msg != "内容类型无效" {
		t.Errorf("类型错误信息 = %q", msg)
	}
	msg = callFail(t, PostAdd, http.MethodPost, "/api/posts",
		`{"title":"x","path":"p1","status":9}`, nil)
	if msg != "状态无效" {
		t.Errorf("状态错误信息 = %q", msg)
	}
}

// TestPostAddRequiresContentWhenPublished 已发布必须填正文，草稿可空
func TestPostAddRequiresContentWhenPublished(t *testing.T) {
	setupKernelDB(t)
	callOk(t, PostAdd, http.MethodPost, "/api/posts", `{"title":"草稿","path":"d1"}`, nil)
	msg := callFail(t, PostAdd, http.MethodPost, "/api/posts",
		`{"title":"发布","path":"p1","status":2}`, nil)
	if !strings.Contains(msg, "正文") {
		t.Errorf("错误信息 = %q", msg)
	}
	// markdown 或 richtext 任一非空即可
	callOk(t, PostAdd, http.MethodPost, "/api/posts",
		`{"title":"发布","path":"p1","status":2,"markdown":"# hi"}`, nil)
	callOk(t, PostAdd, http.MethodPost, "/api/posts",
		`{"title":"发布2","path":"p2","status":2,"richtext":"<p>hi</p>"}`, nil)
}

// TestPostAddPathValidation 路径格式与唯一性
func TestPostAddPathValidation(t *testing.T) {
	setupKernelDB(t)
	for _, path := range []string{"a b", "-abc", "a/b", "中文", ""} {
		msg := callFail(t, PostAdd, http.MethodPost, "/api/posts",
			`{"title":"x","path":"`+path+`"}`, nil)
		if !strings.Contains(msg, "访问路径") && !strings.Contains(msg, "标题") {
			t.Errorf("path=%q 错误信息 = %q", path, msg)
		}
	}
	addPost(t, `{"title":"A","path":"dup"}`)
	msg := callFail(t, PostAdd, http.MethodPost, "/api/posts", `{"title":"B","path":"dup"}`, nil)
	if msg != "访问路径已存在" {
		t.Errorf("重复路径错误信息 = %q", msg)
	}
}

// TestPostAddTitleRequired 标题必填且有长度上限
func TestPostAddTitleRequired(t *testing.T) {
	setupKernelDB(t)
	if msg := callFail(t, PostAdd, http.MethodPost, "/api/posts",
		`{"title":"   ","path":"p1"}`, nil); msg == "" {
		t.Error("空标题应被拒绝")
	}
	long := strings.Repeat("字", 256)
	if msg := callFail(t, PostAdd, http.MethodPost, "/api/posts",
		`{"title":"`+long+`","path":"p1"}`, nil); msg == "" {
		t.Error("超长标题应被拒绝")
	}
}

// TestPostAddPidOnlyForNote 只有笔记支持父子层级
func TestPostAddPidOnlyForNote(t *testing.T) {
	setupKernelDB(t)
	parent := addPost(t, `{"title":"父","path":"parent"}`)
	msg := callFail(t, PostAdd, http.MethodPost, "/api/posts",
		`{"title":"子","path":"child","pid":"`+parent+`"}`, nil)
	if msg != "只有笔记支持父子层级" {
		t.Errorf("错误信息 = %q", msg)
	}
	noteParent := addPost(t, `{"title":"笔记父","path":"note-parent","kind":3}`)
	callOk(t, PostAdd, http.MethodPost, "/api/posts",
		`{"title":"笔记子","path":"note-child","kind":3,"pid":"`+noteParent+`"}`, nil)
}

// TestPostAddFullFields 全字段新增并逐项回读
func TestPostAddFullFields(t *testing.T) {
	setupKernelDB(t)
	cateId := mustAddCate(t, `{"name":"技术","path":"tech"}`)
	// 用内置根账号作为操作者，再建一个内容作者用于验证创建者自动写入
	root, err := repo.UserNum(context.Background(), "root")
	if err != nil {
		t.Fatalf("取根账号失败: %v", err)
	}
	out := callAs(t, root.Id, UserAdd, http.MethodPost, "/api/users",
		`{"num":"author","passwd":"123456","role_id":"`+model.RoleRootId+`"}`, nil)
	if out.Code != reply.CodeSucc {
		t.Fatalf("建用户失败: code=%d msg=%s", out.Code, out.Msg)
	}
	uid, _ := out.Data.(string)
	body := `{
		"title":"完整内容",
		"path":"full",
		"cate_id":"` + cateId + `",
		"kind":2,
		"status":2,
		"summary":"摘要",
		"author":"作者",
		"origin":"来源",
		"icon":"/static/a.png",
		"markdown":"# hi",
		"richtext":"<p>hi</p>",
		"allow":true,
		"show":false,
		"istop":true
	}`
	out2 := callAs(t, uid, PostAdd, http.MethodPost, "/api/posts", body, nil)
	if out2.Code != reply.CodeSucc {
		t.Fatalf("新增失败: code=%d msg=%s", out2.Code, out2.Msg)
	}
	id, _ := out2.Data.(string)
	if id == "" {
		t.Fatal("新增未返回主键")
	}
	mod := postOf(t, id)
	want := map[string]any{
		"title":    "完整内容",
		"path":     "full",
		"cate_id":  cateId,
		"kind":     float64(2),
		"status":   float64(2),
		"summary":  "摘要",
		"author":   "作者",
		"origin":   "来源",
		"icon":     "/static/a.png",
		"markdown": "# hi",
		"richtext": "<p>hi</p>",
		"allow":    true,
		"show":     false,
		"istop":    true,
	}
	for key, expect := range want {
		if mod[key] != expect {
			t.Errorf("%s = %v, want %v", key, mod[key], expect)
		}
	}
	// 创建者应自动写入
	if got, _ := mod["uid"].(string); got != uid {
		t.Errorf("uid = %q, want %q", got, uid)
	}
}

// TestPostEditPartial 局部更新：未提交的字段保持不变
func TestPostEditPartial(t *testing.T) {
	setupKernelDB(t)
	id := addPost(t, `{"title":"原题","path":"p1","summary":"原摘要","author":"原作者"}`)
	callOk(t, PostEdit, http.MethodPut, "/api/posts/"+id, `{"title":"新题"}`, map[string]string{"id": id})
	mod := postOf(t, id)
	if mod["title"] != "新题" {
		t.Errorf("title = %v", mod["title"])
	}
	if mod["path"] != "p1" || mod["summary"] != "原摘要" || mod["author"] != "原作者" {
		t.Errorf("未提交字段被改动: %v", mod)
	}
}

// TestPostEditIdempotent 重复提交相同值不报错
func TestPostEditIdempotent(t *testing.T) {
	setupKernelDB(t)
	id := addPost(t, `{"title":"A","path":"a"}`)
	callOk(t, PostEdit, http.MethodPut, "/api/posts/"+id, `{"path":"a"}`, map[string]string{"id": id})
	callOk(t, PostEdit, http.MethodPut, "/api/posts/"+id, `{}`, map[string]string{"id": id})
}

// TestPostEditPathUnique 路径唯一性校验需排除自身
func TestPostEditPathUnique(t *testing.T) {
	setupKernelDB(t)
	a := addPost(t, `{"title":"A","path":"a"}`)
	addPost(t, `{"title":"B","path":"b"}`)
	callOk(t, PostEdit, http.MethodPut, "/api/posts/"+a, `{"path":"a"}`, map[string]string{"id": a})
	msg := callFail(t, PostEdit, http.MethodPut, "/api/posts/"+a, `{"path":"b"}`, map[string]string{"id": a})
	if msg != "访问路径已存在" {
		t.Errorf("错误信息 = %q", msg)
	}
}

// TestPostEditKindImmutable Kind 不在可改字段内，提交后被忽略
func TestPostEditKindImmutable(t *testing.T) {
	setupKernelDB(t)
	id := addPost(t, `{"title":"A","path":"a","kind":1}`)
	callOk(t, PostEdit, http.MethodPut, "/api/posts/"+id, `{"kind":3,"title":"B"}`, map[string]string{"id": id})
	if got := postOf(t, id)["kind"]; got != float64(1) {
		t.Errorf("kind = %v, want 1（类型不可变）", got)
	}
}

// TestPostEditPublishRequiresContent 改为已发布时必须已有正文
func TestPostEditPublishRequiresContent(t *testing.T) {
	setupKernelDB(t)
	id := addPost(t, `{"title":"草稿","path":"p1"}`)
	msg := callFail(t, PostEdit, http.MethodPut, "/api/posts/"+id, `{"status":2}`, map[string]string{"id": id})
	if !strings.Contains(msg, "正文") {
		t.Errorf("错误信息 = %q", msg)
	}
	// 补正文后可以发布
	callOk(t, PostEdit, http.MethodPut, "/api/posts/"+id,
		`{"status":2,"markdown":"# hi"}`, map[string]string{"id": id})
	if got := postOf(t, id)["status"]; got != float64(2) {
		t.Errorf("status = %v, want 2", got)
	}
}

// TestPostEditNotFound 修改不存在的内容
func TestPostEditNotFound(t *testing.T) {
	setupKernelDB(t)
	msg := callFail(t, PostEdit, http.MethodPut, "/api/posts/nope", `{"title":"x"}`, map[string]string{"id": "nope"})
	if msg != "内容不存在" {
		t.Errorf("错误信息 = %q", msg)
	}
}

// TestPostEditTop 置顶开关
func TestPostEditTop(t *testing.T) {
	setupKernelDB(t)
	id := addPost(t, `{"title":"A","path":"a"}`)
	callOk(t, PostEditTop, http.MethodPut, "/api/posts/"+id+"/top", `{"istop":true}`, map[string]string{"id": id})
	if got := postOf(t, id)["istop"]; got != true {
		t.Errorf("istop = %v, want true", got)
	}
	callOk(t, PostEditTop, http.MethodPut, "/api/posts/"+id+"/top", `{"istop":false}`, map[string]string{"id": id})
	if got := postOf(t, id)["istop"]; got != false {
		t.Errorf("istop = %v, want false", got)
	}
	msg := callFail(t, PostEditTop, http.MethodPut, "/api/posts/nope/top", `{"istop":true}`, map[string]string{"id": "nope"})
	if msg != "内容不存在" {
		t.Errorf("错误信息 = %q", msg)
	}
}

// TestPostDropInner 内置内容不可删除
func TestPostDropInner(t *testing.T) {
	setupKernelDB(t)
	id := addPost(t, `{"title":"关于","path":"about"}`)
	// 内置标记只由服务端维护，这里直接改库模拟
	upd := &model.Post{Id: id, Inner: true}
	if err := repo.PostEdit(context.Background(), upd, "inner"); err != nil {
		t.Fatalf("PostEdit: %v", err)
	}
	msg := callFail(t, PostDrop, http.MethodDelete, "/api/posts/"+id, "", map[string]string{"id": id})
	if msg != "内置内容不可删除" {
		t.Errorf("错误信息 = %q", msg)
	}
}

// TestPostEditInnerPath 内置内容不允许改路径
func TestPostEditInnerPath(t *testing.T) {
	setupKernelDB(t)
	id := addPost(t, `{"title":"关于","path":"about"}`)
	upd := &model.Post{Id: id, Inner: true}
	if err := repo.PostEdit(context.Background(), upd, "inner"); err != nil {
		t.Fatalf("PostEdit: %v", err)
	}
	msg := callFail(t, PostEdit, http.MethodPut, "/api/posts/"+id, `{"path":"new"}`, map[string]string{"id": id})
	if msg != "内置内容不可修改访问路径" {
		t.Errorf("错误信息 = %q", msg)
	}
}

// TestPostDropWithChildren 有子级时先删子级
func TestPostDropWithChildren(t *testing.T) {
	setupKernelDB(t)
	parent := addPost(t, `{"title":"父","path":"parent","kind":3}`)
	child := addPost(t, `{"title":"子","path":"child","kind":3,"pid":"`+parent+`"}`)
	msg := callFail(t, PostDrop, http.MethodDelete, "/api/posts/"+parent, "", map[string]string{"id": parent})
	if msg != "请先删除子级内容" {
		t.Errorf("错误信息 = %q", msg)
	}
	// 先删子级再删父级
	callOk(t, PostDrop, http.MethodDelete, "/api/posts/"+child, "", map[string]string{"id": child})
	callOk(t, PostDrop, http.MethodDelete, "/api/posts/"+parent, "", map[string]string{"id": parent})
}

// TestPostDropEmptyId 缺少主键
func TestPostDropEmptyId(t *testing.T) {
	setupKernelDB(t)
	msg := callFail(t, PostDrop, http.MethodDelete, "/api/posts/", "", map[string]string{"id": ""})
	if msg != "参数有误" {
		t.Errorf("错误信息 = %q", msg)
	}
}

// TestPostListFilters 列表过滤：类型/状态/可见/置顶/关键字
func TestPostListFilters(t *testing.T) {
	setupKernelDB(t)
	addPost(t, `{"title":"文章一","path":"p1","status":2,"markdown":"x"}`)
	addPost(t, `{"title":"文章二","path":"p2"}`)
	addPost(t, `{"title":"页面","path":"pg","kind":2}`)
	hiddenId := addPost(t, `{"title":"隐藏","path":"ph","show":false}`)
	topId := addPost(t, `{"title":"置顶","path":"pt","istop":true}`)

	items, count := callList(t, PostList, http.MethodGet, "/api/posts?page=0&size=0")
	if count != 5 || len(items) != 5 {
		t.Fatalf("全部 = %d, count=%d", len(items), count)
	}

	base := "/api/posts?page=0&size=0"
	items, _ = callList(t, PostList, http.MethodGet, base+"&kind=2")
	if len(items) != 1 {
		t.Errorf("按类型过滤 = %d", len(items))
	}
	items, _ = callList(t, PostList, http.MethodGet, base+"&status=2")
	if len(items) != 1 {
		t.Errorf("按状态过滤 = %d", len(items))
	}
	items, _ = callList(t, PostList, http.MethodGet, base+"&show=false")
	if len(items) != 1 {
		t.Errorf("按可见过滤 = %d", len(items))
	}
	items, _ = callList(t, PostList, http.MethodGet, base+"&istop=true")
	if len(items) != 1 {
		t.Errorf("按置顶过滤 = %d", len(items))
	}
	items, _ = callList(t, PostList, http.MethodGet, base+"&q=文章")
	if len(items) != 2 {
		t.Errorf("按关键字过滤 = %d", len(items))
	}
	// 置顶内容排在最前
	items, _ = callList(t, PostList, http.MethodGet, base)
	if first, ok := items[0].(map[string]any); !ok || first["id"] != topId {
		t.Errorf("首条应为置顶内容, got %v", items[0])
	}
	_ = hiddenId
}

// TestPostGetDetail 详情带标签与栏目
func TestPostGetDetail(t *testing.T) {
	setupKernelDB(t)
	cateId := mustAddCate(t, `{"name":"技术","path":"tech"}`)
	tagId := callOk(t, TagAdd, http.MethodPost, "/api/tags", `{"name":"go"}`, nil).(string)
	id := addPost(t, `{"title":"文章","path":"p1","cate_id":"`+cateId+`","markdown":"# hi"}`)
	callOk(t, PostTagEdit, http.MethodPut, "/api/posts/"+id+"/tags",
		`{"tag_ids":["`+tagId+`"]}`, map[string]string{"id": id})

	mod := postOf(t, id)
	if mod["markdown"] != "# hi" {
		t.Errorf("详情应返回正文: %v", mod["markdown"])
	}
	cate, ok := mod["cate"].(map[string]any)
	if !ok || cate["name"] != "技术" {
		t.Errorf("栏目未返回: %v", mod["cate"])
	}
	tags, ok := mod["tags"].([]any)
	if !ok || len(tags) != 1 {
		t.Errorf("标签未返回: %v", mod["tags"])
	}
}

// TestPostTagInvalid 标签全量替换时校验 id 有效性
func TestPostTagInvalid(t *testing.T) {
	setupKernelDB(t)
	id := addPost(t, `{"title":"A","path":"a"}`)
	msg := callFail(t, PostTagEdit, http.MethodPut, "/api/posts/"+id+"/tags",
		`{"tag_ids":["nope"]}`, map[string]string{"id": id})
	if msg != "存在无效的标签id" {
		t.Errorf("错误信息 = %q", msg)
	}
	// 传空数组等于清空
	callOk(t, PostTagEdit, http.MethodPut, "/api/posts/"+id+"/tags", `{"tag_ids":[]}`, map[string]string{"id": id})
}

// TestPostNeighbors 上下篇接口
func TestPostNeighbors(t *testing.T) {
	setupKernelDB(t)
	older := addPost(t, `{"title":"旧","path":"older","status":2,"markdown":"x"}`)
	newer := addPost(t, `{"title":"新","path":"newer","status":2,"markdown":"x"}`)
	// 修正创建时间，保证顺序可控
	ctx := context.Background()
	if err := repo.PostEdit(ctx, &model.Post{Id: older, Created: 1000, Updated: 1000}, "created", "updated"); err != nil {
		t.Fatalf("PostEdit: %v", err)
	}
	if err := repo.PostEdit(ctx, &model.Post{Id: newer, Created: 2000, Updated: 2000}, "created", "updated"); err != nil {
		t.Fatalf("PostEdit: %v", err)
	}
	mid := addPost(t, `{"title":"中","path":"mid","status":2,"markdown":"x"}`)
	if err := repo.PostEdit(ctx, &model.Post{Id: mid, Created: 1500, Updated: 1500}, "created", "updated"); err != nil {
		t.Fatalf("PostEdit: %v", err)
	}

	data := callOk(t, PostNeighbors, http.MethodGet, "/api/posts/"+mid+"/neighbors", "", map[string]string{"id": mid})
	naver, ok := data.(map[string]any)
	if !ok {
		t.Fatalf("结构异常: %v", data)
	}
	prev, _ := naver["prev"].(map[string]any)
	next, _ := naver["next"].(map[string]any)
	if prev == nil || !strings.Contains(prev["path"].(string), "newer") {
		t.Errorf("上一篇应为更新的内容, got %v", prev)
	}
	if next == nil || !strings.Contains(next["path"].(string), "older") {
		t.Errorf("下一篇应为更旧的内容, got %v", next)
	}
}

// TestPostTreeRoute 树形接口按 Pid 组装
func TestPostTreeRoute(t *testing.T) {
	setupKernelDB(t)
	root := addPost(t, `{"title":"根","path":"root","kind":3}`)
	child := addPost(t, `{"title":"子","path":"child","kind":3,"pid":"`+root+`"}`)

	items, count := callList(t, PostTree, http.MethodGet, "/api/posts/tree?page=0&size=0&kind=3")
	if count != 2 {
		t.Errorf("count = %d, want 2", count)
	}
	if len(items) != 1 {
		t.Fatalf("顶层数量 = %d, want 1", len(items))
	}
	top, ok := items[0].(map[string]any)
	if !ok || top["id"] != root {
		t.Fatalf("顶层节点 = %v", items[0])
	}
	children, ok := top["children"].([]any)
	if !ok || len(children) != 1 {
		t.Fatalf("子级 = %v", top["children"])
	}
	if c, ok := children[0].(map[string]any); !ok || c["id"] != child {
		t.Errorf("子级节点 = %v", children[0])
	}
}
