package appctl

import (
	"blog/model"
	"net/http"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// IndexView 主页面
func IndexView(ctx echo.Context) error {
	//return ctx.HTML(200, `<html><head><meta charset="UTF-8"><title>文档</title></head><body><a href="/swagger/index.html">doc</a></body></html>`)
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

// PostView 文章页面
func ViewPost(ctx echo.Context) error {
	//return ctx.HTML(200, `<html><head><meta charset="UTF-8"><title>文档</title></head><body><a href="/swagger/index.html">doc</a></body></html>`)
	paths := strings.Split(ctx.Param("*"), ".")
	if len(paths) == 2 {
		mod, naver, has := model.PostPath(paths[0])
		if !has {
			return ctx.Redirect(302, "/")
		}
		if paths[1] == "html" {
			mod.Richtext = regImg.ReplaceAllString(mod.Richtext, `<img class="lazy-load" src="data:image/gif;base64,R0lGODlhAQABAAAAACH5BAEKAAEALAAAAAABAAEAAAICTAEAOw==" data-src="$1" alt="$2">`)
			tags, _ := model.PostTags(mod.Id)
			return ctx.Render(http.StatusOK, "post.html", map[string]interface{}{
				"Post":    mod,
				"Naver":   naver,
				"Tags":    tags,
				"HasTag":  len(tags) > 0,
				"HasCate": mod.Cate != nil,
			})
		}
		return ctx.JSON(utils.Succ("", mod))
	}
	return ctx.Redirect(302, "/404")
}

// getTocHTML 生成目录并替换内容
func getTocHTML(html string) string {
	html = strings.Replace(html, `id="`, `id="toc_`, -1)
	hs := regToc.FindAllString(html, -1)
	if len(hs) > 1 {
		sb := strings.Builder{}
		sb.WriteString(`<div class="toc"><ul>`)
		level := 0
		for i := 0; i < len(hs)-1; i++ {
			fg := similar(hs[i], hs[i+1])
			if fg == 0 {
				sb.WriteString(regH.ReplaceAllString(hs[i], `<li><a href="#$1">$2</a></li>`))
			} else if fg == 1 {
				level++
				sb.WriteString(regH.ReplaceAllString(hs[i], `<li><a href="#$1">$2</a><ul>`))
			} else {
				level--
				sb.WriteString(regH.ReplaceAllString(hs[i], `<li><a href="#$1">$2</a></li></ul></li>`))
			}
		}
		fg := similar(hs[len(hs)-2], hs[len(hs)-1])
		if fg == 0 {
			sb.WriteString(regH.ReplaceAllString(hs[len(hs)-1], `<li><a href="#$1">$2</a></li>`))
		} else if fg == 1 {
			level++
			sb.WriteString(regH.ReplaceAllString(hs[len(hs)-1], `<li><a href="#$1">$2</a><ul>`))
		} else {
			level--
			sb.WriteString(regH.ReplaceAllString(hs[len(hs)-1], `<li><a href="#$1">$2</a></li></ul></li>`))
		}
		for level > 0 {
			sb.WriteString(`</ul></li>`)
			level--
		}
		sb.WriteString(`</ul></div>`)
		return sb.String() + html
	}
	if len(hs) == 1 {
		sb := strings.Builder{}
		sb.WriteString(`<div class="toc"><ul>`)
		sb.WriteString(regH.ReplaceAllString(hs[0], `<li><a href="#$1">$2</a></li>`))
		sb.WriteString(`</ul></div>`)
		return sb.String() + html
	}
	return html
}

func atoi(raw string, def int) (int, error) {
	out, err := strconv.Atoi(raw)
	if err != nil {
		return def, err
	}
	return out, nil
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
	return RenderPage(ctx.Param("*"), ctx)
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
