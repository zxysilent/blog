import fetch from "./fetch";
// 获取基本配置项
export const apiMenuTree = () => {
	// return fetch.get(`/api/menu/tree`);
    return fetch({
        method: 'get',
        url: "/api/menu/tree",
    })
};