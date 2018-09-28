import Main from "@/views/Main.vue";

// 不作为Main组件的子页面展示的页面单独写，如下
export const loginRouter = {
  path: "/login",
  name: "login",
  meta: {
    title: "Login - 登录"
  },
  component: () => import("@/views/login.vue")
};

export const page404 = {
  path: "/*",
  name: "error-404",
  meta: {
    title: "404-页面不存在"
  },
  component: () => import("@/views/error-page/404.vue")
};

export const page403 = {
  path: "/403",
  meta: {
    title: "403-权限不足"
  },
  name: "error-403",
  component: () => import("@//views/error-page/403.vue")
};

export const page500 = {
  path: "/500",
  meta: {
    title: "500-服务端错误"
  },
  name: "error-500",
  component: () => import("@/views/error-page/500.vue")
};

export const preview = {
  path: "/preview",
  name: "preview",
  component: () => import("@/views/form/article-publish/preview.vue")
};

export const locking = {
  path: "/locking",
  name: "locking",
  component: () =>
    import("@/views/main-components/lockscreen/components/locking-page.vue")
};

// 作为Main组件的子页面展示但是不在左侧菜单显示的路由写在otherRouter里
export const otherRouter = {
  path: "/",
  name: "otherRouter",
  redirect: "/home",
  component: Main,
  children: [
    {
      path: "home",
      icon: "ios-bus",
      meta: {
        title: "主页"
      },
      name: "home_index",
      component: () => import("@/views/home/home.vue")
    },
    {
      path: "message",
      icon: "ios-bus",
      meta: {
        title: "消息"
      },
      name: "message_index",
      component: () => import("@/views/message/message.vue")
    }
  ]
};

// 作为Main组件的子页面展示并且在左侧菜单显示的路由写在appRouter里
export const appRouter = [
  {
    path: "/post",
    icon: "ios-bus",
    name: "component",
    meta: {
      title: "文章管理"
    },
    component: Main,
    children: [
      {
        path: "list",
        icon: "ios-bus",
        name: "post-list",
        meta: {
          title: "文章列表"
        },
        component: () =>
          import("@/views/post/list.vue")
      },
      {
        path: "add",
        icon: "ios-bus",
        name: "post-add",
        meta: {
          title: "添加文章"
        },
        component: () => import("@/views/post/add.vue")
      },
      {
        path: "image-editor",
        icon: "ios-bus",
        name: "image-editor",
        meta: {
          title: "图片编辑"
        },
        component: () =>
          import("@/views/my-components/image-editor/image-editor.vue")
      }
    ]
  },
  {
    path: "/page",
    icon: "ios-bus",
    name: "page",
    meta: {
      title: "页面管理"
    },
    component: Main,
    children: [
      {
        path: "list",
        meta: {
          title: "页面列表"
        },
        name: "page-list",
        icon: "ios-bus",
        component: () =>
          import("@/views/page/list.vue")
      },
      {
        path: "add",
        meta: {
          title: "添加页面"
        },
        name: "page-add",
        icon: "ios-bus",
        component: () => import("@/views/page/add.vue")
      }
    ]
  },
  {
    path: "/cate",
    icon: "ios-bus",
    name: "cate",
    meta: {
      title: "分类管理"
    },
    component: Main,
    children: [
      {
        path: "list",
        meta: {
          title: "分类列表"
        },
        name: "cate-list",
        icon: "ios-bus",
        component: () =>
          import("@/views/form/article-publish/article-publish.vue")
      },
      {
        path: "add",
        meta: {
          title: "添加分类"
        },
        name: "cate-add",
        icon: "ios-bus",
        component: () => import("@/views/form/work-flow/work-flow.vue")
      }
    ]
  },
  {
    path: "/tag",
    icon: "ios-bus",
    name: "tag",
    meta: {
      title: "标签管理"
    },
    component: Main,
    children: [
      {
        path: "list",
        meta: {
          title: "标签列表"
        },
        name: "tag-list",
        icon: "ios-bus",
        component: () => import("@/views/tag/list.vue")
      },
      {
        path: "add",
        meta: {
          title: "添加标签"
        },
        name: "tag-add",
        icon: "ios-bus",
        component: () => import("@/views/tag/add.vue")
      }
    ]
  },
  // {
  //     path: '/charts',
  //     icon: 'ios-bus',
  //     name: 'charts',
  //     title: '图表',
  //     component: Main,
  //     children: [
  //         { path: 'pie', title: '饼状图', name: 'pie', icon: 'ios-pie', component: resolve => { require('@/views/access/access.vue') },
  //         { path: 'histogram', title: '柱状图', name: 'histogram', icon: 'stats-bars', component: resolve => { require('@/views/access/access.vue') }

  //     ]
  // },
  {
    path: "/setting",
    icon: "ios-cog-outline",
    name: "setting",
    meta: {
      title: "系统设置"
    },
    component: Main,
    children: [
      {
        path: "base",
        meta: {
          title: "基本设置"
        },
        name: "setting-base",
        icon: "ios-cog-outline",
        component: () => import("@/views/setting/base.vue")
      },
      {
        path: "comment",
        meta: {
          title: "评论设置"
        },
        name: "setting-comment",
        icon: "ios-text-outline",
        component: () => import("@/views/setting/comment.vue")
      },
      {
        path: "analytic",
        meta: {
          title: "统计设置"
        },
        name: "setting-analytic",
        icon: "ios-pulse",
        component: () => import("@/views/setting/analytic.vue")
      },
      {
        path: "custom",
        meta: {
          title: "自 定 义"
        },
        name: "setting-custom",
        icon: "ios-code-working",
        component: () => import("@/views/setting/custom.vue")
      }
    ]
  },
  {
    path: "/",
    icon: "ios-contact-outline",
    meta: {
      title: "个人中心"
    },
    component: Main,
    children: [
      {
        path: "self",
        meta: {
          title: "个人中心"
        },
        name: "self",
        icon: "ios-contact-outline",
        component: () => import("@/views/user/self.vue")
      }
    ]
  }
];

// 所有上面定义的路由都要写在下面的routers里
export const routers = [
  loginRouter,
  otherRouter,
  preview,
  locking,
  ...appRouter,
  page500,
  page403,
  page404
];
