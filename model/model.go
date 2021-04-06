package model

import (
	"blog/conf"
	"log"
	"strings"

	"xorm.io/xorm"
	"xorm.io/xorm/caches"

	// 数据库驱动
	_ "github.com/go-sql-driver/mysql"
	"github.com/zxysilent/logs"
	xlog "xorm.io/xorm/log"
)

// Db 数据库操作句柄
var Db *xorm.Engine

func Init() {
	// 初始化数据库操作的 Xorm
	db, err := xorm.NewEngine("mysql", conf.App.Dsn())
	if err != nil {
		logs.Fatal("数据库 dsn:", err.Error())
	}
	if conf.App.OrmHijackLog {
		sl := &xlog.SimpleLogger{
			DEBUG: log.New(logs.Writer(), "", log.Ldate|log.Ltime),
			ERR:   log.New(logs.Writer(), "", log.Ldate|log.Ltime),
			INFO:  log.New(logs.Writer(), "", log.Ldate|log.Ltime),
			WARN:  log.New(logs.Writer(), "", log.Ldate|log.Ltime),
		}
		if conf.App.IsDev() {
			sl.SetLevel(xlog.LOG_DEBUG)
		} else {
			sl.SetLevel(xlog.LOG_WARNING)
		}
		db.SetLogger(sl)
	}
	if err = db.Ping(); err != nil {
		logs.Fatal("数据库 ping:", err.Error())
	}
	db.SetMaxIdleConns(conf.App.OrmIdle)
	db.SetMaxOpenConns(conf.App.OrmOpen)
	// 是否显示sql执行的语句
	db.ShowSQL(conf.App.OrmShow)
	if conf.App.OrmCacheUse {
		// 设置xorm缓存
		cacher := caches.NewLRUCacher(caches.NewMemoryStore(), conf.App.OrmCacheSize)
		db.SetDefaultCacher(cacher)
	}
	// mysql int(11)、tinyint(4)、smallint(6)、mediumint(9)、bigint(20)

	if conf.App.OrmSync {
		err := db.Sync2(
			new(User),
			new(Cate),
			new(Tag),
			new(Post),
			new(PostTag),
			new(Opts),
			new(Role),
			// new(Api),
			// new(RoleApi),
			new(Menu),
			new(RoleMenu),
		)
		if err != nil {
			logs.Fatal("数据库 sync:", err.Error())
		}
	}
	Db = db
	//缓存
	initMap()
	logs.Info("model init")
}

// Page 分页基本数据
type Page struct {
	Pi   int    `json:"pi" form:"pi" query:"pi"`       //分页页码
	Ps   int    `json:"ps" form:"ps" query:"ps"`       //分页大小
	Mult string `json:"mult" form:"mult" query:"mult"` //多条件信息
}

// Trim 去除空白字符
func (p *Page) Trim() string {
	p.Mult = strings.TrimSpace(p.Mult)
	return p.Mult
}

type IptId struct {
	Id int `query:"id" json:"id"` //仅包含Id的请求
}

// Naver 上下页
type Naver struct {
	Prev string
	Next string
}

// State 统计信息
type State struct {
	Post int `json:"post"`
	Page int `json:"page"`
	Cate int `json:"cate"`
	Tag  int `json:"tag"`
}

// Collect 统计信息
func Collect() (*State, bool) {
	mod := &State{}
	has, _ := Db.SQL(`SELECT * FROM(SELECT COUNT(id) as post FROM post WHERE type=0)as a ,(SELECT COUNT(id) as page FROM post WHERE type=1) as b, (SELECT COUNT(id) as cate FROM cate) as c, (SELECT COUNT(id) as tag FROM tag) as d`).Get(mod)
	return mod, has
}
func inOf(goal int, arr []int) bool {
	for idx := range arr {
		if goal == arr[idx] {
			return true
		}
	}
	return false
}

// Reply 生成api文档使用
// 代码里未使用，也不要使用
type Reply struct {
	Code int    `json:"code" example:"200"`
	Msg  string `json:"msg" example:"提示信息"`
}
