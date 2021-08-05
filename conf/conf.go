package conf

import (
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/zxysilent/logs"
)

const (
	dev  = "dev"  //开发模式
	prod = "prod" //线上模式
)

type appcfg struct {
	Title        string `toml:"title"`          //
	Intro        string `toml:"intro"`          //
	Mode         string `toml:"mode"`           //
	Addr         string `toml:"addr"`           //
	Srv          string `toml:"srv"`            //
	TokenKey     string `toml:"token_key"`      //token关键词
	TokenExp     int    `toml:"token_exp"`      //过期时间 h
	TokenKeep    bool   `toml:"token_keep"`     //保持在线
	TokenSso     bool   `toml:"token_sso"`      //单点登录
	TokenSecret  string `toml:"token_secret"`   //加密私钥
	ImageCut     bool   `toml:"image_cut"`      //图片裁剪
	ImageWidth   int    `toml:"image_width"`    //图片宽度
	ImageHeight  int    `toml:"image_height"`   //图片高度
	PageMin      int    `toml:"page_min"`       //最小分页大小
	PageMax      int    `toml:"page_max"`       //最大分页大小
	DbHost       string `toml:"db_host"`        //数据库地址
	DbPort       int    `toml:"db_port"`        //数据库端口
	DbUser       string `toml:"db_user"`        //数据库账号
	DbPasswd     string `toml:"db_passwd"`      //数据库密码
	DbName       string `toml:"db_name"`        //数据库名称
	DbParams     string `toml:"db_params"`      //数据库参数
	OrmIdle      int    `toml:"orm_idle"`       //
	OrmOpen      int    `toml:"orm_open"`       //
	OrmShow      bool   `toml:"orm_show"`       //显示sql
	OrmSync      bool   `toml:"orm_sync"`       //同步表结构
	OrmCacheUse  bool   `toml:"orm_cache_use"`  //是否使用缓存
	OrmCacheSize int    `toml:"orm_cache_size"` //缓存数量
	OrmHijackLog bool   `toml:"orm_hijack_log"` //劫持日志
	Author       struct {
		Name    string `toml:"name"`
		Website string `toml:"website"`
	} `toml:"author"`
	Wechat struct {
		GzhAppid  string `toml:"gzh_appid"`  //公众号
		GzhSecret string `toml:"gzh_secret"` //公众号
		MpgAppid  string `toml:"mpg_appid"`  //小程序
		MpgSecret string `toml:"mpg_secret"` //小程序
		WebAppid  string `toml:"web_appid"`  //web
	} `toml:"wechat"`
	Qq struct {
		WebAppid  string `toml:"web_appid"`  //web
		WebSecret string `toml:"web_secret"` //appkey-申请名称appkey,对接名称secret
	} `toml:"qq"`
}

func (app *appcfg) IsProd() bool {
	return app.Mode == prod
}
func (app *appcfg) IsDev() bool {
	return app.Mode == dev
}

// mb4 most bytes 4
// uid:pass@tcp(host:port)/dbname?charset=utf8mb4&parseTime=true&loc=Local
// 用户名:密码@tcp(主机:端口)/数据库名称?charset=utf8mb4&parseTime=true&loc=Local
const _dsn = "%s:%s@tcp(%s:%d)/%s?%s"

// MySQL链接字符串
func (app *appcfg) Dsn() string {
	return fmt.Sprintf(_dsn, app.DbUser, app.DbPasswd, app.DbHost, app.DbPort, app.DbName, app.DbParams)
}

var (
	App       *appcfg              //运行配置实体
	defConfig = "./conf/conf.toml" //配置文件路径，方便测试
)

func Init() {
	var err error
	App, err = initCfg()
	if err != nil {
		logs.Fatal("config init error : ", err.Error())
	}
	logs.Debug("conf init")
}

func initCfg() (*appcfg, error) {
	app := &appcfg{}
	_, err := toml.DecodeFile(defConfig, &app)
	if err != nil {
		return nil, err
	}
	return app, nil
}
