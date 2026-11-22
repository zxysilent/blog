package repo

import (
	"blog/conf"
	"blog/internal/model"
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strings"

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
	// 显示 SQL 与慢 SQL 阈值改从配置表读取（主键 sql），可在线调整
	// 此时配置缓存尚未加载，hook 会取到默认值，SyncCfg 之后自动生效
	db.AddHook(&repoHook{})
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
	logs.Ctx(ctx).Err(db.Sync(&model.User{})).Info("db.sync.User")
	logs.Ctx(ctx).Err(db.Sync(&model.Role{})).Info("db.sync.Role")
	// 派生表：结构与代码定义不一致时直接重建，不做任何迁移
	resetDerived(ctx)
	logs.Ctx(ctx).Err(db.Sync(&model.Perm{})).Info("db.sync.Perm")
	logs.Ctx(ctx).Err(db.Sync(&model.RolePerm{})).Info("db.sync.RolePerm")
	logs.Ctx(ctx).Err(db.Sync(&model.Dict{})).Info("db.sync.Dict")
	logs.Ctx(ctx).Err(db.Sync(&model.DictEntry{})).Info("db.sync.DictEntry")
	logs.Ctx(ctx).Err(db.Sync(&model.Config{})).Info("db.sync.Config")
	logs.Ctx(ctx).Err(db.Sync(&model.Log{})).Info("db.sync.Log")
	/*-----------WARNING END-----------*/
	// 承载用户内容的表不做重建，但主键类型漂移时必须提示出来：
	// 老库用 int 自增主键，现在是字符串 ULID，SQLite 不做隐式转换，
	// 表现为「数据在库里、接口却查不到」，只靠日志才查得出来
	warnSchemaDrift(ctx)
	// 幂等种子：角色 / 用户 / 权限点 / 角色授权 / 配置 / 字典
	logs.Ctx(ctx).Err(seed(ctx)).Info("repo.seed")
	// 栏目：补齐历史缺失路径、重算层级、重建内存缓存
	logs.Ctx(ctx).Err(CateInit(ctx)).Info("repo.CateInit")
	if err := SyncCfg(ctx); err != nil {
		logs.Ctx(ctx).Err(err).Warn("repo.SyncCfg")
	}
	// 默认私钥是公开可见的，长期不修改等于没有签名保护
	if GetCfg().Auth.TokenSecret == model.AuthSecretDefault {
		logs.Ctx(ctx).Warn("auth.token_secret is still the default value, please change it")
	}
	SyncRolePerms()
	logs.Ctx(ctx).Info("model init")
}

func GetDb() *xorm.Engine {
	return db
}

// resetDerived 重置「可由代码重建」的表
// 判断依据是库中结构与模型定义不一致时，直接删表重建，不做任何迁移代码 ——
// 等价于「每次当新打开」。
//
// 这些表的数据要么完全由路由注册派生（permission / role_permission，主键由
// perm 稳定推导，重建后 id 不变，角色授权不会失效），要么每次启动都会按代码
// 里的默认值幂等补齐（dict / dict_entry 见 seedDict）。
// config 与 log 刻意不在此列：前者承载在线调整过的运行参数，后者是审计记录，
// 它们的新增列 db.Sync 能自己加，不需要重建。
func resetDerived(ctx context.Context) {
	if schemaStale("permission", &model.Perm{}) {
		logs.Ctx(ctx).Warn("permission schema changed, drop and recreate")
		logs.Ctx(ctx).Err(db.DropTables(&model.Perm{})).Info("db.reset.Perm")
	}
	if schemaStale("role_permission", &model.RolePerm{}) {
		logs.Ctx(ctx).Warn("role_permission schema changed, drop and recreate")
		logs.Ctx(ctx).Err(db.DropTables(&model.RolePerm{})).Info("db.reset.RolePerm")
	}
	if schemaStale("dict", &model.Dict{}) {
		logs.Ctx(ctx).Warn("dict schema changed, drop and recreate")
		logs.Ctx(ctx).Err(db.DropTables(&model.Dict{})).Info("db.reset.Dict")
	}
	if schemaStale("dict_entry", &model.DictEntry{}) {
		logs.Ctx(ctx).Warn("dict_entry schema changed, drop and recreate")
		logs.Ctx(ctx).Err(db.DropTables(&model.DictEntry{})).Info("db.reset.DictEntry")
	}
}

