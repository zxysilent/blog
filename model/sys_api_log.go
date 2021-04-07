package model

import "time"

// SysApiLog 调用表
type SysApiLog struct {
	Id    int       `xorm:"pk autoincr INT(11)"`   //主键
	Ctime time.Time `xorm:"DATETIME" json:"ctime"` //时间
}

// SysApiLogGet 单条调用信息
func SysApiLogGet(id int) (*SysApiLog, bool) {
	mod := &SysApiLog{}
	has, _ := Db.ID(id).Get(mod)
	return mod, has
}

// SysApiLogAll 所有调用信息
func SysApiLogAll() ([]SysApiLog, error) {
	mods := make([]SysApiLog, 0, 8)
	err := Db.Find(&mods)
	return mods, err
}

// SysApiLogPage 调用分页信息
func SysApiLogPage(pi int, ps int, cols ...string) ([]SysApiLog, error) {
	mods := make([]SysApiLog, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Desc("Utime").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// SysApiLogCount 调用分页信息总数
func SysApiLogCount() int {
	mod := &SysApiLog{}
	sess := Db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// SysApiLogAdd 添加调用信息
func SysApiLogAdd(mod *SysApiLog) error {
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

// SysApiLogEdit 编辑调用信息
func SysApiLogEdit(mod *SysApiLog, cols ...string) error {
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

// SysApiLogDrop 删除单条调用信息
func SysApiLogDrop(id int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&SysApiLog{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}
