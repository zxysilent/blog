package model

// Cate 分类
type Cate struct {
	Id       int    `xorm:"INT(11) PK AUTOINCR" json:"id"`   //主键
	Kind     int    `xorm:"INT(11)" json:"kind"`             //类型1-文章，2-页面，3-笔记
	Name     string `xorm:"UNIQUE VARCHAR(255)" json:"name"` //分类名
	Intro    string `xorm:"VARCHAR(255)" json:"intro"`       //描述
	Color    string `xorm:"VARCHAR(255)" json:"color"`       //颜色
	Updated  int64  `xorm:"BIGINT" json:"updated,omitempty"` //修改时间
	Created  int64  `xorm:"BIGINT" json:"created,omitempty"` //创建时间
	Children []Cate `xorm:"-" json:"children,omitempty"`     //树形-多级笔记使用

}

func (Cate) TableName() string {
	return "cate"
}

type CateFilter struct {
	Kind *int `query:"kind" form:"kind" json:"kind"`
}

type CateFilterList struct {
	List
	CateFilter
}

type CateFilterPage struct {
	Page
	CateFilter
}
