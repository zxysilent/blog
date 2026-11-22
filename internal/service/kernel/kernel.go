package kernel

import (
	"blog/internal/model"
	"blog/pkg/reply"
	"runtime"

	"github.com/labstack/echo/v4"
)

// Status doc
// @Tags system
// @Summary 获取服务器运行信息
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=model.StateGo} "返回数据"
// @Router /api/status [get]
func Status(ctx echo.Context) error {
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
