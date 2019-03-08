import ajax from "./base";
// 获取基本配置项
export const optsBase = () => {
  return ajax.get(`/api/opts/base`);
};
// 获取某个配置项
export const optsGet = key => {
  return ajax.get(`/api/opts/${key}`);
};
// 编辑某个配置项
export const optsEdit = data => {
  return ajax.post("/api/opts/edit", data);
};
