import { Alova } from '@/utils/alova';
import { getPageList } from './common';
import type { Permission } from './types';

/**
 * 权限标识接口，对应后端 /api/permissions
 * 权限标识由后端路由注册自动同步，接口只读，不提供增删改
 */

/** 权限分页列表 */
export async function getPermissionList(params: Recordable = {}) {
  return getPageList<Permission>('/permissions', params);
}

/** 权限详情 */
export function getPermission(id: string) {
  return Alova.Get<Permission>(`/permissions/${id}`);
}
