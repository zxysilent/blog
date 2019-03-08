import ajax from "./base";
// 分类信息
export const cateAll = () => {
  return ajax.get("/api/cate/all");
};
// 添加分类
export const cateAdd = data => {
  return ajax.post(`/api/cate/add`, data);
};
// 修改分类
export const cateEdit = data => {
  return ajax.post(`/api/cate/edit`, data);
};
// 删除分类
export const cateDel = id => {
  return ajax.get(`/api/cate/del/${id}`);
};
