import ajax from "./base";
// 分类文章列表
export const catePost = (cls, data) => {
  return ajax.get(`/api/cate/post/${cls}`, { params: data });
};
// 页面信息
export const pageAll = () => {
  return ajax.get("/api/page/all");
};
// 一条信息
export const postGet = id => {
  return ajax.get(`/api/post/get/${id}`);
};
// 一条文章的tag信息列表
export const postTagIds = id => {
  return ajax.get(`/api/post/tag/ids/${id}`);
};
//删除
export const postDel = id => {
  return ajax.get(`/api/post/del/${id}`);
};
//文章操作
export const postOpts = data => {
  return ajax.post(`/api/post/opts`, data);
};
