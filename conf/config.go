package conf

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type application struct {
	Applet applet `xml:"applet"`
	Wechat wechat `xml:"wechat"`
	Dsn    dsn    `xml:"dsn"`
	Xorm   xorm   `xml:"xorm"`
	Extra  extra  `xml:"extra"`
}
type applet struct {
	Name   string `xml:"name"`
	Depict string `xml:"depict"`
	Jwtkey string `xml:"jwtkey"`
	Jwtexp int    `xml:"jwtexp"`
	Author string `xml:"author"`
	Mode   string `xml:"mode"`
	Srv    string `xml:"srv"`
	Addr   string `xml:"addr"`
}

func (app *applet) IsProd() bool {
	return app.Mode == "prod"
}
func (app *applet) IsDev() bool {
	return app.Mode == "dev"
}

type wechat struct {
	Appid  string `xml:"appid"`
	Secret string `xml:"secret"`
}
type dsn struct {
	User  string `xml:"user"`
	Pass  string `xml:"pass"`
	Host  string `xml:"host"`
	Port  int    `xml:"port"`
	Name  string `xml:"name"`
	Extra string `xml:"extra"`
}

const _dsn = "%s:%s@tcp(%s:%d)/%s?%s"

func (dsn *dsn) Dsn() string {
	return fmt.Sprintf(_dsn, dsn.User, dsn.Pass, dsn.Host, dsn.Port, dsn.Name, dsn.Extra)
}

type xorm struct {
	Idle  int  `xml:"idle"`
	Open  int  `xml:"open"`
	Show  bool `xml:"show"`
	Sync  bool `xml:"sync"`
	Cache struct {
		Enable bool `xml:"enable"`
		Count  int  `xml:"count"`
	} `xml:"cache"`
}

var _config application
var defConfig = "./conf/conf.xml"

func initConfig() (*application, error) {
	conf := &application{}
	raw, err := os.Open(defConfig)
	if err != nil {
		return nil, err
	}
	defer raw.Close()
	decoder := xml.NewDecoder(raw)
	err = decoder.Decode(&conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

// TestSyntax 测试配置文件
func TestSyntax() error {
	_, err := initConfig()
	return err
}

// Init 加载配置文件
func Init() {
	conf, _ := initConfig()
	App = conf.Applet
	Xorm = conf.Xorm
	Wechat = conf.Wechat
	Dsn = conf.Dsn
	Extra = conf.Extra
}

var (
	// App app config
	App applet
	// Xorm xorm config
	Xorm xorm
	// Wechat wechat config
	Wechat wechat
	// Dsn dsn config
	Dsn dsn
	// Extra extra config
	Extra extra
)

type xmlEntry struct {
	XMLName xml.Name
	Value   string `xml:",chardata"`
}
type extra map[string]string

// UnmarshalXML xml2map
func (ext *extra) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	*ext = extra{}
	for {
		var e xmlEntry
		err := d.Decode(&e)
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		(*ext)[e.XMLName.Local] = e.Value
	}
	return nil
}
func (ext *extra) GetInt(key string) (int, error) {
	return strconv.Atoi((*ext)[key])
}
func (ext *extra) GetBool(key string) (bool, error) {
	return strconv.ParseBool((*ext)[key])
}
func (ext *extra) GetTime(key string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", (*ext)[key], time.Local)
}
func (ext *extra) GetString(key string) (string, bool) {
	val, has := (*ext)[key]
	return val, has
}
