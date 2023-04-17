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
export default setting;
