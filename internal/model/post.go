package model

import (
	"blog/pkg/rndm"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// Post 内容
// blog / wiki / 门户共用的一张内容表，用 Kind 区分内容形态：
//   - KindPost 文章：进首页与归档，按栏目聚合
//   - KindPage 单页：关于、友链等固定页面
//   - KindNote 笔记：用 Pid 自关联成目录树
//
// 与其他内容表的关系：
//   - CateId 指向栏目（Cate），由 Cate 决定前台路由与模板
//   - Tags 走 post_tag 关联表，全量替换
type Post struct {
	Id       string `xorm:"VARCHAR(64) PK" json:"id"`          //主键
	Pid      string `xorm:"INDEX VARCHAR(64)" json:"pid"`      //父级id，笔记用于建目录树
	CateId   string `xorm:"INDEX VARCHAR(64)" json:"cate_id"`  //栏目Id
	Uid      string `xorm:"INDEX VARCHAR(64)" json:"uid"`      //创建者用户id
	Kind     int    `xorm:"INT(11) DEFAULT 1" json:"kind"`     //类型 1文章 2页面 3笔记
	Status   int    `xorm:"INT(11) DEFAULT 1" json:"status"`   //状态 1草稿 2已发布
	Show     bool   `xorm:"BOOL DEFAULT 1" json:"show"`        //前台是否可见
	Istop    bool   `xorm:"BOOL DEFAULT 0" json:"istop"`       //置顶
	Hits     int    `xorm:"INT(11) DEFAULT 0" json:"hits"`     //浏览量
	Title    string `xorm:"INDEX VARCHAR(255)" json:"title"`   //标题
	Path     string `xorm:"VARCHAR(255) UNIQUE" json:"path"`   //访问路径，全站唯一
	Summary  string `xorm:"TEXT" json:"summary"`               //摘要
	Author   string `xorm:"VARCHAR(255)" json:"author"`        //作者
	Origin   string `xorm:"VARCHAR(255)" json:"origin"`        //来源
	Icon     string `xorm:"VARCHAR(255)" json:"icon"`          //封面图
	Markdown string `xorm:"MEDIUMTEXT" json:"markdown"`        //markdown 原文
	Richtext string `xorm:"MEDIUMTEXT" json:"richtext"`        //渲染后的富文本
	Allow    bool   `xorm:"TINYINT(4) DEFAULT 0" json:"allow"` //允许评论
	Inner    bool   `xorm:"BOOL DEFAULT 0" json:"inner"`       //内置内容，禁止删除
	Tags     []Tag  `xorm:"-" json:"tags"`                     //标签
	Cate     *Cate  `xorm:"-" json:"cate"`                     //栏目
	Updated  int64  `xorm:"BIGINT INDEX" json:"updated"`       //修改时间
	Created  int64  `xorm:"BIGINT" json:"created"`             //创建时间
}

func (Post) TableName() string {
	return "post"
}

// BeforeInsert 插入前生成主键
func (s *Post) BeforeInsert() {
	if s.Id == "" {
		s.Id = rndm.ULID()
	}
}

// IsPkZero 主键是否为空
func (s Post) IsPkZero() bool {
	return s.Id == ""
}

// GetId 返回主键
func (s Post) GetId() string {
	return s.Id
}

// PostPart 内容摘要
// 列表与树形结构使用，刻意不包含 markdown/richtext 这两个大字段
type PostPart struct {
	Id       string     `xorm:"VARCHAR(64) PK" json:"id"`    //主键
	Pid      string     `xorm:"VARCHAR(64)" json:"pid"`      //父级id
	CateId   string     `xorm:"VARCHAR(64)" json:"cate_id"`  //栏目Id
	Kind     int        `xorm:"INT(11)" json:"kind"`         //类型
	Status   int        `xorm:"INT(11)" json:"status"`       //状态
	Show     bool       `xorm:"BOOL" json:"show"`            //前台是否可见
	Istop    bool       `xorm:"BOOL" json:"istop"`           //置顶
	Hits     int        `xorm:"INT(11)" json:"hits"`         //浏览量
	Title    string     `xorm:"VARCHAR(255)" json:"title"`   //标题
	Path     string     `xorm:"VARCHAR(255)" json:"path"`    //访问路径
	Summary  string     `xorm:"TEXT" json:"summary"`         //摘要
	Author   string     `xorm:"VARCHAR(255)" json:"author"`  //作者
	Icon     string     `xorm:"VARCHAR(255)" json:"icon"`    //封面图
	Level    int        `xorm:"-" json:"level"`              //层级，笔记树使用
	Tags     []Tag      `xorm:"-" json:"tags"`               //标签
	Updated  int64      `xorm:"BIGINT" json:"updated"`       //修改时间
	Created  int64      `xorm:"BIGINT" json:"created"`       //创建时间
	Children []PostPart `xorm:"-" json:"children,omitempty"` //子级，笔记树使用
}

func (PostPart) TableName() string {
	return "post"
}

// IsPkZero 主键是否为空
func (s PostPart) IsPkZero() bool {
	return s.Id == ""
}

// GetId 返回主键
func (s PostPart) GetId() string {
	return s.Id
}

// PostArchive 归档
type PostArchive struct {
	Time  time.Time  `json:"time"`  // 日期
	Posts []PostPart `json:"posts"` //文章
}

const (
	KindPost = 1 //文章
	KindPage = 2 //页面
	KindNote = 3 //笔记
)

const (
	StatusDraft  = 1 //草稿
	StatusFinish = 2 //已发布
)

// PostPathReg 访问路径格式
// 前台以 /posts/<path>.html 访问，因此只允许 ASCII 的安全字符
var PostPathReg = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]{0,127}$`)

// PostNormalizePath 规范化访问路径：仅去首尾空格
// 不做大小写变换，避免历史链接（Path 是唯一索引）失效
func PostNormalizePath(path string) string {
	return strings.TrimSpace(path)
}

// PostPathValid 访问路径是否合法
func PostPathValid(path string) bool {
	return PostPathReg.MatchString(path)
}

// PostKindValid 内容类型是否合法
func PostKindValid(kind int) bool {
	switch kind {
	case KindPost, KindPage, KindNote:
		return true
	}
	return false
}

// PostStatusValid 状态是否合法
func PostStatusValid(status int) bool {
	switch status {
	case StatusDraft, StatusFinish:
		return true
	}
	return false
}

// PostHasContent 正文是否非空
// markdown 与 richtext 任一非空即可，兼容纯 markdown 与纯富文本两种来源
func PostHasContent(markdown, richtext string) bool {
	return strings.TrimSpace(markdown) != "" || strings.TrimSpace(richtext) != ""
}

// PostValidate 新增/修改时的业务校验
// 草稿允许正文为空，已发布必须填正文
func PostValidate(kind, status int, title, path, markdown, richtext string) error {
	if !PostKindValid(kind) {
		return fmt.Errorf("内容类型无效")
	}
	if !PostStatusValid(status) {
		return fmt.Errorf("状态无效")
	}
	if strings.TrimSpace(title) == "" {
		return fmt.Errorf("标题不能为空")
	}
	if len([]rune(strings.TrimSpace(title))) > 255 {
		return fmt.Errorf("标题不能超过 255 字符")
	}
	if !PostPathValid(path) {
		return fmt.Errorf("访问路径只能包含字母、数字、下划线和短横线，且不能以符号开头")
	}
	if status == StatusFinish && !PostHasContent(markdown, richtext) {
		return fmt.Errorf("已发布的内容必须填写正文")
	}
	return nil
}

// PostNeighbor 相邻内容，管理端预览用于上一篇/下一篇导航
// 与前台 model.Naver 的区别：这里返回结构化数据，而不是已经拼好的 <a> 标签
type PostNeighbor struct {
	Prev *PostPart `json:"prev"` //上一篇
	Next *PostPart `json:"next"` //下一篇
}

// PostFilter 内容过滤条件
type PostFilter struct {
	Kind    *int     `query:"kind" form:"kind" json:"kind"`          //类型
	Pid     *string  `query:"pid" form:"pid" json:"pid"`             //父级id
	CateId  *string  `query:"cate_id" form:"cate_id" json:"cate_id"` //栏目Id
	Fuzzy   bool     `query:"fuzzy" form:"fuzzy" json:"fuzzy"`       //cate_id 是否模糊匹配，即含其所有子栏目
	CateIds []string `query:"-" form:"-" json:"-"`                   //内部使用：栏目子树过滤
	Status  *int     `query:"status" form:"status" json:"status"`    //状态
	Show    *bool    `query:"show" form:"show" json:"show"`          //前台是否可见
	Istop   *bool    `query:"istop" form:"istop" json:"istop"`       //置顶
	Uid     *string  `query:"uid" form:"uid" json:"uid"`             //创建者
}

// PostFilterList 不分页列表参数
type PostFilterList struct {
	List
	PostFilter
}

// PostFilterPage 分页列表参数
type PostFilterPage struct {
	List
	PostFilter
}

// NoteShareArgs 笔记分享参数
type NoteShareArgs struct {
	PostId string `json:"post_id" form:"post_id"` //文章id
	Day    int    `json:"day" form:"day"`         //有效天数
}

// PostListArgs 管理端列表参数
type PostListArgs struct {
	List
	PostFilter
}

// PostAddArgs 内容新增参数
// 布尔字段用指针，未传时按服务端默认值处理
type PostAddArgs struct {
	Pid      string `json:"pid" form:"pid"`           //父级id
	CateId   string `json:"cate_id" form:"cate_id"`   //栏目Id
	Kind     int    `json:"kind" form:"kind"`         //类型，0 按文章
	Status   int    `json:"status" form:"status"`     //状态，0 按草稿
	Title    string `json:"title" form:"title"`       //标题
	Path     string `json:"path" form:"path"`         //访问路径
	Summary  string `json:"summary" form:"summary"`   //摘要
	Author   string `json:"author" form:"author"`     //作者
	Origin   string `json:"origin" form:"origin"`     //来源
	Icon     string `json:"icon" form:"icon"`         //封面图
	Markdown string `json:"markdown" form:"markdown"` //markdown 原文
	Richtext string `json:"richtext" form:"richtext"` //富文本
	Allow    *bool  `json:"allow" form:"allow"`       //允许评论
	Show     *bool  `json:"show" form:"show"`         //前台可见
	Istop    *bool  `json:"istop" form:"istop"`       //置顶
}

// PostEditArgs 内容修改参数
// 全部为指针，nil 表示不修改，支持单字段局部更新
// Kind / Inner / Created / Hits 不在可改范围内：类型创建后固定，
// 内置标记与统计值只允许服务端维护
type PostEditArgs struct {
	Pid      *string `json:"pid" form:"pid"`           //父级id
	CateId   *string `json:"cate_id" form:"cate_id"`   //栏目Id
	Status   *int    `json:"status" form:"status"`     //状态
	Title    *string `json:"title" form:"title"`       //标题
	Path     *string `json:"path" form:"path"`         //访问路径
	Summary  *string `json:"summary" form:"summary"`   //摘要
	Author   *string `json:"author" form:"author"`     //作者
	Origin   *string `json:"origin" form:"origin"`     //来源
	Icon     *string `json:"icon" form:"icon"`         //封面图
	Markdown *string `json:"markdown" form:"markdown"` //markdown 原文
	Richtext *string `json:"richtext" form:"richtext"` //富文本
	Allow    *bool   `json:"allow" form:"allow"`       //允许评论
	Show     *bool   `json:"show" form:"show"`         //前台可见
	Istop    *bool   `json:"istop" form:"istop"`       //置顶
}

// PostEditTopArgs 置顶参数
type PostEditTopArgs struct {
	Istop bool `json:"istop" form:"istop"` //是否置顶
}
