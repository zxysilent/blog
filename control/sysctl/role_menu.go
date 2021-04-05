package sysctl

import (
	"blog/model"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/utils"
)

// RoleMenuAll doc
// @Tags role,menu
// @Summary 通过角色id获取所有角色菜单信息
// @Param id query int true "id"
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=model.Role} "成功数据"
// @Router /adm/role/menu/all [get]
func RoleMenuAll(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mods, err := model.RoleMenuAll(ipt.Id)
	if err != nil {
		return ctx.JSON(utils.Fail("未查询到角色菜单导航信息", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// RoleMenuEdit doc
// @Tags role,menu
// @Summary 修改菜单导航信息菜单导航
// @Param id query int true "id"
// @Param token query string true "token"
// @Success 200 {object} model.Reply{data=string} "成功数据"
// @Router /adm/role/menu/edit [post]
func RoleMenuEdit(ctx echo.Context) error {
	ipt := &struct {
		RoleId  int   `json:"role_id"`
		MenuIds []int `json:"menu_ids"`
	}{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// 查找旧的菜单
	oldMenus, _ := model.RoleMenuAll(ipt.RoleId)
	oldMapIds := make(map[int]bool, len(oldMenus))
	for _, item := range oldMenus {
		oldMapIds[item.Id] = true
	}
	newMapIds := make(map[int]bool, len(ipt.MenuIds))
	for _, item := range ipt.MenuIds {
		newMapIds[item] = true
	}
	// 添加id集合
	addIds := make([]int, 0, len(ipt.MenuIds))
	// 遍历新的集合
	for _, item := range ipt.MenuIds {
		// 如果没有在旧的集合里面,就添加
		if _, ok := oldMapIds[item]; !ok {
			addIds = append(addIds, item)
		}
	}
	// 删除id集合
	dropIds := make([]int, 0, 8)
	// 遍历旧的集合
	for _, item := range oldMenus {
		// 如果没在新集合里面,就添加
		if _, ok := newMapIds[item.Id]; !ok {
			dropIds = append(dropIds, item.Id)
		}
	}
	// 添加新的角色菜单集合
	addRoleMenus := make([]model.RoleMenu, 0, len(addIds))
	for _, item := range addIds {
		addRoleMenus = append(addRoleMenus, model.RoleMenu{
			RoleId: ipt.RoleId,
			MenuId: item,
			Ctime:  time.Now(),
		})
	}
	model.RoleMenuAddMulti(addRoleMenus)
	// 删除旧的角色菜单集合
	model.RoleMenuDropMulti(dropIds)
	// err = model.RoleMenuEdit(ipt)
	// if err != nil {
	// 	return ctx.JSON(utils.Fail("修改失败", err.Error()))
	// }
	return ctx.JSON(utils.Succ("succ"))
}
