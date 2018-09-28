import ajax from "./base";

export const catePost = (cls, data) => {
  return ajax.get(`/adm/cate/post/${cls}`, { params: data });
};

// 页面信息
export const pageAll = () => {
  return ajax.get("/adm/page/all");
};

export const articleGet = id => {
  return ajax.get(`/news/article/${id}`);
};
export const articleChgtop = (id, val) => {
  return ajax.get(`/article/chgtop/${id}`, { params: { val: val } });
};
//删除
export const articleDel = id => {
  return ajax.get(`/article/del/${id}`);
};
//添加新闻
export const articleAdd = data => {
  return ajax.post(`/article/add`, data);
};
//修改新闻
export const articleEdit = data => {
  return ajax.post(`/article/edit`, data);
};
