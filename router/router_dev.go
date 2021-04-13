// +build !prod

package router

import (

	// docs
	_ "blog/docs"
	"bytes"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/zxysilent/logs"
)

const AppJsUrl = "/static/js/app.js"
const AppCssUrl = "/static/css/app.css"

var pool *sync.Pool

func init() {
	logs.SetLevel(logs.DEBUG)
	logs.SetCallInfo(true)
	logs.SetConsole(true)
	pool = &sync.Pool{
		New: func() interface{} {
			return bytes.NewBuffer(make([]byte, 512))
		},
	}
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
		buf.WriteString("\t" + ctx.RealIP())
		buf.WriteString("\t" + ctx.Request().Method + "：" + ctx.Request().RequestURI)
		buf.WriteString(" - " + stop.Sub(start).String())
		// 这里仅仅记录的是 url 访问记录
		// 开发模式直接输出
		// 生产模式中间层会记录
		// os.Stdout.Write(buf.Bytes())
		logs.Debug(buf.String())
		return
	}
}
