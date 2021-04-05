import Vue from "vue";
import Utils from "@/utils.js";
import store from "@/store";
import ViewUI from "view-design";
import VueRouter from "vue-router";
import Layout from "@/views/layout.vue";
const _import = require("./_import"); //获取组件的方法-不这样会失败
Vue.use(VueRouter);
export const dynamicRouter = routers => {
	routers.map(item => {
		if (item.comp) {
			item.component = item.comp === "layout" ? Layout : _import(item.comp);
		} else {
			delete item["component"];
		}
		if (item.children) {
			dynamicRouter(item.children);
		}
	});
};
export const initRouter = [
	{
		// 跳转到默认页面
		path: "/",
		name: "index",
		redirect: "/login"
	},
	{
		path: "/login",
		name: "login",
		title: "登录",
		component: () => import("@/views/login.vue")
	},
	{
		path: "/_home",
		name: "_home",
		redirect: "/home",
		component: Layout,
		children: [
			{
				path: "/home",
				title: "管理主页",
				name: "home",
				component: () => import("@/views/home/home.vue")
			}
		]
	}
];
export const errorRouter = [
	{ path: "/jwt", name: "errjwt", title: "jwt-重新登录", component: () => import("@/views/errors/jwt.vue") },
	{ path: "/401", name: "err401", title: "401-没有权限", component: () => import("@/views/errors/401.vue") },
	{ path: "/50x", name: "err50x", title: "50x-服务异常", component: () => import("@/views/errors/50x.vue") },
	{ path: "/*", name: "err404", title: "404-没发现", component: () => import("@/views/errors/404.vue") }
];

// 路由配置
const RouterConfig = {
	mode: "hash", //history
	routes: initRouter
};
let first = true;

const router = new VueRouter(RouterConfig);

router.beforeEach(async (to, from, next) => {
	ViewUI.LoadingBar.start();
	Utils.title(to.meta.title);
	if (first) {
		await store.dispatch("FetchMenu");
		const asyncRouters = store.getters.getRoutes;
		router.addRoutes(asyncRouters);
		first = false;
		next({ ...to, replace: true });
	} else {
		// 已经登陆 去登陆地方
		if (Utils.getToken() && to.name == "login") {
			Utils.title("主页");
			next({ name: "home" });
		} else if (!Utils.getToken() && !Utils.noAuth(to.name)) {
			// //没有登陆 不是去不需要权限的地方
			Utils.title("登陆");
			next({ name: "login" });
		} else {
			next();
		}
	}
});

router.afterEach(to => {
	ViewUI.LoadingBar.finish();
	window.scrollTo(0, 0);
});
export default router;
