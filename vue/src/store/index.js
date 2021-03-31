import Vue from "vue";
import Vuex from "vuex";
import { apiMenuTree } from "@/api/menu";
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
		async FetchMenu({ commit }) {
			const resp = await apiMenuTree();
			if (resp.code == 200) {
				const cp = resp.data.slice();
				commit("setMenus", resp.data);
				dynamicRouter(cp);
				commit("setRoutes", cp);
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
