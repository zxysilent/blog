package sysctl

import (
	"blog/model"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// StatusGoinfo doc
// @Tags status-状态信息
// @Summary 获取服务器go信息
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=model.Goinfo} "返回数据"
// @Router /adm/status/goinfo [get]
func StatusGoinfo(ctx echo.Context) error {
	mod := model.Goinfo{
		ARCH:    runtime.GOARCH,
		OS:      runtime.GOOS,
		Version: runtime.Version(),
		NumCPU:  runtime.NumCPU(),
	}
	return ctx.JSON(utils.Succ("系统信息", mod))
}

// StatusApp doc
// @Tags status-状态信息
// @Summary 获取服务器go信息
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=model.State} "返回数据"
// @Router /adm/status/app [get]
func StatusAppinfo(ctx echo.Context) error {
	if mod, has := model.Collect(); has {
		return ctx.JSON(utils.Succ(`统计信息`, mod))
	}
	return ctx.JSON(utils.Fail(`未查询到统计信息`))
}

// 上传 下载相关
// Collect 统计信息
