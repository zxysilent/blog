import ajax from "./base";
// 分类信息
export const cateAll = () => {
  return ajax.get("/adm/cate/all");
};
// 添加分类
export const cateAdd = data => {
  return ajax.post(`/adm/cate/add`, data);
};
// 修改分类
export const cateEdit = data => {
  return ajax.post(`/adm/cate/edit`, data);
};
// 删除分类
export const cateDel = id => {
  return ajax.get(`/adm/cate/del/${id}`);
};
