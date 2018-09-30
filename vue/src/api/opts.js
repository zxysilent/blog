import ajax from "./base";
// 获取基本配置项
export const optsBase = () => {
  return ajax.get(`/adm/opts/base`);
};
// 获取某个配置项
export const optsGet = key => {
  return ajax.get(`/adm/opts/${key}`);
};
// 编辑某个配置项
export const optsEdit = data => {
  return ajax.post("/adm/opts/edit", data);
};
