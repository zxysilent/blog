package control

import (
	"blog/model"
	"strconv"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/zxysilent/utils"
)

// UserLogin doc
// @Tags auth
// @Summary 登陆
// @Accept mpfd
// @Param num formData string true "账号" default(super)
// @Param pass formData string true "密码" default(123654)
// @Success 200 {object} utils.Reply "成功数据"
// @Router /login [post]
func UserLogin(ctx echo.Context) error {
	ipt := struct {
		Num  string `json:"num" form:"num"`
		Pass string `json:"pass" form:"pass"`
	}{}
	err := ctx.Bind(&ipt)
	if err != nil {
		return ctx.JSON(utils.NewErrIpt(`请输入用户名和密码`, err.Error()))
	}
	if ipt.Num == "" && len(ipt.Num) > 18 {
		return ctx.JSON(utils.NewErrIpt(`请输入正确的用户名`))
	}
	mod, has := model.UserByNum(ipt.Num)
	if !has {
		return ctx.JSON(utils.NewErrOpt(`用户名输入错误`))
	}
	now := time.Now()
	// 禁止登陆证 5 分钟
	if mod.Ecount == -1 {
		// 登录时间差
		span := 5 - int(now.Sub(mod.Ltime).Minutes())
		if span >= 1 { //「」
			return ctx.JSON(utils.NewFail(`请「` + strconv.Itoa(span) + `」分钟后登录`))
		}
		mod.Ecount = 0
	}
	if mod.Pass != ipt.Pass {
		mod.Ltime = now
		mod.Ecount++
		// 错误次数大于 3 锁定
		if mod.Ecount >= 3 {
			mod.Ecount = -1
			model.UserEditLogin(mod, "Ltime", "Ecount")
			return ctx.JSON(utils.NewFail(`登录锁定请「5」分钟后登录`))
		}
		// 小于3 提示剩余次数
		model.UserEditLogin(mod, "Ltime", "Ecount")
		return ctx.JSON(utils.NewFail(`密码错误,剩于登录次数：` + strconv.Itoa(int(3-mod.Ecount))))
	}
	if !mod.Role.IsAtv() {
		return ctx.JSON(utils.NewFail(`当前账号已被禁用`))
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
		return ctx.JSON(utils.NewFail(`凭证生成失败,请重试`, err.Error()))
	}
	mod.Ltime = now
	mod.Ip = ctx.RealIP()
	model.UserEditLogin(mod, "Ltime", "Ip", "Ecount")
	return ctx.JSON(utils.NewSucc(`登陆成功`, jwtStr))
}

// UserLogout doc
// @Tags auth
// @Summary 注销
// @Success 200 {object} utils.Reply "成功数据"
// @Router /logout [post]
func UserLogout(ctx echo.Context) error {
	return ctx.HTML(200, `hello`)
}

// UserAuth doc
// @Tags auth
// @Summary 获取登录信息
// @Param token query string true "凭证jwt" default(jwt)
// @Success 200 {object} utils.Reply "成功数据"
// @Router /api/auth [get]
func UserAuth(ctx echo.Context) error {
	mod, _ := model.UserGet(ctx.Get("uid").(int))
	return ctx.JSON(utils.NewSucc(`信息`, mod))
}
