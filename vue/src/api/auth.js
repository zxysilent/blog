import fetch from "./fetch";

export const apiAuthLogin = data => {
	return fetch.post("/api/login", data);
};
export const apiAuthVcode = () => {
	return fetch.get("/api/vcode");
};
export const admAuth = () => {
	return fetch.get("/adm/auth");
};
export const apiAuthLogoff = data => {
	//   return fetch.post("/api/logoff", data);
};
// 统计状态
export const admCollect = () => {
	return fetch.get("/adm/collect");
};
// 服务器信息
export const admSys = () => {
	return fetch.get("/adm/sys");
};
