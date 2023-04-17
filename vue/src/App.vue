<template>
    <NConfigProvider :locale="zhCN" :theme="getDarkTheme" :theme-overrides="getThemeOverrides" :date-locale="dateZhCN">
        <AppProvider>
            <RouterView />
        </AppProvider>
    </NConfigProvider>
    <transition v-if="$route.name !== 'login'" name="slide-up"> </transition>
</template>

<script lang="ts" setup>
import { computed } from "vue";
import { zhCN, dateZhCN, darkTheme } from "naive-ui";
import { AppProvider } from "@/components/Application";
import setting from "@/utils/global";
/**
 * @type import('naive-ui').GlobalThemeOverrides
 */
const getThemeOverrides = computed(() => {
    return {
        common: {
            primaryColor: "#2d8cf0ff",
            primaryColorHover: "#57a3f3ff",
            primaryColorPressed: "#2b85e4ff",
            primaryColorSuppl: "#2D8CF0ff",
            infoColorPressed: "#2baee9ff",
            infoColorHover: "#57c5f7ff",
            infoColor: "#2db7f5ff",
            infoColorSuppl: "#2DB7F5ff",
            successColor: "#19be6bff",
            successColorPressed: "#18b566ff",
            successColorHover: "#47cb89ff",
            successColorSuppl: "#19BE6Bff",
            // warningColorPressed: "#f29100ff",
            // warningColorSuppl: "#fcb040ff",
            // warningColorHover: "#ffad33ff",
            // warningColor: "#ff9900ff",
            // errorColorPressed: "#e13d13ff",
            // errorColorHover: "#f16643ff",
            // errorColorSuppl: "#ed4014ff",
            // errorColor: "#ed4014ff",
        },
        LoadingBar: {
            colorLoading: "#2d8cf0ff",
        },
    };
});
// { xs: 0, s: 640, m: 1024, l: 1280, xl: 1536, xxl: 1920 }
const getDarkTheme = computed(() => (setting.darkTheme ? darkTheme : undefined));
</script>

<style lang="less">
@import "styles/index.less";
</style>
