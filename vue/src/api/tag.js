import fetch from "./fetch";
// 标签信息
export const apiTagAll = () => {
	return fetch.get("/api/tag/all");
};
// 添加标签
export const admTagAdd = data => {
	return fetch.post(`/adm/tag/add`, data);
};
// 修改标签
export const admTagEdit = data => {
	return fetch.post(`/adm/tag/edit`, data);
};
// 删除标签
export const admTagDrop = id => {
	return fetch.get(`/adm/tag/drop/${id}`);
};
