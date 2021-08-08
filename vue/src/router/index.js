import Vue from "vue";
import Utils from "@/utils";
import Storage from "@/utils/storage";
import ViewUI from "view-design";
import VueRouter from "vue-router";
import Layout from "@/views/Layout.vue";

Vue.use(VueRouter);

const routes = [
	{ path: "/", name: "init", redirect: "/login" },
	{ path: "/login", name: "login", meta: { title: "登录" }, component: () => import("@/views/Login.vue") },
	{
		path: "/layout",
		meta: { title: "管理主页" },
		name: "layout",
		component: Layout,
		children: [
			{
				path: "/home",
				name: "home",
				meta: { title: "管理主页" },
				component: () => import("@/views/home/index.vue")
			}
		]
	},
	{
		path: "/sys",
		name: "sys",
		meta: { module: "sys", title: "系统模块" },
		component: Layout,
		children: [
			{
				path: "/global",
				name: "global",
				meta: { module: "sys", title: "配置中心" },
				component: () => import("@/views/home/global.vue")
			},
			{
				path: "/profile",
				meta: { module: "sys", title: "个人中心" },
				name: "profile",
				component: () => import("@/views/user/self.vue")
			},

			{
				path: "/dict/list",
				name: "dict-list",
				meta: { module: "sys", title: "字典管理" },
				component: () => import("@/views/dict/list.vue")
			},
			{
				path: "/dict/add",
				name: "dict-add",
				meta: { module: "sys", title: "添加字典" },
				component: () => import("@/views/dict/add.vue")
			},
			{
				path: "/dict/edit/:key",
				name: "dict-edit",
				meta: { module: "sys", title: "编辑字典" },
				component: () => import("@/views/dict/edit.vue")
			}
		]
	},
	{
		path: "/app",
		name: "app",
		meta: { module: "app", title: "博客模块" },
		component: Layout,
		children: [
			{
				path: "/post/list",
				name: "post-list",
				meta: { module: "app", title: "文章列表" },
				component: () => import("@/views/post/list.vue")
			},
			{
				path: "/post/add",
				name: "post-add",
				meta: { module: "app", title: "添加文章" },
				component: () => import("@/views/post/add.vue")
			},
			{
				path: "/post/edit/:id(\\d+)",
				name: "post-edit",
				meta: { module: "app", title: "编辑文章" },
				component: () => import("@/views/post/edit.vue")
			},
			{
				path: "/page/list",
				meta: { module: "app", title: "页面列表" },
				name: "page-list",
				component: () => import("@/views/page/list.vue")
			},
			{
				path: "/page/add",
				meta: { module: "app", title: "添加页面" },
				name: "page-add",
				component: () => import("@/views/page/add.vue")
			},
			{
				path: "/page/edit/:id(\\d+)",
				name: "page-edit",
				meta: { module: "app", title: "编辑页面" },
				component: () => import("@/views/page/edit.vue")
			},
			{
				path: "/cate/list",
				meta: { module: "app", title: "分类列表" },
				name: "cate-list",
				component: () => import("@/views/cate/list.vue")
			},
			{
				path: "/cate/add",
				meta: { module: "app", title: "添加分类" },
				name: "cate-add",
				component: () => import("@/views/cate/add.vue")
			},
			{
				path: "/cate/edit/:id(\\d+)",
				name: "cate-edit",
				meta: { module: "app", title: "修改分类" },
				component: () => import("@/views/cate/edit.vue")
			},
			{
				path: "/tag/list",
				meta: { module: "app", title: "标签列表" },
				name: "tag-list",
				component: () => import("@/views/tag/list.vue")
			},
			{
				path: "/tag/add",
				meta: { module: "app", title: "添加标签" },
				name: "tag-add",
				component: () => import("@/views/tag/add.vue")
			}
		]
	},
	{ path: "/500", name: "500", meta: { title: "500-异常" }, component: () => import("@/components/errors/500.vue") },
	{ path: "/*", name: "404", meta: { title: "404-没发现" }, component: () => import("@/components/errors/404.vue") }
];

// 路由配置
const router = new VueRouter({
	mode: "hash", //hash,history
	routes
});

router.beforeEach(async (to, from, next) => {
	ViewUI.LoadingBar.start();
	Utils.title(to.meta.title);
	// 已经登陆 去登陆地方
	if (Storage.getToken() && to.name == "login") {
		Utils.title("主页");
		next({
			name: "home"
		});
	} else if (!Storage.getToken() && !Utils.noAuth(to.name)) {
		// //没有登陆 不是去不需要权限的地方
		Utils.title("登陆");
		next({
			name: "login"
		});
	} else {
		next();
	}
	ViewUI.LoadingBar.finish();
});

router.afterEach(to => {
	ViewUI.LoadingBar.finish();
	window.scrollTo(0, 0);
});

export default router;
