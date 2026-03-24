package repo

import (
	"blog/internal/model"
	"context"
)

// AdminNum num 查询账号
func AdminNum(num string) (*model.Admin, bool) {
	mod := &model.Admin{Num: num}
	has, _ := db.Get(mod)
	return mod, has
}

// AdminExist 判断是否存在当前账号
func AdminExist(num string) bool {
	has, _ := db.Exist(&model.Admin{
		Num: num,
	})
	return has
}

// AdminPage 查询用户分页信息
func AdminPage(pi int, ps int, mult string, authRid int, aimRid int) ([]model.Admin, int, error) {
	mods := make([]model.Admin, 0, ps)
	sess := db.NewSession()
	defer sess.Close()
	if mult != "" {
		sess.Where("num like concat(?,'%')", mult)
	}
	sess.Where("role_id > ? ", authRid)
	if aimRid > 0 {
		sess.Where("role_id = ? ", aimRid)
	}
	count, err := sess.Desc("Created").Limit(ps, (pi-1)*ps).FindAndCount(&mods)
	// if len(mods) > 0 {
	// 	ids := make([]string, 0, len(mods))
	// 	for idx := range mods {
	// 		if !slices.Contains(ids, mods[idx].RoleId) {
	// 			ids = append(ids, mods[idx].RoleId)
	// 		}
	// 	}
	// 	mapSet := RoleIds(ids)
	// 	for idx := range mods {
	// 		mods[idx].Role = mapSet[mods[idx].RoleId]
	// 	}
	// }
	return mods, int(count), err
}

// AdminAll 查询用户所有
func AdminAll(mult string, authRid int, aimRid int) ([]model.Admin, error) {
	mods := make([]model.Admin, 0, 8)
	sess := db.NewSession()
	defer sess.Close()
	if mult != "" {
		sess.Where("name like concat(?,'%')", mult)
	}
	sess.Where("role_id > ? ", authRid)
	if aimRid > 0 {
		sess.Where("role_id = ? ", aimRid)
	}
	err := sess.Desc("id").Limit(128).Find(&mods)
	// if len(mods) > 0 {
	// 	ids := make([]string, 0, len(mods))
	// 	for idx := range mods {
	// 		if !slices.Contains(ids, mods[idx].RoleId) {
	// 			ids = append(ids, mods[idx].RoleId)
	// 		}
	// 	}
	// 	mapSet := RoleIds(ids)
	// 	for idx := range mods {
	// 		mods[idx].Role = mapSet[mods[idx].RoleId]
	// 	}
	// }
	return mods, err
}
func AdminEdit(ctx context.Context, mod *model.Admin, cols ...string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.Begin()
	if _, err := sess.ID(mod.Id).Cols(cols...).Update(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}
