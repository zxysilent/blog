package appctl

import (
	"blog/model"
	"net/http"
	"regexp"

	"github.com/labstack/echo/v4"
)

var regImg = regexp.MustCompile(`<img src="([^" ]+)" alt="([^" ]*)"\s?\/?>`)
var regToc = regexp.MustCompile("<h[1-6]>.*?</h[1-6]>")
var regH = regexp.MustCompile(`<h[1-6]><a id="(.*?)"></a>(.*?)</h[1-6]>`)

// RenderPage 渲染页面
func RenderPage(path string, ctx echo.Context) error {
	mod, has := model.PostSingle(path)
	if !has {
		return ctx.Redirect(302, "/")
	}
	mod.Richtext = regImg.ReplaceAllString(mod.Richtext, `<img class="lazy-load" src="data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==" data-src="$1" alt="$2">`)
	return ctx.Render(http.StatusOK, "page.html", map[string]interface{}{
		"Page": mod,
		"Show": mod.Status == 2,
	})
}
