package sysctl

import (
	"blog/model"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// RoleGet doc
// @Tags role-角色
// @Summary 通过id获取单条角色信息
// @Param id query int true "id"
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=model.Role} "返回数据"
// @Router /adm/role/get [get]
func RoleGet(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod, has := model.RoleGet(ipt.Id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到角色信息"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// RoleAll doc
// @Tags role-角色
// @Summary 获取所有角色信息
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=[]model.Role} "返回数据"
// @Router /adm/role/all [get]
func RoleAll(ctx echo.Context) error {
	// 只能查看角色id大于等于自己id的角色
	rid, ok := ctx.Get("rid").(int)
	if !ok {
		return ctx.JSON(utils.ErrOpt("角色信息异常"))
	}
	mods, err := model.RoleAll(rid)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到角色信息", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// RolePage doc
// @Tags role-角色
// @Summary 获取角色分页信息
// @Param pi query int true "分页数" default(1)
// @Param ps query int true "每页条数[5,30]" default(5)
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=[]model.Role} "返回数据"
// @Param token query string true "token"
// @Router /adm/role/page [get]
func RolePage(ctx echo.Context) error {
	ipt := &model.Page{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if ipt.Ps > 20 || ipt.Ps < 5 {
		return ctx.JSON(utils.ErrIpt("分页大小输入错误", ipt.Ps))
	}
	count := model.RoleCount()
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := model.RolePage(ipt.Pi, ipt.Ps)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, int(count)))
}

// RoleAdd doc
// @Tags role-角色
// @Summary 添加角色信息
// @Param token query string true "token"
// @Param body body model.Role true "请求数据"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/role/add [post]
func RoleAdd(ctx echo.Context) error {
	ipt := &model.Role{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	ipt.Ctime = time.Now()
	err = model.RoleAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// RoleEdit doc
// @Tags role-角色
// @Summary 修改角色信息
// @Param token query string true "token"
// @Param body body model.Role true "请求数据"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/role/edit [post]
func RoleEdit(ctx echo.Context) error {
	ipt := &model.Role{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if ipt.Id == model.RootRoleId {
		return ctx.JSON(utils.Fail("内置角色不能修改"))
	}
	err = model.RoleEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// RoleDrop doc
// @Tags role-角色
// @Summary 通过id删除单条角色信息
// @Param token query string true "token"
// @Param body body model.IptId true "json"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/role/drop [post]
func RoleDrop(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod, has := model.RoleGet(ipt.Id)
	if !has {
		return ctx.JSON(utils.ErrOpt("删除失败", "不存在当前数据"))
	}
	if mod.Inner {
		return ctx.JSON(utils.ErrOpt("内置数据无法删除"))
	}
	err = model.RoleDrop(ipt.Id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// RoleGrantAll doc
// @Tags role-角色-角色权限
// @Summary 通过角色id获取所有角色授权信息
// @Param id query int true "id"
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=model.Role} "返回数据"
// @Router /adm/role/grant/all [get]
func RoleGrantAll(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mods, err := model.RoleGrantAll(ipt.Id)
	if err != nil {
		return ctx.JSON(utils.Fail("未查询到角色授权信息", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// RoleGrantEdit doc
// @Tags role-角色-角色权限
// @Summary 修改角色授权信息
// @Param id query int true "id"
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=string} "返回数据"
// @Router /adm/role/grant/edit [post]
func RoleGrantEdit(ctx echo.Context) error {
	ipt := &struct {
		RoleId   int   `json:"role_id"`
		GrantIds []int `json:"grant_ids"`
	}{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	rid, _ := ctx.Get("rid").(int)
	if rid == ipt.RoleId {
		return ctx.JSON(utils.Fail("不能修改自己所属角色"))
	}
	// 查找旧的授权
	oldGrants, _ := model.RoleGrantAll(ipt.RoleId)
	oldMapIds := make(map[int]bool, len(oldGrants))
	for _, item := range oldGrants {
		oldMapIds[item.Id] = true
	}
	newMapIds := make(map[int]bool, len(ipt.GrantIds))
	for _, item := range ipt.GrantIds {
		newMapIds[item] = true
	}
	// 添加id集合
	addIds := make([]int, 0, len(ipt.GrantIds))
	// 遍历新的集合
	for _, item := range ipt.GrantIds {
		// 如果没有在旧的集合里面,就添加
		if _, ok := oldMapIds[item]; !ok {
			addIds = append(addIds, item)
		}
	}
	// 删除id集合
	dropIds := make([]int, 0, 8)
	// 遍历旧的集合
	for _, item := range oldGrants {
		// 如果没在新集合里面,就添加
		if _, ok := newMapIds[item.Id]; !ok {
			dropIds = append(dropIds, item.Id)
		}
	}
	// 添加新的角色授权集合
	addRoleGrants := make([]model.RoleGrant, 0, len(addIds))
	for _, item := range addIds {
		addRoleGrants = append(addRoleGrants, model.RoleGrant{
			RoleId:  ipt.RoleId,
			GrantId: item,
		})
	}
	// model.RoleGrantAddMulti(addRoleGrants)
	// // 删除旧的角色授权集合
	// model.RoleGrantDropMulti(dropIds)
	// 合并到一个事务
	err = model.RoleGrantEdit(addRoleGrants, dropIds)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
