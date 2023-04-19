<template>
    <!-- inline-theme-disabled -->
    <n-config-provider :theme="theme" :theme-overrides="themeOverrides" :locale="zhCN" :date-locale="dateZhCN" inline-theme-disabled>
        <n-dialog-provider>
            <n-notification-provider>
                <n-message-provider>
                    <RouterView />
                </n-message-provider>
            </n-notification-provider>
        </n-dialog-provider>
    </n-config-provider>
    <transition v-if="$route.name !== 'login'" name="slide-up"> </transition>
</template>
<script lang="ts" setup>
import { zhCN, dateZhCN } from "naive-ui";
import { NDialogProvider, NNotificationProvider, NMessageProvider } from "naive-ui";
import { theme, themeOverrides } from "@/utils/global";
</script>
<style lang="less">
@import "styles/index.less";
</style>
