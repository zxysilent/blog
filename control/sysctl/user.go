package sysctl

import (
	"blog/model"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// UserGet doc
// @Tags user-用户
// @Summary 通过id获取user信息
// @Param id query int true "id"
// @Success 200 {object} model.Reply{data=model.User} "返回数据"
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
	roleId, _ := ctx.Get("rid").(int)
	if mod.RoleId < roleId {
		return ctx.JSON(utils.ErrDeny("不能查看高角色用户"))
	}
	return ctx.JSON(utils.Succ("用户数据", mod))
}

// UserExist doc
// @Tags user-用户
// @Summary 获取某个用户信息
// @Param num query string true "账号"
// @Success 200 {object} model.Reply "返回数据"
// @Router /api/user/exist [get]
func UserExist(ctx echo.Context) error {
	num := ctx.QueryParam("num")
	if !model.UserExist(num) {
		return ctx.JSON(utils.Fail("不存在"))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// UserAdd doc
// @Tags user-用户
// @Summary 添加user信息
// @Param token query string true "凭证"
// @Param body body model.User true "请求数据"
// @Success 200 {object} model.Reply "返回数据"
// @Router /adm/user/add [post]
func UserAdd(ctx echo.Context) error {
	ipt := &struct {
		model.User
		Pass string `json:"pass"`
	}{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	roleId, _ := ctx.Get("rid").(int)
	if ipt.RoleId <= roleId {
		return ctx.JSON(utils.ErrDeny("不能添加同等角色用户"))
	}
	ipt.User.Passwd = ipt.Pass
	ipt.Ctime = time.Now()
	err = model.UserAdd(&ipt.User)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// UserEdit doc
// @Tags user-用户
// @Summary 修改user信息
// @Param token query string true "凭证"
// @Param body body model.User true "请求数据"
// @Success 200 {object} model.Reply "返回数据"
// @Router /adm/user/edit [post]
func UserEdit(ctx echo.Context) error {
	ipt := &model.User{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	roleId, _ := ctx.Get("rid").(int)
	mod, _ := model.UserGet(ipt.Id)
	if mod.RoleId <= roleId {
		return ctx.JSON(utils.ErrDeny("不能修改同等角色用户"))
	}
	err = model.UserEdit(ipt, "name", "phone", "addr", "role_id")
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// UserEditLock doc
// @Tags user-用户
// @Summary 修改用户锁定状态
// @Param token query string true "token"
// @Param body body object{id=int,lock=bool} true "json"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/user/edit/lock [post]
func UserEditLock(ctx echo.Context) error {
	ipt := &struct {
		Id   int  `json:"id"`
		Lock bool `json:"lock"`
	}{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if ipt.Id == ctx.Get("uid").(int) {
		return ctx.JSON(utils.ErrOpt("不能修改自己的状态"))
	}
	roleId, _ := ctx.Get("rid").(int)
	mod, _ := model.UserGet(ipt.Id)
	if mod.RoleId <= roleId {
		return ctx.JSON(utils.ErrDeny("不能修改同等角色用户"))
	}
	mod = &model.User{Id: ipt.Id, Lock: ipt.Lock}
	err = model.UserEdit(mod, "Lock")
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// UserEditReset doc
// @Tags user-用户
// @Summary 重置密码
// @Param id query int true "id"
// @Param token query string true "凭证"
// @Success 200 {object} model.Reply "返回数据"
// @Router /adm/user/edit/reset [post]
func UserEditReset(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	id, _ := ctx.Get("uid").(int)
	if ipt.Id == id {
		return ctx.JSON(utils.ErrOpt("不能重置自己的密码"))
	}
	roleId, _ := ctx.Get("rid").(int)
	mod, _ := model.UserGet(ipt.Id)
	if mod.RoleId <= roleId {
		return ctx.JSON(utils.ErrDeny("不能重置同等角色用户"))
	}
	mod = &model.User{Id: ipt.Id, Passwd: "fde6bb0541387e4ebdadf7c2ff3112"} //1q2w3e
	err = model.UserEdit(mod, "passwd")
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// UserDrop doc
// @Tags user-用户
// @Summary 删除user信息
// @Param id query int true "id"
// @Param token query string true "凭证"
// @Success 200 {object} model.Reply "返回数据"
// @Router /adm/user/drop [post]
func UserDrop(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if ipt.Id == ctx.Get("uid").(int) {
		return ctx.JSON(utils.ErrOpt("不能删除自己"))
	}
	roleId, _ := ctx.Get("rid").(int)
	mod, _ := model.UserGet(ipt.Id)
	if mod.RoleId <= roleId {
		return ctx.JSON(utils.ErrDeny("不能删除同等角色用户"))
	}
	err = model.UserDrop(ipt.Id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// UserPage doc
// @Tags user-用户
// @Summary 获取分页数据
// @Param pi query int true "分页数"
// @Param ps query int true "每页条数[5,20]" default(5)
// @Success 200 {object} model.Reply "返回数据"
// @Router /api/user/page [get]
func UserPage(ctx echo.Context) error {
	roleId, _ := ctx.Get("rid").(int)
	ipt := &model.Page{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if ipt.Ps > 20 || ipt.Ps < 5 {
		return ctx.JSON(utils.ErrIpt("分页大小输入错误", ipt.Ps))
	}
	count := model.UserCount(roleId)
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := model.UserPage(ipt.Pi, ipt.Ps, roleId)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, int(count)))
}
