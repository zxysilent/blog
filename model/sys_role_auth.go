package model

import "time"

// SysRoleMenu 角色认证
type SysRoleAuth struct {
	Id     int       `xorm:"not null pk autoincr INT(11)"` //主键
	RoleId int       `xorm:"not null default 0 INT(11)"`   //RoleId
	AuthId int       `xorm:"not null default 0 INT(11)"`   //MenuId
	Ctime  time.Time `xorm:"not null DATETIME" json:"ctime"`
}
