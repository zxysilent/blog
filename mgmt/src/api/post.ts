import { Alova } from '@/utils/alova';
import { asList, getPageList, toPageArgs } from './common';
import type { Post, Tag } from './types';

/**
 * 内容接口，对应后端 /api/posts
 * kind：1 文章、2 页面、3 笔记（笔记用 pid 自关联成目录树）
 */

/** 内容分页列表 */
export async function getPostList(params: Recordable = {}) {
  return getPageList<Post>('/posts', params);
}

/** 内容树（笔记目录） */
export async function getPostTree(params: Recordable = {}): Promise<Post[]> {
  const res = await Alova.Get<Post[]>('/posts/tree', { params: toPageArgs(params) });
  return asList<Post>(res);
}

/** 内容详情（含 markdown 与渲染后的富文本） */
export function getPost(id: string) {
  return Alova.Get<Post>(`/posts/${id}`);
}

/** 新增内容 */
export function addPost(params: Recordable) {
  return Alova.Post<string>('/posts', params);
}

/** 修改内容，支持单字段局部更新 */
export function editPost(id: string, params: Recordable) {
  return Alova.Put<null>(`/posts/${id}`, params);
}

/** 置顶/取消置顶 */
export function setPostTop(id: string, istop: boolean) {
  return Alova.Put<null>(`/posts/${id}/top`, { istop });
}

/** 删除内容（同时清理标签关联） */
export function deletePost(id: string) {
  return Alova.Delete<null>(`/posts/${id}`);
}

/** 内容标签 */
export async function getPostTags(id: string): Promise<Tag[]> {
  const res = await Alova.Get<Tag[]>(`/posts/${id}/tags`);
  return asList<Tag>(res);
}

/** 全量替换内容标签 */
export function setPostTags(id: string, tagIds: string[]) {
  return Alova.Put<null>(`/posts/${id}/tags`, { tag_ids: tagIds });
}

/** 上下篇 */
export function getPostNeighbors(id: string) {
  return Alova.Get<{ prev: Post; next: Post }>(`/posts/${id}/neighbors`);
}
