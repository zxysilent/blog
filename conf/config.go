package conf

import (
	"encoding/xml"
	"log"
	"os"
)

var (
	// Name 姓名
	Name string
	// Debug Debug
	Debug bool
	// Conn 连接字符串
	Conn string
	// Addr Addr
	Addr string

	MaxIdle    int
	MaxOpen    int
	ShowSQLl   bool
	UseCache   bool
	CacheCount int
)

// config config
type config struct {
	Name       string `xml:"name"`
	Depict     string `xml:"depict"`
	Debug      bool   `xml:"debug"`
	Conn       string `xml:"conn"`
	Addr       string `xml:"addr"`
	MaxIdle    int    `xml:"max-idle"`
	MaxOpen    int    `xml:"max-open"`
	ShowSQLl   bool   `xml:"show-sql"`
	UseCache   bool   `xml:"use-cache"`
	CacheCount int    `xml:"cache-count"`
}

func init() {
	var conf config
	// 读取配置文件
	xmlRaw, err := os.Open("./conf/conf.xml")
	if err != nil {
		log.Fatalln("配置文件读取失败", err.Error())
	}
	decoder := xml.NewDecoder(xmlRaw)
	err = decoder.Decode(&conf)
	if err != nil {
		log.Fatalln("配置文件无效", err.Error())
	}
	xmlRaw.Close()
	Name = conf.Name
	Debug = conf.Debug
	Conn = conf.Conn
	Addr = conf.Addr
	MaxIdle = conf.MaxIdle
	MaxOpen = conf.MaxOpen
	ShowSQLl = conf.ShowSQLl
	UseCache = conf.UseCache
	CacheCount = conf.CacheCount
}
