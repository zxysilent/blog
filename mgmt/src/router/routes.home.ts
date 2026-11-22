import { RouteRecordRaw } from 'vue-router';
import { Layout } from '@/router/constant';
import { HomeOutlined } from '@vicons/antd';
import { renderIcon } from '@/utils/index';

/**
 * 后台首页路由
 *
 * 只有一个子路由，按 isRootRouter 规则会渲染成一级菜单项（不出现折叠面板）。
 * 不声明 permissions：任何已登录账号都能看到，作为登录后的落地页；
 * 首页里的统计与最近操作再按各自权限标识决定是否请求与展示。
 */
const routes: Array<RouteRecordRaw> = [
  {
    path: '/home',
    name: 'Home',
    redirect: '/home/index',
    component: Layout,
    meta: {
      title: '后台首页',
      icon: renderIcon(HomeOutlined),
      sort: -1,
    },
    children: [
      {
        path: 'index',
        name: 'home_index',
        meta: {
          title: '后台首页',
          affix: true,
        },
        component: () => import('@/views/home/index.vue'),
      },
    ],
  },
];

export default routes;
