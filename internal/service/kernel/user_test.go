package kernel

import (
	"context"
	"net/http"
	"testing"

	"blog/internal/model"
	"blog/internal/repo"
	"blog/pkg/reply"

	"github.com/labstack/echo/v4"
)

// userOf 通过详情接口取用户
func userOf(t *testing.T, id string) map[string]any {
	t.Helper()
	data := callAsRoot(t, UserGet, http.MethodGet, "/api/users/"+id, "", map[string]string{"id": id}).Data
	mod, ok := data.(map[string]any)
	if !ok {
		t.Fatalf("详情结构异常: %v", data)
	}
	return mod
}

// rootId 种子根账号主键
func rootId(t *testing.T) string {
	t.Helper()
	mod, err := repo.UserNum(context.Background(), "root")
	if err != nil {
		t.Fatalf("取根账号失败: %v", err)
	}
	return mod.Id
}

// callAsRoot 以根账号身份调用接口
func callAsRoot(t *testing.T, handler echo.HandlerFunc, method, target, body string, params map[string]string) reply.Reply {
	t.Helper()
	return callAs(t, rootId(t), handler, method, target, body, params)
}

// addUser 以根账号新增用户，返回主键
func addUser(t *testing.T, body string) string {
	t.Helper()
	out := callAsRoot(t, UserAdd, http.MethodPost, "/api/users", body, nil)
	if out.Code != reply.CodeSucc {
		t.Fatalf("新增用户失败: code=%d msg=%s", out.Code, out.Msg)
	}
	id, _ := out.Data.(string)
	if id == "" {
		t.Fatal("新增用户未返回主键")
	}
	return id
}

// TestUserAddIgnoresType 已移除用户类型，客户端多传 type 应被忽略
func TestUserAddIgnoresType(t *testing.T) {
	setupKernelDB(t)
	id := addUser(t, `{"num":"user1","passwd":"123456","role_id":"`+model.RoleRootId+`","type":1}`)

	mod := userOf(t, id)
	if _, has := mod["type"]; has {
		t.Error("响应不应再包含 type 字段")
	}
	if mod["status"] != float64(model.UserStatusNormal) {
		t.Errorf("默认状态 = %v, want %d", mod["status"], model.UserStatusNormal)
	}
	if mod["num"] != "user1" {
		t.Errorf("num = %v", mod["num"])
	}
}

// TestUserEditIgnoresType 修改时多传 type 同样被忽略
func TestUserEditIgnoresType(t *testing.T) {
	setupKernelDB(t)
	id := addUser(t, `{"num":"user2","passwd":"123456","role_id":"`+model.RoleRootId+`"}`)

	callAsRoot(t, UserEdit, http.MethodPut, "/api/users/"+id,
		`{"type":1,"name":"改名"}`, map[string]string{"id": id})
	mod := userOf(t, id)
	if _, has := mod["type"]; has {
		t.Error("响应不应再包含 type 字段")
	}
	if mod["name"] != "改名" {
		t.Errorf("name = %v", mod["name"])
	}
}

// TestUserWithoutRoleHasNoPerm 未分配角色的账号什么权限都没有
// 这正是移除用户类型后的授权模型：能不能做事只看角色授权。
func TestUserWithoutRoleHasNoPerm(t *testing.T) {
	setupKernelDB(t)
	ctx := context.Background()
	id := addUser(t, `{"num":"norole","passwd":"123456","role_id":""}`)

	mod, err := repo.UserGet(ctx, id)
	if err != nil {
		t.Fatalf("UserGet: %v", err)
	}
	if mod.RoleId != "" {
		t.Fatalf("预期无角色, got %q", mod.RoleId)
	}
	// 任何带权限标识的接口都会被拒
	for _, perm := range []string{"user:list", "post:create", "role:update", "status:read"} {
		if repo.RolePermAllow(ctx, mod.RoleId, perm) {
			t.Errorf("无角色账号不应拥有 %q", perm)
		}
	}
	// 不需要权限的接口照常放行
	if !repo.RolePermAllow(ctx, mod.RoleId, "") {
		t.Error("空 perm 应放行")
	}
}

// TestUserOpenidUnique 同一 openid 不能绑两个账号
func TestUserOpenidUnique(t *testing.T) {
	setupKernelDB(t)
	first := addUser(t, `{"num":"wx1","passwd":"123456","role_id":"`+model.RoleRootId+`","openid":"oABC123"}`)
	if first == "" {
		t.Fatal("新增未返回主键")
	}
	// 重复绑定
	out := callAsRoot(t, UserAdd, http.MethodPost, "/api/users",
		`{"num":"wx2","passwd":"123456","role_id":"`+model.RoleRootId+`","openid":"oABC123"}`, nil)
	if out.Code == reply.CodeSucc {
		t.Error("重复 openid 应被拒绝")
	} else if out.Msg != "openid 已被其它账号绑定" {
		t.Errorf("错误信息 = %q", out.Msg)
	}
	// 非法格式
	out = callAsRoot(t, UserAdd, http.MethodPost, "/api/users",
		`{"num":"wx3","passwd":"123456","role_id":"`+model.RoleRootId+`","openid":"bad openid"}`, nil)
	if out.Code == reply.CodeSucc {
		t.Error("非法 openid 应被拒绝")
	}
	// 解绑后可以重新绑定
	callAsRoot(t, UserEdit, http.MethodPut, "/api/users/"+first,
		`{"openid":""}`, map[string]string{"id": first})
	addUser(t, `{"num":"wx2","passwd":"123456","role_id":"`+model.RoleRootId+`","openid":"oABC123"}`)
}

// TestUserListFilterOpenid 列表按 openid 精确过滤
func TestUserListFilterOpenid(t *testing.T) {
	setupKernelDB(t)
	addUser(t, `{"num":"wxa","passwd":"123456","role_id":"`+model.RoleRootId+`","openid":"oWXA"}`)
	addUser(t, `{"num":"wxb","passwd":"123456","role_id":"`+model.RoleRootId+`","openid":"oWXB"}`)

	items, count := callListAs(t, rootId(t), UserList, http.MethodGet, "/api/users?page=0&size=0&openid=oWXB")
	if count != 1 || len(items) != 1 {
		t.Fatalf("过滤结果 = %d, count = %d, want 1", len(items), count)
	}
	mod, _ := items[0].(map[string]any)
	if mod["num"] != "wxb" {
		t.Errorf("命中用户 = %v", mod["num"])
	}
}
