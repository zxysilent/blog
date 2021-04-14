package sysctl

import (
	"blog/model"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// UserGet doc
// @Tags user
// @Summary 通过id获取user信息
// @Param id query int true "id"
// @Success 200 {object} model.Reply{data=model.User} "成功数据"
// @Router /adm/user/get [get]
func UserGet(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod, has := model.UserGet(ipt.Id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到数据"))
	}
	return ctx.JSON(utils.Succ("用户数据", mod))
}

// UserExist doc
// @Tags user
// @Summary 获取某个用户信息
// @Param num query string true "账号"
// @Success 200 {object} model.Reply "成功数据"
// @Router /api/user/exist [get]
func UserExist(ctx echo.Context) error {
	num := ctx.QueryParam("num")
	if !model.UserExist(num) {
		return ctx.JSON(utils.Fail("不存在"))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// UserAdd doc
// @Tags user
// @Summary 添加user信息
// @Param token query string true "凭证"
// @Param body body model.User true "request"
// @Success 200 {object} model.Reply "成功数据"
// @Router /adm/user/add [post]
func UserAdd(ctx echo.Context) error {
	// role, _ := ctx.Get("role").(int)
	// if role < model.RoleAdmin {
	// 	return ctx.JSON(utils.ErrDeny("对不起您没有此权限"))
	// }
	ipt := &struct {
		model.User
		Pass string `json:"pass"`
	}{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	ipt.User.Passwd = ipt.Pass
	ipt.Ltime = time.Now()
	ipt.Ctime = ipt.Ltime
	err = model.UserAdd(&ipt.User)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// UserEdit doc
// @Tags user
// @Summary 修改user信息
// @Param token query string true "凭证"
// @Param body body model.User true "request"
// @Success 200 {object} model.Reply "成功数据"
// @Router /adm/user/edit [post]
func UserEdit(ctx echo.Context) error {
	// role, _ := ctx.Get("role").(int)
	// if role < model.RoleAdmin {
	// 	return ctx.JSON(utils.ErrDeny("对不起您没有此权限"))
	// }
	ipt := &model.User{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	err = model.UserEdit(ipt, "Name", "Phone", "Addr")
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// UserReset doc
// @Tags user
// @Summary 重置密码
// @Param id query int true "id"
// @Param token query string true "凭证"
// @Success 200 {object} model.Reply "成功数据"
// @Router /adm/user/reset [get]
func UserReset(ctx echo.Context) error {
	// role, _ := ctx.Get("role").(int)
	// if role < model.RoleAdmin {
	// 	return ctx.JSON(utils.ErrDeny("对不起您没有此权限"))
	// }
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if ipt.Id == ctx.Get("uid").(int) {
		return ctx.JSON(utils.ErrOpt("不能重置自己d的密码"))
	}
	mod := model.User{Id: ipt.Id, Passwd: "fde6bb0541387e4ebdadf7c2ff3112"} //1q2w3e
	err = model.UserEdit(&mod, "passwd")
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// UserDrop doc
// @Tags user
// @Summary 删除user信息
// @Param id query int true "id"
// @Param token query string true "凭证"
// @Success 200 {object} model.Reply "成功数据"
// @Router /adm/user/drop [post]
func UserDrop(ctx echo.Context) error {
	// role, _ := ctx.Get("role").(int)
	// if role < model.RoleAdmin {
	// 	return ctx.JSON(utils.ErrDeny("对不起您没有此权限"))
	// }
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if ipt.Id == ctx.Get("uid").(int) {
		return ctx.JSON(utils.ErrOpt("不能删除自己"))
	}
	err = model.UserDrop(ipt.Id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// UserPage doc
// @Tags user
// @Summary 获取分页数据
// @Param pi query int true "分页数"
// @Param ps query int true "每页条数[5,20]" default(5)
// @Success 200 {object} model.Reply "成功数据"
// @Router /api/user/page [get]
func UserPage(ctx echo.Context) error {
	// role, _ := ctx.Get("role").(int)
	// if role < model.RoleAdmin {
	// 	return ctx.JSON(utils.ErrDeny("对不起您没有此权限"))
	// }
	ipt := &model.Page{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if ipt.Ps > 20 || ipt.Ps < 5 {
		return ctx.JSON(utils.ErrIpt("分页大小输入错误", ipt.Ps))
	}
	count := model.UserCount()
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := model.UserPage(ipt.Pi, ipt.Ps)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, int(count)))
}
