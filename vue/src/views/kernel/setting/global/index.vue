<template>
    <n-card :bordered="false" class="proCard">
        <!-- <n-tabs type="line" v-model:value="def_val" size="large"> -->
        <n-tabs type="line" size="large">
            <n-tab-pane name="basic" tab="基础配置"><Basic /> </n-tab-pane>
            <n-tab-pane name="login" tab="登录页面"><Login /> </n-tab-pane>
            <n-tab-pane name="blog" tab="博文配置"><Blog /> </n-tab-pane>
            <!-- <n-tab-pane v-for="item in tabAll" :key="item.name" :name="item.name" :tab="item.tab">
                <component :is="item.cname"></component>
            </n-tab-pane> -->
        </n-tabs>
    </n-card>
</template>
<script lang="ts" setup>
import Basic from "./basic.vue";
import Login from "./login.vue";
import Blog from "./blog.vue";

// import Theme from "./theme.vue";
// import { ref, onMounted } from "vue";
// import { useAdminStore } from "@/store/admin";
// const adminStore = useAdminStore();
// const tabAll = Array<any>();
// const def_val = ref<string>("");
// const init = () => {
//     if (adminStore.Authed("home_mgmt")) {
//         tabAll.push({ name: "basic", tab: "基础配置", cname: Basic });
//         tabAll.push({ name: "login", tab: "登录页面", cname: Login });
//     }
//     if (tabAll.length >= 1) {
//         def_val.value = tabAll[0].name;
//     }
// };
// onMounted(() => {
//     init();
// });
</script>
<style lang="less" scoped></style>
