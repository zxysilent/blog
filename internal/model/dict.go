package model

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"
)

// Dict 配置
type Dict struct {
	Key      string `xorm:"PK UNIQUE VARCHAR(64)" json:"key"`      //唯一Key
	Name     string `xorm:"VARCHAR(255)" json:"name"`              //显示
	Value    string `xorm:"TEXT" json:"value"`                     //值
	Min      int    `xorm:"INT(11)" json:"min"`                    //数字最小值，字符串最小长度
	Max      int    `xorm:"INT(11)" json:"max"`                    //数字最大值，字符串最大长度
	Describe string `xorm:"VARCHAR(255)" json:"describe"`          //描述
	Editable bool   `xorm:"TINYINT(4) DEFAULT(0)" json:"editable"` //是否允许修改
	Sort     int    `xorm:"INT(11)" json:"sort"`                   //显示排序
	Type     string `xorm:"VARCHAR(64)" json:"type"`               //类型 JSON,INT,TIME，BOOL,STRING,FLOAT
	Inner    bool   `xorm:"TINYINT(4) DEFAULT 0" json:"inner"`     //内部-禁止删除
	Kind     string `xorm:"VARCHAR(64)" json:"kind"`               //种类
	Intro    string `xorm:"VARCHAR(255)" json:"intro"`             //说明
	Updated  int64  `xorm:"BIGINT" json:"updated"`                 //修改时间
	Created  int64  `xorm:"BIGINT" json:"created"`                 //创建时间
}

func (Dict) TableName() string {
	return "dict"
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

const (
	DictKeyLogin = "login" //登录页面
	DictKeyBasic = "basic" //基础配置
	DictKeyBlog  = "blog"  //应用首页
)

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

func (d *Dict) json(v any) error {
	return json.Unmarshal([]byte(d.Value), v)
}

func (d *Dict) ToBasic() *Basic {
	mod := &Basic{}
	d.json(mod)
	return mod
}

func (d *Dict) ToBlog() *Blog {
	mod := &Blog{}
	d.json(mod)
	return mod
}

func (d *Dict) Check() error {
	switch d.Type {
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

type DictListArgs struct {
	List
	Name string `form:"name" query:"name" json:"name"`
}

type DictEditArgs struct {
	Value    *string `json:"value"`
	Describe *string `json:"describe"`
}

type DictSet struct {
	set   map[string]*Dict `json:"-"`
	Blog  *Blog
	Basic *Basic
}

func (cs *DictSet) Get(id string) *Dict {
	cfg, ok := cs.set[id]
	if !ok {
		cfg = &Dict{Key: id}
	}
	return cfg
}

func NewDictSet() *DictSet {
	return &DictSet{
		set: make(map[string]*Dict),
	}
}

func (cs *DictSet) Set(key string, cfg *Dict) {
	cs.set[key] = cfg
}
