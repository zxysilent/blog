package kernel

import (
	"blog/conf"
	"blog/internal/model"
	"blog/internal/repo"
	"blog/pkg/rate"
	"blog/pkg/reply"
	"blog/pkg/token"
	"blog/pkg/vcode"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/logs"
)

// 登录接口全局限流器
var loginLimiter = rate.NewLimiter(20, 5)

// AuthCaptcha doc
// @Tags auth
// @Summary 获取图形验证码
// @Success 200 {object} reply.Reply{data=vcode.Field} "返回数据"
// @Router /api/captcha [get]
func AuthCaptcha(ctx echo.Context) error {
	field := vcode.New(4, authCfg().TokenSecret)
	// 仅调试模式输出明文答案，方便本地联调，生产不泄露
	if conf.IsDebug() {
		logs.Debug("captcha.code:" + field.Code)
	}
	return reply.Succ(ctx, "ok", field)
}

// AuthLogin doc
// @Tags auth
// @Summary 登录-创建令牌
// @Param body body model.UserLoginArgs true "请求数据"
// @Success 200 {object} reply.Reply{data=string} "返回令牌"
// @Router /api/tokens [post]
func AuthLogin(ctx echo.Context) error {
	ipt := &model.UserLoginArgs{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	ipt.Num = strings.TrimSpace(ipt.Num)
	if ipt.Num == "" || len(ipt.Num) > 64 {
		return reply.Fail(ctx, "请输入正确的账号")
	}
	if ipt.Passwd == "" {
		return reply.Fail(ctx, "请输入密码")
	}
	// 限流
	reserve := loginLimiter.Reserve()
	if delay := reserve.Delay(); delay > 0 {
		time.Sleep(min(delay, 5*time.Second))
	}
	// 验证码
	if err := vcode.Check(authCfg().TokenSecret, ipt.Vreal, ipt.Vcode); err != nil {
		return reply.Fail(ctx, "验证码错误", err.Error())
	}
	rctx := ctx.Request().Context()
	mod, err := repo.UserNum(rctx, ipt.Num)
	if err != nil {
		return reply.Fail(ctx, "账号或密码错误")
	}
	if err := mod.EqualPasswd(ipt.Passwd); err != nil {
		return reply.Fail(ctx, "账号或密码错误")
	}
	if mod.Status != model.UserStatusNormal {
		return reply.ErrDeny(ctx, "账号已被禁用")
	}
	// 不再按用户类型做登录闸门：账号能否操作完全由角色授权决定，
	// 未分配角色或角色无权限时，登录后也只能访问无需权限的接口。
	auth := token.Token[token.Auth]{
		ExpAt: tokenExpAt(),
		Claim: token.Auth{
			UserId: mod.Id,
			RoleId: mod.RoleId,
		},
	}
	return reply.Succ(ctx, "登录成功", auth.Encode(authCfg().TokenSecret))
}

// authCfg 认证配置（令牌过期时间与加密私钥）
// 缓存由 repo.SyncCfg 维护，配置修改后立即生效，无需重启
func authCfg() *model.Auth {
	return repo.GetCfg().Auth
}

// tokenExpAt 计算 token 过期时间戳（秒），<=0 表示永不过期
func tokenExpAt() int64 {
	exp := authCfg().TokenExp
	if exp <= 0 {
		return 0
	}
	return time.Now().Add(time.Duration(exp) * time.Hour).Unix()
}

// AuthMe doc
// @Tags auth
// @Summary 获取当前登录用户
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=model.User} "返回数据"
// @Router /api/auth [get]
func AuthMe(ctx echo.Context) error {
	uid, ok := currentUid(ctx)
	if !ok {
		return reply.ErrToken(ctx, "请重新登陆")
	}
	rctx := ctx.Request().Context()
	mod, err := repo.UserGet(rctx, uid)
	if err != nil {
		logs.Ctx(rctx).Err(err).Str("uid", uid).Error("repo.UserGet")
		return reply.Fail(ctx, "用户不存在")
	}
	return reply.Succ(ctx, "ok", mod)
}

