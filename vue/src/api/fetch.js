import Axios from "axios";
import Storage from "@/utils/storage";
import ViewUI from "view-design";
import Router from "@/router";
const fetch = Axios.create({
	baseURL: process.env.VUE_APP_SRV,
	timeout: 30000
});
//添加一个请求拦截器
fetch.interceptors.request.use(
	function(config) {
		ViewUI.LoadingBar.start();
		//在请求发出之前进行一些操作
		console.log("send");
		// 仅以/adm开始的接口才传递 token
		// 可自行修改为一直携带或者仅登录不携带-由程序内自主控制
		// if (config.url.indexOf("/api/auth/login") == -1) {
		if (config.url.indexOf("/adm") == 0) {
			config.headers.Authorization = Storage.getToken();
		}
		return config;
	},
	function(err) {
		//Do something with request error
		return Promise.reject(err);
	}
);
//添加一个响应拦截器
fetch.interceptors.response.use(
	function(resp) {
		ViewUI.LoadingBar.finish();
		//在这里对返回的数据进行处理
		console.log("recv");
		if (resp.data.code == 330) {
			ViewUI.Notice.error({
				duration: 3,
				title: "系统提醒",
				desc: "对不起你没有权限访问"
			});
			return new Promise(() => {});
		}
		if (resp.data.code == 340) {
			Storage.clear()
			Router.push({ name: "login" });
			return new Promise(() => {});
		}
		if (resp.data.code == 350) {
            ViewUI.Notice.error({
				duration: 3,
				title: "系统提醒",
				desc: resp.data.msg
			});
			Router.push({ name: "500" });
			return new Promise(() => {});
		}
		return resp.data;
	},
	function(err) {
		ViewUI.Notice.error({
			duration: 5,
			title: "系统提醒",
			desc: "网络连接异常,请重试"
		});
		console.log(err);
		return new Promise(() => {});
	}
);
export const apiUploadImage = (file, attr) => {
	let data = new FormData();
	data.append("file", file);
	// data.append("typ", "image");
	// data.append("cut", "yes");
	//   data.append("h","150");
	// data.append("w", "1000");
	attr.cut = true;
	return fetch.post("/api/upload/image", data, { params: attr });
};

export default fetch;
