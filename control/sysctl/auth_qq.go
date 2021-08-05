package sysctl

import (
	"blog/conf"
	"blog/internal/fetch"
	"blog/internal/token"
	"blog/model"
	"fmt"
	"log"
	"net/http"
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
	cookie := new(http.Cookie)
	cookie.Name = "xgapp"
	cookie.HttpOnly = true
	cookie.Path = "/"
	// cookie.MaxAge = conf.App.TokenExp * 60 * 60 //秒
	cookie.Domain = ".xgapp.tuzuu.com"
	cookie.Value = auth.Encode(conf.App.TokenSecret)
	cookie.Expires = time.Now().Add(time.Duration(conf.App.TokenExp) * time.Hour)
	ctx.SetCookie(cookie)
	// return ctx.HTML(http.StatusOK, "登陆成功:")
	return ctx.Redirect(302, "/web")
}
