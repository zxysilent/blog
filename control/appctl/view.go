package appctl

import (
	"blog/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// ------------------------------------------------------ 主页面 ------------------------------------------------------
// ViewIndex 主页面
func ViewIndex(ctx echo.Context) error {
	pi, _ := strconv.Atoi(ctx.FormValue("page"))
	if pi == 0 {
		pi = 1
	}
	ps := model.Gcfg().PageSize
	mods, _ := model.PostPage(-1, model.PostKindPost, pi, ps, "id", "title", "path", "created", "summary")
	total := model.PostCount(-1, model.PostKindPost)
	naver := model.Naver{}
	if pi > 1 {
		naver.Prev = "/?page=" + strconv.Itoa(pi-1)
	}
	if total > (pi * ps) {
		naver.Next = "/?page=" + strconv.Itoa(pi+1)
	}
	return ctx.Render(http.StatusOK, "index.html", map[string]interface{}{
		"Posts": mods,
		"Naver": naver,
	})
}

// ------------------------------------------------------ 文章页面 ------------------------------------------------------
// ViewPost 文章页面
func ViewPost(ctx echo.Context) error {
	paths := strings.Split(ctx.Param("*"), ".")
	if len(paths) == 2 {
		mod, naver, has := model.PostPath(paths[0])
		if !has {
			return ctx.Redirect(302, "/")
		}
		if paths[1] == "html" {
			mod.Richtext = regImg.ReplaceAllString(mod.Richtext, `<img class="lazy-load" src="data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==" data-src="$1" alt="$2">`)
			return ctx.Render(http.StatusOK, "post.html", map[string]interface{}{
				"Post":  mod,
				"Show":  mod.Status == 2,
				"Naver": naver,
			})
		}
		return ctx.JSON(utils.Succ("", mod))
	}
	return ctx.Redirect(302, "/404")
}

// ------------------------------------------------------ 单个页面 ------------------------------------------------------
// ViewAbout 关于页面
func ViewAbout(ctx echo.Context) error {
	return RenderPage("about", ctx)
}

// ViewLinks 友链页面
func ViewLinks(ctx echo.Context) error {
	return RenderPage("links", ctx)
}

// ViewPage 其它页面
func ViewPage(ctx echo.Context) error {
	paths := strings.Split(ctx.Param("*"), ".")
	if len(paths) == 2 {
		return RenderPage(paths[0], ctx)
	}
	return ctx.Redirect(302, "/404")

}

// ------------------------------------------------------ 归档页面 ------------------------------------------------------
// ViewArchives 归档页面
func ViewArchives(ctx echo.Context) error {
	archives, err := model.PostArchive()
	if err != nil {
		return ctx.Redirect(302, "/")
	}
	return ctx.Render(http.StatusOK, "archive.html", map[string]interface{}{
		"Archives": archives,
	})
}

// ------------------------------------------------------ 分类页面 ------------------------------------------------------
// ViewCatePost 分类文章列表
func ViewCatePost(ctx echo.Context) error {
	cate := ctx.Param("cate")
	if cate == "" {
		return ctx.Redirect(302, "/")
	}
	mod, has := model.CateGet(cate)
	if !has {
		return ctx.Redirect(302, "/")
	}
	pi, _ := strconv.Atoi(ctx.FormValue("page"))
	if pi == 0 {
		pi = 1
	}
	ps := model.Gcfg().PageSize
	mods, err := model.PostPage(mod.Id, model.PostKindPost, pi, ps, "id", "title", "path", "created", "summary", "updated", "status")
	if err != nil || len(mods) < 1 {
		return ctx.Redirect(302, "/")
	}
	total := model.PostCount(mod.Id, model.PostKindPost)
	naver := model.Naver{}
	if pi > 1 {
		naver.Prev = "/cate/" + mod.Name + "?page=1"
	}
	if total > (pi * ps) {
		naver.Next = "/cate/" + mod.Name + "?page=" + strconv.Itoa(pi+1)
	}
	return ctx.Render(http.StatusOK, "post-cate.html", map[string]interface{}{
		"Cate":      mod,
		"CatePosts": mods,
		"Naver":     naver,
	})
}

// ------------------------------------------------------ 标签页面 ------------------------------------------------------
// ViewTags 标签页面
func ViewTags(ctx echo.Context) error {
	mods, err := model.TagStateAll()
	if err != nil {
		return ctx.Redirect(302, "/")
	}
	return ctx.Render(http.StatusOK, "tags.html", map[string]interface{}{
		"Tags": mods,
	})
}

// ViewTagPost 标签下的文章列表
func ViewTagPost(ctx echo.Context) error {
	tag := ctx.Param("tag")
	if tag == "" {
		return ctx.Redirect(302, "/tags")
	}
	mod, has := model.TagGet(tag)
	if !has {
		return ctx.Redirect(302, "/tags")
	}
	pi, _ := strconv.Atoi(ctx.FormValue("page"))
	if pi == 0 {
		pi = 1
	}
	ps := model.Gcfg().PageSize
	mods, err := model.TagPostPage(mod.Id, pi, ps)
	if err != nil || len(mods) < 1 {
		return ctx.Redirect(302, "/tags")
	}
	total := model.TagPostCount(mod.Id)
	naver := model.Naver{}
	if pi > 1 {
		naver.Prev = "/tag/" + mod.Name + "?page=1"
	}
	if total > (pi * ps) {
		naver.Next = "/tag/" + mod.Name + "?page=" + strconv.Itoa(pi+1)
	}
	return ctx.Render(http.StatusOK, "post-tag.html", map[string]interface{}{
		"Tag":      mod,
		"TagPosts": mods,
		"Naver":    naver,
	})
}
