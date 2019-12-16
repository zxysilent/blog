package conf

import (
	"log"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	defConfig = "./conf.xml"
	conf, err := initConfig()
	if err != nil {
		log.Fatalln(err)
	} else {
		log.Println(conf)
		log.Println(conf.Dsn.Dsn())
	}
}
