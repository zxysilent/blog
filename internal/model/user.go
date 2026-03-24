package model

import (
	"golang.org/x/crypto/bcrypt"
)

// User 管理用户
type User struct {
	Id      string `xorm:"VARCHAR(64) PK" json:"id"`           //主键
	Num     string `xorm:"INDEX VARCHAR(64)" json:"num"`       //账号
	Name    string `xorm:"VARCHAR(64)" json:"name"`            //名称
	Passwd  string `xorm:"VARCHAR(64)" json:"-"`               //密码
	Avatar  string `xorm:"VARCHAR(255)" json:"avatar"`         //头像
	Phone   string `xorm:"VARCHAR(64)" json:"phone"`           //手机
	Email   string `xorm:"VARCHAR(64)" json:"email"`           //邮箱
	RoleId  string `xorm:"INT DEFAULT 0" json:"role_id"`       //角色
	Role    *Role  `xorm:"-" swaggerignore:"true" json:"role"` //用户角色
	Type    int    `xorm:"INT(11)" json:"type"`                //用户类型,统一用户表
	Updated int64  `xorm:"BIGINT" json:"updated"`              //修改时间
	Created int64  `xorm:"BIGINT" json:"created"`              //创建时间
}

func (s *User) Enc(raw string) (string, error) {
	passwd, err := bcrypt.GenerateFromPassword([]byte(s.Passwd), bcrypt.DefaultCost)
	return string(passwd), err
}

func (User) TableName() string {
	return "user"
}

type UserLoginArgs struct {
	Num    string `json:"num" form:"num"`       //账号
	Vcode  string `form:"vcode" json:"vcode"`   //验证码
	Vreal  string `form:"vreal" json:"vreal"`   //验证码原始值
	Passwd string `json:"passwd" form:"passwd"` //密码
}

type UserRolePageArgs struct {
	RoleId string `query:"role_id" form:"role_id" json:"role_id"`
	List
}

type UserMultArgs struct {
	Mult   string `json:"mult" form:"mult" query:"mult"` //多条件信息
	RoleId string `query:"role_id" form:"role_id" json:"role_id"`
}

type UserListArgs struct {
	List
}
