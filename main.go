//go:generate go test -timeout 30s -run ^TestGenRouter$ blog/internal/utils -count=1 -v -msg=路径相对于utils -rp=../router/router_api.go -rs=../service
package main

import (
	"blog/conf"
	"blog/internal/repo"
	"blog/internal/router"
	"blog/internal/service"
	"blog/internal/utils"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/zxysilent/logs"
)

// @Title blog’s Api文档
// @Version 1.0
// @Description 凭证传递方式包括 get、post、header、cookie
// @Description /api/* 前后端api
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description					Description for what is this security definition being used
func main() {
	ctx := logs.TraceCtx(context.Background())
	defer logs.Close()
	logs.Ctx(ctx).Info("app initializing")
	conf.Init()
	repo.Init(ctx)
	defer repo.Close()
	service.Init(ctx)
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	logs.Info("app running")
	if !utils.FileExist("./static") {
		os.MkdirAll("./static", 0755)
	}
	if !utils.FileExist("./dist") {
		os.MkdirAll("./dist", 0755)
	}
	app := router.Init()
	// json.NewEncoder(os.Stdout).Encode(app.Routes())
	go func() {
		if err := app.Start(conf.Addr()); err != nil {
			if err != http.ErrServerClosed { //Normal exit
				logs.Ctx(ctx).Error("Server Shutdown:", err)
			}
		}
	}()
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		if err != http.ErrServerClosed { //Normal exit
			logs.Ctx(ctx).Error("Server Shutdown:", err)
		}
	}
	logs.Ctx(ctx).Info("app quitted")
}

//
