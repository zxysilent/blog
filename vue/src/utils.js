import { Base64 } from "js-base64";

const storage = process.env.NODE_ENV === "development" ? localStorage : sessionStorage;

let utils = {};
utils.title = function(title) {
	title = title || "zadmin";
	window.document.title = title;
};

//保存数据
utils.setItem = (k, v) => {
	storage.setItem(k, v);
};
utils.getItem = k => {
	return storage.getItem(k);
};
utils.clearItem = k => {
	storage.removeItem(k);
};
utils.clearData = () => {
	storage.clear();
};
utils.setToken = token => {
	utils.setItem("bearer", token);
};
utils.getToken = () => {
	return utils.getItem("bearer");
};
// 获取保存的用户信息
utils.getAuth = () => {
	try {
		let token = Base64.decode(utils.getItem("bearer").split(".")[1]);
		let auth = JSON.parse(token);
		if (!auth.hasOwnProperty("id")) {
			utils.removeItem("bearer");
			location.href = "/login";
		}
		return auth;
	} catch (e) {
		utils.removeItem("bearer");
		location.href = "/login";
	}
};
// 不需要也可登录页面集合
utils.noAuth = r => {
	return ["login", "errjwt", "err401", "err50x", "err404"].indexOf(r) > -1;
};
utils.Role = {
	RSup: 30, //超级管理员
	RAtv: 20, //启用/禁用
	RBas: 10, //基本权限
	//判断指定位置权限
	getRole: (rl, r) => {
		if ((rl & (1 << r)) >> r == 1) {
			return true;
		}
		return false;
	}
};
utils.Role.isSup = rl => {
	return utils.Role.getRole(rl, utils.Role.RSup);
};
utils.Role.isAtv = rl => {
	return utils.Role.getRole(rl, utils.Role.RAtv);
};
// 权限路由相关
utils.Role.allow = (role, arr) => {
	for (let i = 0; i < arr.length; i++) {
		if (utils.Role.getRole(role, arr[i])) {
			return true;
		}
	}
	return false;
};
export default utils;
