package repo

import (
	"context"
	"strings"
	"testing"

	"blog/internal/model"
)

// addPost 新增一条内容
func addPost(t *testing.T, ctx context.Context, mod model.Post) string {
	t.Helper()
	if mod.Kind == 0 {
		mod.Kind = model.KindPost
	}
	if mod.Status == 0 {
		mod.Status = model.StatusDraft
	}
	if mod.Title == "" {
		mod.Title = mod.Path
	}
	mod.Show = true
	if err := PostAdd(ctx, &mod); err != nil {
		t.Fatalf("PostAdd(%s): %v", mod.Path, err)
	}
	return mod.Id
}

// TestPostGetPath 公开读取只允许「已发布 + 可见」
func TestPostGetPath(t *testing.T) {
	ctx := setupCateDB(t)

	draft := addPost(t, ctx, model.Post{Path: "draft", Status: model.StatusDraft})
	published := addPost(t, ctx, model.Post{Path: "pub", Status: model.StatusFinish})
	hidden := model.Post{Path: "hidden", Status: model.StatusFinish, Title: "hidden", Show: false}
	if err := PostAdd(ctx, &hidden); err != nil {
		t.Fatalf("PostAdd: %v", err)
	}

	if _, err := PostGetPath(ctx, "pub"); err != nil {
		t.Errorf("已发布内容应可读取: %v", err)
	}
	if _, err := PostGetPath(ctx, "draft"); err == nil {
		t.Errorf("草稿不应可读取 (id=%s)", draft)
	}
	if _, err := PostGetPath(ctx, "hidden"); err == nil {
		t.Errorf("隐藏内容不应可读取 (id=%s)", hidden.Id)
	}
	if _, err := PostGetPath(ctx, "nope"); err == nil {
		t.Error("不存在的内容不应可读取")
	}
	// 已发布内容应带出栏目信息
	mod, _ := PostGetPath(ctx, "pub")
	if mod == nil {
		t.Fatal("未取到内容")
	}
	if mod.Id != published {
		t.Errorf("id = %s, want %s", mod.Id, published)
	}
}

// TestPostGetPathFillsCate 详情从栏目缓存填充，栏目删除后不再带出
func TestPostGetPathFillsCate(t *testing.T) {
	ctx := setupCateDB(t)
	cate := mustCateAdd(t, ctx, model.Cate{Name: "技术", Path: "tech", Use: true, Show: true})
	addPost(t, ctx, model.Post{Path: "p1", Status: model.StatusFinish, CateId: cate.Id})

	mod, err := PostGetPath(ctx, "p1")
	if err != nil {
		t.Fatalf("PostGetPath: %v", err)
	}
	if mod.Cate == nil || mod.Cate.Name != "技术" {
		t.Fatalf("栏目未填充: %+v", mod.Cate)
	}
	// 缓存里的副本不能带子级，避免详情接口返回整棵子树
	if mod.Cate.Children != nil {
		t.Error("填充的栏目不应带子级")
	}
}

// TestPostHitsIncrement 浏览量自增并发安全
func TestPostHitsIncrement(t *testing.T) {
	ctx := setupCateDB(t)
	id := addPost(t, ctx, model.Post{Path: "hits"})

	const n = 20
	done := make(chan struct{}, n)
	for i := 0; i < n; i++ {
		go func() {
			_ = PostHits(ctx, id)
			done <- struct{}{}
		}()
	}
	for i := 0; i < n; i++ {
		<-done
	}
	mod, err := PostGet(ctx, id)
	if err != nil {
		t.Fatalf("PostGet: %v", err)
	}
	if mod.Hits != n {
		t.Errorf("hits = %d, want %d", mod.Hits, n)
	}
}

// TestPostExistPath 路径唯一性校验，编辑时排除自身
func TestPostExistPath(t *testing.T) {
	ctx := setupCateDB(t)
	id := addPost(t, ctx, model.Post{Path: "exists"})

	if !PostExistPath(ctx, "exists", "") {
		t.Error("已存在路径应判为占用")
	}
	if PostExistPath(ctx, "exists", id) {
		t.Error("编辑自身时应排除自身")
	}
	if PostExistPath(ctx, "free", "") {
		t.Error("不存在的路径不应判为占用")
	}
}

