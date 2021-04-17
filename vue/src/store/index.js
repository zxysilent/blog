import Vue from "vue";
import Vuex from "vuex";
import { admAuthMenu } from "@/api/auth";
Vue.use(Vuex);
import { dynamicRouter } from "@/router";

const Store = {
	state: {
		menus: [], //左侧导航
		routes: [] //路由集合
	},
	mutations: {
		setMenus(state, menus) {
			//设置导航
			state.menus = menus;
		},
		setRoutes(state, routes) {
			//设置路由
			state.routes = routes;
		}
	},
	actions: {
		// 从后台获取菜单
		async authMenu({ commit }) {
			const resp = await admAuthMenu();
			if (resp.code == 200) {
				// deep clone
				// arr.slice 仅复制一层,不能处理数组对象
				// const routes = resp.data.slice();
				const routes = JSON.parse(JSON.stringify(resp.data));
				commit("setMenus", resp.data);
				console.log(resp.data);
				// 仅放入路由,菜单没有
				routes.push({
					path: "/*",
					name: "err404",
					title: "404-没发现",
					comp: "views/errors/404.vue"
				});
				dynamicRouter(routes);
				// routes.push(errorRouter);
				// console.log(errorRouter);
				console.log(routes);
				commit("setRoutes", routes);
			}
			return true;
		}
	},
	getters: {
		getMenus(state) {
			return state.menus;
		},
		getRoutes(state) {
			return state.routes;
		}
	}
};
const store = new Vuex.Store(Store);
export default store;
