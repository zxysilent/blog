export const urlJwt = "/core/jwt"; //需要重新登陆
export const url404 = "/core/404"; //未发现
export const url401 = "/core/401"; //没有权限
export const url50x = "/core/50x"; //服务器异常
export const urlLogin = "/core/login";
export const urlServer = process.env.NODE_ENV === "development" ? "http://127.0.0.1:88" : "http://127.0.0.1:88";
export const urlUpload = urlServer + "/upload";
