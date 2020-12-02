package model

import (
	"strconv"
	"time"

	"github.com/zxysilent/utils"
)

// // Post 文章
// type Post struct {
// 	Id              int       `xorm:"pk autoincr INT(11)" json:"id" form:"id"`
// 	UserId          int       `xorm:"INT(11)" json:"user_id" form:"user_id"`
// 	CateId          int       `xorm:"INT(11)" json:"cate_id" form:"cate_id"`
// 	Cate            *Cate     `xorm:"- <- ->" json:"cate,omitempty" form:"cate"`
// 	Type            int       `xorm:"TINYINT(11)" json:"type" form:"type"`     //0 为文章，1 为页面
// 	Status          int       `xorm:"TINYINT(11)" json:"status" form:"status"` //'0 为草稿，1 为待审核，2 为已拒绝，3 为已经发布
// 	Title           string    `xorm:"VARCHAR(255)" json:"title" form:"title"`
// 	Path            string    `xorm:"unique VARCHAR(255)" json:"path" form:"path"`
// 	Summary         string    `xorm:"LONGTEXT" json:"summary,omitempty" form:"summary"`
// 	MarkdownContent string    `xorm:"LONGTEXT" json:"markdown_content,omitempty" form:"markdown_content"`
// 	Content         string    `xorm:"LONGTEXT" json:"content,omitempty" form:"content"`
// 	AllowComment    bool      `xorm:"TINYINT(4)" json:"allow_comment" form:"allow_comment"`
// 	CreateTime      time.Time `xorm:"index DATETIME" json:"create_time" form:"create_time"`
// 	UpdateTime      time.Time `xorm:"DATETIME" json:"update_time" form:"update_time"`
// 	IsPublic        bool      `xorm:"TINYINT(4)" json:"is_public" form:"is_public"`
// 	CommentNum      int       `xorm:"INT(11)" json:"comment_num" form:"comment_num"`
// 	Options         string    `xorm:"TEXT" json:"options,omitempty" form:"options"`
// }

// Post 文章
type Post struct {
	Id              int       `xorm:"not null pk autoincr INT(11)" json:"id"`
	CateId          int       `xorm:"not null default 0 INT(11)" json:"cate_id"`
	Cate            *Cate     `xorm:"-" json:"cate,omitempty" form:"cate"`
	UserId          int       `xorm:"not null INT(11)" json:"user_id"`
	Type            int       `xorm:"not null default 0 comment('0 为文章，1 为页面') TINYINT(4)" json:"type"`
	Status          int       `xorm:"not null default 0 comment('0 为草稿，1 为待审核，2 为已拒绝，3 为已经发布') TINYINT(4)" json:"status"`
	Title           string    `xorm:"not null VARCHAR(255)" json:"title"`
	Path            string    `xorm:"not null  comment('URL 的 path') VARCHAR(255)" json:"path"`
	Summary         string    `xorm:"not null comment('摘要') LONGTEXT" json:"summary"`
	MarkdownContent string    `xorm:"not null LONGTEXT" json:"markdown_content"`
	Content         string    `xorm:"not null LONGTEXT" json:"content"`
	AllowComment    bool      `xorm:"not null default 1 comment('1 为允许， 0 为不允许') TINYINT(4)" json:"allow_comment"`
	CreateTime      time.Time `xorm:"default 'NULL' index DATETIME" json:"create_time"`
	UpdateTime      time.Time `xorm:"not null DATETIME" json:"update_time"`
	IsPublic        bool      `xorm:"not null default 1 comment('1 为公开，0 为不公开') TINYINT(4)" json:"is_public"`
	CommentNum      int       `xorm:"not null default 0 INT(11)" json:"comment_num"`
	Options         string    `xorm:"not null  comment('一些选项，JSON 结构') VARCHAR(4096)" json:"options"`
}

// Archive 归档
type Archive struct {
	Time  time.Time // 日期
	Posts []Post    //文章
}

