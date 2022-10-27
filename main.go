package main

import (
	"blog/conf"
	"blog/model"
	"blog/router"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zxysilent/logs"
)

// @Title Blog’s 接口文档
// @Version 1.0
// @Description token传递方式包括[get/post]token、[header]Authorization
// @Description /api/* 公共访问
// @Description /adm/* 必须传入 token
// @Host 127.0.0.1:8085
// @BasePath /
func main() {
	logs.Info("app initializing")
	conf.Init()
	model.Init()
	defer model.Close()
	defer logs.Flush()
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	logs.Info("app running")
	app := router.Init()
	go func() {
		if err := app.Start(conf.App.Addr); err != nil {
			logs.Fatal(err.Error())
		}
		// if err := app.StartTLS(conf.App.Addr, "./server.crt", "./server.key"); err != nil {
		// 	logs.Fatal(err.Error())
		// }
	}()
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		if err != http.ErrServerClosed { //Normal exit
			logs.Error("Server Shutdown:", err)
		}
	}
	logs.Info("app quitted")
}
