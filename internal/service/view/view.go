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
	"github.com/zxysilent/logs"
)

// publicFilter 构造前台列表过滤器
// 前台只展示「已发布 + 可见」的内容，草稿与隐藏内容一律不出现
// 分页参数缺省值统一由 model.List.Normalize 处理
func publicFilter(ctx echo.Context, ps int) *model.PostFilterPage {
	filter := &model.PostFilterPage{}
	filter.Pi, _ = strconv.Atoi(ctx.FormValue("page"))
	filter.Normalize(ps)
	status := model.StatusFinish
	show := true
	filter.Status = &status
	filter.Show = &show
	return filter
}

// ------------------------------------------------------ 主页面 ------------------------------------------------------
// ViewIndex 主页面
func ViewIndex(ctx echo.Context) error {
	filter := publicFilter(ctx, repo.GetCfg().Blog.PageSize)
	kind := model.KindPost
	filter.Kind = &kind
	mods, total, _ := repo.PostPage(filter, "id", "title", "path", "created", "summary", "author", "icon", "hits", "istop")
	naver := model.Naver{}
	if filter.Pi > 1 {
		naver.Prev = "/?page=" + strconv.Itoa(filter.Pi-1)
	}
	if total > (filter.Pi * filter.Ps) {
		naver.Next = "/?page=" + strconv.Itoa(filter.Pi+1)
	}
	return ctx.Render(http.StatusOK, "index.html", map[string]any{
		"Posts": mods,
		"Naver": naver,
	})
}

// ViewSearch 搜索页面
func ViewSearch(ctx echo.Context) error {
	rctx := ctx.Request().Context()
	mult := ctx.FormValue("mult")
	count := 0
	tags := []model.Tag{}
	naver := model.Naver{}
	mods := []model.PostPart{}
	if mult != "" {
		filter := publicFilter(ctx, repo.GetCfg().Blog.PageSize)
		kind := model.KindPost
		filter.Kind = &kind
		filter.Mult = mult
		mods, count, _ = repo.PostPage(filter, "id", "title", "path", "created", "summary")
		if filter.Pi > 1 {
			naver.Prev = "/search?mult=" + mult + "&page=" + strconv.Itoa(filter.Pi-1)
		}
		if count > (filter.Pi * filter.Ps) {
			naver.Next = "/search?mult=" + mult + "&page=" + strconv.Itoa(filter.Pi+1)
		}
	} else {
		tags, _, _ = repo.TagList(rctx, &model.TagListArgs{})
	}
	return ctx.Render(http.StatusOK, "search.html", map[string]any{
		"Posts": mods,
		"Mult":  mult,
		"Count": count,
		"Naver": naver,
		"Tags":  tags,
	})
}

// ------------------------------------------------------ 文章页面 ------------------------------------------------------
// ViewPost 文章页面
// 按公开路径读取已发布且可见的内容，读取成功即累加浏览量
func ViewPost(ctx echo.Context) error {
	rctx := ctx.Request().Context()
	paths := strings.Split(ctx.Param("*"), ".")
	if len(paths) >= 1 && paths[0] != "" {
		mod, err := repo.PostGetPath(rctx, model.PostNormalizePath(paths[0]))
		if err != nil {
			return ctx.Render(http.StatusOK, "error.html", map[string]any{
				"Msg": "内容不存在",
			})
		}
		if mod.Kind != model.KindPost && mod.Kind != model.KindPage {
			return ctx.Redirect(http.StatusFound, "/")
		}
		naver := repo.PostNaver(mod.CateId, mod.Created)
		mod.Richtext = utils.GenLazyLoad(mod.Richtext)
		// 浏览量自增，失败不影响页面渲染
		_ = repo.PostHits(rctx, mod.Id)
		mod.Hits++
		return ctx.Render(http.StatusOK, "post.html", map[string]any{
			"Post":  mod,
			"Show":  mod.Status == model.StatusFinish,
			"Naver": naver,
		})
	}
	return ctx.Redirect(302, "/404")
}

