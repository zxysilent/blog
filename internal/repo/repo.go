package repo

import (
	"blog/conf"
	"blog/internal/model"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"reflect"
	"runtime"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/zxysilent/logs"
	_ "modernc.org/sqlite"
	"xorm.io/xorm"
	"xorm.io/xorm/contexts"
	"xorm.io/xorm/log"
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
		err = fmt.Errorf("unsupported dbtype (%s)", conf.App.Db.Type)
	}
	if err != nil {
		panic("数据库 dsn:" + err.Error())
	}
	logs.Ctx(ctx).Info("use db :", conf.App.Db.Type)
	db.SetLogLevel(log.LOG_OFF)
	db.AddHook(&repoHook{showSql: conf.App.Db.ShowSql})
	if err = db.Ping(); err != nil {
		panic("数据库 ping:" + err.Error())
	}
	db.SetMaxIdleConns(max(runtime.NumCPU()/2, 2))
	db.SetMaxOpenConns(max(runtime.NumCPU()/2, 4))
	// mysql int(11)、tinyint(4)、smallint(6)、mediumint(9)、bigint(20)
	logs.Ctx(ctx).Err(db.Sync(&model.Tag{})).Info("db.sync.Tag")
	logs.Ctx(ctx).Err(db.Sync(&model.Cate{})).Info("db.sync.Cate")
	logs.Ctx(ctx).Err(db.Sync(&model.Post{})).Info("db.sync.Post")
	logs.Ctx(ctx).Err(db.Sync(&model.PostTag{})).Info("db.sync.PostTag")
	/*-----------WARNING BEGIN-----------*/
	logs.Ctx(ctx).Err(db.Sync(&model.Dict{})).Info("db.sync.Dict")
	logs.Ctx(ctx).Err(db.Sync(&model.Role{})).Info("db.sync.Role")
	logs.Ctx(ctx).Err(db.Sync(&model.Grant{})).Info("db.sync.Grant")
	logs.Ctx(ctx).Err(db.Sync(&model.Admin{})).Info("db.sync.Admin")
	logs.Ctx(ctx).Err(db.Sync(&model.RoleGrant{})).Info("db.sync.RoleGrant")
	/*-----------WARNING END-----------*/
	// // 预加载
	// if !AdminExist("iadmin") {
	// 	logs.Ctx(ctx).Info("init user iadmin")
	// 	milli := time.Now().UnixMilli()
	// 	AdminAdd(&model.Admin{Num: "iadmin", Name: "管理员", Passwd: "33d7be2196ff70efaf6913fc8bdcab", RoleId: 1, Updated: milli, Created: milli})
	// }
	// if !DictExist(model.DictKeyBasic) {
	// 	logs.Ctx(ctx).Info("init dict basic")
	// 	milli := time.Now().UnixMilli()
	// 	DictAdd(&model.Dict{Key: model.DictKeyBasic, Value: `{"title":"","site_url":"","logo_url":"","author":"","copyright":"","keywords":"","description":"","beian_miit":"","beian_nism":""}`, Type: "json", Inner: true, Intro: "基础配置", Updated: milli, Created: milli})
	// }
	// if !DictExist(model.DictKeyLogin) {
	// 	logs.Info("init dict login")
	// 	milli := time.Now().UnixMilli()
	// 	DictAdd(&model.Dict{Key: model.DictKeyLogin, Value: `{"title":"欢迎登录","desc":"便捷、高效、可靠"}`, Type: "json", Inner: true, Intro: "登录页面", Updated: milli, Created: milli})
	// }
	// if !DictExist(model.DictKeyBlog) {
	// 	logs.Ctx(ctx).Info("init dict blog")
	// 	milli := time.Now().UnixMilli()
	// 	DictAdd(&model.Dict{Key: model.DictKeyBlog, Value: `{"page_size":0,"analytic":"","comment":"","github_url":"","weibo_url":""}`, Type: "json", Inner: true, Intro: "Blog配置", Updated: milli, Created: milli})
	// }
	SyncCfg(ctx)
	logs.Ctx(ctx).Info("model init")
}

