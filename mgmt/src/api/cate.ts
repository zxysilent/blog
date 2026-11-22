import { Alova } from '@/utils/alova';
import { asList, getPageList, toPageArgs } from './common';
import type { Cate } from './types';

/**
 * 栏目接口，对应后端 /api/cates
 * 层级由服务端按父链自动维护，最多 3 级
 */

/** 栏目分页列表 */
export async function getCateList(params: Recordable = {}) {
  return getPageList<Cate>('/cates', params);
}

/** 栏目树，供下拉与筛选使用 */
export async function getCateTree(params: Recordable = {}): Promise<Cate[]> {
  const res = await Alova.Get<Cate[]>('/cates/tree', { params: toPageArgs(params) });
  return asList<Cate>(res);
}

/** 栏目全量列表（带内容计数） */
export async function getCateAll(): Promise<Cate[]> {
  const res = await getPageList<Cate>('/cates', { page: 1, size: 1000 });
  return res.list;
}

/** 栏目详情 */
export function getCate(id: string) {
  return Alova.Get<Cate>(`/cates/${id}`);
}

/** 新增栏目 */
export function addCate(params: Recordable) {
  return Alova.Post<string>('/cates', params);
}

/** 修改栏目 */
export function editCate(id: string, params: Recordable) {
  return Alova.Put<null>(`/cates/${id}`, params);
}

/** 删除栏目 */
export function deleteCate(id: string) {
  return Alova.Delete<null>(`/cates/${id}`);
}
