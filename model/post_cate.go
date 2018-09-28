package model

//PostCate 文章分类
type PostCate struct {
	Id     int   `xorm:"not null pk autoincr INT(11)" json:"id" form:"id"`
	PostId int   `xorm:"not null unique(post_cate) INT(11)" json:"post_id" form:"post_id"`
	Post   *Post `xorm:"- <- ->" json:"post" form:"post"`
	CateId int   `xorm:"not null unique(post_cate) INT(11)" json:"cate_id" form:"cate_id"`
	Cate   *Cate `xorm:"- <- ->" json:"cate" form:"cate"`
}

// PostCates 文章对应的标签
func PostCates(pid int) ([]PostCate, error) {
	mods := make([]PostCate, 0, 4)
	err := DB.Where("post_id = ? ", pid).Find(&mods)
	if err == nil {
		ids := make([]int, 0, len(mods))
		for idx := range mods {
			if !inOf(mods[idx].CateId, ids) {
				ids = append(ids, mods[idx].CateId)
			}
		}
		mapCate := cateIds(ids)
		if mapCate != nil {
			for idx := range mods {
				mods[idx].Cate = mapCate[mods[idx].CateId]
			}
		}
	}
	return mods, nil
}

// CatePostCount 通过标签查询文章分页总数
func CatePostCount(cid int) int {
	var count int
	if cid < 1 {
		DB.SQL(`SELECT count(post_id) as count FROM post_cate LEFT JOIN post ON post_id=post.id`).Get(&count)
		return count
	}
	DB.SQL(`SELECT count(post_id) as count FROM post_cate LEFT JOIN post ON post_id=post.id WHERE cate_id=?`, cid).Get(&count)
	return count
}

// CatePostList 通过标签查询文章分页
func CatePostList(cid, pi, ps int) ([]PostCate, error) {
	mods := make([]PostCate, 0, ps)
	var err error
	if cid < 1 {
		err = DB.SQL(`SELECT post.id id,post_id,cate_id FROM post_cate LEFT JOIN post ON post_id=post.id order by id desc limit ?,? `, (pi-1)*ps, ps).Find(&mods)
	} else {
		err = DB.SQL(`SELECT post.id id,post_id,cate_id FROM post_cate LEFT JOIN post ON post_id=post.id WHERE cate_id=? order by id desc limit ?,? `, cid, (pi-1)*ps, ps).Find(&mods)
	}
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
