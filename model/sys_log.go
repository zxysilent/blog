package model

import "time"

// SysLog 调用表
type SysLog struct {
	Id    int       `xorm:"pk autoincr INT(11)"`   //主键
	Ctime time.Time `xorm:"DATETIME" json:"ctime"` //时间
}

// SysLogGet 单条调用信息
func SysLogGet(id int) (*SysLog, bool) {
	mod := &SysLog{}
	has, _ := Db.ID(id).Get(mod)
	return mod, has
}

// SysLogAll 所有调用信息
func SysLogAll() ([]SysLog, error) {
	mods := make([]SysLog, 0, 8)
	err := Db.Find(&mods)
	return mods, err
}

// SysLogPage 调用分页信息
func SysLogPage(pi int, ps int, cols ...string) ([]SysLog, error) {
	mods := make([]SysLog, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Desc("Utime").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// SysLogCount 调用分页信息总数
func SysLogCount() int {
	mod := &SysLog{}
	sess := Db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// SysLogAdd 添加调用信息
func SysLogAdd(mod *SysLog) error {
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

// SysLogEdit 编辑调用信息
func SysLogEdit(mod *SysLog, cols ...string) error {
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

// SysLogDrop 删除单条调用信息
func SysLogDrop(id int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&SysLog{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}
