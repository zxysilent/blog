import { Base64 } from "js-base64";
let util = {};
util.title = function(title) {
  title = title || "blog";
  window.document.title = title;
};
//dev prod
util.cfgSvrURL =
  process.env.NODE_ENV === "development"
    ? "http://127.0.0.1:88"
    : "http://127.0.0.1:88";

//文件上传地址
//dev prod
util.cfgUpload =
  process.env.NODE_ENV === "development"
    ? "http://127.0.0.1:88/upload"
    : "http://127.0.0.1:88/upload";

// 登陆失效默认页面
util.defURL = "/core/";
//保存数据
util.setData = (k, v) => {
  localStorage.setItem(k, v);
};
util.getData = k => {
  return localStorage.getItem(k);
};
util.clearData = k => {
  localStorage.removeItem(k);
};
util.setToken = token => {
  util.setData("bearer", token);
};
util.getToken = () => {
  return util.getData("bearer");
};
// 获取保存的用户信息
util.getAuth = () => {
  try {
    let token = Base64.decode(localStorage.getItem("bearer").split(".")[1]);
    let auth = JSON.parse(token);
    if (!auth.hasOwnProperty("id")) {
      localStorage.removeItem("bearer");
      location.href = util.defURL;
    }
    return auth;
  } catch (e) {
    localStorage.removeItem("bearer");
    location.href = util.defURL;
  }
};
util.Role = {
  RSup: 30, //超级管理员
  RAtv: 20, //启用/禁用
  RBas: 10, //基本权限
  //判断指定位置权限
  getRole: (rl, r) => {
    if ((rl & (1 << r)) >> r == 1) {
      return true;
    }
    return false;
  }
};
util.Role.isSup = rl => {
  return util.Role.getRole(rl, util.Role.RSup);
};
util.Role.isAtv = rl => {
  return util.Role.getRole(rl, util.Role.RAtv);
};

util.Role.allow = (role, arr) => {
  for (let i = 0; i < arr.length; i++) {
    if (util.Role.getRole(role, arr[i])) {
      return true;
    }
  }
  return false;
};
export default util;
