package model

import (
	"strings"
)

// Page 分页基本数据
type List struct {
	Pi   int    `json:"pi" form:"pi" query:"pi"`       //分页页码
	Ps   int    `json:"ps" form:"ps" query:"ps"`       //分页大小
	Mult string `json:"mult" form:"mult" query:"mult"` //多条件信息
}

func (p *List) SafeMult() string {
	return safeMult(p.Mult)
}

func (l *List) GetLimit() int {
	if l.Ps < 1 {
		l.Ps = 1
	}
	return l.Ps
}

func (l *List) GetOffset() int {
	if l.Ps < 1 {
		l.Ps = 1
	}
	if l.Pi < 1 {
		l.Pi = 1
	}
	return (l.Pi - 1) * l.Ps
}

func (l *List) GetCap() int {
	if l.Ps >= 32 {
		return 32
	}
	if l.Ps < 8 {
		return 8
	}
	return l.Ps
}

func (l *List) IsZero() bool {
	return l.Pi == 0 && l.Ps == 0
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
