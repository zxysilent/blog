<template>
    <n-loading-bar-provider>
        <n-dialog-provider>
            <n-notification-provider>
                <n-message-provider>
                    <Provider />
                    <slot slot="default"></slot>
                </n-message-provider>
            </n-notification-provider>
        </n-dialog-provider>
    </n-loading-bar-provider>
</template>
<script lang="ts" setup>
import { NDialogProvider, NNotificationProvider, NMessageProvider, NLoadingBarProvider } from "naive-ui";
import { Provider } from "@/components/Provider";
</script>
