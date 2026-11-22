package model

import (
	"blog/pkg/rndm"
	"strings"
)

// Log 操作日志
// 只记录敏感操作（写接口 + 登录），由中间件自动写入，不提供手工新增。
// Data 存的是脱敏后的请求参数，密码类字段一律替换为 ***。
type Log struct {
	Id      string `xorm:"VARCHAR(64) PK" json:"id"`           //主键
	UserId  string `xorm:"INDEX VARCHAR(64)" json:"user_id"`   //操作人id
	User    *User  `xorm:"-" swaggerignore:"true" json:"user"` //操作人-非持久化
	Module  string `xorm:"INDEX VARCHAR(128)" json:"module"`   //模块（分组中文名）
	Action  string `xorm:"VARCHAR(128)" json:"action"`         //动作（权限中文名）
	Perm    string `xorm:"INDEX VARCHAR(128)" json:"perm"`     //权限标识
	Method  string `xorm:"VARCHAR(16)" json:"method"`          //请求方法
	Path    string `xorm:"VARCHAR(255)" json:"path"`           //请求路径
	Ip      string `xorm:"VARCHAR(64)" json:"ip"`              //来源ip
	Status  int    `xorm:"INT(11)" json:"status"`              //业务码 200成功 其它失败
	Data    string `xorm:"TEXT" json:"data"`                   //请求参数（已脱敏）
	Created int64  `xorm:"BIGINT INDEX" json:"created"`        //创建时间
}

func (Log) TableName() string {
	return "log"
}

// BeforeInsert 插入前生成主键
func (s *Log) BeforeInsert() {
	if s.Id == "" {
		s.Id = rndm.ULID()
	}
}

// IsPkZero 主键是否为空
func (s Log) IsPkZero() bool {
	return s.Id == ""
}

// GetId 返回主键
func (s Log) GetId() string {
	return s.Id
}

// LogSensitiveFields 请求参数中需要脱敏的字段名
// 命中即整体替换为 ***，不保留任何原文片段
var LogSensitiveFields = []string{
	"passwd", "password", "opasswd", "npasswd", "old_password", "new_password",
	"secret", "token", "vcode", "vreal", "private_key", "access_key",
}

// LogMaskedValue 脱敏后的占位值
const LogMaskedValue = "***"

// LogNeedMask 字段是否需要脱敏
func LogNeedMask(key string) bool {
	k := strings.ToLower(strings.TrimSpace(key))
	for _, it := range LogSensitiveFields {
		if k == it {
			return true
		}
	}
	return false
}

// LogFilter 操作日志过滤条件
type LogFilter struct {
	UserId *string `query:"user_id" form:"user_id" json:"user_id"` //操作人
	Module string  `query:"module" form:"module" json:"module"`    //模块
	Perm   string  `query:"perm" form:"perm" json:"perm"`          //权限标识
	Method string  `query:"method" form:"method" json:"method"`    //请求方法
	Status *int    `query:"status" form:"status" json:"status"`    //业务码
	From   int64   `query:"from" form:"from" json:"from"`          //起始时间（毫秒）
	To     int64   `query:"to" form:"to" json:"to"`                //结束时间（毫秒）
}

// LogListArgs 操作日志列表参数
type LogListArgs struct {
	List
	LogFilter
}
