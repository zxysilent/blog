import ajax from "./base";
// 获取基本配置项
export const apiOptsBase = () => {
	return ajax.get(`/api/opts/base`);
};
// 获取某个配置项
export const apiOptsGet = key => {
	return ajax.get(`/api/opts/${key}`);
};
// 编辑某个配置项
export const admOptsEdit = data => {
	return ajax.post("/adm/opts/edit", data);
};
