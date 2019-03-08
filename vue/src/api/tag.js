import ajax from "./base";
// 标签信息
export const tagAll = () => {
  return ajax.get("/api/tag/all");
};
// 添加标签
export const tagAdd = data => {
  return ajax.post(`/api/tag/add`, data);
};
// 修改标签
export const tagEdit = data => {
  return ajax.post(`/api/tag/edit`, data);
};
// 删除标签
export const tagDel = id => {
  return ajax.get(`/api/tag/del/${id}`);
};
