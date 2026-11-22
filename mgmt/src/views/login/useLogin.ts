/**
 * 登录页共享逻辑
 *
 * 登录页有两套外观（mgmt / vue），但表单、验证码、提交与跳转完全一致，
 * 因此逻辑集中在这里，两套外观只负责布局与皮肤，避免改一处漏一处。
 */
import { computed, reactive, ref } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useMessage } from 'naive-ui';
import type { FormInst, FormRules } from 'naive-ui';
import { ResultEnum } from '@/enums/httpEnum';
import { PageEnum } from '@/enums/pageEnum';
import { websiteConfig } from '@/config/website.config';
import { useUserStore } from '@/store/user';
import { getCaptcha } from '@/api/auth';
import { getSite } from '@/api/site';

/** 登录页风格：与后端 model.LoginStyle* 保持一致 */
export type LoginStyle = 'mgmt' | 'vue';

export interface FormState {
  username: string;
  password: string;
  vcode: string;
}

/** 登录页可配置项：标题、描述、风格、Logo */
export interface LoginSite {
  style: LoginStyle;
  /** 登录页标题（config: login.title） */
  title: string;
  /** 登录页描述（config: login.desc） */
  desc: string;
  /** 站点标题（config: basic.title），品牌区展示 */
  siteTitle: string;
  /** Logo 地址，空值回落到内置图 */
  logo: string;
}

/** 表单状态与提交逻辑 */
export function useLoginForm() {
  const message = useMessage();
  const router = useRouter();
  const route = useRoute();
  const userStore = useUserStore();

  const formRef = ref<FormInst | null>(null);
  const loading = ref(false);
  const rememberMe = ref(true);

  const formInline = reactive<FormState>({
    username: 'root',
    password: '',
    vcode: '',
  });

  // 图形验证码：image 为 base64，real 为服务端签名，登录时原样回传
  const captcha = reactive({ image: '', real: '' });

  const rules = computed<FormRules>(() => ({
    username: { required: true, message: '请输入账号', trigger: 'blur' },
    password: { required: true, message: '请输入密码', trigger: 'blur' },
    vcode: { required: true, message: '请输入验证码', trigger: 'blur' },
  }));

  async function refreshCaptcha() {
    try {
      const res = await getCaptcha();
      captcha.image = res?.image || '';
      captcha.real = res?.real || '';
    } catch {
      // 验证码获取失败不阻断页面，用户可点击重试
      captcha.image = '';
      captcha.real = '';
    }
    formInline.vcode = '';
  }

  async function handleSubmit() {
    if (loading.value) return;
    try {
      await formRef.value?.validate();
    } catch {
      return;
    }

    loading.value = true;
    message.loading('登录中...');
    try {
      const { code, msg } = await userStore.login({
        num: formInline.username,
        passwd: formInline.password,
        vcode: formInline.vcode,
        vreal: captcha.real,
      });
      if (code === ResultEnum.SUCCESS) {
        const toPath = decodeURIComponent((route.query?.redirect as string) || '/');
        message.success('登录成功，即将进入系统');
        router.replace(route.name === PageEnum.BASE_LOGIN_NAME ? '/' : toPath);
      } else {
        message.info(msg || '登录失败');
        // 验证码一次性有效，失败后必须换一张
        refreshCaptcha();
      }
    } catch {
      message.error('登录失败，请稍后重试');
      refreshCaptcha();
    } finally {
      message.destroyAll();
      loading.value = false;
    }
  }

  return {
    message,
    formRef,
    loading,
    rememberMe,
    formInline,
    captcha,
    rules,
    refreshCaptcha,
    handleSubmit,
  };
}

/**
 * 读取登录页配置
 *
 * `/api/site` 是公开接口，无需登录即可读取 basic / login 两组配置，
 * 登录页据此决定渲染哪套外观。读取失败时静默回落到默认值：
 * 登录本身不能因为一个展示配置不可用而失效。
 */
export function useLoginSite() {
  const site = reactive<LoginSite>({
    style: 'mgmt',
    title: '欢迎登录',
    desc: '便捷、高效、可靠',
    siteTitle: websiteConfig.title,
    logo: websiteConfig.logo,
  });
  const ready = ref(false);

  async function loadSite() {
    try {
      const res = await getSite();
      const login = (res?.login || {}) as Recordable;
      const basic = (res?.basic || {}) as Recordable;
      if (login.style === 'vue' || login.style === 'mgmt') site.style = login.style;
      if (login.title) site.title = login.title;
      if (login.desc) site.desc = login.desc;
      if (basic.title) site.siteTitle = basic.title;
      if (basic.logo_url) site.logo = basic.logo_url;
    } catch {
      // 后端不可用时保持默认（mgmt 风格 + 内置文案）
    } finally {
      ready.value = true;
    }
  }

  return { site, ready, loadSite };
}
