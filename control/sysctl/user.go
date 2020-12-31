package sysctl

import (
	"blog/conf"
	"blog/internal/jwt"
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

// 防止暴力破解,每秒20次登录限制
var loginLimiter = rate.NewLimiter(20, 5)

// UserLogin doc
// @Tags auth
// @Summary 登陆
// @Accept mpfd
// @Param num formData string true "账号" default(super)
// @Param pass formData string true "密码" default(123654)
// @Router /login [post]
func Login(ctx echo.Context) error {
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
	mod, has := model.UserByNum(ipt.Num)
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
	if mod.Pass != ipt.Passwd {
		mod.Ltime = now
		mod.Ecount++
		// 错误次数大于 3 锁定
		if mod.Ecount >= 3 {
			mod.Ecount = -1
			model.UserEditLogin(mod, "Ltime", "Ecount")
			return ctx.JSON(utils.Fail(`登录锁定请「5」分钟后登录`))
		}
		// 小于3 提示剩余次数
		model.UserEditLogin(mod, "Ltime", "Ecount")
		return ctx.JSON(utils.Fail(`密码错误,剩于登录次数：` + strconv.Itoa(int(3-mod.Ecount))))
	}
	if !mod.Role.IsAtv() {
		return ctx.JSON(utils.Fail(`当前账号已被禁用`))
	}
	auth := jwt.JwtAuth{
		Id:    mod.Id,
		Role:  int(mod.Role),
		ExpAt: time.Now().Add(time.Hour * 2).Unix(),
	}
	mod.Ltime = now
	mod.Ip = ctx.RealIP()
	model.UserEditLogin(mod, "Ltime", "Ip", "Ecount")
	return ctx.JSON(utils.Succ(`登陆成功`, auth.Encode(conf.App.TokenSecret)))
}

// UserLogout doc
// @Tags auth
// @Summary 注销
// @Router /logout [post]
func UserLogout(ctx echo.Context) error {
	return ctx.HTML(200, `hello`)
}

// UserAuth doc
// @Tags auth
// @Summary 获取登录信息
// @Param token query string true "凭证jwt" default(jwt)
// @Router /api/auth [get]
func UserAuth(ctx echo.Context) error {
	mod, _ := model.UserGet(ctx.Get("uid").(int))
	return ctx.JSON(utils.Succ(`信息`, mod))
}

// Login doc
// @Tags auth
// @Summary 登陆
// @Accept mpfd
// @Param num formData string true "账号" default(super)
// @Param passwd formData string true "密码" default(123654)
// @Router /api/login [post]
func Vcode(ctx echo.Context) error {
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

//---------- 管理相关

// SysUserGet doc
// @Tags sysuser
// @Summary 通过id获取单条用户信息
// @Param id path int true "pk id" default(1)
// @Router /sys/user/get/{id} [get]
func SysUserGet(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	mod, has := model.SysUserGet(id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到用户信息"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// SysUserAll doc
// @Tags sysuser
// @Summary 获取所有用户信息
// @Router /sys/user/all [get]
func SysUserAll(ctx echo.Context) error {
	mods, err := model.SysUserAll()
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到用户信息", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// SysUserPage doc
// @Tags sysuser
// @Summary 获取用户分页信息
// @Param cid path int true "分类id" default(1)
// @Param pi query int true "分页数" default(1)
// @Param ps query int true "每页条数[5,20]" default(5)
// @Router /sys/user/page/{cid} [get]
func SysUserPage(ctx echo.Context) error {
	// cid, err := strconv.Atoi(ctx.Param("cid"))
	// if err != nil {
	//  return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	// }
	ipt := &model.Page{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if ipt.Ps > 20 || ipt.Ps < 5 {
		return ctx.JSON(utils.ErrIpt("分页大小输入错误", ipt.Ps))
	}
	count := model.SysUserCount()
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := model.SysUserPage(ipt.Pi, ipt.Ps)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, int(count)))
}

// SysUserAdd doc
// @Tags sysuser
// @Summary 添加用户信息
// @Param token query string true "hmt" default(token)
// @Router /sys/user/add [post]
func SysUserAdd(ctx echo.Context) error {
	ipt := &model.SysUser{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// ipt.Utime = time.Now()
	err = model.SysUserAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// SysUserEdit doc
// @Tags sysuser
// @Summary 修改用户信息
// @Param token query string true "hmt" default(token)
// @Router /sys/user/edit [post]
func SysUserEdit(ctx echo.Context) error {
	ipt := &model.SysUser{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// ipt.Utime = time.Now()
	err = model.SysUserEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// SysUserDrop doc
// @Tags sysuser
// @Summary 通过id删除单条用户信息
// @Param id path int true "pk id" default(1)
// @Param token query string true "hmt" default(token)
// @Router /sys/user/drop/{id} [get]
func SysUserDrop(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	err = model.SysUserDrop(id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
