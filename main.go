//go:generate go test -timeout 30s -run ^TestGenRouter$ blog/internal/utils -count=1 -v -msg=路径相对于utils -rp=../router/api.go -rs=../service
package main

import (
	"blog/conf"
	"blog/internal/repo"
	"blog/internal/router"
	"blog/internal/service"
	"blog/internal/utils"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/zxysilent/logs"
)

var (
	version   string //
	buildSha  string //
	buildTime string //
)

// @Title blog’s Api文档
// @Version 1.0
// @Description 凭证传递方式包括 get、post、header、cookie
// @Description /api/* 前后端api
//
// @securityDefinitions.apikey	KeyAuth
// @in							header
// @name						Authorization
// @description					Description for what is this security definition being used
func main() {
	fmt.Println("Golang version :", runtime.Version())
	fmt.Println("Git commit sha :", buildSha)
	fmt.Println("Build     time :", buildTime)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGSEGV)
	defer stop()
	ctx = logs.TraceCtx(ctx)
	defer logs.Close()
	logs.Ctx(ctx).Info("app initializing")
	conf.Init(ctx)
	repo.Init(ctx)
	defer repo.Close()
	service.Init(ctx)
	logs.Ctx(ctx).Info("app running")
	if !utils.FileExist("./static") {
		os.MkdirAll("./static", 0755)
	}
	if !utils.FileExist("./dist") {
		os.MkdirAll("./dist", 0755)
	}
	app := router.Init()
	go func() {
		if err := app.Start(conf.Addr()); err != nil {
			if err != http.ErrServerClosed { //Normal exit
				logs.Ctx(ctx).Err(err).Error("Server Shutdown")
			}
		}
	}()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := app.Shutdown(ctx); err != nil {
		if err != http.ErrServerClosed { //Normal exit
			logs.Ctx(ctx).Err(err).Error("Server Shutdown")
		}
	}
	logs.Ctx(ctx).Info("app quitted")
}
