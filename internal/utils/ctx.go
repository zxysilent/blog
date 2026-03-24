package utils

import (
	"cmp"
	"context"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/logs"
)

func WrapCtx(rctx echo.Context) context.Context {
	header := rctx.Request().Header
	ctx := logs.TraceCtx(context.Background(), cmp.Or(header.Get(echo.HeaderXRequestID), logs.TraceId()))
	return ctx
}
