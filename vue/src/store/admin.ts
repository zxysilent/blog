import { toRaw } from "vue";
import { defineStore } from "pinia";
import appRoutes from "@/router/router";
import { apiAuthGrant } from "@/api";
import { RouteRecordRaw } from "vue-router";
const filter = <T = RouteRecordRaw>(tree: T[], func: (n: T) => boolean): T[] => {
    const children = "children";
    const listFilter = (list: T[]) => {
        return list
            .map((node: any) => ({ ...node }))
            .filter((node) => {
                node[children] = node[children] && listFilter(node[children]);
                return func(node) || (node[children] && node[children].length);
            });
    }

    return listFilter(tree);
}
const filterRoutes = (grants) => {
    let accessedRouters;
    const routeFilter = (route) => {
        if (!route.meta.auth) return true;
        return grants.includes(route.meta.auth);
    };
    try {
        //过滤账户是否拥有某一个权限，并将菜单从加载列表移除
        accessedRouters = filter(appRoutes, routeFilter);
    } catch (error) {
        console.log(error);
        return []
    }
    accessedRouters = accessedRouters.filter(routeFilter);
    return toRaw(accessedRouters);
};
export interface IAdminState {
    id: number;
    num: string;
    name: string;
    phone: string;
    email: string;
    avatar: string;
    role_id: number;
    role: object;
    grants: any[];
    granted: boolean;
    menus: any[];
    token: string;
}

export const useAdminStore = defineStore({
    id: "user",
    state: (): IAdminState => ({
        id: 0,
        num: "",
        name: "",
        phone: "",
        email: "",
        avatar: "",
        role_id: 0,
        role: {},
        grants: [],
        granted: false,
        menus: [],
        token: "",
    }),
    getters: {
        AuthGrant(state) {
            return state.grants;
        },
        AuthGranted(state) {
            return state.granted;
        },
        Authed(state) {
            return (access: string): boolean => {
                const grants = state.grants;
                const accesses = access.split(",")
                return accesses.some(v => grants.includes(v));
            }
        }
    },
    actions: {
        setToken(token: string) {
            this.token = token;
        },
        async fetchGrant() {
            const resp = await apiAuthGrant({});
            if (resp.code === 200) {
                this.grants = resp.data;
                this.granted = true;
                const routes = await filterRoutes(this.grants);
                console.log(routes);
                this.menus = routes;
                return routes;
            } else {
                console.log("fetchGrant=>", resp.msg);
            }
        },
        emitGranted(val) {
            console.log("emitGranted", val);
            this.granted = val;
        },
        // 登出
        async logout() {
            return Promise.resolve("");
        },
    },
});