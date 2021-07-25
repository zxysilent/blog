package appctl

import (
	"blog/model"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// PostGet doc
// @Tags post-文章页面
// @Summary 通过id获取单条文章
// @Param id query int true "id"
// @Success 200 {object} model.Reply{data=model.Post} "返回数据"
// @Router /api/post/get [get]
func PostGet(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod, has := model.PostGet(ipt.Id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到文章"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// PostPage doc
// @Tags post-文章页面
// @Summary 获取文章分页
// @Param cate_id path int true "分类id" default(1)
// @Param pi query int true "分页数" default(1)
// @Param ps query int true "每页条数[5,30]" default(5)
// @Success 200 {object} model.Reply{data=[]model.Post} "返回数据"
// @Router /api/post/page [get]
func PostPage(ctx echo.Context) error {
	cateId, err := strconv.Atoi(ctx.QueryParam("cate_id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误", err.Error()))
	}
	ipt := &model.Page{}
	err = ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if ipt.Ps > 30 || ipt.Ps < 1 {
		return ctx.JSON(utils.ErrIpt("分页大小输入错误", ipt.Ps))
	}
	count := model.PostCount(cateId, model.PostKindPost)
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := model.PostPage(cateId, model.PostKindPost, ipt.Pi, ipt.Ps, "id", "title", "path", "created", "summary", "updated", "status")
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, int(count)))
}

// PostTagGet 通过文章id 获取 标签ids
func PostTagGet(ctx echo.Context) error {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		return ctx.JSON(utils.ErrIpt("数据输入错误,请重试", err.Error()))
	}
	mods := model.PostTagGet(id)
	if mods == nil {
		return ctx.JSON(utils.ErrOpt("未查询到标签信息"))
	}
	return ctx.JSON(utils.Succ("标签ids", mods))
}

// PostOpts 文章操作
func PostOpts(ctx echo.Context) error {
	// ipt := &struct {
	// 	Post model.Post "json:"post" form:"post"" // 文章信息
	// 	Type int        "json:"type" form:"type"" // 0 文章 1 页面
	// 	Tags []int      "json:"tags" form:"tags"" // 标签
	// 	Edit bool       "json:"edit" form:"edit"" // 是否编辑
	// }{}
	// err := ctx.Bind(ipt)
	// if err != nil {
	// 	return ctx.JSON(utils.ErrIpt("数据输入错误,请重试", err.Error()))
	// }
	// if !ipt.Edit && model.PostExist(ipt.Post.Path) {
	// 	return ctx.JSON(utils.ErrIpt("当前访问路径已经存在,请重新输入"))
	// }
	// // 同步类型
	// ipt.Post.Kind = ipt.Type
	// if strings.Contains(ipt.Post.Richtext, "<!--more-->") {
	// 	ipt.Post.Summary = strings.Split(ipt.Post.Richtext, "<!--more-->")[0]
	// }
	// // 生成目录
	// if ipt.Type == 0 {
	// 	ipt.Post.Richtext = getTocHTML(ipt.Post.Richtext)
	// }
	// // 编辑 文章/页面
	// if ipt.Edit {
	// 	// 修改日期在发布日期之前
	// 	if ipt.Post.Updated.Before(ipt.Post.Created) {
	// 		// 修改时间再发布时间后1分钟
	// 		ipt.Post.Updated = ipt.Post.Created.Add(time.Minute * 2)
	// 	}
	// 	if model.PostEdit(&ipt.Post, "cate_id", "status", "title", "summary", "markdown", "richtext", "allow", "created", "updated") {
	// 		if ipt.Type == 0 {
	// 			// 处理变动标签
	// 			old := model.PostTagGet(ipt.Post.Id)
	// 			new := ipt.Tags
	// 			add := make([]int, 0)
	// 			del := make([]int, 0)
	// 			for _, itm := range old {
	// 				if !inOf(itm, new) {
	// 					del = append(del, itm)
	// 				}
	// 			}
	// 			for _, itm := range new {
	// 				if !inOf(itm, old) {
	// 					add = append(add, itm)
	// 				}
	// 			}
	// 			tagAdds := make([]model.PostTag, 0, len(add))
	// 			for _, itm := range add {
	// 				tagAdds = append(tagAdds, model.PostTag{
	// 					TagId:  itm,
	// 					PostId: ipt.Post.Id,
	// 				})
	// 			}
	// 			// 删除标签
	// 			model.PostTagDrops(ipt.Post.Id, del)
	// 			// 添加标签
	// 			model.TagPostAdds(&tagAdds)
	// 			return ctx.JSON(utils.Succ("文章修改成功"))
	// 		}
	// 		return ctx.JSON(utils.Succ("页面修改成功"))
	// 	}
	// 	if ipt.Type == 0 {
	// 		return ctx.JSON(utils.Fail("文章修改失败,请重试"))
	// 	}
	// 	return ctx.JSON(utils.Fail("页面修改失败,请重试"))
	// }
	// // 添加 文章/页面
	// ipt.Post.Updated = time.Now()
	// if model.PostAdd(&ipt.Post) {
	// 	// 添加标签
	// 	// 文章
	// 	if ipt.Type == 0 {
	// 		//添加标签
	// 		tagPosts := make([]model.PostTag, 0, len(ipt.Tags))
	// 		for _, itm := range ipt.Tags {
	// 			tagPosts = append(tagPosts, model.PostTag{
	// 				TagId:  itm,
	// 				PostId: ipt.Post.Id,
	// 			})
	// 		}
	// 		model.TagPostAdds(&tagPosts)
	// 		return ctx.JSON(utils.Succ("文章添加成功"))
	// 	}
	// 	return ctx.JSON(utils.Succ("页面添加成功"))
	// }
	// if ipt.Type == 0 {
	// 	return ctx.JSON(utils.Fail("文章添加失败,请重试"))
	// }
	return ctx.JSON(utils.Fail("页面添加失败,请重试"))
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

// PostDrop doc
// @Tags post-文章页面
// @Summary 通过id删除单条文章
// @Param id query int true "id"
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/post/drop [post]
func PostDrop(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	err = model.PostDrop(ipt.Id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	// 删除 文章对应的标签信息
	model.PostTagDrop(ipt.Id)
	return ctx.JSON(utils.Succ("succ"))
}

func inOf(goal int, arr []int) bool {
	for idx := range arr {
		if goal == arr[idx] {
			return true
		}
	}
	return false
}

// ------------------------------------------------------ 页面相关 ------------------------------------------------------
// PostPageAll 页面列表
func PostPageAll(ctx echo.Context) error {
	mods, err := model.PostAll(-1, model.PostKindPage)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到页面信息", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到页面信息", "len"))
	}
	return ctx.JSON(utils.Succ("页面信息", mods))
}
