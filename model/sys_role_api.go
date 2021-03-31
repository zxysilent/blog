package model

import "time"

// RoleMenu 角色认证
type RoleApi struct {
	Id     int       `xorm:"pk autoincr INT(11)"` //主键
	RoleId int       `xorm:"default 0 INT(11)"`   //RoleId
	ApiId  int       `xorm:"default 0 INT(11)"`   //ApiId
	Ctime  time.Time `xorm:"DATETIME" json:"ctime"`
}

func (RoleApi) TableName() string {
	return "sys_role_api"
}

// RoleApiGet 单条角色认证信息
func RoleApiGet(id int) (*RoleApi, bool) {
	mod := &RoleApi{}
	has, _ := Db.ID(id).Get(mod)
	return mod, has
}

// RoleApiAll 所有角色认证信息
func RoleApiAll() ([]RoleApi, error) {
	mods := make([]RoleApi, 0, 8)
	err := Db.Find(&mods)
	return mods, err
}

// RoleApiPage 角色认证分页信息
func RoleApiPage(pi int, ps int, cols ...string) ([]RoleApi, error) {
	mods := make([]RoleApi, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Desc("Utime").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// RoleApiCount 角色认证分页信息总数
func RoleApiCount() int {
	mod := &RoleApi{}
	sess := Db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// RoleApiAdd 添加角色认证信息
func RoleApiAdd(mod *RoleApi) error {
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

// RoleApiEdit 编辑角色认证信息
func RoleApiEdit(mod *RoleApi, cols ...string) error {
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

// RoleApiIds 返回角色认证信息-ids
func RoleApiIds(ids []int) map[int]*RoleApi {
	mods := make([]RoleApi, 0, len(ids))
	Db.In("id", ids).Find(&mods)
	if len(mods) > 0 {
		mapMods := make(map[int]*RoleApi, len(mods))
		for idx := range mods {
			mapMods[mods[idx].Id] = &mods[idx]
		}
		return mapMods
	}
	return nil
}

// RoleApiDrop 删除单条角色认证信息
func RoleApiDrop(id int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&RoleApi{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}
