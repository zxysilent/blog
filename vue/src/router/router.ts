import { RouteRecordRaw } from "vue-router";
import { SettingsOutline, HomeOutline, PeopleOutline, PaperPlaneOutline, PencilOutline } from "@vicons/ionicons5";
import { renderIcon, } from "@/utils/index";//renderNew
const Layout = () => import("@/components/Layout/index.vue");
const Empty = () => import("@/components/Layout/empty.vue");
/**
 * @param name 路由名称, 必须设置,且不能重名
 * @param meta 路由元信息（路由附带扩展信息）
 * @param redirect 重定向地址, 访问这个路由时,自定进行重定向
 * @param meta.title 菜单名称
 * @param meta.icon 菜单图标
 * */
const rootRoutes: RouteRecordRaw[] = [
    { path: "/", name: "root", redirect: "/home" },
    { path: "/login", name: "login", meta: { title: "登录" }, component: () => import("@/views/kernel/login/index.vue") },
    // { path: "/admin", name: "admin", meta: { title: "登录" }, component: () => import("@/views/login/index.vue") },
    { path: "/403", name: "403", meta: { title: "403" }, component: () => import("@/components/Exception/403.vue") },
    { path: "/500", name: "500", meta: { title: "500" }, component: () => import("@/components/Exception/500.vue") },
    // { path: "/succ", name: "result-succ", meta: { title: "成功页", }, component: () => import("@/views/result/succ.vue") },
    // { path: "/fail", name: "result-fail", meta: { title: "失败页", }, component: () => import("@/views/result/fail.vue") },
    // { path: "/info", name: "result-info", meta: { title: "信息页", }, component: () => import("@/views/result/info.vue") },
];
const appRoutes: Array<RouteRecordRaw> = [
    {
        path: "/layout",
        name: "layout",
        redirect: "/home",
        component: Layout,
        meta: { title: "主控界面", icon: renderIcon(HomeOutline) },
        children: [
            { path: "/home", name: "home", meta: { title: "系统主页" }, component: () => import("@/views/home/index.vue") },
            { path: "/log", name: "log", meta: { title: "系统日志", auth: "home_log" }, component: () => import("@/views/home/log.vue") },
            { path: "/about", name: "about", meta: { title: "系统关于", }, component: () => import("@/views/home/about.vue") },
        ],
    },
    {
        path: "/note",
        name: "note",
        component: () => import("@/views/applet/note/index.vue"),
        meta: { title: "笔记模式", hidden: true, icon: renderIcon(PencilOutline) },
    },
    /*
     * applet
     */
    {
        path: "/applet",
        component: Layout,
        meta: { title: "博文模块", icon: renderIcon(PaperPlaneOutline), auth: "home_applet" },
        name: "applet",
        children: [
            {
                path: "/post",
                component: Empty,
                meta: { title: "博文管理", auth: "home_applet" },
                name: "post",
                children: [
                    { path: "/post/index", name: "post-index", meta: { title: "博文列表", auth: "home_applet" }, component: () => import("@/views/applet/post/index.vue") },
                    { path: "/post/add", name: "post-add", meta: { title: "添加博文", auth: "home_mgmt", hidden: true, active: "post-index" }, component: () => import("@/views/applet/post/add.vue") },
                    {
                        path: "/post/edit/:id(\\d+)",
                        name: "post-edit", meta: { title: "修改博文", auth: "home_mgmt", hidden: true, active: "post-index" }, component: () => import("@/views/applet/post/edit.vue")
                    },
                    {
                        path: "/post/detail/:id(\\d+)",
                        name: "post-detail", meta: { title: "查看博文", auth: "home_applet", hidden: true, active: "post-index" }, component: () => import("@/views/applet/post/detail.vue")
                    }
                ]
            },
            {
                path: "/cate",
                component: Empty,
                meta: { title: "分类管理", auth: "home_applet" },
                name: "cate",
                children: [
                    { path: "/cate/index", name: "cate-index", meta: { title: "分类列表", auth: "home_applet" }, component: () => import("@/views/applet/cate/index.vue"), },
                    { path: "/cate/add", name: "cate-add", meta: { title: "添加分类", auth: "home_mgmt", hidden: true, active: "cate-index" }, component: () => import("@/views/applet/cate/add.vue"), },
                    {
                        path: "/cate/edit/:id(\\d+)",
                        name: "cate-edit", meta: { title: "修改分类", auth: "home_mgmt", hidden: true, active: "cate-index" }, component: () => import("@/views/applet/cate/edit.vue"),
                    },
                    {
                        path: "/cate/detail/:id(\\d+)",
                        name: "cate-detail", meta: { title: "查看分类", auth: "home_applet", hidden: true, active: "cate-index" }, component: () => import("@/views/applet/cate/detail.vue"),
                    }
                ]
            },
            {
                path: "/tag",
                component: Empty,
                meta: { title: "标签管理", auth: "home_applet" },
                name: "tag",
                children: [
                    { path: "/tag/index", name: "tag-index", meta: { title: "标签列表", auth: "home_applet" }, component: () => import("@/views/applet/tag/index.vue"), },
                    { path: "/tag/add", name: "tag-add", meta: { title: "添加标签", auth: "home_mgmt", hidden: true, active: "tag-index" }, component: () => import("@/views/applet/tag/add.vue"), },
                    {
                        path: "/tag/edit/:id(\\d+)",
                        name: "tag-edit", meta: { title: "修改标签", auth: "home_mgmt", hidden: true, active: "tag-index" }, component: () => import("@/views/applet/tag/edit.vue"),
                    },
                    {
                        path: "/tag/detail/:id(\\d+)",
                        name: "tag-detail", meta: { title: "查看标签", auth: "home_applet", hidden: true, active: "tag-index" }, component: () => import("@/views/applet/tag/detail.vue"),
                    }
                ]
            },
        ]
    },
    /*
     * applet
     */
    {
        path: "/admin",
        name: "admin",
        component: Layout,
        meta: { title: "角色权限", icon: renderIcon(PeopleOutline), auth: "admin_view" },
        children: [
            { path: "/admin/index", name: "admin-index", meta: { title: "管理用户", auth: "admin_view" }, component: () => import("@/views/kernel/admin/index.vue") },
            { path: "/admin/add", name: "admin-add", meta: { title: "添加管理", auth: "admin_add", hidden: true, active: "admin-index" }, component: () => import("@/views/kernel/admin/add.vue") },
            { path: "/admin/edit/:id(\\d+)", name: "admin-edit", meta: { title: "修改管理", auth: "admin_edit", hidden: true, active: "admin-index" }, component: () => import("@/views/kernel/admin/edit.vue") },
            { path: "/admin/detail/:id(\\d+)", name: "admin-detail", meta: { title: "查看管理", auth: "admin_view", hidden: true, active: "admin-index" }, component: () => import("@/views/kernel/admin/detail.vue") },
            { path: "/grant/index", name: "grant-index", meta: { title: "权限管理", auth: "grant_view" }, component: () => import("@/views/kernel/admin/grant/index.vue") },
            { path: "/grant/add", name: "grant-add", meta: { title: "添加权限", auth: "grant_add", hidden: true, active: "grant-index" }, component: () => import("@/views/kernel/admin/grant/add.vue") },
            { path: "/grant/edit/:id(\\d+)", name: "grant-edit", meta: { title: "修改权限", auth: "grant_edit", hidden: true, active: "grant-index" }, component: () => import("@/views/kernel/admin/grant/edit.vue") },
            {
                path: "/grant/detail/:id(\\d+)",
                name: "grant-detail", meta: { title: "查看权限", auth: "grant_view", hidden: true, active: "grant-index" }, component: () => import("@/views/kernel/admin/grant/detail.vue")
            },
            { path: "/role/index", name: "role-index", meta: { title: "角色管理", auth: "grant_view" }, component: () => import("@/views/kernel/admin/role/index.vue") },
            { path: "/role/add", name: "role-add", meta: { title: "添加角色", auth: "role_add", hidden: true, active: "role-index" }, component: () => import("@/views/kernel/admin/role/add.vue") },
            { path: "/role/edit/:id(\\d+)", name: "role-edit", meta: { title: "修改角色", auth: "role_edit", hidden: true, active: "role-index" }, component: () => import("@/views/kernel/admin/role/edit.vue") },
            { path: "/role/detail/:id(\\d+)", name: "role-detail", meta: { title: "查看角色", auth: "role_view", hidden: true, active: "role-index" }, component: () => import("@/views/kernel/admin/role/detail.vue") }
        ],
    },
    {
        path: "/setting",
        name: "setting",
        redirect: "/setting/account",
        component: Layout,
        meta: { title: "系统设置", icon: renderIcon(SettingsOutline), },
        children: [
            { path: "/account", name: "account", meta: { title: "个人设置", }, component: () => import("@/views/kernel/setting/account/account.vue") },
            { path: "/global", name: "global", meta: { title: "全局配置", }, component: () => import("@/views/kernel/setting/global/index.vue") },
            { path: "/dict/index", name: "dict-index", meta: { title: "配置字典", active: "dict-index", auth: "dict_view" }, component: () => import("@/views/kernel/setting/dict/index.vue") },
            { path: "/dict/add", name: "dict-add", meta: { title: "添加字典", hidden: true, active: "dict-index", auth: "dict_add" }, component: () => import("@/views/kernel/setting/dict/add.vue") },
            { path: "/edit/:key", name: "dict-edit", meta: { title: "修改字典", hidden: true, active: "dict-index", auth: "dict_edit" }, component: () => import("@/views/kernel/setting/dict/edit.vue") },
            { path: "/detail/:key", name: "dict-detail", meta: { title: "查看字典", hidden: true, active: "dict-index", auth: "dict_view" }, component: () => import("@/views/kernel/setting/dict/detail.vue") },
        ],
    },
    { path: "/:path(.*)*", name: "404", meta: { title: "404 not found", hidden: true }, component: () => import("@/components/Exception/404.vue") },
    // {
    //     path: "/external",
    //     name: "https://naive-ui-admin-docs.vercel.app",
    //     component: Layout,
    //     meta: { title: "项目文档", icon: renderIcon(ExclamationCircleOutlined), },
    // },
];
export { rootRoutes };
export default appRoutes;
