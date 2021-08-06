package router

import (
	"blog/conf"
	"blog/internal/token"
	"blog/model"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zxysilent/utils"
)

var funcMap template.FuncMap

func init() {
	funcMap = template.FuncMap{"str2html": Str2html, "str2js": Str2js, "date": Date, "md5": Md5}
}

// midRecover 恢复中间件
func midRecover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				stack := make([]byte, 1<<10)
				length := runtime.Stack(stack, false)
				// stdlog.Println(string(stack[:length]))
				os.Stdout.Write(stack[:length])
				ctx.Error(err)
			}
		}()
		return next(ctx)
	}
}

// HTTPErrorHandler 全局错误捕捉
func HTTPErrorHandler(err error, ctx echo.Context) {
	if !ctx.Response().Committed {
		if he, ok := err.(*echo.HTTPError); ok {
			if he.Code == 404 {
				if strings.HasPrefix(ctx.Request().URL.Path, "/static") || strings.HasPrefix(ctx.Request().URL.Path, "/dist") {
					ctx.NoContent(404)
				} else if strings.HasPrefix(ctx.Request().URL.Path, "/api") || strings.HasPrefix(ctx.Request().URL.Path, "/adm") {
					ctx.JSON(utils.ErrSvr("系统错误", he.Message))
				} else {
					ctx.HTML(404, html404)
				}
			} else {
				ctx.JSON(utils.ErrSvr("系统错误", he.Message))
			}
		} else {
			ctx.JSON(utils.ErrSvr("系统错误", err.Error()))
		}
	}
}

// 跨越配置
var crosConfig = middleware.CORSConfig{
	AllowOrigins: []string{"*"},
	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
}

// TplRender is a custom html/template renderer for Echo framework
type TplRender struct {
	templates *template.Template
}

// Render renders a template document
func (t *TplRender) Render(w io.Writer, name string, data interface{}, ctx echo.Context) error {
	// 获取数据配置项
	if mp, is := data.(map[string]interface{}); is {
		mp["appjs"] = AppJsUrl
		mp["appcss"] = AppCssUrl
		mp["global"] = model.Gcfg()
	}
	//开发模式
	//每次强制读取模板
	//每次强制加载函数
	if conf.App.IsDev() {
		t.templates, _ = utils.LoadTmpl("./views", funcMap)
	}
	return t.templates.ExecuteTemplate(w, name, data)
}

// Str2html Convert string to template.HTML type.
func Str2html(raw string) template.HTML {
	return template.HTML(raw)
}

// Str2js Convert string to template.JS type.
func Str2js(raw string) template.JS {
	return template.JS(raw)
}

// Date Date
func Date(t time.Time, format string) string {
	return t.Format(format) //"2006-01-02 15:04:05"
}

// Md5 Md5
func Md5(str string) string {
	ctx := md5.New()
	ctx.Write([]byte(str))
	return hex.EncodeToString(ctx.Sum(nil))
}

// initRender 初始化模板和函数
func initRender() *TplRender {
	tpl, _ := utils.LoadTmpl("./views", funcMap)
	return &TplRender{
		templates: tpl,
	}
}

// midAuth 登录认证中间件
func midAuth(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		tokenRaw := ctx.FormValue(conf.App.TokenKey) // query/form 查找 token
		if tokenRaw == "" {
			tokenRaw = ctx.Request().Header.Get(echo.HeaderAuthorization) // header 查找token
			if tokenRaw == "" {
				return ctx.JSON(utils.ErrToken("请重新登陆", "token为空"))
			}
		}
		auth := token.Auth{}
		err := auth.Decode(tokenRaw, conf.App.TokenSecret)
		if err != nil {
			return ctx.JSON(utils.ErrToken("请重新登陆", err.Error()))
		}
		// 验证通过，保存信息
		ctx.Set("auth", auth)
		ctx.Set("uid", auth.Id)
		ctx.Set("rid", auth.RoleId)
		// 后续流程
		return next(ctx)
	}
}

