<template>
    <n-card :bordered="false" class="proCard">
        <n-tabs type="line" default-value="basic" size="large">
            <n-tab-pane name="basic" tab="基本设置"><BasicSetting /> </n-tab-pane>
            <n-tab-pane name="safety" tab="安全设置"> <SafetySetting /> </n-tab-pane>
        </n-tabs>
    </n-card>
</template>
<script lang="ts" setup>
import BasicSetting from "./BasicSetting.vue";
import SafetySetting from "./SafetySetting.vue";
</script>
<style lang="less" scoped></style>
