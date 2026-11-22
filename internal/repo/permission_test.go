package repo

import (
	"testing"

	"blog/internal/model"
)

// TestPermTableNames 表名保持原样，不随内部缩写变化
func TestPermTableNames(t *testing.T) {
	if got := (model.Perm{}).TableName(); got != "permission" {
		t.Errorf("Perm.TableName() = %q, want permission", got)
	}
	if got := (model.RolePerm{}).TableName(); got != "role_permission" {
		t.Errorf("RolePerm.TableName() = %q, want role_permission", got)
	}
}

// TestRolePermColumn 列名与模型字段一致使用 perm_id
// 表名保持 role_permission，列名走 perm 缩写；
// 列名与模型不一致时 resetDerived 会重建该表，历史授权随之清空。
func TestRolePermColumn(t *testing.T) {
	setupCateDB(t)
	ti, err := db.TableInfo(&model.RolePerm{})
	if err != nil {
		t.Fatalf("TableInfo: %v", err)
	}
	cols := make(map[string]bool, 4)
	for _, col := range ti.Columns() {
		cols[col.Name] = true
	}
	if !cols["perm_id"] {
		t.Fatalf("role_permission 缺少 perm_id 列，实际列 = %v", cols)
	}
	if cols["permission_id"] {
		t.Error("不应出现 permission_id 列")
	}
	// 库里真实建出的表同样应是 perm_id
	metas, err := db.DBMetas()
	if err != nil {
		t.Fatalf("DBMetas: %v", err)
	}
	for _, meta := range metas {
		if meta.Name != "role_permission" {
			continue
		}
		found := false
		for _, col := range meta.Columns() {
			if col.Name == "perm_id" {
				found = true
			}
		}
		if !found {
			t.Error("库中 role_permission 缺少 perm_id 列")
		}
	}
}

// TestPermIdStable 主键由 perm 稳定派生：同 perm 恒等，长度固定
func TestPermIdStable(t *testing.T) {
	a := model.PermId("user:list")
	b := model.PermId("user:list")
	if a != b {
		t.Errorf("同 perm 派生出不同 id: %q vs %q", a, b)
	}
	if len(a) != 26 {
		t.Errorf("id 长度 = %d, want 26", len(a))
	}
	if c := model.PermId("user:create"); c == a {
		t.Error("不同 perm 不应得到相同 id")
	}
	// BeforeInsert 未指定 id 时使用派生值
	mod := &model.Perm{Perm: "tag:list"}
	mod.BeforeInsert()
	if mod.Id != model.PermId("tag:list") {
		t.Errorf("BeforeInsert id = %q, want %q", mod.Id, model.PermId("tag:list"))
	}
	// 已指定主键时不覆盖
	fixed := &model.Perm{Id: "fixed", Perm: "tag:list"}
	fixed.BeforeInsert()
	if fixed.Id != "fixed" {
		t.Errorf("指定主键被覆盖: %q", fixed.Id)
	}
}

// TestPermSyncUpsert 同步：新增、更新名称与排序、删除下线项
func TestPermSyncUpsert(t *testing.T) {
	ctx := setupCateDB(t)

	defs := []model.PermInitArgs{
		{Perm: "user:list", Name: "用户列表", GroupKey: "user", Group: "用户管理", Sort: 100},
		{Perm: "user:grant", Name: "用户管理授权", GroupKey: "user", Group: "用户管理", Sort: 100},
		{Perm: "post:list", Name: "内容列表", GroupKey: "content", Group: "内容管理", Sort: 600},
	}
	if err := PermSync(ctx, defs); err != nil {
		t.Fatalf("PermSync: %v", err)
	}
	all, err := PermAll(ctx)
	if err != nil {
		t.Fatalf("PermAll: %v", err)
	}
	if len(all) != 3 {
		t.Fatalf("权限点数量 = %d, want 3", len(all))
	}
	// 主键与派生值一致
	if got, err := PermGet(ctx, model.PermId("user:list")); err != nil {
		t.Fatalf("PermGet: %v", err)
	} else if got.Name != "用户列表" {
		t.Errorf("Name = %q", got.Name)
	}

	// 改名 + 改排序
	defs[0].Name = "用户列表2"
	defs[0].Sort = 90
	if err := PermSync(ctx, defs); err != nil {
		t.Fatalf("PermSync: %v", err)
	}
	got, _ := PermGetKey(ctx, "user:list")
	if got.Name != "用户列表2" || got.Sort != 90 {
		t.Errorf("更新未生效: name=%q sort=%d", got.Name, got.Sort)
	}

	// 下线内容分组
	if err := PermSync(ctx, defs[:2]); err != nil {
		t.Fatalf("PermSync: %v", err)
	}
	if _, err := PermGetKey(ctx, "post:list"); err == nil {
		t.Error("下线的权限点未删除")
	}
	if n := len(func() []model.Perm { m, _ := PermAll(ctx); return m }()); n != 2 {
		t.Errorf("剩余权限点 = %d, want 2", n)
	}
}

