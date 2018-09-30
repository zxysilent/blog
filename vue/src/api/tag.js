import ajax from "./base";
// 标签信息
export const tagAll = () => {
  return ajax.get("/adm/tag/all");
};
// 添加标签
export const tagAdd = data => {
  return ajax.post(`/adm/tag/add`, data);
};
// 修改标签
export const tagEdit = data => {
  return ajax.post(`/adm/tag/edit`, data);
};
// 删除标签
export const tagDel = id => {
  return ajax.get(`/adm/tag/del/${id}`);
};
