package kernel

import (
	"blog/internal/model"
	"blog/internal/repo"
	"blog/pkg/reply"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/zxysilent/logs"
)

// UserList doc
// @Tags user
// @Summary 用户列表
// @Security KeyAuth
// @Param query query model.UserListArgs true "请求数据"
// @Success 200 {object} reply.Reply{data=[]model.User} "返回数据"
// @Router /api/users [get]
func UserList(ctx echo.Context) error {
	actor, err := currentRole(ctx)
	if err != nil {
		return reply.ErrToken(ctx, "请重新登陆", "角色不存在")
	}
	args := &model.UserListArgs{}
	if err := ctx.Bind(args); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	rctx := ctx.Request().Context()
	// 只能看到角色层级低于自己的用户，且看不到自己；根角色可看全部（自己除外）
	self, _ := currentUid(ctx)
	mods, count, err := repo.UserList(rctx, args, actor.Level, self, actor.Id == model.RoleRootId)
	if err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.UserList")
		return reply.Fail(ctx, "查询失败")
	}
	return reply.List(ctx, "ok", mods, count)
}

// UserGet doc
// @Tags user
// @Summary 用户详情
// @Security KeyAuth
// @Param id path string true "用户id"
// @Success 200 {object} reply.Reply{data=model.User} "返回数据"
// @Router /api/users/{id} [get]
func UserGet(ctx echo.Context) error {
	actor, err := currentRole(ctx)
	if err != nil {
		return reply.ErrToken(ctx, "请重新登陆", "角色不存在")
	}
	rctx := ctx.Request().Context()
	mod, err := repo.UserGet(rctx, ctx.Param("id"))
	if err != nil {
		logs.Ctx(rctx).Err(err).Str("id", ctx.Param("id")).Error("repo.UserGet")
		return reply.Fail(ctx, "用户不存在")
	}
	// 自己的资料不在用户管理里回显，与列表口径保持一致
	if selfId, ok := currentUid(ctx); ok && selfId == mod.Id {
		return reply.Fail(ctx, "查看自己的资料请使用个人设置")
	}
	// 只能查看层级低于自己的用户
	if !canView(actor, userRole(ctx, mod.RoleId)) {
		return reply.Fail(ctx, "无权查看该用户")
	}
	return reply.Succ(ctx, "ok", mod)
}

// UserAdd doc
// @Tags user
// @Summary 用户新增
// @Security KeyAuth
// @Param body body model.UserAddArgs true "请求数据"
// @Success 200 {object} reply.Reply{data=string} "返回数据"
// @Router /api/users [post]
func UserAdd(ctx echo.Context) error {
	ipt := &model.UserAddArgs{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	ipt.Num = strings.TrimSpace(ipt.Num)
	if len(ipt.Num) < 3 || len(ipt.Num) > 64 {
		return reply.Fail(ctx, "账号长度应在 3~64 之间")
	}
	if ipt.Passwd == "" {
		return reply.Fail(ctx, "请输入密码")
	}
	rctx := ctx.Request().Context()
	if repo.UserExistNum(rctx, ipt.Num) {
		return reply.Fail(ctx, "账号已存在")
	}
	actor, err := currentRole(ctx)
	if err != nil {
		return reply.ErrToken(ctx, "请重新登陆", "角色不存在")
	}
	// 允许不分配角色：这类账号登录后没有任何权限，
	// 需要授权时再给角色即可（不再有「用户类型」这一层门槛）
	if ipt.RoleId != "" {
		role, err := repo.RoleGet(rctx, ipt.RoleId)
		if err != nil {
			return reply.Fail(ctx, "角色不存在")
		}
		// 只能分配层级低于自己的角色
		if !canManage(actor, role) {
			return reply.Fail(ctx, "只能分配层级低于自己的角色")
		}
	}
	if ipt.Status == 0 {
		ipt.Status = model.UserStatusNormal
	}
	openid := model.UserNormalizeOpenid(ipt.Openid)
	if !model.UserOpenidValid(openid) {
		return reply.Fail(ctx, "openid 只能包含字母、数字、下划线和短横线，且不超过 128 字符")
	}
	if repo.UserExistOpenid(rctx, openid, "") {
		return reply.Fail(ctx, "openid 已被其它账号绑定")
	}
	milli := time.Now().UnixMilli()
	mod := &model.User{
		Num:     ipt.Num,
		Name:    strings.TrimSpace(ipt.Name),
		Avatar:  ipt.Avatar,
		Phone:   strings.TrimSpace(ipt.Phone),
		Email:   strings.TrimSpace(ipt.Email),
		Openid:  openid,
		RoleId:  ipt.RoleId,
		Status:  ipt.Status,
		Updated: milli,
		Created: milli,
	}
	if err := mod.SetPasswd(ipt.Passwd); err != nil {
		return reply.Fail(ctx, "密码加密失败")
	}
	if err := repo.UserAdd(rctx, mod); err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.UserAdd")
		return reply.Fail(ctx, "添加失败")
	}
	return reply.Succ(ctx, "添加成功", mod.Id)
}

