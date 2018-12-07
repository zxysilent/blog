package model

// Opts 配置
type Opts struct {
	Key   string `xorm:"pk unique VARCHAR(64)" json:"key" form:"key"`
	Value string `xorm:"VARCHAR(1024)" json:"value" form:"value"`
	Desc  string `xorm:"VARCHAR(255)" json:"desc" form:"desc"`
}

type mapOpts map[string]string

func (m mapOpts) Get(key string) (string, bool) {
	value, ok := m[key]
	return value, ok
}
func (m mapOpts) Set(key, val string) {
	if m[key] != val {
		m[key] = val
	}
}

func (m mapOpts) MustGet(key string) string {
	return m[key]
}

//MapOpts 数据库配置信息
var MapOpts mapOpts

func initMap() {
	mods := make([]Opts, 0, 8)
	DB.Find(&mods)
	m := make(map[string]string, len(mods))
	for _, v := range mods {
		m[v.Key] = v.Value
	}
	MapOpts = m
}

func init() {
	initMap()
}

// OptsGet 获取某个 配置
func OptsGet(key string) (string, bool) {
	return MapOpts.Get(key)
}

// OptsMustGet 获取某个 配置
func OptsMustGet(key string) string {
	return MapOpts.MustGet(key)
}

// OptsEdit 编辑配置
func OptsEdit(mod *Opts) bool {
	sess := DB.NewSession()
	defer sess.Close()
	sess.Begin()
	affect, err := sess.ID(mod.Key).Cols("Value").Update(mod)
	if affect >= 0 && err == nil {
		sess.Commit()
		//	initMap()
		MapOpts.Set(mod.Key, mod.Value)
		return true
	}
	sess.Rollback()
	return false
}
