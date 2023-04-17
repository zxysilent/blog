<template>
    <n-card :bordered="false" class="proCard">
        <n-tabs type="line" default-value="basic" size="large">
            <n-tab-pane name="basic" tab="基本设置"><BasicSetting /> </n-tab-pane>
            <n-tab-pane name="reveal" tab="显示设置"> <RevealSetting /> </n-tab-pane>
            <n-tab-pane name="email" tab="邮件设置"> <EmailSetting /> </n-tab-pane>
        </n-tabs>
    </n-card>
</template>
<script lang="ts" setup>
import BasicSetting from "./BasicSetting.vue";
import RevealSetting from "./RevealSetting.vue";
import EmailSetting from "./EmailSetting.vue";
</script>
<style lang="less" scoped></style>
