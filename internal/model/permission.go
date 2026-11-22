package model

import "blog/pkg/rndm"

// Perm 权限点
// 只读：完全由路由注册汇总生成，不提供手工增删改
// 表名保持 permission，对外接口路径保持 /permissions，
// 其余内部标识统一使用 perm 缩写。
type Perm struct {
	Id       string `xorm:"VARCHAR(64) PK" json:"id"`           //主键
	Perm     string `xorm:"UNIQUE VARCHAR(128)" json:"perm"`    //权限标识 资源:动作 如 user:create
	Name     string `xorm:"VARCHAR(255)" json:"name"`           //权限名称
	GroupKey string `xorm:"VARCHAR(64) INDEX" json:"group_key"` //分组键 如 user，用于推导分组授权权限
	Group    string `xorm:"VARCHAR(255)" json:"group"`          //权限分组中文名
	Sort     int    `xorm:"INT DEFAULT 1000" json:"sort"`       //排序
	Created  int64  `xorm:"BIGINT" json:"created"`              //创建时间
}

func (Perm) TableName() string {
	return "permission"
}

// BeforeInsert 插入前生成主键
// 权限点完全由路由派生，主键必须由 perm 稳定推导：
// 同一个权限标识在任何库、任何时刻都得到同一个 id，
// 这样角色授权（role_permission.perm_id）不会因重建库、换环境而失效。
func (s *Perm) BeforeInsert() {
	if s.Id == "" {
		s.Id = PermId(s.Perm)
	}
}

// PermId 由权限标识推导稳定主键
// perm 唯一且非空，因此不存在歧义
func PermId(perm string) string {
	return rndm.ULIDFrom("permission:" + perm)
}

// IsPkZero 主键是否为空
func (s Perm) IsPkZero() bool {
	return s.Id == ""
}

// GetId 返回主键
func (s Perm) GetId() string {
	return s.Id
}

// PermNode 权限树节点
// Perm 为空表示分组节点，非空表示可授权的权限叶子
type PermNode struct {
	Id       string     `json:"id,omitempty"`       //权限id，分组节点为空
	Name     string     `json:"name"`               //中文名称
	GroupKey string     `json:"group_key"`          //分组键
	Group    string     `json:"group"`              //分组名称
	Perm     string     `json:"perm,omitempty"`     //权限标识，分组节点为空
	Sort     int        `json:"sort"`               //排序
	Children []PermNode `json:"children,omitempty"` //子节点
}

// PermListArgs 权限列表参数
type PermListArgs struct {
	List
	Group string `query:"group" form:"group" json:"group"` //分组过滤
}

// PermInitArgs 权限初始化参数，由路由汇总后启动时幂等灌入
type PermInitArgs struct {
	Perm     string //权限标识
	Name     string //中文名称
	GroupKey string //分组键
	Group    string //中文分组
	Sort     int    //排序
}

// GroupGrantPerm 分组授权权限标识
// 修改某个分组下的任意权限，操作者必须拥有该分组的授权权限
// 即「每个分组都有一个允许授权的权限」
func GroupGrantPerm(groupKey string) string {
	return groupKey + ":grant"
}
