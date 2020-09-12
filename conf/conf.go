package conf

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/zxysilent/logs"
)

type appconf struct {
	Title   string `toml:"title"`
	Explain string `toml:"explain"`
	Mode    string `toml:"mode"`
	Addr    string `toml:"addr"`
	Srv     string `toml:"srv"`
	Jwtkey  string `toml:"jwtkey"`
	Jwtexp  int    `toml:"jwtexp"`
	Author  struct {
		Name    string `toml:"name"`
		Website string `toml:"website"`
	} `toml:"author"`
	Wechat struct {
		Appid  string `toml:"appid"`
		Secret string `toml:"secret"`
	} `toml:"wechat"`
	Database struct {
		Host   string `toml:"host"`
		Port   int    `toml:"port"`
		User   string `toml:"user"`
		Passwd string `toml:"passwd"`
		Dbname string `toml:"dbname"`
		Params string `toml:"params"`
	} `toml:"database"`
	Xorm struct {
		Idle        int  `toml:"idle"`
		Open        int  `toml:"open"`
		Show        bool `toml:"show"`
		Sync        bool `toml:"sync"`
		CacheEnable bool `toml:"cache_enable"`
		CacheCount  int  `toml:"cache_count"`
	} `toml:"xorm"`
}

func (app *appconf) IsProd() bool {
	return app.Mode == "prod"
}
func (app *appconf) IsDev() bool {
	return app.Mode == "dev"
}

const _dsn = "%s:%s@tcp(%s:%d)/%s?%s"

func (app *appconf) Dsn() string {
	return fmt.Sprintf(_dsn, app.Database.User, app.Database.Passwd, app.Database.Host, app.Database.Port, app.Database.Dbname, app.Database.Params)
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
