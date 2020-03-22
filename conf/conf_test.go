package conf

import (
	"testing"
)

func TestConf(t *testing.T) {
	defConfig = "./conf.toml"
	conf, err := initConf()
	if err != nil {
		t.Error(err)
	} else {
		t.Log(conf)
		t.Log(conf.Dsn())
	}
}
