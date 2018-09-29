import Vue from "vue";
import iView from "iview";
import { router } from "./router/index";
import { appRouter } from "./router/router";
import store from "./store";
import App from "./app.vue";
import "iview/dist/styles/iview.css";
import util from "./libs/util";

Vue.use(iView, {
  transfer: true,
//   size: "large"
});

new Vue({
  el: "#app",
  router: router,
  store: store,
  render: h => h(App),
  data: {
    currentPageName: ""
  },
  mounted() {
    this.currentPageName = this.$route.name;
    // 显示打开的页面的列表
    this.$store.commit("setOpenedList");
    this.$store.commit("initCachepage");
    // 权限菜单过滤相关
    this.$store.commit("updateMenulist");
    // iview-admin检查更新
    util.checkUpdate(this);
  },
  created() {
  }
});
