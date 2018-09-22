package control

import (
	"blog/model"
	"net/http"

	"github.com/labstack/echo"
)

// Tags 标签列表
func Tags(ctx echo.Context) error {
	mods, err := model.TagAll()
	if err != nil {
		return ctx.Redirect(302, "/")
	}
	return ctx.Render(http.StatusOK, "tags.html", map[string]interface{}{
		"Tags": mods,
	})
}

// Tag 标签文章列表
func Tag(ctx echo.Context) error {
	mods, err := model.TagAll()
	if err != nil {
		return ctx.Redirect(302, "/")
	}
	return ctx.Render(http.StatusOK, "tags.html", map[string]interface{}{
		"Tags": mods,
	})
}
