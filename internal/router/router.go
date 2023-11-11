package router

import (
	"blog/conf"
	"blog/internal/model"
	"blog/internal/repo"
	"blog/internal/service/view"
	"blog/internal/utils"
	"blog/pkg/page"
	"blog/pkg/token"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"regexp"
	"runtime/debug"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zxysilent/logs"
)

// Init 入口
func Init() *echo.Echo {
	engine := echo.New()
	// 初始渲染引擎
	engine.Renderer = initRender()
	// 恢复 日志记录
	engine.Use(midRecover, midLogger)
	// engine.Use(middleware.LoggerWithConfig(logConfig))
	// 跨域设置
	engine.Use(middleware.CORSWithConfig(crosConfig))
	engine.Use(middleware.RateLimiterWithConfig(rateConfig))
	// 不显示横幅
	engine.HideBanner = true
	// 自定义错误处理
	engine.HTTPErrorHandler = HTTPErrorHandler
	engine.Debug = conf.IsDebug()
	// ico
	engine.File("/favicon.ico", "favicon.ico")
	// 注册文档
	RegDocs(engine)
	// 静态目录
	engine.Static("/dist", "dist")
	engine.Static("/static", "static")
	engine.File("/favicon.ico", "favicon.ico") // ico
	engine.GET("/", view.ViewIndex)            // 首页
	engine.GET("/archives", view.ViewArchives) // 归档
	engine.GET("/tags", view.ViewTags)
	// 标签
	engine.GET("/tags/:tag", view.ViewTagPost)    // 具体某个标签
	engine.GET("/cates/:cate", view.ViewCatePost) // 分类
	engine.GET("/notes/*", view.ViewNote)         // 笔记
	engine.GET("/search", view.ViewSearch)        // 搜索
	engine.GET("/about", view.ViewAbout)          // 关于
	engine.GET("/links", view.ViewLinks)          // 友链
	engine.GET("/pages/*", view.ViewPage)         // 具体某个页面
	engine.GET("/posts/*", view.ViewPost)         // 具体某个文章
	engine.File("/mgmt*", "dist/index.html")      // 前后端分离页面
	engine.GET("/mgmt", func(ctx echo.Context) error {
		return ctx.Redirect(http.StatusFound, "/mgmt/")
	})
	engine.GET("/ping", func(ctx echo.Context) error {
		return ctx.JSON(utils.Succ("pong", "service is running"))
	})
	// 后台登录
	engine.GET("/login.html", func(ctx echo.Context) error {
		// 302 临时重定向
		return ctx.Redirect(302, "/mgmt/")
	})
	apiRouter(engine)
	return engine
}

// func init() {
// 	logs.SetLevel(logs.LDEBUG)
// 	logs.SetCaller(true)
// 	logs.SetConsole(true)
// 	pool = &sync.Pool{
// 		New: func() interface{} {
// 			return bytes.NewBuffer(make([]byte, 512))
// 		},
// 	}
// }

func midRecover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				logs.With().Caller(false).Str("trace", ctx.Request().Header.Get(echo.HeaderXRequestID)).Str("caller", "echo").Raw("panic", debug.Stack()).Err(err).Error("recover")
				ctx.Error(err)
			}
		}()
		return next(ctx)
	}
}

// midLogger 中间件-日志记录
func midLogger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) (err error) {
		start := time.Now()
		rid := ctx.Request().Header.Get(echo.HeaderXRequestID)
		if rid == "" {
			rid = logs.TraceId()
		}
		ctx.Request().Header.Set(echo.HeaderXRequestID, rid)
		ctx.Response().Header().Set(echo.HeaderXRequestID, rid)
		if err = next(ctx); err != nil {
			ctx.Error(err)
		}
		stop := time.Now()
		logs.With().Caller(false).Str("trace", rid).Str("caller", "echo").Str("ip", ctx.RealIP()).Str(ctx.Request().Method, ctx.Request().RequestURI).Str("span", stop.Sub(start).String()).Debug()
		return
	}
}

