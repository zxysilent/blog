import { createDiscreteApi, darkTheme, ConfigProviderProps } from 'naive-ui';
import { computed } from 'vue';
import setting from "@/utils/global";
/**
 * @type import('naive-ui').GlobalThemeOverrides
 */
const getThemeOverrides = computed(() => {
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
/**
 * 挂载 Naive-ui 脱离上下文的 API
 * 如果你想在 setup 外使用 useDialog、useMessage、useNotification、useLoadingBar，可以通过 createDiscreteApi 来构建对应的 API。
 *https://www.naiveui.com/zh-CN/light/components/discrete
 */
// { xs: 0, s: 640, m: 1024, l: 1280, xl: 1536, xxl: 1920 }
export function setupNaiveDiscreteApi() {
    const configProviderPropsRef = computed<ConfigProviderProps>(() => ({
        theme: setting.darkTheme ? darkTheme : undefined,
        themeOverrides: {
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
            },
            LoadingBar: {
                colorLoading: "#2d8cf0ff",
            },
        },
    }));
    const { loadingBar } = createDiscreteApi(
        ['loadingBar'],
        {
            configProviderProps: configProviderPropsRef,
        }
    );
    window['$loading'] = loadingBar;
}
