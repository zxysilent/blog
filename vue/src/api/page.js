import fetch from "./fetch";
// 通过id获取单条页面
export const apiPageGet = (data) => {
    return fetch.request({
		url: "/api/page/get",
		method: "get",
		params: data,
	});
};
// 获取所有页面
export const apiPageAll = (data) => {
    return fetch.request({
		url: "/api/page/all",
		method: "get",
		params: data,
	});
};
// 获取页面分页
export const apiPagePage = (data) => {
    return fetch.request({
		url: "/api/page/page",
		method: "get",
		params: data,
	});
};
// 添加页面
export const admPageAdd = (data) => {
    return fetch.request({
		url: "/adm/page/add",
		method: "post",
		data: data,
	});
};
// 修改页面
export const admPageEdit = (data) => {
	return fetch.request({
		url: "/adm/page/edit",
		method: "post",
		data: data,
	});
};
// 通过id删除单条页面
export const admPageDrop = (data) => {
	return fetch.request({
		url: "/adm/page/drop",
		method: "post",
		data: data,
	});
};
