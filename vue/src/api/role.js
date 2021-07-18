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

// ------------------------------------------------------ 角色授权 ------------------------------------------------------
export const admGrantTree = data => {
    return fetch.request({
		url: "/adm/grant/tree",
		method: "get",
		params: data
	});
};

export const admRoleGrantAll = data => {
    return fetch.request({
		url: "/adm/role/grant/all",
		method: "get",
		params: data
	});
};

export const adadmRoleGrantEdit = data => {
    return fetch.request({
		url: "/adm/role/grant/edit",
		method: "post",
		data: data
	});
};
