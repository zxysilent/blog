import fetch from "./fetch";
// 获取所有标签
export const apiTagAll = (data) => {
    return fetch.request({
		url: "/api/tag/all",
		method: "get",
		params: data,
	});
};
// 获取标签分页
export const apiTagPage = (data) => {
    return fetch.request({
		url: "/api/tag/page",
		method: "get",
		params: data,
	});
};
// 添加标签
export const admTagAdd = (data) => {
    return fetch.request({
		url: "/adm/tag/add",
		method: "post",
		data: data,
	});
};
// 修改标签
export const admTagEdit = (data) => {
	return fetch.request({
		url: "/adm/tag/edit",
		method: "post",
		data: data,
	});
};
// 通过id删除单条标签
export const admTagDrop = (data) => {
	return fetch.request({
		url: "/adm/tag/drop",
		method: "post",
		data: data,
	});
};

