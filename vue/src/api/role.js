import fetch from "./fetch";

export const admRoleAll = data => {
    return fetch.request({
		url: "/adm/role/all",
		method: "get",
		params: data
	});
};

export const admRoleGet = data => {
    return fetch.request({
		url: "/adm/role/get",
		method: "get",
		params: data
	});
};

export const admRoleAdd = data => {
    return fetch.request({
		url: "/adm/role/add",
		method: "post",
		data: data
	});
};

export const admRoleEdit = data => {
    return fetch.request({
		url: "/adm/role/edit",
		method: "post",
		data: data
	});
};

export const admRoleDrop = data => {
    return fetch.request({
		url: "/adm/role/drop",
		method: "post",
		data: data
	});
};

// ------------------------------------------------------ 角色菜单 ------------------------------------------------------

export const admRoleMenuAll = data => {
    return fetch.request({
		url: "/adm/role/menu/all",
		method: "get",
		params: data
	});
};

export const adadmRoleMenuEdit = data => {
    return fetch.request({
		url: "/adm/role/menu/edit",
		method: "post",
		data: data
	});
};
// ------------------------------------------------------ 角色接口 ------------------------------------------------------

export const admRoleApiAll = data => {
    return fetch.request({
		url: "/adm/role/get",
		method: "get",
		params: data
	});
};

export const adadmRoleApiEdit = data => {
    return fetch.request({
		url: "/adm/role/api/edit",
		method: "post",
		data: data
	});
};