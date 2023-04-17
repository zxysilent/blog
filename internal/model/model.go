package model

import (
	"blog/conf"
	"errors"
	"strings"
)

type IptId struct {
	Id int `query:"id" form:"id" json:"id"` //仅包含Id的请求
}

// Page 分页基本数据
type Page struct {
	Pi   int    `json:"pi" form:"pi" query:"pi"`       //分页页码
	Ps   int    `json:"ps" form:"ps" query:"ps"`       //分页大小
	Mult string `json:"mult" form:"mult" query:"mult"` //多条件信息
}

func (p *Page) SafeMult() string {
	return safeMult(p.Mult)
}

// Stat 检查状态
func (p *Page) Stat() error {
	if p.Ps < conf.App.PageMin {
		return errors.New("page size 过小")
	}
	if p.Ps > conf.App.PageMax {
		return errors.New("page size 过大")
	}
	return nil
}

const (
	StatusEnable  = 1 //启用
	StatusDisable = 2 //禁用
)

const Zero = 0 //忽略
const All = -1 //所有

type List struct {
	Mult  string `form:"mult" query:"mult" json:"mult"`    //多条件信息
	Limit int    `form:"limit" query:"limit" json:"limit"` //条数限制(<=0不限制)
}

func (l *List) SafeMult() string {
	return safeMult(l.Mult)
}

func safeMult(s string) string {
	s = strings.TrimSpace(s)
	return "%" + s + "%"
}

// Naver 上下页
type Naver struct {
	Prev string
	Next string
}
