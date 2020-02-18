import fetch from "./fetch";
// 分类文章列表
export const apiCatePost = (cls, data) => {
	return fetch.get(`/api/cate/post/${cls}`, { params: data });
};
// 页面信息
export const apiPageAll = () => {
	return fetch.get("/api/page/all");
};
// 一条信息
export const apiPostGet = id => {
	return fetch.get(`/api/post/get/${id}`);
};
// 一条文章的tag信息列表
export const apiPostTagGet = id => {
	return fetch.get(`/api/post/tag/get/${id}`);
};
//删除
export const admPostDrop = id => {
	return fetch.get(`/adm/post/drop/${id}`);
};
//文章操作
export const admPostOpts = data => {
	return fetch.post(`/adm/post/opts`, data);
};