// AuthPerms doc
// @Tags auth
// @Summary 获取当前用户权限标识集合
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=[]string} "返回数据"
// @Router /api/auth/permissions [get]
func AuthPerms(ctx echo.Context) error {
	rid, ok := currentRid(ctx)
	if !ok {
		return reply.ErrToken(ctx, "请重新登陆")
	}
	rctx := ctx.Request().Context()
	mods, err := repo.RolePermAll(rctx, rid)
	if err != nil {
		logs.Ctx(rctx).Err(err).Str("rid", rid).Error("repo.RolePermAll")
		return reply.Fail(ctx, "获取权限失败")
	}
	perms := make([]string, 0, len(mods))
	for _, item := range mods {
		perms = append(perms, item.Perm)
	}
	return reply.Succ(ctx, "ok", perms)
}

// AuthEditMe doc
// @Tags auth
// @Summary 修改当前用户资料
// @Security KeyAuth
// @Param body body model.MeEditArgs true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/auth [put]
func AuthEditMe(ctx echo.Context) error {
	uid, ok := currentUid(ctx)
	if !ok {
		return reply.ErrToken(ctx, "请重新登陆")
	}
	ipt := &model.MeEditArgs{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	rctx := ctx.Request().Context()
	mod, err := repo.UserGet(rctx, uid)
	if err != nil {
		return reply.Fail(ctx, "用户不存在")
	}
	mod.Name = strings.TrimSpace(ipt.Name)
	mod.Avatar = strings.TrimSpace(ipt.Avatar)
	mod.Phone = strings.TrimSpace(ipt.Phone)
	mod.Email = strings.TrimSpace(ipt.Email)
	mod.Updated = time.Now().UnixMilli()
	if err := repo.UserEdit(rctx, mod, "name", "avatar", "phone", "email", "updated"); err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.UserEdit")
		return reply.Fail(ctx, "修改失败")
	}
	return reply.Succ(ctx, "修改成功")
}

// AuthVerifyPasswd doc
// @Tags auth
// @Summary 校验当前账号密码
// @Security KeyAuth
// @Param body body object true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/auth/verify [post]
//
// 用于锁屏解锁这类「已登录状态下再次确认身份」的场景：
// 走登录接口会强制图形验证码，体验过重；这里只比对当前账号密码。
func AuthVerifyPasswd(ctx echo.Context) error {
	uid, ok := currentUid(ctx)
	if !ok {
		return reply.ErrToken(ctx, "请重新登陆")
	}
	ipt := &struct {
		Passwd string `json:"passwd" form:"passwd"`
	}{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	if ipt.Passwd == "" {
		return reply.Fail(ctx, "请输入密码")
	}
	mod, err := repo.UserGet(ctx.Request().Context(), uid)
	if err != nil {
		return reply.Fail(ctx, "用户不存在")
	}
	if err := mod.EqualPasswd(ipt.Passwd); err != nil {
		return reply.Fail(ctx, "密码错误")
	}
	return reply.Succ(ctx, "身份校验通过")
}

// AuthEditPasswd doc
// @Tags auth
// @Summary 修改自己的密码
// @Security KeyAuth
// @Param body body model.MePasswdArgs true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/auth/password [put]
func AuthEditPasswd(ctx echo.Context) error {
	uid, ok := currentUid(ctx)
	if !ok {
		return reply.ErrToken(ctx, "请重新登陆")
	}
	ipt := &model.MePasswdArgs{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	if len(ipt.Npasswd) < 6 {
		return reply.Fail(ctx, "新密码长度不能少于 6 位")
	}
	rctx := ctx.Request().Context()
	mod, err := repo.UserGet(rctx, uid)
	if err != nil {
		return reply.Fail(ctx, "用户不存在")
	}
	if err := mod.EqualPasswd(ipt.Opasswd); err != nil {
		return reply.Fail(ctx, "原密码错误")
	}
	if err := mod.SetPasswd(ipt.Npasswd); err != nil {
		logs.Ctx(rctx).Err(err).Error("SetPasswd")
		return reply.Fail(ctx, "密码加密失败")
	}
	mod.Updated = time.Now().UnixMilli()
	if err := repo.UserEdit(rctx, mod, "passwd", "updated"); err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.UserEdit")
		return reply.Fail(ctx, "修改失败")
	}
	return reply.Succ(ctx, "修改成功")
}

// currentUid 从上下文取当前用户id
func currentUid(ctx echo.Context) (string, bool) {
	uid, ok := ctx.Get("uid").(string)
	return uid, ok && uid != ""
}

// currentRid 从上下文取当前角色id
func currentRid(ctx echo.Context) (string, bool) {
	rid, ok := ctx.Get("rid").(string)
	return rid, ok && rid != ""
}
