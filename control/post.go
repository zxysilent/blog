package control

import (
	"blog/model"
	"blog/util"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo"
)

var reg = regexp.MustCompile(`<img src="([^" ]+)" alt="([^" ]*)"\s?\/?>`)

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
		return ctx.Res(util.NewSucc("", mod))
	}
	return nil
}

// PostGet 一个
// id int
func PostGet(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.Res(util.NewErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	mod, has := model.PostGet(id)
	if !has {
		return ctx.Res(util.NewErrOpt(`未查询信息`))
	}
	return ctx.Res(util.NewSucc(`信息`, mod))
}

// PostPageAll 页面列表
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

// PostTagIds 通过文章id 获取 标签ids
func PostTagIds(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.Res(util.NewErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	mods := model.PostTagIds(id)
	if mods == nil {
		return ctx.Res(util.NewErrOpt(`未查询到标签信息`))
	}
	return ctx.Res(util.NewSucc(`标签ids`, mods))
}

// PostOpts 文章操作
func PostOpts(ctx echo.Context) error {
	ipt := &struct {
		Post model.Post `json:"post" form:"post"` // 文章信息
		Type int        `json:"type" form:"type"` // 0 文章 1 页面
		Tags []int      `json:"tags" form:"tags"` // 标签
		Edit bool       `json:"edit" form:"edit"` // 是否编辑
	}{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.Res(util.NewErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	if !ipt.Edit && model.PostExist(ipt.Post.Path) {
		return ctx.Res(util.NewErrIpt(`当前访问路径已经存在,请重新输入`))
	}
	// 同步类型
	ipt.Post.Type = ipt.Type
	if strings.Contains(ipt.Post.Content, "<!--more-->") {
		ipt.Post.Summary = strings.Split(ipt.Post.Content, "<!--more-->")[0]
	}
	// 编辑 文章/页面
	if ipt.Edit {
		// 修改日期在发布日期之前
		if ipt.Post.UpdateTime.Before(ipt.Post.CreateTime) {
			// 修改时间再发布时间后1分钟
			ipt.Post.UpdateTime = ipt.Post.CreateTime.Add(time.Minute * 2)
		}
		if model.PostEdit(&ipt.Post) {
			if ipt.Type == 0 {
				// 处理变动标签

				return ctx.Res(util.NewSucc(`文章修改成功`))
			}
			return ctx.Res(util.NewSucc(`页面修改成功`))
		}
		if ipt.Type == 0 {
			return ctx.Res(util.NewFail(`文章修改失败,请重试`))
		}
		return ctx.Res(util.NewFail(`页面修改失败,请重试`))
	}
	// 添加 文章/页面
	ipt.Post.UpdateTime = time.Now()
	if model.PostAdd(&ipt.Post) {
		// 添加标签
		// 文章
		if ipt.Type == 0 {
			//添加标签
			tagPosts := make([]model.PostTag, 0, len(ipt.Tags))
			for _, itm := range ipt.Tags {
				tagPosts = append(tagPosts, model.PostTag{
					TagId:  itm,
					PostId: ipt.Post.Id,
				})
			}
			model.TagPostAdds(&tagPosts)
			return ctx.Res(util.NewSucc(`文章添加成功`))
		}
		return ctx.Res(util.NewSucc(`页面添加成功`))
	}
	if ipt.Type == 0 {
		return ctx.Res(util.NewFail(`文章添加失败,请重试`))
	}
	return ctx.Res(util.NewFail(`页面添加失败,请重试`))
}
