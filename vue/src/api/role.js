import fetch from "./fetch";

export const admRoleAll = data => {
	return fetch.get("/adm/role/all", { params: data });
};

export const admRoleGet = data => {
	return fetch.get("/adm/role/get", { params: data });
};

export const admRoleAdd = data => {
	return fetch.post("/adm/role/add", data);
};

export const admRoleEdit = data => {
	return fetch.post("/adm/role/edit", data);
};

export const admRoleDrop = data => {
	return fetch.post("/adm/role/drop", data);
};

// ------------------------------------------------------ 角色菜单 ------------------------------------------------------

export const admRoleMenuAll = data => {
	return fetch.get("/adm/role/menu/all", { params: data });
};

export const adadmRoleMenuEdit = data => {
	return fetch.post("/adm/role/menu/edit", data);
};
// ------------------------------------------------------ 角色接口 ------------------------------------------------------

export const admRoleApiAll = data => {
	return fetch.get("/adm/role/get", { params: data });
};

export const adadmRoleApiEdit = data => {
	return fetch.post("/adm/role/api/edit", data);
};