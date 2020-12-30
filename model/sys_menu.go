package model

import "time"

// SysMenu 菜单导航
type SysMenu struct {
	Id    int       `xorm:"not null pk autoincr INT(11)"`   //主键
	Pid   int       `xorm:"not null default 0 INT(11)"`     //父id
	Path  string    `xorm:"not null VARCHAR(255)"`          //路径
	Intro string    `xorm:"not null VARCHAR(255)"`          //介绍
	Sort  int       `xorm:"not null default 999 INT(11)"`   //排序
	Ctime time.Time `xorm:"not null DATETIME" json:"ctime"` //时间
}
