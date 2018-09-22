package control

import (
	"blog/model"
	"net/http"

	"github.com/labstack/echo"
)

// About 关于
func About(ctx echo.Context) error {
	mod, has := model.PostSingle("about")
	if !has {
		return ctx.Redirect(302, "/")
	}
	mod.Content = reg.ReplaceAllString(mod.Content, `<img class="lazy-load" src="data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==" data-src="$1" alt="$2">`)
	return ctx.Render(http.StatusOK, "page.html", map[string]interface{}{
		"Page": mod,
	})
}

// Links 关于
func Links(ctx echo.Context) error {
	mod, has := model.PostSingle("links")
	if !has {
		return ctx.Redirect(302, "/")
	}
	mod.Content = reg.ReplaceAllString(mod.Content, `<img class="lazy-load" src="data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==" data-src="$1" alt="$2">`)
	return ctx.Render(http.StatusOK, "page.html", map[string]interface{}{
		"Page": mod,
	})
}

// Page 页面
func Page(ctx echo.Context) error {
	mod, has := model.PostSingle(ctx.Param("*"))
	if !has {
		return ctx.Redirect(302, "/")
	}
	mod.Content = reg.ReplaceAllString(mod.Content, `<img class="lazy-load" src="data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==" data-src="$1" alt="$2">`)
	return ctx.Render(http.StatusOK, "page.html", map[string]interface{}{
		"Page": mod,
	})
}
