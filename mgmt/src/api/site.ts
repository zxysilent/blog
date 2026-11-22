import { Alova } from '@/utils/alova';
import type { Site } from './types';

/** 站点信息接口，对应后端 /api/site 与 /api/status */

/** 站点公开配置：基础信息 + 博客运行参数，匿名可读 */
export function getSite() {
  return Alova.Get<Site>('/site');
}

/** 服务器状态（运行时、内存等），用于首页概览 */
export function getStatus() {
  return Alova.Get<Record<string, any>>('/status');
}

/** 上传文件，返回可访问地址 */
export function uploadFile(file: File) {
  const form = new FormData();
  form.append('file', file);
  return Alova.Post<string>('/files', form);
}

/** 上传图片，返回可访问地址 */
export function uploadImage(file: File) {
  const form = new FormData();
  form.append('file', file);
  return Alova.Post<string>('/images', form);
}