// TestPostListOmitsContent 管理端列表默认不返回大字段
func TestPostListOmitsContent(t *testing.T) {
	ctx := setupCateDB(t)
	addPost(t, ctx, model.Post{Path: "p1", Markdown: "# big", Richtext: "<p>big</p>"})

	mods, count, err := PostList(ctx, &model.PostListArgs{})
	if err != nil {
		t.Fatalf("PostList: %v", err)
	}
	if count != 1 || len(mods) != 1 {
		t.Fatalf("列表数量 = %d, count = %d", len(mods), count)
	}
	if mods[0].Markdown != "" || mods[0].Richtext != "" {
		t.Errorf("列表不应返回正文字段: md=%q rt=%q", mods[0].Markdown, mods[0].Richtext)
	}
	if mods[0].Title == "" {
		t.Error("列表应返回标题")
	}

	// 显式指定列时按需返回
	mods, _, err = PostList(ctx, &model.PostListArgs{}, "id", "title", "markdown")
	if err != nil {
		t.Fatalf("PostList: %v", err)
	}
	if len(mods) != 1 || mods[0].Markdown != "# big" {
		t.Errorf("显式列未生效: %+v", mods)
	}
}

// TestPostListOrder 置顶优先，其次按创建时间倒序
func TestPostListOrder(t *testing.T) {
	ctx := setupCateDB(t)
	old := model.Post{Path: "old", Title: "old", Created: 1000, Updated: 1000}
	if err := PostAdd(ctx, &old); err != nil {
		t.Fatalf("PostAdd: %v", err)
	}
	newer := model.Post{Path: "newer", Title: "newer", Created: 2000, Updated: 2000, Istop: true}
	if err := PostAdd(ctx, &newer); err != nil {
		t.Fatalf("PostAdd: %v", err)
	}

	mods, _, err := PostList(ctx, &model.PostListArgs{})
	if err != nil {
		t.Fatalf("PostList: %v", err)
	}
	if len(mods) != 2 {
		t.Fatalf("数量 = %d", len(mods))
	}
	// 置顶的在最前，即使创建时间更早
	if !mods[0].Istop || mods[0].Path != "newer" {
		t.Errorf("首条应为置顶内容, got %+v", mods[0])
	}
}

// TestPostCountFilter 计数与过滤条件保持一致
func TestPostCountFilter(t *testing.T) {
	ctx := setupCateDB(t)
	addPost(t, ctx, model.Post{Path: "a", Kind: model.KindPost, Status: model.StatusFinish})
	addPost(t, ctx, model.Post{Path: "b", Kind: model.KindPage, Status: model.StatusFinish})
	addPost(t, ctx, model.Post{Path: "c", Kind: model.KindPost, Status: model.StatusDraft})

	if n := PostCount(ctx, &model.PostFilter{}); n != 3 {
		t.Errorf("全部计数 = %d, want 3", n)
	}
	kind := model.KindPost
	if n := PostCount(ctx, &model.PostFilter{Kind: &kind}); n != 2 {
		t.Errorf("按类型计数 = %d, want 2", n)
	}
	status := model.StatusFinish
	if n := PostCount(ctx, &model.PostFilter{Status: &status}); n != 2 {
		t.Errorf("按状态计数 = %d, want 2", n)
	}
	show := true
	if n := PostCount(ctx, &model.PostFilter{Show: &show}); n != 3 {
		t.Errorf("按可见计数 = %d, want 3", n)
	}
}