// ViewPost 文章页面
func ViewNote(ctx echo.Context) error {
	// share := token.Token[token.Share]{}
	// err := share.Decode(ctx.QueryParam("s"), repo.GetCfg().Auth.TokenSecret)
	// if err != nil {
	// 	return ctx.Render(http.StatusOK, "error.html", map[string]any{
	// 		"Msg": "分享失败",
	// 	})
	// }
	// paths := strings.Split(ctx.Param("*"), ".")
	// if len(paths) >= 1 {
	// 	mod := &model.Post{Path: paths[0], Kind: model.KindNote}
	// 	err := repo.PostGet(mod)
	// 	if err != nil {
	// 		return ctx.Render(http.StatusOK, "error.html", map[string]any{
	// 			"Msg": err.Error(),
	// 		})
	// 	}
	// 	mod.Richtext = utils.GenLazyLoad(mod.Richtext)
	// 	return ctx.Render(http.StatusOK, "note.html", map[string]any{
	// 		"Note":  mod,
	// 		"ExpAt": share.ExpAt * 1000,
	// 		"Show":  share.Claim.PostId == mod.Id,
	// 	})
	// }
	return ctx.Redirect(302, "/404")
}

// // ------------------------------------------------------ 单个页面 ------------------------------------------------------
// ViewAbout 关于页面
func ViewAbout(ctx echo.Context) error {
	// mod := &model.Post{Path: "about", Kind: model.KindPage}
	// err := repo.PostGet(mod)
	// if err != nil {
	// 	return ctx.Redirect(302, "/")
	// }
	// mod.Richtext = utils.GenLazyLoad(mod.Richtext)
	// return ctx.Render(http.StatusOK, "page.html", map[string]any{
	// 	"Page": mod,
	// 	"Show": mod.Status == 2,
	// })
	return nil
}

// ViewLinks 友链页面
func ViewLinks(ctx echo.Context) error {
	// mod := &model.Post{Path: "links", Kind: model.KindPage}
	// err := repo.PostGet(mod)
	// if err != nil {
	// 	return ctx.Redirect(302, "/")
	// }
	// mod.Richtext = utils.GenLazyLoad(mod.Richtext)
	// return ctx.Render(http.StatusOK, "page.html", map[string]any{
	// 	"Page": mod,
	// 	"Show": mod.Status == 2,
	// })
	return nil
}

// ViewPage 其它页面
func ViewPage(ctx echo.Context) error {
	// paths := strings.Split(ctx.Param("*"), ".")
	// if len(paths) >= 1 {
	// 	mod := &model.Post{Path: paths[0], Kind: model.KindPage}
	// 	err := repo.PostGet(mod)
	// 	if err != nil {
	// 		return ctx.Render(http.StatusOK, "error.html", map[string]any{
	// 			"Msg": err.Error(),
	// 		})
	// 	}
	// 	mod.Richtext = utils.GenLazyLoad(mod.Richtext)
	// 	return ctx.Render(http.StatusOK, "page.html", map[string]any{
	// 		"Page": mod,
	// 		"Show": mod.Status == 2,
	// 	})
	// }
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
	// filter := &model.PostFilterList{}
	// post := model.KindPost
	// filter.Kind = &post
	// status := model.StatusFinish
	// filter.Status = &status
	// mods, err := repo.PostList(filter, "id", "title", "path", "created")
	// if err != nil {
	// 	return ctx.Redirect(302, "/")
	// }
	// archives := make([]model.PostArchive, 0, 8)
	// for _, v := range mods {
	// 	created := time.UnixMilli(v.Created)
	// 	if idx := archInOf(created, archives); idx == -1 { //没有
	// 		archives = append(archives, model.PostArchive{
	// 			Time:  created,
	// 			Posts: []model.PostPart{v},
	// 		})
	// 	} else { //有
	// 		archives[idx].Posts = append(archives[idx].Posts, v)
	// 	}
	// }
	// return ctx.Render(http.StatusOK, "archive.html", map[string]any{
	// 	"Archives": archives,
	// })
	return nil
}

