import fetch from "./fetch";
// 获取基本配置项
export const apiOptsBase = () => {
	return fetch.request({
		url: "/api/opts/base",
		method: "get"
	});
};
// 获取某个配置项
export const apiOptsGet = key => {
	return fetch.request({
		url: `/api/opts/${key}`,
		method: "get"
	});
};
// 编辑某个配置项
export const admOptsEdit = data => {
	return fetch.request({
		url: "/adm/opts/edit",
		method: "post",
		data: data
	});
};
