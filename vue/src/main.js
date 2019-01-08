import Vue from "vue";
import Vuex from "vuex";
import iView from "iview";
import App from "./app.vue";
import Store from "./store.js";
import Util from "./util/util.js";
import Routers from "./router.js";
import VueRouter from "vue-router";
import "iview/dist/styles/iview.css";

Vue.use(Vuex);
//状态管理
const store = new Vuex.Store(Store);
Vue.use(iView, {
  transfer: true
  // size: "large"
});

Vue.use(VueRouter);
// 路由配置
const RouterConfig = {
  mode: "history",
  routes: Routers
};

const router = new VueRouter(RouterConfig);
router.beforeEach((to, from, next) => {
  iView.LoadingBar.start();
  Util.title(to.meta.title);
  // 没有登陆并且不是跳转到登陆页面
  if (!localStorage.getItem("bearer") && to.name !== "login") {
    Util.title("登陆");
    next({
      name: "login"
    });
  } else if (localStorage.getItem("bearer") && to.name === "login") {
    // 判断是否已经登录且前往的是登录页
    next({
      name: "home"
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
  store: store,
  render: h => h(App),
  data: {},
  mounted() {},
  created() {}
});
