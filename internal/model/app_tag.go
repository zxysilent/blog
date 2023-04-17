package model

// Tag 标签
type Tag struct {
	Id      int    `xorm:"INT(11) PK AUTOINCR" json:"id"`   //主键
	Name    string `xorm:"UNIQUE VARCHAR(255)" json:"name"` //标签名
	Intro   string `xorm:"VARCHAR(255)" json:"intro"`       //描述
	Color   string `xorm:"VARCHAR(255)" json:"color"`       //颜色
	Updated int64  `xorm:"BIGINT" json:"updated"`           //修改时间
	Created int64  `xorm:"BIGINT" json:"created"`           //创建时间
}

func (Tag) TableName() string {
	return "tag"
}

// TagState 统计
type TagState struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Intro string `json:"intro"`
}

type TagFilter struct {
}

type TagFilterList struct {
	List
}

type TagFilterPage struct {
	Page
}
