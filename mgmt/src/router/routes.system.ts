import { RouteRecordRaw } from 'vue-router';
import { Layout } from '@/router/constant';
import { OptionsSharp } from '@vicons/ionicons5';
import { renderIcon } from '@/utils/index';

/**
 * 角色权限路由（用户 / 角色 / 权限标识，对应旧版的「角色权限」分组）
 * meta.permissions 写后端权限标识（与路由 Name 第一段一致），
 * 缺少对应权限时该菜单不会出现在侧边栏，也无法直接访问
 */
const routes: Array<RouteRecordRaw> = [
  {
    path: '/system',
    name: 'System',
    redirect: '/system/user',
    component: Layout,
    meta: {
      title: '角色权限',
      icon: renderIcon(OptionsSharp),
      sort: 1,
    },
    children: [
      {
        path: 'user',
        name: 'system_user',
        meta: {
          title: '用户管理',
          permissions: ['user:list'],
        },
        component: () => import('@/views/system/user/user.vue'),
      },
      {
        path: 'role',
        name: 'system_role',
        meta: {
          title: '角色权限',
          permissions: ['role:list'],
        },
        component: () => import('@/views/system/role/role.vue'),
      },
      {
        path: 'permission',
        name: 'system_permission',
        meta: {
          title: '权限标识',
          permissions: ['permission:list'],
        },
        component: () => import('@/views/system/permission/permission.vue'),
      },
    ],
  },
];

export default routes;
