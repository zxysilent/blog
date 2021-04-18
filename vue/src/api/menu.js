import fetch from "./fetch";

// 通过id获取单条菜单
export const admMenuGet = data => {
	return fetch.request({
		url: "/adm/menu/get",
		method: "get",
		params: data
	});
};
// 获取所有菜单
export const admMenuAll = data => {
	return fetch.request({
		url: "/adm/menu/all",
		method: "get",
		params: data
	});
};
// 获取所有菜单树
export const admMenuTree = data => {
	return fetch.request({
		url: "/adm/menu/tree",
		method: "get",
		params: data
	});
};
// 获取菜单分页
export const admMenuPage = data => {
	return fetch.request({
		url: "/adm/menu/page",
		method: "get",
		params: data
	});
};
// 添加菜单
export const admMenuAdd = data => {
	return fetch.request({
		url: "/adm/menu/add",
		method: "post",
		data: data
	});
};
// 修改菜单
export const admMenuEdit = data => {
	return fetch.request({
		url: "/adm/menu/edit",
		method: "post",
		data: data
	});
};
// 修改菜单显隐
export const admMenuEditShow = data => {
	return fetch.request({
		url: "/adm/menu/edit/show",
		method: "post",
		data: data
	});
};
// 通过id删除单条菜单
export const admMenuDrop = data => {
	return fetch.request({
		url: "/adm/menu/drop",
		method: "post",
		data: data
	});
};
