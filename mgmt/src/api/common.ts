import { Alova } from '@/utils/alova';
import type { PageArgs, PageResult, Reply } from './types';

/**
 * 后端分页参数换算
 * 表格组件用 page/pageSize，后端用 page/size，这里统一转换，
 * 并把常见的过滤字段原样透传（后端字段名与前端一致）
 *
 * 注意：表格的 reload 若被直接当作事件处理器绑定（`@click="reload"`），
 * 参数位会拿到 MouseEvent，展开后会污染查询参数（isTrusted/currentTarget 等），
 * 所以这里只放行字符串、数字、布尔三类原始值。
 */
export function toPageArgs(params: Recordable = {}): PageArgs {
  const safe = params || {};
  const page = Number(safe.page) > 0 ? Number(safe.page) : 1;
  const size = Number(safe.pageSize) > 0 ? Number(safe.pageSize) : 20;
  const args: PageArgs = { page, size };
  Object.keys(safe).forEach((key) => {
    if (key === 'page' || key === 'pageSize') return;
    const val = safe[key];
    // 空串与 undefined 不参与过滤，避免后端把它们当成有效条件
    if (val === undefined || val === null || val === '') return;
    const type = typeof val;
    if (type !== 'string' && type !== 'number' && type !== 'boolean') return;
    args[key] = val;
  });
  return args;
}

/**
 * 把后端 { data, count } 组装成表格需要的分页结果
 * 后端不返回页码，这里按请求参数补上，便于表格计算总页数
 */
export function toPageResult<T>(res: { data?: T[]; count?: number }, params: PageArgs = {}): PageResult<T> {
  const list = Array.isArray(res?.data) ? (res.data as T[]) : [];
  const itemCount = typeof res?.count === 'number' ? res.count : list.length;
  const size = params.size && params.size > 0 ? params.size : 20;
  const page = params.page && params.page > 0 ? params.page : 1;
  return {
    list,
    page,
    itemCount,
    pageCount: itemCount > 0 ? Math.ceil(itemCount / size) : 0,
  };
}

/**
 * 分页列表请求
 * 响应拦截器把分页结果解成 { data, count }，这里统一组装成表格需要的结构，
 * 避免每个接口都写一遍临时类型
 */
export async function getPageList<T>(url: string, params: Recordable = {}): Promise<PageResult<T>> {
  const args = toPageArgs(params);
  const res: any = await Alova.Get(url, { params: args });
  return toPageResult<T>(res, args);
}

/**
 * 取出数组结果
 * 响应拦截器已把信封拆开：非分页接口直接给 data（数组），分页接口给 { data, count }。
 * 这里统一成数组，调用方不必关心接口是哪种形态。
 */
export function asList<T>(res: any): T[] {
  if (Array.isArray(res)) return res as T[];
  if (Array.isArray(res?.data)) return res.data as T[];
  return [];
}

/** 判断业务码是否成功 */
export function isOk(res: Reply | undefined): boolean {
  return !!res && res.code === 200;
}
/** 毫秒时间戳转 yyyy-MM-dd HH:mm:ss，用于表格展示 */
export function fmtTime(ms?: number): string {
  if (!ms) return '-';
  const d = new Date(ms);
  const pad = (n: number) => String(n).padStart(2, '0');
  return `${d.getFullYear()}-${pad(d.getMonth() + 1)}-${pad(d.getDate())} ${pad(d.getHours())}:${pad(
    d.getMinutes()
  )}:${pad(d.getSeconds())}`;
}
