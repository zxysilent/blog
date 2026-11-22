import { Alova } from '@/utils/alova';
import { asList } from './common';
import type { Captcha, LoginArgs, PermissionNode, Reply, User } from './types';

/**
 * 认证相关接口
 * 对应后端 /api/captcha、/api/tokens、/api/auth*
 */

/** 图形验证码：返回 base64 图片与服务端签名 */
export function getCaptcha() {
  return Alova.Get<Captcha>('/captcha');
}

/**
 * 登录-创建令牌
 * 后端校验验证码与管理员身份，成功时返回令牌字符串
 * 这里保留原生信封，登录页需要区分具体失败原因
 */
export function login(params: LoginArgs) {
  return Alova.Post<Reply<string>>('/tokens', params, {
    meta: { isTransformResponse: false },
  });
}

/** 当前登录用户 */
export function getUserInfo() {
  return Alova.Get<User>('/auth');
}

/** 我的权限标识集合，用于前端菜单与按钮控制 */
export async function getPermissions(): Promise<string[]> {
  const res = await Alova.Get<string[]>('/auth/permissions');
  return asList<string>(res);
}

/** 修改自己的资料 */
export function updateProfile(params: { name?: string; avatar?: string; phone?: string; email?: string }) {
  return Alova.Put<null>('/auth', params);
}

/** 修改自己的密码 */
export function changePassword(params: { opasswd: string; npasswd: string }) {
  return Alova.Put<null>('/auth/password', params);
}

/**
 * 校验当前账号密码，用于锁屏解锁这类二次确认场景
 * 保留原生信封，调用方据 code 判断是否通过
 */
export function verifyPasswd(params: { passwd: string }) {
  return Alova.Post<Reply<null>>('/auth/verify', params, {
    meta: { isTransformResponse: false },
  });
}

/** 权限标识树（只读，按分组聚合），角色授权时使用 */
export async function getPermissionTree(): Promise<PermissionNode[]> {
  const res = await Alova.Get<PermissionNode[]>('/permissions/tree');
  return asList<PermissionNode>(res);
}

export type { Captcha, LoginArgs, PermissionNode, User };