// warnSchemaDrift 提示「承载用户内容、不能自动重建」的表结构漂移
// 这些表的数据无法从代码重建，因此只告警不处理；
// 出现告警时应手工迁移（或换新库）后再使用，否则查询会静默查不到数据。
func warnSchemaDrift(ctx context.Context) {
	for _, item := range []struct {
		table string
		bean  any
	}{
		{"tag", &model.Tag{}},
		{"cate", &model.Cate{}},
		{"post", &model.Post{}},
		{"post_tag", &model.PostTag{}},
		{"user", &model.User{}},
		{"role", &model.Role{}},
		{"config", &model.Config{}},
		{"log", &model.Log{}},
	} {
		if schemaStale(item.table, item.bean) {
			logs.Ctx(ctx).Warn("table schema differs from model, please migrate or use a fresh database: " + item.table)
		}
	}
}

// schemaStale 判断库中表结构与模型定义不一致
// 三种情况视为不一致：
//  1. 库里还留着模型已删除的列（Sync 只会加列，不会删列）
//  2. 库里缺少模型的主键列 —— SQLite 无法给已有表补主键列，
//     这种情况 Sync 会直接失败，必须删表重建
//  3. 主键列类型不同（例如老版本 int 主键、现在改成字符串 ULID）——
//     类型不一致不会报错，但按字符串条件查 int 列永远查不到（SQLite 不做隐式转换），
//     表现为「数据在库里，接口却查不到」，必须重建
//
// 表不存在或取不到元信息时返回 false，交由 db.Sync 建表处理。
func schemaStale(table string, bean any) bool {
	ti, err := db.TableInfo(bean)
	if err != nil {
		return false
	}
	want := make(map[string]struct{}, len(ti.Columns()))
	// 主键：列名 -> 期望的 SQL 类型名，用于发现老库的 int 主键
	wantPk := make(map[string]string, 1)
	for _, col := range ti.Columns() {
		want[col.Name] = struct{}{}
		if col.IsPrimaryKey {
			wantPk[col.Name] = col.SQLType.Name
		}
	}
	metas, err := db.DBMetas()
	if err != nil {
		return false
	}
	for _, meta := range metas {
		if meta.Name != table {
			continue
		}
		have := make(map[string]struct{}, len(meta.Columns()))
		for _, col := range meta.Columns() {
			have[col.Name] = struct{}{}
			// 多余列：模型已删字段，库还留着
			if _, ok := want[col.Name]; !ok {
				return true
			}
		}
		// 缺失主键列，或主键类型与模型不符
		for name, typ := range wantPk {
			if _, ok := have[name]; !ok {
				return true
			}
			for _, col := range meta.Columns() {
				if col.Name == name && typeFamily(col.SQLType.Name) != typeFamily(typ) {
					return true
				}
			}
		}
		return false
	}
	return false
}

// typeFamily 把数据库类型名归到大类
// 不同方言对同一 Go 类型给出的类型名差别很大（模型 VARCHAR → SQLite 读回来是 TEXT），
// 直接比字符串会误报，因此只比较「文本 / 整数 / 其它」这几类，
// 足以识别老库里 int 自增主键改成字符串 ULID 这种真正会出问题的漂移。
func typeFamily(name string) string {
	up := strings.ToUpper(strings.TrimSpace(name))
	switch {
	case strings.Contains(up, "CHAR"), strings.Contains(up, "TEXT"), strings.Contains(up, "CLOB"):
		return "text"
	case strings.Contains(up, "INT"):
		return "int"
	}
	return "other"
}

func Close() {
	db.Close()
}

// repoHook SQL 日志 hook
// 打印策略全部来自配置表（主键 sql），不持有快照：
// 配置修改后下一次查询立即生效，无需重启
//   - 出错：始终按 error 打印，与配置无关
//   - 慢 SQL：超过 SlowMs 按 warn 打印，SlowMs <= 0 时不记录
//   - 全部 SQL：Show 为真时按 debug 打印
type repoHook struct{}

func (rh *repoHook) BeforeProcess(ctx *contexts.ContextHook) (context.Context, error) {
	return ctx.Ctx, nil
}

func (rh *repoHook) AfterProcess(ctx *contexts.ContextHook) error {
	if ctx.Err != nil {
		logs.Ctx(ctx.Ctx).Caller(false).Err(ctx.Err).Str("SQL", ctx.SQL).Any("args", ctx.Args).Dur("dur", ctx.ExecuteTime).Error()
		return ctx.Err
	}
	cfg := GetCfg().Sql
	if slow := cfg.SlowDuration(); slow > 0 && ctx.ExecuteTime > slow {
		logs.Ctx(ctx.Ctx).Caller(false).Str("SlowSQL", ctx.SQL).Any("args", ctx.Args).Dur("dur", ctx.ExecuteTime).Warn()
		return nil
	}
	if cfg.Show {
		logs.Ctx(ctx.Ctx).Caller(false).Str("SQL", ctx.SQL).Any("args", ctx.Args).Dur("dur", ctx.ExecuteTime).Debug()
	}
	return nil
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
