package model

// Tag 标签
type Tag struct {
	Id    int    `xorm:"INT(11) PK AUTOINCR comment('主键')" json:"id"`     //主键
	Name  string `xorm:"VARCHAR(255) UNIQUE comment('标签名称')" json:"name"` //标签名称
	Intro string `xorm:"VARCHAR(255) comment('标签描述')" json:"intro"`       //标签描述
}

// TagState 统计
type TagState struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
	Intro string `json:"intro"`
}

// TagStateAll 所有标签统计 当前标签下有文章才显示
func TagStateAll() ([]TagState, error) {
	mods := make([]TagState, 0, 8)
	err := Db.SQL("SELECT `name`,intro,count(tag_id) as count FROM post_tag ,tag WHERE tag.id=tag_id GROUP BY tag_id HAVING count>0").Find(&mods)
	return mods, err
}

// TagGetName 通过name 查询标签
func TagGetName(name string) (*Tag, bool) {
	mod := &Tag{
		Name: name,
	}
	has, _ := Db.Get(mod)
	return mod, has
}

// TagGet 单条标签信息
func TagGet(id int) (*Tag, bool) {
	mod := &Tag{}
	has, _ := Db.ID(id).Get(mod)
	return mod, has
}

// TagAll 所有标签信息
func TagAll() ([]Tag, error) {
	mods := make([]Tag, 0, 8)
	sess := Db.NewSession()
	defer sess.Close()
	err := sess.Find(&mods)
	return mods, err
}

// TagPage 标签分页信息
func TagPage(pi int, ps int, cols ...string) ([]Tag, error) {
	mods := make([]Tag, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// TagCount 标签分页信息总数
func TagCount() int {
	mod := &Tag{}
	sess := Db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// TagAdd 添加标签信息
func TagAdd(mod *Tag) error {
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

// TagEdit 编辑标签信息
func TagEdit(mod *Tag, cols ...string) error {
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

// TagDrop 删除单条标签信息
func TagDrop(id int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&Tag{}); err != nil {
		sess.Rollback()
		return err
	}
	// 删除标签接口
	// sess.Exec("DELETE FROM sys_role_api WHERE role_id = ?", id)
	sess.Commit()
	return nil
}

// TagMapIds 通过id集合返回标签信息
func TagMapIds(ids []int) map[int]*Tag {
	mods := make([]Tag, 0, len(ids))
	Db.In("id", ids).Find(&mods)
	mapSet := make(map[int]*Tag, len(mods))
	for idx := range mods {
		mapSet[mods[idx].Id] = &mods[idx]
	}
	return mapSet
}
