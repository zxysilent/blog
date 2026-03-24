package conf

import (
	"context"

	"github.com/BurntSushi/toml"
	"github.com/zxysilent/logs"
)

const (
	TokenKey = "token"
)

type config struct {
	Title       string `toml:"title"`
	Mode        string `toml:"mode"`
	Host        string `toml:"host"`
	Port        int    `toml:"port"`
	Endpoint    string `toml:"endpoint"`     //
	TokenExp    int    `toml:"token_exp"`    //过期时间 h
	TokenSecret string `toml:"token_secret"` //加密私钥
	Db          struct {
		Type    string `toml:"type"`     //数据库类型 mysql sqlite
		Mysql   string `toml:"mysql"`    //uid:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=true&loc=Local
		Sqlite  string `toml:"sqlite"`   //"file:sqlite.db?journal_mode=WAL"
		ShowSql bool   `toml:"show_sql"` //显示sql
	} `toml:"db"`
	Author struct {
		Name    string `toml:"name"`
		Website string `toml:"website"`
	} `toml:"author"`
	ClientAppid  string `toml:"client_appid"`
	ClientSecret string `toml:"client_secret"`
}

const (
	DbTypeMysql  = "mysql"
	DbTypeSqlite = "sqlite"
)

var (
	App       *config
	defConfig = "./conf.toml"
)

func Init(ctx context.Context, path ...string) {
	var err error
	// 方便调试
	if len(path) > 0 {
		defConfig = path[0]
	}
	App, err = initConfig()
	if err != nil {
		panic("config init error : " + err.Error())
	}
	logs.Ctx(ctx).Info("conf init")
}

func initConfig() (*config, error) {
	app := &config{}
	_, err := toml.DecodeFile(defConfig, &app)
	if err != nil {
		return nil, err
	}
	return app, nil
}
