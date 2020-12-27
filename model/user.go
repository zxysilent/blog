package model

import (
	"strconv"
	"time"
)

// tips int(11)、tinyint(4)、smallint(6)、mediumint(9)、bigint(20)

// User 用户
type User struct {
	Id     int       `xorm:"not null pk autoincr INT(11)" json:"id"`
	Name   string    `xorm:"not null  VARCHAR(255)" json:"name"`
	Num    string    `xorm:"not null  unique VARCHAR(255)" json:"num"`
	Pass   string    `xorm:"not null  VARCHAR(255)" json:"pass"`
	Role   Role      `xorm:"not null default 0 INT(11)" json:"role"`
	Email  string    `xorm:"not null  unique VARCHAR(255)" json:"email"`
	Phone  string    `xorm:"not null  VARCHAR(255)" json:"phone"`
	Ip     string    `xorm:"not null  VARCHAR(32)" json:"ip"`
	Ecount int       `xorm:"not null default 0 INT(11)" json:"ecount"`
	Ltime  time.Time `xorm:"not null DATETIME" json:"ltime"`
	Ctime  time.Time `xorm:"not null DATETIME" json:"ctime"`
}

// Role 权限角色
type Role uint32

// Role 30-0
// 30-21 角色
// 20-0  权限
const (
	RSup uint32 = 30 //super 		超级管理员
	RAtv uint32 = 20 //active		启用/禁用
	RBas uint32 = 10 //base 		基本权限
)

//Role 返回指定位的权限
func (rl Role) Role(r uint32) bool {
	if rl&(1<<r)>>r == 1 {
		return true
	}
	return false
}

// UserBaseRole 返回基本权限
func UserBaseRole() Role {
	return (1 << RBas) + (1 << RAtv)
}

//IsAtv 是否可用
func (rl Role) IsAtv() bool {
	return rl.Role(RAtv)
}

//UserByNum 通过账号获取用户信息
func UserByNum(num string) (*User, bool) {
	mod := &User{}
	has, _ := Db.Where("num=?", num).Get(mod)
	return mod, has
}

//UserExist 判断是否存在当前账号
func UserExist(num string) bool {
	has, _ := Db.Exist(&User{
		Num: num,
	})
	return has
}

//UserGet 查询一条用户信息
func UserGet(id int) (*User, bool) {
	mod := &User{}
	has, _ := Db.ID(id).Get(mod)
	return mod, has
}

//UserEditLogin 更新用户登陆信息
func UserEditLogin(mod *User, cols ...string) bool {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(mod.Id).Cols(cols...).Update(mod); err == nil {
		sess.Commit()
		return true
	}
	sess.Rollback()
	return false
}

// UserAdd 添加用户信息
func UserAdd(mod *User) bool {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.InsertOne(mod); err == nil {
		sess.Commit()
		return true
	}
	sess.Rollback()
	return false
}

// UserPage 通过用户类型和分页信息返回用户信息
func UserPage(rl int, lmtRl Role, pi int, ps int) ([]User, error) {
	mods := make([]User, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	sess.Where("Uid < 1 and role < ?", lmtRl)
	sess.Where("role < ?", 1<<uint32(rl+1))
	err := sess.Desc("Ctime").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// UserCount 通过用户类型返回用户信息总数
func UserCount(rl int, lmtRl Role) int {
	mod := &User{}
	sess := Db.NewSession()
	defer sess.Close()
	sess.Where("Uid < 1 and role < ?", lmtRl)
	sess.Where("role < ?", 1<<uint32(rl+1))
	count, _ := sess.Count(mod)
	return int(count)
}

//UserChgatv 更新用户状态
func UserChgatv(id int, rl ...Role) bool {
	mod := new(User)
	if has, _ := Db.ID(id).Get(mod); !has {
		return false
	}
	//mod.Role = mod.Role ^ (1 << RAtv)
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if len(rl) > 0 {
		sess.Where("role < ?", rl[0])
	}
	if _, err := sess.ID(id).Cols("Role").Update(mod); err == nil {
		sess.Commit()
		return true
	}
	sess.Rollback()
	return false
}

//UserPass 修改密码
func UserPass(id int, pass string, rl ...Role) bool {
	mod := &User{Pass: pass}
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	sess.ID(id)
	if len(rl) > 0 {
		sess.Where("role < ?", rl[0])
	}
	affect, err := sess.Update(mod)
	if affect >= 0 && err == nil {
		sess.Commit()
		return true
	}
	sess.Rollback()
	return false
}

//UserEdit 修改用户信息 不包括密码
func UserEdit(mod *User, rl Role, cols ...string) bool {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if rl > 0 {
		sess.Where("role < ?", rl)
	}
	if _, err := sess.ID(mod.Id).Cols(cols...).Update(mod); err == nil {
		sess.Commit()
		return true
	}
	sess.Rollback()
	return false
}

// UserDrop 删除用户信息
func UserDrop(id int, rl Role) bool {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if rl > 0 {
		sess.Where("role < ?", rl)
	}
	if _, err := sess.ID(id).Delete(&User{}); err == nil {
		Db.ClearCacheBean(&User{}, strconv.Itoa(id))
		sess.Commit()
		return true
	}
	sess.Rollback()
	return false
}
