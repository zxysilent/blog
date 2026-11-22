import { RouteRecordRaw } from 'vue-router';
import { Layout } from '@/router/constant';
import { SettingsOutline } from '@vicons/ionicons5';
import { renderIcon } from '@/utils/index';

/**
 * 系统设置路由（导航顺序参考旧版 vue：后台首页 → 内容管理 → 角色权限 → 系统设置）
 * meta.permissions 写后端权限标识，缺少权限时菜单不出现、也无法直接访问；
 * 个人设置不声明权限，任何已登录账号都能访问。
 */
const routes: Array<RouteRecordRaw> = [
  {
    path: '/setting',
    name: 'Setting',
    redirect: '/setting/profile',
    component: Layout,
    meta: {
      title: '系统设置',
      icon: renderIcon(SettingsOutline),
      sort: 2,
    },
    children: [
      {
        path: 'profile',
        name: 'setting_profile',
        meta: {
          title: '个人设置',
          affix: true,
        },
        component: () => import('@/views/setting/profile.vue'),
      },
      {
        // 内置配置的表单化编辑（站点/登录/博文/认证/上传/SQL）
        // 与「配置管理」（清单与自建配置）分开，职责更单一
        path: 'global',
        name: 'setting_global',
        meta: {
          title: '全局配置',
          permissions: ['config:list'],
        },
        component: () => import('@/views/setting/global/index.vue'),
      },
      {
        path: 'config',
        name: 'setting_config',
        meta: {
          title: '配置管理',
          permissions: ['config:list'],
        },
        component: () => import('@/views/setting/config/config.vue'),
      },
      {
        path: 'dict',
        name: 'setting_dict',
        meta: {
          title: '配置字典',
          permissions: ['dict:list'],
        },
        component: () => import('@/views/setting/dict/dict.vue'),
      },
      {
        // 审计日志：只读查询 + 清理，写入由后端中间件自动完成
        path: 'log',
        name: 'setting_log',
        meta: {
          title: '操作日志',
          permissions: ['log:list'],
        },
        component: () => import('@/views/setting/log/log.vue'),
      },
    ],
  },
];

export default routes;
