import fetch from "./fetch";

export const apiMenuAll = () => {
	return fetch.get("/api/menu/all");
};
export const apiMenuTree = data => {
	return fetch.get("/api/menu/tree", { params: data });
};

export const apiMenuGet = data => {
	return fetch.get("/api/menu/get", { params: data });
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
