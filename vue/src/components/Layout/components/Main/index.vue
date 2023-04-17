<template>
    <RouterView>
        <template #default="{ Component, route }">
            <transition :name="getTransitionName" mode="out-in" appear>
                <component :is="Component" :key="route.fullPath" />
            </transition>
        </template>
    </RouterView>
</template>

<script lang="ts" setup>
import { computed } from "vue";
import setting from "@/utils/global";
const getTransitionName = computed(() => {
    return setting.isPageAnimate ? setting.pageAnimateType : "";
});
defineProps({
    notNeedKey: { type: Boolean, default: false },
    animate: { type: Boolean, default: true },
});
</script>
<style lang="less" scoped></style>
