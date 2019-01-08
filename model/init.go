package model

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	jwt "github.com/dgrijalva/jwt-go"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

// DB 数据库操作句柄
var DB *xorm.Engine

// Conf 配置信息
var Conf config

func init() {
	// 读取配置文件
	conf, err := os.Open("./conf.json")
	if err != nil {
		log.Fatalln("配置文件读取失败", err.Error())
	}
	decoder := json.NewDecoder(conf)
	err = decoder.Decode(&Conf)
	if err != nil {
		log.Fatalln("配置文件无效", err.Error())
	}
	// 初始化数据库操作的 Xorm
	DB, err = xorm.NewEngine("mysql", Conf.Conn)
	if err != nil {
		log.Fatalln("数据库连接失败", err.Error())
	}
	DB.SetMaxIdleConns(100)
	DB.SetMaxOpenConns(200)
	if err = DB.Ping(); err != nil {
		log.Fatalln("数据库不能正常工作", err.Error())
	}
	// 更加开发环境是否显示sql执行的语句
	DB.ShowSQL(Conf.Debug)
	// 设置xorm缓存
	cacher := xorm.NewLRUCacher(xorm.NewMemoryStore(), 1000)
	DB.SetDefaultCacher(cacher)
	// DB.Logger().SetLevel(0)
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

// config config
type config struct {
	Debug bool   `json:"debug"`
	Conn  string `json:"conn"`
	Addr  string `json:"addr"`
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
