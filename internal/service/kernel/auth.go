package kernel

import (
	"blog/conf"
	"blog/internal/model"
	"blog/internal/repo"
	"blog/internal/utils"
	"blog/pkg/rate"
	"blog/pkg/reply"
	"blog/pkg/token"
	"blog/pkg/vcode"
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
// @Success 200 {object} reply.Reply{data=string} "返回数据"
// @Router /api/auth/login [post]
func AuthLogin(ctx echo.Context) error {
	rctx := utils.WrapCtx(ctx)
	ct, cancel := context.WithTimeout(rctx, 5*time.Second)
	defer cancel()
	if err := loginLimiter.Wait(ct); err != nil {
		return reply.Fail(ctx, "当前登录人数过多,请等待", err)
	}
	ipt := &model.AdminLoginArgs{}
	err := ctx.Bind(ipt)
	if err != nil {
		return reply.Fail(ctx, "请输入用户名和密码", err)
	}
	logs.With().Any("args", ipt).Debug("login")
	if err = vcode.Check(conf.App.TokenSecret, ipt.Vreal, ipt.Vcode); err != nil {
		return reply.Fail(ctx, "验证码错误")
	}
	if ipt.Num == "" && len(ipt.Num) > 18 {
		return reply.Fail(ctx, "请输入正确的用户名")
	}
	mod, has := repo.AdminNum(ipt.Num)
	if !has {
		return reply.Fail(ctx, "用户名输入错误")
	}
	if mod.Passwd != ipt.Enc() {
		return reply.Fail(ctx, "密码输入错误")
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
	return reply.Succ(ctx, "登陆成功", token)
}

// AuthVcode doc
// @Tags kernel,auth
// @Summary 验证码
// @Accept mpfd
// @Success 200 {object} reply.Reply{data=string} "返回数据"
// @Router /api/auth/vcode [get]
func AuthVcode(ctx echo.Context) error {
	return reply.Succ(ctx, "succ", vcode.New(4, conf.App.TokenSecret))
}

// AuthGet doc
// @Tags kernel
// @Auth mgmt
// @Summary 当前登录信息
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=model.Admin} "返回数据"
// @Router /api/auth/get [get]
func AuthGet(ctx echo.Context) error {
	rctx := utils.WrapCtx(ctx)
	uid, _ := ctx.Get("uid").(string)
	mod, err := repo.GetAny[model.Admin](rctx, uid)
	if err != nil {
		return reply.Fail(ctx, "未查询到信息", err)
	}
	return reply.Succ(ctx, "succ", mod)
}

// AuthGrant doc
// @Tags kernel,auth
// @Auth
// @Summary 获取当前用户的授权
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=[]string} "返回数据"
// @Router /api/auth/grant [get]
func AuthGrant(ctx echo.Context) error {
	rctx := utils.WrapCtx(ctx)
	roleId, _ := ctx.Get("rid").(string)
	mods, err := repo.RoleGrantAll(rctx, roleId)
	if err != nil {
		return reply.Fail(ctx, "未查询到角色授权信息", err)
	}
	grants := make([]string, 0, len(mods))
	for _, val := range mods {
		grants = append(grants, val.Guid)
	}
	return reply.Succ(ctx, "succ", grants)
}

// AuthEditPasswd doc
// @Tags kernel,auth
// @Auth mgmt
// @Summary 修改自己的密码
// @Security KeyAuth
// @Param body body model.PasswdArgs true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/auth/edit/passwd [post]
func AuthEditPasswd(ctx echo.Context) error {
	rctx := utils.WrapCtx(ctx)
	ipt := &model.PasswdArgs{}
	err := ctx.Bind(ipt)
	if err != nil {
		return reply.Fail(ctx, "输入数据有误", err)
	}
	uid, _ := ctx.Get("uid").(string)
	mod, err := repo.GetAny[model.Admin](rctx, uid)
	if err != nil {
		return reply.Fail(ctx, "未查询到信息", err)
	}
	ipt.Enc()
	if mod.Passwd != ipt.Opasswd {
		return reply.Fail(ctx, "原始密码输入错误,请重试")
	}
	mod.Passwd = ipt.Npasswd
	if err := repo.AdminEdit(rctx, &mod, "passwd"); err != nil {
		return reply.Fail(ctx, "密码修改失败", err.Error())
	}
	return reply.Succ(ctx, "succ")
}

// AuthEdit doc
// @Tags kernel,auth
// @Auth mgmt
// @Summary 修改自己的信息
// @Security KeyAuth
// @Success 200 {object} reply.Reply "返回数据"
// @Param body body model.Admin true "request"
// @Router /api/auth/edit [post]
func AuthEdit(ctx echo.Context) error {
	rctx := utils.WrapCtx(ctx)
	ipt := &model.Admin{}
	err := ctx.Bind(ipt)
	if err != nil {
		return reply.Fail(ctx, "输入数据有误", err.Error())
	}
	ipt.Id, _ = ctx.Get("uid").(string)
	if err := repo.AdminEdit(rctx, ipt, "name", "email", "phone", "avatar"); err != nil {
		return reply.Fail(ctx, "修改失败", err)
	}
	return reply.Succ(ctx, "succ")
}
