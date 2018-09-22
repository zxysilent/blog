package model

import (
	"fmt"
)

// Opts 配置
type Opts struct {
	Key   string `xorm:"not null pk default '''' unique VARCHAR(64)"`
	Value string `xorm:"default 'NULL' VARCHAR(1024)"`
	Desc  string `xorm:"default 'NULL' VARCHAR(255)"`
}

type mapOpts map[string]string

func (m mapOpts) Get(key string) string {
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
	fmt.Println(MapOpts.Get("site_url"))
}
