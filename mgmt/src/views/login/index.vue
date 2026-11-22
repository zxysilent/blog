<template>
  <!--
    登录页按后端配置 login.style 选择外观：
      mgmt 管理端风格（默认）—— 左侧品牌区 + 右侧登录卡片
      vue  旧版前台风格      —— 浅灰底 + 居中窄表单 + 底部链接
    配置读取完成前先铺一层空白底色，避免两套布局来回闪一下。
  -->
  <div v-if="!ready" class="login-loading" />
  <VueLogin
    v-else-if="site.style === 'vue'"
    :title="site.title"
    :desc="site.desc"
    :site-title="site.siteTitle"
    :logo="site.logo"
  />
  <MgmtLogin
    v-else
    :title="site.title"
    :desc="site.desc"
    :site-title="site.siteTitle"
    :logo="site.logo"
  />
</template>

<script lang="ts" setup>
  import { onMounted } from 'vue';
  import MgmtLogin from './MgmtLogin.vue';
  import VueLogin from './VueLogin.vue';
  import { useLoginSite } from './useLogin';

  const { site, ready, loadSite } = useLoginSite();

  onMounted(loadSite);
</script>

<style lang="less" scoped>
  .login-loading {
    min-height: 100vh;
    background: #f7f8fa;
  }
</style>
