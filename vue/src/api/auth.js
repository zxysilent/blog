import fetch from "./fetch";

export const apiAuthLogin = data => {
	return fetch.post("/api/auth/login", data);
};

export const apiAuthVcode = () => {
	return fetch.get("/api/auth/vcode");
};

export const admAuth = () => {
	return fetch.get("/adm/auth");
};

export const admAuthMenu = data => {
	return fetch.get("/adm/auth/menu", { params: data });
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
