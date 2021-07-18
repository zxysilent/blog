package model

import (
	"errors"

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
	Id              int    `xorm:"INT(11) PK AUTOINCR comment('主键')"` //主键
	SiteUrl         string `xorm:"VARCHAR(255) comment('网站地址')"`      //网站地址
	SiteLogoUrl     string `xorm:"VARCHAR(255) comment('Logo地址')"`    //Logo地址
	SiteTitle       string `xorm:"VARCHAR(255) comment('网站标题')"`      //网站标题
	SiteKeywords    string `xorm:"VARCHAR(255) comment('网站关键词')"`     //网站关键词
	SiteDescription string `xorm:"VARCHAR(255) comment('网站描述')"`      //网站描述
	SiteFaviconUrl  string `xorm:"VARCHAR(255) comment('Favicon地址')"` //Favicon地址
	SiteBeianMiit   string `xorm:"VARCHAR(255) comment('ICP备案')"`     //ICP备案
	SiteBeianNism   string `xorm:"VARCHAR(255) comment('公安备案')"`      //公安备案
	SiteCopyright   string `xorm:"VARCHAR(255) comment('版权')"`        //版权
	SitePowered     string `xorm:"VARCHAR(255) comment('技术支持')"`      //技术支持
	SiteJs          string `xorm:"VARCHAR(255) comment('全局js')"`      //全局js
	SiteCss         string `xorm:"VARCHAR(255) comment('全局css')"`     //全局css--以上为基本属性
	PageSize        int    `xorm:"INT(11) DEFAULT 6 comment('分页大小')"` //分页大小
	Comment         string `xorm:"VARCHAR(255) comment('评论脚本')"`      //评论脚本
	GithubUrl       string `xorm:"VARCHAR(255) comment('githu地址')"`   //githu地址
	WeiboUrl        string `xorm:"VARCHAR(255) comment('微博地址')"`      //微博地址
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