// HTTPErrorHandler 全局错误捕捉
func HTTPErrorHandler(err error, ctx echo.Context) {
	if !ctx.Response().Committed {
		if he, ok := err.(*echo.HTTPError); ok {
			if he.Code == 404 {
				if strings.HasPrefix(ctx.Request().URL.Path, "/static") || strings.HasPrefix(ctx.Request().URL.Path, "/dist") {
					ctx.NoContent(404)
				} else if strings.HasPrefix(ctx.Request().URL.Path, "/api") {
					ctx.JSON(utils.ErrSvr("系统错误", he.Message))
				} else {
					ctx.HTML(404, page.Err404().SetExtra("404").Render())
				}
			} else {
				if strings.HasPrefix(ctx.Request().URL.Path, "/api") {
					ctx.JSON(utils.ErrSvr("系统错误", he.Message))
				} else {
					ctx.HTML(502, page.Error().SetExtra(fmt.Sprint(he.Message)).Render())
				}
			}
		} else {
			ctx.HTML(502, page.Error().SetExtra(err.Error()).Render())
		}
	}
}

// 跨越配置
var crosConfig = middleware.CORSConfig{
	AllowOrigins: []string{"*"},
	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "user-token", "accept-serial", "token"},
	AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
}
var rateConfig = middleware.RateLimiterConfig{
	Skipper: func(ctx echo.Context) bool {
		path := ctx.Request().URL.Path
		if strings.HasPrefix(path, "/api") {
			return false
		}
		return true
	},
	Store: middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{Rate: 20, Burst: 5, ExpiresIn: 3 * time.Minute}),
	IdentifierExtractor: func(ctx echo.Context) (string, error) {
		path := ctx.Request().URL.Path
		return path, nil
	},
	ErrorHandler: func(ctx echo.Context, err error) error {
		return ctx.JSON(utils.ErrOpt("请求错误,请稍后重试"))
	},
	DenyHandler: func(ctx echo.Context, identifier string, err error) error {
		return ctx.JSON(utils.ErrOpt("当前请求过多,请稍后重试", identifier))
	},
}

// TplRender is a custom html/template renderer for Echo framework
type TplRender struct {
	templates *template.Template
}

// Render renders a template document
func (tmpl *TplRender) Render(w io.Writer, name string, data interface{}, ctx echo.Context) error {
	// 获取数据配置项
	if mp, is := data.(map[string]interface{}); is {
		mp["appjs"] = AppJsUrl
		mp["appcss"] = AppCssUrl
		mp["basic"] = repo.CfgBasic()
		mp["blog"] = repo.CfgBlog()
	}
	//开发模式
	//每次强制读取模板
	//每次强制加载函数
	if conf.IsDebug() {
		tmpl.templates, _ = utils.LoadTmpl("./views", funcMap)
	}
	return tmpl.templates.ExecuteTemplate(w, name, data)
}

var funcMap template.FuncMap

func init() {
	funcMap = template.FuncMap{"str2html": Str2html, "str2js": Str2js, "date": Date, "md5": Md5, "unix": Unix}
}

// Str2html Convert string to template.HTML type.
func Str2html(raw string) template.HTML {
	return template.HTML(raw)
}

