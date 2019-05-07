package model

import (
	"blog/conf"
	"log"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	// 数据库驱动
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

// DB 数据库操作句柄
var DB *xorm.Engine

func init() {
	var err error
	// 初始化数据库操作的 Xorm
	DB, err = xorm.NewEngine("mysql", conf.Conn)
	if err != nil {
		log.Fatalln("数据库连接失败", err.Error())
	}
	DB.SetMaxIdleConns(conf.MaxIdle)
	DB.SetMaxOpenConns(conf.MaxOpen)
	if err = DB.Ping(); err != nil {
		log.Fatalln("数据库不能正常工作", err.Error())
	}
	// 是否显示sql执行的语句
	DB.ShowSQL(conf.ShowSQLl)
	if conf.UseCache {
		// 设置xorm缓存
		cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), conf.CacheCount)
		DB.SetDefaultCacher(cacher)
	}
	// 同步表
	DB.Sync2(new(User))
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

// JwtClaims jwt
type JwtClaims struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Num  string `json:"num"`
	Role Role   `json:"role"`
	jwt.StandardClaims
}

// Naver 上下页
type Naver struct {
	Prev string
	Next string
}

// State 统计信息
type State struct {
	Post int `json:"post" form:"post"`
	Page int `json:"page" form:"page"`
	Cate int `json:"cate" form:"cate"`
	Tag  int `json:"tag" form:"tag"`
}

// Collect 统计信息
func Collect() (*State, bool) {
	mod := &State{}
	has, _ := DB.SQL(`SELECT * FROM(SELECT COUNT(id) as post FROM post WHERE type=0)as a ,(SELECT COUNT(id) as page FROM post WHERE type=1) as b, (SELECT COUNT(id) as cate FROM cate) as c, (SELECT COUNT(id) as tag FROM tag) as d`).Get(mod)
	return mod, has
}
