import { defineStore } from 'pinia';
import projectSetting from '@/settings/projectSetting';
import { getConfigById, saveConfigValue, setLoginStyle } from '@/api/config';
import type { IHeaderSetting, IMenuSetting, IMultiTabsSetting, ICrumbsSetting } from '/#/config';

const {
  navMode,
  navTheme,
  isMobile,
  headerSetting,
  showFooter,
  menuSetting,
  multiTabsSetting,
  crumbsSetting,
  permissionMode,
  isPageAnimate,
  pageAnimateType,
} = projectSetting;

interface ProjectSettingState {
  navMode: string; //导航模式
  navTheme: string; //导航风格
  /** 登录页风格 mgmt/vue：与导航风格同属外观开关，存于 layout 配置 */
  loginStyle: string;
  headerSetting: IHeaderSetting; //顶部设置
  showFooter: boolean; //页脚
  menuSetting: IMenuSetting; //多标签
  multiTabsSetting: IMultiTabsSetting; //多标签
  crumbsSetting: ICrumbsSetting; //面包屑
  permissionMode: string; //权限模式
  isPageAnimate: boolean; //是否开启路由动画
  pageAnimateType: string; //路由动画类型
  isMobile: boolean; // 是否处于移动端模式
  /** 是否已从后端同步过配置，避免重复请求 */
  loaded: boolean;
}

/**
 * 后端字段为 snake_case，前端状态为驼峰，集中在这两个映射函数里转换，
 * 新增字段时两边同时改，不用在组件里各写一遍。
 */
function toBackend(state: ProjectSettingState): Recordable {
  return {
    nav_mode: state.navMode,
    nav_theme: state.navTheme,
    // login_style 原样带回：它需要 login:style 权限，值没变时后端放行，
    // 避免整份配置保存时把风格字段丢掉
    login_style: state.loginStyle,
    header_fixed: state.headerSetting?.fixed ?? true,
    header_reload: state.headerSetting?.isReload ?? true,
    multi_tabs_show: state.multiTabsSetting?.show ?? true,
    multi_tabs_fixed: state.multiTabsSetting?.fixed ?? true,
    menu_fixed: state.menuSetting?.fixed ?? true,
    menu_mix: state.menuSetting?.mixMenu ?? false,
    menu_collapsed: state.menuSetting?.collapsed ?? false,
    menu_width: state.menuSetting?.menuWidth ?? 200,
    menu_min_width: state.menuSetting?.minMenuWidth ?? 64,
    crumb_show: state.crumbsSetting?.show ?? true,
    crumb_show_icon: state.crumbsSetting?.showIcon ?? false,
    show_footer: state.showFooter,
    page_animate: state.isPageAnimate,
    page_animate_type: state.pageAnimateType,
  };
}

/**
 * 把后端配置合并进状态
 *
 * 逐项判断类型：后端可能因为版本差异缺少某些字段，
 * 直接 Object.assign 会把缺失字段覆盖成 undefined，导致布局崩掉。
 */
function applyBackend(state: ProjectSettingState, raw: Recordable) {
  if (!raw || typeof raw !== 'object') return;
  if (typeof raw.nav_mode === 'string') state.navMode = raw.nav_mode;
  if (typeof raw.nav_theme === 'string') state.navTheme = raw.nav_theme;
  if (raw.login_style === 'mgmt' || raw.login_style === 'vue') state.loginStyle = raw.login_style;
  if (typeof raw.header_fixed === 'boolean') state.headerSetting.fixed = raw.header_fixed;
  if (typeof raw.header_reload === 'boolean') state.headerSetting.isReload = raw.header_reload;
  if (typeof raw.multi_tabs_show === 'boolean') state.multiTabsSetting.show = raw.multi_tabs_show;
  if (typeof raw.multi_tabs_fixed === 'boolean') state.multiTabsSetting.fixed = raw.multi_tabs_fixed;
  if (typeof raw.menu_fixed === 'boolean') state.menuSetting.fixed = raw.menu_fixed;
  if (typeof raw.menu_mix === 'boolean') state.menuSetting.mixMenu = raw.menu_mix;
  if (typeof raw.menu_collapsed === 'boolean') state.menuSetting.collapsed = raw.menu_collapsed;
  if (typeof raw.menu_width === 'number') state.menuSetting.menuWidth = raw.menu_width;
  if (typeof raw.menu_min_width === 'number') state.menuSetting.minMenuWidth = raw.menu_min_width;
  if (typeof raw.crumb_show === 'boolean') state.crumbsSetting.show = raw.crumb_show;
  if (typeof raw.crumb_show_icon === 'boolean') state.crumbsSetting.showIcon = raw.crumb_show_icon;
  if (typeof raw.show_footer === 'boolean') state.showFooter = raw.show_footer;
  if (typeof raw.page_animate === 'boolean') state.isPageAnimate = raw.page_animate;
  if (typeof raw.page_animate_type === 'string') state.pageAnimateType = raw.page_animate_type;
}

