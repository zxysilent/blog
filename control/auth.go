package control

import (
	"blog/model"
	"blog/util"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

// UserLogin 登陆
func UserLogin(ctx echo.Context) error {
	ipt := struct {
		Num  string `json:"num" form:"num"`
		Pass string `json:"pass" form:"pass"`
		Isu  bool   `json:"isu" form:"isu"`
	}{}
	err := ctx.Bind(&ipt)
	if err != nil {
		return ctx.Res(util.NewErrIpt(`请输入用户名和密码`, err.Error()))
	}
	fmt.Println(ipt)
	if ipt.Num == "" && len(ipt.Num) > 18 {
		return ctx.Res(util.NewErrIpt(`请输入正确的用户名`))
	}
	mod, has := model.UserByNum(ipt.Num)
	if !has {
		return ctx.Res(util.NewErrOpt(`用户名输入错误`))
	}
	if !mod.Role.IsAtv() {
		return ctx.Res(util.NewFail(`当前账号已被禁用`))
	}
	if mod.Pass != ipt.Pass {
		return ctx.Res(util.NewFail(`密码输入错误`))
	}
	claims := model.JwtClaims{
		Id:   mod.Id,
		Name: mod.Name,
		Num:  mod.Num,
		Role: mod.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(),
		},
	}
	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Generate encoded token and send it as response.
	jwtStr, err := token.SignedString([]byte("zxy.sil.ent"))
	if err != nil {
		return ctx.Res(util.NewFail(`凭证生成失败,请重试`, err.Error()))
	}
	model.UserEditLogin(mod.Id, ctx.RealIP())
	return ctx.Res(util.NewSucc(`登陆成功`, jwtStr))
}

// UserLogout doc
// @Tags auth
// @Summary 注销
// @Success 200 {object} util.Result "成功数据"
// @Router /logout [post]
func UserLogout(ctx echo.Context) error {
	return ctx.HTML(200, `hello`)
}

// UserAuth 登陆信息
func UserAuth(ctx echo.Context) error {
	auth := ctx.Get("auth").(*model.JwtClaims)
	mod, _ := model.UserGet(auth.Id)
	return ctx.Res(util.NewSucc(`信息`, mod))
}
