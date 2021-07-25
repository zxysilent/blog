import fetch from "./fetch";
// 通过id获取单条分类
export const apiCateGet = (data) => {
    return fetch.request({
		url: "/api/cate/get",
		method: "get",
		params: data,
	});
};
// 所有分类
export const apiCateAll = () => {
	return fetch.request({
		url: "/api/cate/all",
		method: "get"
	});
};
// 分类分页
export const apiCatePage = data => {
	return fetch.request({
		url: "/api/cate/page",
		method: "get",
		params: data
	});
};
// 添加分类
export const admCateAdd = data => {
	return fetch.request({
		url: "/adm/cate/add",
		method: "post",
		data: data
	});
};
// 修改分类
export const admCateEdit = data => {
	return fetch.request({
		url: "/adm/cate/edit",
		method: "post",
		data: data
	});
};
// 删除分类
export const admCateDrop = data => {
	return fetch.request({
		url: `/adm/cate/drop`,
		method: "post",
		data: data
	});
};
