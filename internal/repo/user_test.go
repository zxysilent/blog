package repo

import (
	"context"
	"testing"

	"blog/internal/model"
)

// addUser 新增用户，返回主键
func addUser(t *testing.T, ctx context.Context, num, openid string) string {
	t.Helper()
	mod := &model.User{
		Num:    num,
		Name:   num,
		Openid: openid,
		RoleId: model.RoleRootId,
		Status: model.UserStatusNormal,
	}
	mod.BeforeInsert()
	if err := UserAdd(ctx, mod); err != nil {
		t.Fatalf("UserAdd(%s): %v", num, err)
	}
	return mod.Id
}

// TestUserOpenidLookup 按 openid 反查用户，空串不命中
func TestUserOpenidLookup(t *testing.T) {
	ctx := setupCateDB(t)
	id := addUser(t, ctx, "wxuser", "oABC123def")

	mod, err := UserOpenid(ctx, "oABC123def")
	if err != nil {
		t.Fatalf("UserOpenid: %v", err)
	}
	if mod.Id != id {
		t.Errorf("id = %q, want %q", mod.Id, id)
	}
	// 前后空格应被规范化
	if mod, err = UserOpenid(ctx, "  oABC123def  "); err != nil {
		t.Errorf("带空格查询失败: %v", err)
	} else if mod.Id != id {
		t.Errorf("带空格命中错误用户: %q", mod.Id)
	}
	// 未绑定 / 空串直接当作不存在
	if _, err := UserOpenid(ctx, ""); err == nil {
		t.Error("空 openid 应返回不存在")
	}
	if _, err := UserOpenid(ctx, "not-exist"); err == nil {
		t.Error("不存在的 openid 应返回不存在")
	}
}

// TestUserExistOpenid 唯一性校验：空串不占用，编辑时排除自身
func TestUserExistOpenid(t *testing.T) {
	ctx := setupCateDB(t)
	bound := addUser(t, ctx, "bound", "oBOUND1")
	addUser(t, ctx, "free1", "")
	addUser(t, ctx, "free2", "")

	if !UserExistOpenid(ctx, "oBOUND1", "") {
		t.Error("已绑定的 openid 应判为占用")
	}
	if UserExistOpenid(ctx, "oBOUND1", bound) {
		t.Error("编辑自身时应排除自身")
	}
	if UserExistOpenid(ctx, "oFREE1", "") {
		t.Error("未使用的 openid 不应判为占用")
	}
	// 多个用户未绑定时互不影响（空串不参与唯一性）
	if UserExistOpenid(ctx, "", "") {
		t.Error("空 openid 不应判为占用")
	}
}

// TestUserOpenidUniqueAtService 同一 openid 不能绑到两个账号
func TestUserOpenidMultipleEmpty(t *testing.T) {
	ctx := setupCateDB(t)
	// 三个未绑定用户共存，说明索引不是 UNIQUE
	addUser(t, ctx, "u1", "")
	addUser(t, ctx, "u2", "")
	addUser(t, ctx, "u3", "")

	// all=true 表示不做层级限制，测试只关心 openid 唯一性
	mods, count, err := UserList(ctx, &model.UserListArgs{}, 0, "", true)
	if err != nil {
		t.Fatalf("UserList: %v", err)
	}
	// 3 个新建 + 种子根账号
	if count != 4 || len(mods) != 4 {
		t.Errorf("用户数 = %d, count = %d, want 4", len(mods), count)
	}
	empty := 0
	for _, it := range mods {
		if it.Openid == "" {
			empty++
		}
	}
	if empty != 4 {
		t.Errorf("未绑定数 = %d, want 4（空 openid 不参与唯一性）", empty)
	}
}

// TestUserListFilterOpenid 列表按 openid 精确过滤
func TestUserListFilterOpenid(t *testing.T) {
	ctx := setupCateDB(t)
	addUser(t, ctx, "wx1", "oWX1")
	addUser(t, ctx, "wx2", "oWX2")
	addUser(t, ctx, "plain", "")

	mods, count, err := UserList(ctx, &model.UserListArgs{Openid: "oWX2"}, 0, "", true)
	if err != nil {
		t.Fatalf("UserList: %v", err)
	}
	if count != 1 || len(mods) != 1 {
		t.Fatalf("按 openid 过滤 = %d, count = %d, want 1", len(mods), count)
	}
	if mods[0].Num != "wx2" {
		t.Errorf("命中用户 = %q, want wx2", mods[0].Num)
	}

	// 带空格同样能命中
	mods, count, err = UserList(ctx, &model.UserListArgs{Openid: "  oWX1 "}, 0, "", true)
	if err != nil {
		t.Fatalf("UserList: %v", err)
	}
	if count != 1 || len(mods) != 1 || mods[0].Num != "wx1" {
		t.Errorf("带空格过滤结果 = %+v, count = %d", mods, count)
	}
}

// TestUserOpenidPersist 写入后能完整读回，长 openid 不截断
func TestUserOpenidPersist(t *testing.T) {
	ctx := setupCateDB(t)
	want := "oABCdefGHIjklMNOpqrSTUvwxYZ0123456789_-ABC"
	id := addUser(t, ctx, "long", want)

	mod, err := UserGet(ctx, id)
	if err != nil {
		t.Fatalf("UserGet: %v", err)
	}
	if mod.Openid != want {
		t.Errorf("Openid = %q, want %q", mod.Openid, want)
	}
}

// TestUserOpenidUnbind 清空 openid 后可重新绑定给别的账号
func TestUserOpenidUnbind(t *testing.T) {
	ctx := setupCateDB(t)
	a := addUser(t, ctx, "a", "oSHARED")
	b := addUser(t, ctx, "b", "")

	// 未解绑前冲突
	if !UserExistOpenid(ctx, "oSHARED", b) {
		t.Error("应检测到冲突")
	}
	// 解绑 A
	if err := UserEdit(ctx, &model.User{Id: a, Openid: ""}, "openid"); err != nil {
		t.Fatalf("UserEdit: %v", err)
	}
	if UserExistOpenid(ctx, "oSHARED", b) {
		t.Error("解绑后不应再冲突")
	}
	if _, err := UserOpenid(ctx, "oSHARED"); err == nil {
		t.Error("解绑后不应再查到")
	}
	// 绑定到 B
	if err := UserEdit(ctx, &model.User{Id: b, Openid: "oSHARED"}, "openid"); err != nil {
		t.Fatalf("UserEdit: %v", err)
	}
	mod, err := UserOpenid(ctx, "oSHARED")
	if err != nil {
		t.Fatalf("UserOpenid: %v", err)
	}
	if mod.Id != b {
		t.Errorf("命中 id = %q, want %q", mod.Id, b)
	}
}