// PostPage 分页
func PostPage(pi, ps int) ([]Post, error) {
	mods := make([]Post, 0, ps)
	err := Db.Cols("id", "title", "path", "create_time", "summary", "comment_num", "options").Where("Type = 0 and Is_Public = 1 and Status = 3 ").Desc("create_time").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// PostCount 返回总数
func PostCount() int {
	mod := &Post{
		Type:     0,
		IsPublic: true,
	}
	count, _ := Db.UseBool("is_public").Count(mod)
	return int(count)
}

// PostArchive 归档
func PostArchive() ([]Archive, error) {
	posts := make([]Post, 0, 8)
	err := Db.Cols("id", "title", "path", "create_time").Where("Type = 0 and Is_Public = 1 and Status = 3 ").Desc("create_time").Find(&posts)
	if err != nil {
		return nil, err
	}
	mods := make([]Archive, 0, 8)
	for _, v := range posts {
		v.MarkdownContent = ""
		v.Content = ""
		v.Summary = ""
		v.Options = ""
		if idx := archInOf(v.CreateTime, mods); idx == -1 { //没有
			mods = append(mods, Archive{
				Time:  v.CreateTime,
				Posts: []Post{v},
			})
		} else { //有
			mods[idx].Posts = append(mods[idx].Posts, v)
		}
	}
	return mods, nil
}

func archInOf(time time.Time, mods []Archive) int {
	for idx := 0; idx < len(mods); idx++ {
		if time.Year() == mods[idx].Time.Year() && time.Month() == mods[idx].Time.Month() {
			return idx
		}
	}
	return -1
}

// PostPath 一条post
func PostPath(path string) (*Post, *Naver, bool) {
	mod := &Post{
		Path:     path,
		Type:     0,
		IsPublic: true,
	}
	has, _ := Db.UseBool("is_public").Get(mod)
	if has {
		mod.Cate, _ = CateGet(mod.CateId)
		naver := &Naver{}
		p := Post{}
		b, _ := Db.Where("Type = 0 and Is_Public = 1 and Status = 3 and Create_Time <?", mod.CreateTime.Format(utils.FormatDateTime)).Desc("Create_Time").Get(&p)
		if b {
			// <a href="{{.Naver.Prev}}" class="prev">&laquo; 上一页</a>
			naver.Prev = `<a href="/post/` + p.Path + `.html" class="prev">&laquo; ` + p.Title + `</a>`
		}
		n := Post{}
		b1, _ := Db.Where("Type = 0 and Is_Public = 1 and Status = 3 and Create_Time >?", mod.CreateTime.Format(utils.FormatDateTime)).Asc("Create_Time").Get(&n)
		if b1 {
			//<a href="{{.Naver.Next}}" class="next">下一页 &raquo;</a>
			naver.Next = `<a href="/post/` + n.Path + `.html" class="next"> ` + n.Title + ` &raquo;</a>`
		}
		return mod, naver, true
	}
	return nil, nil, has
}

//PostSingle 单页面 page
func PostSingle(path string) (*Post, bool) {
	mod := &Post{
		Path: path,
		Type: 1,
	}
	has, _ := Db.Get(mod)
	return mod, has
}

// PostPageAll 所有页面
func PostPageAll() ([]Post, error) {
	mods := make([]Post, 0, 4)
	err := Db.Cols("id", "title", "path", "create_time", "summary", "comment_num", "options", "update_time", "is_public", "status").Where("Type = 1").Desc("create_time").Find(&mods)
	return mods, err
}

//PostGet 一个
func PostGet(id int) (*Post, bool) {
	mod := &Post{
		Id: id,
	}
	has, _ := Db.Get(mod)
	if has {
		mod.Summary = ""
	}
	return mod, has
}

// postIds 通过id返回文章集合
func postIds(ids []int) map[int]*Post {
	mods := make([]Post, 0, 6)
	Db.Cols("id", "title", "path", "cate_id", "create_time", "summary", "comment_num", "options").In("id", ids).Find(&mods)
	if len(mods) > 0 {
		mapSet := make(map[int]*Post, len(mods))
		for idx := range mods {
			mapSet[mods[idx].Id] = &mods[idx]
		}
		return mapSet
	}
	return nil
}

//PostExist 判断是否存在
func PostExist(ptah string) bool {
	has, _ := Db.Exist(&Post{
		Path: ptah,
	})
	return has
}

// PostEdit 修改文章/页面
func PostEdit(mod *Post) bool {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	affect, err := sess.ID(mod.Id).Cols("Cate_id", "Status", "Title", "Summary", "Markdown_Content", "Content", "allow_comment", "Create_Time", "Comment_Num", "update_time", "is_public").Update(mod)
	if affect >= 0 && err == nil {
		sess.Commit()
		return true
	}
	sess.Rollback()
	return false
}

// PostAdd 添加文章/页面
func PostAdd(mod *Post) bool {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	affect, _ := sess.InsertOne(mod)
	if affect != 1 {
		sess.Rollback()
		return false
	}
	sess.Commit()
	return true
}

// PostDrop 删除
func PostDrop(id int) bool {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if affect, err := sess.ID(id).Delete(&Post{}); affect > 0 && err == nil {
		sess.Commit()
		Db.ClearCacheBean(&Post{}, strconv.Itoa(id))
		return true
	}
	sess.Rollback()
	return false
}
