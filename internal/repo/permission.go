package repo

import (
	"blog/internal/model"
	"context"
	"errors"
	"sort"
	"time"
)

// PermExist 权限是否已存在
func PermExist(ctx context.Context, mod *model.Perm) bool {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	exist, _ := sess.Exist(mod)
	return exist
}

// PermGet 权限详情
func PermGet(ctx context.Context, id string) (*model.Perm, error) {
	mod, err := GetAny[model.Perm](ctx, id)
	if err != nil {
		return nil, err
	}
	return &mod, nil
}

// PermGetKey 按权限标识查
func PermGetKey(ctx context.Context, perm string) (*model.Perm, error) {
	mod := &model.Perm{}
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	has, err := sess.Where("perm = ?", perm).Get(mod)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, ErrNotExist
	}
	return mod, nil
}

// PermList 权限分页列表
func PermList(ctx context.Context, args *model.PermListArgs, cols ...string) ([]model.Perm, int, error) {
	mods := make([]model.Perm, 0, args.GetCap())
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	if args.Mult != "" {
		sess.And("(perm like ? or name like ?)", args.SafeMult(), args.SafeMult())
	}
	// group 前端传的是分组键（如 content），库里 group 存中文名（如 内容管理），
	// 两者都接受，避免只按中文名过滤时前端分组卡片点击后为空
	if args.Group != "" {
		sess.And("(`group` = ? or group_key = ?)", args.Group, args.Group)
	}
	sess.Asc("`group`").Asc("sort").Asc("id")
	if args.IsZero() {
		err := sess.Find(&mods)
		return mods, len(mods), err
	}
	count, err := sess.Limit(args.GetLimit(), args.GetOffset()).FindAndCount(&mods)
	return mods, int(count), err
}

// PermAll 全部权限
func PermAll(ctx context.Context) ([]model.Perm, error) {
	mods := make([]model.Perm, 0, 32)
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	err := sess.Asc("`group`").Asc("sort").Asc("id").Find(&mods)
	return mods, err
}

// PermAdd 新增权限
func PermAdd(ctx context.Context, mod *model.Perm) error {
	return AddAny(ctx, mod)
}

// PermEdit 修改权限
func PermEdit(ctx context.Context, mod *model.Perm, cols ...string) error {
	return EditAny(ctx, mod, cols...)
}

// PermDrop 删除权限（同时清理角色授权）
func PermDrop(ctx context.Context, id string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if err := sess.Begin(); err != nil {
		return err
	}
	if _, err := sess.ID(id).NoAutoCondition().Delete(&model.Perm{}); err != nil {
		sess.Rollback()
		return err
	}
	if _, err := sess.Where("perm_id = ?", id).Delete(&model.RolePerm{}); err != nil {
		sess.Rollback()
		return err
	}
	return sess.Commit()
}

