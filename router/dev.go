// +build !prod

package router

import (

	// docs
	_ "blog/docs"
	"bytes"
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

const AppJsUrl = "/static/js/app.js"
const AppCssUrl = "/static/css/app.css"

func init() {
	log.SetFlags(log.Lshortfile | log.Ltime)
	log.SetPrefix("[app] ")
}

// RegDocs 注册文档
// dev[开发] 模式需要文档
func RegDocs(engine *echo.Echo) {
	docUrl := echoSwagger.URL("/swagger/doc.json")
	engine.GET("/swagger/*", echoSwagger.EchoWrapHandler(docUrl))
}

// midLogger 中间件-日志记录
func midLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) (err error) {
		start := time.Now()
		if err = next(ctx); err != nil {
			ctx.Error(err)
		}
		stop := time.Now()
		buf := pool.Get().(*bytes.Buffer)
		buf.Reset()
		defer pool.Put(buf)
		buf.WriteString("[" + start.Format("2006-01-02 15:04:05") + "] ")
		buf.WriteString("\tip：" + ctx.RealIP())
		buf.WriteString("\tmethod：" + ctx.Request().Method)
		buf.WriteString("\tpath：" + ctx.Request().RequestURI)
		buf.WriteString("\tspan：" + stop.Sub(start).String())
		buf.WriteString("\n")
		// 开发模式直接输出
		// 生产模式中间层会记录
		os.Stdout.Write(buf.Bytes())
		return
	}
}
