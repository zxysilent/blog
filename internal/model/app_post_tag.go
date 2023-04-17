package model

//PostTag 文章标签
type PostTag struct {
	Id      int   `xorm:"INT(11) PK AUTOINCR" json:"id"`           //主键
	PostId  int   `xorm:"INT(11) UNIQUE(post_tag)" json:"post_id"` //文章Id
	TagId   int   `xorm:"INT(11) UNIQUE(post_tag)" json:"tag_id"`  //标签Id
	Post    *Post `xorm:"-" swaggerignore:"true" json:"post"`      //文章
	Tag     *Tag  `xorm:"-" swaggerignore:"true" json:"tag"`       //标签
	Updated int64 `xorm:"BIGINT" json:"-"`                         //修改时间
	Created int64 `xorm:"BIGINT" json:"created"`                   //创建时间
}

func (PostTag) TableName() string {
	return "post_tag"
}

type PostTagFilter struct {
	TagId int `query:"tag_id" form:"tag_id" json:"tag_id"` //标签Id
	Page
}
