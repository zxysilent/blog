package model

// Tag 标签
type Tag struct {
	Id       int    `xorm:"not null pk autoincr INT(11)"`
	Name     string `xorm:"default 'NULL' unique VARCHAR(255)"`
	Pathname string `xorm:"default 'NULL' VARCHAR(255)"`
}

// tagIds 通过id返回新闻类别信息集合
func tagIds(ids []int) map[int]*Tag {
	mods := make([]Tag, 0, 6)
	DB.In("id", ids).Find(&mods)
	if len(mods) > 0 {
		mapTag := make(map[int]*Tag, len(mods))
		for idx := range mods {
			mapTag[mods[idx].Id] = &mods[idx]
		}
		return mapTag
	}
	return nil
}

// TagAll 所有标签
func TagAll() ([]Tag, error) {
	mods := make([]Tag, 0, 8)
	err := DB.Asc("id").Find(&mods)
	return mods, err
}
