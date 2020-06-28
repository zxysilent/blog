package control

import (
	"blog/model"
	"net/http"
	"regexp"
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
	ps, _ := atoi(model.MapOpts.MustGet("page_size"), 6)
	mods, _ := model.PostPage(pi, ps)
	total := model.PostCount()
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

// ArchivesView 归档
func ArchivesView(ctx echo.Context) error {
	mods, err := model.PostArchive()
	if err != nil {
		return ctx.Redirect(302, "/")
	}
	return ctx.Render(http.StatusOK, "archive.html", map[string]interface{}{
		"Posts": mods,
	})
}
func ArchivesJson(ctx echo.Context) error {
	mods, err := model.PostArchive()
	if err != nil {
		return ctx.JSON(utils.Fail("未查询到数据", err))
	}
	return ctx.JSON(utils.Succ("归档", mods))
}

// CatePostView 分类文章列表
func CatePostView(ctx echo.Context) error {
	cate := ctx.Param("cate")
	if cate == "" {
		return ctx.Redirect(302, "/")
	}
	mod, has := model.CateName(cate)
	if !has {
		return ctx.Redirect(302, "/")
	}
	pi, _ := strconv.Atoi(ctx.FormValue("page"))
	if pi == 0 {
		pi = 1
	}
	ps, _ := atoi(model.MapOpts.MustGet("page_size"), 6)
	mods, err := model.CatePostList(mod.Id, pi, ps, true)
	if err != nil || len(mods) < 1 {
		return ctx.Redirect(302, "/")
	}
	total := model.CatePostCount(mod.Id, true)
	naver := model.Naver{}
	if pi > 1 {
		naver.Prev = "/cate/" + mod.Name + "?page=1"
	}
	if total > (pi * ps) {
		naver.Next = "/cate/" + mod.Name + "?page=" + strconv.Itoa(pi+1)
	}
	return ctx.Render(http.StatusOK, "cate-post.html", map[string]interface{}{
		"Cate":      mod,
		"CatePosts": mods,
		"Naver":     naver,
	})
}

// TagsView 标签列表
func TagsView(ctx echo.Context) error {
	mods, err := model.TagStateAll()
	if err != nil {
		return ctx.Redirect(302, "/")
	}
	return ctx.Render(http.StatusOK, "tags.html", map[string]interface{}{
		"Tags": mods,
	})
}
func TagsJson(ctx echo.Context) error {
	mods, err := model.TagStateAll()
	if err != nil {
		return ctx.JSON(utils.Fail("未查询到标签", err))
	}
	return ctx.JSON(utils.Succ("标签", mods))
}

// TagPostView 标签文章列表
func TagPostView(ctx echo.Context) error {
	tag := ctx.Param("tag")
	if tag == "" {
		return ctx.Redirect(302, "/tags")
	}
	mod, has := model.TagName(tag)
	if !has {
		return ctx.Redirect(302, "/tags")
	}
	pi, _ := strconv.Atoi(ctx.FormValue("page"))
	if pi == 0 {
		pi = 1
	}
	ps, _ := atoi(model.MapOpts.MustGet("page_size"), 6)
	mods, err := model.TagPostList(mod.Id, pi, ps)
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
	return ctx.Render(http.StatusOK, "tag-post.html", map[string]interface{}{
		"Tag":      mod,
		"TagPosts": mods,
		"Naver":    naver,
	})
}

// PostView 文章页面
func PostView(ctx echo.Context) error {
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

var reg = regexp.MustCompile(`<img src="([^" ]+)" alt="([^" ]*)"\s?\/?>`)

// 生成目录并替换内容
func getTocHTML(html string) string {
	html = strings.Replace(html, `id="`, `id="toc_`, -1)
	regToc := regexp.MustCompile("<h[1-6]>.*?</h[1-6]>")
	regH := regexp.MustCompile(`<h[1-6]><a id="(.*?)"></a>(.*?)</h[1-6]>`)
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
