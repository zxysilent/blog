import ajax from "./base";
// 分类文章列表
export const catePost = (cls, data) => {
  return ajax.get(`/adm/cate/post/${cls}`, { params: data });
};
// 页面信息
export const pageAll = () => {
  return ajax.get("/adm/page/all");
};
// 一条信息
export const postGet = id => {
  return ajax.get(`/adm/post/get/${id}`);
};
// 一条文章的tag信息列表
export const postTagIds = id => {
  return ajax.get(`/adm/post/tag/ids/${id}`);
};
//删除
export const postDel = id => {
  return ajax.get(`/adm/post/del/${id}`);
};
//文章操作
export const postOpts = data => {
  return ajax.post(`/adm/post/opts`, data);
};
