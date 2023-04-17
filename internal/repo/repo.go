package repo

import (
	"blog/conf"
	"blog/internal/model"
	"context"
	"errors"
	"time"

	"xorm.io/xorm"
	"xorm.io/xorm/caches"

	// 数据库驱动
	_ "github.com/go-sql-driver/mysql"
	"github.com/zxysilent/logs"
	_ "modernc.org/sqlite"
	xlog "xorm.io/xorm/log"
)

// db 数据库操作句柄
var db *xorm.Engine

// Init 初始化数据库连接
func Init(ctx context.Context) {
	var err error
	switch conf.App.Db.Type {
	case conf.DbTypeMysql:
		db, err = xorm.NewEngine("mysql", conf.App.Db.Mysql)
	case conf.DbTypeSqlite:
		db, err = xorm.NewEngine("sqlite", conf.App.Db.Sqlite)
	default:
		err = errors.New("unsupported dbtype")
	}
	if err != nil {
		panic("数据库 dsn:" + err.Error())
	}
	logs.Ctx(ctx).Info("use db :", conf.App.Db.Type)
	// 劫持xorm日志
	logs.Ctx(ctx).Info("orm hijack log :", conf.App.Orm.HijackLog)
	if conf.App.Orm.HijackLog {
		sl := logs.Xorm(xlog.LOG_INFO)
		if conf.IsDebug() {
			sl.SetLevel(xlog.LOG_DEBUG)
		}
		db.SetLogger(sl)
	}
	if err = db.Ping(); err != nil {
		panic("数据库 ping:" + err.Error())
	}
	db.SetMaxIdleConns(conf.App.Orm.Idle)
	db.SetMaxOpenConns(conf.App.Orm.Open)
	// 是否显示sql执行的语句
	logs.Ctx(ctx).Info("show sql ", conf.App.Orm.Show)
	db.ShowSQL(conf.App.Orm.Show)
	logs.Ctx(ctx).Info("use cache ", conf.App.Orm.Show)
	if conf.App.Orm.CacheUse {
		// 设置xorm缓存
		cacher := caches.NewLRUCacher(caches.NewMemoryStore(), conf.App.Orm.CacheSize)
		db.SetDefaultCacher(cacher)
	}
	// mysql int(11)、tinyint(4)、smallint(6)、mediumint(9)、bigint(20)
	if conf.App.Orm.Sync {
		err = db.Sync(
			new(model.Tag),
			new(model.Cate),
			new(model.Post),
			new(model.PostTag),
			/*-----------WARNING BEGIN-----------*/
			new(model.Dict),
			new(model.Role),
			new(model.Grant),
			new(model.Admin),
			new(model.RoleGrant),
			/*-----------WARNING END-----------*/
		)
		if err != nil {
			logs.With().Err(err).Error("db sync")
		}
	}
	// 预加载
	if !AdminExist("iadmin") {
		logs.Ctx(ctx).Info("init user iadmin")
		milli := time.Now().UnixMilli()
		AdminAdd(&model.Admin{Num: "iadmin", Name: "管理员", Passwd: "33d7be2196ff70efaf6913fc8bdcab", RoleId: 1, Updated: milli, Created: milli})
	}
	if !DictExist(model.DictKeyBasic) {
		logs.Ctx(ctx).Info("init dict basic")
		milli := time.Now().UnixMilli()
		DictAdd(&model.Dict{Key: model.DictKeyBasic, Value: `{"title":"","site_url":"","logo_url":"","author":"","copyright":"","keywords":"","description":"","beian_miit":"","beian_nism":""}`, Kind: "json", Inner: true, Intro: "基础配置", Updated: milli, Created: milli})
	}
	if !DictExist(model.DictKeyLogin) {
		logs.Info("init dict login")
		milli := time.Now().UnixMilli()
		DictAdd(&model.Dict{Key: model.DictKeyLogin, Value: `{"title":"欢迎登录","desc":"便捷、高效、可靠"}`, Kind: "json", Inner: true, Intro: "登录页面", Updated: milli, Created: milli})
	}
	if !DictExist(model.DictKeyBlog) {
		logs.Ctx(ctx).Info("init dict blog")
		milli := time.Now().UnixMilli()
		DictAdd(&model.Dict{Key: model.DictKeyBlog, Value: `{"page_size":0,"analytic":"","comment":"","github_url":"","weibo_url":""}`, Kind: "json", Inner: true, Intro: "Blog配置", Updated: milli, Created: milli})
	}
	initBasic()
	initBlog()
	logs.Ctx(ctx).Info("model init")
}

func GetDb() *xorm.Engine {
	return db
}
func Close() {
	db.Close()
}

func inOf(goal int, arr []int) bool {
	for idx := range arr {
		if goal == arr[idx] {
			return true
		}
	}
	return false
}

func Exec(sqlOrArgs ...interface{}) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Begin()
	if _, err := sess.Exec(sqlOrArgs...); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

type IRepo interface {
	Exist(mod *model.Admin) bool
	Get(mod *model.Admin) error
	List(args model.Page, filter model.Admin, cols ...string) ([]model.Admin, int, error)
	Add(mod *model.Admin) error
	Edit(mod *model.Admin, cols ...string) error
	Drop(id string) error
}
