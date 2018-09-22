package model

//PostTag 文章标签
type PostTag struct {
	Id     int  `xorm:"not null pk autoincr INT(11)"`
	PostId int  `xorm:"not null unique(post_tag) INT(11)"`
	TagId  int  `xorm:"not null unique(post_tag) INT(11)"`
	Tag    *Tag `xorm:"- <- ->"`
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
