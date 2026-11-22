import { Alova } from '@/utils/alova';
import { getPageList } from './common';
import type { Log } from './types';

/**
 * 操作日志接口，对应后端 /api/logs
 * 日志由后端审计中间件自动写入，这里只提供查询与清理
 */

/** 日志分页列表，支持按模块/权限标识/方法/状态码/时间范围过滤 */
export async function getLogList(params: Recordable = {}) {
  return getPageList<Log>('/logs', params);
}

/** 日志详情 */
export function getLog(id: string) {
  return Alova.Get<Log>(`/logs/${id}`);
}

/** 清空日志 */
export function clearLogs() {
  return Alova.Delete<number>('/logs');
}

/** 删除单条日志 */
export function deleteLog(id: string) {
  return Alova.Delete<null>(`/logs/${id}`);
}
