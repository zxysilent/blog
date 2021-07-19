import Vue from "vue";
import App from "./App.vue";
import router from "./router";
import store from "./store";
import ViewUI from "view-design";
// import style
import "view-design/dist/styles/iview.css";
import { auth } from "@/directive/auth";
// 按钮权限指令
auth(Vue);
Vue.prototype.$auth = function(needs) {
	const need = needs.toString().split(",");
	const access = store.getters && store.getters.AuthGrant;
	return need.some(v => access.includes(v));
};
Vue.use(ViewUI, {
	transfer: true
	// size: "large"
});
//设置为 false 以阻止 vue 在启动时生成生产提示
Vue.config.productionTip = false;
console.log(process.env);

new Vue({
	router,
	store,
	render: h => h(App)
}).$mount("#app");
