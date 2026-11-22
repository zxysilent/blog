package model

import (
	"blog/pkg/rndm"
	"math"
)

// Role 角色
type Role struct {
	Id      string `xorm:"VARCHAR(64) PK" json:"id"`          //主键
	Name    string `xorm:"VARCHAR(255)" json:"name"`          //角色名称
	Intro   string `xorm:"VARCHAR(255)" json:"intro"`         //角色描述
	Level   int    `xorm:"INT(11) DEFAULT 1" json:"level"`    //角色层级-数值越小层级越高，0 仅根角色
	Inner   bool   `xorm:"TINYINT(4) DEFAULT 0" json:"inner"` //内置-禁止删除
	Updated int64  `xorm:"BIGINT" json:"updated"`             //修改时间
	Created int64  `xorm:"BIGINT" json:"created"`             //创建时间

	UserCount int `xorm:"-" json:"user_count"` //关联用户数，仅列表填充
	PermCount int `xorm:"-" json:"perm_count"` //已授权权限数，仅列表填充
}

func (Role) TableName() string {
	return "role"
}

// BeforeInsert 插入前生成主键
func (s *Role) BeforeInsert() {
	if s.Id == "" {
		s.Id = rndm.ULID()
	}
}

// IsPkZero 主键是否为空
func (s Role) IsPkZero() bool {
	return s.Id == ""
}

// GetId 返回主键
func (s Role) GetId() string {
	return s.Id
}

const (
	// RoleRootId 根角色唯一主键，拥有所有权限且无法删除
	// 26 位最小合法 ULID，字典序最小
	RoleRootId = "00000000000000000000000000"
	// RoleRootLevel 根角色层级，0 为最小值即最高权限
	// 非根角色层级必须严格大于其创建者层级，故 0 永远只属于根角色
	RoleRootLevel = 0
	// RoleNilLevel 角色缺失时的兜底层级，取最大值即权限最低
	// 层级数值越小权限越高，因此用最大值保证“无角色”既不拥有任何权限，
	// 也不会因为零值（=根角色层级）被误判为最高权限
	RoleNilLevel = math.MaxInt
)

// RoleListArgs 角色列表参数
type RoleListArgs struct {
	List
}

// RoleAddArgs 角色新增参数
type RoleAddArgs struct {
	Name  string `json:"name" form:"name"`   //角色名称
	Intro string `json:"intro" form:"intro"` //角色描述
	Level *int   `json:"level" form:"level"` //角色层级-不传则默认比当前角色大一级（权限更低）
}

// RoleEditArgs 角色修改参数
type RoleEditArgs struct {
	Name  string `json:"name" form:"name"`   //角色名称
	Intro string `json:"intro" form:"intro"` //角色描述
	Level *int   `json:"level" form:"level"` //角色层级-不传则不变，必须大于当前角色层级
}

// RolePermArgs 角色授权参数
// 对外字段名保持 permission_ids，不跟随内部缩写变化
type RolePermArgs struct {
	PermIds []string `json:"permission_ids" form:"permission_ids"` //权限id集合
}
