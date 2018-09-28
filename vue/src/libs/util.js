import axios from "axios";
import env from "../../build/env";
import { Base64 } from "js-base64";

let util = {};
util.title = function(title) {
  title = title || "iView admin";
  window.document.title = title;
};

const ajaxUrl =
  env === "development"
    ? "http://127.0.0.1:8888"
    : env === "production"
      ? "https://www.url.com"
      : "https://debug.url.com";

util.ajax = axios.create({
  baseURL: ajaxUrl,
  timeout: 30000
});

util.inOf = function(arr, targetArr) {
  let res = true;
  arr.forEach(item => {
    if (targetArr.indexOf(item) < 0) {
      res = false;
    }
  });
  return res;
};

util.oneOf = function(ele, targetArr) {
  if (targetArr.indexOf(ele) >= 0) {
    return true;
  } else {
    return false;
  }
};

util.showThisRoute = function(itAccess, currentAccess) {
  if (typeof itAccess === "object" && Array.isArray(itAccess)) {
    return util.oneOf(currentAccess, itAccess);
  } else {
    return itAccess === currentAccess;
  }
};

util.getRouterObjByName = function(routers, name) {
  if (!name || !routers || !routers.length) {
    return null;
  }
  // debugger;
  let routerObj = null;
  for (let item of routers) {
    if (item.name === name) {
      return item;
    }
    routerObj = util.getRouterObjByName(item.children, name);
    if (routerObj) {
      return routerObj;
    }
  }
  return null;
};

util.handleTitle = function(vm, item) {
  return item.meta.title;
};

util.toDefaultPage = function(routers, name, route, next) {
  let len = routers.length;
  let i = 0;
  let notHandle = true;
  while (i < len) {
    if (
      routers[i].name === name &&
      routers[i].children &&
      routers[i].redirect === undefined
    ) {
      route.replace({
        name: routers[i].children[0].name
      });
      notHandle = false;
      next();
      break;
    }
    i++;
  }
  if (notHandle) {
    next();
  }
};

util.fullscreenEvent = function(vm) {
  vm.$store.commit("initCachepage");
  // 权限菜单过滤相关
  vm.$store.commit("updateMenulist");
  // 全屏相关
};

util.checkUpdate = function(vm) {};
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
   //   localStorage.removeItem("bearer");
   //   location.href = "/#/login";
    }
    return auth;
  } catch (e) {
    localStorage.removeItem("bearer");
   // location.href = "/#/login";
  }
};
//  RSup uint32 = 30 //super 		超级管理员
// 	RAdm uint32 = 29 //admin 		管理员
// 	RHmt uint32 = 28 //headmaster 	教师-班主任 可以发新闻
// 	RTea uint32 = 27 //teacher 		教师-课程老师
// 	RStu uint32 = 26 //student 		学生
// 	RAtv uint32 = 10 //active		启用/禁用
// 	RBas uint32 = 8  //base 		基本权限
util.Role = {
  RSup: 30, //超级管理员
  RAdm: 29, //管理员
  RHmt: 28, //教师-班主任 可以发新闻
  RTea: 27, //教师-课程老师
  RStu: 26, //学生
  RAtv: 10, //启用/禁用
  RBas: 8, //基本权限
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
util.Role.isAdm = rl => {
  return util.Role.getRole(rl, util.Role.RAdm);
};
util.Role.isHmt = rl => {
  return util.Role.getRole(rl, util.Role.RHmt);
};
util.Role.isTea = rl => {
  return util.Role.getRole(rl, util.Role.RTea);
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
