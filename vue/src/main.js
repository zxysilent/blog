import Vue from "vue";
// 移除 vuex
// import Vuex from "vuex";
import iView from "iview";
import App from "./app.vue";
// 移除 vuex
// import Store from "./store.js";
import Util from "./init/util.js";
import Routers from "./router.js";
import VueRouter from "vue-router";
import "iview/dist/styles/iview.css";

//设置为 false 以阻止 vue 在启动时生成生产提示
Vue.config.productionTip = false;
console.log(process.env)
// 移除 vuex
// 状态管理
// Vue.use(Vuex);
// const store = new Vuex.Store(Store);
Vue.use(iView, {
	transfer: true
	// size: "large"
});

Vue.use(VueRouter);
// 路由配置
const RouterConfig = {
	mode: "hash",//history
	routes: Routers
};

const router = new VueRouter(RouterConfig);
router.beforeEach((to, from, next) => {
	iView.LoadingBar.start();
	Util.title(to.meta.title);
	// 已经登陆 去登陆地方
	if (Util.getToken() && to.name == "login") {
		Util.title("主页");
			next({
			name: "home"
		});
	} else if (!Util.getToken() && !Util.noAuth(to.name)) {
		// //没有登陆 不是去不需要权限的地方
		Util.title("登陆");
		next({
			name: "login"
		});
	} else {
		next();
	}
	iView.LoadingBar.finish();
});

router.afterEach(to => {
	iView.LoadingBar.finish();
	window.scrollTo(0, 0);
});

new Vue({
	el: "#app",
	router: router,
	// 移除 vuex
	//store: store,
	render: h => h(App)
});
