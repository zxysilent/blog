package utils

// Reply  format
// 统一返回格式
type Reply struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

// page format Message
// 数据分页附加数据
type page struct {
	Count int         `json:"count"`
	Items interface{} `json:"items"`
}

const (
	codeSucc     int = 200 //正常
	codeFail     int = 300 //失败
	codeErrIpt   int = 310 //输入数据有误
	codeErrOpt   int = 320 //无数据返回
	codeErrDeny  int = 330 //没有权限
	codeErrToken int = 340 //token错误
	codeErrSvr   int = 350 //服务端错误
	codeExt      int = 400 //其他约定,eg 更新token
)

func newReply(code int, msg string, data ...interface{}) (int, Reply) {
	if len(data) > 0 {
		return 200, Reply{
			Code: code,
			Msg:  msg,
			Data: data[0],
		}
	}
	return 200, Reply{
		Code: code,
		Msg:  msg,
	}
}

// Succ 返回一个成功标识的结果格式
func Succ(msg string, data ...interface{}) (int, Reply) {
	return newReply(codeSucc, msg, data...)
}

// Fail 返回一个失败标识的结果格式
func Fail(msg string, data ...interface{}) (int, Reply) {
	return newReply(codeFail, msg, data...)
}

// Page 返回一个带有分页数据的结果格式
func Page(msg string, items interface{}, count int) (int, Reply) {
	return 200, Reply{
		Code: codeSucc,
		Msg:  msg,
		Data: page{
			Items: items,
			Count: count,
		},
	}
}

// ErrIpt 返回一个输入错误的结果
func ErrIpt(msg string, data ...interface{}) (int, Reply) {
	return newReply(codeErrIpt, msg, data...)
}

// ErrOpt 返回一个输出错误的结果
func ErrOpt(msg string, data ...interface{}) (int, Reply) {
	return newReply(codeErrOpt, msg, data...)
}

// ErrDeny 返回一个没有权限的结果
func ErrDeny(msg string, data ...interface{}) (int, Reply) {
	return newReply(codeErrDeny, msg, data...)
}

// ErrToken 返回一个验证失败的结果
func ErrToken(msg string, data ...interface{}) (int, Reply) {
	return newReply(codeErrToken, msg, data...)
}

// ErrSvr 返回一个服务端错误的结果
func ErrSvr(msg string, data ...interface{}) (int, Reply) {
	return newReply(codeErrSvr, msg, data...)
}

// Ext 返回一个其他约定的结果
func Ext(msg string, data ...interface{}) (int, Reply) {
	return newReply(codeExt, msg, data...)
}
