package model

import "time"

// Post 文章
type Post struct {
	Id       int    `xorm:"INT(11) PK AUTOINCR" json:"id"`   //主键
	CateId   int    `xorm:"INT(11)" json:"cate_id"`          //分类Id
	Kind     int    `xorm:"INT(11)" json:"kind"`             //类型1-文章，2-页面，3-笔记
	Status   int    `xorm:"INT(11)" json:"status"`           //状态1-草稿，2-已发布
	Title    string `xorm:"VARCHAR(255)" json:"title"`       //标题
	Path     string `xorm:"VARCHAR(255) UNIQUE" json:"path"` //访问路径
	Summary  string `xorm:"TEXT" json:"summary"`             //摘要
	Markdown string `xorm:"MEDIUMTEXT" json:"markdown"`      //markdown内容
	Richtext string `xorm:"MEDIUMTEXT" json:"richtext"`      //富文本内容
	Allow    bool   `xorm:"TINYINT(4)" json:"allow"`         //允许评论
	Tags     []Tag  `xorm:"-" json:"tags"`                   //标签
	Cate     *Cate  `xorm:"-" json:"cate"`                   //分类
	Updated  int64  `xorm:"BIGINT INDEX" json:"updated"`     //修改时间
	Created  int64  `xorm:"BIGINT" json:"created"`           //创建时间
}

func (Post) TableName() string {
	return "post"
}

type PostPart struct {
	Id      int    `xorm:"INT(11) PK AUTOINCR" json:"id"`   //主键
	Kind    int    `xorm:"INT(11)" json:"kind"`             //类型1-文章，2-页面
	Status  int    `xorm:"INT(11)" json:"status"`           //状态1-草稿，2-已发布
	Title   string `xorm:"VARCHAR(255)" json:"title"`       //标题
	Path    string `xorm:"VARCHAR(255) UNIQUE" json:"path"` //访问路径
	Summary string `xorm:"TEXT" json:"summary"`             //摘要
	Tags    []Tag  `xorm:"-" json:"tags"`                   //标签
	Updated int64  `xorm:"BIGINT INDEX" json:"updated"`     //修改时间
	Created int64  `xorm:"BIGINT" json:"created"`           //创建时间
}

// PostArchive 归档
type PostArchive struct {
	Time  time.Time  `json:"time"`  // 日期
	Posts []PostPart `json:"posts"` //文章
}

func (PostPart) TableName() string {
	return "post"
}

const (
	KindPost = 1 //文章
	KindPage = 2 //页面
	KindNote = 3 //笔记
)
const (
	PostStatusDraft  = 1 //草稿
	PostStatusFinish = 2 //完成
)

type PostFilter struct {
	Kind   *int `query:"kind" form:"kind" json:"kind"`
	CateId *int `query:"cate_id" form:"cate_id" json:"cate_id"`
	Status *int `query:"status" form:"status" json:"status"`
}

type PostFilterList struct {
	List
	PostFilter
}

type PostFilterPage struct {
	Page
	PostFilter
}

type PostShareArgs struct {
	PostId int `json:"post_id"`
	Day    int `json:"day"`
}
