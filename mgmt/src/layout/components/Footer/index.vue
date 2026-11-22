<template>
  <div class="page-footer">
    <div class="links">
      <a href="https://github.com/zxysilent/blog" target="_blank">源码</a>
      <a href="https://golang.google.cn" target="_blank">Golang</a>
      <a href="https://www.naiveui.com" target="_blank">Naive UI</a>
    </div>
    <div class="copyright">blog · 一个 go 语言实现的博客系统</div>
  </div>
</template>

<script lang="ts" setup></script>

<style lang="less" scoped>
  .page-footer {
    padding: 20px 0;
    text-align: center;
    font-size: 13px;
    color: var(--text-color-3);
  }

  .links a {
    margin: 0 8px;
    color: var(--text-color-3);
    text-decoration: none;

    &:hover {
      color: var(--primary-color);
    }
  }

  .copyright {
    margin-top: 6px;
  }
</style>
