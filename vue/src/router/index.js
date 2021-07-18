import Vue from "vue";
import { getToken } from "@/utils/token";
import store from "@/store";
import ViewUI from "view-design";
import VueRouter from "vue-router";
import Layout from "@/views/Layout.vue";
const _import = require("./_import"); //获取组件的方法-不这样会失败
Vue.use(VueRouter);
export const dynamicRouter = routers => {
	routers.map(item => {
		item.meta = {
			title: item.title //标题显示
		};
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
		meta: {
			title: "登录"
		},
		component: () => import("@/views/Login.vue")
	},
	{
		path: "/_home",
		name: "_home",
		redirect: "/home",
		component: Layout,
		children: [
			{
				path: "/home",
				meta: { title: "管理主页" },
				name: "home",
				component: () => import("@/views/home/index.vue")
			},
			{
				path: "/auth/self",
				name: "auth-self",
				meta: { title: "个人信息" },
				component: () => import("@/views/user/self.vue")
			}
		]
	},
	{ path: "/500", name: "500", meta: { title: "500-异常" }, component: () => import("@/components/errors/500.vue") },
	{ path: "/*", name: "404", meta: { title: "404-没发现" }, component: () => import("@/components/errors/404.vue") }
];

const toTitle = title => {
	title = title || "zadmin";
	window.document.title = title;
};
// 路由配置
const RouterConfig = {
	mode: "hash", //history
	routes: initRouter
};
let first = true;

const router = new VueRouter(RouterConfig);
const noAuth = ["login", "errjwt", "err50x", "err401"];
router.beforeEach(async (to, from, next) => {
	ViewUI.LoadingBar.start();
	toTitle(to.meta.title);
	const token = getToken();
	console.log("token", token);
	console.log("to.name", to.name);
	// 去不需要登录的地方
	if (noAuth.indexOf(to.name) > -1) {
		if (token) {
			toTitle("主页");
			next({ name: "home" });
		} else {
			next();
		}
	} else {
		if (token) {
			if (first) {
				await store.dispatch("authMenu");
				const asyncRouters = store.getters.getRoutes;
				router.addRoutes(asyncRouters);
				first = false;
				next({ ...to, replace: true });
			} else {
				next();
			}
		} else {
			//没有登陆 不是去不需要权限的地方
			toTitle("登陆");
			next({ name: "login" });
		}
	}
	// next()
});

router.afterEach(to => {
	ViewUI.LoadingBar.finish();
	window.scrollTo(0, 0);
});
export default router;
