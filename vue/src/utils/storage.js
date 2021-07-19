// 开发模式使用 localStorage
// 正式上线使用 sessionStorage
const storage = process.env.NODE_ENV === "development" ? localStorage : sessionStorage;
// token key
const tokenKey = "zblog-token";
const grantKey = "zblog-grant";
// 设置token
storage.__proto__.setToken = token => {
	storage.setItem(tokenKey, token);
};
// 读取token
storage.__proto__.getToken = () => {
	return storage.getItem(tokenKey);
};
// 设置token
storage.__proto__.setGrant = token => {
	storage.setItem(grantKey, token);
};
// 读取token
storage.__proto__.getGrant = () => {
	return storage.getItem(grantKey);
};
export default storage;
