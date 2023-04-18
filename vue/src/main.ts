import App from "./App.vue";
import router from "./router";
import "./styles/tailwind.css";
import { createApp } from "vue";
import naive from "@/utils/naive";
import { createPinia } from 'pinia';
import { auth } from "@/directives/auth";
import { discrete } from "./utils/global";
const app = createApp(App);
const store = createPinia();
// Vue.prototype.$auth = function(needs) {
// 	const need = needs.toString().split(",");
// 	const access = store.getters && store.getters.AuthGrant;
// 	return need.some(v => access.includes(v));
// };
const bootstrap = async () => {
    // 注册全局自定义指令
    app.directive("auth", auth);
    app.use(naive);
    app.use(store)
    discrete()
    // 挂载路由
    await app.use(router);
    // 路由准备就绪后挂载APP实例
    await router.isReady();
    const meta = document.createElement('meta')
    meta.name = 'naive-ui-style'
    document.head.appendChild(meta)
    app.mount("#app", true);
}
bootstrap();
