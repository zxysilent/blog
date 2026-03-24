package kernel

import (
	"blog/internal/model"
	"blog/internal/repo"
	"blog/internal/utils"
	"blog/pkg/reply"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/logs"
)

// AdminGet doc
// @Auth mgmt
// @Tags kernel,admin
// @Summary 通过id获取admin信息
// @Param id query int true "id"
// @Security KeyAuth
// @Success 200 {object} reply.Reply{data=model.Admin} "返回数据"
// @Router /api/admin/{id} [get]
func AdminGet(ctx echo.Context) error {
	rctx := utils.WrapCtx(ctx)
	pk := ctx.Param("id")
	mod, err := repo.GetAny[model.Admin](rctx, pk)
	if err != nil {
		logs.Ctx(rctx).Str("model", "admin").Str("pk", pk).Error("repo.GetAny")
		return reply.Fail(ctx, "未查询到数据", err)
	}
	return reply.Succ(ctx, "用户数据", mod)
}

// AdminExist doc
// @Auth mgmt
// @Tags kernel,admin
// @Summary 获取某个用户信息
// @Param num query string true "账号"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/admin/{num}/exist [get]
func AdminExist(ctx echo.Context) error {
	rctx := utils.WrapCtx(ctx)
	num := ctx.Param("num")
	if !repo.AdminExist(num) {
		logs.Ctx(rctx).Error("repo.AdminExist")
		return reply.Fail(ctx, "未查询到数据")
	}
	return reply.Succ(ctx, "succ")
}

// AdminAdd doc
// @Auth mgmt
// @Tags kernel,admin
// @Summary 添加admin信息
// @Security KeyAuth
// @Param body body model.Admin true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/admin/add [post]
func AdminAdd(ctx echo.Context) error {
	// ipt := &struct {
	// 	model.Admin
	// 	Passwd string `json:"passwd"`
	// }{}
	// err := ctx.Bind(ipt)
	// if err != nil {
	// 	return reply.Fail(ctx, "输入有误", err)
	// }
	// roleId, _ := ctx.Get("rid").(string)
	// if roleId != "" {
	// 	return ctx.JSON(utils.ErrDeny("不能添加同等角色用户"))
	// }
	// // ipt.Admin.Passwd = ipt.Admin.Enc(ipt.Passwd)
	// ipt.Created = time.Now().UnixMilli()
	// ipt.Updated = ipt.Created
	// err = repo.AdminAdd(&ipt.Admin)
	// if err != nil {
	// 	return reply.Fail(ctx, "添加失败", err)
	// }
	// repo.LOG(ctx, "admin", "添加用户", ipt)
	return reply.Succ(ctx, "succ")
}

// AdminEdit doc
// @Auth mgmt
// @Tags kernel,admin
// @Summary 修改admin信息
// @Security KeyAuth
// @Param body body model.Admin true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/admin/edit [post]
func AdminEdit(ctx echo.Context) error {
	ipt := &struct {
		model.Admin
		Passwd string `json:"passwd"`
	}{}
	err := ctx.Bind(ipt)
	if err != nil {
		return reply.Fail(ctx, "输入有误", err)
	}
	// roleId, _ := ctx.Get("rid").(int)
	// mod, _ := repo.AdminGet(ipt.Id)
	// if mod.RoleId <= roleId {
	// 	return ctx.JSON(utils.ErrDeny("不能修改同等角色用户"))
	// }
	// ipt.Updated = time.Now().UnixMilli()
	// cols := []string{"name", "phone", "avatar", "email", "role_id", "updated"}
	// if ipt.Passwd != "" {
	// 	cols = append(cols, "passwd")
	// 	// ipt.Admin.Passwd = ipt.Admin.Enc(ipt.Passwd)
	// }
	// err = repo.AdminEdit(&ipt.Admin, cols...)
	// if err != nil {
	// 	return reply.Fail(ctx, "修改失败", err)
	// }
	// repo.LOG(ctx, "admin", "修改用户", ipt)
	return reply.Succ(ctx, "succ")
}

