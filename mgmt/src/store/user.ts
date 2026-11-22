import { defineStore } from 'pinia';
import { store } from '@/store';
import { ACCESS_TOKEN, CURRENT_USER, TABS_ROUTES } from '@/store/mutation-types';
import { ResultEnum } from '@/enums/httpEnum';

import { getPermissions, getUserInfo, login } from '@/api/auth';
import type { LoginArgs, User } from '@/api/types';
import { storage } from '@/utils/Storage';

export type UserInfoType = User | Recordable;

export interface IUserState {
  token: string;
  /** 展示用名称 */
  username: string;
  avatar: string;
  /** 权限标识集合，形如 ['user:list', 'role:list'] */
  permissions: string[];
  info: UserInfoType;
}

export const useUserStore = defineStore({
  id: 'app-user',
  state: (): IUserState => ({
    token: storage.get(ACCESS_TOKEN, ''),
    username: '',
    avatar: '',
    permissions: [],
    info: storage.get(CURRENT_USER, {}),
  }),
  getters: {
    getToken(): string {
      return this.token;
    },
    getAvatar(): string {
      return this.avatar;
    },
    getNickname(): string {
      return this.username;
    },
    getPermissions(): string[] {
      return this.permissions;
    },
    getUserInfo(): UserInfoType {
      return this.info;
    },
  },
  actions: {
    setToken(token: string) {
      this.token = token;
    },
    setPermissions(permissions: string[]) {
      this.permissions = permissions;
    },
    setUserInfo(info: UserInfoType) {
      this.info = info;
      this.username = (info as User)?.name || (info as User)?.num || '';
      this.avatar = (info as User)?.avatar || '';
    },
    /**
     * 登录
     * 后端 /api/tokens 成功时 data 即令牌字符串，失败时返回业务码与提示
     */
    async login(params: LoginArgs) {
      const res = await login(params);
      if (res.code === ResultEnum.SUCCESS && res.data) {
        const token = res.data;
        const ex = 7 * 24 * 60 * 60;
        storage.set(ACCESS_TOKEN, token, ex);
        // 清掉上一个账号残留的标签页，避免切号后看到无权访问的页面入口
        storage.remove(TABS_ROUTES);
        this.setToken(token);
      }
      return res;
    },

    /**
     * 获取登录用户信息与权限标识
     * 权限是后端按角色下发的权限标识集合，前端据此渲染菜单与按钮
     */
    async getInfo() {
      const [info, permissions] = await Promise.all([getUserInfo(), getPermissions()]);
      if (!info) {
        throw new Error('获取用户信息失败');
      }
      this.setUserInfo(info);
      this.setPermissions(permissions || []);
      const ex = 7 * 24 * 60 * 60;
      storage.set(CURRENT_USER, info, ex);
      return { ...info, permissions: permissions || [] };
    },

    /** 登出：清空本地状态（后端令牌无状态，无需通知服务端） */
    async logout() {
      this.setPermissions([]);
      this.setUserInfo({});
      this.username = '';
      this.avatar = '';
      this.token = '';
      storage.remove(ACCESS_TOKEN);
      storage.remove(CURRENT_USER);
    },
  },
});

// Need to be used outside the setup
export function useUser() {
  return useUserStore(store);
}
