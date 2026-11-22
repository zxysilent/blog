package model

import "blog/pkg/rndm"

// PostTag 文章标签
type PostTag struct {
	Id      string `xorm:"VARCHAR(64) PK" json:"id"`                    //主键
	PostId  string `xorm:"UNIQUE(post_tag) VARCHAR(64)" json:"post_id"` //文章Id
	TagId   string `xorm:"UNIQUE(post_tag) VARCHAR(64)" json:"tag_id"`  //标签Id
	Post    *Post  `xorm:"-" swaggerignore:"true" json:"post"`          //文章
	Tag     *Tag   `xorm:"-" swaggerignore:"true" json:"tag"`           //标签
	Updated int64  `xorm:"BIGINT" json:"-"`                             //修改时间
	Created int64  `xorm:"BIGINT" json:"created"`                       //创建时间
}

func (PostTag) TableName() string {
	return "post_tag"
}

// BeforeInsert 插入前生成主键
func (s *PostTag) BeforeInsert() {
	if s.Id == "" {
		s.Id = rndm.ULID()
	}
}

// IsPkZero 主键是否为空
func (s PostTag) IsPkZero() bool {
	return s.Id == ""
}

// GetId 返回主键
func (s PostTag) GetId() string {
	return s.Id
}

type PostTagFilter struct {
	TagId string `query:"tag_id" form:"tag_id" json:"tag_id"` //标签Id
	List
}
