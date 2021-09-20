package appctl

import (
	"blog/model"
	"net/http"
	"regexp"
	"strings"

	"github.com/labstack/echo/v4"
)

// lazy-load 懒加载图片
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

func similar(a, b string) int {
	if a[:4] == b[:4] {
		return 0
	}
	if a[:4] < b[:4] {
		return 1
	}
	return -1
}

func inOf(goal int, arr []model.Tag) bool {
	for idx := range arr {
		if goal == arr[idx].Id {
			return true
		}
	}
	return false
}