// TestPostPageCateSub 栏目子树过滤
func TestPostPageCateSub(t *testing.T) {
	ctx := setupCateDB(t)
	top := mustCateAdd(t, ctx, model.Cate{Name: "技术", Path: "tech", Use: true, Show: true})
	sub := mustCateAdd(t, ctx, model.Cate{Name: "Go", Path: "go", Pid: top.Id, Use: true, Show: true})

	addPost(t, ctx, model.Post{Path: "p-top", Status: model.StatusFinish, Kind: model.KindPost, CateId: top.Id})
	addPost(t, ctx, model.Post{Path: "p-sub", Status: model.StatusFinish, Kind: model.KindPost, CateId: sub.Id})

	filter := &model.PostFilterPage{List: model.List{Ps: 10}, PostFilter: model.PostFilter{CateIds: CateSubIds(top.Id)}}
	mods, count, err := PostPage(filter)
	if err != nil {
		t.Fatalf("PostPage: %v", err)
	}
	if count != 2 || len(mods) != 2 {
		t.Errorf("子树过滤 = %d, count = %d, want 2", len(mods), count)
	}

	// 只查父栏目本身
	filter = &model.PostFilterPage{List: model.List{Ps: 10}, PostFilter: model.PostFilter{CateId: &top.Id}}
	mods, count, err = PostPage(filter)
	if err != nil {
		t.Fatalf("PostPage: %v", err)
	}
	if count != 1 || len(mods) != 1 || mods[0].Path != "p-top" {
		t.Errorf("单栏目过滤 = %+v, count = %d", mods, count)
	}
}

// TestPostTree 按 Pid 组装成树
func TestPostTree(t *testing.T) {
	ctx := setupCateDB(t)
	root := model.Post{Path: "root", Title: "root", Kind: model.KindNote, Status: model.StatusFinish}
	if err := PostAdd(ctx, &root); err != nil {
		t.Fatalf("PostAdd: %v", err)
	}
	child := model.Post{Path: "child", Title: "child", Kind: model.KindNote, Status: model.StatusFinish, Pid: root.Id}
	if err := PostAdd(ctx, &child); err != nil {
		t.Fatalf("PostAdd: %v", err)
	}
	grand := model.Post{Path: "grand", Title: "grand", Kind: model.KindNote, Status: model.StatusFinish, Pid: child.Id}
	if err := PostAdd(ctx, &grand); err != nil {
		t.Fatalf("PostAdd: %v", err)
	}

	kind := model.KindNote
	mods, count, err := PostTree(&model.PostFilterPage{PostFilter: model.PostFilter{Kind: &kind}})
	if err != nil {
		t.Fatalf("PostTree: %v", err)
	}
	if count != 3 {
		t.Errorf("count = %d, want 3", count)
	}
	if len(mods) != 1 {
		t.Fatalf("顶层数量 = %d, want 1", len(mods))
	}
	top := mods[0]
	if top.Level != 1 || top.Path != "root" {
		t.Errorf("顶层节点 = %+v", top)
	}
	if len(top.Children) != 1 || top.Children[0].Level != 2 {
		t.Fatalf("二级节点 = %+v", top.Children)
	}
	grandNode := top.Children[0].Children
	if len(grandNode) != 1 || grandNode[0].Level != 3 || grandNode[0].Path != "grand" {
		t.Errorf("三级节点 = %+v", grandNode)
	}
}

// TestPostTreeOrphanPromoted 父级不在结果集中的节点提升为顶层
func TestPostTreeOrphanPromoted(t *testing.T) {
	ctx := setupCateDB(t)
	// 只插入子节点，父节点不存在
	child := model.Post{Path: "orphan-child", Title: "child", Kind: model.KindNote, Status: model.StatusFinish, Pid: "not-exist"}
	if err := PostAdd(ctx, &child); err != nil {
		t.Fatalf("PostAdd: %v", err)
	}
	kind := model.KindNote
	mods, _, err := PostTree(&model.PostFilterPage{PostFilter: model.PostFilter{Kind: &kind}})
	if err != nil {
		t.Fatalf("PostTree: %v", err)
	}
	if len(mods) != 1 || mods[0].Id != child.Id {
		t.Errorf("孤儿节点未提升为顶层: %+v", mods)
	}
}

