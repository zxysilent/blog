package kernel

import (
	"blog/conf"
	"blog/internal/model"
	"blog/pkg/rate"
	"blog/pkg/token"
	"blog/pkg/vcode"

	"blog/internal/repo"
	"blog/internal/utils"

	"context"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/logs"
)

// 防止暴力破解,每秒20次登录限制
var loginLimiter = rate.NewLimiter(20, 5)

// AuthLogin doc
// @Tags kernel,auth
// @Summary 登陆
// @Accept json
// @Param body body model.AdminLoginArgs true "请求数据"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/auth/login [post]
func AuthLogin(ctx echo.Context) error {
	ct, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := loginLimiter.Wait(ct); err != nil {
		return ctx.JSON(utils.Fail("当前登录人数过多,请等待", err.Error()))
	}
	ipt := &model.AdminLoginArgs{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("请输入用户名和密码", err.Error()))
	}
	logs.With().Any("args", ipt).Debug("login")
	if err = vcode.Check(conf.App.TokenSecret, ipt.Vreal, ipt.Vcode); err != nil {
		return ctx.JSON(utils.ErrIpt(err.Error()))
	}
	if ipt.Num == "" && len(ipt.Num) > 18 {
		return ctx.JSON(utils.ErrIpt("请输入正确的用户名"))
	}
	mod, has := repo.AdminNum(ipt.Num)
	if !has {
		return ctx.JSON(utils.ErrOpt("用户名输入错误"))
	}
	if mod.Passwd != ipt.Enc() {
		return ctx.JSON(utils.Fail("密码输入错误"))
	}
	auth := token.Token[token.Auth]{
		ExpAt: int64(conf.App.TokenExp),
		Claim: token.Auth{
			UserId: mod.Id,
			RoleId: mod.RoleId,
			Admin:  true,
		},
	}
	if conf.App.TokenExp > 0 {
		auth.ExpAt = time.Now().Add(time.Hour * time.Duration(conf.App.TokenExp)).Unix()
	}
	token := auth.Encode(conf.App.TokenSecret)
	// logs.Info("登陆成功", auth.Id, token)
	// cookie := new(http.Cookie)
	// cookie.Name = conf.App.TokenKey
	// cookie.HttpOnly = true
	// cookie.Path = "/"
	// // cookie.MaxAge = conf.App.Hmtexp * 60 * 60 //秒
	// cookie.Domain = "127.0.0.1"
	// cookie.Value = token
	// cookie.Expires = time.Now().Add(time.Duration(conf.App.TokenExp) * time.Hour)
	// ctx.SetCookie(cookie)
	ctx.Set("uid", mod.Id)
	repo.LOG(ctx, "auth", "登录", auth)
	return ctx.JSON(utils.Succ("登陆成功", token))
}

// AuthVcode doc
// @Tags kernel,auth
// @Summary 验证码
// @Accept mpfd
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/auth/vcode [get]
func AuthVcode(ctx echo.Context) error {
	return ctx.JSON(utils.Succ("succ", vcode.New(4, conf.App.TokenSecret)))
}

// AuthGet doc
// @Tags kernel
// @Auth mgmt
// @Summary 当前登录信息
// @Security ApiKeyAuth
// @Success 200 {object} utils.Reply{data=model.Admin} "返回数据"
// @Router /api/auth/get [get]
func AuthGet(ctx echo.Context) error {
	mod, _ := repo.AdminGet(ctx.Get("uid").(int))
	return ctx.JSON(utils.Succ("succ", mod))
}

// AuthGrant doc
// @Tags kernel,auth
// @Auth
// @Summary 获取当前用户的授权
// @Security ApiKeyAuth
// @Success 200 {object} utils.Reply{data=[]string} "返回数据"
// @Router /api/auth/grant [get]
func AuthGrant(ctx echo.Context) error {
	// 登录信息获取roleId
	roleId, _ := ctx.Get("rid").(int)
	mods, err := repo.RoleGrantAll(roleId)
	if err != nil {
		return ctx.JSON(utils.Fail("未查询到角色授权信息", err.Error()))
	}
	grants := make([]string, 0, len(mods))
	for _, val := range mods {
		grants = append(grants, val.Guid)
	}
	return ctx.JSON(utils.Succ("succ", grants))
}

// AuthEditPasswd doc
// @Tags kernel,auth
// @Auth mgmt
// @Summary 修改自己的密码
// @Security ApiKeyAuth
// @Param body body model.PasswdArgs true "请求数据"
// @Success 200 {object} utils.Reply "返回数据"
// @Router /api/auth/edit/passwd [post]
func AuthEditPasswd(ctx echo.Context) error {
	ipt := &model.PasswdArgs{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("输入数据有误", err.Error()))
	}
	mod, has := repo.AdminGet(ctx.Get("uid").(int))
	if !has {
		return ctx.JSON(utils.Fail("输入数据有误,请重试"))
	}
	ipt.Enc()
	if mod.Passwd != ipt.Opasswd {
		return ctx.JSON(utils.Fail("原始密码输入错误,请重试"))
	}
	mod.Passwd = ipt.Npasswd
	if err := repo.AdminEdit(mod, "passwd"); err != nil {
		return ctx.JSON(utils.Fail("密码修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// AuthEdit doc
// @Tags kernel,auth
// @Auth mgmt
// @Summary 修改自己的信息
// @Security ApiKeyAuth
// @Success 200 {object} utils.Reply "返回数据"
// @Param body body model.Admin true "request"
// @Router /api/auth/edit [post]
func AuthEdit(ctx echo.Context) error {
	ipt := &model.Admin{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("输入数据有误", err.Error()))
	}
	ipt.Id = ctx.Get("uid").(int)
	if err := repo.AdminEdit(ipt, "name", "email", "phone", "avatar"); err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
