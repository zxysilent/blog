package model

import "time"

// SysRoleMenu 角色菜单导航
type SysRoleMenu struct {
	Id     int       `xorm:"not null pk autoincr INT(11)"` //主键
	RoleId int       `xorm:"not null default 0 INT(11)"`   //RoleId
	MenuId int       `xorm:"not null default 0 INT(11)"`   //MenuId
	Add    bool      `xorm:"not null BOOL default false"`  //添加
	Edit   bool      `xorm:"not null BOOL default false"`  //编辑
	Drop   bool      `xorm:"not null BOOL default false"`  //删除
	Fetch  bool      `xorm:"not null BOOL default false"`  //查询
	Export bool      `xorm:"not null BOOL default false"`  //导出
	Import bool      `xorm:"not null BOOL default false"`  //导入
	Ctime  time.Time `xorm:"not null DATETIME" json:"ctime"`
}
