package model

//PostTag 文章标签
type PostTag struct {
	Id     int   `xorm:"not null pk autoincr INT(11)"`
	PostId int   `xorm:"not null unique(post_tag) INT(11)"`
	Post   *Post `xorm:"- <- ->"`
	TagId  int   `xorm:"not null unique(post_tag) INT(11)"`
	Tag    *Tag  `xorm:"- <- ->"`
}

// PostTags 文章对应的标签
func PostTags(pid int) ([]PostTag, error) {
	mods := make([]PostTag, 0, 4)
	err := DB.Where("post_id = ? ", pid).Find(&mods)
	if err == nil {
		ids := make([]int, 0, len(mods))
		for idx := range mods {
			if !inOf(mods[idx].TagId, ids) {
				ids = append(ids, mods[idx].TagId)
			}
		}
		mapTag := tagIds(ids)
		if mapTag != nil {
			for idx := range mods {
				mods[idx].Tag = mapTag[mods[idx].TagId]
			}
		}
	}
	return mods, nil
}

// PostTagIds 文章对应的标签id
func PostTagIds(pid int) []int {
	mods := make([]PostTag, 0, 4)
	DB.Where("post_id = ? ", pid).Find(&mods)
	if len(mods) > 0 {
		ids := make([]int, 0, len(mods))
		for i := range mods {
			ids = append(ids, mods[i].TagId)
		}
		return ids
	}
	return nil
}

// TagPostCount 通过标签查询文章分页总数
func TagPostCount(tid int) int {
	var count int
	DB.SQL(`SELECT count(post_id) as count FROM post_tag LEFT JOIN post ON post_id=post.id WHERE is_public=1 and tag_id=?`, tid).Get(&count)
	return count
}

// TagPostList 通过标签查询文章分页
func TagPostList(tid, pi, ps int) ([]PostTag, error) {
	mods := make([]PostTag, 0, ps)
	err := DB.SQL(`SELECT post.id id,post_id,tag_id FROM post_tag LEFT JOIN post ON post_id=post.id WHERE is_public=1 and tag_id=? limit ?,? `, tid, (pi-1)*ps, ps).Find(&mods)
	if len(mods) > 0 {
		ids := make([]int, 0, len(mods))
		for idx := range mods {
			if !inOf(mods[idx].PostId, ids) {
				ids = append(ids, mods[idx].PostId)
			}
		}
		mapSet := postIds(ids)
		if mapSet != nil {
			for idx := range mods {
				mods[idx].Post = mapSet[mods[idx].PostId]
			}
		}
	}
	return mods, err
}

// TagPostAdds 添加文章标签[]
func TagPostAdds(mod *[]PostTag) bool {
	sess := DB.NewSession()
	defer sess.Close()
	sess.Begin()
	affect, _ := sess.InsertOne(mod)
	if affect < 1 {
		sess.Rollback()
		return false
	}
	sess.Commit()
	return true
}
