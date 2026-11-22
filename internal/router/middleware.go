package router

import (
	"blog/internal/repo"
	"blog/pkg/page"
	"blog/pkg/reply"
	"blog/pkg/token"
	"context"
	"errors"
	"fmt"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/zxysilent/logs"
)

// Recover 中间件-panic 恢复
func Recover(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}
				tctx := logs.TraceCtx(context.Background(), ctx.Request().Header.Get(echo.HeaderXRequestID))
				logs.Ctx(tctx).Raw("panic", debug.Stack()).Err(err).Error("recover")
				// v4 的错误处理顺序是 (err, ctx)
				ctx.Echo().HTTPErrorHandler(err, ctx)
			}
		}()
		return next(ctx)
	}
}

// Logger 中间件-日志记录
func Logger(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) (err error) {
		start := time.Now()
		req := ctx.Request()
		res := ctx.Response()
		rid := req.Header.Get(echo.HeaderXRequestID)
		if rid == "" {
			rid = logs.TraceId()
		}
		req.Header.Set(echo.HeaderXRequestID, rid)
		res.Header().Set(echo.HeaderXRequestID, rid)
		// 把 trace 上下文注入 request，下游统一用 ctx.Request().Context()
		rctx := logs.TraceCtx(context.Background(), rid)
		ctx.SetRequest(req.WithContext(rctx))
		logs.With().Caller(false).Str("trace", rid).Str("caller", "api-begin").Str("ip", ctx.RealIP()).Str(req.Method, req.RequestURI).Debug()
		if err = next(ctx); err != nil {
			ctx.Echo().HTTPErrorHandler(err, ctx)
		}
		logs.With().Caller(false).Str("trace", rid).Str("caller", "api-end").Int("status", res.Status).Int64("size", res.Size).Dur("dur", time.Since(start)).Info()
		return
	}
}

// HTTPErrorHandler 全局错误捕捉
func HTTPErrorHandler(err error, ctx echo.Context) {
	if ctx.Response().Committed {
		return
	}
	path := ctx.Request().URL.Path
	isApi := strings.HasPrefix(path, "/api")
	msg := err.Error()
	code := http.StatusInternalServerError
	var he *echo.HTTPError
	if errors.As(err, &he) {
		code = he.Code
		if m, ok := he.Message.(string); ok {
			msg = m
		}
	}
	if code == http.StatusNotFound {
		switch {
		case strings.HasPrefix(path, "/static"), strings.HasPrefix(path, "/dist"):
			ctx.NoContent(http.StatusNotFound)
		case isApi:
			reply.Fail(ctx, "接口不存在", msg)
		default:
			ctx.HTML(http.StatusNotFound, page.Err404().SetExtra("404").Render())
		}
		return
	}
	if isApi {
		reply.Fail(ctx, "系统错误", msg)
		return
	}
	ctx.HTML(http.StatusBadGateway, page.Error().SetExtra(msg).Render())
}

// crosConfig 跨域配置
var crosConfig = middleware.CORSConfig{
	AllowOrigins: []string{"*"},
	AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization, "user-token", "accept-serial", "token"},
	AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
}

// rateConfig 限流配置，仅对 /api 生效
var rateConfig = middleware.RateLimiterConfig{
	Skipper: func(ctx echo.Context) bool {
		return !strings.HasPrefix(ctx.Request().URL.Path, "/api")
	},
	Store: middleware.NewRateLimiterMemoryStoreWithConfig(middleware.RateLimiterMemoryStoreConfig{Rate: 20, Burst: 5, ExpiresIn: 3 * time.Minute}),
	IdentifierExtractor: func(ctx echo.Context) (string, error) {
		return ctx.Request().URL.Path, nil
	},
	ErrorHandler: func(ctx echo.Context, err error) error {
		return reply.Fail(ctx, "请求错误,请稍后重试", err)
	},
	DenyHandler: func(ctx echo.Context, identifier string, err error) error {
		return reply.Fail(ctx, "当前请求过多,请稍后重试", identifier)
	},
}

// AuthToken 中间件-登录认证 + 基于路由权限标识的校验
// 路由未声明权限的接口只校验登录态；根角色直通
func AuthToken(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		raw := extractToken(ctx)
		if raw == "" {
			return reply.ErrToken(ctx, "请重新登陆", "未发现token")
		}
		tk := token.Token[token.Auth]{}
		if err := tk.Decode(raw, repo.GetCfg().Auth.TokenSecret); err != nil {
			return reply.ErrToken(ctx, "请重新登陆", err.Error())
		}
		claim := tk.Claim
		ctx.Set("auth", claim)
		ctx.Set("uid", claim.UserId)
		ctx.Set("rid", claim.RoleId)
		if !repo.RolePermAllow(ctx.Request().Context(), claim.RoleId, routePerm(ctx)) {
			return reply.ErrDeny(ctx, "没有权限")
		}
		return next(ctx)
	}
}
