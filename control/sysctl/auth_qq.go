package sysctl

import (
	"blog/conf"
	"blog/internal/fetch"
	"blog/internal/token"
	"blog/model"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

const redirectUrl = "https%3A%2F%2Fblog.zxysilent.com%2Fauth%2Fqq.html"

func ViewLoginQq(ctx echo.Context) error {
	return ctx.Redirect(302, fmt.Sprintf("https://graph.qq.com/oauth2.0/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=zxysilent&", conf.App.Qq.WebAppid, redirectUrl)) //临时
}

func ViewAuthQq(ctx echo.Context) error {
	loginUrl := fmt.Sprintf("https://graph.qq.com/oauth2.0/token?grant_type=authorization_code&client_id=%s&client_secret=%s&code=%s&redirect_uri=%s&fmt=json", conf.App.Qq.WebAppid, conf.App.Qq.WebSecret, ctx.FormValue("code"), redirectUrl)
	atoken := struct {
		AccessToken  string `json:"access_token"`
		ExpiresIn    string `json:"expires_in"`
		RefreshToken string `json:"refresh_token"`
	}{}
	err := fetch.Get(loginUrl, &atoken)
	if err != nil {
		//w.Write([]byte(err.Error()))
	}
	log.Println(atoken)
	reply := struct {
		ClientId string `json:"client_id"`
		Openid   string `json:"openid"`
	}{}
	// 获取openid
	fetch.Get(fmt.Sprintf("https://graph.qq.com/oauth2.0/me?access_token=%s&fmt=json", atoken.AccessToken), &reply)
	log.Println(reply)
	// return ctx.JSON(utils.Succ("succ", reply))
	// 判断 openid 是否存在
	mod, has := model.UserOpenidQq(reply.Openid)
	if !has {
		//登陆失败
	}
	auth := token.Auth{
		Id:     mod.Id,
		RoleId: -2, //web
		ExpAt:  time.Now().Add(time.Hour * time.Duration(conf.App.TokenExp)).Unix(),
	}
	// return ctx.HTML(http.StatusOK, "登陆成功:")
	return ctx.HTML(200, strings.Replace(htmlLoginSucc, "{{token}}", auth.Encode(conf.App.TokenSecret), 1))
}

const htmlLoginSucc = `<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8" />
        <meta http-equiv="X-UA-Compatible" content="IE=edge" />
        <meta name="viewport" content="width=device-width, initial-scale=1.0" />
		<meta http-equiv="refresh" content="2; url=/login.html" />  
        <title>登录成功正在跳转</title>
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
                            <use xlink:href="#icon-success">
                                <svg id="icon-success" viewBox="0 0 1024 1024" version="1.1" xmlns="http://www.w3.org/2000/svg"><path d="M512 96C282.2 96 96 282.2 96 512s186.2 416 416 416 416-186.2 416-416S741.8 96 512 96z m213 301L457.6 665.6h-0.2c-3.4 3.4-12.6 11-23.2 11-7.6 0-16.2-4.2-23.4-11.4l-112-112c-3.2-3.2-3.2-8.2 0-11.4l35.6-35.6c1.6-1.6 3.6-2.4 5.6-2.4 2 0 4 0.8 5.6 2.4l88.8 88.8 244-245.8c1.6-1.6 3.6-2.4 5.6-2.4 2.2 0 4.2 0.8 5.6 2.4l35 36.2c3.6 3.4 3.6 8.4 0.4 11.6z" fill="#19be6b"></path></svg>
                            </use>
                        </svg>
                        <div class="result-title">登录成功</div>
                        <div class="result-desc">
                            <div>QQ登录成功，正在跳转</div>
                        </div>
                        <div class="result-extra">
                            <div>暂无额外信息</div>
                        </div>
                        <div class="result-actions">
                            <div>
                                <a href="/" class="btn btn-success" title="返回主页"> 返回主页 </a>
                                <a href="javascript:history.go(-1)" class="btn btn-primary" title="返回主页"> 上一页 </a>
                                <!-- <a href="https://github.com/zxysilent" target="_blank" title="查看作者" class="btn btn-default"> 查看作者 </a> -->
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
		<script>sessionStorage.setItem("zblog-token","{{token}}")</script>
    </body>
</html>`
