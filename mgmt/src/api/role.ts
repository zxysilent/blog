import { Alova } from '@/utils/alova';
import { asList, getPageList } from './common';
import type { Permission, Role } from './types';

/**
 * 角色管理接口，对应后端 /api/roles
 * 层级规则：数值越小权限越高，只能操作数值更大（权限更低）的角色
 */

/** 角色分页列表 */
export async function getRoleList(params: Recordable = {}) {
  return getPageList<Role>('/roles', params);
}

/** 角色详情 */
export function getRole(id: string) {
  return Alova.Get<Role>(`/roles/${id}`);
}

/** 新增角色 */
export function addRole(params: Recordable) {
  return Alova.Post<string>('/roles', params);
}

/** 修改角色 */
export function editRole(id: string, params: Recordable) {
  return Alova.Put<null>(`/roles/${id}`, params);
}

/** 删除角色 */
export function deleteRole(id: string) {
  return Alova.Delete<null>(`/roles/${id}`);
}

/** 角色已授权的权限详情 */
export async function getRolePermissions(id: string): Promise<Permission[]> {
  const res = await Alova.Get<Permission[]>(`/roles/${id}/permissions`);
  return asList<Permission>(res);
}

/** 角色授权-全量替换 */
export function setRolePermissions(id: string, permissionIds: string[]) {
  return Alova.Put<null>(`/roles/${id}/permissions`, { permission_ids: permissionIds });
}
