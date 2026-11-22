import { Alova } from '@/utils/alova';
import { asList, getPageList } from './common';
import type { Config, ConfigGroup } from './types';

/**
 * 配置接口，对应后端 /api/configs
 * 内置配置（auth/basic/blog/login/upload/sql/layout）不可删除，值可在线修改并立即生效
 *
 * 内置配置都有独立路径与独立权限点（config:basic 等），
 * 便于按配置分权；自建配置走通用的 /configs/:id。
 */

/** 配置分页列表，支持按分组/名称/类型过滤 */
export async function getConfigList(params: Recordable = {}) {
  return getPageList<Config>('/configs', params);
}

/** 配置分组及数量，供分组标签页展示 */
export async function getConfigGroups(): Promise<ConfigGroup[]> {
  const res = await Alova.Get<ConfigGroup[]>('/configs/groups');
  return asList<ConfigGroup>(res);
}

/**
 * 配置详情
 * 内置配置命中后端的独立路由（如 /configs/basic），需要对应 config:<id> 权限
 */
export function getConfigById(id: string) {
  return Alova.Get<Config>(`/configs/${id}`);
}

/** 配置详情（保持旧名，等价于 getConfigById） */
export function getConfig(id: string) {
  return getConfigById(id);
}

/** 保存配置值，内置配置同样走独立路径 */
export function saveConfigValue(id: string, value: string) {
  return Alova.Put<null>(`/configs/${id}`, { value });
}

/** 新增配置，id 留空由后端生成 */
export function addConfig(params: Recordable) {
  return Alova.Post<string>('/configs', params);
}

/** 修改配置：值、描述、分组 */
export function editConfig(id: string, params: Recordable) {
  return Alova.Put<null>(`/configs/${id}`, params);
}

/** 删除配置 */
export function deleteConfig(id: string) {
  return Alova.Delete<null>(`/configs/${id}`);
}

/**
 * 切换登录页风格
 * 单独成接口：风格改动需要 login:style 权限，与登录页配置修改（config:login）分开
 */
export function setLoginStyle(style: 'mgmt' | 'vue') {
  return Alova.Put<null>('/configs/login/style', { style });
}

/** 站点公开配置（匿名可读），前台渲染用 */
export function getSiteConfig() {
  return Alova.Get<Record<string, any>>('/site');
}
