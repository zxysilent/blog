package model

import "time"

// SysAuth 认证
type SysAuth struct {
	Id    int       `xorm:"pk autoincr INT(11)"`   //主键
	Path  string    `xorm:"VARCHAR(255)"`          //路径
	Intro string    `xorm:"VARCHAR(255)"`          //介绍
	Ctime time.Time `xorm:"DATETIME" json:"ctime"` //时间
}

// SysAuthGet 单条认证信息
func SysAuthGet(id int) (*SysAuth, bool) {
	mod := &SysAuth{}
	has, _ := Db.ID(id).Get(mod)
	return mod, has
}

// SysAuthAll 所有认证信息
func SysAuthAll() ([]SysAuth, error) {
	mods := make([]SysAuth, 0, 8)
	err := Db.Find(&mods)
	return mods, err
}

// SysAuthPage 认证分页信息
func SysAuthPage(pi int, ps int, cols ...string) ([]SysAuth, error) {
	mods := make([]SysAuth, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// SysAuthCount 认证分页信息总数
func SysAuthCount() int {
	mod := &SysAuth{}
	sess := Db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// SysAuthAdd 添加认证信息
func SysAuthAdd(mod *SysAuth) error {
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

// SysAuthEdit 编辑认证信息
func SysAuthEdit(mod *SysAuth, cols ...string) error {
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

// SysAuthIds 返回认证信息-ids
func SysAuthIds(ids []int) map[int]*SysAuth {
	mods := make([]SysAuth, 0, len(ids))
	Db.In("id", ids).Find(&mods)
	if len(mods) > 0 {
		mapMods := make(map[int]*SysAuth, len(mods))
		for idx := range mods {
			mapMods[mods[idx].Id] = &mods[idx]
		}
		return mapMods
	}
	return nil
}

// SysAuthDrop 删除单条认证信息
func SysAuthDrop(id int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&SysAuth{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}
