import fetch from "./fetch";
// 通过id获取单条文章
export const apiPostGet = data => {
	return fetch.request({
		url: "/api/post/get",
		method: "get",
		params: data
	});
};
// 获取所有文章
export const apiPostAll = data => {
	return fetch.request({
		url: "/api/post/all",
		method: "get",
		params: data
	});
};
// 获取文章分页
export const apiPostPage = data => {
	return fetch.request({
		url: "/api/post/page",
		method: "get",
		params: data
	});
};
// 添加文章
export const admPostAdd = data => {
	return fetch.request({
		url: "/adm/post/add",
		method: "post",
		data: data
	});
};
// 修改文章
export const admPostEdit = data => {
	return fetch.request({
		url: "/adm/post/edit",
		method: "post",
		data: data
	});
};
// 通过id删除单条文章
export const admPostDrop = data => {
	return fetch.request({
		url: "/adm/post/drop",
		method: "post",
		data: data
	});
};
