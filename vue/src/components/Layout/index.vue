<template>
    <n-layout class="layout" :position="fixedMenu" has-sider>
        <n-layout-sider
            v-if="navMode === 'vertical'"
            show-trigger="bar"
            @collapse="collapsed = true"
            :position="fixedMenu"
            @expand="collapsed = false"
            :collapsed="collapsed"
            collapse-mode="width"
            :collapsed-width="64"
            :width="leftMenuWidth"
            :native-scrollbar="false"
            :inverted="inverted"
            class="layout-sider">
            <Logo :collapsed="collapsed" />
            <AsideMenu v-model:collapsed="collapsed" />
        </n-layout-sider>
        <n-layout :inverted="inverted">
            <n-layout-header :inverted="getHeaderInverted" :position="fixedHeader">
                <PageHeader v-model:collapsed="collapsed" :inverted="inverted" />
            </n-layout-header>
            <n-layout-content class="layout-content" :class="{ 'layout-default-background': setting.darkTheme === false }">
                <div class="layout-content-main" :class="{ 'fluid-header': fixedHeader === 'static' }">
                    <div class="main-view mt-3">
                        <MainView />
                    </div>
                </div>
                <!--1.15废弃，没啥用，占用操作空间-->
                <n-layout-footer v-if="setting.showFooter" position="static">
                    <PageFooter />
                </n-layout-footer>
            </n-layout-content>
            <n-back-top :bottom="60" />
        </n-layout>
    </n-layout>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted } from "vue";
import { Logo } from "./components/Logo";
import { MainView } from "./components/Main";
import { AsideMenu } from "./components/Menu";
import { PageHeader } from "./components/Header";
import { PageFooter } from "./components/Footer";
import setting from "@/utils/global";

const navMode = setting.navMode;

const collapsed = ref<boolean>(false);

const fixedHeader = computed(() => {
    return setting.headerSetting.fixed ? "absolute" : "static";
});

const fixedMenu = computed(() => {
    return setting.menuSetting.fixed ? "absolute" : "static";
});

const inverted = computed(() => {
    return ["dark", "header-dark"].includes(setting.navTheme);
});

const getHeaderInverted = computed(() => {
    return ["light", "header-dark"].includes(setting.navTheme) ? inverted : !inverted;
});

const leftMenuWidth = computed(() => {
    const { minMenuWidth, menuWidth } = setting.menuSetting;
    return collapsed.value ? minMenuWidth : menuWidth;
});

const watchWidth = () => {
    const Width = document.body.clientWidth;
    if (Width <= 950) {
        collapsed.value = true;
    } else collapsed.value = false;
};

onMounted(() => {
    window.addEventListener("resize", watchWidth);
    //挂载在 window 方便与在js中使用
});
</script>

<style lang="less" scoped>
.layout {
    display: flex;
    flex-direction: row;
    flex: auto;
    &-default-background {
        background: #f5f7f9;
    }

    .layout-sider {
        min-height: 100vh;
        box-shadow: 2px 0 8px 0 rgb(29 35 41 / 5%);
        position: relative;
        z-index: 13;
        transition: all 0.2s ease-in-out;
    }

    .layout-sider-fix {
        position: fixed;
        top: 0;
        left: 0;
    }

    .ant-layout {
        overflow: hidden;
    }

    .layout-right-fix {
        overflow-x: hidden;
        padding-left: 200px;
        min-height: 100vh;
        transition: all 0.2s ease-in-out;
    }

    .layout-content {
        flex: auto;
        min-height: 100vh;
    }

    .n-layout-header.n-layout-header--absolute-positioned {
        z-index: 11;
    }

    .n-layout-footer {
        background: none;
    }
}

.layout-content-main {
    margin: 0 10px 10px;
    position: relative;
    padding-top: 64px;
}

.layout-content-main-fix {
    padding-top: 64px;
}

.fluid-header {
    padding-top: 0px;
}
</style>
