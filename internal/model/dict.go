package model

import (
	"blog/pkg/rndm"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
)

// 字典值类型
const (
	DictTypeString = "string" //字符串
	DictTypeBool   = "bool"   //布尔
	DictTypeInt    = "int"    //整数
	DictTypeTime   = "time"   //时间
)

// 内置字典主键（主键即业务键，代码中引用可避免拼写错误）
const (
	DictIdGender   = "gender"    //性别
	DictIdPageSize = "page_size" //分页大小，列表页可选的分页条数
)

var (
	dictTimeFmts = []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02",
		"15:04:05",
	}
)

// Dict 字典（字典类型主表）
// 主键 Id 同时作为唯一业务键：创建时允许自定义（如 gender），
// 未指定时由 BeforeInsert 生成 ULID
// 字典项存放在独立的 dict_entry 表，Entries 仅为读写时组装的结果，不持久化
// Inner=true 为内置字典，禁止删除
type Dict struct {
	Id      string      `xorm:"VARCHAR(64) PK" json:"id"`     //主键，兼作字典键
	Name    string      `xorm:"VARCHAR(128)" json:"name"`     //字典名称
	Type    string      `xorm:"VARCHAR(16)" json:"type"`      //值类型 string/bool/int/time
	Status  bool        `xorm:"BOOL DEFAULT 1" json:"status"` //是否启用
	Inner   bool        `xorm:"BOOL DEFAULT 0" json:"inner"`  //内置-禁止删除
	Remark  string      `xorm:"VARCHAR(255)" json:"remark"`   //备注
	Entries []DictEntry `xorm:"-" json:"entries"`             //选项集合-非持久化
	Created int64       `xorm:"BIGINT" json:"created"`        //创建时间
	Updated int64       `xorm:"BIGINT" json:"updated"`        //修改时间
}

func (Dict) TableName() string {
	return "dict"
}

// BeforeInsert 插入前生成主键
func (s *Dict) BeforeInsert() {
	if s.Id == "" {
		s.Id = rndm.ULID()
	}
}

// IsPkZero 主键是否为空
func (s Dict) IsPkZero() bool {
	return s.Id == ""
}

// GetId 返回主键
func (s Dict) GetId() string {
	return s.Id
}

// DictEntry 字典项（独立表）
// dict_id 关联字典，pid 表达层级（支持树形字典）
// 唯一索引 (dict_id, pid, field) 保证同一层级字段名不重复
type DictEntry struct {
	Id       string      `xorm:"VARCHAR(64) PK" json:"id"`                            //主键
	DictId   string      `xorm:"UNIQUE(dict_entry) INDEX VARCHAR(64)" json:"dict_id"` //所属字典
	Pid      string      `xorm:"UNIQUE(dict_entry) VARCHAR(64)" json:"pid"`           //父项id，空为顶层
	Field    string      `xorm:"UNIQUE(dict_entry) VARCHAR(64)" json:"field"`         //字段名，同层唯一
	Label    string      `xorm:"VARCHAR(128)" json:"label"`                           //展示名
	Value    string      `xorm:"VARCHAR(255)" json:"value"`                           //字典值
	Sort     int         `xorm:"INT DEFAULT 1000" json:"sort"`                        //排序
	Status   bool        `xorm:"BOOL DEFAULT 1" json:"status"`                        //是否启用
	Created  int64       `xorm:"BIGINT" json:"created"`                               //创建时间
	Updated  int64       `xorm:"BIGINT" json:"updated"`                               //修改时间
	Type     string      `xorm:"-" json:"type"`                                       //值类型，由所属字典下发-非持久化
	Children []DictEntry `xorm:"-" json:"children,omitempty"`                         //子级-非持久化
}

func (DictEntry) TableName() string {
	return "dict_entry"
}

// BeforeInsert 插入前生成主键
func (s *DictEntry) BeforeInsert() {
	if s.Id == "" {
		s.Id = rndm.ULID()
	}
}

// IsPkZero 主键是否为空
func (s DictEntry) IsPkZero() bool {
	return s.Id == ""
}

// GetId 返回主键
func (s DictEntry) GetId() string {
	return s.Id
}

// Valid 校验字典主表字段
func (d *Dict) Valid() error {
	d.Id = strings.TrimSpace(d.Id)
	d.Name = strings.TrimSpace(d.Name)
	d.Remark = strings.TrimSpace(d.Remark)
	d.Type = strings.ToLower(strings.TrimSpace(d.Type))
	if d.Type == "" {
		d.Type = DictTypeString
	}
	if err := ValidateBizId(d.Id, "字典主键"); err != nil {
		return err
	}
	if d.Name == "" || len(d.Name) > 128 {
		return errors.New("字典名称不能为空且长度不超过 128")
	}
	switch d.Type {
	case DictTypeString, DictTypeBool, DictTypeInt, DictTypeTime:
	default:
		return fmt.Errorf("不支持的字典类型 %s", d.Type)
	}
	if len(d.Remark) > 255 {
		return errors.New("备注长度不能超过 255")
	}
	return nil
}

