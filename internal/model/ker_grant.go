package model

// Grant 授权
type Grant struct {
	Id      int      `xorm:"INT PK AUTOINCR" json:"id"`         //主键
	Guid    string   `xorm:"VARCHAR(64) UNIQUE" json:"guid"`    //全局唯一标识符（GUID，Globally Unique Identifier）
	Name    string   `xorm:"VARCHAR(255)" json:"name"`          //名称
	Group   string   `xorm:"VARCHAR(255)" json:"group"`         //组
	Routes  []string `xorm:"VARCHAR(512)" json:"routes"`        //api 路由
	Sort    int      `xorm:"INT DEFAULT 1000" json:"sort"`      //排序id
	Inner   bool     `xorm:"TINYINT(4) DEFAULT 0" json:"inner"` //内部禁止删除
	Updated int64    `xorm:"BIGINT" json:"updated"`             //修改时间
	Created int64    `xorm:"BIGINT" json:"created"`             //创建时间
}

// 根角色
const RootRoleId = 1

func (Grant) TableName() string {
	return sysPrefix + "grant"
}
