package control

import (
	"blog/model"
	"blog/util"
	"net/http"
	"strconv"

	"github.com/labstack/echo"
)

// Index 主页面
func Index(ctx echo.Context) error {
	//return ctx.HTML(200, `<html><head><meta charset="UTF-8"><title>文档</title></head><body><a href="/swagger/index.html">doc</a></body></html>`)
	pi, _ := strconv.Atoi(ctx.FormValue("page"))
	if pi == 0 {
		pi = 1
	}
	ps := util.Atoi(model.MapOpts.MustGet("page_size"), 6)
	mods, _ := model.PostPage(pi, ps)
	total := model.PostCount()
	naver := model.Naver{}
	if pi > 1 {
		naver.Prev = "/?page=1"
	}
	if total > (pi * ps) {
		naver.Next = "/?page=" + strconv.Itoa(pi+1)
	}
	return ctx.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Posts": mods,
		"Naver": naver,
	})
}

// Archives 归档
func Archives(ctx echo.Context) error {
	mods, err := model.PostArchive()
	if err != nil {
		return ctx.Redirect(302, "/")
	}
	return ctx.Render(http.StatusOK, "archive.html", map[string]interface{}{
		"Posts": mods,
	})
}
