package model

import (
	"errors"
	"time"

	"github.com/zxysilent/logs"
)

// Goinfo go information
type Goinfo struct {
	ARCH    string `json:"arch"`
	OS      string `json:"os"`
	Version string `json:"version"`
	NumCPU  int    `json:"num_cpu"`
}

// ------------------------------------------------------ Global 全局配置 ------------------------------------------------------
// 结构体配置

// Global 全局配置
type Global struct {
	Id int `xorm:"INT(11) PK AUTOINCR comment('主键')" json:"id"`
}

const globalId = 1

func (Global) TableName() string {
	return "sys_global"
}

var globalCache Global

func initGlobal() error {
	mod := Global{}
	has, _ := Db.ID(1).Get(&mod)
	if !has {
		return errors.New("no")
	}
	globalCache = mod
	logs.Debug("global cache")
	return nil
}

func GlobalGet() (*Global, bool) {
	// mod := &Global{}
	// has, _ := Db.ID(globalId).Get(mod)
	return &globalCache, true
}

// GlobalEdit 编辑global信息
func GlobalEdit(mod *Global, cols ...string) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(globalId).Cols(cols...).Update(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	initGlobal()
	return nil
}

// ------------------------------------------------------ Dict 配置 ------------------------------------------------------
// 字典配置

// Dict 配置
type Dict struct {
	Key   string    `xorm:"PK UNIQUE VARCHAR(64) comment('唯一Key')" json:"key"`
	Value string    `xorm:"VARCHAR(255) comment('值')" json:"value"`
	Inner bool      `xorm:"TINYINT(4) DEFAULT 0 comment('内部禁止删除')" json:"inner"`
	Intro string    `xorm:"VARCHAR(255) comment('说明')" json:"intro"`
	Ctime time.Time `xorm:"DATETIME comment('时间')" swaggerignore:"true" json:"ctime"`
}

func (Dict) TableName() string {
	return "sys_dict"
}

var dictCache map[string]*Dict

func initDict() error {
	mods := make([]Dict, 0, 8)
	err := Db.Find(&mods)
	if err != nil {
		return err
	}
	dictMap := make(map[string]*Dict, 8)
	for idx := range mods {
		mod := mods[idx]
		dictMap[mod.Key] = &mod

	}
	dictCache = dictMap
	logs.Debug("Dict cache")
	return nil
}
func DictCache() error {
	return initDict()
}

// DictGet 单条字典
func DictGet(key string) (*Dict, bool) {
	// mod := &Dict{}
	// has, _ := Db.ID(key).Get(mod)
	mod, has := dictCache[key]
	return mod, has
}

// DictAll 所有字典
func DictAll() ([]Dict, error) {
	mods := make([]Dict, 0, 8)
	err := Db.Find(&mods)
	return mods, err
}

// DictPage 字典分页
func DictPage(pi int, ps int, cols ...string) ([]Dict, error) {
	mods := make([]Dict, 0, ps)
	sess := Db.NewSession()
	defer sess.Close()
	if len(cols) > 0 {
		sess.Cols(cols...)
	}
	err := sess.Desc("Ctime").Limit(ps, (pi-1)*ps).Find(&mods)
	return mods, err
}

// DictCount 字典分页总数
func DictCount() int {
	mod := &Dict{}
	sess := Db.NewSession()
	defer sess.Close()
	count, _ := sess.Count(mod)
	return int(count)
}

// DictAdd 添加字典
func DictAdd(mod *Dict) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.InsertOne(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	initDict()
	return nil
}

// DictEdit 编辑字典
func DictEdit(mod *Dict, cols ...string) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(mod.Key).Cols(cols...).Update(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	initDict()
	return nil
}

// DictDrop 删除单条字典
func DictDrop(key string) error {
	sess := Db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.ID(key).Delete(&Dict{}); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	initDict()
	return nil
}
