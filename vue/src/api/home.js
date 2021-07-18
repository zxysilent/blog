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