func GetDb() *xorm.Engine {
	return db
}

func Close() {
	db.Close()
}

type repoHook struct {
	showSql bool
}

func (rh *repoHook) BeforeProcess(ctx *contexts.ContextHook) (context.Context, error) {
	return ctx.Ctx, nil
}

func (rh *repoHook) AfterProcess(ctx *contexts.ContextHook) error {
	slog.Error(logs.TraceOf(ctx.Ctx))
	if ctx.Err != nil {
		logs.Ctx(ctx.Ctx).Caller(false).Err(ctx.Err).Str("SQL", ctx.SQL).Any("args", ctx.Args).Dur("dur", ctx.ExecuteTime).Error()
	} else if ctx.ExecuteTime > 300*time.Millisecond {
		logs.Ctx(ctx.Ctx).Caller(false).Str("SlowSQL", ctx.SQL).Any("args", ctx.Args).Dur("dur", ctx.ExecuteTime).Warn()
	} else if rh.showSql {
		logs.Ctx(ctx.Ctx).Caller(false).Str("SQL", ctx.SQL).Any("args", ctx.Args).Dur("dur", ctx.ExecuteTime).Debug()
	}
	return ctx.Err
}

var (
	ErrExist    = errors.New("db row exist")
	ErrNotExist = errors.New("db row not exist")
	ErrPkZero   = errors.New("pk cannot be empty")
)

func GetAny[T any](ctx context.Context, pk any) (T, error) {
	var mod T
	if reflect.ValueOf(pk).IsZero() {
		return mod, ErrPkZero
	}
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.ID(pk)
	sess.NoAutoCondition()
	if has, _ := sess.Get(&mod); !has {
		return mod, ErrNotExist
	}
	return mod, nil
}

func DelAny[T any](ctx context.Context, pk any) error {
	var mod T
	if reflect.ValueOf(pk).IsZero() {
		return ErrPkZero
	}
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.ID(pk)
	sess.NoAutoCondition()
	if _, err := sess.Delete(&mod); err != nil {
		return err
	}
	return nil
}

func ExistAny[T any](ctx context.Context, mod T) bool {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	exist, _ := sess.Exist(&mod)
	return exist
}

func AddAny[T any](ctx context.Context, mod *T) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	if _, err := sess.InsertOne(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

type EAny interface {
	GetId() string
}

func EditAny[T EAny](ctx context.Context, mod *T, cols ...string) error {
	sess := db.NewSession()
	defer sess.Close()
	sess.Context(ctx)
	sess.Begin()
	if _, err := sess.ID((*mod).GetId()).Cols(cols...).Update(mod); err != nil {
		sess.Rollback()
		return err
	}
	sess.Commit()
	return nil
}

func fixlike(str string) string {
	str = strings.ReplaceAll(str, "_", "\\_")
	return strings.ReplaceAll(str, "%", "\\%")
}

// LOG for ext
func LOG(any interface{}, modeule string, action string, data ...any) {
	// adminId := 0
	// switch ctx := any.(type) {
	// case int:
	// 	adminId = ctx
	// case echo.Context:
	// 	adminId, _ = ctx.Get("uid").(int)
	// default:
	// }
	// mod := model.Log{
	// 	AdminId: adminId,
	// 	Module:  modeule,
	// 	Action:  action,
	// }
	// if len(data) > 0 {
	// 	switch ext := data[0].(type) {
	// 	case int:
	// 		mod.Data = strconv.Itoa(ext)
	// 	case string:
	// 		mod.Data = ext
	// 	default:
	// 		b, _ := json.Marshal(data[0])
	// 		mod.Data = string(b)
	// 	}
	// }
	// mod.Created = time.Now().Unix()
}
