import ajax from "./base";

export const apiLogin = data => {
	return ajax.post("/api/login", data);
};
export const admAuth = () => {
	return ajax.get("/adm/auth");
};
export const apiLogoff = data => {
	//   return ajax.post("/api/logoff", data);
};
// 统计状态
export const admCollect = () => {
	return ajax.get("/adm/collect");
};
// 服务器信息
export const admSys = () => {
	return ajax.get("/adm/sys");
};
