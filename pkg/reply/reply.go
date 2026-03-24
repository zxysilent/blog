package reply

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Reply  format
// 统一返回格式
type Reply struct {
	Code  int    `json:"code"`
	Msg   string `json:"msg"`
	Data  any    `json:"data,omitempty"`
	Count int    `json:"count,omitempty"`
}

const (
	CodeSucc     int = 200 //正常
	CodeFail     int = 300 //失败
	CodeErrDeny  int = 340 //没有权限
	CodeErrToken int = 360 //token错误
)

func newReply(code int, msg string, data ...any) Reply {
	reply := Reply{
		Code: code,
		Msg:  msg,
	}
	if len(data) > 0 {
		reply.Data = data[0]
	}
	if err, ok := reply.Data.(error); ok && err != nil {
		reply.Data = err.Error()
	}
	return reply
}

// Succ 返回一个成功标识的结果格式
func Succ(ctx echo.Context, msg string, data ...any) error {
	return ctx.JSON(http.StatusOK, newReply(CodeSucc, msg, data...))
}

// Fail 返回一个失败标识的结果格式
func Fail(ctx echo.Context, msg string, data ...any) error {
	return ctx.JSON(http.StatusOK, newReply(CodeFail, msg, data...))
}

// Page 返回一个带有分页数据的结果格式
func List(ctx echo.Context, msg string, data any, count int) error {
	return ctx.JSON(http.StatusOK, Reply{
		Code:  CodeSucc,
		Msg:   msg,
		Data:  data,
		Count: count,
	})
}

// ErrDeny 返回一个没有权限的结果
func ErrDeny(ctx echo.Context, msg string, data ...any) error {
	return ctx.JSON(http.StatusOK, newReply(CodeErrDeny, msg, data...))
}

// ErrToken 返回一个验证失败的结果
func ErrToken(ctx echo.Context, msg string, data ...any) error {
	return ctx.JSON(http.StatusOK, newReply(CodeErrToken, msg, data...))
}
