package model

import (
	"blog/pkg/rndm"
	"regexp"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// 用户状态
const (
	UserStatusNormal  = 1 //正常
	UserStatusDisable = 2 //禁用
)

// User 用户（统一用户表，管理端与前台共用）
// 不区分「管理员/普通用户」类型：能否操作完全由角色授权决定，
// 没分配角色（或角色无权限）的账号登录后只能访问无需权限的接口。
type User struct {
	Id      string `xorm:"VARCHAR(64) PK" json:"id"`           //主键
	Num     string `xorm:"UNIQUE VARCHAR(64)" json:"num"`      //账号
	Name    string `xorm:"VARCHAR(64)" json:"name"`            //名称
	Passwd  string `xorm:"VARCHAR(255)" json:"-"`              //密码 bcrypt
	Avatar  string `xorm:"VARCHAR(255)" json:"avatar"`         //头像
	Phone   string `xorm:"VARCHAR(64)" json:"phone"`           //手机
	Email   string `xorm:"VARCHAR(64)" json:"email"`           //邮箱
	Openid  string `xorm:"INDEX VARCHAR(128)" json:"openid"`   //微信 openid，预留公众号/小程序登录
	RoleId  string `xorm:"INDEX VARCHAR(64)" json:"role_id"`   //角色id
	Role    *Role  `xorm:"-" swaggerignore:"true" json:"role"` //角色-非持久化
	Status  int    `xorm:"INT(11) DEFAULT 1" json:"status"`    //状态 1-正常 2-禁用
	Updated int64  `xorm:"BIGINT" json:"updated"`              //修改时间
	Created int64  `xorm:"BIGINT" json:"created"`              //创建时间
}

func (User) TableName() string {
	return "user"
}

// BeforeInsert 插入前生成主键
func (s *User) BeforeInsert() {
	if s.Id == "" {
		s.Id = rndm.ULID()
	}
}

// IsPkZero 主键是否为空
func (s User) IsPkZero() bool {
	return s.Id == ""
}

// GetId 返回主键
func (s User) GetId() string {
	return s.Id
}

// SetPasswd 用 bcrypt 生成密码哈希
func (u *User) SetPasswd(raw string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Passwd = string(hash)
	return nil
}

// EqualPasswd 校验明文密码
func (u *User) EqualPasswd(raw string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Passwd), []byte(raw))
}

// UserOpenidReg openid 格式：微信 openid / unionid 都是 base64url 字符集
var UserOpenidReg = regexp.MustCompile(`^[A-Za-z0-9_-]{1,128}$`)

// UserNormalizeOpenid 规范化 openid：去首尾空格
// 不改变大小写，openid 对大小写敏感
func UserNormalizeOpenid(openid string) string {
	return strings.TrimSpace(openid)
}

// UserOpenidValid openid 是否合法
// 空串表示未绑定，合法；非空必须是 base64url 字符集且不超过 128 字符
// （微信 openid 28 位，unionid 28 位，留足余量兼容其它平台）
func UserOpenidValid(openid string) bool {
	if openid == "" {
		return true
	}
	return UserOpenidReg.MatchString(openid)
}

// UserLoginArgs 登录参数
type UserLoginArgs struct {
	Num    string `json:"num" form:"num" query:"num"`          //账号
	Passwd string `json:"passwd" form:"passwd" query:"passwd"` //密码
	Vcode  string `json:"vcode" form:"vcode" query:"vcode"`    //验证码
	Vreal  string `json:"vreal" form:"vreal" query:"vreal"`    //验证码签名
}

// UserListArgs 用户列表参数
type UserListArgs struct {
	List
	Num    string `query:"num" form:"num" json:"num"`             //账号精确匹配
	RoleId string `query:"role_id" form:"role_id" json:"role_id"` //角色过滤
	Status int    `query:"status" form:"status" json:"status"`    //状态过滤
	Openid string `query:"openid" form:"openid" json:"openid"`    //openid 精确匹配
}

// UserAddArgs 用户新增参数
type UserAddArgs struct {
	Num    string `json:"num" form:"num"`         //账号
	Name   string `json:"name" form:"name"`       //名称
	Passwd string `json:"passwd" form:"passwd"`   //密码
	Avatar string `json:"avatar" form:"avatar"`   //头像
	Phone  string `json:"phone" form:"phone"`     //手机
	Email  string `json:"email" form:"email"`     //邮箱
	Openid string `json:"openid" form:"openid"`   //微信 openid
	RoleId string `json:"role_id" form:"role_id"` //角色id
	Status int    `json:"status" form:"status"`   //状态
}

// UserEditArgs 用户修改参数（指针字段用于局部更新）
type UserEditArgs struct {
	Num    *string `json:"num" form:"num"`         //账号
	Name   *string `json:"name" form:"name"`       //名称
	Passwd *string `json:"passwd" form:"passwd"`   //密码
	Avatar *string `json:"avatar" form:"avatar"`   //头像
	Phone  *string `json:"phone" form:"phone"`     //手机
	Email  *string `json:"email" form:"email"`     //邮箱
	Openid *string `json:"openid" form:"openid"`   //微信 openid，传空串表示解绑
	RoleId *string `json:"role_id" form:"role_id"` //角色id
	Status *int    `json:"status" form:"status"`   //状态
}

// MeEditArgs 修改自己资料
type MeEditArgs struct {
	Name   string `json:"name" form:"name"`     //名称
	Avatar string `json:"avatar" form:"avatar"` //头像
	Phone  string `json:"phone" form:"phone"`   //手机
	Email  string `json:"email" form:"email"`   //邮箱
}

// MePasswdArgs 修改自己密码
type MePasswdArgs struct {
	Opasswd string `json:"opasswd" form:"opasswd"` //旧密码
	Npasswd string `json:"npasswd" form:"npasswd"` //新密码
}
