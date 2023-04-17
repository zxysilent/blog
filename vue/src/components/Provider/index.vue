<template></template>
<script lang="ts" setup>
import { useLoadingBar, useDialog, useMessage, useNotification } from "naive-ui";
window.$loadingBar = useLoadingBar();
window.$dialog = useDialog();
window.$message = useMessage();
window.$notify = useNotification();
console.log("全局消息组件注册完毕");
</script>
