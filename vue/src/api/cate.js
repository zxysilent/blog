import fetch from "./fetch";
// 分类信息
export const apiCateAll = () => {
	return fetch.request({
		url: "/api/cate/all",
		method: "get"
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
export const admCateDrop = id => {
	return fetch.request({
		url: `/adm/cate/drop/${id}`,
		method: "get"
	});
};
