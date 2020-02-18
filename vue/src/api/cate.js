import fetch from "./fetch";
// 分类信息
export const apiCateAll = () => {
	return fetch.get("/api/cate/all");
};
// 添加分类
export const admCateAdd = data => {
	return fetch.post(`/adm/cate/add`, data);
};
// 修改分类
export const admCateEdit = data => {
	return fetch.post(`/adm/cate/edit`, data);
};
// 删除分类
export const admCateDrop = id => {
	return fetch.get(`/adm/cate/drop/${id}`);
};
