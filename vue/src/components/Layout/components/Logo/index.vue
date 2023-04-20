<template>
    <div class="logo">
        <img src="~@/assets/logo.svg" alt="" :class="{ 'mr-2': !collapsed }" />
        <h2 v-show="!collapsed" class="title">BLOG</h2>
    </div>
</template>

<script lang="ts" setup>
defineProps({ collapsed: Boolean });
</script>
<style lang="less" scoped>
.logo {
    display: flex;
    align-items: center;
    justify-content: center;
    height: 64px;
    line-height: 64px;
    overflow: hidden;
    white-space: nowrap;
    img {
        width: auto;
        height: 32px;
    }
    .title {
        margin: 0;
        font-size: inherit;
        font-weight: inherit;
    }
}
</style>
