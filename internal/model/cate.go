package model

import (
	"blog/pkg/rndm"
	"regexp"
	"strings"
)

// CateMaxLevel 栏目最大层级
// 层级由服务端按父链自动维护，创建/移动时不允许突破该上限
const CateMaxLevel = 3

// CatePathReg 公开访问路径格式
// 小写字母或数字开头，只含小写字母、数字、下划线、中划线，长度 1~64
var CatePathReg = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]{0,63}$`)

// Cate 栏目 / 分类
// 通用后台的导航与内容容器，同时服务于 blog、wiki、门户等场景：
//   - 树形结构：Pid 自关联，最多 CateMaxLevel 级
//   - 公开路径：Path 全站唯一，前台按路径访问，改名不影响已有链接
//   - 外链跳转：Url 非空时直接 302，不渲染列表
//   - 访问控制：Auth 为真时需登录才能访问
//   - 模板覆盖：Itpl/Ltpl/Dtpl 分别覆盖栏目首页、列表页、详情页模板
type Cate struct {
	Id       string `xorm:"VARCHAR(64) PK" json:"id"`         //主键
	Pid      string `xorm:"INDEX VARCHAR(64)" json:"pid"`     //父级id，空串为顶层
	Kind     int    `xorm:"INT(11) DEFAULT 1" json:"kind"`    //归属内容类型 1文章 2页面 3笔记
	Name     string `xorm:"INDEX VARCHAR(255)" json:"name"`   //名称，同级唯一
	Path     string `xorm:"INDEX VARCHAR(64)" json:"path"`    //公开访问路径，全站唯一
	Url      string `xorm:"VARCHAR(255)" json:"url"`          //跳转地址，非空则直接跳转
	Intro    string `xorm:"VARCHAR(500)" json:"intro"`        //描述
	Icon     string `xorm:"VARCHAR(255)" json:"icon"`         //图标
	Color    string `xorm:"VARCHAR(64)" json:"color"`         //颜色
	Sort     int    `xorm:"INT(11) DEFAULT 1000" json:"sort"` //排序，越小越前
	Level    int    `xorm:"INT(11) DEFAULT 1" json:"level"`   //层级，服务端维护
	Use      bool   `xorm:"BOOL DEFAULT 1" json:"use"`        //可被内容选中
	Show     bool   `xorm:"BOOL DEFAULT 1" json:"show"`       //前台导航显示
	Inner    bool   `xorm:"BOOL DEFAULT 0" json:"inner"`      //内置，禁止删除
	Auth     bool   `xorm:"BOOL DEFAULT 0" json:"auth"`       //需登录才能访问
	Itpl     string `xorm:"VARCHAR(255)" json:"itpl"`         //栏目首页模板
	Ltpl     string `xorm:"VARCHAR(255)" json:"ltpl"`         //列表模板
	Dtpl     string `xorm:"VARCHAR(255)" json:"dtpl"`         //详情模板
	Updated  int64  `xorm:"BIGINT" json:"updated"`            //修改时间
	Created  int64  `xorm:"BIGINT" json:"created"`            //创建时间
	Count    int    `xorm:"-" json:"count"`                   //内容数量，仅列表填充
	Children []Cate `xorm:"-" json:"children"`                //子级，树形返回
}

func (Cate) TableName() string {
	return "cate"
}

// BeforeInsert 插入前生成主键
func (s *Cate) BeforeInsert() {
	if s.Id == "" {
		s.Id = rndm.ULID()
	}
}

// IsPkZero 主键是否为空
func (s Cate) IsPkZero() bool {
	return s.Id == ""
}

// GetId 返回主键
func (s Cate) GetId() string {
	return s.Id
}

// CateNormalizePath 规范化公开路径：去空格并转小写
// 仅做规范化，格式校验见 CatePathValid
func CateNormalizePath(path string) string {
	return strings.ToLower(strings.TrimSpace(path))
}

// CatePathValid 公开路径是否合法
func CatePathValid(path string) bool {
	return CatePathReg.MatchString(path)
}

// CateFilter 栏目过滤条件
type CateFilter struct {
	Pid  *string `query:"pid" form:"pid" json:"pid"`    //父级id
	Kind *int    `query:"kind" form:"kind" json:"kind"` //归属内容类型
	Use  *bool   `query:"use" form:"use" json:"use"`    //仅可选中
	Show *bool   `query:"show" form:"show" json:"show"` //仅前台显示
}

// CateListArgs 栏目分页列表参数
type CateListArgs struct {
	List
	CateFilter
}

// CateTreeArgs 栏目树参数
type CateTreeArgs struct {
	CateFilter
	Root bool `query:"root" form:"root" json:"root"` //顶部附加「全部」节点
	Flat bool `query:"flat" form:"flat" json:"flat"` //扁平返回（带 level，便于下拉缩进）
}

// CateAddArgs 栏目新增参数
// 布尔字段用指针，未传时按「可用、显示」处理，而不是被零值覆盖
type CateAddArgs struct {
	Pid   string `json:"pid"`   //父级id，空串为顶层
	Kind  int    `json:"kind"`  //归属内容类型，0 时按文章
	Name  string `json:"name"`  //名称，同级唯一
	Path  string `json:"path"`  //公开路径，留空自动生成
	Url   string `json:"url"`   //跳转地址
	Intro string `json:"intro"` //描述
	Icon  string `json:"icon"`  //图标
	Color string `json:"color"` //颜色
	Sort  int    `json:"sort"`  //排序
	Use   *bool  `json:"use"`   //可被内容选中
	Show  *bool  `json:"show"`  //前台导航显示
	Auth  *bool  `json:"auth"`  //需登录访问
	Itpl  string `json:"itpl"`  //栏目首页模板
	Ltpl  string `json:"ltpl"`  //列表模板
	Dtpl  string `json:"dtpl"`  //详情模板
}

// CateEditArgs 栏目修改参数
// 全部为指针，nil 表示不修改，支持单字段局部更新
type CateEditArgs struct {
	Pid   *string `json:"pid"`   //父级id，传空串表示移到顶层
	Kind  *int    `json:"kind"`  //归属内容类型
	Name  *string `json:"name"`  //名称，同级唯一
	Path  *string `json:"path"`  //公开路径，全站唯一
	Url   *string `json:"url"`   //跳转地址
	Intro *string `json:"intro"` //描述
	Icon  *string `json:"icon"`  //图标
	Color *string `json:"color"` //颜色
	Sort  *int    `json:"sort"`  //排序
	Use   *bool   `json:"use"`   //可被内容选中
	Show  *bool   `json:"show"`  //前台导航显示
	Auth  *bool   `json:"auth"`  //需登录访问
	Itpl  *string `json:"itpl"`  //栏目首页模板
	Ltpl  *string `json:"ltpl"`  //列表模板
	Dtpl  *string `json:"dtpl"`  //详情模板
}
