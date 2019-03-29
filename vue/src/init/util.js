import { Base64 } from "js-base64";
import { urlLogin } from "@/init/conf";
const storage = process.env.NODE_ENV === "development" ? localStorage : sessionStorage;

let util = {};
util.title = function(title) {
	title = title || "blog";
	window.document.title = title;
};

//保存数据
util.setItem = (k, v) => {
	storage.setItem(k, v);
};
util.getItem = k => {
	return storage.getItem(k);
};
util.clearItem = k => {
	storage.removeItem(k);
};
util.clearData = () => {
	storage.clear();
};
util.setToken = token => {
	util.setItem("bearer", token);
};
util.getToken = () => {
	return util.getItem("bearer");
};
// 获取保存的用户信息
util.getAuth = () => {
	try {
		let token = Base64.decode(util.getItem("bearer").split(".")[1]);
		let auth = JSON.parse(token);
		if (!auth.hasOwnProperty("id")) {
			util.removeItem("bearer");
			location.href = urlLogin;
		}
		return auth;
	} catch (e) {
		util.removeItem("bearer");
		location.href = urlLogin;
	}
};
util.Role = {
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
util.Role.isSup = rl => {
	return util.Role.getRole(rl, util.Role.RSup);
};
util.Role.isAtv = rl => {
	return util.Role.getRole(rl, util.Role.RAtv);
};

util.Role.allow = (role, arr) => {
	for (let i = 0; i < arr.length; i++) {
		if (util.Role.getRole(role, arr[i])) {
			return true;
		}
	}
	return false;
};
export default util;
