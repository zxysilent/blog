import { RouteRecordRaw } from 'vue-router';
import { Layout } from '@/router/constant';
import { DocumentTextOutline } from '@vicons/ionicons5';
import { renderIcon } from '@/utils/index';

/**
 * 内容管理路由
 * meta.permissions 写后端权限标识（与路由 Name 第一段一致），
 * 缺少对应权限时该菜单不会出现在侧边栏，也无法直接访问
 */
const routes: Array<RouteRecordRaw> = [
  {
    path: '/content',
    name: 'Content',
    redirect: '/content/post',
    component: Layout,
    meta: {
      title: '内容管理',
      icon: renderIcon(DocumentTextOutline),
      sort: 0,
    },
    children: [
      {
        path: 'post',
        name: 'content_post',
        meta: {
          title: '博文管理',
          permissions: ['post:list'],
        },
        component: () => import('@/views/content/post/post.vue'),
      },
      {
        path: 'post/add',
        name: 'content_post_add',
        meta: {
          title: '新增博文',
          // 表单页不进侧边栏，只从列表页跳转过来
          hidden: true,
          permissions: ['post:create'],
        },
        component: () => import('@/views/content/post/form.vue'),
      },
      {
        path: 'post/:id/edit',
        name: 'content_post_edit',
        meta: {
          title: '编辑博文',
          hidden: true,
          permissions: ['post:update'],
        },
        component: () => import('@/views/content/post/form.vue'),
      },
      {
        path: 'post/:id',
        name: 'content_post_detail',
        meta: {
          title: '内容预览',
          // 只读预览页，不进侧边栏，从列表页跳转
          hidden: true,
          permissions: ['post:read'],
        },
        component: () => import('@/views/content/post/detail.vue'),
      },
      {
        path: 'cate',
        name: 'content_cate',
        meta: {
          title: '分类管理',
          permissions: ['cate:list'],
        },
        component: () => import('@/views/content/cate/cate.vue'),
      },
      {
        path: 'tag',
        name: 'content_tag',
        meta: {
          title: '标签管理',
          permissions: ['tag:list'],
        },
        component: () => import('@/views/content/tag/tag.vue'),
      },
    ],
  },
];

export default routes;
