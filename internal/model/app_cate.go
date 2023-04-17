package model

// Cate 分类
type Cate struct {
	Id      int    `xorm:"INT(11) PK AUTOINCR" json:"id"`   //主键
	Name    string `xorm:"UNIQUE VARCHAR(255)" json:"name"` //分类名
	Intro   string `xorm:"VARCHAR(255)" json:"intro"`       //描述
	Color   string `xorm:"VARCHAR(255)" json:"color"`       //颜色
	Updated int64  `xorm:"BIGINT" json:"updated"`           //修改时间
	Created int64  `xorm:"BIGINT" json:"created"`           //创建时间
}

func (Cate) TableName() string {
	return "cate"
}

type CateFilter struct {
}

type CateFilterList struct {
	List
}

type CateFilterPage struct {
	Page
}
