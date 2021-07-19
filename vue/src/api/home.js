import fetch from "./fetch";
// 统计状态
export const admStatusAppinfo = () => {
	return fetch.request({
		url: "/adm/status/appinfo",
		method: "get"
	});
};
// 服务器信息
export const admStatusGoinfo = () => {
	return fetch.request({
		url: "/adm/status/goinfo",
		method: "get"
	});
};
// ------------------------------------------------------ 配置中心 ------------------------------------------------------
// 获取配置global
export const apiGlobalGet = () => {
	return fetch.get("/api/global/get");
};
// 修改配置global
export const admGlobalEdit = data => {
	return fetch.post("/adm/global/edit", data);
};
// 修改配置global
export const admGlobalEditUsage = data => {
	return fetch.post("/adm/global/edit/usage", data);
};
// 清空数据
export const admGlobalClear = data => {
	return fetch.post("/adm/global/clear", data);
};
