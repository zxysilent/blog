import ajax from "./base";
// 标签信息
export const apiTagAll = () => {
	return ajax.get("/api/tag/all");
};
// 添加标签
export const admTagAdd = data => {
	return ajax.post(`/adm/tag/add`, data);
};
// 修改标签
export const admTagEdit = data => {
	return ajax.post(`/adm/tag/edit`, data);
};
// 删除标签
export const admTagDrop = id => {
	return ajax.get(`/adm/tag/drop/${id}`);
};
