import { rootRoutes } from "./router";
import { createRouter, createWebHashHistory, } from "vue-router";
import storage from "@/utils/storage";
import { useAdminStore } from "@/store/admin";
import { RouteRecordRaw } from "vue-router";
const whiteRouter = ["500", "404", "403", "login", "admin"];
const router = createRouter({
    history: createWebHashHistory(),
    routes: rootRoutes,
    strict: true,
    scrollBehavior: () => ({ left: 0, top: 0 }),
});

router.beforeEach(async (to, _, next) => {
    console.log("to ", to)
    window.$loadingBar && window.$loadingBar.start();
    // 已经登陆 去登陆地方
    if (storage.getToken() && to.name == "login") {
        // next({ name: "home" })
        // return
        to.name = "home"
    }
    if (storage.getToken() && !whiteRouter.includes(to.name as string)) {
        const adminStore = useAdminStore();
        if (!adminStore.granted) {
            const routes = await adminStore.fetchGrant();
            // 动态添加可访问路由表
            routes.forEach((item) => {
                router.addRoute(item as unknown as RouteRecordRaw);
            });
            console.log(router.getRoutes())
            console.log(adminStore.grants);
            adminStore.emitGranted(true)
            next({ ...to, replace: true })
            return
        }
    }
    if (!storage.getToken() && !whiteRouter.includes(to.name as string)) {
        // 没有登陆 不是去不需要权限的地方
        next({ name: "login" })
        return
    }
    next();
});

router.afterEach((to, _) => {
    document.title = (to?.meta?.title as string) || document.title;
    window.scrollTo(0, 0);
    window.$loadingBar && window.$loadingBar.finish();
});
router.onError((error) => {
    console.log("路由错误=>", error);
});

export default router;
