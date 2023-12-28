package kernel

import (
	"blog/internal/model"
	"blog/internal/repo"
	"blog/internal/utils"
	"time"

	"github.com/labstack/echo/v4"
)

// AdminGet doc
// @Auth mgmt
// @Tags kernel
// @Summary 通过id获取admin信息
// @Param id query int true "id"
// @Security ApiKeyAuth
// @Success 200 {object} utils.Reply{data=model.Admin} "返回数据"
// @Router /api/admin/get [get]
func AdminGet(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod, has := repo.AdminGet(ipt.Id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到数据"))
	}
	return ctx.JSON(utils.Succ(`用户数据`, mod))
}

// AdminExist doc
// @Auth mgmt
// @Tags kernel
// @Summary 获取某个用户信息
// @Param num query string true "账号"
// @Success 200 {object} utils.Reply "返回数据"
// @Router /api/admin/exist [get]
func AdminExist(ctx echo.Context) error {
	num := ctx.QueryParam("num")
	if !repo.AdminExist(num) {
		return ctx.JSON(utils.Fail("不存在"))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// AdminAdd doc
// @Auth mgmt
// @Tags kernel
// @Summary 添加admin信息
// @Security ApiKeyAuth
// @Param body body model.Admin true "请求数据"
// @Success 200 {object} utils.Reply "返回数据"
// @Router /api/admin/add [post]
func AdminAdd(ctx echo.Context) error {
	ipt := &struct {
		model.Admin
		Passwd string `json:"passwd"`
	}{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	roleId, _ := ctx.Get("rid").(int)
	if ipt.RoleId <= roleId {
		return ctx.JSON(utils.ErrDeny("不能添加同等角色用户"))
	}
	ipt.Admin.Passwd = ipt.Admin.Enc(ipt.Passwd)
	ipt.Created = time.Now().UnixMilli()
	ipt.Updated = ipt.Created
	err = repo.AdminAdd(&ipt.Admin)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	repo.LOG(ctx, "admin", "添加用户", ipt)
	return ctx.JSON(utils.Succ("succ"))
}

// AdminEdit doc
// @Auth mgmt
// @Tags kernel
// @Summary 修改admin信息
// @Security ApiKeyAuth
// @Param body body model.Admin true "请求数据"
// @Success 200 {object} utils.Reply "返回数据"
// @Router /api/admin/edit [post]
func AdminEdit(ctx echo.Context) error {
	ipt := &struct {
		model.Admin
		Passwd string `json:"passwd"`
	}{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	roleId, _ := ctx.Get("rid").(int)
	mod, _ := repo.AdminGet(ipt.Id)
	if mod.RoleId <= roleId {
		return ctx.JSON(utils.ErrDeny("不能修改同等角色用户"))
	}
	ipt.Updated = time.Now().UnixMilli()
	cols := []string{"name", "phone", "avatar", "email", "role_id", "updated"}
	if ipt.Passwd != "" {
		cols = append(cols, "passwd")
		ipt.Admin.Passwd = ipt.Admin.Enc(ipt.Passwd)
	}
	err = repo.AdminEdit(&ipt.Admin, cols...)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	repo.LOG(ctx, "admin", "修改用户", ipt)
	return ctx.JSON(utils.Succ("succ"))
}

// AdminEditLock doc
// @Auth mgmt mgmt
// @Tags kernel
// @Summary 修改用户锁定状态
// @Security ApiKeyAuth
// @Param body body object{id=int,lock=bool} true "json"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/admin/edit/lock [post]
func AdminEditLock(ctx echo.Context) error {
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
	mod, _ := repo.AdminGet(ipt.Id)
	if mod.RoleId <= roleId {
		return ctx.JSON(utils.ErrDeny("不能修改高角色用户"))
	}
	//	mod = model.Admin{Id: ipt.Id, Lock: ipt.Lock}
	err = repo.AdminEdit(mod, "Lock")
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// AdminDrop doc
// @Auth mgmt
// @Tags kernel
// @Summary 删除admin信息
// @Param id query int true "id"
// @Security ApiKeyAuth
// @Success 200 {object} utils.Reply "返回数据"
// @Router /api/admin/drop [post]
func AdminDrop(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if ipt.Id == ctx.Get("uid").(int) {
		return ctx.JSON(utils.ErrOpt("不能删除自己"))
	}
	roleId, _ := ctx.Get("rid").(int)
	mod, _ := repo.AdminGet(ipt.Id)
	if mod.RoleId <= roleId {
		return ctx.JSON(utils.ErrDeny("不能删除高角色用户"))
	}
	err = repo.AdminDrop(ipt.Id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	repo.LOG(ctx, "admin", "删除用户", ipt)
	return ctx.JSON(utils.Succ("succ"))
}

// AdminPage doc
// @Auth mgmt
// @Tags kernel
// @Summary 获取admin分页数据
// @Security ApiKeyAuth
// @Param query query model.AdminRolePageArgs true "请求数据"
// @Success 200 {object} utils.Reply "返回数据"
// @Router /api/admin/page [get]
func AdminPage(ctx echo.Context) error {
	roleId, _ := ctx.Get("rid").(int)
	ipt := &model.AdminRolePageArgs{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if err = ipt.Stat(); err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mods, count, err := repo.AdminPage(ipt.Pi, ipt.Ps, ipt.Mult, roleId, ipt.RoleId)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "items eq 0"))
	}
	return ctx.JSON(utils.Page("succ", mods, count))
}

// AdminPage doc
// @Auth mgmt
// @Tags kernel
// @Summary 获取admin分页数据
// @Param query query model.AdminMultArgs true "请求数据"
// @Security ApiKeyAuth
// @Success 200 {object} utils.Reply "返回数据"
// @Router /api/admin/all [get]
func AdminAll(ctx echo.Context) error {
	roleId, _ := ctx.Get("rid").(int)
	ipt := &model.AdminMultArgs{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mods, err := repo.AdminAll(ipt.Mult, roleId, ipt.RoleId)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "items eq 0"))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}
