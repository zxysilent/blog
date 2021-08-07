package sysctl

import (
	"blog/conf"
	"blog/internal/fetch"
	"blog/internal/page"
	"blog/internal/token"
	"blog/model"
	"fmt"
	"log"
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
		html := page.Error()
		html.SetTitle("登录错误")
		html.SetDesc("QQ登录失败，请查看错误信息。")
		html.SetExtra(err.Error())
		return ctx.HTML(200, html.Render())
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
		html := page.Error()
		html.SetTitle("登录错误")
		html.SetDesc("QQ登录失败，请查看错误信息。")
		html.SetExtra("当前QQ未绑定账号")
		return ctx.HTML(200, html.Render())
	}
	auth := token.Auth{
		Id:     mod.Id,
		RoleId: -2, //web
		ExpAt:  time.Now().Add(time.Hour * time.Duration(conf.App.TokenExp)).Unix(),
	}
	html := page.Succ()
	html.SetTitle("登录成功")
	html.SetDesc("QQ登录成功，正在跳转")
	html.SetExtra("登录账号：" + mod.Num)
	html.SetExecJs(fmt.Sprintf(`sessionStorage.setItem("zblog-token","%s");setTimeout(function(){location.href="%s"},2000)`, auth.Encode(conf.App.TokenSecret), "/login.html"))
	return ctx.HTML(200, html.Render())
}
