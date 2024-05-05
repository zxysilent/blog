package kernel

import (
	"blog/conf"
	"blog/internal/model"
	"blog/internal/repo"
	"blog/internal/utils"

	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/labstack/echo/v4"
	_ "golang.org/x/image/bmp"
)

// UploadFile doc
// @Auth
// @Tags kernel,upload
// @Summary 上传文件
// @Accept  mpfd
// @Param file formData file true "file"
// @Security ApiKeyAuth
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/upload/file [post]
func UploadFile(ctx echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.JSON(utils.Fail("未发现文件", err.Error()))
	}
	src, err := file.Open()
	if err != nil {
		return ctx.JSON(utils.Fail("文件打开失败", err.Error()))

	}
	defer src.Close()
	dir := time.Now().Format("200601/02")
	os.MkdirAll("./static/upload/"+dir[:6], 0755)
	name := "static/upload/" + dir + utils.SUID() + path.Ext(file.Filename)
	dst, err := os.Create(name)
	if err != nil {
		return ctx.JSON(utils.Fail("文件打创建文件失败", err.Error()))

	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		return ctx.JSON(utils.Fail("文件保存失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("上传成功", conf.App.Endpoint+"/"+name))
}

// UploadImage doc
// @Auth
// @Tags kernel,upload
// @Summary 上传图片并裁剪
// @Accept  mpfd
// @Param file formData file true "file"
// @Security ApiKeyAuth
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/upload/image [post]
func UploadImage(ctx echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.JSON(utils.ErrIpt("未发现文件,请重试", err.Error()))
	}
	if !strings.Contains(file.Header.Get("Content-Type"), "image") {
		return ctx.JSON(utils.ErrIpt("请选择图片文件", file.Header.Get("Content-Type")))
	}
	src, err := file.Open()
	if err != nil {
		return ctx.JSON(utils.ErrIpt("文件打开失败,请重试", err.Error()))
	}
	defer src.Close()
	attr := &struct {
		Cut bool `json:"cut" query:"cut" form:"cut"`
		Wd  int  `json:"wd" query:"wd" form:"wd"`
		Hd  int  `json:"hd" query:"hd" form:"hd"`
	}{}
	ctx.Bind(attr)
	dir := time.Now().Format("200601/02")
	os.MkdirAll("./static/upload/"+dir[:6], 0755)
	ext := path.Ext(file.Filename)
	if conf.App.ImageCut || attr.Cut {
		ext = ".jpg"
	}
	name := "static/upload/" + dir + utils.SUID() + ext
	dst, err := os.Create(name)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("目标文件创建失败,请重试", err.Error()))
	}
	defer dst.Close()
	if conf.App.ImageCut || attr.Cut { //图片裁剪
		imgSrc, _, err := image.Decode(src)
		// 图片解码
		if err != nil {
			return ctx.JSON(utils.ErrIpt("读取图片失败,请重试", err.Error()))
		}
		if attr.Wd <= 0 {
			attr.Wd = conf.App.ImageWidth
		}
		if attr.Hd <= 0 {
			attr.Hd = conf.App.ImageHeight
		}
		// 宽度>指定宽度 防止负调整
		dx := imgSrc.Bounds().Dx()
		if dx > attr.Wd {
			imgSrc = imaging.Resize(imgSrc, attr.Wd, 0, imaging.Lanczos)
		}
		//高度>指定高度 防止负调整
		dy := imgSrc.Bounds().Dy()
		if dy > attr.Hd {
			imgSrc = imaging.Resize(imgSrc, 0, attr.Hd, imaging.Lanczos)
		}
		err = jpeg.Encode(dst, imgSrc, nil)
		if err != nil {
			return ctx.JSON(utils.ErrIpt("文件写入失败,请重试", err.Error()))
		}
	} else {
		_, err = io.Copy(dst, src)
		if err != nil {
			return ctx.JSON(utils.ErrIpt("文件写入失败,请重试", err.Error()))
		}
	}
	return ctx.JSON(utils.Succ("上传成功", conf.App.Endpoint+"/"+name))
}

// StatusGo doc
// @Auth
// @Tags kernel,app
// @Summary 获取服务器go信息
// @Security ApiKeyAuth
// @Success 200 {object} utils.Reply{data=model.StateGo} "返回数据"
// @Router /api/status/go [get]
// @Router /api/status/app [get]
func StatusGo(ctx echo.Context) error {
	mod := model.StateGo{
		ARCH:         runtime.GOARCH,
		OS:           runtime.GOOS,
		Version:      runtime.Version(),
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
	}

	return ctx.JSON(utils.Succ("系统信息", mod))
}

// ------------------------------------------------------ 记录 ------------------------------------------------------
// LogPage doc
// @Auth mgmt
// @Tags kernel,log
// @Summary 获取日志分页
// @Param query query model.Page true "请求数据"
// @Success 200 {object} utils.Reply{data=[]model.Log} "返回数据"
// @Router /api/log/page [get]
func LogPage(ctx echo.Context) error {
	ipt := &model.Page{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if err := ipt.Stat(); err != nil {
		return ctx.JSON(utils.ErrIpt("分页大小输入错误", err.Error()))
	}
	count := repo.LogCount()
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "count eq 0"))
	}
	mods, err := repo.LogPage(ipt.Pi, ipt.Ps)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "items eq 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, count))
}

func Init() {
	initCron()
}
