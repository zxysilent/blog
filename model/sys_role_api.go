package model

import "time"

// RoleApi 角色接口
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
