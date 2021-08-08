import fetch from "./fetch";
// 通过id获取单条字典
export const apiDictGet = data => {
	return fetch.request({
		url: "/api/dict/get",
		method: "get",
		params: data
	});
};
// 获取字典分页
export const apiDictPage = data => {
	return fetch.request({
		url: "/api/dict/page",
		method: "get",
		params: data
	});
};
// 添加字典
export const admDictAdd = data => {
	return fetch.request({
		url: "/adm/dict/add",
		method: "post",
		data: data
	});
};
// 修改字典
export const admDictEdit = data => {
	return fetch.request({
		url: "/adm/dict/edit",
		method: "post",
		data: data
	});
};
// 通过id删除单条字典
export const admDictDrop = data => {
	return fetch.request({
		url: "/adm/dict/drop",
		method: "post",
		data: data
	});
};