// PermSync 按路由汇总结果同步权限点
// 权限点以路由注册为唯一来源，完全只读：
//   - 路由中存在的 perm 不存在则新增，中文名/分组/排序有变化则更新
//   - 主键由 perm 稳定派生（model.PermId），同一 perm 在任何库、任何时刻
//     都得到同一个 id，因此角色授权（role_permission.perm_id）永不因
//     重建库、换环境、回滚版本而失效
//   - 库中存在但路由已下线的权限会被删除（连带清理角色关联）
func PermSync(ctx context.Context, defs []model.PermInitArgs) error {
	keep := make(map[string]struct{}, len(defs))
	for _, item := range defs {
		keep[item.Perm] = struct{}{}
		old, err := PermGetKey(ctx, item.Perm)
		if err != nil {
			if !errors.Is(err, ErrNotExist) {
				return err
			}
			mod := &model.Perm{
				Id:       model.PermId(item.Perm),
				Perm:     item.Perm,
				Name:     item.Name,
				GroupKey: item.GroupKey,
				Group:    item.Group,
				Sort:     item.Sort,
				Created:  time.Now().UnixMilli(),
			}
			if err := AddAny(ctx, mod); err != nil {
				return err
			}
			continue
		}
		// 主键漂移（旧版本随机主键）：校正为派生 id，并迁移角色授权引用
		if want := model.PermId(item.Perm); old.Id != want {
			if err := PermRealign(ctx, old.Id, want); err != nil {
				return err
			}
			old.Id = want
		}
		// 无变化时跳过，避免每次启动都写库
		if old.Name == item.Name && old.Group == item.Group &&
			old.GroupKey == item.GroupKey && old.Sort == item.Sort {
			continue
		}
		old.Name = item.Name
		old.GroupKey = item.GroupKey
		old.Group = item.Group
		old.Sort = item.Sort
		if err := EditAny(ctx, old, "name", "group_key", "group", "sort"); err != nil {
			return err
		}
	}
	// 清理已下线权限点
	all, err := PermAll(ctx)
	if err != nil {
		return err
	}
	for idx := range all {
		if _, ok := keep[all[idx].Perm]; ok {
			continue
		}
		if err := PermDrop(ctx, all[idx].Id); err != nil {
			return err
		}
	}
	// 顺带清理指向已不存在权限点的授权（仅旧版本升级时可能存在）
	return rolePermPrune(ctx)
}

// PermRealign 把权限点主键校正为按 perm 推导的稳定 id
// 同时把角色授权中的引用一并迁移，保证升级过程中已有授权不丢失
func PermRealign(ctx context.Context, from, to string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if err := sess.Begin(); err != nil {
		return err
	}
	if _, err := sess.Table("role_permission").Where("perm_id = ?", from).
		Update(map[string]any{"perm_id": to}); err != nil {
		sess.Rollback()
		return err
	}
	if _, err := sess.Table("permission").Where("id = ?", from).
		Update(map[string]any{"id": to}); err != nil {
		sess.Rollback()
		return err
	}
	return sess.Commit()
}

// PermTree 权限树，按分组聚合
// 返回两级结构：分组节点（Perm 为空）+ 权限叶子节点
func PermTree(ctx context.Context) ([]model.PermNode, error) {
	mods, err := PermAll(ctx)
	if err != nil {
		return nil, err
	}
	// group -> 子节点
	groups := make(map[string][]model.PermNode, 8)
	groupNames := make(map[string]string, 8)
	for idx := range mods {
		it := mods[idx]
		groups[it.GroupKey] = append(groups[it.GroupKey], model.PermNode{
			Id:       it.Id,
			Name:     it.Name,
			GroupKey: it.GroupKey,
			Group:    it.Group,
			Perm:     it.Perm,
			Sort:     it.Sort,
		})
		groupNames[it.GroupKey] = it.Group
	}
	nodes := make([]model.PermNode, 0, len(groups))
	for key, children := range groups {
		sort.Slice(children, func(i, j int) bool {
			if children[i].Sort != children[j].Sort {
				return children[i].Sort < children[j].Sort
			}
			return children[i].Perm < children[j].Perm
		})
		name := groupNames[key]
		nodes = append(nodes, model.PermNode{
			Name:     name,
			GroupKey: key,
			Group:    name,
			Sort:     children[0].Sort,
			Children: children,
		})
	}
	sort.Slice(nodes, func(i, j int) bool {
		if nodes[i].Sort != nodes[j].Sort {
			return nodes[i].Sort < nodes[j].Sort
		}
		return nodes[i].Name < nodes[j].Name
	})
	return nodes, nil
}

// PermMapIds 权限id集合查询
func PermMapIds(ctx context.Context, ids []string) map[string]*model.Perm {
	set := make(map[string]*model.Perm, len(ids))
	if len(ids) < 1 {
		return set
	}
	mods := make([]model.Perm, 0, len(ids))
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.In("id", ids).Find(&mods)
	for idx := range mods {
		set[mods[idx].Id] = &mods[idx]
	}
	return set
}
