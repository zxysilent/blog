package router

import (
	"blog/internal/model"
	"blog/internal/repo"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// 审计日志只记录「敏感操作」：
//
//	1. 声明了权限标识的写请求（POST/PUT/PATCH/DELETE）——改数据必然敏感
//	2. 登录接口 POST /api/tokens         —— 登录成功/失败都要留痕
//
// 普通查询、公开配置读取等只读接口不记录，避免日志表被淹没。
// 写入走异步队列（repo.LogPush），日志表故障不会影响业务。

// auditMaxBody 请求参数最多保留的字节数，超出截断
// 大字段（如文章正文）没必要进日志，只保留足够定位的信息
const auditMaxBody = 2 << 10

// auditMaxResp 响应体最多解析的字节数，只需要够读到业务码
const auditMaxResp = 512

// auditSkipPaths 永不记录的路径
var auditSkipPaths = map[string]struct{}{
	"/api/captcha": {}, //验证码
}

// auditEnabled 审计开关
// 取自 config.auth 的 audit 字段，走内存缓存，热更新即时生效
func auditEnabled() bool {
	return repo.GetCfg().Auth.Audit
}

// Audit 中间件-操作审计
// 由配置 config.auth 的 audit 字段控制开关（后台可热切换）：
// 关闭时直接放行，不读请求体、不建索引、无任何额外开销。
func Audit(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		if !auditEnabled() {
			return next(ctx)
		}
		req := ctx.Request()
		path := ctx.Path()
		if _, skip := auditSkipPaths[path]; skip {
			return next(ctx)
		}
		perm := routePerm(ctx)
		login := req.Method == http.MethodPost && req.URL.Path == "/api/tokens"
		// 只记录敏感的写操作与登录
		if !login && (perm == "" || !isWriteMethod(req.Method)) {
			return next(ctx)
		}

		// 读取请求体后必须放回，否则下游 ctx.Bind 拿不到数据
		// multipart（上传）不读：整包读进内存代价太大，且长度限制要靠
		// http.MaxBytesReader 在真正的读取阶段生效，提前读完等于绕过了它。
		// 这类请求的审计内容由业务处理器通过 model.CtxAuditData 给出摘要。
		var body []byte
		if !isMultipart(req) {
			body = readAndRestoreBody(req)
		}
		// 捕获响应体用于提取业务码
		// v4 的 Response() 是 *echo.Response，直接换掉内层 Writer 即可
		res := ctx.Response()
		writer := &auditWriter{ResponseWriter: res.Writer}
		res.Writer = writer

		err := next(ctx)

		mod := &model.Log{
			Perm:    perm,
			Method:  req.Method,
			Path:    req.URL.Path,
			Ip:      ctx.RealIP(),
			Status:  writer.bizCode(),
			Data:    auditBody(body),
			Created: time.Now().UnixMilli(),
		}
		// 上传等非 JSON 请求体无法解析，由业务处理器给出摘要，优先采用
		if extra, ok := ctx.Get(model.CtxAuditData).(string); ok && extra != "" {
			mod.Data = extra
		}
		if uid, ok := ctx.Get("uid").(string); ok {
			mod.UserId = uid
		}
		if perm != "" {
			if def, ok := parsePermName(routeName(ctx)); ok {
				mod.Module = def.Group
				mod.Action = def.Name
			}
		} else if login {
			mod.Module = "账号"
			mod.Action = "登录"
		}
		repo.LogPush(mod)
		return err
	}
}

// isWriteMethod 是否为写方法
func isWriteMethod(method string) bool {
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		return true
	}
	return false
}

// routeName 取当前请求匹配到的路由 Name
// 请求期只需要分组与中文名，复用 BuildRouteNames 已经建好的索引
func routeName(ctx echo.Context) string {
	return routeNames[ctx.Path()][ctx.Request().Method]
}

// isMultipart 是否为 multipart 请求（上传）
func isMultipart(req *http.Request) bool {
	return strings.HasPrefix(req.Header.Get(echo.HeaderContentType), echo.MIMEMultipartForm)
}

// readAndRestoreBody 读取请求体并写回
// 审计只读一次，但下游 handler 仍要能正常 Bind，因此必须完整还原。
// 这里不做长度截断：截断会导致下游解析失败，长度限制只作用于入库内容。
func readAndRestoreBody(req *http.Request) []byte {
	if req.Body == nil {
		return nil
	}
	full, err := io.ReadAll(req.Body)
	if err != nil {
		return nil
	}
	req.Body = io.NopCloser(bytes.NewReader(full))
	return full
}

// auditWriter 包装响应写入器，顺手截留前若干字节
type auditWriter struct {
	http.ResponseWriter
	buf []byte
}

func (w *auditWriter) Write(b []byte) (int, error) {
	if len(w.buf) < auditMaxResp {
		remain := auditMaxResp - len(w.buf)
		if remain > len(b) {
			remain = len(b)
		}
		w.buf = append(w.buf, b[:remain]...)
	}
	return w.ResponseWriter.Write(b)
}

// bizCode 从响应体里提取业务码
// 本项目所有接口都是 HTTP 200 + body.code，取不到时返回 0
func (w *auditWriter) bizCode() int {
	if len(w.buf) < 1 {
		return 0
	}
	// 目标形态固定为 {"code":200,...}，直接找键名比完整解析更省
	idx := strings.Index(string(w.buf), `"code":`)
	if idx < 0 {
		return 0
	}
	rest := string(w.buf)[idx+len(`"code":`):]
	end := strings.IndexAny(rest, ",}")
	if end < 0 {
		return 0
	}
	code, err := strconv.Atoi(strings.TrimSpace(rest[:end]))
	if err != nil {
		return 0
	}
	return code
}

// auditBody 生成入库用的参数快照
// 非 JSON 或空体不入库；命中敏感字段的值整体替换为 ***
func auditBody(raw []byte) string {
	// 纯空白体等价于空体，不必记成「非 JSON」
	raw = bytes.TrimSpace(raw)
	if len(raw) < 1 {
		return ""
	}
	var parsed any
	if err := json.Unmarshal(raw, &parsed); err != nil {
		// 非 JSON（如上传场景）只记长度，不记内容
		return "非 JSON 请求体，长度 " + strconv.Itoa(len(raw))
	}
	masked := auditMask(parsed)
	out, err := json.Marshal(masked)
	if err != nil {
		return ""
	}
	if len(out) > auditMaxBody {
		return string(out[:auditMaxBody])
	}
	return string(out)
}

// auditMask 递归脱敏
func auditMask(v any) any {
	switch val := v.(type) {
	case map[string]any:
		out := make(map[string]any, len(val))
		for k, sub := range val {
			if model.LogNeedMask(k) {
				out[k] = model.LogMaskedValue
				continue
			}
			out[k] = auditMask(sub)
		}
		return out
	case []any:
		out := make([]any, 0, len(val))
		for _, sub := range val {
			out = append(out, auditMask(sub))
		}
		return out
	default:
		return v
	}
}
