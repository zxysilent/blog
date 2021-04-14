package appctl

import (
	"blog/model"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// UserExist doc
// @Tags 用户
// @Summary 判断当前用户账号是否存在
// @Param num path string true "用户账号"" default('')
// @Router /user/exist/{num} [get]
func UserExist(ctx echo.Context) error {
	num := ctx.Param("num")
	if !model.UserExist(num) {
		return ctx.JSON(utils.ErrOpt(`当前账号不存在`))
	}
	return ctx.JSON(utils.Succ(`当前账号存在`))
}

// UserAdd doc
// @Tags 用户
// @Summary 添加用户信息
// @Accept mpfd
// @Param mod formData {object} true "用户信息mod"
// @Router /user/add [post]
// func UserAdd(ctx echo.Context) error {
// 	ipt := &struct {
// 		model.User
// 		Pass  string   `json:"pass" form:"pass"`
// 		Roles []uint32 `json:"roles" form:"roles"`
// 	}{}
// 	err := ctx.Bind(ipt)
// 	if err != nil {
// 		return ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
// 	}

// 	if len(ipt.Roles) == 0 {
// 		return ctx.JSON(utils.ErrIpt(`请至少选择一个权限`))
// 	}
// 	ipt.Ctime = time.Now()
// 	role := model.UserBaseRole()
// 	for idx := range ipt.Roles {
// 		role += (1 << ipt.Roles[idx])
// 	}
// 	ipt.User.Role = role
// 	ipt.User.Pass = ipt.Pass
// 	if !model.UserAdd(&ipt.User) {
// 		return ctx.JSON(utils.Fail(`用户信息添加失败`))
// 	}
// 	return ctx.JSON(utils.Succ(`用户信息添加成功`))
// }

// UserPage doc
// @Tags 用户
// @Summary 分页信息
// @Description
// @Param rl path int true "权限类型" default(27)
// @Param pi query int true "分页页数pi" default(1)
// @Param ps query int true "分页大小ps" default(6)
// @Router /user/page/{rl}?pi={pi}&ps={ps}[get]
// func UserPage(ctx echo.Context) error {
// 	auth := ctx.Get("auth").(*model.JwtClaims)
// 	rl, err := ctx.ParamInt("rl")
// 	if err != nil {
// 		return ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
// 	}
// 	if rl > 29 || rl < 27 {
// 		return ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, "范围27-29"))
// 	}
// 	ipt := &model.Page{}
// 	err = ctx.Bind(ipt)
// 	if err != nil {
// 		return ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
// 	}
// 	count := model.UserCount(rl, auth.Role)
// 	if count == 0 {
// 		return ctx.JSON(utils.ErrOpt(`未查询到用户信息,请重试`))
// 	}
// 	mods, err := model.UserPage(rl, auth.Role, ipt.Pi, ipt.Ps)
// 	if err != nil {
// 		return ctx.JSON(utils.ErrOpt(`未查询到用户信息,请重试`, err.Error()))
// 	}
// 	return ctx.JSON(utils.Page(`用户信息`, mods, count))
// }

// UserChgatv doc
// @Tags 用户
// @Summary 更新用户状态
// @Description
// @Param id path int true "用户id" default(1)
// @Router /user/chgatv/{id}[get]
// func UserChgatv(ctx echo.Context) error {
// 	id, err := ctx.ParamInt("id")
// 	if err != nil {
// 		return ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
// 	}
// 	auth := ctx.Get("auth").(*model.JwtClaims)
// 	if model.UserChgatv(id, auth.Role) {
// 		return ctx.JSON(utils.Succ("用户状态更新成功"))
// 	}
// 	return ctx.JSON(utils.Fail("用户状态更新失败"))
// }

// UserResetPass doc
// @Tags 用户
// @Summary 重置密码[123654]
// @Description
// @Param id path int true "用户id" default(1)
// @Router /user/reset/pass/{id}[get]
// func UserResetPass(ctx echo.Context) error {
// 	id, err := ctx.ParamInt("id")
// 	if err != nil {
// 		return ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
// 	}
// 	auth := ctx.Get("auth").(*model.JwtClaims)
// 	if model.UserPass(id, "33d7be2196ff70efaf6913fc8bdcab", auth.Role) {
// 		return ctx.JSON(utils.Succ("用户密码重置成功"))
// 	}
// 	return ctx.JSON(utils.Fail("用户密码重置失败"))
// }

// UserEdit doc
func UserEdit(ctx echo.Context) error {
	ipt := &struct {
		model.User
		Roles []uint32 `json:"roles" form:"roles"`
	}{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
	}

	if len(ipt.Roles) == 0 {
		return ctx.JSON(utils.ErrIpt(`请至少选择一个权限`))
	}
	// role := model.UserBaseRole()
	// for idx := range ipt.Roles {
	// 	role += (1 << ipt.Roles[idx])
	// }
	//ipt.User.Role = role
	//auth := ctx.Get("auth").(*model.JwtClaims)
	// if !model.UserEdit(&ipt.User, auth.Role, "Name", "Phone", "Email", "Desc", "Role") {
	// 	return ctx.JSON(utils.Fail(`用户信息修改失败`))
	// }
	return ctx.JSON(utils.Succ(`用户信息修改成功`))
}

// UserDrop doc
// @Tags 用户
// @Summary 删除用户
// @Description
// @Param id path int true "用户id" default(1)
// @Router /user/drop/{id} [get]
// func UserDrop(ctx echo.Context) error {
// 	id, err := ctx.ParamInt("id")
// 	if err != nil {
// 		return ctx.JSON(utils.ErrIpt(`数据输入错误,请重试`, err.Error()))
// 	}
// 	auth := ctx.Get("auth").(*model.JwtClaims)
// 	if !model.UserDrop(id, auth.Role) {
// 		return ctx.JSON(utils.Fail(`用户信息删除失败,请重试`))
// 	}
// 	return ctx.JSON(utils.Succ(`用户信息删除成功`))
// }
