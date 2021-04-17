// 开发模式使用 localStorage
// 正式上线使用 sessionStorage
const storage = process.env.NODE_ENV === "development" ? localStorage : sessionStorage;
// token key
const tokenKey = "bearer";
// 保存数据
export const setItem = (k, v) => {
	storage.setItem(k, v);
};
// 读取数据
export const getItem = k => {
	return storage.getItem(k);
};
// 清空数据
export const clearData = () => {
	storage.clear();
};
// 设置token
export const setToken = token => {
	storage.setItem(tokenKey, token);
};
// 读取token
export const getToken = () => {
	return storage.getItem(tokenKey);
};
