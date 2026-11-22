package model

import "blog/pkg/rndm"

// Tag 标签
type Tag struct {
	Id      string `xorm:"VARCHAR(64) PK" json:"id"`        //主键
	Name    string `xorm:"UNIQUE VARCHAR(255)" json:"name"` //标签名
	Intro   string `xorm:"VARCHAR(255)" json:"intro"`       //描述
	Color   string `xorm:"VARCHAR(255)" json:"color"`       //颜色
	Updated int64  `xorm:"BIGINT" json:"updated"`           //修改时间
	Created int64  `xorm:"BIGINT" json:"created"`           //创建时间
}

func (Tag) TableName() string {
	return "tag"
}

// BeforeInsert 插入前生成主键
func (s *Tag) BeforeInsert() {
	if s.Id == "" {
		s.Id = rndm.ULID()
	}
}

// IsPkZero 主键是否为空
func (s Tag) IsPkZero() bool {
	return s.Id == ""
}

// GetId 返回主键
func (s Tag) GetId() string {
	return s.Id
}

// TagState 统计
type TagState struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Intro string `json:"intro"`
}

type TagListArgs struct {
	List
}
