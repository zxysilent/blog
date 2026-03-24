package model

//PostTag 文章标签
type PostTag struct {
	Id      string `xorm:"VARCHAR(64) PK" json:"id"`                    //主键
	PostId  string `xorm:"VARCHAR(64) UNIQUE(post_tag)" json:"post_id"` //文章Id
	TagId   string `xorm:"VARCHAR(64) UNIQUE(post_tag)" json:"tag_id"`  //标签Id
	Post    *Post  `xorm:"-" swaggerignore:"true" json:"post"`          //文章
	Tag     *Tag   `xorm:"-" swaggerignore:"true" json:"tag"`           //标签
	Updated int64  `xorm:"BIGINT" json:"-"`                             //修改时间
	Created int64  `xorm:"BIGINT" json:"created"`                       //创建时间
}

func (PostTag) TableName() string {
	return "post_tag"
}

type PostTagFilter struct {
	TagId string `query:"tag_id" form:"tag_id" json:"tag_id"` //标签Id
	List
}
