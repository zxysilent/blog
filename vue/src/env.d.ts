/// <reference types="vite/client" />

declare module "*.vue" {
    import { DefineComponent } from "vue";
    // eslint-disable-next-line @typescript-eslint/no-explicit-any, @typescript-eslint/ban-types
    const component: DefineComponent<{}, {}, any>;
    export default component;
}

declare interface Window {
    $loadingBar: any;
    // $message: any;
    // $notify: any;
    // $dialog: any;
}
interface Meta {
    title: string;//标题
    auth: string;//权限标识
    icon?: string;//图标
    externalLink?: string;// 外链跳转地址
    hidden?: boolean;  //隐藏
    active?: string//
}



// 环境变量相关
interface ImportMetaEnv {
    readonly VITE_APP_SRV: string;
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}
// interface Menu {
//     title: string;
//     label: string;
//     key: string;
//     meta: Meta;
//     name: string;
//     children?: Menu[];
//     fullPath?: string;
//     icon?: any;
//     path: string;
//     permissions?: string[];
//     redirect?: string;
// }

// export interface ProjectSettingState {
//     //导航模式
//     navMode: string;
//     //导航风格
//     navTheme: string;
//     //页脚
//     showFooter: boolean;
//     //菜单设置
//     menuSetting: object;
// }


// export interface IheaderSetting {
//     fixed: boolean;
// }

// export interface ImenuSetting {
//     minMenuWidth: number;
//     menuWidth: number;
//     fixed: boolean;
// }


// export interface GlobConfig {
//     title: string;
//     apiUrl: string;
//     shortName: string;
//     urlPrefix?: string;
//     uploadUrl?: string;
//     prodMock: boolean;
//     imgUrl?: string;
// }

declare module "*.css";
declare module "*.svg";
declare module "*.png";
declare module "*.jpg";
declare module "*.jpeg";
declare module "*.gif";
declare module "*.bmp";
declare module "*.tiff";
