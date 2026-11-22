import { Alova } from '@/utils/alova';
import { getPageList } from './common';
import type { Role, User } from './types';

/**
 * 用户管理接口，对应后端 /api/users
 * 列表只返回角色层级不高于自己的用户（后端按当前登录角色过滤）
 */

/** 用户分页列表 */
export async function getUserList(params: Recordable = {}) {
  return getPageList<User>('/users', params);
}

/** 用户详情 */
export function getUser(id: string) {
  return Alova.Get<User>(`/users/${id}`);
}

/** 新增用户 */
export function addUser(params: Recordable) {
  return Alova.Post<string>('/users', params);
}

/** 修改用户，支持单字段局部更新 */
export function editUser(id: string, params: Recordable) {
  return Alova.Put<null>(`/users/${id}`, params);
}

/** 删除用户 */
export function deleteUser(id: string) {
  return Alova.Delete<null>(`/users/${id}`);
}

/** 角色下拉用：可分配的角色列表（不分页） */
export async function getRoleOptions(): Promise<Role[]> {
  const res: any = await Alova.Get('/roles', { params: { page: 1, size: 1000 } });
  return Array.isArray(res?.data) ? res.data : [];
}
