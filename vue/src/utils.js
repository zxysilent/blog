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

utils.clearData = () => {
	storage.clear();
};
utils.setToken = token => {
	utils.setItem("bearer", token);
};
utils.getToken = () => {
	return utils.getItem("bearer");
};

export default utils;
