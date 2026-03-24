package kernel

import (
	"blog/conf"
	"blog/internal/model"
	"blog/internal/utils"
	"blog/pkg/reply"
	_ "image/gif"
	_ "image/png"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/logs"
	_ "golang.org/x/image/bmp"
)

// UploadFile doc
// @Auth
// @Tags kernel,upload
// @Summary 上传文件
// @Accept  mpfd
// @Param file formData file true "file"
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=string} "返回数据"
// @Router /api/upload/file [post]
func UploadFile(ctx echo.Context) error {
	rctx := utils.WrapCtx(ctx)
	file, err := ctx.FormFile("file")
	if err != nil {
		return reply.Fail(ctx, "未发现文件", err)
	}
	src, err := file.Open()
	if err != nil {
		return reply.Fail(ctx, "文件打开失败", err)

	}
	defer src.Close()
	dir := time.Now().Format("200601/02")
	os.MkdirAll("./static/upload/"+dir[:6], 0755)
	name := "static/upload/" + dir + utils.SUID() + path.Ext(file.Filename)
	dst, err := os.Create(name)
	if err != nil {
		logs.Ctx(rctx).Err(err).Error("os.Create")
		return reply.Fail(ctx, "文件打创建文件失败", err)
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	if err != nil {
		logs.Ctx(rctx).Err(err).Error("io.Copy")
		return reply.Fail(ctx, "文件保存失败", err)
	}
	return reply.Succ(ctx, "上传成功", conf.App.Endpoint+"/"+name)
}

// UploadImage doc
// @Auth
// @Tags kernel,upload
// @Summary 上传图片并裁剪
// @Accept  mpfd
// @Param file formData file true "file"
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=string} "返回数据"
// @Router /api/upload/image [post]
func UploadImage(ctx echo.Context) error {
	file, err := ctx.FormFile("file")
	if err != nil {
		return reply.Fail(ctx, "未发现文件,请重试", err)
	}
	if !strings.Contains(file.Header.Get("Content-Type"), "image") {
		return reply.Fail(ctx, "请选择图片文件", file.Header.Get("Content-Type"))
	}
	src, err := file.Open()
	if err != nil {
		return reply.Fail(ctx, "文件打开失败,请重试", err)
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
	// if conf.App.ImageCut || attr.Cut {
	// 	ext = ".jpg"
	// }
	name := "static/upload/" + dir + utils.SUID() + ext
	dst, err := os.Create(name)
	if err != nil {
		return reply.Fail(ctx, "目标文件创建失败,请重试", err)
	}
	defer dst.Close()
	// if conf.App.ImageCut || attr.Cut { //图片裁剪
	// 	imgSrc, _, err := image.Decode(src)
	// 	// 图片解码
	// 	if err != nil {
	// 		return reply.Fail(ctx, "读取图片失败,请重试", err)
	// 	}
	// 	if attr.Wd <= 0 {
	// 		attr.Wd = conf.App.ImageWidth
	// 	}
	// 	if attr.Hd <= 0 {
	// 		attr.Hd = conf.App.ImageHeight
	// 	}
	// 	// 宽度>指定宽度 防止负调整
	// 	dx := imgSrc.Bounds().Dx()
	// 	if dx > attr.Wd {
	// 		imgSrc = imaging.Resize(imgSrc, attr.Wd, 0, imaging.Lanczos)
	// 	}
	// 	//高度>指定高度 防止负调整
	// 	dy := imgSrc.Bounds().Dy()
	// 	if dy > attr.Hd {
	// 		imgSrc = imaging.Resize(imgSrc, 0, attr.Hd, imaging.Lanczos)
	// 	}
	// 	err = jpeg.Encode(dst, imgSrc, nil)
	// 	if err != nil {
	// 		return reply.Fail(ctx, "文件写入失败,请重试", err)
	// 	}
	// } else {
	_, err = io.Copy(dst, src)
	if err != nil {
		return reply.Fail(ctx, "文件写入失败,请重试", err)
	}
	// }
	return reply.Succ(ctx, "上传成功", conf.App.Endpoint+"/"+name)
}

// StatusGo doc
// @Auth
// @Tags kernel,app
// @Summary 获取服务器go信息
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=model.StateGo} "返回数据"
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
	return reply.Succ(ctx, "系统信息", mod)
}

func Init() {
	initCron()
}
