import ajax from "./base";
// 分类文章列表
export const apiCatePost = (cls, data) => {
	return ajax.get(`/api/cate/post/${cls}`, { params: data });
};
// 页面信息
export const apiPageAll = () => {
	return ajax.get("/api/page/all");
};
// 一条信息
export const apiPostGet = id => {
	return ajax.get(`/api/post/get/${id}`);
};
// 一条文章的tag信息列表
export const apiPostTagGet = id => {
	return ajax.get(`/api/post/tag/get/${id}`);
};
//删除
export const admPostDrop = id => {
	return ajax.get(`/adm/post/drop/${id}`);
};
//文章操作
export const admPostOpts = data => {
	return ajax.post(`/adm/post/opts`, data);
};
