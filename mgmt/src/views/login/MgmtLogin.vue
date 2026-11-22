<template>
  <div class="login-page" :style="cssVars">
    <!-- 左侧品牌区：静态渐变 + 网格纹理 -->
    <aside class="login-page__brand">
      <div class="brand">
        <div class="flex items-center">
          <img class="brand__logo" :src="logo" :alt="siteTitle" />
          <h1 class="brand__title">{{ siteTitle }}</h1>
        </div>
        <p class="brand__desc">{{ desc || 'Go 语言实现的博客系统，支持多用户与权限管理' }}</p>
        <ul class="brand__features">
          <li v-for="feature in features" :key="feature" class="brand__feature">
            <span class="brand__dot" />
            <span>{{ feature }}</span>
          </li>
        </ul>
      </div>
      <p class="brand__copyright">© {{ currentYear }} {{ siteTitle }} · Go + Echo + XORM</p>
    </aside>

    <!-- 右侧表单区 -->
    <main class="login-page__main">
      <div class="login-card">
        <h2 class="login-card__title">账号登录</h2>
        <p class="login-card__subtitle">{{ title || '欢迎回来，请登录您的账号' }}</p>

        <n-form
          ref="formRef"
          class="login-card__form"
          :model="formInline"
          :rules="rules"
          size="large"
          @keydown.enter="handleSubmit"
        >
          <n-form-item path="username">
            <n-input
              ref="usernameRef"
              v-model:value="formInline.username"
              placeholder="请输入账号"
              clearable
            >
              <template #prefix>
                <n-icon size="18">
                  <PersonOutline />
                </n-icon>
              </template>
            </n-input>
          </n-form-item>

          <n-form-item path="password">
            <n-input
              v-model:value="formInline.password"
              type="password"
              show-password-on="click"
              placeholder="请输入密码"
            >
              <template #prefix>
                <n-icon size="18">
                  <LockClosedOutline />
                </n-icon>
              </template>
            </n-input>
          </n-form-item>

          <!-- 图形验证码：后端登录接口强制校验，点击图片可刷新 -->
          <n-form-item path="vcode">
            <div class="login-card__vcode">
              <n-input v-model:value="formInline.vcode" placeholder="请输入验证码" maxlength="6">
                <template #prefix>
                  <n-icon size="18">
                    <ShieldCheckmarkOutline />
                  </n-icon>
                </template>
              </n-input>
              <img
                v-if="captcha.image"
                class="login-card__captcha"
                :src="'data:image/png;base64,' + captcha.image"
                title="点击刷新验证码"
                alt="captcha"
                @click="refreshCaptcha"
              />
            </div>
          </n-form-item>

          <div class="login-card__row">
            <n-checkbox v-model:checked="rememberMe">记住我</n-checkbox>
            <a class="login-card__link" href="javascript:" @click="refreshCaptcha">看不清？换一张</a>
          </div>

          <n-button type="primary" size="large" block :loading="loading" @click="handleSubmit">
            登录
          </n-button>
        </n-form>

        <p class="login-card__footer">默认账号 root，密码见部署文档</p>
      </div>
    </main>
  </div>
</template>

<script lang="ts" setup>
  /**
   * 管理端风格登录页：左侧品牌区 + 右侧登录卡片
   *
   * 由 config 里的 login.style = mgmt 选中，是默认外观。
   * 表单逻辑在 useLogin.ts，本组件只负责布局与皮肤。
   */
  import { computed, onMounted, ref } from 'vue';
  import { useThemeVars } from 'naive-ui';
  import type { InputInst } from 'naive-ui';
  import { LockClosedOutline, PersonOutline, ShieldCheckmarkOutline } from '@vicons/ionicons5';
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

  const { formRef, loading, rememberMe, formInline, captcha, rules, refreshCaptcha, handleSubmit } =
    useLoginForm();

  const usernameRef = ref<InputInst | null>(null);
  const currentYear = new Date().getFullYear();

  const features = [
    '博文 / 页面 / 笔记 统一内容管理',
    '基于路由的权限标识自动同步',
    '用户、角色、栏目、标签全量管理',
    '配置与字典在线调整，免重启',
  ];

  // 主题变量透传给 CSS，登录页跟随系统主题色与深色模式
  const themeVars = useThemeVars();
  const cssVars = computed<Record<string, string>>(() => {
    const v = themeVars.value;
    return {
      '--login-bg': v.bodyColor,
      '--login-surface': v.cardColor,
      '--login-border': v.dividerColor,
      '--login-text': v.textColor1,
      '--login-text-secondary': v.textColor3,
      '--login-primary': v.primaryColor,
      '--login-primary-deep': v.primaryColorPressed,
    };
  });

  onMounted(() => {
    refreshCaptcha();
    usernameRef.value?.focus();
  });