// TestPermTree 按分组键聚合两级树，分组在前叶子在后
func TestPermTree(t *testing.T) {
	ctx := setupCateDB(t)
	defs := []model.PermInitArgs{
		{Perm: "user:list", Name: "用户列表", GroupKey: "user", Group: "用户管理", Sort: 100},
		{Perm: "user:create", Name: "用户新增", GroupKey: "user", Group: "用户管理", Sort: 100},
		{Perm: "user:grant", Name: "用户管理授权", GroupKey: "user", Group: "用户管理", Sort: 100},
		{Perm: "post:list", Name: "内容列表", GroupKey: "content", Group: "内容管理", Sort: 600},
	}
	if err := PermSync(ctx, defs); err != nil {
		t.Fatalf("PermSync: %v", err)
	}
	nodes, err := PermTree(ctx)
	if err != nil {
		t.Fatalf("PermTree: %v", err)
	}
	if len(nodes) != 2 {
		t.Fatalf("分组数量 = %d, want 2", len(nodes))
	}
	// 排序小的分组在前
	if nodes[0].Group != "用户管理" || nodes[1].Group != "内容管理" {
		t.Errorf("分组顺序错误: %q, %q", nodes[0].Group, nodes[1].Group)
	}
	if nodes[0].GroupKey != "user" {
		t.Errorf("分组键 = %q, want user", nodes[0].GroupKey)
	}
	if len(nodes[0].Children) != 3 {
		t.Fatalf("用户管理叶子数量 = %d, want 3", len(nodes[0].Children))
	}
	// 叶子节点带 id 与 perm，分组节点两者都为空
	for _, leaf := range nodes[0].Children {
		if leaf.Id == "" || leaf.Perm == "" {
			t.Errorf("叶子缺少 id 或 perm: %+v", leaf)
		}
		if leaf.GroupKey != "user" {
			t.Errorf("叶子分组键 = %q", leaf.GroupKey)
		}
	}
	if nodes[0].Id != "" || nodes[0].Perm != "" {
		t.Errorf("分组节点不应带 id/perm: %+v", nodes[0])
	}
}

// TestPermListFilter 列表按分组与关键字过滤
func TestPermListFilter(t *testing.T) {
	ctx := setupCateDB(t)
	if err := PermSync(ctx, []model.PermInitArgs{
		{Perm: "user:list", Name: "用户列表", GroupKey: "user", Group: "用户管理", Sort: 100},
		{Perm: "user:create", Name: "用户新增", GroupKey: "user", Group: "用户管理", Sort: 100},
		{Perm: "post:list", Name: "内容列表", GroupKey: "content", Group: "内容管理", Sort: 600},
	}); err != nil {
		t.Fatalf("PermSync: %v", err)
	}

	mods, count, err := PermList(ctx, &model.PermListArgs{Group: "用户管理"})
	if err != nil {
		t.Fatalf("PermList: %v", err)
	}
	if count != 2 || len(mods) != 2 {
		t.Errorf("按分组过滤 = %d, count=%d, want 2", len(mods), count)
	}

	mods, _, err = PermList(ctx, &model.PermListArgs{List: model.List{Mult: "list"}})
	if err != nil {
		t.Fatalf("PermList: %v", err)
	}
	if len(mods) != 2 {
		t.Errorf("按关键字过滤 = %d, want 2", len(mods))
	}
}

