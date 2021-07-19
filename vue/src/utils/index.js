const utils = {};
utils.title = function(title) {
	title = title || "--";
	window.document.title = title;
};

// 不需要也可登录页面集合
utils.noAuth = r => {
	return ["login", "500", "404"].indexOf(r) > -1;
};

// yyyy-MM-dd hh:mm:ss
utils.fmtDate = (date, fmt) => {
	var o = {
		"M+": date.getMonth() + 1, //月份
		"d+": date.getDate(), //日
		"h+": date.getHours(), //小时
		"m+": date.getMinutes(), //分
		"s+": date.getSeconds() //秒
	};
	if (/(y+)/.test(fmt)) fmt = fmt.replace(RegExp.$1, (date.getFullYear() + "").substr(4 - RegExp.$1.length));
	for (var k in o)
		if (new RegExp("(" + k + ")").test(fmt))
			fmt = fmt.replace(RegExp.$1, RegExp.$1.length == 1 ? o[k] : ("00" + o[k]).substr(("" + o[k]).length));
	return fmt;
};
export default utils;
