package model

import (
	"crypto/md5"
	"encoding/hex"
)

// Admin 管理用户
type Admin struct {
	Id      int    `xorm:"PK AUTOINCR INT" json:"id"`          //主键
	Num     string `xorm:"INDEX VARCHAR(64)" json:"num"`       //账号
	Name    string `xorm:"VARCHAR(64)" json:"name"`            //名称
	Passwd  string `xorm:"VARCHAR(64)" json:"-"`               //密码
	Avatar  string `xorm:"VARCHAR(255)" json:"avatar"`         //头像
	Salt    string `xorm:"VARCHAR(64)" json:"salt"`            //
	Phone   string `xorm:"VARCHAR(64)" json:"phone"`           //手机
	Email   string `xorm:"VARCHAR(64)" json:"email"`           //邮箱
	RoleId  int    `xorm:"INT DEFAULT 0" json:"role_id"`       //角色
	Role    *Role  `xorm:"-" swaggerignore:"true" json:"role"` //用户角色
	Updated int64  `xorm:"BIGINT" json:"updated"`              //修改时间
	Created int64  `xorm:"BIGINT" json:"created"`              //创建时间
}

func (s *Admin) Enc(raw string) string {
	hash := md5.Sum([]byte(raw))
	return hex.EncodeToString(hash[:])[1:31]
}

func (Admin) TableName() string {
	return sysPrefix + "admin"
}

type AdminLoginArgs struct {
	Num    string `json:"num" form:"num"`       //账号
	Vcode  string `form:"vcode" json:"vcode"`   //验证码
	Vreal  string `form:"vreal" json:"vreal"`   //验证码原始值
	Passwd string `json:"passwd" form:"passwd"` //密码
}

func (s *AdminLoginArgs) Enc() string {
	hash := md5.Sum([]byte(s.Passwd))
	return hex.EncodeToString(hash[:])[1:31]
}

type PasswdArgs struct {
	Opasswd string `form:"opasswd" json:"opasswd"` //旧密码
	Npasswd string `form:"npasswd" json:"npasswd"` //新密码
}

func (s *PasswdArgs) Enc() {
	old := md5.Sum([]byte(s.Opasswd))
	s.Opasswd = hex.EncodeToString(old[:])[1:31]
	new := md5.Sum([]byte(s.Npasswd))
	s.Npasswd = hex.EncodeToString(new[:])[1:31]
}

type AdminRolePageArgs struct {
	RoleId int `query:"role_id" form:"role_id" json:"role_id"`
	Page
}

type AdminMultArgs struct {
	Mult   string `json:"mult" form:"mult" query:"mult"` //多条件信息
	RoleId int    `query:"role_id" form:"role_id" json:"role_id"`
}