// ------------------------------------------------------ 分类页面 ------------------------------------------------------
// ViewCatePost 栏目页面
// 按公开路径（Path）访问，而不是名称，改名不影响已有链接：
//   - Url 非空：直接 302 跳转，不渲染
//   - Auth 为真：未登录时跳到首页
//   - 有子栏目且配置了 Itpl：用栏目首页模板（通常用于展示子栏目入口）
//   - 其余情况：用 Ltpl（默认 post-cate.html）
func ViewCatePost(ctx echo.Context) error {
	raw := strings.TrimSuffix(ctx.Param("cate"), ".html")
	if raw == "" {
		return ctx.Redirect(http.StatusFound, "/")
	}
	mod, ok := repo.CateByPath(model.CateNormalizePath(raw))
	if !ok {
		// 兼容旧链接：早期前台用 /cates/<名称>，仍可命中
		mod, ok = repo.CateByName(raw)
	}
	if !ok {
		return ctx.Redirect(http.StatusFound, "/")
	}
	if mod.Url != "" {
		return ctx.Redirect(http.StatusFound, mod.Url)
	}
	if mod.Auth && !viewAuthed(ctx) {
		return ctx.Redirect(http.StatusFound, "/")
	}
	children := repo.CateChildren(mod.Id)
	filter := publicFilter(ctx, repo.GetCfg().Blog.PageSize)
	// 父栏目聚合子栏目的内容，子栏目只查自己
	filter.CateIds = repo.CateSubIds(mod.Id)
	kind := mod.Kind
	filter.Kind = &kind
	mods, total, err := repo.PostPage(filter, "id", "title", "path", "created", "summary", "updated", "status", "author", "icon", "hits", "istop")
	if err != nil {
		logs.Ctx(ctx.Request().Context()).Err(err).Error("repo.PostPage")
	}
	naver := model.Naver{}
	base := "/cates/" + mod.Path
	if filter.Pi > 1 {
		naver.Prev = base + "?page=" + strconv.Itoa(filter.Pi-1)
	}
	if total > (filter.Pi * filter.Ps) {
		naver.Next = base + "?page=" + strconv.Itoa(filter.Pi+1)
	}
	name := mod.Ltpl
	if len(children) > 0 && mod.Itpl != "" {
		name = mod.Itpl
	}
	if name == "" {
		name = "post-cate.html"
	}
	return ctx.Render(http.StatusOK, name, map[string]any{
		"Cate":      mod,
		"SubCates":  children,
		"CatePosts": mods,
		"Total":     total,
		"Naver":     naver,
	})
}

// viewAuthed 前台登录态判定：能解出合法令牌即视为已登录
// 前台没有权限点概念，只区分匿名与登录
func viewAuthed(ctx echo.Context) bool {
	raw := ctx.Request().Header.Get(echo.HeaderAuthorization)
	if raw == "" {
		raw = ctx.FormValue(conf.TokenKey)
	}
	if raw == "" {
		if cookie, err := ctx.Cookie(conf.TokenKey); err == nil && cookie != nil {
			raw = cookie.Value
		}
	}
	if raw == "" {
		return false
	}
	tk := token.Token[token.Auth]{}
	return tk.Decode(raw, repo.GetCfg().Auth.TokenSecret) == nil
}

// ------------------------------------------------------ 标签页面 ------------------------------------------------------
// ViewTags 标签页面
func ViewTags(ctx echo.Context) error {
	mods, err := repo.TagStatus()
	if err != nil {
		return ctx.Render(http.StatusOK, "error.html", map[string]any{
			"Msg": err.Error(),
		})
	}
	return ctx.Render(http.StatusOK, "tags.html", map[string]any{
		"Tags": mods,
	})
}

// ViewTagPost 标签下的文章列表
func ViewTagPost(ctx echo.Context) error {
	// tagName := ctx.Param("tag")
	// if tagName == "" {
	// 	return ctx.Redirect(302, "/tags")
	// }
	// tag := &model.Tag{
	// 	Name: tagName,
	// }
	// err := repo.TagGet(tag)
	// if err != nil {
	// 	return ctx.Redirect(302, "/tags")
	// }
	// filter := &model.PostTagFilter{}
	// filter.TagId = tag.Id
	// filter.Pi, _ = strconv.Atoi(ctx.FormValue("page"))
	// if filter.Pi == 0 {
	// 	filter.Pi = 1
	// }
	// filter.Ps = repo.GetCfg().Blog.PageSize
	// mods, err := repo.TagPostPage(filter)
	// if err != nil || len(mods) < 1 {
	// 	return ctx.Redirect(302, "/tags")
	// }
	// total := repo.TagPostCount(tag.Id)
	// naver := model.Naver{}
	// if filter.Pi > 1 {
	// 	naver.Prev = "/tags/" + tag.Name + "?page=" + strconv.Itoa(filter.Pi-1)
	// }
	// if total > (filter.Pi * filter.Ps) {
	// 	naver.Next = "/tags/" + tag.Name + "?page=" + strconv.Itoa(filter.Pi+1)
	// }
	// return ctx.Render(http.StatusOK, "post-tag.html", map[string]any{
	// 	"Tag":      tag,
	// 	"TagPosts": mods,
	// 	"Naver":    naver,
	// })
	return nil
}
