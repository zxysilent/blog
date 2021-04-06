package sysctl

import (
	"blog/conf"
	"blog/internal/hwt"
	"blog/internal/rate"
	"blog/internal/vcode"
	"blog/model"
	"context"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// auth 相关-登录，注销、个人信息、角色菜单、角色接口

//RBAC 模型
// sysctl 独立权限管理
// role 管理auth(api访问)
// role 管理menu(导航菜单)

// 防止暴力破解,每秒20次登录限制
var loginLimiter = rate.NewLimiter(20, 5)

// UserLogin doc
// @Tags auth
// @Summary 登陆
// @Accept mpfd
// @Param num formData string true "账号" default(zxysilent)
// @Param pass formData string true "密码" default(zxyslt)
// @Success 200 {object} model.Reply{data=string} "成功数据"
// @Router /api/auth/login [post]
func AuthLogin(ctx echo.Context) error {
	ct, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := loginLimiter.Wait(ct); err != nil {
		return ctx.JSON(utils.Fail("当前登录人数过多,请等待", err.Error()))
	}
	ipt := struct {
		Num    string `json:"num" form:"num"`
		Vcode  string `form:"vcode" json:"vcode"`
		Vreal  string `form:"vreal" json:"vreal"`
		Passwd string `json:"passwd" form:"passwd"`
	}{}
	err := ctx.Bind(&ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("请输入用户名和密码", err.Error()))
	}
	if ipt.Vreal != hmc(ipt.Vcode, "v.c.o.d.e") {
		return ctx.JSON(utils.ErrIpt("请输入正确的验证码"))
	}
	if ipt.Num == "" && len(ipt.Num) > 18 {
		return ctx.JSON(utils.ErrIpt(`请输入正确的账号`))
	}
	mod, has := model.UserLogin(ipt.Num)
	if !has {
		return ctx.JSON(utils.ErrOpt(`账号输入错误`))
	}
	now := time.Now()
	// 禁止登陆证 5 分钟
	if mod.Ecount == -1 {
		// 登录时间差
		span := 5 - int(now.Sub(mod.Ltime).Minutes())
		if span >= 1 { //「」
			return ctx.JSON(utils.Fail(`请「` + strconv.Itoa(span) + `」分钟后登录`))
		}
		mod.Ecount = 0
	}
	if mod.Passwd != ipt.Passwd {
		mod.Ltime = now
		mod.Ecount++
		// 错误次数大于 3 锁定
		if mod.Ecount >= 3 {
			mod.Ecount = -1
			//model.UserEditLogin(mod, "Ltime", "Ecount")
			return ctx.JSON(utils.Fail(`登录锁定请「5」分钟后登录`))
		}
		// 小于3 提示剩余次数
		//model.UserEditLogin(mod, "Ltime", "Ecount")
		return ctx.JSON(utils.Fail(`密码错误,剩于登录次数：` + strconv.Itoa(int(3-mod.Ecount))))
	}
	// if !mod.Role1.IsAtv() {
	// 	return ctx.JSON(utils.Fail(`当前账号已被禁用`))
	// }
	auth := hwt.Auth{
		Id:     mod.Id,
		RoleId: mod.RoleId,
		ExpAt:  time.Now().Add(time.Hour * time.Duration(conf.App.TokenExp)).Unix(),
	}
	mod.Ltime = now
	// mod.Ip = ctx.RealIP()
	// model.UserEditLogin(mod, "Ltime", "Ip", "Ecount")
	return ctx.JSON(utils.Succ(`登陆成功`, auth.Encode(conf.App.TokenSecret)))
}

// UserLogout doc
// @Tags auth
// @Summary 注销
// @Router /api/auth/logout [post]
func UserLogout(ctx echo.Context) error {
	return ctx.HTML(200, `hello`)
}

// UserAuth doc
// @Tags auth
// @Summary 获取登录信息
// @Param token query string true "凭证"
// @Success 200 {object} model.Reply{data=model.User} "成功数据"
// @Router /adm/auth [get]
func Auth(ctx echo.Context) error {
	mod, _ := model.UserGet(ctx.Get("uid").(int))
	return ctx.JSON(utils.Succ(`auth`, mod))
}

// Login doc
// @Tags auth
// @Summary 登陆
// @Accept mpfd
// @Success 200 {object} model.Reply "成功数据"
// @Router /api/auth/vcode [post]
func AuthVcode(ctx echo.Context) error {
	rnd := utils.RandDigitStr(5)
	out := struct {
		Vcode string `json:"vcode"`
		Vreal string `json:"vreal"`
	}{
		Vcode: vcode.NewImage(rnd).Base64(),
		Vreal: hmc(rnd, "v.c.o.d.e"),
	}
	return ctx.JSON(utils.Succ("succ", out))
}

func hmc(raw, key string) string {
	hm := hmac.New(sha1.New, []byte(key))
	hm.Write([]byte(raw))
	return base64.RawURLEncoding.EncodeToString(hm.Sum(nil))
}
