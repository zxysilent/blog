package model

// RoleGrant 角色授权
type RoleGrant struct {
	Id      int `xorm:"INT(11) PK AUTOINCR comment('主键')" json:"id"`
	RoleId  int `xorm:"INT(11) DEFAULT 0 comment('角色id')" json:"role_id"`
	GrantId int `xorm:"INT(11) DEFAULT 0 comment('授权id')" json:"grant_id"`
}

func (RoleGrant) TableName() string {
	return "sys_role_grant"
}

// RoleGrantAddMulti 批量添加角色授权信息
func RoleGrantAddMulti(mods []RoleGrant) error {
	sess := Db.NewSession()
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
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.In("grant_id", grantIds).Delete(&RoleGrant{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// RoleGrantEdit 编辑角色授权
// param	mods	"待添加集合"
// param	grantIds	"待删除集合"
func RoleGrantEdit(mods []RoleGrant, grantIds []int) error {
	if len(mods) == 0 && len(grantIds) == 0 {
		return nil
	}
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if len(mods) > 0 {
		if _, err := sess.InsertMulti(&mods); err != nil {
			sess.Rollback()
			return err
		}
	}
	if len(grantIds) > 0 {
		if _, err := sess.In("grant_id", grantIds).Delete(&RoleGrant{}); err != nil {
			sess.Rollback()
			return err
		}
	}
	sess.Commit()
	return nil
}

func getGrantsByRole(roleId int) []int {
	mods := make([]int, 0, 8)
	Db.Table("sys_role_grant").Cols("grant_id").Where("role_id = ?", roleId).Find(&mods)
	return mods
}
