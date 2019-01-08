import Main from "@/views/Main.vue";
// 分小模块单独写路由
const initRouter = [
  {
    // 跳转到默认页面
    path: "/",
    name: "index",
    redirect: "/core/login"
  },
  {
    path: "/core/login",
    name: "login",
    meta: {
      title: "登录"
    },
    component: () => import("@/views/login.vue")
  },
  {
    path: "/core/",
    name: "core",
    redirect: "/core/home",
    component: Main,
    children: [
      {
        path: "home",
        meta: {
          title: "管理主页"
        },
        name: "home",
        component: () => import("@/views/home/home.vue")
      },
      {
        path: "message",
        meta: {
          title: "消息"
        },
        name: "message_index",
        component: () => import("@/views/message/message.vue")
      }
    ]
  }
];

const errorRouter = [
  {
    path: "/core/*",
    name: "error-404",
    meta: {
      title: "页面不存在"
    },
    component: () => import("@/views/error/404.vue")
  },
  {
    path: "/core/40x",
    meta: {
      title: "登陆失效"
    },
    name: "error-403",
    component: () => import("@/views/error/40x.vue")
  },
  {
    path: "/core/50x",
    meta: {
      title: "50x-服务端错误"
    },
    name: "error-50x",
    component: () => import("@/views/error/50x.vue")
  }
];
// 作为Main组件的子页面展示并且在左侧菜单显示的路由写在appRouter里
const appRouter = [
  {
    path: "/core/post",
    name: "post",
    meta: {
      title: "文章管理"
    },
    component: Main,
    children: [
      {
        path: "list",
        name: "post-list",
        meta: {
          title: "文章列表"
        },
        component: () => import("@/views/post/list.vue")
      },
      {
        path: "add",
        name: "post-add",
        meta: {
          title: "添加文章"
        },
        component: () => import("@/views/article/article.vue")
      },
      {
        path: "edit/:id(\\d+)",
        name: "post-edit",
        meta: {
          title: "编辑文章"
        },
        component: () => import("@/views/article/article.vue")
      }
    ]
  },
  {
    path: "/core/page",
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
        component: () => import("@/views/page/list.vue")
      },
      {
        path: "add",
        meta: {
          title: "添加页面"
        },
        name: "page-add",
        component: () => import("@/views/article/article.vue")
      },
      {
        path: "edit/:id(\\d+)",
        name: "page-edit",
        meta: {
          title: "编辑页面"
        },
        component: () => import("@/views/article/article.vue")
      }
    ]
  },
  {
    path: "/core/cate",
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
        component: () => import("@/views/cate/list.vue")
      },
      {
        path: "add",
        meta: {
          title: "添加分类"
        },
        name: "cate-add",
        component: () => import("@/views/cate/add.vue")
      }
    ]
  },
  {
    path: "/core/tag",
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
        component: () => import("@/views/tag/list.vue")
      },
      {
        path: "add",
        meta: {
          title: "添加标签"
        },
        name: "tag-add",
        component: () => import("@/views/tag/add.vue")
      }
    ]
  },
  {
    path: "/core/setting",
    name: "setting",
    meta: {
      title: "系统设置",
      icon: "ios-cog-outline"
    },
    component: Main,
    children: [
      {
        path: "base",
        meta: {
          title: "基本设置"
        },
        name: "setting-base",
        component: () => import("@/views/setting/base.vue")
      },
      {
        path: "/core/test",
        meta: {
          title: "基本设置"
        },
        name: "test",
        component: () => import("@/views/my-components/image-editor/image-editor.vue")
      },
      {
        path: "comment",
        meta: {
          title: "评论设置",
          icon: "ios-text-outline"
        },
        name: "setting-comment",
        component: () => import("@/views/setting/comment.vue")
      },
      {
        path: "analytic",
        meta: {
          icon: "ios-pulse",
          title: "统计设置"
        },
        name: "setting-analytic",
        component: () => import("@/views/setting/analytic.vue")
      },
      {
        path: "custom",
        meta: {
          icon: "ios-code-working",
          title: "自 定 义"
        },
        name: "setting-custom",
        component: () => import("@/views/setting/custom.vue")
      }
    ]
  },
  {
    path: "/core/",
    icon: "ios-contact-outline",
    meta: {
      title: "个人中心"
    },
    component: Main,
    children: [
      {
        path: "self",
        meta: {
          title: "个人中心",
          icon: "ios-contact-outline"
        },
        name: "self",
        component: () => import("@/views/user/self.vue")
      }
    ]
  }
];
// 所有定义的路由都要写在下面的routers里
const routers = [...initRouter, ...appRouter, ...errorRouter];
export default routers;
