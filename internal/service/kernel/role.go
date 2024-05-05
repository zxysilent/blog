package kernel

import (
	"blog/internal/model"
	"blog/internal/repo"
	"blog/internal/utils"
	"time"

	"github.com/labstack/echo/v4"
)

// RoleGet doc
// @Auth mgmt
// @Tags kernel,role
// @Summary 通过id获取单条角色信息
// @Param id query int true "id"
// @Security ApiKeyAuth
// @Success 200 {object} utils.Reply{data=model.Role} "返回数据"
// @Router /api/role/get [get]
func RoleGet(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod, has := repo.RoleGet(ipt.Id)
	if !has {
		return ctx.JSON(utils.ErrOpt("未查询到角色信息"))
	}
	return ctx.JSON(utils.Succ("succ", mod))
}

// RoleAll doc
// @Auth mgmt
// @Tags kernel,role
// @Summary 获取所有角色信息
// @Security ApiKeyAuth
// @Success 200 {object} utils.Reply{data=[]model.Role} "返回数据"
// @Router /api/role/all [get]
func RoleAll(ctx echo.Context) error {
	// 只能查看角色id大于等于自己id的角色
	rid, ok := ctx.Get("rid").(int)
	if !ok {
		return ctx.JSON(utils.ErrOpt("角色信息异常"))
	}
	mods, err := repo.RoleAll(rid)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("未查询到角色信息", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// RolePage doc
// @Auth mgmt
// @Tags kernel,role
// @Summary 获取角色分页信息
// @Security ApiKeyAuth
// @Param query query model.Page true "请求数据"
// @Success 200 {object} utils.Reply{data=[]model.Role} "返回数据"
// @Security ApiKeyAuth
// @Router /api/role/page [get]
func RolePage(ctx echo.Context) error {
	ipt := &model.Page{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	if err = ipt.Stat(); err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	// 只能查看角色id大于等于自己id的角色
	rid, ok := ctx.Get("rid").(int)
	if !ok {
		return ctx.JSON(utils.ErrOpt("角色信息异常"))
	}
	count := repo.RoleCount(rid)
	if count < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", " count < 1"))
	}
	mods, err := repo.RolePage(rid, ipt.Pi, ipt.Ps)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("查询数据错误", err.Error()))
	}
	if len(mods) < 1 {
		return ctx.JSON(utils.ErrOpt("未查询到数据", "len(mods) < 1"))
	}
	return ctx.JSON(utils.Page("succ", mods, count))
}

// RoleAdd doc
// @Auth mgmt
// @Tags kernel,role
// @Summary 添加角色信息
// @Security ApiKeyAuth
// @Param body body model.Role true "请求数据"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/role/add [post]
func RoleAdd(ctx echo.Context) error {
	ipt := &model.Role{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	ipt.Created = time.Now().UnixMilli()
	ipt.Updated = ipt.Created
	err = repo.RoleAdd(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("添加失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// RoleEdit doc
// @Auth mgmt
// @Tags kernel,role
// @Summary 修改角色信息
// @Security ApiKeyAuth
// @Param body body model.Role true "请求数据"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/role/edit [post]
func RoleEdit(ctx echo.Context) error {
	ipt := &model.Role{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	ipt.Updated = time.Now().UnixMilli()
	err = repo.RoleEdit(ipt)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// RoleDrop doc
// @Auth mgmt
// @Tags kernel,role
// @Summary 通过id删除单条角色信息
// @Security ApiKeyAuth
// @Param body body model.IptId true "json"
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/role/drop [post]
func RoleDrop(ctx echo.Context) error {
	ipt := &model.IptId{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mod, has := repo.RoleGet(ipt.Id)
	if !has {
		return ctx.JSON(utils.ErrOpt("删除失败", "不存在当前数据"))
	}
	if mod.Inner {
		return ctx.JSON(utils.ErrOpt("内置数据无法删除"))
	}
	err = repo.RoleDrop(ipt.Id)
	if err != nil {
		return ctx.JSON(utils.ErrOpt("删除失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}

// RoleGrantAll doc
// @Auth mgmt
// @Tags kernel,role
// @Summary 通过角色id获取所有角色授权信息
// @Param role_id query int true "id"
// @Security ApiKeyAuth
// @Success 200 {object} utils.Reply{data=model.Role} "返回数据"
// @Router /api/role/grant/all [get]
func RoleGrantAll(ctx echo.Context) error {
	ipt := &struct {
		RoleId int `form:"role_id" query:"role_id" json:"role_id"`
	}{}
	err := ctx.Bind(ipt)
	if err != nil {
		return ctx.JSON(utils.ErrIpt("输入有误", err.Error()))
	}
	mods, err := repo.RoleGrantAll(ipt.RoleId)
	if err != nil {
		return ctx.JSON(utils.Fail("未查询到角色授权信息", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ", mods))
}

// RoleGrantEdit doc
// @Auth mgmt
// @Tags kernel,role
// @Summary 修改角色授权信息
// @Param id query int true "id"
// @Security ApiKeyAuth
// @Success 200 {object} utils.Reply{data=string} "返回数据"
// @Router /api/role/grant/edit [post]
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
	// if ipt.RoleId < 10 {
	// 	return ctx.JSON(utils.Fail("内置角色不能修改"))
	// }
	// 查找旧的授权
	oldGrants, _ := repo.RoleGrantAll(ipt.RoleId)
	oldIds := make(map[int]bool, len(oldGrants))
	for _, item := range oldGrants {
		oldIds[item.Id] = true
	}
	newIds := make(map[int]bool, len(ipt.GrantIds))
	for _, item := range ipt.GrantIds {
		newIds[item] = true
	}
	// 添加id集合
	addIds := make([]int, 0, len(ipt.GrantIds))
	// 遍历新的集合
	for _, item := range ipt.GrantIds {
		// 如果没有在旧的集合里面,就添加
		if _, ok := oldIds[item]; !ok {
			addIds = append(addIds, item)
		}
	}
	// 删除id集合
	dropIds := make([]int, 0, 8)
	// 遍历旧的集合
	for _, item := range oldGrants {
		// 如果没在新集合里面,就添加
		if _, ok := newIds[item.Id]; !ok {
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
	// repo.RoleGrantAddMulti(addRoleGrants)
	// // 删除旧的角色授权集合
	// repo.RoleGrantDropMulti(dropIds)
	// 合并到一个事务
	err = repo.RoleGrantEdit(addRoleGrants, dropIds)
	if err != nil {
		return ctx.JSON(utils.Fail("修改失败", err.Error()))
	}
	return ctx.JSON(utils.Succ("succ"))
}
