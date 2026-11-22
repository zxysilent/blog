import { toRaw, unref } from 'vue';
import { defineStore } from 'pinia';
import { RouteRecordRaw } from 'vue-router';
import { store } from '@/store';
import { asyncRoutes, constantRouter } from '@/router/index';
import { useProjectSetting } from '@/hooks/useProjectSetting';

interface TreeHelperConfig {
  id: string;
  children: string;
  pid: string;
}

const DEFAULT_CONFIG: TreeHelperConfig = {
  id: 'id',
  children: 'children',
  pid: 'pid',
};

const getConfig = (config: Partial<TreeHelperConfig>) => Object.assign({}, DEFAULT_CONFIG, config);

export interface IAsyncRouteState {
  menus: RouteRecordRaw[];
  routers: any[];
  routersAdded: any[];
  keepAliveComponents: string[];
  isDynamicRouteAdded: boolean;
}

function filter<T = any>(
  tree: T[],
  func: (n: T) => boolean,
  config: Partial<TreeHelperConfig> = {}
): T[] {
  config = getConfig(config);
  const children = config.children as string;

  function listFilter(list: T[]) {
    return list
      .map((node: any) => ({ ...node }))
      .filter((node) => {
        node[children] = node[children] && listFilter(node[children]);
        return func(node) || (node[children] && node[children].length);
      });
  }

  return listFilter(tree);
}

export const useAsyncRouteStore = defineStore({
  id: 'app-async-route',
  state: (): IAsyncRouteState => ({
    menus: [],
    routers: constantRouter,
    routersAdded: [],
    keepAliveComponents: [],
    // Whether the route has been dynamically added
    isDynamicRouteAdded: false,
  }),
  getters: {
    getMenus(): RouteRecordRaw[] {
      return this.menus;
    },
    getIsDynamicRouteAdded(): boolean {
      return this.isDynamicRouteAdded;
    },
  },
  actions: {
    getRouters() {
      return toRaw(this.routersAdded);
    },
    setDynamicRouteAdded(added: boolean) {
      this.isDynamicRouteAdded = added;
    },
    // 设置动态路由
    setRouters(routers: RouteRecordRaw[]) {
      this.routersAdded = routers;
      this.routers = constantRouter.concat(routers);
    },
    setMenus(menus: RouteRecordRaw[]) {
      // 设置动态路由
      this.menus = menus;
    },
    setKeepAliveComponents(compNames: string[]) {
      // 设置需要缓存的组件
      this.keepAliveComponents = compNames;
    },
    async generateRoutes(data) {
      // 权限来自后端 /api/auth/permissions，形如 ['user:list', 'role:list']
      // 这里统一转成字符串集合，路由 meta.permissions 也写字符串，两边对得上即可
      const permissionsList: string[] = (data?.permissions ?? []).map((item) =>
        typeof item === 'string' ? item : item?.value
      );
      const routeFilter = (route) => {
        const { meta } = route;
        const { permissions } = meta || {};
        // 未声明权限的页面默认可见
        if (!permissions) return true;
        return permissions.some((item) => permissionsList.includes(item));
      };
      const { permissionMode } = useProjectSetting();
      if (unref(permissionMode) === 'BACK') {
        // 后端未提供菜单接口（菜单由前端静态定义），BACK 模式退化为全量路由
        console.warn('permissionMode=BACK 已不再支持，已按 FIXED 处理');
      }
      // 过滤账户是否拥有某一个权限，并将菜单从加载列表移除
      const accessedRouters = filter(asyncRoutes, routeFilter)
        .filter(routeFilter)
        // 父级菜单本身不声明权限，子级全被过滤掉时会剩下一个空菜单，这里一并剪掉
        .filter((route: any) => !Array.isArray(route.children) || route.children.length > 0);
      this.setRouters(accessedRouters);
      this.setMenus(accessedRouters);
      return toRaw(accessedRouters);
    },
  },
});

// Need to be used outside the setup
export function useAsyncRoute() {
  return useAsyncRouteStore(store);
}
