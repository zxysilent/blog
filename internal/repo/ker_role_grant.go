package repo

import "blog/internal/model"

// RoleGrantAddMulti 批量添加角色授权信息
func RoleGrantAddMulti(mods []model.RoleGrant) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.InsertMulti(&mods); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// RoleGrantDropMulti 批量删除角色授权信息
func RoleGrantDropMulti(grantIds []int) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.In("grant_id", grantIds).Delete(&model.RoleGrant{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// RoleGrantEdit 编辑角色授权
// param	mods	"待添加集合"
// param	grantIds	"待删除集合"
func RoleGrantEdit(mods []model.RoleGrant, grantIds []int) error {
	if len(mods) == 0 && len(grantIds) == 0 {
		return nil
	}
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	if len(mods) > 0 {
		if _, err := sess.InsertMulti(&mods); err != nil {
			sess.Rollback()
			return err
		}
	}
	if len(grantIds) > 0 {
		if _, err := sess.In("grant_id", grantIds).Delete(&model.RoleGrant{}); err != nil {
			sess.Rollback()
			return err
		}
	}
	sess.Commit()
	return nil
}

func getGrantsByRole(roleId int) []int {
	mods := make([]int, 0, 8)
	db.Table("role_grant").Cols("grant_id").Where("role_id = ?", roleId).Find(&mods)
	return mods
}
