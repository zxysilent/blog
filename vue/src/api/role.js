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
