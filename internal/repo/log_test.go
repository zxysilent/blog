package repo

import (
	"context"
	"testing"
	"time"

	"blog/internal/model"
)

// mustLogAdd 同步写入一条日志
func mustLogAdd(t *testing.T, ctx context.Context, mod *model.Log) string {
	t.Helper()
	if mod.Created == 0 {
		mod.Created = time.Now().UnixMilli()
	}
	mod.BeforeInsert()
	if err := LogAdd(ctx, mod); err != nil {
		t.Fatalf("LogAdd: %v", err)
	}
	return mod.Id
}

// TestLogAddAndGet 写入后能读回全部字段
func TestLogAddAndGet(t *testing.T) {
	ctx := setupCateDB(t)
	mod := &model.Log{
		UserId: model.RoleRootId,
		Module: "用户管理",
		Action: "用户删除",
		Perm:   "user:delete",
		Method: "DELETE",
		Path:   "/api/users/xxx",
		Ip:     "127.0.0.1",
		Status: 200,
		Data:   `{"id":"xxx"}`,
	}
	mustLogAdd(t, ctx, mod)

	got, err := LogGet(ctx, mod.Id)
	if err != nil {
		t.Fatalf("LogGet: %v", err)
	}
	if got.Module != "用户管理" || got.Action != "用户删除" || got.Perm != "user:delete" {
		t.Errorf("字段不一致: %+v", got)
	}
	if got.Method != "DELETE" || got.Ip != "127.0.0.1" || got.Status != 200 {
		t.Errorf("字段不一致: %+v", got)
	}
	if got.Data != `{"id":"xxx"}` {
		t.Errorf("Data = %q", got.Data)
	}
	if got.Created <= 0 {
		t.Errorf("Created = %d", got.Created)
	}
}

// TestLogListFilters 各种过滤条件
func TestLogListFilters(t *testing.T) {
	ctx := setupCateDB(t)
	base := time.Now().UnixMilli()
	mustLogAdd(t, ctx, &model.Log{Perm: "user:delete", Module: "用户管理", Method: "DELETE", Status: 200, Created: base})
	mustLogAdd(t, ctx, &model.Log{Perm: "user:create", Module: "用户管理", Method: "POST", Status: 300, Created: base + 10})
	mustLogAdd(t, ctx, &model.Log{Perm: "post:delete", Module: "内容管理", Method: "DELETE", Status: 200, Created: base + 20})

	all, count, err := LogList(ctx, &model.LogListArgs{})
	if err != nil {
		t.Fatalf("LogList: %v", err)
	}
	if count != 3 || len(all) != 3 {
		t.Fatalf("全部 = %d, count = %d, want 3", len(all), count)
	}
	// 默认按时间倒序
	if all[0].Path == all[2].Path && all[0].Created < all[2].Created {
		t.Error("应按时间倒序")
	}

	mods, n, err := LogList(ctx, &model.LogListArgs{LogFilter: model.LogFilter{Module: "用户管理"}})
	if err != nil {
		t.Fatalf("LogList: %v", err)
	}
	if n != 2 || len(mods) != 2 {
		t.Errorf("按模块过滤 = %d, count = %d, want 2", len(mods), n)
	}

	mods, n, err = LogList(ctx, &model.LogListArgs{LogFilter: model.LogFilter{Perm: "post:delete"}})
	if err != nil {
		t.Fatalf("LogList: %v", err)
	}
	if n != 1 || len(mods) != 1 {
		t.Errorf("按权限过滤 = %d, count = %d, want 1", len(mods), n)
	}

	mods, n, err = LogList(ctx, &model.LogListArgs{LogFilter: model.LogFilter{Method: "DELETE"}})
	if err != nil {
		t.Fatalf("LogList: %v", err)
	}
	if n != 2 || len(mods) != 2 {
		t.Errorf("按方法过滤 = %d, count = %d, want 2", len(mods), n)
	}

	status := 300
	mods, n, err = LogList(ctx, &model.LogListArgs{LogFilter: model.LogFilter{Status: &status}})
	if err != nil {
		t.Fatalf("LogList: %v", err)
	}
	if n != 1 || len(mods) != 1 || mods[0].Status != 300 {
		t.Errorf("按状态过滤 = %+v, count = %d", mods, n)
	}

	// 时间区间
	mods, n, err = LogList(ctx, &model.LogListArgs{LogFilter: model.LogFilter{From: base + 5, To: base + 15}})
	if err != nil {
		t.Fatalf("LogList: %v", err)
	}
	if n != 1 || len(mods) != 1 {
		t.Errorf("按时间区间过滤 = %d, count = %d, want 1", len(mods), n)
	}

	// 关键字命中动作或路径
	mods, _, err = LogList(ctx, &model.LogListArgs{List: model.List{Mult: "用户"}})
	if err != nil {
		t.Fatalf("LogList: %v", err)
	}
	if len(mods) != 2 {
		t.Errorf("关键字过滤 = %d, want 2", len(mods))
	}
}

