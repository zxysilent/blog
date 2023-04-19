<template>
    <n-menu
        :options="menus"
        :inverted="inverted"
        :mode="mode"
        :collapsed="collapsed"
        :collapsed-width="64"
        :collapsed-icon-size="20"
        :indent="24"
        :expanded-keys="state.openKeys"
        :value="getSelectedKeys"
        @update:value="clickMenuItem"
        @update:expanded-keys="menuExpanded" />
</template>

<script lang="ts" setup>
import { ref, onMounted, reactive, computed, watch, unref } from "vue";
import { useRoute, useRouter } from "vue-router";
import { generatorMenu } from "@/utils";
import { useAdminStore } from "@/store/admin";
import setting from "@/utils/global";
const adminStore = useAdminStore();
const props = defineProps({
    // 菜单模式
    mode: { type: String, default: "vertical" },
    // 侧边栏菜单是否收起
    collapsed: { type: Boolean },
});
console.log(props.mode);
// 当前路由
const currentRoute = useRoute();
const router = useRouter();
const menus = ref<any[]>([]);
const selectedKeys = ref<string>((currentRoute.meta?.active as string) ? (currentRoute.meta?.active as string) : (currentRoute.name as string));
// 获取当前打开的子菜单
const matched = currentRoute.matched;

const getOpenKeys = matched && matched.length ? matched.map((item) => item.name) : [];

const state = reactive({
    openKeys: getOpenKeys,
});

const inverted = computed(() => {
    return ["dark", "header-dark"].includes(setting.navTheme);
});

const getSelectedKeys = computed(() => {
    return selectedKeys.value;
});

// 监听分割菜单
// watch(
//     () => setting.menuSetting.mixMenu,
//     () => {
//         updateMenu();
//         if (props.collapsed) {
//             emit("update:collapsed", !props.collapsed);
//         }
//     }
// );

// 跟随页面路由变化，切换菜单选中状态
watch(
    () => currentRoute.fullPath,
    () => {
        updateMenu();
        const matched = currentRoute.matched;
        state.openKeys = matched.map((item) => item.name);
        const activeMenu: string = (currentRoute.meta?.active as string) || "";
        selectedKeys.value = activeMenu ? (activeMenu as string) : (currentRoute.name as string);
    }
);

function updateMenu() {
    menus.value = generatorMenu(adminStore.menus);
}

// 点击菜单
function clickMenuItem(key: string) {
    if (/http(s)?:/.test(key)) {
        window.open(key);
    } else {
        router.push({ name: key });
    }
}

//展开菜单
function menuExpanded(openKeys: string[]) {
    if (!openKeys) return;
    const latestOpenKey = openKeys.find((key) => state.openKeys.indexOf(key) === -1);
    const isExistChildren = findChildrenLen(latestOpenKey as string);
    state.openKeys = isExistChildren ? (latestOpenKey ? [latestOpenKey] : []) : openKeys;
}

//查找是否存在子路由
function findChildrenLen(key: string) {
    if (!key) return false;
    const subRouteChildren: string[] = [];
    for (const { children, key } of unref(menus)) {
        if (children && children.length) {
            subRouteChildren.push(key as string);
        }
    }
    return subRouteChildren.includes(key);
}

onMounted(() => {
    updateMenu();
});
</script>
