package model

import (
	"time"
)

// User 用户
type User struct {
	Id       int       `xorm:"INT(11) PK AUTOINCR comment('主键')" json:"id"`              //主键
	Name     string    `xorm:"VARCHAR(255) comment('姓名')" json:"name"`                   //姓名
	Num      string    `xorm:"VARCHAR(255) UNIQUE comment('账号')" json:"num"`             //账号
	OpenidQq string    `xorm:"VARCHAR(64) UNIQUE comment('qq_openid')" json:"openid_qq"` //qq
	Passwd   string    `xorm:"VARCHAR(255) comment('密码')" json:"passwd"`                 //密码
	Email    string    `xorm:"VARCHAR(255) comment('邮箱')" json:"email"`                  //邮箱
	Phone    string    `xorm:"VARCHAR(255) comment('电话')" json:"phone"`                  //电话号码
	Ecount   int       `xorm:"INT(11) DEFAULT 0 comment('错误次数')" json:"ecount"`          //错误次数
	Ltime    time.Time `xorm:"DATETIME comment('上次登录时间')" json:"ltime"`                  //上次登录时间
	Ctime    time.Time `xorm:"DATETIME comment('创建时间')" json:"ctime"`                    //创建时间
}

func (User) TableName() string {
	return "sys_user"
}

//UserLogin 用户登录
func UserLogin(num string) (*User, bool) {
	mod := &User{}
	has, _ := db.Where("num = ?", num).Get(mod)
	return mod, has
}

// UserGet 单条用户信息
func UserGet(id int) (*User, bool) {
	mod := &User{}
	has, _ := db.ID(id).Get(mod)
	return mod, has
}

// UserOpenidQq 单条用户信息
func UserOpenidQq(openid string) (*User, bool) {
	mod := &User{}
	has, _ := db.Where("openid_qq = ?", openid).Get(mod)
	return mod, has
}

// UserEdit 编辑用户信息
func UserEdit(mod *User, cols ...string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(mod.Id).Cols(cols...).Update(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}
