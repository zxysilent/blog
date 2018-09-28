import ajax from "./base";
// 分类信息
export const cateAll = () => {
  return ajax.get("/adm/cate/all");
};
