import { Alova } from '@/utils/alova';
import { asList, getPageList } from './common';
import type { Dict, DictEntry } from './types';

/**
 * 字典接口，对应后端 /api/dicts
 * 字典主键即业务键（如 gender、page_size），创建时可自定义
 */

/** 字典分页列表 */
export async function getDictList(params: Recordable = {}) {
  return getPageList<Dict>('/dicts', params);
}

/** 字典详情（含选项树） */
export function getDict(id: string) {
  return Alova.Get<Dict>(`/dicts/${id}`);
}

/** 新增字典（含选项） */
export function addDict(params: Recordable) {
  return Alova.Post<string>('/dicts', params);
}

/** 修改字典，entries 传 null 表示不动选项，传 [] 表示清空 */
export function editDict(id: string, params: Recordable) {
  return Alova.Put<null>(`/dicts/${id}`, params);
}

/** 删除字典（级联删除选项） */
export function deleteDict(id: string) {
  return Alova.Delete<null>(`/dicts/${id}`);
}

/**
 * 字典选项（匿名可读，已过滤禁用项）
 * 用于下拉框等场景，不参与权限校验
 */
export async function getDictEntries(id: string): Promise<DictEntry[]> {
  const res = await Alova.Get<DictEntry[]>(`/dicts/${id}/entries`);
  return asList<DictEntry>(res);
}

/** 字典值转下拉选项 */
export async function getDictOptions(id: string) {
  const rows = await getDictEntries(id);
  return (rows || []).map((it) => ({ label: it.label, value: it.value }));
}