// TestRolePermSyncAndAllow 授权写入后缓存立即生效
func TestRolePermSyncAndAllow(t *testing.T) {
	ctx := setupCateDB(t)
	if err := PermSync(ctx, []model.PermInitArgs{
		{Perm: "user:list", Name: "用户列表", GroupKey: "user", Group: "用户管理", Sort: 100},
		{Perm: "user:create", Name: "用户新增", GroupKey: "user", Group: "用户管理", Sort: 100},
	}); err != nil {
		t.Fatalf("PermSync: %v", err)
	}

	role := &model.Role{Name: "运营", Level: 1}
	role.BeforeInsert()
	if err := RoleAdd(ctx, role); err != nil {
		t.Fatalf("RoleAdd: %v", err)
	}

	// 初始无授权
	if RolePermAllow(ctx, role.Id, "user:list") {
		t.Error("未授权时不应放行")
	}
	if err := RolePermSync(ctx, role.Id, []string{model.PermId("user:list")}); err != nil {
		t.Fatalf("RolePermSync: %v", err)
	}
	if !RolePermAllow(ctx, role.Id, "user:list") {
		t.Error("授权后应放行")
	}
	if RolePermAllow(ctx, role.Id, "user:create") {
		t.Error("未授权的权限不应放行")
	}
	// 空 perm 表示接口无需授权
	if !RolePermAllow(ctx, role.Id, "") {
		t.Error("空 perm 应直接放行")
	}
	// 根角色直通
	if !RolePermAllow(ctx, model.RoleRootId, "any:perm") {
		t.Error("根角色应放行任意权限")
	}

	// 重复提交去重
	if err := RolePermSync(ctx, role.Id, []string{
		model.PermId("user:list"),
		model.PermId("user:list"),
		"",
	}); err != nil {
		t.Fatalf("RolePermSync: %v", err)
	}
	if ids := RolePermIds(ctx, role.Id); len(ids) != 1 {
		t.Errorf("去重后授权数 = %d, want 1", len(ids))
	}

	// 清空
	if err := RolePermSync(ctx, role.Id, nil); err != nil {
		t.Fatalf("RolePermSync: %v", err)
	}
	if RolePermAllow(ctx, role.Id, "user:list") {
		t.Error("清空后不应放行")
	}
}

// TestRolePermSyncScoped 只替换目标角色的授权，不影响其它角色
func TestRolePermSyncScoped(t *testing.T) {
	ctx := setupCateDB(t)
	if err := PermSync(ctx, []model.PermInitArgs{
		{Perm: "user:list", Name: "用户列表", GroupKey: "user", Group: "用户管理", Sort: 100},
		{Perm: "post:list", Name: "内容列表", GroupKey: "content", Group: "内容管理", Sort: 600},
	}); err != nil {
		t.Fatalf("PermSync: %v", err)
	}
	a := &model.Role{Name: "A", Level: 1}
	a.BeforeInsert()
	b := &model.Role{Name: "B", Level: 2}
	b.BeforeInsert()
	for _, it := range []*model.Role{a, b} {
		if err := RoleAdd(ctx, it); err != nil {
			t.Fatalf("RoleAdd: %v", err)
		}
	}
	if err := RolePermSync(ctx, a.Id, []string{model.PermId("user:list")}); err != nil {
		t.Fatalf("RolePermSync: %v", err)
	}
	if err := RolePermSync(ctx, b.Id, []string{model.PermId("post:list")}); err != nil {
		t.Fatalf("RolePermSync: %v", err)
	}
	// 再改 A 不应影响 B
	if err := RolePermSync(ctx, a.Id, []string{model.PermId("post:list")}); err != nil {
		t.Fatalf("RolePermSync: %v", err)
	}
	if !RolePermAllow(ctx, b.Id, "post:list") {
		t.Error("B 的授权被误删")
	}
	if RolePermAllow(ctx, a.Id, "user:list") {
		t.Error("A 的旧授权未清除")
	}
	if !RolePermAllow(ctx, a.Id, "post:list") {
		t.Error("A 的新授权未生效")
	}
}

// TestPermDropCascade 删除权限点时连带清理角色授权
func TestPermDropCascade(t *testing.T) {
	ctx := setupCateDB(t)
	if err := PermSync(ctx, []model.PermInitArgs{
		{Perm: "user:list", Name: "用户列表", GroupKey: "user", Group: "用户管理", Sort: 100},
	}); err != nil {
		t.Fatalf("PermSync: %v", err)
	}
	role := &model.Role{Name: "运营", Level: 1}
	role.BeforeInsert()
	if err := RoleAdd(ctx, role); err != nil {
		t.Fatalf("RoleAdd: %v", err)
	}
	permId := model.PermId("user:list")
	if err := RolePermSync(ctx, role.Id, []string{permId}); err != nil {
		t.Fatalf("RolePermSync: %v", err)
	}

	if err := PermDrop(ctx, permId); err != nil {
		t.Fatalf("PermDrop: %v", err)
	}
	if ids := RolePermIds(ctx, role.Id); len(ids) != 0 {
		t.Errorf("授权未随权限点清理: %v", ids)
	}
}
