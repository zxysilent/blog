package model

//PostCate 文章分类
type PostCate struct {
	Id     int   `xorm:"not null pk autoincr INT(11)"`
	PostId int   `xorm:"not null unique(post_cate) INT(11)"`
	CateId int   `xorm:"not null unique(post_cate) INT(11)"`
	Cate   *Cate `xorm:"- <- ->"`
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
