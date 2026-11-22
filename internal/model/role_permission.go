package model

import "blog/pkg/rndm"

// RolePerm 角色-权限关联
// 表名保持 role_permission（对外的授权关系语义），
// 其余内部标识统一使用 perm 缩写。
type RolePerm struct {
	Id     string `xorm:"VARCHAR(64) PK" json:"id"`                             //主键
	RoleId string `xorm:"UNIQUE(role_perm) VARCHAR(64)" json:"role_id"`         //角色id
	PermId string `xorm:"perm_id UNIQUE(role_perm) VARCHAR(64)" json:"perm_id"` //权限id
}

func (RolePerm) TableName() string {
	return "role_permission"
}

// BeforeInsert 插入前生成主键
func (s *RolePerm) BeforeInsert() {
	if s.Id == "" {
		s.Id = rndm.ULID()
	}
}

// IsPkZero 主键是否为空
func (s RolePerm) IsPkZero() bool {
	return s.Id == ""
}

// GetId 返回主键
func (s RolePerm) GetId() string {
	return s.Id
}
