package model

// RoleGrant 角色授权
type RoleGrant struct {
	Id      int `xorm:"INT PK AUTOINCR" json:"id"`     // 主键
	RoleId  int `xorm:"INT DEFAULT 0" json:"role_id"`  // 角色id
	GrantId int `xorm:"INT DEFAULT 0" json:"grant_id"` // 授权id
}

func (RoleGrant) TableName() string {
	return sysPrefix + "role_grant"
}
