import fetch from "./fetch";
// 登录
export const apiAuthLogin = data => {
	return fetch.request({
		url: "/api/auth/login",
		method: "post",
		data: data
	});
};
// 获取验证码
export const apiAuthVcode = () => {
	return fetch.request({
		url: "/api/auth/vcode",
		method: "get"
	});
};
// 获取登录信息
export const admAuthGet = () => {
	return fetch.request({
		url: "/adm/auth/get",
		method: "get"
	});
};
// 修改自己信息
export const admAuthEdit = data => {
	return fetch.request({
		url: "/adm/auth/edit",
		method: "post",
		data: data
	});
};
// 修改自己密码
export const admAuthPasswd = data => {
	return fetch.request({
		url: "/adm/auth/passwd",
		method: "post",
		data: data
	});
};

export const apiAuthLogoff = data => {
	//   return fetch.post("/api/logoff", data);
};
