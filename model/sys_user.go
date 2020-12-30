package model

import "time"

// SysUser 用户
type SysUser struct {
	Id     int       `xorm:"not null pk autoincr INT(11)" json:"id"`
	Name   string    `xorm:"not null  VARCHAR(255)" json:"name"`
	Num    string    `xorm:"not null  unique VARCHAR(255)" json:"num"`
	Passwd string    `xorm:"not null  VARCHAR(255)" json:"passwd"`
	RoleId int       `xorm:"not null default 0 INT(11)" json:"role_id"`
	Email  string    `xorm:"not null  unique VARCHAR(255)" json:"email"`
	Phone  string    `xorm:"not null  VARCHAR(255)" json:"phone"`
	Ip     string    `xorm:"not null  VARCHAR(32)" json:"ip"`
	Ecount int       `xorm:"not null default 0 INT(11)" json:"ecount"`
	Ltime  time.Time `xorm:"not null DATETIME" json:"ltime"`
	Ctime  time.Time `xorm:"not null DATETIME" json:"ctime"`
}