// html404 404错误页面
const html404 = `<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8" />
        <meta http-equiv="X-UA-Compatible" content="IE=edge" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <title>404 Not Found</title>
        <style>
        body,html{background:#f0f0f0;overflow:hidden;font-family:Helvetica Neue,Helvetica,PingFang SC,Hiragino Sans GB,Microsoft YaHei,微软雅黑,Arial,sans-serif}
        .app{display:block;background:#fff;border-radius:4px;font-size:14px;position:relative;transition:all .2s ease-in-out}
        .app-body{padding:16px}
        .gap{padding-bottom:20px;padding-top:16px;margin-bottom:20px;margin-top:16px}
        .result-actions{margin-top:32px}
        @media screen and (max-width:480px){.result{width:100%}.result-extra{padding:18px 20px}}
        .result{width:72%;margin:0 auto;text-align:center}
        .result-title{margin-bottom:16px;color:#17233d;font-weight:500;font-size:24px;line-height:32px}
        .result-icon-success{background-color:#19be6b}
        .result-icon .icon{color:#fff;font-size:150px;border-radius:50%}
        .icon{display:inline-block;font-style:normal;font-weight:400;font-variant:normal;text-transform:none;text-rendering:optimizeLegibility;line-height:1;-webkit-font-smoothing:antialiased;-moz-osx-font-smoothing:grayscale;vertical-align:-.125em;text-align:center}
        .result-icon{display:inline-block;width:150px;border-radius:50%}
        .result-desc{margin-bottom:24px;color:#808695;font-size:14px;line-height:22px}
        .result-extra{padding:24px 40px;text-align:left;background:#f8f8f9;border-radius:4px}
        .result-actions .btn:not(:last-child){margin-right:8px}
        .btn{display:inline-block;margin-bottom:0;font-weight:400;text-align:center;vertical-align:middle;touch-action:manipulation;cursor:pointer;background-image:none;border:1px solid transparent;white-space:nowrap;-webkit-user-select:none;-ms-user-select:none;user-select:none;height:32px;line-height:32px;text-decoration-line:none;padding:0 15px;font-size:14px;border-radius:4px;transition:color .2s linear,background-color .2s linear,border .2s linear,box-shadow .2s linear;color:#515a6e;background-color:#fff;border-color:#dcdee2}
        .btn-success{color:#fff;background-color:#19be6b;border-color:#19be6b}
        .btn-success:hover{color:#fff;background-color:#47cb89;border-color:#47cb89}
        .btn-info{color:#fff;background-color:#2db7f5;border-color:#2db7f5}
        .btn-info:hover{color:#fff;background-color:#57c5f7;border-color:#57c5f7}
        .btn-warning{color:#fff;background-color:#f90;border-color:#f90}
        .btn-warning:hover{color:#fff;background-color:#ffad33;border-color:#ffad33}
        .btn-primary{color:#fff;background-color:#2d8cf0;border-color:#2d8cf0}
        .btn-primary:hover{color:#fff;background-color:#57a3f3;border-color:#57a3f3}
        .btn-default:hover{color:#57a3f3;border-color:#57a3f3}
        </style>
    </head>
    <body>
        <div class="app">
            <div class="app-body">
                <div class="gap">
                    <div class="result">
                        <svg class="result-icon">
                            <use xlink:href="#icon-404">
                                <svg id="icon-success" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M512 96C282.2 96 96 282.2 96 512s186.2 416 416 416 416-186.2 416-416S741.8 96 512 96z m213 301L457.6 665.6h-0.2c-3.4 3.4-12.6 11-23.2 11-7.6 0-16.2-4.2-23.4-11.4l-112-112c-3.2-3.2-3.2-8.2 0-11.4l35.6-35.6c1.6-1.6 3.6-2.4 5.6-2.4 2 0 4 0.8 5.6 2.4l88.8 88.8 244-245.8c1.6-1.6 3.6-2.4 5.6-2.4 2.2 0 4.2 0.8 5.6 2.4l35 36.2c3.6 3.4 3.6 8.4 0.4 11.6z" fill="#19be6b"></path></svg>
                                <svg id="icon-info" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M470.8 344.4c0-22.8 18.6-39.8 41-39.8 22.8 0 41.4 17 41.4 39.8s-18.6 40-41.4 40c-22.4 0-41-17.2-41-40z m2.8 71.4H550V704h-76.4V415.8z" fill="#2db7f5"></path><path d="M512 152c96.2 0 186.6 37.4 254.6 105.4S872 415.8 872 512s-37.4 186.6-105.4 254.6S608.2 872 512 872c-96.2 0-186.6-37.4-254.6-105.4S152 608.2 152 512s37.4-186.6 105.4-254.6S415.8 152 512 152m0-56C282.2 96 96 282.2 96 512s186.2 416 416 416 416-186.2 416-416S741.8 96 512 96z" fill="#2db7f5"></path></svg>
                                <svg id="icon-warn" class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M512 96c-229.8 0-416 186.2-416 416s186.2 416 416 416 416-186.2 416-416-186.2-416-416-416z m254.6 670.6c-68 68-158.4 105.4-254.6 105.4-96.2 0-186.6-37.4-254.6-105.4S152 608.2 152 512c0-96.2 37.4-186.6 105.4-254.6s158.4-105.4 254.6-105.4c96.2 0 186.6 37.4 254.6 105.4S872 415.8 872 512c0 96.2-37.4 186.6-105.4 254.6z" fill="#ff9900"></path><path d="M512 575c16.4 0 27.8-13.2 28.2-29.6L546.4 315.2c0.4-19.4-15.2-35.2-34.4-35.2-19.4 0-34.8 16-34.4 35.2l6.2 230.2c0.4 16.4 11.8 29.6 28.2 29.6zM512 634.4c-21.4 0-38.2 16.2-38.2 36.8s16.8 36.8 38.2 36.8 38.2-16.2 38.2-36.8-16.8-36.8-38.2-36.8z" fill="#ff9900"></path></svg>
                                <svg id="icon-error"  class="icon" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M512 96C282.2 96 96 282.2 96 512s186.2 416 416 416 416-186.2 416-416S741.8 96 512 96z m105.4 566.6L512 557.2l-105.4 105.4c-12.4 12.4-32.8 12.4-45.2 0-6.2-6.2-9.4-14.4-9.4-22.6 0-8.2 3.2-16.4 9.4-22.6l105.4-105.4-105.4-105.4c-6.2-6.2-9.4-14.4-9.4-22.6 0-8.2 3.2-16.4 9.4-22.6 12.4-12.4 32.8-12.4 45.2 0l105.4 105.4 105.4-105.4c12.4-12.4 32.8-12.4 45.2 0 12.4 12.4 12.4 32.8 0 45.2L557.2 512l105.4 105.4c12.4 12.4 12.4 32.8 0 45.2-12.4 12.6-32.8 12.6-45.2 0z" fill="#ed4014"></path></svg>
                                <svg id="icon-404"class="icon" viewBox="0 0 2486 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M1276.342857 15.798857c239.908571 0 397.897143 200.265143 409.014857 454.656l0.585143 31.012572c0 267.849143-160.914286 485.814857-409.6 485.814857-254.244571 0-412.964571-200.850286-423.643428-455.68l-0.585143-30.134857c0-267.702857 175.542857-485.668571 424.228571-485.668572z m-658.285714 1.316572c58.514286 0 117.028571 51.638857 117.028571 118.198857v434.907428h29.257143c43.885714 0 102.4 46.957714 102.4 105.472 0 59.099429-58.514286 104.301714-102.4 104.301715h-29.257143v81.334857c0 66.56-58.514286 118.198857-117.028571 118.198857s-117.028571-52.077714-117.028572-118.198857v-81.334857h-336.457142c-58.514286 0-117.028571-47.835429-117.028572-115.712 0-41.984 21.504-75.483429 46.08-107.52l12.434286-15.945143L486.4 92.16c14.628571-21.942857 29.257143-41.252571 58.514286-54.710857 14.628571-13.750857 43.885714-20.333714 73.142857-20.333714z m1623.771428 0c58.514286 0 117.028571 51.638857 117.028572 118.198857v434.907428h14.628571c58.514286 0 102.4 46.957714 102.4 105.472 0 59.099429-43.885714 104.301714-102.4 104.301715h-14.628571v81.334857c0 66.56-58.514286 118.198857-117.028572 118.198857-73.142857 0-117.028571-52.077714-117.028571-118.198857v-81.334857h-336.457143c-58.514286 0-117.028571-47.835429-117.028571-115.712 0-40.96 10.24-73.581714 30.427428-104.740572l13.458286-18.724571 380.342857-448.658286c11.702857-8.777143 21.065143-17.115429 28.964572-24.868571l29.549714-29.842286c29.257143-13.750857 58.514286-20.333714 87.771428-20.333714z m-965.485714 226.304c-58.514286 0-87.771429 24.283429-117.028571 68.315428-25.6 39.058286-51.2 93.622857-57.197715 158.72l-1.316571 31.012572c0 77.677714 29.257143 143.36 58.514286 188.562285 29.257143 44.763429 58.514286 69.632 117.028571 69.632 43.885714 0 87.771429-24.576 117.028572-68.754285 25.6-39.204571 39.936-93.769143 43.154285-158.427429l0.731429-31.012571c0-76.8-14.628571-142.482286-43.885714-187.977143-29.257143-45.202286-73.142857-70.070857-117.028572-70.070857z m-775.314286 143.945142l-146.285714 182.857143h146.285714v-182.857143z m1623.771429 0l-160.914286 182.857143h160.914286v-182.857143z" fill="#ed4014"></path></svg>
                            </use>
                        </svg>
                        <div class="result-title">404 Not Found</div>
                        <div class="result-desc">
                            <div>页面不存在,请确认后访问。</div>
                        </div>
                        <div class="result-extra">
                            <div>暂无额外信息</div>
                        </div>
                        <div class="result-actions">
                            <div>
                                <a href="/" class="btn btn-success" title="返回主页"> 返回主页 </a>
                                <a href="javascript:history.go(-1)" class="btn btn-primary" title="返回主页"> 上一页 </a>
                                <a href="https://github.com/zxysilent" target="_blank" title="查看作者" class="btn btn-default"> 查看作者 </a>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </body>
</html>
`
