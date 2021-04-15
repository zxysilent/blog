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
export default utils;
