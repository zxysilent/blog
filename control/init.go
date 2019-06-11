package control

import (
	"blog/model"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/zxysilent/utils"
)

// IndexView 主页面
func IndexView(ctx echo.Context) error {
	//return ctx.HTML(200, `<html><head><meta charset="UTF-8"><title>文档</title></head><body><a href="/swagger/index.html">doc</a></body></html>`)
	pi, _ := strconv.Atoi(ctx.FormValue("page"))
	if pi == 0 {
		pi = 1
	}
	ps, _ := utils.Atoi(model.MapOpts.MustGet("page_size"), 6)
	mods, _ := model.PostPage(pi, ps)
	total := model.PostCount()
	naver := model.Naver{}
	if pi > 1 {
		naver.Prev = "/?page=" + strconv.Itoa(pi-1)
	}
	if total > (pi * ps) {
		naver.Next = "/?page=" + strconv.Itoa(pi+1)
	}
	return ctx.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Posts": mods,
		"Naver": naver,
	})
}

// ArchivesView 归档
func ArchivesView(ctx echo.Context) error {
	mods, err := model.PostArchive()
	if err != nil {
		return ctx.Redirect(302, "/")
	}
	return ctx.Render(http.StatusOK, "archive.html", map[string]interface{}{
		"Posts": mods,
	})
}

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
