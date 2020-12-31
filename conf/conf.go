package conf

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/zxysilent/logs"
)

type appconf struct {
	Title       string `toml:"title"`
	Intro       string `toml:"intro"`
	Mode        string `toml:"mode"`
	Addr        string `toml:"addr"`
	Srv         string `toml:"srv"`
	TokenKey    string `toml:"token_key"`
	TokenExp    int    `toml:"token_exp"`
	TokenKeep   bool   `toml:"token_keep"`
	TokenSecret string `toml:"token_secret"`
	Author      struct {
		Name    string `toml:"name"`
		Website string `toml:"website"`
	} `toml:"author"`
	Wechat struct {
		Appid  string `toml:"appid"`
		Secret string `toml:"secret"`
	} `toml:"wechat"`
	ImageCut     bool   `toml:"image_cut"`
	ImageWidth   int    `toml:"image_width"`
	ImageHeight  int    `toml:"image_height"`
	PageMin      int    `toml:"page_min"`
	PageMax      int    `toml:"page_max"`
	DbHost       string `toml:"db_host"`
	DbPort       int    `toml:"db_port"`
	DbUser       string `toml:"db_user"`
	DbPasswd     string `toml:"db_passwd"`
	DbName       string `toml:"db_name"`
	DbParams     string `toml:"db_params"`
	OrmIdle      int    `toml:"orm_idle"`
	OrmOpen      int    `toml:"orm_open"`
	OrmShow      bool   `toml:"orm_show"`
	OrmSync      bool   `toml:"orm_sync"`
	OrmCacheUse  bool   `toml:"orm_cache_use"`
	OrmCacheSize int    `toml:"orm_cache_size"`
	OrmHijackLog bool   `toml:"orm_hijack_log"`
}

func (app *appconf) IsProd() bool {
	return app.Mode == "prod"
}
func (app *appconf) IsDev() bool {
	return app.Mode == "dev"
}

//uid:pass@tcp(host:port)/dbname?charset=utf8&parseTime=true
//用户名:密码@tcp(主机:端口)/数据库名称?charset=utf8&parseTime=true
const _dsn = "%s:%s@tcp(%s:%d)/%s?%s"

func (app *appconf) Dsn() string {
	return fmt.Sprintf(_dsn, app.DbUser, app.DbPasswd, app.DbHost, app.DbPort, app.DbName, app.DbParams)
}

var (
	App       *appconf
	defConfig = "./conf/conf.toml"
)

func Init() {
	var err error
	App, err = initConf()
	if err != nil {
		logs.Fatal("config init error : ", err.Error())
	}
	logs.Debug("conf init")
}

func initConf() (*appconf, error) {
	app := &appconf{}
	_, err := toml.DecodeFile(defConfig, &app)
	if err != nil {
		return nil, err
	}
	return app, nil
}