// AdminEditLock doc
// @Auth mgmt
// @Tags kernel,admin
// @Summary 修改用户锁定状态
// @Security KeyAuth
// @Param body body object{id=int,lock=bool} true "json"
// @Success 200 {object} reply.Reply{data=string} "返回数据"
// @Router /api/admin/edit/lock [post]
func AdminEditLock(ctx echo.Context) error {
	ipt := &struct {
		Id   int  `json:"id"`
		Lock bool `json:"lock"`
	}{}
	err := ctx.Bind(ipt)
	if err != nil {
		return reply.Fail(ctx, "输入有误", err)
	}
	// if ipt.Id == ctx.Get("uid").(int) {
	// 	return ctx.JSON(utils.ErrOpt("不能修改自己的状态"))
	// }
	// roleId, _ := ctx.Get("rid").(int)
	// mod, _ := repo.AdminGet(ipt.Id)
	// if mod.RoleId <= roleId {
	// 	return ctx.JSON(utils.ErrDeny("不能修改高角色用户"))
	// }
	// //	mod = model.Admin{Id: ipt.Id, Lock: ipt.Lock}
	// err = repo.AdminEdit(mod, "Lock")
	// if err != nil {
	// 	return reply.Fail(ctx, "修改失败", err)
	// }
	return reply.Succ(ctx, "succ")
}

// AdminDrop doc
// @Auth mgmt
// @Tags kernel,admin
// @Summary 删除admin信息
// @Param id query int true "id"
// @Security KeyAuth
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/admin/drop [post]
func AdminDrop(ctx echo.Context) error {
	// ipt := &model.IptId{}
	// err := ctx.Bind(ipt)
	// if err != nil {
	// 	return reply.Fail(ctx, "输入有误", err)
	// }
	// if ipt.Id == ctx.Get("uid").(int) {
	// 	return ctx.JSON(utils.ErrOpt("不能删除自己"))
	// }
	// roleId, _ := ctx.Get("rid").(int)
	// mod, _ := repo.AdminGet(ipt.Id)
	// if mod.RoleId <= roleId {
	// 	return ctx.JSON(utils.ErrDeny("不能删除高角色用户"))
	// }
	// err = repo.AdminDrop(ipt.Id)
	// if err != nil {
	// 	return reply.Fail(ctx, "删除失败", err)
	// }
	// repo.LOG(ctx, "admin", "删除用户", ipt)
	return reply.Succ(ctx, "succ")
}

// AdminPage doc
// @Auth mgmt
// @Tags kernel,admin
// @Summary 获取admin分页数据
// @Security KeyAuth
// @Param query query model.AdminRolePageArgs true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/admin/page [get]
func AdminPage(ctx echo.Context) error {
	// roleId, _ := ctx.Get("rid").(int)
	// ipt := &model.AdminRolePageArgs{}
	// err := ctx.Bind(ipt)
	// if err != nil {
	// 	return reply.Fail(ctx, "输入有误", err)
	// }
	// if err = ipt.Stat(); err != nil {
	// 	return reply.Fail(ctx, "输入有误", err)
	// }
	// mods, count, err := repo.AdminPage(ipt.Pi, ipt.Ps, ipt.Mult, roleId, ipt.RoleId)
	// if err != nil {
	// 	return reply.Fail(ctx, "查询数据错误", err)
	// }
	// if len(mods) < 1 {
	// 	return ctx.JSON(utils.ErrOpt("未查询到数据", "items eq 0"))
	// }
	// return reply.List(ctx, "succ", mods, count)
	return nil
}

// AdminPage doc
// @Auth mgmt
// @Tags kernel,admin
// @Summary 获取admin分页数据
// @Param query query model.AdminMultArgs true "请求数据"
// @Security KeyAuth
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/admin/all [get]
func AdminAll(ctx echo.Context) error {
	// roleId, _ := ctx.Get("rid").(int)
	// ipt := &model.AdminMultArgs{}
	// err := ctx.Bind(ipt)
	// if err != nil {
	// 	return reply.Fail(ctx, "输入有误", err)
	// }
	// mods, err := repo.AdminAll(ipt.Mult, roleId, ipt.RoleId)
	// if err != nil {
	// 	return reply.Fail(ctx, "查询数据错误", err)
	// }
	// if len(mods) < 1 {
	// 	return ctx.JSON(utils.ErrOpt("未查询到数据", "items eq 0"))
	// }
	// return ctx.JSON(utils.Succ("succ", mods))
	return nil
}
