import ajax from "./base";
// 分类信息
export const apiCateAll = () => {
	return ajax.get("/api/cate/all");
};
// 添加分类
export const admCateAdd = data => {
	return ajax.post(`/adm/cate/add`, data);
};
// 修改分类
export const admCateEdit = data => {
	return ajax.post(`/adm/cate/edit`, data);
};
// 删除分类
export const admCateDrop = id => {
	return ajax.get(`/adm/cate/drop/${id}`);
};
