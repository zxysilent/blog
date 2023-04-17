package model

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

// Dict 配置
type Dict struct {
	Key     string `xorm:"PK UNIQUE VARCHAR(64)" json:"key"`  //唯一Key
	Value   string `xorm:"VARCHAR(2048)" json:"value"`        //值
	Kind    string `xorm:"VARCHAR(255)" json:"kind"`          //类型 JSON,INT，TIME，BOOL,STRING
	Inner   bool   `xorm:"TINYINT(4) DEFAULT 0" json:"inner"` //内部-禁止删除
	Intro   string `xorm:"VARCHAR(255)" json:"intro"`         //说明
	Updated int64  `xorm:"BIGINT" json:"updated"`             //修改时间
	Created int64  `xorm:"BIGINT" json:"created"`             //创建时间
}

func (Dict) TableName() string {
	return sysPrefix + "dict"
}

// ------------------------------------------------------ Basic 基础配置 ------------------------------------------------------
// Basic 基础配置
type Basic struct {
	Title       string `json:"title"`       //标题
	SiteUrl     string `json:"site_url"`    //网站地址
	LogoUrl     string `json:"logo_url"`    //Logo地址
	Author      string `json:"author"`      //网站作者
	Copyright   string `json:"copyright"`   //版权
	Keywords    string `json:"keywords"`    //网站关键词
	Description string `json:"description"` //描述
	BeianMiit   string `json:"beian_miit"`  //ICP备案
	BeianNism   string `json:"beian_nism"`  //公安备案

}
type Blog struct {
	PageSize  int    `json:"page_size"`  //分页大小
	Analytic  string `json:"analytic"`   //统计脚本
	Comment   string `json:"comment"`    //评论脚本
	GithubUrl string `json:"github_url"` //githu地址
	WeiboUrl  string `json:"weibo_url"`  //微博地址
}

const DictKeyLogin = "login" //登录页面
const DictKeyBasic = "basic" //基础配置
const DictKeyBlog = "blog"   //应用首页

func (d *Dict) Int() int {
	v, _ := d.int()
	return v
}

func (d *Dict) int() (int, error) {
	return strconv.Atoi(d.Value)
}

func (d *Dict) Time() time.Time {
	v, _ := d.time()
	return v
}

func (d *Dict) time() (time.Time, error) {
	v, err := strconv.ParseInt(d.Value, 10, 64)
	if err != nil {
		return time.Time{}, err
	}
	return time.UnixMilli(v), nil
}

func (d *Dict) Bool() bool {
	v, _ := d.bool()
	return v
}

func (d *Dict) bool() (bool, error) {
	return strconv.ParseBool(d.Value)
}

func (d *Dict) json(v interface{}) error {
	return json.Unmarshal([]byte(d.Value), v)
}

func (d *Dict) ToBasic() *Basic {
	mod := &Basic{}
	json.Unmarshal([]byte(d.Value), mod)
	return mod
}
func (d *Dict) ToBlog() *Blog {
	mod := &Blog{}
	json.Unmarshal([]byte(d.Value), mod)
	return mod
}

func (d *Dict) Check() error {
	switch d.Kind {
	case "int":
		_, err := d.int()
		return err
	case "bool":
		_, err := d.bool()
		return err
	case "time":
		_, err := d.time()
		return err
	case "json":
		if !json.Valid([]byte(d.Value)) {
			return errors.New("invalid json")
		}
	default:
	}
	return nil
}
