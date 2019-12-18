import ajax from "./base";

// 判断当前账号是否已经存在
export const apiUserExist = num => {
	return ajax.get(`/api/user/exist/${num}`);
};
//添加用户
export const admUserAdd = data => {
	return ajax.post("/adm/user/add", data);
};
//分页
export const apiUserPage = (rl, data) => {
	return ajax.get(`/api/user/page/${rl}`, { params: data });
};
//更新用户状态
export const admUserChgatv = id => {
	return ajax.get(`/adm/user/chgatv/${id}`);
};
//重置用户密码
export const admUserResetPass = id => {
	return ajax.get(`/adm/user/reset/pass/${id}`);
};
//修改用户
export const admUserEdit = data => {
	return ajax.post("/adm/user/edit", data);
};
//删除用户
export const admUserDrop = id => {
	return ajax.get(`/adm/user/drop/${id}`);
};
//修改自己密码
export const admUserPass = data => {
	return ajax.post("/adm/user/pass", data);
};
//修改自己信息
export const admUserEditSelf = data => {
	return ajax.post("/adm/user/edit/self", data);
};
