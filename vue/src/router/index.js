import Vue from "vue";
import iView from "iview";
import Util from "../libs/util";
import VueRouter from "vue-router";
import { routers, otherRouter, appRouter } from "./router";

Vue.use(VueRouter);

// 路由配置
const RouterConfig = {
  // mode: 'history',
  routes: routers
};

export const router = new VueRouter(RouterConfig);

router.beforeEach((to, from, next) => {
  iView.LoadingBar.start();
  Util.title(to.meta.title);
  // 没有登陆并且不是跳转到登陆页面
  if (!localStorage.getItem("bearer") && to.name !== "login") {
    next({
      name: "login"
    });
  } else if (localStorage.getItem("bearer") && to.name === "login") {
    // 判断是否已经登录且前往的是登录页
    Util.title();
    next({
      name: "home_index"
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
