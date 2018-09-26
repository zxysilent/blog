package control

import (
	"blog/model"
	"blog/util"

	"github.com/labstack/echo"
)

// OptsGet 获取某个配置项
func OptsGet(ctx echo.Context) error {
	key := ctx.Param("key")
	if key == "" {
		return ctx.Res(util.NewErrIpt(`请填写key值`))
	}
	if val, ok := model.OptsGet(key); ok {
		return ctx.Res(util.NewSucc(``, val))
	}
	return ctx.Res(util.NewErrIpt(`错误的key值`))
}

// OptsEdit 编辑某个配置项
func OptsEdit(ctx echo.Context) error {
	ipt := &model.Opts{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.Res(util.NewErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	if !model.OptsEdit(ipt) {
		return ctx.Res(util.NewFail(`配置项修改失败`))
	}
	return ctx.Res(util.NewSucc(`配置项修改成功`))
}

// OptsBase 基本配置项目
func OptsBase(ctx echo.Context) error {
	// ipt := &model.Opts{}
	// err := ctx.Bind(ipt)
	// if err != nil {
	// 	return ctx.Res(util.NewErrIpt(`数据输入错误,请重试`, err.Error()))
	// }
	// if !model.OptsEdit(ipt) {
	// 	return ctx.Res(util.NewFail(`配置项修改失败`))
	// }
	mp := model.MapOpts
	// delete(mp, "analytic")
	// delete(mp, "comment")
	return ctx.Res(util.NewSucc(`基本配置项目`, mp))
}
