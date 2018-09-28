package model

// Cate 分类
type Cate struct {
	Id    int    `xorm:"not null pk autoincr INT(11)" json:"id" form:"id"`
	Name  string `xorm:"not null unique VARCHAR(64)" json:"name" form:"name"`
	Intro string `xorm:"default 'NULL' VARCHAR(64)" json:"intro" form:"intro"`
}

// CateIds 通过id返回新闻类别信息集合
func cateIds(ids []int) map[int]*Cate {
	mods := make([]Cate, 0, 6)
	DB.In("id", ids).Find(&mods)
	if len(mods) > 0 {
		mapSet := make(map[int]*Cate, len(mods))
		for idx := range mods {
			mapSet[mods[idx].Id] = &mods[idx]
		}
		return mapSet
	}
	return nil
}

// CateAll 所有分类
func CateAll() ([]Cate, error) {
	mods := make([]Cate, 0, 4)
	err := DB.Find(&mods)
	return mods, err
}
