import { Alova } from '@/utils/alova';
import { getPageList } from './common';
import type { Tag } from './types';

/**
 * 标签接口，对应后端 /api/tags
 */

/** 标签分页列表 */
export async function getTagList(params: Recordable = {}) {
  return getPageList<Tag>('/tags', params);
}

/** 标签详情 */
export function getTag(id: string) {
  return Alova.Get<Tag>(`/tags/${id}`);
}

/** 新增标签 */
export function addTag(params: Recordable) {
  return Alova.Post<string>('/tags', params);
}

/** 修改标签 */
export function editTag(id: string, params: Recordable) {
  return Alova.Put<null>(`/tags/${id}`, params);
}

/** 删除标签 */
export function deleteTag(id: string) {
  return Alova.Delete<null>(`/tags/${id}`);
}