// BuildTree 把扁平选项列表组装成树填入 Entries，按 sort、field 升序
func (d *Dict) BuildTree(rows []DictEntry) {
	if len(rows) < 1 {
		d.Entries = []DictEntry{}
		return
	}
	group := make(map[string][]DictEntry, len(rows))
	for idx := range rows {
		it := rows[idx]
		it.Type = d.Type
		it.Children = nil
		group[it.Pid] = append(group[it.Pid], it)
	}
	var build func(pid string, depth int) []DictEntry
	build = func(pid string, depth int) []DictEntry {
		if depth > 16 { // 防御脏数据造成的环
			return nil
		}
		items, ok := group[pid]
		if !ok {
			return nil
		}
		sort.Slice(items, func(i, j int) bool {
			if items[i].Sort != items[j].Sort {
				return items[i].Sort < items[j].Sort
			}
			return items[i].Field < items[j].Field
		})
		out := make([]DictEntry, len(items))
		copy(out, items)
		for i := range out {
			out[i].Children = build(out[i].Id, depth+1)
		}
		return out
	}
	d.Entries = build("", 0)
}

// Flatten 展开为平铺选项列表（过滤禁用项，清空 children）
// 字典整体禁用时返回空列表
func (d *Dict) Flatten() []DictEntry {
	if !d.Status {
		return []DictEntry{}
	}
	out := make([]DictEntry, 0, len(d.Entries))
	var walk func(entries []DictEntry)
	walk = func(entries []DictEntry) {
		for idx := range entries {
			it := entries[idx]
			if !it.Status {
				continue
			}
			item := it
			item.Children = nil
			out = append(out, item)
			walk(it.Children)
		}
	}
	walk(d.Entries)
	return out
}

// FlatDictEntries 校验树形选项并展开为可直接入库的扁平列表
// 自动补齐主键、dict_id、pid；同一层级 field 不可重复
func FlatDictEntries(dictId, typ string, entries []DictEntry) ([]DictEntry, error) {
	out := make([]DictEntry, 0, len(entries))
	var walk func(pid string, items []DictEntry) error
	walk = func(pid string, items []DictEntry) error {
		seen := make(map[string]struct{}, len(items))
		for idx := range items {
			it := items[idx]
			it.Field = strings.TrimSpace(it.Field)
			it.Label = strings.TrimSpace(it.Label)
			it.Value = strings.TrimSpace(it.Value)
			if it.Field == "" {
				return errors.New("选项字段名不能为空")
			}
			if len(it.Field) > 64 {
				return fmt.Errorf("选项字段名 %s 长度不能超过 64", it.Field)
			}
			if _, has := seen[it.Field]; has {
				return fmt.Errorf("同一层级选项字段名重复 %s", it.Field)
			}
			seen[it.Field] = struct{}{}
			if it.Label == "" {
				it.Label = it.Value
			}
			if len(it.Label) > 128 {
				return fmt.Errorf("选项 %s 的展示名长度不能超过 128", it.Field)
			}
			if len(it.Value) > 255 {
				return fmt.Errorf("选项 %s 的值长度不能超过 255", it.Field)
			}
			if err := ValidateDictValue(typ, it.Value); err != nil {
				return fmt.Errorf("选项 %s 的值非法: %w", it.Field, err)
			}
			id := it.Id
			if id == "" {
				id = rndm.ULID()
			}
			out = append(out, DictEntry{
				Id:      id,
				DictId:  dictId,
				Pid:     pid,
				Field:   it.Field,
				Label:   it.Label,
				Value:   it.Value,
				Sort:    it.Sort,
				Status:  it.Status,
				Created: it.Created,
				Updated: it.Updated,
			})
			if err := walk(id, it.Children); err != nil {
				return err
			}
		}
		return nil
	}
	if err := walk("", entries); err != nil {
		return nil, err
	}
	return out, nil
}

// ValidateDictValue 按类型校验字典值
func ValidateDictValue(typ, val string) error {
	switch typ {
	case DictTypeString:
		return nil
	case DictTypeBool:
		if val != "true" && val != "false" {
			return errors.New("布尔值只能是 true 或 false")
		}
		return nil
	case DictTypeInt:
		if _, err := strconv.ParseInt(val, 10, 64); err != nil {
			return errors.New("必须是整数")
		}
		return nil
	case DictTypeTime:
		for _, fmtStr := range dictTimeFmts {
			if _, err := time.Parse(fmtStr, val); err == nil {
				return nil
			}
		}
		return errors.New("时间格式不合法")
	}
	return nil
}

// DictListArgs 字典列表参数
type DictListArgs struct {
	List
	Status *bool `query:"status" form:"status" json:"status"` //状态过滤
}

// DictAddArgs 字典新增参数
type DictAddArgs struct {
	Id      string      `json:"id" form:"id"`           //主键，可自定义，留空自动生成
	Name    string      `json:"name" form:"name"`       //字典名称
	Type    string      `json:"type" form:"type"`       //值类型
	Status  *bool       `json:"status" form:"status"`   //是否启用
	Remark  string      `json:"remark" form:"remark"`   //备注
	Entries []DictEntry `json:"entries" form:"entries"` //选项集合
}

// DictEditArgs 字典修改参数
// Entries 为 nil 表示不修改选项，指向空切片表示清空选项
type DictEditArgs struct {
	Name    *string      `json:"name" form:"name"`       //字典名称
	Type    *string      `json:"type" form:"type"`       //值类型
	Status  *bool        `json:"status" form:"status"`   //是否启用
	Remark  *string      `json:"remark" form:"remark"`   //备注
	Entries *[]DictEntry `json:"entries" form:"entries"` //选项集合
}
