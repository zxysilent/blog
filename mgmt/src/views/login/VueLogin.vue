<template>
  <div class="legacy-login" :style="cssVars">
    <!-- 顶部品牌：Logo + 标题 + 描述，居中排布 -->
    <div class="legacy-login__top">
      <div class="legacy-login__header">
        <img class="legacy-login__logo" :src="logo" :alt="siteTitle" />
        <span class="legacy-login__title">{{ title || siteTitle }}</span>
      </div>
      <p class="legacy-login__desc">{{ desc }}</p>
    </div>

    <!-- 居中表单：368px 宽，左对齐标签 -->
    <div class="legacy-login__main">
      <n-form
        ref="formRef"
        :model="formInline"
        :rules="rules"
        size="large"
        label-placement="left"
        label-width="72"
        @keydown.enter="handleSubmit"
      >
        <n-form-item label="账号" path="username">
          <n-input v-model:value="formInline.username" placeholder="请输入账号">
            <template #prefix>
              <n-icon size="18" color="#808695">
                <PersonOutline />
              </n-icon>
            </template>
          </n-input>
        </n-form-item>

        <n-form-item label="密码" path="password">
          <n-input
            v-model:value="formInline.password"
            type="password"
            show-password-on="click"
            placeholder="请输入密码"
          >
            <template #prefix>
              <n-icon size="18" color="#808695">
                <LockClosedOutline />
              </n-icon>
            </template>
          </n-input>
        </n-form-item>

        <!-- 验证码与图片无缝拼接，图片点击可刷新 -->
        <n-form-item label="验证码" path="vcode">
          <div class="legacy-login__vcode">
            <n-input v-model:value="formInline.vcode" placeholder="请输入验证码" maxlength="6" />
            <img
              v-if="captcha.image"
              class="legacy-login__captcha"
              :src="'data:image/png;base64,' + captcha.image"
              title="点击刷新"
              alt="点击刷新"
              @click="refreshCaptcha"
            />
          </div>
        </n-form-item>

        <n-form-item>
          <n-button type="primary" attr-type="button" block :loading="loading" @click="handleSubmit">
            登&nbsp;&nbsp;&nbsp;&nbsp;录
          </n-button>
        </n-form-item>
      </n-form>

      <div class="legacy-login__tips">
        <span>系统推荐浏览器</span>
        <div class="legacy-login__browsers">
          <span class="legacy-login__browser"><LogoEdge /> Edge</span>
          <span class="legacy-login__browser"><LogoChrome /> Chrome</span>
          <span class="legacy-login__browser"><LogoFirefox /> Firefox</span>
        </div>
      </div>
    </div>

    <!-- 底部：站点链接与版权，固定在窗口底部 -->
    <div class="legacy-login__footer">
      <div class="legacy-login__links">
        <a href="/">首页</a>
        <a href="https://github.com/zxysilent" target="_blank">技术支持</a>
      </div>
      <div class="legacy-login__copyright">
        Copyright &copy; {{ currentYear }} {{ siteTitle }}
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
  /**
   * 旧版前台风格登录页：浅灰底 + 居中窄表单 + 底部链接
   *
   * 由 config 里的 login.style = vue 选中。
   * 表单逻辑在 useLogin.ts，本组件只负责布局与皮肤。
   */
  import { computed, onMounted } from 'vue';
  import { useThemeVars } from 'naive-ui';
  import {
    LockClosedOutline,
    PersonOutline,
    LogoEdge,
    LogoChrome,
    LogoFirefox,
  } from '@vicons/ionicons5';
  import { useLoginForm } from './useLogin';

  defineProps<{
    /** 登录页标题（config: login.title） */
    title?: string;
    /** 登录页描述（config: login.desc） */
    desc?: string;
    /** 站点标题 */
    siteTitle?: string;
    /** Logo 地址 */
    logo?: string;
  }>();

  const { formRef, loading, formInline, captcha, rules, refreshCaptcha, handleSubmit } =
    useLoginForm();

  const currentYear = new Date().getFullYear();

  // 只取主题色与文字色：旧版是固定浅色底，不跟随深色模式，避免观感走样
  const themeVars = useThemeVars();
  const cssVars = computed<Record<string, string>>(() => ({
    '--login-primary': themeVars.value.primaryColor,
  }));

  onMounted(refreshCaptcha);
</script>

<style lang="less" scoped>
  .legacy-login {
    position: relative;
    box-sizing: border-box;
    width: 100%;
    min-height: 100vh;
    padding: 110px 0 144px;
    overflow: hidden;
    background:
      radial-gradient(circle at 18% 12%, rgb(37 99 235 / 8%) 0%, transparent 45%),
      radial-gradient(circle at 82% 4%, rgb(14 165 233 / 10%) 0%, transparent 42%),
      #f0f2f5;

    &__top {
      text-align: center;
    }

    &__header {
      display: flex;
      height: 44px;
      line-height: 44px;
      align-items: center;
      justify-content: center;
    }

    &__logo {
      height: 44px;
      margin-right: 16px;
      border-radius: 50%;
      vertical-align: top;
    }

    &__title {
      position: relative;
      top: 2px;
      font-size: 30px;
      font-weight: 600;
      color: rgb(0 0 0 / 85%);
    }

    &__desc {
      margin: 12px 0 30px;
      font-size: 16px;
      color: rgb(0 0 0 / 45%);
    }

    &__main {
      width: 368px;
      min-width: 260px;
      margin: 0 auto;
    }

    // 验证码：与输入框拼成一体，不留缝隙
    &__vcode {
      display: flex;
      align-items: stretch;
      width: 100%;

      :deep(.n-input) {
        flex: 1 1 auto;
        min-width: 0;
        --n-border-radius: 6px 0 0 6px !important;
      }
    }

    // 原图 200x80，40px 高即 100x40，等比不裁切
    &__captcha {
      flex: 0 0 auto;
      display: block;
      width: 100px;
      height: 40px;
      margin-left: -1px;
      cursor: pointer;
      border: 1px solid #d9d9d9;
      border-radius: 0 6px 6px 0;
      object-fit: cover;
      background: #fff;

      &:hover {
        border-color: var(--login-primary);
      }
    }

    &__tips {
      margin-top: 10px;
      font-size: 16px;
      line-height: 22px;
      text-align: center;
      color: rgb(0 0 0 / 45%);
    }

    &__browsers {
      display: flex;
      align-items: center;
      justify-content: center;
      gap: 32px;
      margin-top: 12px;
    }

    &__browser {
      display: inline-flex;
      gap: 6px;
      align-items: center;
      font-size: 14px;
      color: rgb(0 0 0 / 45%);
    }

    &__footer {
      position: absolute;
      bottom: 0;
      width: 100%;
      padding: 0 16px;
      margin: 48px 0 24px;
      text-align: center;
    }

    &__links {
      margin-bottom: 8px;
      font-size: 14px;

      a {
        color: rgb(0 0 0 / 45%);
        text-decoration: none;
        transition: all 0.3s;

        &:not(:last-child) {
          margin-right: 40px;
        }

        &:hover {
          color: var(--login-primary);
        }
      }
    }

    &__copyright {
      font-size: 14px;
      color: rgb(0 0 0 / 45%);
    }
  }

  // 小屏时收窄外边距，避免 368px 表单把页面撑出横向滚动条
  @media (max-width: 420px) {
    .legacy-login__main {
      width: 100%;
      padding: 0 16px;
    }
  }
</style>
