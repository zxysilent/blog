package conf

import (
	"os"
	"testing"
)

func TestConf(t *testing.T) {
	defer os.RemoveAll("logs")
	defConfig = "../conf.toml"
	conf, err := initConfig()
	if err != nil {
		t.Fatal(err)
	} else {
		t.Logf("%+v", conf)
		t.Log(conf.Db.Mysql)
		t.Log(conf.Db.Sqlite)
	}
}
