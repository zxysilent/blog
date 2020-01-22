package control

import (
	"blog/model"

	"github.com/labstack/echo"
	"github.com/zxysilent/utils"
)

// OptsGet 获取某个配置项
func OptsGet(ctx echo.Context) error {
	key := ctx.Param("key")
	if key == "" {
		return ctx.JSON(utils.NewErrIpt(`请填写key值`))
	}
	if val, ok := model.OptsGet(key); ok {
		return ctx.JSON(utils.NewSucc(``, val))
	}
	return ctx.JSON(utils.NewErrIpt(`错误的key值`))
}

// OptsEdit 编辑某个配置项
func OptsEdit(ctx echo.Context) error {
	ipt := &model.Opts{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.NewErrIpt(`数据输入错误,请重试`, err.Error()))
	}
	if !model.OptsEdit(ipt) {
		return ctx.JSON(utils.Fail(`配置项修改失败`))
	}
	return ctx.JSON(utils.NewSucc(`配置项修改成功`))
}

// OptsBase 基本配置项目
func OptsBase(ctx echo.Context) error {
	// ipt := &model.Opts{}
	// err := ctx.Bind(ipt)
	// if err != nil {
	// 	return ctx.JSON(utils.NewErrIpt(`数据输入错误,请重试`, err.Error()))
	// }
	// if !model.OptsEdit(ipt) {
	// 	return ctx.JSON(utils.Fail(`配置项修改失败`))
	// }
	mp := model.MapOpts
	// delete(mp, "analytic")
	// delete(mp, "comment")
	return ctx.JSON(utils.NewSucc(`基本配置项目`, mp))
}
