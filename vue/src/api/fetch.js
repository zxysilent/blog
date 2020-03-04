import axios from "axios";
import Util from "@/utils.js";
import ViewUI from "view-design";

const fetch = axios.create({
	baseURL: process.env.VUE_APP_SRV,
	timeout: 30000
});
//添加一个请求拦截器
fetch.interceptors.request.use(
	function(config) {
		ViewUI.LoadingBar.start();
		//在请求发出之前进行一些操作
		console.log("send");
		if (config.url.indexOf("/api/login") == -1) {
			config.headers.Authorization = "Bearer " + Util.getItem("bearer"); //Bearer
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
	function(res) {
		ViewUI.LoadingBar.finish();
		//在这里对返回的数据进行处理
		console.log("recv");
		if (res.data.code == 330) {
			// ViewUI.Notice.error({
			// 	duration: 3,
			// 	title: "系统提醒",
			// 	desc: "对不起你没有权限访问"
			// });
			// return new Promise(() => {});
			location.href = "/#/401";//没有权限
		}
		if (res.data.code == 340) {
			// ViewUI.Notice.error({
			// 	duration: 2,
			// 	title: "系统提醒",
			// 	desc: "登陆失效,请重新登陆",
			// 	onClose: () => {
			// 		localStorage.removeItem("bearer");
			// 		//没有登陆信息默认跳转到登陆页面
			// 		// location.reload();
			// 	}
			// });
			// return new Promise(() => {});
			Util.clearData();
			location.href = "/#/jwt";//需要重新登陆
		}
		if (res.data.code == 350) {
			// ViewUI.Notice.error({
			// 	duration: 3,
			// 	title: "系统提醒",
			// 	desc: "服务端发生错误,请重试"
			// });
			// return new Promise(() => {});
			location.href = "/#/50x";//服务器异常
		}
		return res.data;
	},
	function(err) {
		ViewUI.Notice.error({
			duration: 5,
			title: "系统提醒",
			desc: "网络连接异常,请重试"
		});
		return new Promise(() => {});
		//return Promise.reject(err);
	}
);
export default fetch;