// TestPostDropCascadeTags 删除内容时同步清理标签关联
func TestPostDropCascadeTags(t *testing.T) {
	ctx := setupCateDB(t)
	tag := model.Tag{Name: "go"}
	if err := TagAdd(ctx, &tag); err != nil {
		t.Fatalf("TagAdd: %v", err)
	}
	id := addPost(t, ctx, model.Post{Path: "p1"})
	if err := PostTagSync(id, []string{tag.Id}); err != nil {
		t.Fatalf("PostTagSync: %v", err)
	}

	if err := PostDrop(ctx, id); err != nil {
		t.Fatalf("PostDrop: %v", err)
	}
	if _, err := PostGet(ctx, id); err == nil {
		t.Error("内容未删除")
	}
	if tags := PostTags(id); len(tags) != 0 {
		t.Errorf("标签关联未清理: %+v", tags)
	}
	// 标签本身不应被删除
	keep := &model.Tag{Id: tag.Id}
	if err := TagGet(ctx, keep); err != nil {
		t.Errorf("标签被误删: %v", err)
	}
}

// TestPostEditTop 置顶可反复切换
func TestPostEditTop(t *testing.T) {
	ctx := setupCateDB(t)
	id := addPost(t, ctx, model.Post{Path: "p1"})

	if err := PostEditTop(ctx, id, true); err != nil {
		t.Fatalf("PostEditTop: %v", err)
	}
	mod, _ := PostGet(ctx, id)
	if !mod.Istop {
		t.Error("置顶未生效")
	}
	if err := PostEditTop(ctx, id, false); err != nil {
		t.Fatalf("PostEditTop: %v", err)
	}
	mod, _ = PostGet(ctx, id)
	if mod.Istop {
		t.Error("取消置顶未生效")
	}
}

// TestPostNaverScope 上下篇只取「已发布 + 可见」且同栏目
func TestPostNaverScope(t *testing.T) {
	ctx := setupCateDB(t)
	cate := mustCateAdd(t, ctx, model.Cate{Name: "技术", Path: "tech", Use: true, Show: true})
	other := mustCateAdd(t, ctx, model.Cate{Name: "生活", Path: "life", Use: true, Show: true})

	base := int64(2000)
	seed := []struct {
		path   string
		cateId string
		status int
		show   bool
		at     int64
	}{
		// 比基准旧的同栏目内容 -> 下一篇
		{"older-same", cate.Id, model.StatusFinish, true, 1000},
		// 比基准新的同栏目内容 -> 上一篇
		{"newer-same", cate.Id, model.StatusFinish, true, 3000},
		// 其他栏目，不应参与
		{"older-other", other.Id, model.StatusFinish, true, 1500},
		// 草稿与隐藏，不应参与
		{"newer-draft", cate.Id, model.StatusDraft, true, 2500},
		{"newer-hidden", cate.Id, model.StatusFinish, false, 2600},
	}
	mid := model.Post{Path: "mid", Title: "mid", Kind: model.KindPost, Status: model.StatusFinish, Show: true, CateId: cate.Id, Created: base, Updated: base}
	if err := PostAdd(ctx, &mid); err != nil {
		t.Fatalf("PostAdd: %v", err)
	}
	for _, it := range seed {
		mod := model.Post{
			Path: it.path, Title: it.path, Kind: model.KindPost,
			Status: it.status, Show: it.show, CateId: it.cateId,
			Created: it.at, Updated: it.at,
		}
		if err := PostAdd(ctx, &mod); err != nil {
			t.Fatalf("PostAdd(%s): %v", it.path, err)
		}
	}

	naver := PostNaver(cate.Id, base)
	if naver.Prev == "" {
		t.Error("应有上一篇")
	}
	if naver.Next == "" {
		t.Error("应有下一篇")
	}
	// 草稿、隐藏、跨栏目的都不应被选中
	for _, bad := range []string{"older-other", "newer-draft", "newer-hidden"} {
		if strings.Contains(naver.Prev, bad) || strings.Contains(naver.Next, bad) {
			t.Errorf("上下篇命中了不应出现的内容 %q: prev=%q next=%q", bad, naver.Prev, naver.Next)
		}
	}
	// 上一篇 = 比当前新的同栏目内容，下一篇 = 比当前旧的
	if !strings.Contains(naver.Prev, "newer-same") {
		t.Errorf("上一篇应为更新的同栏目内容, got %q", naver.Prev)
	}
	if !strings.Contains(naver.Next, "older-same") {
		t.Errorf("下一篇应为更旧的同栏目内容, got %q", naver.Next)
	}
}