func Unix(ms int64, fmt string) string {
	return time.UnixMilli(ms).Format(fmt)
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

// 初始化模板和函数
func initRender() *TplRender {
	tmpl, err := utils.LoadTmpl("./views", funcMap)
	if err != nil {
		logs.Error("模板加载", err.Error())
		os.Exit(0)
	}
	return &TplRender{
		templates: tmpl,
	}
}

func midToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		tokenRaw := ctx.Request().Header.Get(echo.HeaderAuthorization) // header 查找token--方便前端约定
		if tokenRaw == "" {
			tokenRaw = ctx.FormValue(conf.App.TokenKey) // query/form 查找 token
			if tokenRaw == "" {
				cookie, _ := ctx.Cookie(conf.App.TokenKey)
				if cookie != nil {
					tokenRaw = cookie.Value
				}
			}
		}
		if tokenRaw == "" {
			ctx.JSON(utils.ErrToken("请重新登陆", "未发现token"))
			return nil
		}
		auth := token.Token[token.Auth]{}
		err := auth.Decode(tokenRaw, conf.App.TokenSecret)
		if err == nil {
			claim := auth.Claim
			ctx.Set("auth", claim)
			ctx.Set("uid", claim.UserId)
			ctx.Set("rid", claim.RoleId)
			if claim.RoleId != model.RootRoleId {
				// 判断权限
				allow := false
				grants, er := repo.RoleGrantAll(claim.RoleId)
				logs.Debug(er)
				path := ctx.Request().URL.Path
				for _, grant := range grants {
					for _, route := range grant.Routes {
						if routeMatch(path, route) {
							allow = true
							break
						}
					}
				}
				if !allow {
					return ctx.JSON(utils.ErrDeny("没有权限"))
				}
			}
		} else {
			return ctx.JSON(utils.ErrToken("请重新登陆", "token验证失败,"+err.Error()))
		}
		return next(ctx)
	}
}
func midAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		tokenRaw := ctx.Request().Header.Get(echo.HeaderAuthorization) // header 查找token--方便前端约定
		if tokenRaw == "" {
			tokenRaw = ctx.FormValue(conf.App.TokenKey) // query/form 查找 token
			if tokenRaw == "" {
				cookie, _ := ctx.Cookie(conf.App.TokenKey)
				if cookie != nil {
					tokenRaw = cookie.Value
				}
			}
		}
		if tokenRaw == "" {
			ctx.JSON(utils.ErrToken("请重新登陆", "未发现token"))
			return nil
		}
		auth := token.Token[token.Auth]{}
		err := auth.Decode(tokenRaw, conf.App.TokenSecret)
		if err == nil {
			claim := auth.Claim
			if !claim.Admin {
				return ctx.JSON(utils.ErrDeny("没有权限"))
			}
			ctx.Set("auth", claim)
			ctx.Set("uid", claim.UserId)
			ctx.Set("rid", claim.RoleId)
		} else {
			return ctx.JSON(utils.ErrToken("请重新登陆", "token验证失败,"+err.Error()))
		}
		return next(ctx)
	}
}
func routeMatch(path string, route string) bool {
	if path == route {
		return true
	}
	logs.Debug(path, route)
	route = strings.Replace(route, "/*", "/.*", -1)
	re := regexp.MustCompile(`:[^/]+`)
	route = re.ReplaceAllString(route, "$1[^/]+$2")
	res, _ := regexp.MatchString("^"+route+"$", path)
	return res
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
        body,html{background:#f0f0f0;font-family:Helvetica Neue,Helvetica,PingFang SC,Hiragino Sans GB,Microsoft YaHei,微软雅黑,Arial,sans-serif}
        .app{display:block;background:#fff;border-radius:4px;font-size:14px;position:relative;transition:all .2s ease-in-out}
        .app-body{padding:16px}
        .gap{padding-bottom:20px;padding-top:16px;margin-bottom:20px;margin-top:16px}
        .result-actions{margin-top:32px}
        @media screen and (max-width:480px){.result{width:100%}.result-extra{padding:18px 20px}}
        .result{width:72%;margin:0 auto;text-align:center}
        .result-title{margin-bottom:16px;color:#17233d;font-weight:500;font-size:24px;line-height:32px}
        .result-icon{display:inline-block;width:150px;border-radius:50%}
        .result-desc{margin-bottom:24px;color:#808695;font-size:14px;line-height:22px}
        .result-extra{padding:24px 40px;text-align:center;background:#f8f8f9;border-radius:4px}
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
        .footer{margin:28px 0 24px;padding:0 28px;text-align:center}
        .footer-links{margin-bottom:8px}
        .footer-copyright{color:#808695;font-size:14px}
        .footer-links a{margin-right:20px;margin-left:20px;font-size:14px;color:#808695;-webkit-transition:all .2s ease-in-out;transition:all .2s ease-in-out}
        a{color:#2d8cf0;background:0 0;text-decoration:none;outline:0;cursor:pointer;transition:color .2s ease}
        </style>
    </head>
    <body>
        <div class="app">
            <div class="app-body">
                <div class="gap">
                    <div class="result">
                        <svg class="result-icon">
                            <use xlink:href="#icon-404">
                                <svg id="icon-404"viewBox="0 0 2486 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M1276.342857 15.798857c239.908571 0 397.897143 200.265143 409.014857 454.656l0.585143 31.012572c0 267.849143-160.914286 485.814857-409.6 485.814857-254.244571 0-412.964571-200.850286-423.643428-455.68l-0.585143-30.134857c0-267.702857 175.542857-485.668571 424.228571-485.668572z m-658.285714 1.316572c58.514286 0 117.028571 51.638857 117.028571 118.198857v434.907428h29.257143c43.885714 0 102.4 46.957714 102.4 105.472 0 59.099429-58.514286 104.301714-102.4 104.301715h-29.257143v81.334857c0 66.56-58.514286 118.198857-117.028571 118.198857s-117.028571-52.077714-117.028572-118.198857v-81.334857h-336.457142c-58.514286 0-117.028571-47.835429-117.028572-115.712 0-41.984 21.504-75.483429 46.08-107.52l12.434286-15.945143L486.4 92.16c14.628571-21.942857 29.257143-41.252571 58.514286-54.710857 14.628571-13.750857 43.885714-20.333714 73.142857-20.333714z m1623.771428 0c58.514286 0 117.028571 51.638857 117.028572 118.198857v434.907428h14.628571c58.514286 0 102.4 46.957714 102.4 105.472 0 59.099429-43.885714 104.301714-102.4 104.301715h-14.628571v81.334857c0 66.56-58.514286 118.198857-117.028572 118.198857-73.142857 0-117.028571-52.077714-117.028571-118.198857v-81.334857h-336.457143c-58.514286 0-117.028571-47.835429-117.028571-115.712 0-40.96 10.24-73.581714 30.427428-104.740572l13.458286-18.724571 380.342857-448.658286c11.702857-8.777143 21.065143-17.115429 28.964572-24.868571l29.549714-29.842286c29.257143-13.750857 58.514286-20.333714 87.771428-20.333714z m-965.485714 226.304c-58.514286 0-87.771429 24.283429-117.028571 68.315428-25.6 39.058286-51.2 93.622857-57.197715 158.72l-1.316571 31.012572c0 77.677714 29.257143 143.36 58.514286 188.562285 29.257143 44.763429 58.514286 69.632 117.028571 69.632 43.885714 0 87.771429-24.576 117.028572-68.754285 25.6-39.204571 39.936-93.769143 43.154285-158.427429l0.731429-31.012571c0-76.8-14.628571-142.482286-43.885714-187.977143-29.257143-45.202286-73.142857-70.070857-117.028572-70.070857z m-775.314286 143.945142l-146.285714 182.857143h146.285714v-182.857143z m1623.771429 0l-160.914286 182.857143h160.914286v-182.857143z" fill="#ed4014"></path></svg>
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
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
        <div class="footer">
            <div class="footer-links">
                <a href="https://blog.zxysilent.com"  target="_blank" title="blog">blog</a>
                <a href="https://github.com/zxysilent/blog" target="_blank" title="github">github</a>
                <a href="https://github.com/zxysilent" target="_blank" title="查看作者">作者</a>
            </div>
            <div class="footer-copyright">Copyright &copy; <script>document.write(new Date().getFullYear());</script>&nbsp;<a target="_blank" href="https://github.com/zxysilent">github.com/zxysilent</a></div>
        </div>
    </body>
</html>`
