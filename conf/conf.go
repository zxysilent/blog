package conf

import (
	"github.com/BurntSushi/toml"
	"github.com/zxysilent/logs"
)

type config struct {
	Title       string `toml:"title"`
	Mode        string `toml:"mode"`
	Host        string `toml:"host"`
	Port        int    `toml:"port"`
	RpcHost     string `toml:"rpc_host"`
	RpcPort     int    `toml:"rpc_port"`
	Endpoint    string `toml:"endpoint"`     //
	TokenKey    string `toml:"token_key"`    //token关键词
	TokenExp    int    `toml:"token_exp"`    //过期时间 h
	TokenSecret string `toml:"token_secret"` //加密私钥
	ImageCut    bool   `toml:"image_cut"`    //图片裁剪
	ImageWidth  int    `toml:"image_width"`  //图片宽度
	ImageHeight int    `toml:"image_height"` //图片高度
	PageMin     int    `toml:"page_min"`     //最小分页大小
	PageMax     int    `toml:"page_max"`     //最大分页大小
	Db          struct {
		Type   string `toml:"type"`   //数据库类型 mysql sqlite
		Mysql  string `toml:"mysql"`  //uid:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=true&loc=Local
		Sqlite string `toml:"sqlite"` //"file:sqlite.db?journal_mode=WAL"
	} `toml:"db"`
	Orm struct {
		Idle      int  `toml:"idle"`       //
		Open      int  `toml:"open"`       //
		Show      bool `toml:"show"`       //显示sql
		Sync      bool `toml:"sync"`       //同步表结构
		CacheUse  bool `toml:"cache_use"`  //是否使用缓存
		CacheSize int  `toml:"cache_size"` //缓存数量
		HijackLog bool `toml:"hijack_log"` //劫持日志
	} `toml:"orm"`
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

func Init(path ...string) {
	var err error
	// 方便调试
	if len(path) > 0 {
		defConfig = path[0]
	}
	App, err = initConfig()
	if err != nil {
		panic("config init error : " + err.Error())
	}
	logs.Info("conf init")
}

func initConfig() (*config, error) {
	app := &config{}
	_, err := toml.DecodeFile(defConfig, &app)
	if err != nil {
		return nil, err
	}
	return app, nil
}
