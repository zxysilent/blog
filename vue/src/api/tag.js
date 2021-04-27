import fetch from "./fetch";
// 标签信息
export const apiTagAll = () => {
	return fetch.request({
		url: "/api/tag/all",
		method: "get"
	});
};
// 添加标签
export const admTagAdd = data => {
	return fetch.request({
		url: "/adm/tag/add",
		method: "post",
		data: data
	});
};
// 修改标签
export const admTagEdit = data => {
	return fetch.request({
		url: "/adm/tag/edit",
		method: "post",
		data: data
	});
};
// 删除标签
export const admTagDrop = id => {
	return fetch.request({
		url: `/adm/tag/drop/${id}`,
		method: "get"
	});
};
