package model

import (
	"time"
)

type User struct {
	Id            int       `xorm:"not null pk autoincr INT(10)"`
	Name          string    `xorm:"not null default '''' unique VARCHAR(255)"`
	DisplayName   string    `xorm:"default 'NULL' VARCHAR(255)"`
	Password      string    `xorm:"not null default '''' VARCHAR(255)"`
	Type          int       `xorm:"not null default 1 comment('1 为管理员  2 为编辑') TINYINT(11)"`
	Email         string    `xorm:"not null default '''' unique VARCHAR(255)"`
	Status        int       `xorm:"not null default 1 comment('1 为有效 2 为禁用') TINYINT(11)"`
	CreateTime    time.Time `xorm:"not null DATETIME"`
	CreateIp      string    `xorm:"not null default '''' VARCHAR(20)"`
	LastLoginTime time.Time `xorm:"not null DATETIME"`
	LastLoginIp   string    `xorm:"not null default '''' VARCHAR(20)"`
	AppKey        string    `xorm:"default 'NULL' VARCHAR(255)"`
	AppSecret     string    `xorm:"default 'NULL' VARCHAR(255)"`
}
