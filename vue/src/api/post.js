import fetch from "./fetch";
// 分类文章列表
export const apiCatePost = (cls, data) => {
	return fetch.request({
		url: `/api/cate/post/${cls}`,
		method: "get",
		params: data
	});
};
// 页面信息
export const apiPageAll = () => {
	return fetch.request({
		url: "/api/page/all",
		method: "get"
	});
};
// 一条信息
export const apiPostGet = id => {
	return fetch.request({
		url: `/api/post/get/${id}`,
		method: "get"
	});
};
// 一条文章的tag信息列表
export const apiPostTagGet = id => {
	return fetch.request({
		url: `/api/post/tag/get/${id}`,
		method: "get"
	});
};
//删除
export const admPostDrop = id => {
	return fetch.request({
		url: `/adm/post/drop/${id}`,
		method: "get"
	});
};
//文章操作
export const admPostOpts = data => {
	return fetch.request({
		url: "/adm/post/opts",
		method: "post",
		data: data
	});
};
