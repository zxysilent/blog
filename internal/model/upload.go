package model

import (
	"fmt"
	"path"
	"strings"
)

// 上传限制与白名单
//
// 上传是少数直接把外部数据落到磁盘的入口，风险集中在三处：
//  1. 类型：扩展名决定静态服务器的 Content-Type，放行 .html/.svg 等于给自己开了存储型 XSS
//  2. 体积：不限长的话一个请求就能把磁盘写满
//  3. 覆盖：文件名必须由服务端生成，绝不能用客户端提供的名字拼路径

const (
	// UploadMaxSize 单次上传的最大字节数（32MB）
	// 超过时由 http.MaxBytesReader 在读取阶段就中断，不会先落盘再判断
	UploadMaxSize int64 = 32 << 20

	// UploadNameLen 生成文件名的长度（ULID，26 位）
	UploadNameLen = 26

	// UploadNameTextLen 审计记录中原始文件名的最大保留长度（字符数）
	UploadNameTextLen = 96
)

// UploadExtAllow 允许的扩展名白名单
//
// 刻意不放行的类型及原因：
//   - html / htm / xhtml / svg：同源静态目录下会被浏览器当页面执行，可直接窃取 token
//   - js / mjs / css：同上，可被当作脚本或样式加载
//   - php / jsp / asp / aspx / cgi / sh：一旦静态目录被反向代理交给解释器就等于 RCE
//   - exe / bat / cmd / com / dll / so：可执行文件，诱导下载后运行
//   - jar / war / class：可执行归档
var UploadExtAllow = map[string]struct{}{
	// 图片
	".jpg": {}, ".jpeg": {}, ".png": {}, ".gif": {}, ".webp": {}, ".bmp": {},
	".ico": {}, ".avif": {}, ".tiff": {}, ".tif": {},
	// 文档
	".pdf": {}, ".txt": {}, ".md": {}, ".csv": {},
	".doc": {}, ".docx": {}, ".xls": {}, ".xlsx": {}, ".ppt": {}, ".pptx": {},
	// 压缩包
	".zip": {}, ".rar": {}, ".7z": {}, ".tar": {}, ".gz": {}, ".tgz": {}, ".bz2": {}, ".xz": {},
	// 音视频
	".mp3": {}, ".wav": {}, ".flac": {}, ".aac": {}, ".ogg": {},
	".mp4": {}, ".webm": {}, ".mov": {}, ".m4v": {}, ".mkv": {},
}

// UploadImageExt 允许的图片扩展名
// 与 UploadExtAllow 的图片子集保持一致，用于图片上传接口的二次校验
var UploadImageExt = map[string]struct{}{
	".jpg": {}, ".jpeg": {}, ".png": {}, ".gif": {}, ".webp": {}, ".bmp": {},
	".ico": {}, ".avif": {}, ".tiff": {}, ".tif": {},
}

// UploadDangerExt 明确拒绝的扩展名，仅用于给出更准确的错误提示
var UploadDangerExt = map[string]struct{}{
	".html": {}, ".htm": {}, ".xhtml": {}, ".shtml": {}, ".svg": {}, ".svgz": {},
	".js": {}, ".mjs": {}, ".css": {},
	".php": {}, ".php3": {}, ".php4": {}, ".php5": {}, ".phtml": {}, ".pht": {},
	".jsp": {}, ".jspx": {}, ".asp": {}, ".aspx": {}, ".ashx": {}, ".cgi": {},
	".sh": {}, ".bash": {}, ".zsh": {}, ".py": {}, ".rb": {}, ".pl": {},
	".exe": {}, ".bat": {}, ".cmd": {}, ".com": {}, ".dll": {}, ".so": {}, ".msi": {},
	".jar": {}, ".war": {}, ".class": {},
}

// UploadExt 取文件名的小写扩展名
// 只取最后一段，避免 "a.tar.gz" 被当成 ".gz" 之外的其它值
func UploadExt(filename string) string {
	ext := strings.ToLower(path.Ext(strings.TrimSpace(filename)))
	return ext
}

// UploadExtValid 扩展名是否在通用白名单内
func UploadExtValid(ext string) bool {
	_, ok := UploadExtAllow[ext]
	return ok
}

// UploadImageExtValid 扩展名是否在图片白名单内
func UploadImageExtValid(ext string) bool {
	_, ok := UploadImageExt[ext]
	return ok
}

// UploadExtError 生成可读的错误提示
// 危险类型单独提示，便于排查是「选错文件」还是「有人在试」
func UploadExtError(ext string) error {
	if ext == "" {
		return fmt.Errorf("文件缺少扩展名，无法识别类型")
	}
	if _, ok := UploadDangerExt[ext]; ok {
		return fmt.Errorf("出于安全考虑，不允许上传 %s 类型的文件", ext)
	}
	return fmt.Errorf("不支持的文件类型 %s", ext)
}

// UploadCheck 通用上传校验：扩展名必须在白名单内
func UploadCheck(filename string) (string, error) {
	ext := UploadExt(filename)
	if !UploadExtValid(ext) {
		return "", UploadExtError(ext)
	}
	return ext, nil
}

// UploadImageCheck 图片上传校验：扩展名必须在图片白名单内
// 注意这里只校验名字，文件内容是否为真图片由 image.DecodeConfig 在服务层判定
func UploadImageCheck(filename string) (string, error) {
	ext := UploadExt(filename)
	if !UploadImageExtValid(ext) {
		return "", UploadExtError(ext)
	}
	return ext, nil
}

// UploadFileNameText 把客户端提供的文件名整理成可安全记录的形式
// 原始文件名完全由调用方控制，含换行/控制字符时可能污染日志展示甚至伪造日志行
func UploadFileNameText(filename string) string {
	name := strings.Map(func(r rune) rune {
		if r < 0x20 || r == 0x7f {
			return -1
		}
		return r
	}, strings.TrimSpace(filename))
	name = strings.TrimSpace(name)
	if name == "" {
		return "(未命名)"
	}
	runes := []rune(name)
	if len(runes) > UploadNameTextLen {
		name = string(runes[:UploadNameTextLen]) + "..."
	}
	return name
}

// UploadSizeText 把字节数转成便于阅读的文本
func UploadSizeText(size int64) string {
	const unit = 1024
	if size < unit {
		return fmt.Sprintf("%dB", size)
	}
	div, exp := int64(unit), 0
	for n := size / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f%cB", float64(size)/float64(div), "KMGTPE"[exp])
}

// CtxAuditData 上下文键：业务处理器写入的审计摘要
// 上传等非 JSON 请求无法从请求体提取有效信息，由处理器自己给出摘要，
// 审计中间件优先使用它作为 Data。
const CtxAuditData = "audit.data"
