package control

import (
	"blog/model"
	"net/http"

	"github.com/labstack/echo/v4"
)

// AboutView 关于
func AboutView(ctx echo.Context) error {
	mod, has := model.PostSingle("about")
	if !has {
		return ctx.Redirect(302, "/")
	}
	mod.Content = reg.ReplaceAllString(mod.Content, `<img class="lazy-load" src="data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==" data-src="$1" alt="$2">`)
	return ctx.Render(http.StatusOK, "page.html", map[string]interface{}{
		"Page": mod,
		"Show": mod.IsPublic && mod.Status == 3,
	})
}

// LinksView 友链
func LinksView(ctx echo.Context) error {
	mod, has := model.PostSingle("links")
	if !has {
		return ctx.Redirect(302, "/")
	}
	mod.Content = reg.ReplaceAllString(mod.Content, `<img class="lazy-load" src="data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==" data-src="$1" alt="$2">`)
	return ctx.Render(http.StatusOK, "page.html", map[string]interface{}{
		"Page": mod,
		"Show": mod.IsPublic && mod.Status == 3,
	})
}

// PageView 页面
func PageView(ctx echo.Context) error {
	mod, has := model.PostSingle(ctx.Param("*"))
	if !has {
		return ctx.Redirect(302, "/")
	}
	mod.Content = reg.ReplaceAllString(mod.Content, `<img class="lazy-load" src="data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==" data-src="$1" alt="$2">`)
	return ctx.Render(http.StatusOK, "page.html", map[string]interface{}{
		"Page": mod,
		"Show": mod.IsPublic && mod.Status == 3,
	})
}
