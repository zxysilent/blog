import { App } from 'vue';
import { createRouter, createWebHashHistory, RouteRecordRaw } from 'vue-router';
import { RedirectRoute } from '@/router/base';
import { PageEnum } from '@/enums/pageEnum';
import { createRouterGuards } from './guards';
import type { IModuleType } from './types';

const modules = import.meta.glob<IModuleType>('./routes.*.ts', { eager: true });

/**
 * 已启用的路由模块白名单（对应 routes.<name>.ts）
 *
 * 项目主线与旧版 vue 一致：后台首页 → 内容管理（博文/分类/标签）→ 系统管理（用户/角色/权限）→ 系统设置（个人/配置/字典/日志）；
 * 模板自带的演示页（仪表盘、组件示例、表单示例、iframe、多级菜单、结果页等）
 * 先隐藏起来，避免菜单冗余；文件全部保留，后续要用时把名字加回数组即可。
 */
const enabledModules = ['home', 'content', 'setting', 'system'];

/** 从 './routes.system.ts' 取出 'system' */
function moduleName(key: string): string {
  return key.replace(/^\.\/routes\./, '').replace(/\.ts$/, '');
}

const routeModuleList: RouteRecordRaw[] = Object.keys(modules)
  .filter((key) => enabledModules.includes(moduleName(key)))
  .reduce((list, key) => {
    const mod = modules[key].default ?? {};
    const modList = Array.isArray(mod) ? [...mod] : [mod];
    return [...list, ...modList];
  }, []);

function sortRoute(a, b) {
  return (a.meta?.sort ?? 0) - (b.meta?.sort ?? 0);
}

routeModuleList.sort(sortRoute);

export const RootRoute: RouteRecordRaw = {
  path: '/',
  name: 'Root',
  redirect: PageEnum.BASE_HOME,
  meta: {
    title: 'Root',
  },
};

export const LoginRoute: RouteRecordRaw = {
  path: '/login',
  name: 'Login',
  component: () => import('@/views/login/index.vue'),
  meta: {
    title: '登录',
  },
};

//需要验证权限
export const asyncRoutes = [...routeModuleList];

//普通路由 无需验证权限
export const constantRouter: RouteRecordRaw[] = [LoginRoute, RootRoute, RedirectRoute];

const router = createRouter({
  // hash 模式：地址形如 /#/system/user（部署在 /mgmt/ 下则为 /mgmt/#/system/user）
  // 静态托管不需要再配 history fallback，刷新任意子页面都只会请求 index.html
  //
  // 这里刻意不传 BASE_URL：createWebHashHistory 的默认行为是取当前 pathname 作为 base
  // （开发环境 /、生产环境 /mgmt/），正好匹配本项目两种部署方式；
  // 若传 BASE_URL(='/')，部署到 /mgmt/ 后链接会变成 /#/system/user，直接跳出 SPA。
  history: createWebHashHistory(),
  routes: constantRouter,
  strict: true,
  scrollBehavior: () => ({ left: 0, top: 0 }),
});

export function setupRouter(app: App) {
  app.use(router);
  // 创建路由守卫
  createRouterGuards(router);
}

export default router;
