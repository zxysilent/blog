package sysctl

import (
	"blog/model"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// StatusGoinfo doc
// @Tags status
// @Summary 获取服务器go信息
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=model.Goinfo} "成功数据"
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

// 上传 下载相关
