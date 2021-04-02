import fetch from "./fetch";

export const admMenuAll = data => {
	return fetch.get("/adm/menu/all", { params: data });
};

export const admMenuTree = data => {
	return fetch.get("/adm/menu/tree", { params: data });
};

export const admMenuGet = data => {
	return fetch.get("/adm/menu/get", { params: data });
};

export const admMenuAdd = data => {
	return fetch.post("/adm/menu/add", data);
};

export const admMenuEdit = data => {
	return fetch.post("/adm/menu/edit", data);
};
export const admMenuEditShow = data => {
	return fetch.post("/adm/menu/edit/show", data);
};

export const admMenuDrop = data => {
	return fetch.post("/adm/menu/drop", data);
};
