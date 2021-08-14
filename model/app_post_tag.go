package model

//PostTag 文章标签
type PostTag struct {
	Id     int   `xorm:"INT(11) PK AUTOINCR comment('主键')" json:"id"` //主键
	PostId int   `xorm:"INT(11) UNIQUE(post_tag)" json:"post_id"`     //文章Id
	TagId  int   `xorm:"INT(11) UNIQUE(post_tag)" json:"tag_id"`      //标签Id
	Post   *Post `xorm:"-" swaggerignore:"true" json:"post"`          //文章
	Tag    *Tag  `xorm:"-" swaggerignore:"true" json:"tag"`           //标签
}

// TagPostCount 通过标签查询文章分页总数
func TagPostCount(tagId int) int {
	var count int
	Db.SQL(`SELECT count(post_id) as count FROM post_tag WHERE tag_id=?`, tagId).Get(&count)
	return count
}

// TagPostPage 通过标签查询文章分页
func TagPostPage(tagId, pi, ps int) ([]PostTag, error) {
	mods := make([]PostTag, 0, ps)
	err := Db.SQL(`SELECT id,post_id,tag_id FROM post_tag WHERE tag_id=? ORDER BY post_id DESC LIMIT ?,? `, tagId, (pi-1)*ps, ps).Find(&mods)
	if len(mods) > 0 {
		ids := make([]int, 0, len(mods))
		for idx := range mods {
			if !inOf(mods[idx].PostId, ids) {
				ids = append(ids, mods[idx].PostId)
			}
		}
		mapSet := PostIds(ids)
		for idx := range mods {
			mods[idx].Post = mapSet[mods[idx].PostId]
		}
	}
	return mods, err
}

// TagPostAdds 添加文章标签[]
func TagPostAdds(mods *[]PostTag) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.Insert(mods); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

// TagPostDrop 删除标签时候
// 删除对应的标签_文章
func TagPostDrop(tagId int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	sess.Where("tag_id = ?", tagId)
	if _, err := sess.Delete(&PostTag{}); err != nil {
		sess.Commit()
		return err
	}
	sess.Rollback()
	return nil
}

// PostTagDrops 修改文章时候
// 删除对应的标签_文章
func PostTagDrops(postId int, tagIds []int) error {
	if len(tagIds) < 1 {
		return nil
	}
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	sess.Where("post_id = ?", postId).In("tag_id", tagIds)
	if _, err := sess.Delete(&PostTag{}); err != nil {
		sess.Commit()
		return err
	}
	sess.Rollback()
	return nil
}

// PostTagDrop 删除文章时候
// 删除对应的标签_文章
func PostTagDrop(postId int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	sess.Where("post_id = ?", postId)
	if _, err := sess.Delete(&PostTag{}); err != nil {
		sess.Commit()
		return err
	}
	sess.Rollback()
	return nil
}
