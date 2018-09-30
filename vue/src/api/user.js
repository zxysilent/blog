import ajax from "./base";

// 判断当前账号是否已经存在
export const userExist = num => {
  return ajax.get(`/user/exist/${num}`);
};

//添加用户
export const userAdd = data => {
  return ajax.post("/user/add", data);
};
//分页
export const userPage = (rl, data) => {
  return ajax.get(`/user/page/${rl}`, { params: data });
};
//更新用户状态
export const userChgatv = id => {
  return ajax.get(`/user/chgatv/${id}`);
};
//重置用户密码
export const userResetPass = id => {
  return ajax.get(`/user/reset/pass/${id}`);
};
//修改用户
export const userEdit = data => {
  return ajax.post("/user/edit", data);
};
//删除用户
export const userDel = id => {
  return ajax.get(`/user/del/${id}`);
};
//修改自己密码
export const userPass = data => {
  return ajax.post("/user/pass", data);
};
//修改自己信息
export const userEditSelf = data => {
  return ajax.post("/user/edit/self", data);
};
