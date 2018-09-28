package model

// Tag 标签
type Tag struct {
	Id    int    `xorm:"not null pk autoincr INT(11)"`
	Name  string `xorm:"default 'NULL' unique VARCHAR(64)"`
	Intro string `xorm:"default 'NULL' VARCHAR(64)"`
}

// tagIds 通过id返回标签集合
func tagIds(ids []int) map[int]*Tag {
	mods := make([]Tag, 0, 6)
	DB.In("id", ids).Find(&mods)
	if len(mods) > 0 {
		mapSet := make(map[int]*Tag, len(mods))
		for idx := range mods {
			mapSet[mods[idx].Id] = &mods[idx]
		}
		return mapSet
	}
	return nil
}

// TagName 通过name 查询标签
func TagName(nam string) (*Tag, bool) {
	mod := &Tag{
		Name: nam,
	}
	has, _ := DB.Get(mod)
	return mod, has
}

// TagAll 所有标签
func TagAll() ([]Tag, error) {
	mods := make([]Tag, 0, 8)
	err := DB.Asc("id").Find(&mods)
	return mods, err
}
