package model

import (
	"fmt"
	"regexp"
	"strings"
)

// 分页参数默认值
// 所有列表接口的分页缺省值统一在这里定义，调用方不再各自写 if 判空
const (
	DefaultPage     = 1  //默认页码
	DefaultPageSize = 32 //默认每页条数
)

// List 分页基本数据
// 缺省值由 Normalize / GetLimit / GetOffset 统一处理，服务层无需重复校验
// 注意：Get* 系列不修改接收者，只有 Normalize 会写回字段
type List struct {
	Pi   int    `json:"page" form:"page" query:"page"` //分页页码
	Ps   int    `json:"size" form:"size" query:"size"` //分页大小
	Mult string `json:"q" form:"q" query:"q"`          //模糊查询关键字
}

func (p *List) SafeMult() string {
	return safeMult(p.Mult)
}

// Normalize 归一化分页参数，把未传或非法的值回落为默认值
// ps 指定默认每页条数，<=0 时用 DefaultPageSize
// 已显式传入的合法值不会被改动
// 归一化后 IsZero 恒为 false，即一定走分页查询
func (l *List) Normalize(ps int) *List {
	if ps <= 0 {
		ps = DefaultPageSize
	}
	if l.Pi < DefaultPage {
		l.Pi = DefaultPage
	}
	if l.Ps < 1 {
		l.Ps = ps
	}
	return l
}

// GetLimit 每页条数，未传或非法时取 DefaultPageSize
func (l *List) GetLimit() int {
	if l.Ps < 1 {
		return DefaultPageSize
	}
	return l.Ps
}

// GetOffset 偏移量，由页码与每页条数推导
func (l *List) GetOffset() int {
	pi := l.Pi
	if pi < DefaultPage {
		pi = DefaultPage
	}
	return (pi - DefaultPage) * l.GetLimit()
}

// GetCap 预分配容量，控制在 [8, DefaultPageSize]
func (l *List) GetCap() int {
	ps := l.GetLimit()
	if ps < 8 {
		return 8
	}
	if ps > DefaultPageSize {
		return DefaultPageSize
	}
	return ps
}

// IsZero 未传任何分页参数：此时列表查询不分页，返回全部
func (l *List) IsZero() bool {
	return l.Pi == 0 && l.Ps == 0
}

func safeMult(s string) string {
	s = strings.TrimSpace(s)
	return "%" + s + "%"
}

// bizIdReg 业务主键格式：允许指定，字母或数字开头，只含字母数字、下划线、中划线
// 未指定时由 BeforeInsert 生成 ULID（小写 base32，首字符可能是数字 0-7）
var bizIdReg = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9_-]*$`)

// ValidateBizId 校验可指定的业务主键
// 用于 dict / config 这类「主键即业务键」的表：创建时允许自定义 id
func ValidateBizId(id string, label string) error {
	if id == "" {
		return fmt.Errorf("%s不能为空", label)
	}
	if len(id) > 64 {
		return fmt.Errorf("%s长度不能超过 64", label)
	}
	if !bizIdReg.MatchString(id) {
		return fmt.Errorf("%s只能包含字母、数字、下划线和短横线，且不能以符号开头", label)
	}
	return nil
}

// Naver 上下页
type Naver struct {
	Prev string
	Next string
}