</script>

<style lang="less" scoped>
  .login-page {
    display: flex;
    min-height: 100vh;
    background: var(--login-bg);
    color: var(--login-text);

    &__brand {
      display: none;
      flex: 1 1 45%;
      flex-direction: column;
      justify-content: space-between;
      padding: 64px 56px 40px;
      overflow: hidden;
      color: #fff;
      background:
        linear-gradient(rgb(255 255 255 / 6%) 1px, transparent 1px) 0 0 / 40px 40px,
        linear-gradient(90deg, rgb(255 255 255 / 6%) 1px, transparent 1px) 0 0 / 40px 40px,
        linear-gradient(150deg, #14213d 0%, var(--login-primary-deep) 55%, var(--login-primary) 100%);
    }

    &__main {
      display: flex;
      flex: 1 1 55%;
      align-items: center;
      justify-content: center;
      padding: 48px 24px;
    }
  }

  @media (min-width: 900px) {
    .login-page__brand {
      display: flex;
    }
  }

  .brand {
    max-width: 420px;

    &__logo {
      display: block;
      width: 32px;
      margin-right: 10px;
    }

    &__title {
      font-size: 30px;
      font-weight: 600;
      letter-spacing: 0.5px;
    }

    &__desc {
      margin: 0 0 40px;
      font-size: 15px;
      line-height: 1.7;
      color: rgb(255 255 255 / 72%);
    }

    &__features {
      padding: 0;
      margin: 0;
      list-style: none;
    }

    &__feature {
      display: flex;
      align-items: center;
      margin-bottom: 14px;
      font-size: 14px;
      color: rgb(255 255 255 / 82%);
    }

    &__dot {
      width: 6px;
      height: 6px;
      margin-right: 10px;
      background: rgb(255 255 255 / 85%);
      border-radius: 50%;
    }

    &__copyright {
      margin: 0;
      font-size: 13px;
      color: rgb(255 255 255 / 55%);
    }
  }

  .login-card {
    width: 100%;
    max-width: 400px;
    padding: 40px 36px 32px;
    background: var(--login-surface);
    border: 1px solid var(--login-border);
    border-radius: 12px;
    box-shadow: 0 8px 24px rgb(15 23 42 / 6%);

    &__title {
      margin: 0 0 8px;
      font-size: 24px;
      font-weight: 600;
      color: var(--login-text);
    }

    &__subtitle {
      margin: 0 0 32px;
      font-size: 14px;
      color: var(--login-text-secondary);
    }

    // 把输入框与按钮统一抬到 48px：验证码原图是 200x80（1:2.5），
    // 48px 高对应 120x48，正好等比铺满不裁切，也与右侧图片严格对齐。
    &__form {
      :deep(.n-input),
      :deep(.n-button) {
        --n-height: 48px !important;
      }

      :deep(.n-input__input-el) {
        font-size: 15px;
      }
    }

    // 验证码行：输入框与图片拼成一个整体控件（无间隙，边框合并）
    &__vcode {
      display: flex;
      align-items: stretch;
      width: 100%;

      :deep(.n-input) {
        flex: 1 1 auto;
        min-width: 0;
        // 只保留左侧圆角，右侧与验证码图片无缝拼接
        --n-border-radius: 6px 0 0 6px !important;
      }
    }

    // 验证码图片：紧贴输入框右侧，靠负边距吃掉重叠的 1px 边框
    &__captcha {
      flex: 0 0 auto;
      display: block;
      width: 120px;
      height: 48px;
      margin-left: -1px;
      cursor: pointer;
      border: 1px solid var(--login-border);
      border-radius: 0 6px 6px 0;
      object-fit: cover;
      background: #fff;
      transition: border-color 0.2s;

      &:hover {
        border-color: var(--login-primary);
      }
    }

    &__row {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin: 4px 0 20px;
    }

    &__link {
      font-size: 14px;
      color: var(--login-primary);
      text-decoration: none;

      &:hover {
        opacity: 0.8;
      }
    }

    &__footer {
      margin: 24px 0 0;
      font-size: 13px;
      text-align: center;
      color: var(--login-text-secondary);
    }
  }
</style>