// UserEdit doc
// @Tags user
// @Summary 用户修改
// @Security KeyAuth
// @Param id path string true "用户id"
// @Param body body model.UserEditArgs true "请求数据"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/users/{id} [put]
func UserEdit(ctx echo.Context) error {
	ipt := &model.UserEditArgs{}
	if err := ctx.Bind(ipt); err != nil {
		return reply.Fail(ctx, "输入有误", err.Error())
	}
	rctx := ctx.Request().Context()
	mod, err := repo.UserGet(rctx, ctx.Param("id"))
	if err != nil {
		return reply.Fail(ctx, "用户不存在")
	}
	actor, err := currentRole(ctx)
	if err != nil {
		return reply.ErrToken(ctx, "请重新登陆", "角色不存在")
	}
	// 只能修改层级低于自己的用户（自己的账号走 /auth 接口）
	uid, _ := currentUid(ctx)
	if uid != mod.Id && !canManage(actor, userRole(ctx, mod.RoleId)) {
		return reply.Fail(ctx, "只能修改层级低于自己的用户")
	}
	cols := make([]string, 0, 8)
	if ipt.Num != nil && strings.TrimSpace(*ipt.Num) != mod.Num {
		num := strings.TrimSpace(*ipt.Num)
		if len(num) < 3 || len(num) > 64 {
			return reply.Fail(ctx, "账号长度应在 3~64 之间")
		}
		if repo.UserExistNum(rctx, num) {
			return reply.Fail(ctx, "账号已存在")
		}
		mod.Num = num
		cols = append(cols, "num")
	}
	if ipt.Name != nil {
		mod.Name = strings.TrimSpace(*ipt.Name)
		cols = append(cols, "name")
	}
	if ipt.Avatar != nil {
		mod.Avatar = *ipt.Avatar
		cols = append(cols, "avatar")
	}
	if ipt.Phone != nil {
		mod.Phone = strings.TrimSpace(*ipt.Phone)
		cols = append(cols, "phone")
	}
	if ipt.Email != nil {
		mod.Email = strings.TrimSpace(*ipt.Email)
		cols = append(cols, "email")
	}
	if ipt.Openid != nil {
		openid := model.UserNormalizeOpenid(*ipt.Openid)
		if !model.UserOpenidValid(openid) {
			return reply.Fail(ctx, "openid 只能包含字母、数字、下划线和短横线，且不超过 128 字符")
		}
		// 空串表示解绑，不参与唯一性判定
		if repo.UserExistOpenid(rctx, openid, mod.Id) {
			return reply.Fail(ctx, "openid 已被其它账号绑定")
		}
		mod.Openid = openid
		cols = append(cols, "openid")
	}
	if ipt.RoleId != nil {
		roleId := strings.TrimSpace(*ipt.RoleId)
		// 传空串表示解除角色绑定，账号仍在但没有任何权限
		if roleId != "" {
			role, err := repo.RoleGet(rctx, roleId)
			if err != nil {
				return reply.Fail(ctx, "角色不存在")
			}
			// 只能分配层级低于自己的角色
			if !canManage(actor, role) {
				return reply.Fail(ctx, "只能分配层级低于自己的角色")
			}
		}
		mod.RoleId = roleId
		cols = append(cols, "role_id")
	}
	if ipt.Status != nil {
		mod.Status = *ipt.Status
		cols = append(cols, "status")
	}
	if ipt.Passwd != nil && *ipt.Passwd != "" {
		if err := mod.SetPasswd(*ipt.Passwd); err != nil {
			return reply.Fail(ctx, "密码加密失败")
		}
		cols = append(cols, "passwd")
	}
	if len(cols) < 1 {
		return reply.Fail(ctx, "没有需要修改的数据")
	}
	mod.Updated = time.Now().UnixMilli()
	cols = append(cols, "updated")
	if err := repo.UserEdit(rctx, mod, cols...); err != nil {
		logs.Ctx(rctx).Err(err).Error("repo.UserEdit")
		return reply.Fail(ctx, "修改失败")
	}
	return reply.Succ(ctx, "修改成功")
}

// UserDrop doc
// @Tags user
// @Summary 用户删除
// @Security KeyAuth
// @Param id path string true "用户id"
// @Success 200 {object} reply.Reply "返回数据"
// @Router /api/users/{id} [delete]
func UserDrop(ctx echo.Context) error {
	id := ctx.Param("id")
	uid, ok := currentUid(ctx)
	if !ok {
		return reply.ErrToken(ctx, "请重新登陆")
	}
	if id == uid {
		return reply.Fail(ctx, "不能删除自己")
	}
	if id == "" {
		return reply.Fail(ctx, "参数有误")
	}
	actor, err := currentRole(ctx)
	if err != nil {
		return reply.ErrToken(ctx, "请重新登陆", "角色不存在")
	}
	rctx := ctx.Request().Context()
	mod, err := repo.UserGet(rctx, id)
	if err != nil {
		return reply.Fail(ctx, "用户不存在")
	}
	// 只能删除层级低于自己的用户
	if !canManage(actor, userRole(ctx, mod.RoleId)) {
		return reply.Fail(ctx, "只能删除层级低于自己的用户")
	}
	if err := repo.UserDrop(rctx, id); err != nil {
		logs.Ctx(rctx).Err(err).Str("id", id).Error("repo.UserDrop")
		return reply.Fail(ctx, "删除失败")
	}
	return reply.Succ(ctx, "删除成功")
}
