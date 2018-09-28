package control

import (
	"blog/model"
	"blog/util"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo"
)

var reg = regexp.MustCompile(`<img src="([^" ]+)" alt="([^" ]*)"\s?\/?>`)

// Post 主页面
func Post(ctx echo.Context) error {
	//return ctx.HTML(200, `<html><head><meta charset="UTF-8"><title>文档</title></head><body><a href="/swagger/index.html">doc</a></body></html>`)
	paths := strings.Split(ctx.Param("*"), ".")
	if len(paths) == 2 {
		mod, naver, has := model.PostPath(paths[0])
		if !has {
			return ctx.Redirect(302, "/")
		}
		if paths[1] == "html" {
			mod.Content = reg.ReplaceAllString(mod.Content, `<img class="lazy-load" src="data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==" data-src="$1" alt="$2">`)
			tags, _ := model.PostTags(mod.Id)
			cates, _ := model.PostCates(mod.Id)
			return ctx.Render(http.StatusOK, "post.html", map[string]interface{}{
				"Post":    mod,
				"Naver":   naver,
				"Tags":    tags,
				"HasTag":  len(tags) > 0,
				"Cates":   cates,
				"HasCate": len(cates) > 0,
			})
		}
		return ctx.Res(util.NewSucc("", mod))
	}
	return nil
}

// PostPageAll 页面列表
// @Success 200 {object} util.Result "成功数据"
func PostPageAll(ctx echo.Context) error {
	mods, err := model.PostPageAll()
	if err != nil {
		return ctx.Res(util.NewErrOpt(`未查询到页面信息`, err.Error()))
	}
	if len(mods) < 1 {
		return ctx.Res(util.NewErrOpt(`未查询到页面信息`, "len"))
	}
	return ctx.Res(util.NewSucc(`页面信息`, mods))
}
