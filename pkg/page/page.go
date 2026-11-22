package page

import (
	"bytes"
	"html/template"
)

type _page struct {
	Kind   string //类型
	Circle bool   //图标
	Title  string //标题
	Desc   HTML   //描述
	Extra  HTML   //额外信息
	ExecJs JS     //可能执行的js
}
type JS = template.JS
type HTML = template.HTML

const (
	k404   string = "404"
	kSucc         = "succ"
	kError        = "error"
	kInfo         = "info"
	kWarn         = "warn"
)

type context struct {
	kind   string //类型
	circle bool   //图标
	title  string //标题
	desc   HTML   //描述
	extra  HTML   //额外信息
	execJs JS     //可能执行的js
}

func newPage(kind string) *context {
	ctx := &context{}
	ctx.circle = true
	ctx.kind = kind
	return ctx
}

// SetCircle 设置图标样式
// b bool circle
func (ctx *context) SetCircle(b bool) *context {
	ctx.circle = b
	return ctx
}

// SetTitle 设置标题
// title string 标题
func (ctx *context) SetTitle(title string) *context {
	ctx.title = title
	return ctx
}

// SetDesc 设置描述
// desc string 描述
func (ctx *context) SetDesc(desc string) *context {
	ctx.desc = HTML(desc)
	return ctx
}

// SetExtra 设置附加信息
// extra string 附加信息
func (ctx *context) SetExtra(extra string) *context {
	ctx.extra = HTML(extra)
	return ctx
}

// SetExecJs 设置执行js
// execJs string 执行js
func (ctx *context) SetExecJs(execJs string) *context {
	ctx.execJs = JS(execJs)
	return ctx
}
func (ctx *context) Render() string {
	var buf bytes.Buffer
	data := _page{
		Kind:   ctx.kind,
		Circle: ctx.circle,
		Title:  ctx.title,
		Desc:   ctx.desc,
		Extra:  ctx.extra,
		ExecJs: ctx.execJs,
	}
	tmpl.Execute(&buf, data)
	return buf.String()
}
func init() {
	tmpl, _ = template.New("tmpl").Parse(tmplHtml)
}

// Err404 404页面
func Err404() *context {
	ctx := newPage(k404)
	ctx.title = "404 Not Found"
	ctx.desc = "页面不存在,请确认后访问。"
	ctx.extra = "暂无额外信息"
	return ctx
}

// Succ 成功页面
func Succ() *context {
	ctx := newPage(kSucc)
	ctx.title = "成功"
	ctx.desc = "成功"
	ctx.extra = "暂无额外信息"
	return ctx
}

// Info 信息页面
func Info() *context {
	ctx := newPage(kInfo)
	ctx.title = "提示"
	ctx.desc = "提示信息"
	ctx.extra = "暂无额外信息"
	return ctx
}

// Warn 警告页面
func Warn() *context {
	ctx := newPage(kWarn)
	ctx.title = "警告"
	ctx.desc = "警告信息"
	ctx.extra = "暂无额外信息"
	return ctx
}

// Error 错误页面
func Error() *context {
	ctx := newPage(kError)
	ctx.title = "错误"
	ctx.desc = "错误"
	ctx.extra = "暂无额外信息"
	return ctx
}

