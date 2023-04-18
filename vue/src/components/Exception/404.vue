<template>
    <div class="flex flex-col justify-center page-container">
        <div class="text-center">
            <img src="~@/assets/images/exception/404.svg" alt="" />
        </div>
        <div class="text-center">
            <h1 class="text-base text-gray-500">抱歉，你访问的页面不存在</h1>
            <n-space justify="center">
                <n-button type="info" @click="$router.push({ name: 'home' })">回到首页</n-button>
                <n-button type="success" @click="$router.push({ name: 'login' })">登录</n-button>
            </n-space>
        </div>
    </div>
</template>
<script lang="ts" setup></script>
<style lang="less" scoped>
.page-container {
    width: 100%;
    border-radius: 4px;
    padding: 50px 0;
    height: 100vh;
    .text-center {
        h1 {
            color: #666;
            padding: 20px 0;
        }
    }
    img {
        width: 520px;
        margin: 0 auto;
    }
}
</style>
