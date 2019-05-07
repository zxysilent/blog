package model

import (
	"time"
)

// tips int(11)、tinyint(4)、smallint(6)、mediumint(9)、bigint(20)

// User 用户
type User struct {
	Id     int       `xorm:"pk autoincr INT(11) not null" json:"id" form:"id"`              //主键
	Num    string    `xorm:"unique VARCHAR(32) default('') not null" json:"num" form:"num"` //账号
	Name   string    `xorm:"VARCHAR(32) default('') not null" json:"name" form:"name"`      //名称
	Pass   string    `xorm:"VARCHAR(32) default('') not null" json:"-" form:"-"`            //密码
	Role   Role      `xorm:"INT(11) default(0) not null" json:"role" form:"role"`           //权限
	Phone  string    `xorm:"VARCHAR(32) default('') not null" json:"phone" form:"phone"`    //手机
	Email  string    `xorm:"VARCHAR(255)  default('') not null" json:"email" form:"email"`  //邮箱
	Ip     string    `xorm:"VARCHAR(32)  default('') not null" json:"ip" form:"ip"`         //登录ip
	Remark string    `xorm:"VARCHAR(255) default('') not null" json:"remark" form:"remark"` //备注
	Ecount int8      `xorm:"TINYINT default(0) not null" json:"ecount" form:"ecount"`       //登录错误次数
	Ltime  time.Time `xorm:"DATETIME null" json:"ltime" form:"ltime"`                       //上次登录时间
	Ctime  time.Time `xorm:"DATETIME null" json:"ctime" form:"ctime"`                       //创建时间
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
	has, _ := DB.Where("num=?", num).Get(mod)
	return mod, has
}

//UserExist 判断是否存在当前账号
func UserExist(num string) bool {
	has, _ := DB.Exist(&User{
		Num: num,
	})
	return has
}

//UserGet 查询一条用户信息
func UserGet(id int) (*User, bool) {
	mod := &User{}
	has, _ := DB.ID(id).Get(mod)
	return mod, has
}

//UserEditLogin 更新用户登陆信息
func UserEditLogin(mod *User, cols ...string) bool {
	sess := DB.NewSession()
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
	sess := DB.NewSession()
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
	sess := DB.NewSession()
	defer sess.Close()
	sess.Where("Uid < 1 and role < ?", lmtRl)
	sess.Where("role < ?", 1<<uint32(rl+1))
	err := sess.Desc("Ctime").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// UserCount 通过用户类型返回用户信息总数
func UserCount(rl int, lmtRl Role) int {
	mod := &User{}
	sess := DB.NewSession()
	defer sess.Close()
	sess.Where("Uid < 1 and role < ?", lmtRl)
	sess.Where("role < ?", 1<<uint32(rl+1))
	count, _ := sess.Count(mod)
	return int(count)
}

//UserChgatv 更新用户状态
func UserChgatv(id int, rl ...Role) bool {
	mod := new(User)
	if has, _ := DB.ID(id).Get(mod); !has {
		return false
	}
	//mod.Role = mod.Role ^ (1 << RAtv)
	sess := DB.NewSession()
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
	sess := DB.NewSession()
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
	sess := DB.NewSession()
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

// UserDel 删除用户信息
func UserDel(id int, rl Role) bool {
	sess := DB.NewSession()
	defer sess.Close()
	sess.Begin()
	if rl > 0 {
		sess.Where("role < ?", rl)
	}
	if _, err := sess.ID(id).Delete(&User{}); err == nil {
		DB.ClearCacheBean(&User{}, string(id))
		sess.Commit()
		return true
	}
	sess.Rollback()
	return false
}
