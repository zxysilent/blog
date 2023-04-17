import axios, { AxiosRequestConfig, AxiosResponse, AxiosRequestHeaders } from "axios";
import storage from "./storage";
import Router from "@/router";
// const fetch = axios.create({ baseURL: import.meta.env.VITE_APP_SRV, timeout: 30000 });
const fetch = axios.create({ baseURL: "", timeout: 30000 });
export interface Reply {
    code: number;
    msg: string;
    data?: string | ReplyPage | any;
}
export interface ReplyPage {
    items: any
    count: number
}
//添加一个请求拦截器
//@ts-ignore
fetch.interceptors.request.use((config: AxiosRequestConfig) => {
    window.$loadingBar && window.$loadingBar.start();
    //在请求发出之前进行一些操作
    console.log("send");
    //@ts-ignore
    if (config.url.indexOf("/api") == 0) {
        (config.headers as AxiosRequestHeaders).Authorization = storage.getToken();
    }
    return config;
},
    (err) => {
        return Promise.reject(err);
    }
);
//添加一个响应拦截器
//@ts-ignore
fetch.interceptors.response.use((resp: AxiosResponse<Reply>): Promise<Reply> => {
    window.$loadingBar && window.$loadingBar.finish();
    //在这里对返回的数据进行处理
    console.log("recv");
    if (resp.data.code == 330) {
        Router.push({ name: "403" });
        return new Promise(() => { });
        // return Promise.resolve(resp.data);
    }
    if (resp.data.code == 340) {
        storage.clear();
        Router.push({ name: "login" });
        return new Promise(() => { });
    }
    if (resp.data.code == 350) {
        Router.push({ name: "500" });
        return new Promise(() => { });
    }
    return Promise.resolve<Reply>(resp.data);
},
    (err: any): any => {
        Router.push({ name: "500" });
        console.log(err);
        return new Promise(() => { });
    }
);
export default fetch;