var tmpl *template.Template
var tmplHtml = `<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8" /><meta http-equiv="X-UA-Compatible" content="IE=edge" /><meta name="viewport" content="width=device-width, initial-scale=1.0" /><title>{{.Title}}</title><style>body,html{background:#f0f0f0;font-family:Helvetica Neue,Helvetica,PingFang SC,Hiragino Sans GB,Microsoft YaHei,微软雅黑,Arial,sans-serif}.app{display:block;background:#fff;border-radius:4px;font-size:14px;position:relative;transition:all .2s ease-in-out}.app-body{padding:16px}.gap{padding-bottom:20px;padding-top:16px;margin-bottom:20px;margin-top:16px}.result-actions{margin-top:32px}@media screen and (max-width:480px){.result{width:100%}.result-extra{padding:18px 20px}}.result{width:72%;margin:0 auto;text-align:center}.result-title{margin-bottom:16px;color:#17233d;font-weight:500;font-size:24px;line-height:32px}.result-icon{display:inline-block;width:150px;border-radius:50%}.result-desc{margin-bottom:24px;color:#808695;font-size:14px;line-height:22px}.result-extra{padding:24px 40px;text-align:center;background:#f8f8f9;border-radius:4px}.result-actions .btn:not(:last-child){margin-right:8px}.btn{display:inline-block;margin-bottom:0;font-weight:400;text-align:center;vertical-align:middle;touch-action:manipulation;cursor:pointer;background-image:none;border:1px solid transparent;white-space:nowrap;-webkit-user-select:none;-ms-user-select:none;user-select:none;height:32px;line-height:32px;text-decoration-line:none;padding:0 15px;font-size:14px;border-radius:4px;transition:color .2s linear,background-color .2s linear,border .2s linear,box-shadow .2s linear;color:#515a6e;background-color:#fff;border-color:#dcdee2}.btn-success{color:#fff;background-color:#19be6b;border-color:#19be6b}.btn-success:hover{color:#fff;background-color:#47cb89;border-color:#47cb89}.btn-info{color:#fff;background-color:#2db7f5;border-color:#2db7f5}.btn-info:hover{color:#fff;background-color:#57c5f7;border-color:#57c5f7}.btn-warning{color:#fff;background-color:#f90;border-color:#f90}.btn-warning:hover{color:#fff;background-color:#ffad33;border-color:#ffad33}.btn-primary{color:#fff;background-color:#2d8cf0;border-color:#2d8cf0}.btn-primary:hover{color:#fff;background-color:#57a3f3;border-color:#57a3f3}.btn-default:hover{color:#57a3f3;border-color:#57a3f3}.footer{margin:28px 0 24px;padding:0 28px;text-align:center}.footer-links{margin-bottom:8px}.footer-copyright{color:#808695;font-size:14px}.footer-links a{margin-right:20px;margin-left:20px;font-size:14px;color:#808695;-webkit-transition:all .2s ease-in-out;transition:all .2s ease-in-out}a{color:#2d8cf0;background:0 0;text-decoration:none;outline:0;cursor:pointer;transition:color .2s ease}</style></head><body><div class="app"><div class="app-body"><div class="gap"><div class="result"><svg class="result-icon"><use xlink:href="#icon-{{.Kind}}">{{if .Circle}}<!-- circle -->{{if eq .Kind "succ"}}<svg id="icon-succ" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M725.2 385.8L690 349.6c-1.4-1.6-3.6-2.4-5.6-2.4-2.2 0-4.2 0.8-5.6 2.4l-244 245.8-88.8-88.8c-1.6-1.6-3.6-2.4-5.6-2.4-2 0-4 0.8-5.6 2.4l-35.6 35.6c-3.2 3.2-3.2 8.2 0 11.4l112 112c7.2 7.2 16 11.4 23.4 11.4 10.6 0 19.8-7.8 23.2-11h0.2l267.4-268.8c2.8-3.4 2.8-8.4-0.2-11.4z" fill="#19be6b"></path><path d="M512 152c96.2 0 186.6 37.4 254.6 105.4S872 415.8 872 512s-37.4 186.6-105.4 254.6S608.2 872 512 872c-96.2 0-186.6-37.4-254.6-105.4S152 608.2 152 512s37.4-186.6 105.4-254.6S415.8 152 512 152m0-56C282.2 96 96 282.2 96 512s186.2 416 416 416 416-186.2 416-416S741.8 96 512 96z" fill="#19be6b"></path></svg>{{else if eq .Kind "info"}}<svg id="icon-info" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M470.8 344.4c0-22.8 18.6-39.8 41-39.8 22.8 0 41.4 17 41.4 39.8s-18.6 40-41.4 40c-22.4 0-41-17.2-41-40z m2.8 71.4H550V704h-76.4V415.8z" fill="#2db7f5"></path><path d="M512 152c96.2 0 186.6 37.4 254.6 105.4S872 415.8 872 512s-37.4 186.6-105.4 254.6S608.2 872 512 872c-96.2 0-186.6-37.4-254.6-105.4S152 608.2 152 512s37.4-186.6 105.4-254.6S415.8 152 512 152m0-56C282.2 96 96 282.2 96 512s186.2 416 416 416 416-186.2 416-416S741.8 96 512 96z" fill="#2db7f5"></path></svg>{{else if eq .Kind "warn"}}<svg id="icon-warn" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M512 96c-229.8 0-416 186.2-416 416s186.2 416 416 416 416-186.2 416-416-186.2-416-416-416z m254.6 670.6c-68 68-158.4 105.4-254.6 105.4-96.2 0-186.6-37.4-254.6-105.4S152 608.2 152 512c0-96.2 37.4-186.6 105.4-254.6s158.4-105.4 254.6-105.4c96.2 0 186.6 37.4 254.6 105.4S872 415.8 872 512c0 96.2-37.4 186.6-105.4 254.6z" fill="#ff9900"></path><path d="M512 575c16.4 0 27.8-13.2 28.2-29.6L546.4 315.2c0.4-19.4-15.2-35.2-34.4-35.2-19.4 0-34.8 16-34.4 35.2l6.2 230.2c0.4 16.4 11.8 29.6 28.2 29.6zM512 634.4c-21.4 0-38.2 16.2-38.2 36.8s16.8 36.8 38.2 36.8 38.2-16.2 38.2-36.8-16.8-36.8-38.2-36.8z" fill="#ff9900"></path></svg>{{else if eq .Kind "error"}}<svg id="icon-error"  viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M662.6 617.4L557.2 512l105.4-105.4c12.4-12.4 12.4-32.8 0-45.2-12.4-12.4-32.8-12.4-45.2 0L512 466.8l-105.4-105.4c-12.4-12.4-31.2-14.2-45.2 0-14.2 14.2-12 33.2 0 45.2l105.4 105.4-105.4 105.4c-13.4 13.4-12.8 32.6 0 45.2 12.8 12.8 32.8 12.4 45.2 0l105.4-105.4 105.4 105.4c12.4 12.4 32.8 12.4 45.2 0 12.6-12.4 12.6-32.8 0-45.2z" fill="#ed4014"></path><path d="M512 152c96.2 0 186.6 37.4 254.6 105.4S872 415.8 872 512s-37.4 186.6-105.4 254.6S608.2 872 512 872c-96.2 0-186.6-37.4-254.6-105.4S152 608.2 152 512s37.4-186.6 105.4-254.6S415.8 152 512 152m0-56C282.2 96 96 282.2 96 512s186.2 416 416 416 416-186.2 416-416S741.8 96 512 96z" fill="#ed4014"></path></svg>{{end}}{{else}}<!-- not circle -->{{if eq .Kind "succ"}}<svg id="icon-succ" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M512 96C282.2 96 96 282.2 96 512s186.2 416 416 416 416-186.2 416-416S741.8 96 512 96z m213 301L457.6 665.6h-0.2c-3.4 3.4-12.6 11-23.2 11-7.6 0-16.2-4.2-23.4-11.4l-112-112c-3.2-3.2-3.2-8.2 0-11.4l35.6-35.6c1.6-1.6 3.6-2.4 5.6-2.4 2 0 4 0.8 5.6 2.4l88.8 88.8 244-245.8c1.6-1.6 3.6-2.4 5.6-2.4 2.2 0 4.2 0.8 5.6 2.4l35 36.2c3.6 3.4 3.6 8.4 0.4 11.6z" fill="#19be6b"></path></svg>{{else if eq .Kind "info"}}<svg id="icon-info" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M512 96C282.2 96 96 282.2 96 512s186.2 416 416 416 416-186.2 416-416S741.8 96 512 96z m38 608h-76.4V415.8H550V704z m-38.2-319.6c-22.6 0-41-17.2-41-40s18.6-39.8 41-39.8c22.8 0 41.4 17 41.4 39.8s-18.6 40-41.4 40z" fill="#2db7f5"></path></svg>{{else if eq .Kind "warn"}}<svg id="icon-warn" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M512 96C282.2 96 96 282.2 96 512s186.2 416 416 416 416-186.2 416-416S741.8 96 512 96z m34.4 219.2l-6.2 230.2c-0.4 16.4-11.8 29.6-28.2 29.6s-27.8-13.2-28.2-29.6l-6.2-230.2c-0.4-19.2 15-35.2 34.4-35.2 19.2 0 34.8 15.8 34.4 35.2zM512 708c-21.4 0-38.2-16.2-38.2-36.8s16.8-36.8 38.2-36.8c21.4 0 38.2 16.2 38.2 36.8S533.4 708 512 708z" fill="#ff9900"></path></svg>{{else if eq .Kind "error"}}<svg id="icon-error" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M512 96C282.2 96 96 282.2 96 512s186.2 416 416 416 416-186.2 416-416S741.8 96 512 96z m105.4 566.6L512 557.2l-105.4 105.4c-12.4 12.4-32.8 12.4-45.2 0-6.2-6.2-9.4-14.4-9.4-22.6 0-8.2 3.2-16.4 9.4-22.6l105.4-105.4-105.4-105.4c-6.2-6.2-9.4-14.4-9.4-22.6 0-8.2 3.2-16.4 9.4-22.6 12.4-12.4 32.8-12.4 45.2 0l105.4 105.4 105.4-105.4c12.4-12.4 32.8-12.4 45.2 0 12.4 12.4 12.4 32.8 0 45.2L557.2 512l105.4 105.4c12.4 12.4 12.4 32.8 0 45.2-12.4 12.6-32.8 12.6-45.2 0z" fill="#ed4014"></path></svg>{{end}}{{end}}{{if eq .Kind "404"}}<svg id="icon-404"viewBox="0 0 2486 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M1276.342857 15.798857c239.908571 0 397.897143 200.265143 409.014857 454.656l0.585143 31.012572c0 267.849143-160.914286 485.814857-409.6 485.814857-254.244571 0-412.964571-200.850286-423.643428-455.68l-0.585143-30.134857c0-267.702857 175.542857-485.668571 424.228571-485.668572z m-658.285714 1.316572c58.514286 0 117.028571 51.638857 117.028571 118.198857v434.907428h29.257143c43.885714 0 102.4 46.957714 102.4 105.472 0 59.099429-58.514286 104.301714-102.4 104.301715h-29.257143v81.334857c0 66.56-58.514286 118.198857-117.028571 118.198857s-117.028571-52.077714-117.028572-118.198857v-81.334857h-336.457142c-58.514286 0-117.028571-47.835429-117.028572-115.712 0-41.984 21.504-75.483429 46.08-107.52l12.434286-15.945143L486.4 92.16c14.628571-21.942857 29.257143-41.252571 58.514286-54.710857 14.628571-13.750857 43.885714-20.333714 73.142857-20.333714z m1623.771428 0c58.514286 0 117.028571 51.638857 117.028572 118.198857v434.907428h14.628571c58.514286 0 102.4 46.957714 102.4 105.472 0 59.099429-43.885714 104.301714-102.4 104.301715h-14.628571v81.334857c0 66.56-58.514286 118.198857-117.028572 118.198857-73.142857 0-117.028571-52.077714-117.028571-118.198857v-81.334857h-336.457143c-58.514286 0-117.028571-47.835429-117.028571-115.712 0-40.96 10.24-73.581714 30.427428-104.740572l13.458286-18.724571 380.342857-448.658286c11.702857-8.777143 21.065143-17.115429 28.964572-24.868571l29.549714-29.842286c29.257143-13.750857 58.514286-20.333714 87.771428-20.333714z m-965.485714 226.304c-58.514286 0-87.771429 24.283429-117.028571 68.315428-25.6 39.058286-51.2 93.622857-57.197715 158.72l-1.316571 31.012572c0 77.677714 29.257143 143.36 58.514286 188.562285 29.257143 44.763429 58.514286 69.632 117.028571 69.632 43.885714 0 87.771429-24.576 117.028572-68.754285 25.6-39.204571 39.936-93.769143 43.154285-158.427429l0.731429-31.012571c0-76.8-14.628571-142.482286-43.885714-187.977143-29.257143-45.202286-73.142857-70.070857-117.028572-70.070857z m-775.314286 143.945142l-146.285714 182.857143h146.285714v-182.857143z m1623.771429 0l-160.914286 182.857143h160.914286v-182.857143z" fill="#ed4014"></path></svg>{{end}}</use></svg><div class="result-title">{{.Title}}</div><div class="result-desc"><div>{{.Desc}}</div></div><div class="result-extra"><div>{{.Extra}}</div></div>
<div class="result-actions"><div><a href="/" class="btn btn-success" title="返回主页"> 返回主页 </a> <a href="https://github.com/zxysilent" target="_blank" title="查看作者" class="btn btn-default"> 查看作者 </a></div></div></div></div></div></div><div class="footer"><div class="footer-links"><a href="https://blog.zxysilent.com"  target="_blank" title="blog">blog</a><a href="https://github.com/zxysilent/blog" target="_blank" title="github">github</a><a href="https://github.com/zxysilent" target="_blank" title="查看作者">作者</a></div><div class="footer-copyright">Copyright &copy; <script>document.write(new Date().getFullYear());</script>&nbsp;<a target="_blank" href="https://github.com/zxysilent">github.com/zxysilent</a></div></div>{{if ne .ExecJs ""}}<script>{{.ExecJs}}</script>{{end}}</body></html>`
var tmplHtmlBak = `<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8" />
        <meta http-equiv="X-UA-Compatible" content="IE=edge" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
        <title>{{.Title}}</title>
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
                            <use xlink:href="#icon-{{.Kind}}">
							{{if .Circle}}<!-- circle -->
								{{if eq .Kind "succ"}}<svg id="icon-succ" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M725.2 385.8L690 349.6c-1.4-1.6-3.6-2.4-5.6-2.4-2.2 0-4.2 0.8-5.6 2.4l-244 245.8-88.8-88.8c-1.6-1.6-3.6-2.4-5.6-2.4-2 0-4 0.8-5.6 2.4l-35.6 35.6c-3.2 3.2-3.2 8.2 0 11.4l112 112c7.2 7.2 16 11.4 23.4 11.4 10.6 0 19.8-7.8 23.2-11h0.2l267.4-268.8c2.8-3.4 2.8-8.4-0.2-11.4z" fill="#19be6b"></path><path d="M512 152c96.2 0 186.6 37.4 254.6 105.4S872 415.8 872 512s-37.4 186.6-105.4 254.6S608.2 872 512 872c-96.2 0-186.6-37.4-254.6-105.4S152 608.2 152 512s37.4-186.6 105.4-254.6S415.8 152 512 152m0-56C282.2 96 96 282.2 96 512s186.2 416 416 416 416-186.2 416-416S741.8 96 512 96z" fill="#19be6b"></path></svg>
                                {{else if eq .Kind "info"}}<svg id="icon-info" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M470.8 344.4c0-22.8 18.6-39.8 41-39.8 22.8 0 41.4 17 41.4 39.8s-18.6 40-41.4 40c-22.4 0-41-17.2-41-40z m2.8 71.4H550V704h-76.4V415.8z" fill="#2db7f5"></path><path d="M512 152c96.2 0 186.6 37.4 254.6 105.4S872 415.8 872 512s-37.4 186.6-105.4 254.6S608.2 872 512 872c-96.2 0-186.6-37.4-254.6-105.4S152 608.2 152 512s37.4-186.6 105.4-254.6S415.8 152 512 152m0-56C282.2 96 96 282.2 96 512s186.2 416 416 416 416-186.2 416-416S741.8 96 512 96z" fill="#2db7f5"></path></svg>
                                {{else if eq .Kind "warn"}}<svg id="icon-warn" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M512 96c-229.8 0-416 186.2-416 416s186.2 416 416 416 416-186.2 416-416-186.2-416-416-416z m254.6 670.6c-68 68-158.4 105.4-254.6 105.4-96.2 0-186.6-37.4-254.6-105.4S152 608.2 152 512c0-96.2 37.4-186.6 105.4-254.6s158.4-105.4 254.6-105.4c96.2 0 186.6 37.4 254.6 105.4S872 415.8 872 512c0 96.2-37.4 186.6-105.4 254.6z" fill="#ff9900"></path><path d="M512 575c16.4 0 27.8-13.2 28.2-29.6L546.4 315.2c0.4-19.4-15.2-35.2-34.4-35.2-19.4 0-34.8 16-34.4 35.2l6.2 230.2c0.4 16.4 11.8 29.6 28.2 29.6zM512 634.4c-21.4 0-38.2 16.2-38.2 36.8s16.8 36.8 38.2 36.8 38.2-16.2 38.2-36.8-16.8-36.8-38.2-36.8z" fill="#ff9900"></path></svg>
                                {{else if eq .Kind "error"}}<svg id="icon-error"  viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M662.6 617.4L557.2 512l105.4-105.4c12.4-12.4 12.4-32.8 0-45.2-12.4-12.4-32.8-12.4-45.2 0L512 466.8l-105.4-105.4c-12.4-12.4-31.2-14.2-45.2 0-14.2 14.2-12 33.2 0 45.2l105.4 105.4-105.4 105.4c-13.4 13.4-12.8 32.6 0 45.2 12.8 12.8 32.8 12.4 45.2 0l105.4-105.4 105.4 105.4c12.4 12.4 32.8 12.4 45.2 0 12.6-12.4 12.6-32.8 0-45.2z" fill="#ed4014"></path><path d="M512 152c96.2 0 186.6 37.4 254.6 105.4S872 415.8 872 512s-37.4 186.6-105.4 254.6S608.2 872 512 872c-96.2 0-186.6-37.4-254.6-105.4S152 608.2 152 512s37.4-186.6 105.4-254.6S415.8 152 512 152m0-56C282.2 96 96 282.2 96 512s186.2 416 416 416 416-186.2 416-416S741.8 96 512 96z" fill="#ed4014"></path></svg>{{end}}
							{{else}}<!-- not circle -->
								{{if eq .Kind "succ"}}<svg id="icon-succ" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M512 96C282.2 96 96 282.2 96 512s186.2 416 416 416 416-186.2 416-416S741.8 96 512 96z m213 301L457.6 665.6h-0.2c-3.4 3.4-12.6 11-23.2 11-7.6 0-16.2-4.2-23.4-11.4l-112-112c-3.2-3.2-3.2-8.2 0-11.4l35.6-35.6c1.6-1.6 3.6-2.4 5.6-2.4 2 0 4 0.8 5.6 2.4l88.8 88.8 244-245.8c1.6-1.6 3.6-2.4 5.6-2.4 2.2 0 4.2 0.8 5.6 2.4l35 36.2c3.6 3.4 3.6 8.4 0.4 11.6z" fill="#19be6b"></path></svg>
								{{else if eq .Kind "info"}}<svg id="icon-info" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M512 96C282.2 96 96 282.2 96 512s186.2 416 416 416 416-186.2 416-416S741.8 96 512 96z m38 608h-76.4V415.8H550V704z m-38.2-319.6c-22.6 0-41-17.2-41-40s18.6-39.8 41-39.8c22.8 0 41.4 17 41.4 39.8s-18.6 40-41.4 40z" fill="#2db7f5"></path></svg>
								{{else if eq .Kind "warn"}}<svg id="icon-warn" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M512 96C282.2 96 96 282.2 96 512s186.2 416 416 416 416-186.2 416-416S741.8 96 512 96z m34.4 219.2l-6.2 230.2c-0.4 16.4-11.8 29.6-28.2 29.6s-27.8-13.2-28.2-29.6l-6.2-230.2c-0.4-19.2 15-35.2 34.4-35.2 19.2 0 34.8 15.8 34.4 35.2zM512 708c-21.4 0-38.2-16.2-38.2-36.8s16.8-36.8 38.2-36.8c21.4 0 38.2 16.2 38.2 36.8S533.4 708 512 708z" fill="#ff9900"></path></svg>
								{{else if eq .Kind "error"}}<svg id="icon-error" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M512 96C282.2 96 96 282.2 96 512s186.2 416 416 416 416-186.2 416-416S741.8 96 512 96z m105.4 566.6L512 557.2l-105.4 105.4c-12.4 12.4-32.8 12.4-45.2 0-6.2-6.2-9.4-14.4-9.4-22.6 0-8.2 3.2-16.4 9.4-22.6l105.4-105.4-105.4-105.4c-6.2-6.2-9.4-14.4-9.4-22.6 0-8.2 3.2-16.4 9.4-22.6 12.4-12.4 32.8-12.4 45.2 0l105.4 105.4 105.4-105.4c12.4-12.4 32.8-12.4 45.2 0 12.4 12.4 12.4 32.8 0 45.2L557.2 512l105.4 105.4c12.4 12.4 12.4 32.8 0 45.2-12.4 12.6-32.8 12.6-45.2 0z" fill="#ed4014"></path></svg>{{end}}
							{{end}}
                            {{if eq .Kind "404"}}<svg id="icon-404"viewBox="0 0 2486 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M1276.342857 15.798857c239.908571 0 397.897143 200.265143 409.014857 454.656l0.585143 31.012572c0 267.849143-160.914286 485.814857-409.6 485.814857-254.244571 0-412.964571-200.850286-423.643428-455.68l-0.585143-30.134857c0-267.702857 175.542857-485.668571 424.228571-485.668572z m-658.285714 1.316572c58.514286 0 117.028571 51.638857 117.028571 118.198857v434.907428h29.257143c43.885714 0 102.4 46.957714 102.4 105.472 0 59.099429-58.514286 104.301714-102.4 104.301715h-29.257143v81.334857c0 66.56-58.514286 118.198857-117.028571 118.198857s-117.028571-52.077714-117.028572-118.198857v-81.334857h-336.457142c-58.514286 0-117.028571-47.835429-117.028572-115.712 0-41.984 21.504-75.483429 46.08-107.52l12.434286-15.945143L486.4 92.16c14.628571-21.942857 29.257143-41.252571 58.514286-54.710857 14.628571-13.750857 43.885714-20.333714 73.142857-20.333714z m1623.771428 0c58.514286 0 117.028571 51.638857 117.028572 118.198857v434.907428h14.628571c58.514286 0 102.4 46.957714 102.4 105.472 0 59.099429-43.885714 104.301714-102.4 104.301715h-14.628571v81.334857c0 66.56-58.514286 118.198857-117.028572 118.198857-73.142857 0-117.028571-52.077714-117.028571-118.198857v-81.334857h-336.457143c-58.514286 0-117.028571-47.835429-117.028571-115.712 0-40.96 10.24-73.581714 30.427428-104.740572l13.458286-18.724571 380.342857-448.658286c11.702857-8.777143 21.065143-17.115429 28.964572-24.868571l29.549714-29.842286c29.257143-13.750857 58.514286-20.333714 87.771428-20.333714z m-965.485714 226.304c-58.514286 0-87.771429 24.283429-117.028571 68.315428-25.6 39.058286-51.2 93.622857-57.197715 158.72l-1.316571 31.012572c0 77.677714 29.257143 143.36 58.514286 188.562285 29.257143 44.763429 58.514286 69.632 117.028571 69.632 43.885714 0 87.771429-24.576 117.028572-68.754285 25.6-39.204571 39.936-93.769143 43.154285-158.427429l0.731429-31.012571c0-76.8-14.628571-142.482286-43.885714-187.977143-29.257143-45.202286-73.142857-70.070857-117.028572-70.070857z m-775.314286 143.945142l-146.285714 182.857143h146.285714v-182.857143z m1623.771429 0l-160.914286 182.857143h160.914286v-182.857143z" fill="#ed4014"></path></svg>{{end}}
                            </use>
                        </svg>
                        <div class="result-title">{{.Title}}</div>
                        <div class="result-desc">
                            <div>{{.Desc}}</div>
                        </div>
                        <div class="result-extra">
                            <div>{{.Extra}}</div>
                        </div>
                        <div class="result-actions">
                            <div>
                                <a href="/" class="btn btn-success" title="返回主页"> 返回主页 </a>
								<!-- <a href="http://app.helaodash.com/app/#/login" class="btn btn-primary" title="返回主页"> 重新登录 </a> -->
                                <a href="https://github.com/zxysilent" target="_blank" title="查看作者" class="btn btn-default"> 查看作者 </a>
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
        {{if ne .ExecJs ""}}<script>{{.ExecJs}}</script>{{end}}
    </body>
</html>`
