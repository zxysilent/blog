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

// Archive 归档
type Archive struct {
	Time  time.Time // 日期
	Posts []Post    //文章
}

// PostPage 分页
func PostPage(pi, ps int) ([]Post, error) {
	mods := make([]Post, 0, ps)
	err := Db.Cols("id", "title", "path", "created", "summary").Where("kind = 1  and status = 2 ").Desc("created").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// PostCount 返回总数
func PostCount() int {
	mod := &Post{
		Kind: 1,
	}
	count, _ := Db.Count(mod)
	return int(count)
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
		naver := &Naver{}
		p := Post{}
		b, _ := Db.Where("kind = 1 and status = 2 and created <?", mod.Created.Format(utils.StdDateTime)).Desc("Created").Get(&p)
		if b {
			// <a href="{{.Naver.Prev}}" class="prev">&laquo; 上一页</a>
			naver.Prev = `<a href="/post/` + p.Path + `.html" class="prev">&laquo; ` + p.Title + `</a>`
		}
		n := Post{}
		b1, _ := Db.Where("kind = 1  and status = 2 and created >?", mod.Created.Format(utils.StdDateTime)).Asc("Created").Get(&n)
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
		Kind: 2,
	}
	has, _ := Db.Get(mod)
	return mod, has
}

// PostPageAll 所有页面
func PostPageAll() ([]Post, error) {
	mods := make([]Post, 0, 4)
	err := Db.Cols("id", "title", "path", "created", "summary", "comment_num", "options", "updated", "is_public", "status").Where("Type = 1").Desc("created").Find(&mods)
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
	Db.Cols("id", "title", "path", "cate_id", "created", "summary").In("id", ids).Find(&mods)
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
	affect, err := sess.ID(mod.Id).Cols("cate_id", "status", "title", "summary", "markdown", "richtext", "allow", "created", "updated").Update(mod)
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
