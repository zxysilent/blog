import { h } from 'vue';
import { NIcon, NTag } from 'naive-ui';
/**
 * render 图标
 */
export const renderIcon = (icon) => {
    return () => h(NIcon, null, { default: () => h(icon) });
}

/**
 * render new Tag
 */
const newTagColors = { color: '#f90', textColor: '#fff', borderColor: '#f90' };
export const renderNew = (type = 'warning', text = 'New', color: object = newTagColors) => {
    return () =>
        h(
            NTag as any,
            { type, round: true, size: 'small', color, }, { default: () => text }
        );
}

const excludes = ['/:path(.*)*', '/', '/redirect', '/login']
const filterRouter = (routers: Array<any>) => {
    return routers.filter((item) => {
        return !item.meta?.hidden && !excludes.includes(item.path)
    });
}
const visibleRouter = (routers: Array<any>) => {
    let visbles = 0;
    if (!routers) {
        return false;
    }
    routers.forEach((item) => {
        (!item.meta?.hidden) && visbles++;
    })
    return visbles == 1;
}
/**
 * 递归组装菜单格式
 */
export const generatorMenu = (routers: Array<any>) => {
    return filterRouter(routers).map((item) => {
        const isRoot = visibleRouter(item.children)
        const info = isRoot ? item.children[0] : item;
        const currentMenu = {
            ...info,
            ...info.meta,
            label: info.meta?.title,
            key: info.name,
            icon: isRoot ? item.meta?.icon : info.meta?.icon,
        };
        // 是否有子菜单，并递归处理
        if (info.children && info.children.length > 0) {
            // Recursion
            currentMenu.children = generatorMenu(info.children);
        } else {
            currentMenu.children = undefined;
        }
        return currentMenu;
    });
}



/**
 * 判断是否 url
 */
export const isUrl = (url: string) => {
    return /(^http|https:\/\/)/g.test(url);
}