export const useProjectSettingStore = defineStore({
  id: 'app-project-setting',
  state: (): ProjectSettingState => ({
    navMode: navMode,
    navTheme,
    // 默认值与管理端风格一致；后端下发后会覆盖
    loginStyle: 'mgmt',
    isMobile,
    // 深拷贝：settings/projectSetting.ts 导出的对象是模块级单例，
    // 直接引用会让「加载失败回落默认值」拿到被上次会话改过的值
    headerSetting: { ...headerSetting },
    showFooter,
    menuSetting: { ...menuSetting },
    multiTabsSetting: { ...multiTabsSetting },
    crumbsSetting: { ...crumbsSetting },
    permissionMode,
    isPageAnimate,
    pageAnimateType,
    loaded: false,
  }),
  getters: {
    getNavMode(): string {
      return this.navMode;
    },
    getNavTheme(): string {
      return this.navTheme;
    },
    getLoginStyle(): string {
      return this.loginStyle;
    },
    getIsMobile(): boolean {
      return this.isMobile;
    },
    getHeaderSetting(): object {
      return this.headerSetting;
    },
    getShowFooter(): boolean {
      return this.showFooter;
    },
    getMenuSetting(): object {
      return this.menuSetting;
    },
    getMultiTabsSetting(): object {
      return this.multiTabsSetting;
    },
    getCrumbsSetting(): object {
      return this.crumbsSetting;
    },
    getPermissionMode(): string {
      return this.permissionMode;
    },
    getIsPageAnimate(): boolean {
      return this.isPageAnimate;
    },
    getPageAnimateType(): string {
      return this.pageAnimateType;
    },
  },
  actions: {
    setNavTheme(value: string): void {
      this.navTheme = value;
    },
    setIsMobile(value: boolean): void {
      this.isMobile = value;
    },

    /**
     * 从后端加载布局配置
     *
     * 配置存于 config 表（主键 layout），由独立接口 /api/configs/layout 下发；
     * 读取失败（未登录或后端异常）时保持前端默认值，
     * 管理端不能因为一个展示配置不可用而打不开。
     */
    async loadLayoutSetting(force = false) {
      if (this.loaded && !force) return;
      try {
        const mod = await getConfigById('layout');
        const raw = mod?.value ? JSON.parse(mod.value) : null;
        if (raw) applyBackend(this, raw);
        this.loaded = true;
      } catch {
        // 未登录或后端异常时沿用本地默认值
      }
    },

    /**
     * 当前配置的快照（后端字段序列化）
     * 用于判断状态是否真的变化，避免加载后立刻把同一个值写回去
     */
    snapshot(): string {
      return JSON.stringify(toBackend(this));
    },

    /**
     * 保存登录页风格
     * 风格存在 layout 配置里，但改动需要 login:style 权限，
     * 所以走专用接口而不是随整份布局一起提交
     */
    async saveLoginStyle(style: 'mgmt' | 'vue') {
      await setLoginStyle(style);
      this.loginStyle = style;
    },

    /**
     * 保存布局配置（整体替换）
     * 后端会做枚举与宽度区间校验，非法值直接返回错误
     */
    async saveLayoutSetting() {
      await saveConfigValue('layout', this.snapshot());
    },
  },
});
