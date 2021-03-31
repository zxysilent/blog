import Vue from "vue";
import App from "./App.vue";

import router from "./router";
import store from "./store";
import ViewUI from "view-design";

// import style
import "view-design/dist/styles/iview.css";

Vue.use(ViewUI, { transfer: true });

//设置为 false 以阻止 vue 在启动时生成生产提示
Vue.config.productionTip = false;
console.log(process.env);

new Vue({
	router,
	store,
	render: h => h(App)
}).$mount("#app");
