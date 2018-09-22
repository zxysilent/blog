import axios from "axios";
import env from "../../build/env";
import semver from "semver";
import packjson from "../../package.json";

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

export default util;
