import ajax from "./base";

export const login = data => {
	return ajax.post("/login", data);
};
export const auth = () => {
	return ajax.get("/api/auth");
};
export const logoff = data => {
	//   return ajax.post("/api/logoff", data);
};
// 统计状态
export const collect = () => {
	return ajax.get("/api/collect");
};
// 服务器信息
export const sys = () => {
	return ajax.get("/api/sys");
};
