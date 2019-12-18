package control

import (
	"blog/model"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/labstack/echo"
	"github.com/zxysilent/utils"
)

// Upload 上传文件
func Upload(ctx echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.JSON(utils.NewErrIpt(`未发现文件,请重试`, err.Error()))
	}
	src, err := file.Open()
	if err != nil {
		return ctx.JSON(utils.NewErrIpt(`文件打开失败,请重试`, err.Error()))
	}
	defer src.Close()
	basePath := "res/upimg/" + time.Now().Format(utils.FmtyyyyMMdd) + "/"
	//确保文件夹存在
	os.MkdirAll(basePath, 0777)
	fileName := utils.RandStr(16) + filepath.Ext(file.Filename)
	filePathName := basePath + fileName
	dst, err := os.Create(filePathName)
	if err != nil {
		return ctx.JSON(utils.NewErrIpt(`目标文件创建失败,请重试`, err.Error()))
	}
	defer dst.Close()
	if _, err = io.Copy(dst, src); err != nil {
		return ctx.JSON(utils.NewErrIpt(`文件写入失败,请重试`, err.Error()))
	}
	return ctx.JSON(utils.NewSucc(`文件上传成功`, "/"+filePathName))
}

// Core 重定向
func Core(ctx echo.Context) error {
	// 301 永久
	// 302 临时
	return ctx.Redirect(301, "/core/")
}

// Sys 系统信息
func Sys(ctx echo.Context) error {
	info := struct {
		ARCH    string `json:"arch" form:"arch"`
		OS      string `json:"os" form:"os"`
		Version string `json:"version" form:"version"`
		NumCPU  int    `json:"num_cpu" form:"num_cpu"`
	}{
		ARCH:    runtime.GOARCH,
		OS:      runtime.GOOS,
		Version: runtime.Version(),
		NumCPU:  runtime.NumCPU(),
	}
	return ctx.JSON(utils.NewSucc(`系统信息`, info))
}

// Collect 统计信息
func Collect(ctx echo.Context) error {
	if mod, has := model.Collect(); has {
		return ctx.JSON(utils.NewSucc(`统计信息`, mod))
	}
	return ctx.JSON(utils.NewFail(`未查询到统计信息`))
}

// Models doc
// @Tags doc
// @Summary 所有定义的数据结构体
// @Param body body model.Page true "Page struct"
// @Param body body model.Cate true "Cate struct"
// @Param body body model.Opts true "Opts struct"
// @Param body body model.PostTag true "PostTag struct"
// @Param body body model.Post true "Post struct"
// @Param body body model.Tag true "Tag struct"
// @Param body body model.User true "User struct"
func Models() {
}
