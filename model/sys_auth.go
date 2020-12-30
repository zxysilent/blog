package model

import "time"

// SysAuth 认证
type SysAuth struct {
	Id    int       `xorm:"not null pk autoincr INT(11)"`   //主键
	Path  string    `xorm:"not null VARCHAR(255)"`          //路径
	Intro string    `xorm:"not null VARCHAR(255)"`          //介绍
	Ctime time.Time `xorm:"not null DATETIME" json:"ctime"` //时间
}
