package model

// Role 角色表
type Role struct {
	Id      int    `xorm:"INT PK AUTOINCR" json:"id"`         //主键
	Name    string `xorm:"VARCHAR(255)" json:"name"`          //角色名称
	Intro   string `xorm:"VARCHAR(255)" json:"intro"`         //角色描述
	Inner   bool   `xorm:"TINYINT(4) DEFAULT 0" json:"inner"` //内部禁止删除
	Updated int64  `xorm:"BIGINT" json:"-"`                   //修改时间
	Created int64  `xorm:"BIGINT" json:"created"`             //创建时间
}

func (Role) TableName() string {
	return sysPrefix + "role"
}
