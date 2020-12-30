package model

import "time"

// SysRole 角色表
type SysRole struct {
	Id    int       `xorm:"not null pk autoincr INT(11)"`   //主键
	Name  string    `xorm:"not null VARCHAR(255)"`          //名称
	Intro string    `xorm:"not null VARCHAR(255)"`          //介绍
	Super bool      `xorm:"not null BOOL default false"`    //超级
	Ctime time.Time `xorm:"not null DATETIME" json:"ctime"` //时间
}
