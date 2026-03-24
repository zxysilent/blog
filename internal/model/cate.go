package model

// Cate 分类/空间
type Cate struct {
	Id      string `xorm:"VARCHAR(64) PK" json:"id"`        //主键
	Kind    int    `xorm:"INT(11)" json:"kind"`             //类型1-文章，2-页面，3-笔记
	Name    string `xorm:"UNIQUE VARCHAR(255)" json:"name"` //分类名
	Intro   string `xorm:"VARCHAR(255)" json:"intro"`       //描述
	Color   string `xorm:"VARCHAR(255)" json:"color"`       //颜色
	Updated int64  `xorm:"BIGINT" json:"updated"`           //修改时间
	Created int64  `xorm:"BIGINT" json:"created"`           //创建时间
	Icon    string `xorm:"VARCHAR(255)" json:"icon"`        //图标
}

func (Cate) TableName() string {
	return "cate"
}

type CateFilter struct {
	Kind *int `query:"kind" form:"kind" json:"kind"`
	Pid  *int `query:"pid" form:"pid" json:"pid"`
}

type CateListArgs struct {
	List
	CateFilter
}
type CateFilterTree struct {
	CateFilter
}
type CateFilterPage struct {
	List
	CateFilter
}
