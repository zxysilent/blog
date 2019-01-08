package control

import (
	"blog/model"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/zxysilent/util"
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
		return ctx.JSON(util.NewErrIpt(`请输入用户名和密码`, err.Error()))
	}
	fmt.Println(ipt)
	if ipt.Num == "" && len(ipt.Num) > 18 {
		return ctx.JSON(util.NewErrIpt(`请输入正确的用户名`))
	}
	mod, has := model.UserByNum(ipt.Num)
	if !has {
		return ctx.JSON(util.NewErrOpt(`用户名输入错误`))
	}
	if !mod.Role.IsAtv() {
		return ctx.JSON(util.NewFail(`当前账号已被禁用`))
	}
	if mod.Pass != ipt.Pass {
		return ctx.JSON(util.NewFail(`密码输入错误`))
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
		return ctx.JSON(util.NewFail(`凭证生成失败,请重试`, err.Error()))
	}
	model.UserEditLogin(mod.Id, ctx.RealIP())
	return ctx.JSON(util.NewSucc(`登陆成功`, jwtStr))
}

// UserLogout doc
// @Tags auth
// @Summary 注销
// @Success 200 {object} util.Result "成功数据"
// @Router /logout [post]
func UserLogout(ctx echo.Context) error {
	return ctx.HTML(200, `hello`)
}

// Core 重定向
func Core(ctx echo.Context) error {
	// 301 永久
	// 302 临时
	return ctx.Redirect(301, "/core/")
}

// UserAuth 登陆信息
func UserAuth(ctx echo.Context) error {
	mod, _ := model.UserGet(ctx.Get("uid").(int))
	return ctx.JSON(util.NewSucc(`信息`, mod))
}

// Upload 上传文件
func Upload(ctx echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.JSON(util.NewErrIpt(`未发现文件,请重试`, err.Error()))
	}
	src, err := file.Open()
	if err != nil {
		return ctx.JSON(util.NewErrIpt(`文件打开失败,请重试`, err.Error()))
	}
	defer src.Close()
	basePath := "res/upimg/" + time.Now().Format(util.FmtyyyyMMdd) + "/"
	//确保文件夹存在
	os.MkdirAll(basePath, 0777)
	fileName := util.RandStr(16) + filepath.Ext(file.Filename)
	filePathName := basePath + fileName
	dst, err := os.Create(filePathName)
	if err != nil {
		return ctx.JSON(util.NewErrIpt(`目标文件创建失败,请重试`, err.Error()))
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return ctx.JSON(util.NewErrIpt(`文件写入失败,请重试`, err.Error()))
	}
	return ctx.JSON(util.NewSucc(`文件上传成功`, "/"+filePathName))
}
