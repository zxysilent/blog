import Vue from "vue";
import Vuex from "vuex";
import { admAuthGrant } from "@/api/auth";
Vue.use(Vuex);

const Store = {
	state: {
		grants: [] //授权
	},
	mutations: {
		setGrants(state, grants) {
			state.grants = grants;
		}
	},
	actions: {
		// 从后台获取授权
		async fetchGrant({ commit }) {
			const resp = await admAuthGrant();
			console.log("admAuthGrant");
			console.log(resp);
			if (resp.code == 200) {
				commit("setGrants", resp.data);
			} else {
				commit("setGrants", []);
			}
			return true;
		}
	},
	getters: {
		AuthGrant(state) {
			return state.grants;
		}
	}
};

const store = new Vuex.Store(Store);
export default store;
