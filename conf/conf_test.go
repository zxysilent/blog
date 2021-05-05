package conf

import (
	"os"
	"testing"
)

func TestConf(t *testing.T) {
	defConfig = "./conf.toml"
	conf, err := initCfg()
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(conf)
		t.Log(conf.Dsn())
	}
	t.Cleanup(func() {
		os.RemoveAll("logs")
	})
}