// TestLogListPaging 分页
func TestLogListPaging(t *testing.T) {
	ctx := setupCateDB(t)
	for i := 0; i < 5; i++ {
		mustLogAdd(t, ctx, &model.Log{Perm: "user:delete", Module: "用户管理", Status: 200})
	}
	mods, count, err := LogList(ctx, &model.LogListArgs{List: model.List{Pi: 1, Ps: 2}})
	if err != nil {
		t.Fatalf("LogList: %v", err)
	}
	if count != 5 || len(mods) != 2 {
		t.Errorf("分页 = %d, count = %d, want 2/5", len(mods), count)
	}
}

// TestLogFillUser 列表与详情都应带出操作人
func TestLogFillUser(t *testing.T) {
	ctx := setupCateDB(t)
	uid := addUser(t, ctx, "operator", "")
	mustLogAdd(t, ctx, &model.Log{UserId: uid, Perm: "user:delete", Module: "用户管理", Status: 200})

	mods, _, err := LogList(ctx, &model.LogListArgs{})
	if err != nil {
		t.Fatalf("LogList: %v", err)
	}
	if len(mods) != 1 || mods[0].User == nil || mods[0].User.Num != "operator" {
		t.Fatalf("未填充操作人: %+v", mods)
	}

	got, err := LogGet(ctx, mods[0].Id)
	if err != nil {
		t.Fatalf("LogGet: %v", err)
	}
	if got.User == nil || got.User.Num != "operator" {
		t.Errorf("详情未填充操作人: %+v", got.User)
	}
	// 操作人被删除后不应报错，只是不带出
	if err := UserDrop(ctx, uid); err != nil {
		t.Fatalf("UserDrop: %v", err)
	}
	mods, _, err = LogList(ctx, &model.LogListArgs{})
	if err != nil {
		t.Fatalf("LogList: %v", err)
	}
	if len(mods) != 1 || mods[0].User != nil {
		t.Errorf("操作人已删除，应留空: %+v", mods[0].User)
	}
}

// TestLogDrop 删除单条
func TestLogDrop(t *testing.T) {
	ctx := setupCateDB(t)
	id := mustLogAdd(t, ctx, &model.Log{Perm: "user:delete", Module: "用户管理", Status: 200})
	if err := LogDrop(ctx, id); err != nil {
		t.Fatalf("LogDrop: %v", err)
	}
	if _, err := LogGet(ctx, id); err == nil {
		t.Error("删除后仍能查到")
	}
}

// TestLogClear 按时间清理与全量清空
func TestLogClear(t *testing.T) {
	ctx := setupCateDB(t)
	base := time.Now().UnixMilli()
	mustLogAdd(t, ctx, &model.Log{Perm: "a", Created: base - 10000})
	mustLogAdd(t, ctx, &model.Log{Perm: "b", Created: base - 5000})
	mustLogAdd(t, ctx, &model.Log{Perm: "c", Created: base})

	// 只删 8 秒前的
	n, err := LogClear(ctx, base-8000)
	if err != nil {
		t.Fatalf("LogClear: %v", err)
	}
	if n != 1 {
		t.Errorf("删除条数 = %d, want 1", n)
	}
	mods, count, err := LogList(ctx, &model.LogListArgs{})
	if err != nil {
		t.Fatalf("LogList: %v", err)
	}
	if count != 2 || len(mods) != 2 {
		t.Errorf("剩余 = %d, count = %d, want 2", len(mods), count)
	}

	// 全部清空
	n, err = LogClear(ctx, 0)
	if err != nil {
		t.Fatalf("LogClear: %v", err)
	}
	if n != 2 {
		t.Errorf("清空条数 = %d, want 2", n)
	}
	mods, count, err = LogList(ctx, &model.LogListArgs{})
	if err != nil {
		t.Fatalf("LogList: %v", err)
	}
	if count != 0 || len(mods) != 0 {
		t.Errorf("清空后 = %d, count = %d", len(mods), count)
	}
}

// TestLogPushAsync 异步写入最终一致
func TestLogPushAsync(t *testing.T) {
	ctx := setupCateDB(t)
	mod := &model.Log{Perm: "user:delete", Module: "用户管理", Status: 200, Created: time.Now().UnixMilli()}
	mod.BeforeInsert()
	LogPush(mod)
	LogFlush(2 * time.Second)

	got, err := LogGet(ctx, mod.Id)
	if err != nil {
		t.Fatalf("异步写入未落库: %v", err)
	}
	if got.Perm != "user:delete" {
		t.Errorf("Perm = %q", got.Perm)
	}
}

// TestLogPushNil 入队 nil 不应 panic
func TestLogPushNil(t *testing.T) {
	LogPush(nil)
	LogFlush(time.Second)
}
