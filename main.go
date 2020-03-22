package main

import (
	"blog/conf"
	"blog/model"
	"blog/router"
	"log"
	"os"
	"os/signal"
	"syscall"
)

// @Title Blog’s Api文档
// @Version 1.0
// @Description token传递方式包括 [get/post]token 、[header] Authorization=Bearer xxxx
// @Description 数据传递方式包括 json、formData 推荐使用 json
// @Description /api/* 公共访问
// @Description /adm/* 必须传入 token
// @Host 127.0.0.1:88
// @BasePath /
func main() {
	log.Println("app initializing")
	conf.Init()
	model.Init()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	log.Println("app running")
	go router.RunApp()
	<-quit
	log.Println("app quitted")
}
