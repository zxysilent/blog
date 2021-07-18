import Vue from "vue";
import Utils from "@/utils";
import store from "@/store";
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
			},
			{
				path: "/record",
				name: "record",
				meta: { title: "操作记录" },
				component: () => import("@/views/home/record.vue")
			},
			{
				path: "/global",
				name: "global",
				meta: { title: "配置中心" },
				component: () => import("@/views/home/global.vue")
			},
			{
				path: "/setting",
				name: "setting",
				meta: { title: "信息配置" },
				component: () => import("@/views/home/setting.vue")
			}
		]
	},
	{
		path: "/child",
		name: "child",
		meta: { title: "小孩管理" },
		component: Layout,
		children: [
			{
				path: "list",
				name: "child-list",
				meta: { title: "小孩列表" },
				component: () => import("@/views/child/list.vue")
			},
			{
				path: "add",
				name: "child-add",
				meta: { title: "添加小孩" },
				component: () => import("@/views/child/add.vue")
			},
			{
				path: "edit/:id(\\d+)",
				name: "child-edit",
				meta: { title: "修改小孩" },
				component: () => import("@/views/child/edit.vue")
			},
			{
				path: "view/:id(\\d+)",
				name: "child-view",
				meta: { title: "查看小孩" },
				component: () => import("@/views/child/view.vue")
			}
		]
	},
	{
		path: "/adult",
		name: "adult",
		meta: { title: "大人管理" },
		component: Layout,
		children: [
			{
				path: "list",
				name: "adult-list",
				meta: { title: "大人列表" },
				component: () => import("@/views/adult/list.vue")
			},
			{
				path: "add",
				name: "adult-add",
				meta: { title: "添加大人" },
				component: () => import("@/views/adult/add.vue")
			},
			{
				path: "edit/:id(\\d+)",
				name: "adult-edit",
				meta: { title: "修改大人" },
				component: () => import("@/views/adult/edit.vue")
			}
		]
	},
	{
		path: "/guest",
		name: "guest",
		meta: { title: "微信用户" },
		component: Layout,
		children: [
			{
				path: "list",
				name: "guest-list",
				meta: { title: "微信用户列表" },
				component: () => import("@/views/guest/list.vue")
			},
			{
				path: "edit/:id(\\d+)",
				name: "guest-edit",
				meta: { title: "微信用户" },
				component: () => import("@/views/guest/edit.vue")
			}
		]
	},
	{
		path: "/dict",
		name: "dict",
		meta: { title: "字典管理" },
		component: Layout,
		children: [
			{
				path: "list",
				name: "dict-list",
				meta: { title: "字典列表" },
				component: () => import("@/views/dict/list.vue")
			},
			{
				path: "add",
				name: "dict-add",
				meta: { title: "添加字典" },
				component: () => import("@/views/dict/add.vue")
			},
			{
				path: "edit/:key",
				name: "dict-edit",
				meta: { title: "修改字典" },
				component: () => import("@/views/dict/edit.vue")
			}
		]
	},
	{
		path: "/role",
		name: "role",
		meta: { title: "角色管理" },
		component: Layout,
		children: [
			{
				path: "list",
				name: "role-list",
				meta: { title: "角色列表" },
				component: () => import("@/views/role/list.vue")
			},
			{
				path: "add",
				name: "role-add",
				meta: { title: "添加角色" },
				component: () => import("@/views/role/add.vue")
			},
			{
				path: "edit/:id(\\d+)",
				name: "role-edit",
				meta: { title: "修改角色" },
				component: () => import("@/views/role/edit.vue")
			}
		]
	},
	{
		path: "/user",
		name: "user",
		meta: { title: "用户管理" },
		component: Layout,
		children: [
			{
				path: "list",
				name: "user-list",
				meta: { title: "用户列表" },
				component: () => import("@/views/user/list.vue")
			},
			{
				path: "add",
				name: "user-add",
				meta: { title: "添加用户" },
				component: () => import("@/views/user/add.vue")
			},
			{
				path: "edit/:id(\\d+)",
				name: "user-edit",
				meta: { title: "修改用户" },
				component: () => import("@/views/user/edit.vue")
			}
		]
	},
	{
		path: "/auth",
		meta: { title: "个人中心" },
		component: Layout,
		children: [
			{
				path: "self",
				meta: { title: "个人中心" },
				name: "self",
				component: () => import("@/views/user/self.vue")
			}
		]
	},
	{ path: "/500", name: "500", meta: { title: "500-异常" }, component: () => import("@/components/errors/500.vue") },
	{ path: "/*", name: "404", meta: { title: "404-没发现" }, component: () => import("@/components/errors/404.vue") }
];

// 路由配置
const router = new VueRouter({
	mode: "hash", //hash/history
	routes
});
let init = true;
router.beforeEach(async (to, from, next) => {
	ViewUI.LoadingBar.start();
	Utils.title(to.meta.title);
	if (Utils.getToken() && !Utils.noAuth(to.name)) {
		if (init) {
			await store.dispatch("fetchGrant");
			console.log(store.getters.AuthGrant);
			init = false;
		}
	}
	// 已经登陆 去登陆地方
	if (Utils.getToken() && to.name == "login") {
		Utils.title("主页");
		next({
			name: "home"
		});
	} else if (!Utils.getToken() && !Utils.noAuth(to.name)) {
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
