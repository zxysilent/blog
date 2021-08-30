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
	Id          int    `xorm:"INT(11) PK AUTOINCR comment('主键')" json:"id"`          //主键
	SiteUrl     string `xorm:"VARCHAR(255) comment('网站地址')" json:"site_url"`         //网站地址
	LogoUrl     string `xorm:"VARCHAR(255) comment('Logo地址')" json:"logo_url"`       //Logo地址
	Author      string `xorm:"VARCHAR(255) comment('网站作者')" json:"author"`           //网站作者
	Title       string `xorm:"VARCHAR(255) comment('网站标题')" json:"title"`            //网站标题
	Keywords    string `xorm:"VARCHAR(255) comment('网站关键词')" json:"keywords"`        //网站关键词
	Description string `xorm:"VARCHAR(255) comment('网站描述')" json:"description"`      //网站描述
	FaviconUrl  string `xorm:"VARCHAR(255) comment('Favicon地址')" json:"favicon_url"` //Favicon地址
	BeianMiit   string `xorm:"VARCHAR(255) comment('ICP备案')" json:"beian_miit"`      //ICP备案
	BeianNism   string `xorm:"VARCHAR(255) comment('公安备案')" json:"beian_nism"`       //公安备案
	Copyright   string `xorm:"VARCHAR(255) comment('版权')" json:"copyright"`          //版权
	SiteJs      string `xorm:"VARCHAR(512) comment('全局js')" json:"site_js"`          //全局js
	SiteCss     string `xorm:"VARCHAR(512) comment('全局css')" json:"site_css"`        //全局css--以上为基本属性
	PageSize    int    `xorm:"INT(11) DEFAULT 6 comment('分页大小')" json:"page_size"`   //分页大小
	Analytic    string `xorm:"VARCHAR(1024) comment('统计脚本')" json:"analytic"`        //统计脚本
	Comment     string `xorm:"VARCHAR(1024) comment('评论脚本')" json:"comment"`         //评论脚本
	GithubUrl   string `xorm:"VARCHAR(255) comment('githu地址')" json:"github_url"`    //githu地址
	WeiboUrl    string `xorm:"VARCHAR(255) comment('微博地址')" json:"weibo_url"`        //微博地址
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
func Gcfg() Global {
	cache := globalCache
	return cache
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
