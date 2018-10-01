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
var Conf Config

func init() {
	conf, err := os.Open("./conf.json")
	if err != nil {
		log.Fatalln("配置文件读取失败", err.Error())
	}
	decoder := json.NewDecoder(conf)
	err = decoder.Decode(&Conf)
	if err != nil {
		log.Fatalln("配置文件无效", err.Error())
	}
	DB, err = xorm.NewEngine("mysql", Conf.Conn)
	if err != nil {
		log.Fatalln("数据库连接失败", err.Error())
	}
	DB.SetMaxIdleConns(100)
	DB.SetMaxOpenConns(200)
	if err = DB.Ping(); err != nil {
		log.Fatalln("数据库不能正常工作", err.Error())
	}
	DB.ShowSQL(Conf.Debug)
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

// Config config
type Config struct {
	Debug bool   `json:"debug"`
	Conn  string `json:"conn"`
	Addr  string `json:"addr"`
}

// Stat 统计信息
type Stat struct {
	Course int `xorm:"INT(11)" json:"course" form:"course"` //课程数
	Period int `xorm:"INT(11)" json:"period" form:"period"` //课时数
	Lib    int `xorm:"INT(11)" json:"lib" form:"lib"`       //题库数
}

// SysStat 统计信息
func SysStat() *Stat {
	mod := &Stat{}
	DB.Sql("SELECT * FROM (SELECT COUNT(id) course FROM course)  a JOIN (SELECT COUNT(id) period FROM period ) b JOIN (SELECT COUNT(id) lib FROM lib ) c").Get(mod)
	return mod
}

//Naver 上下页
type Naver struct {
	Prev string
	Next string
}
