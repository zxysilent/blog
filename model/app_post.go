package model

import (
	"strconv"
	"time"

	"github.com/zxysilent/utils"
)

// Post 文章
type Post struct {
	Id       int       `xorm:"INT(11) PK AUTOINCR comment('主键')" json:"id"`            //主键
	CateId   int       `xorm:"INT(11)" json:"cate_id"`                                 //分类Id
	Kind     int       `xorm:"INT(11) comment('1-文章，2-页面')" json:"kind"`               //类型1-文章，2-页面
	Status   int       `xorm:"INT(11) comment('1-草稿，2-已发布')" json:"status"`            //状态1-草稿，2-已发布
	Title    string    `xorm:"VARCHAR(255)" json:"title"`                              //标题
	Path     string    `xorm:"VARCHAR(255) UNIQUE comment('URL 的 path') " json:"path"` //访问路径
	Summary  string    `xorm:"TEXT comment('摘要')" json:"summary"`                      //摘要
	Markdown string    `xorm:"MEDIUMTEXT comment('markdown')" json:"markdown"`         //markdown内容
	Richtext string    `xorm:"MEDIUMTEXT comment('富文本')" json:"richtext"`              //富文本内容
	Allow    bool      `xorm:"TINYINT(4) comment('允许评论')" json:"allow"`                //允许评论
	Created  time.Time `xorm:"DATETIME INDEX" json:"created"`                          //创建时间
	Updated  time.Time `xorm:"DATETIME" json:"updated"`                                //修改时间
	Tags     []Tag     `xorm:"-" json:"tags"`                                          //标签
	Cate     *Cate     `xorm:"-" json:"cate"`                                          //分类
}

const (
	PostKindPost = 1 //文章
	PostKindPage = 2 //页面
)
const (
	PostStatusDraft  = 1 //草稿
	PostStatusFinish = 2 //完成
)

// PostGet 单条文章/页面
func PostGet(id int) (*Post, bool) {
	mod := &Post{}
	has, _ := Db.ID(id).Get(mod)
	if has && mod.Kind == PostKindPost {
		tags := make([]Tag, 0, 4)
		Db.SQL("SELECT * FROM tag WHERE id IN (SELECT tag_id FROM post_tag WHERE post_id = ?)", mod.Id).Find(&tags)
		mod.Tags = tags
	}
	return mod, has
}

// PostAll 所有文章/页面
func PostAll(cateId int, kind int, cols ...string) ([]Post, error) {
	mods := make([]Post, 0, 4)
	sess := Db.NewSession()
	defer sess.Close()
	if cateId > 0 {
		sess.Where("cate_id = ?", cateId)
	}
	if kind > 0 {
		sess.Where("kind = ?", kind)
	}
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := Db.Desc("created").Find(&mods)
	return mods, err
}

//PostExist 判断是否存在
func PostExist(ptah string) bool {
	has, _ := Db.Exist(&Post{
		Path: ptah,
	})
	return has
}

// PostPage 文章/页面分页
func PostPage(cateId int, kind int, pi int, ps int, cols ...string) ([]Post, error) {
	mods := make([]Post, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if cateId > 0 {
		sess.Where("cate_id = ?", cateId)
	}
	if kind > 0 {
		sess.Where("kind = ?", kind)
	}
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Desc("created").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// PostCount 返回总数
func PostCount(cateId int, kind int) int {
	mod := &Post{}
	sess := Db.NewSession()
	defer sess.Close()
	if cateId > 0 {
		sess.Where("cate_id = ?", cateId)
	}
	if kind > 0 {
		sess.Where("kind = ?", kind)
	}
	count, _ := sess.Count(mod)
	return int(count)
}

// PostEdit 编辑文章
func PostEdit(mod *Post, cols ...string) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(mod.Id).Cols(cols...).Update(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// PostAdd 添加文章/页面
func PostAdd(mod *Post) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.InsertOne(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// PostDrop 删除单条文章
func PostDrop(id int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&Post{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	Db.ClearCacheBean(&Post{}, strconv.Itoa(id))
	return nil
}

// PostIds 通过id集合返回文章
func PostIds(ids []int) map[int]*Post {
	mods := make([]Post, 0, len(ids))
	Db.Cols("id", "title", "path", "cate_id", "created", "summary").In("id", ids).Find(&mods)
	mapSet := make(map[int]*Post, len(mods))
	for idx := range mods {
		mapSet[mods[idx].Id] = &mods[idx]
	}
	return mapSet
}

// ------------------------------------------------------ 页面使用 ------------------------------------------------------
//PostSingle 单页面 page
func PostSingle(path string) (*Post, bool) {
	mod := &Post{}
	has, _ := Db.Where("kind = 2 and path = ?", path).Get(mod)
	return mod, has
}

// ------------------------------------------------------ 归档使用 ------------------------------------------------------
// Archive 归档
type Archive struct {
	Time  time.Time `json:"time"`  // 日期
	Posts []Post    `json:"posts"` //文章
}

// PostArchive 归档
func PostArchive() ([]Archive, error) {
	posts := make([]Post, 0, 8)
	err := Db.Cols("id", "title", "path", "created").Where("kind = 1  and status = 2 ").Desc("created").Find(&posts)
	if err != nil {
		return nil, err
	}
	mods := make([]Archive, 0, 8)
	for _, v := range posts {
		v.Markdown = ""
		v.Richtext = ""
		v.Summary = ""
		if idx := archInOf(v.Created, mods); idx == -1 { //没有
			mods = append(mods, Archive{
				Time:  v.Created,
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
		Path: path,
		Kind: 1,
	}
	has, _ := Db.Get(mod)
	if has {
		mod.Cate, _ = CateGet(mod.CateId)
		if mod.Kind == PostKindPost {
			tags := make([]Tag, 0, 4)
			Db.SQL("SELECT * FROM tag WHERE id IN (SELECT tag_id FROM post_tag WHERE post_id = ?)", mod.Id).Find(&tags)
			mod.Tags = tags
		}
		naver := &Naver{}
		p := Post{}
		b, _ := Db.Where("kind = 1 and status = 2 and created >?", mod.Created.Format(utils.StdDateTime)).Asc("created").Get(&p)
		if b {
			// <a href="{{.Naver.Prev}}" class="prev">&laquo; 上一页</a>
			naver.Prev = `<a href="/post/` + p.Path + `.html" class="prev">&laquo; ` + p.Title + `</a>`
		}
		n := Post{}
		b1, _ := Db.Where("kind = 1  and status = 2 and created <?", mod.Created.Format(utils.StdDateTime)).Desc("created").Get(&n)
		if b1 {
			//<a href="{{.Naver.Next}}" class="next">下一页 &raquo;</a>
			naver.Next = `<a href="/post/` + n.Path + `.html" class="next"> ` + n.Title + ` &raquo;</a>`
		}
		return mod, naver, true
	}
	return nil, nil, has
}
