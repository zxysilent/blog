export const urlJwt = "/jwt"; //需要重新登陆
export const url404 = "/404"; //未发现
export const url401 = "/401"; //没有权限
export const url50x = "/50x"; //服务器异常
export const urlLogin = "/login";
export const urlServer = process.env.NODE_ENV === "development" ? "http://127.0.0.1:88" : "http://127.0.0.1:88";
export const urlUpload = urlServer + "/upload";
