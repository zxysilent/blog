package kernel

import (
	"blog/conf"
	"blog/internal/model"
	"blog/pkg/reply"
	"blog/pkg/rndm"
	"errors"
	"image"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/logs"

	// 注册图片解码器，供 image.DecodeConfig 识别真实格式
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
	_ "golang.org/x/image/tiff"
	_ "golang.org/x/image/webp"
)

// 上传相关的安全约束集中在 model/upload.go：
//   - 扩展名白名单（拒绝 html/svg/js/php 等可在同源执行的类型）
//   - 体积上限（读取阶段就中断，不先落盘）
//   - 文件名一律由服务端生成，绝不使用客户端提供的名字拼路径

// uploadLimit 用 MaxBytesReader 包住请求体，超限时读取即失败
// 必须在下游 ctx.FormFile 之前调用，否则 multipart 已被完整读入内存
func uploadLimit(ctx echo.Context) {
	ctx.Request().Body = http.MaxBytesReader(ctx.Response(), ctx.Request().Body, model.UploadMaxSize)
}

// uploadErrSize 统一超限提示
func uploadErrSize(ctx echo.Context) error {
	return reply.Fail(ctx, "文件过大，单个文件不能超过 "+model.UploadSizeText(model.UploadMaxSize))
}

// uploadDir 生成按年月分层的目标目录
// 先建到「年月」这一级：早期实现只建到 yyyyMM 一级却写 yyyyMM/dd，
// 日期子目录不存在时会直接 os.Create 失败
func uploadDir(t time.Time) (string, string) {
	rel := t.Format("200601/02")
	abs := filepath.Join("static", "upload", filepath.FromSlash(rel))
	return rel, abs
}

// uploadName 生成落盘文件名：目录 + ULID + 扩展名
// ULID 单调递增，同毫秒并发也不会撞名，比原来的 8 位随机串可靠
func uploadName(rel, ext string) string {
	return "static/upload/" + rel + "/" + strings.ToLower(rndm.ULID()) + ext
}

// uploadSave 把上传内容写到 dst，失败时清理半成品
// 用显式 Close 而不是 defer，才能确定性感知写入错误（如磁盘写满）
func uploadSave(dst string, src io.Reader) error {
	os.Remove(dst) // 同名残留先清掉，避免追加写入
	file, err := os.Create(dst)
	if err != nil {
		return err
	}
	if _, err = io.Copy(file, src); err != nil {
		file.Close()
		os.Remove(dst)
		return err
	}
	if err = file.Close(); err != nil {
		os.Remove(dst)
		return err
	}
	return nil
}

// uploadTooLarge 判断是否为 MaxBytesReader 触发的超限错误
func uploadTooLarge(err error) bool {
	var maxErr *http.MaxBytesError
	return errors.As(err, &maxErr)
}

// uploadLog 把上传摘要写进上下文，供审计中间件记录
// 上传是 multipart 请求，审计中间件拿不到可读的请求体，只能由这里给出摘要
func uploadLog(ctx echo.Context, filename, ext, saved string, size int64) {
	ctx.Set(model.CtxAuditData, strings.Join([]string{
		"原始文件名=" + filename,
		"扩展名=" + ext,
		"大小=" + model.UploadSizeText(size),
		"存储路径=" + saved,
	}, "; "))
}

// UploadFile doc
// @Tags system
// @Summary 上传文件
// @Accept  mpfd
// @Param file formData file true "文件"
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=string} "返回数据"
// @Router /api/files [post]
func UploadFile(ctx echo.Context) error {
	rctx := ctx.Request().Context()
	uploadLimit(ctx)
	file, err := ctx.FormFile("file")
	if err != nil {
		if uploadTooLarge(err) {
			return uploadErrSize(ctx)
		}
		return reply.Fail(ctx, "未发现文件", err)
	}
	if file.Size > model.UploadMaxSize {
		return uploadErrSize(ctx)
	}
	// 客户端提供的文件名只用于取扩展名，绝不参与路径拼接
	ext, err := model.UploadCheck(file.Filename)
	if err != nil {
		return reply.Fail(ctx, err.Error())
	}
	src, err := file.Open()
	if err != nil {
		return reply.Fail(ctx, "文件打开失败", err)
	}
	defer src.Close()

	rel, absDir := uploadDir(time.Now())
	if err := os.MkdirAll(absDir, 0o755); err != nil {
		logs.Ctx(rctx).Err(err).Str("dir", absDir).Error("upload.mkdir")
		return reply.Fail(ctx, "创建目录失败")
	}
	name := uploadName(rel, ext)
	if err := uploadSave(name, src); err != nil {
		logs.Ctx(rctx).Err(err).Str("name", name).Error("upload.save")
		if uploadTooLarge(err) {
			return uploadErrSize(ctx)
		}
		return reply.Fail(ctx, "文件保存失败")
	}
	uploadLog(ctx, file.Filename, ext, name, file.Size)
	return reply.Succ(ctx, "上传成功", conf.App.Endpoint+"/"+name)
}

// UploadImage doc
// @Tags system
// @Summary 上传图片
// @Accept  mpfd
// @Param file formData file true "图片"
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=string} "返回数据"
// @Router /api/images [post]
func UploadImage(ctx echo.Context) error {
	rctx := ctx.Request().Context()
	uploadLimit(ctx)
	file, err := ctx.FormFile("file")
	if err != nil {
		if uploadTooLarge(err) {
			return uploadErrSize(ctx)
		}
		return reply.Fail(ctx, "未发现文件,请重试", err)
	}
	if file.Size > model.UploadMaxSize {
		return uploadErrSize(ctx)
	}
	ext, err := model.UploadImageCheck(file.Filename)
	if err != nil {
		return reply.Fail(ctx, err.Error())
	}
	src, err := file.Open()
	if err != nil {
		return reply.Fail(ctx, "文件打开失败,请重试", err)
	}
	defer src.Close()

	// 客户端声明的 Content-Type 完全可控，不能作为判据；
	// 这里用解码器读文件头，能真正确认是不是图片，同时拿到尺寸
	cfg, format, err := image.DecodeConfig(src)
	if err != nil {
		return reply.Fail(ctx, "文件内容不是有效的图片")
	}
	if _, err := src.Seek(0, io.SeekStart); err != nil {
		return reply.Fail(ctx, "文件读取失败,请重试", err)
	}
	if cfg.Width <= 0 || cfg.Height <= 0 {
		return reply.Fail(ctx, "图片尺寸异常")
	}

	rel, absDir := uploadDir(time.Now())
	if err := os.MkdirAll(absDir, 0o755); err != nil {
		logs.Ctx(rctx).Err(err).Str("dir", absDir).Error("upload.mkdir")
		return reply.Fail(ctx, "创建目录失败")
	}
	name := uploadName(rel, ext)
	if err := uploadSave(name, src); err != nil {
		logs.Ctx(rctx).Err(err).Str("name", name).Error("upload.save")
		if uploadTooLarge(err) {
			return uploadErrSize(ctx)
		}
		return reply.Fail(ctx, "文件写入失败,请重试")
	}
	uploadLog(ctx, file.Filename, ext, name, file.Size)
	logs.Ctx(rctx).Str("format", format).Int("w", cfg.Width).Int("h", cfg.Height).Debug("upload.image")
	return reply.Succ(ctx, "上传成功", conf.App.Endpoint+"/"+name)
}
