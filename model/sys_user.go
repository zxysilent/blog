package model

import "time"

// SysUser 用户
type SysUser struct {
	Id     int       `xorm:" INT(11) pk autoincr" json:"id"`   //主键
	Name   string    `xorm:"VARCHAR(255)" json:"name"`         //姓名
	Num    string    `xorm:"VARCHAR(255) unique" json:"num"`   //账号
	Passwd string    `xorm:"VARCHAR(255)" json:"passwd"`       //密码
	RoleId int       `xorm:"default 0 INT(11)" json:"role_id"` //角色
	Email  string    `xorm:"VARCHAR(255)" json:"email"`        //邮箱
	Phone  string    `xorm:"VARCHAR(255)" json:"phone"`        //电话
	Ip     string    `xorm:"VARCHAR(32)" json:"ip"`            //IP
	Lock   bool      `xorm:"default 0 TINYINT(4)" json:"lock"` //锁定
	Ecount int       `xorm:"default 0 INT(11)" json:"ecount"`  //错误次数
	Ltime  time.Time `xorm:"DATETIME" json:"ltime"`            //上次登录时间
	Ctime  time.Time `xorm:"DATETIME" json:"ctime"`            //创建时间
	Auths  []Api     `xorm:"-" json:"auths"`                   //权限集合
	Menus  []Menu    `xorm:"-" json:"menus"`                   //菜单导航
}

//SysUserByNum 通过账号获取用户信息
func SysUserByNum(num string) (*User, bool) {
	mod := &User{}
	has, _ := Db.Where("num=?", num).Get(mod)
	return mod, has
}

//SysUserExist 判断是否存在当前账号
func SysUserExist(num string) bool {
	has, _ := Db.Exist(&User{
		Num: num,
	})
	return has
}

// SysUserGet 单条用户信息
func SysUserGet(id int) (*SysUser, bool) {
	mod := &SysUser{}
	has, _ := Db.ID(id).Get(mod)
	return mod, has
}

// SysUserAll 所有用户信息
func SysUserAll() ([]SysUser, error) {
	mods := make([]SysUser, 0, 8)
	err := Db.Find(&mods)
	return mods, err
}

// SysUserPage 用户分页信息
func SysUserPage(pi int, ps int, cols ...string) ([]SysUser, error) {
	mods := make([]SysUser, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Desc("Utime").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// SysUserCount 用户分页信息总数
func SysUserCount() int {
	mod := &SysUser{}
	sess := Db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// SysUserAdd 添加用户信息
func SysUserAdd(mod *SysUser) error {
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

// SysUserEdit 编辑用户信息
func SysUserEdit(mod *SysUser, cols ...string) error {
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

// SysUserIds 返回用户信息-ids
func SysUserIds(ids []int) map[int]*SysUser {
	mods := make([]SysUser, 0, len(ids))
	Db.In("id", ids).Find(&mods)
	if len(mods) > 0 {
		mapMods := make(map[int]*SysUser, len(mods))
		for idx := range mods {
			mapMods[mods[idx].Id] = &mods[idx]
		}
		return mapMods
	}
	return nil
}

// SysUserDrop 删除单条用户信息
func SysUserDrop(id int) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(id).Delete(&SysUser{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}
