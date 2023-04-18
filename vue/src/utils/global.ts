import { computed } from 'vue';
import { createDiscreteApi, darkTheme, ConfigProviderProps, GlobalThemeOverrides } from 'naive-ui';
const setting = {
    //导航模式 vertical 左侧菜单模式;horizontal 顶部菜单模式
    navMode: "vertical",
    //导航风格 dark 暗色侧边栏;light 白色侧边栏;header-dark 暗色侧边栏+顶栏
    navTheme: "dark",
    //深色主题
    darkTheme: false,
    //顶部
    headerSetting: {
        //固定顶部
        fixed: true,
    },
    //页脚
    showFooter: false,
    //菜单
    menuSetting: {
        //最小宽度
        minMenuWidth: 64,
        //菜单宽度
        menuWidth: 200,
        //固定菜单
        fixed: true,
    },
    //是否开启路由动画
    isPageAnimate: true,
    //路由动画类型
    pageAnimateType: "zoom-fade",
};
// const animates = [
//     { value: 'zoom-fade', label: '渐变' },
//     { value: 'zoom-out', label: '闪现' },
//     { value: 'fade-slide', label: '滑动' },
//     { value: 'fade', label: '消退' },
//     { value: 'fade-bottom', label: '底部消退' },
//     { value: 'fade-scale', label: '缩放消退' },
// ];

/**
 * @type import('naive-ui').GlobalThemeOverrides
 */
const themeOverrides = computed<GlobalThemeOverrides>(() => {
    return {
        common: {
            primaryColor: "#2d8cf0ff",
            primaryColorHover: "#57a3f3ff",
            primaryColorPressed: "#2b85e4ff",
            primaryColorSuppl: "#2D8CF0ff",
            infoColorPressed: "#2baee9ff",
            infoColorHover: "#57c5f7ff",
            infoColor: "#2db7f5ff",
            infoColorSuppl: "#2DB7F5ff",
            successColor: "#19be6bff",
            successColorPressed: "#18b566ff",
            successColorHover: "#47cb89ff",
            successColorSuppl: "#19BE6Bff",
            // warningColorPressed: "#f29100ff",
            // warningColorSuppl: "#fcb040ff",
            // warningColorHover: "#ffad33ff",
            // warningColor: "#ff9900ff",
            // errorColorPressed: "#e13d13ff",
            // errorColorHover: "#f16643ff",
            // errorColorSuppl: "#ed4014ff",
            // errorColor: "#ed4014ff",
        },
        LoadingBar: {
            colorLoading: "#2d8cf0ff",
        },
    };
});
// { xs: 0, s: 640, m: 1024, l: 1280, xl: 1536, xxl: 1920 }
const theme = computed(() => (setting.darkTheme ? darkTheme : undefined));
/**
 * 挂载 Naive-ui 脱离上下文的 API
 * 如果你想在 setup 外使用 useDialog、useMessage、useNotification、useLoadingBar，可以通过 createDiscreteApi 来构建对应的 API。
 *https://www.naiveui.com/zh-CN/light/components/discrete
 */
const configProviderPropsRef = computed<ConfigProviderProps>(() => {
    return {
        theme: theme.value,
        themeOverrides: themeOverrides.value,
    }
});
const discrete = () => {
    const { loadingBar } = createDiscreteApi(
        ['loadingBar'],
        { configProviderProps: configProviderPropsRef }
    );
    window.$loadingBar = loadingBar;
}
export { discrete, theme, themeOverrides };
export default setting;
