package view

import (
	"blog/conf"
	"blog/internal/model"
	"blog/internal/repo"
	"blog/internal/utils"
	"blog/pkg/token"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

// ------------------------------------------------------ 主页面 ------------------------------------------------------
// ViewIndex 主页面
func ViewIndex(ctx echo.Context) error {
	filter := &model.PostFilterPage{}
	filter.Pi, _ = strconv.Atoi(ctx.FormValue("page"))
	if filter.Pi == 0 {
		filter.Pi = 1
	}
	filter.Ps = repo.CfgBlog().PageSize
	post := model.KindPost
	filter.Kind = &post
	mods, total, _ := repo.PostPage(filter, "id", "title", "path", "created", "summary")
	naver := model.Naver{}
	if filter.Pi > 1 {
		naver.Prev = "/?page=" + strconv.Itoa(filter.Pi-1)
	}
	if total > (filter.Pi * filter.Ps) {
		naver.Next = "/?page=" + strconv.Itoa(filter.Pi+1)
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
	if len(paths) >= 1 {
		mod := &model.Post{Path: paths[0], Kind: model.KindPost}
		err := repo.PostGet(mod)
		if err != nil {
			return ctx.Render(http.StatusOK, "error.html", map[string]interface{}{
				"Msg": err.Error(),
			})
		}
		naver := model.Naver{}
		mod.Richtext = utils.GenLazyLoad(mod.Richtext)
		return ctx.Render(http.StatusOK, "post.html", map[string]interface{}{
			"Post":  mod,
			"Show":  mod.Status == 2,
			"Naver": naver,
		})
	}
	return ctx.Redirect(302, "/404")
}

// ViewPost 文章页面
func ViewNote(ctx echo.Context) error {
	share := token.Token[token.Share]{}
	err := share.Decode(ctx.QueryParam("s"), conf.App.TokenSecret)
	if err != nil {
		return ctx.Render(http.StatusOK, "error.html", map[string]interface{}{
			"Msg": "分享失败",
		})
	}
	paths := strings.Split(ctx.Param("*"), ".")
	if len(paths) >= 1 {
		mod := &model.Post{Path: paths[0], Kind: model.KindNote}
		err := repo.PostGet(mod)
		if err != nil {
			return ctx.Render(http.StatusOK, "error.html", map[string]interface{}{
				"Msg": err.Error(),
			})
		}
		mod.Richtext = utils.GenLazyLoad(mod.Richtext)
		return ctx.Render(http.StatusOK, "note.html", map[string]interface{}{
			"Note":  mod,
			"ExpAt": share.ExpAt * 1000,
			"Show":  share.Claim.PostId == mod.Id,
		})
	}
	return ctx.Redirect(302, "/404")
}

// // ------------------------------------------------------ 单个页面 ------------------------------------------------------
// ViewAbout 关于页面
func ViewAbout(ctx echo.Context) error {
	mod := &model.Post{Path: "about", Kind: model.KindPage}
	err := repo.PostGet(mod)
	if err != nil {
		return ctx.Redirect(302, "/")
	}
	mod.Richtext = utils.GenLazyLoad(mod.Richtext)
	return ctx.Render(http.StatusOK, "page.html", map[string]interface{}{
		"Page": mod,
		"Show": mod.Status == 2,
	})
}

// ViewLinks 友链页面
func ViewLinks(ctx echo.Context) error {
	mod := &model.Post{Path: "links", Kind: model.KindPage}
	err := repo.PostGet(mod)
	if err != nil {
		return ctx.Redirect(302, "/")
	}
	mod.Richtext = utils.GenLazyLoad(mod.Richtext)
	return ctx.Render(http.StatusOK, "page.html", map[string]interface{}{
		"Page": mod,
		"Show": mod.Status == 2,
	})
}

// ViewPage 其它页面
func ViewPage(ctx echo.Context) error {
	paths := strings.Split(ctx.Param("*"), ".")
	if len(paths) >= 1 {
		mod := &model.Post{Path: paths[0], Kind: model.KindPage}
		err := repo.PostGet(mod)
		if err != nil {
			return ctx.Render(http.StatusOK, "error.html", map[string]interface{}{
				"Msg": err.Error(),
			})
		}
		mod.Richtext = utils.GenLazyLoad(mod.Richtext)
		return ctx.Render(http.StatusOK, "page.html", map[string]interface{}{
			"Page": mod,
			"Show": mod.Status == 2,
		})
	}
	return ctx.Redirect(302, "/404")

}
func archInOf(time time.Time, mods []model.PostArchive) int {
	for idx := 0; idx < len(mods); idx++ {
		if time.Year() == mods[idx].Time.Year() && time.Month() == mods[idx].Time.Month() {
			return idx
		}
	}
	return -1
}

// ------------------------------------------------------ 归档页面 ------------------------------------------------------
// ViewArchives 归档页面
func ViewArchives(ctx echo.Context) error {
	filter := &model.PostFilterList{}
	post := model.KindPost
	filter.Kind = &post
	status := model.PostStatusFinish
	filter.Status = &status
	mods, err := repo.PostList(filter, "id", "title", "path", "created")
	if err != nil {
		return ctx.Redirect(302, "/")
	}
	archives := make([]model.PostArchive, 0, 8)
	for _, v := range mods {
		created := time.UnixMilli(v.Created)
		if idx := archInOf(created, archives); idx == -1 { //没有
			archives = append(archives, model.PostArchive{
				Time:  created,
				Posts: []model.PostPart{v},
			})
		} else { //有
			archives[idx].Posts = append(archives[idx].Posts, v)
		}
	}
	return ctx.Render(http.StatusOK, "archive.html", map[string]interface{}{
		"Archives": archives,
	})
}

// ------------------------------------------------------ 分类页面 ------------------------------------------------------
// ViewCatePost 分类文章列表
func ViewCatePost(ctx echo.Context) error {
	cateName := ctx.Param("cate")
	if cateName == "" {
		return ctx.Redirect(302, "/")
	}
	mod := &model.Cate{
		Name: cateName,
	}
	err := repo.CateGet(mod)
	if err != nil {
		return ctx.Redirect(302, "/")
	}
	filter := &model.PostFilterPage{}
	filter.Pi, _ = strconv.Atoi(ctx.FormValue("page"))
	if filter.Pi == 0 {
		filter.Pi = 1
	}
	filter.Ps = repo.CfgBlog().PageSize
	post := model.KindPost
	filter.CateId = &mod.Id
	filter.Kind = &post
	mods, total, err := repo.PostPage(filter, "id", "title", "path", "created", "summary", "updated", "status")
	if err != nil || len(mods) < 1 {
		return ctx.Redirect(302, "/")
	}
	naver := model.Naver{}
	if filter.Pi > 1 {
		naver.Prev = "/cate/" + mod.Name + "?page=" + strconv.Itoa(filter.Pi-1)
	}
	if total > (filter.Pi * filter.Ps) {
		naver.Next = "/cate/" + mod.Name + "?page=" + strconv.Itoa(filter.Pi+1)
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
	mods, err := repo.TagStatus()
	if err != nil {
		return ctx.Render(http.StatusOK, "error.html", map[string]interface{}{
			"Msg": err.Error(),
		})
	}
	return ctx.Render(http.StatusOK, "tags.html", map[string]interface{}{
		"Tags": mods,
	})
}

// ViewTagPost 标签下的文章列表
func ViewTagPost(ctx echo.Context) error {
	tagName := ctx.Param("tag")
	if tagName == "" {
		return ctx.Redirect(302, "/tags")
	}
	tag := &model.Tag{
		Name: tagName,
	}
	err := repo.TagGet(tag)
	if err != nil {
		return ctx.Redirect(302, "/tags")
	}
	filter := &model.PostTagFilter{}
	filter.TagId = tag.Id
	filter.Pi, _ = strconv.Atoi(ctx.FormValue("page"))
	if filter.Pi == 0 {
		filter.Pi = 1
	}
	filter.Ps = repo.CfgBlog().PageSize
	mods, err := repo.TagPostPage(filter)
	if err != nil || len(mods) < 1 {
		return ctx.Redirect(302, "/tags")
	}
	total := repo.TagPostCount(tag.Id)
	naver := model.Naver{}
	if filter.Pi > 1 {
		naver.Prev = "/tag/" + tag.Name + "?page=" + strconv.Itoa(filter.Pi-1)
	}
	if total > (filter.Pi * filter.Ps) {
		naver.Next = "/tag/" + tag.Name + "?page=" + strconv.Itoa(filter.Pi+1)
	}
	return ctx.Render(http.StatusOK, "post-tag.html", map[string]interface{}{
		"Tag":      tag,
		"TagPosts": mods,
		"Naver":    naver,
	})
}
