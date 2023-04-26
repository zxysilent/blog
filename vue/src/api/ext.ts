// =-= 一些不好生成的代码
import fetch, { Reply } from '@/utils/fetch';

/**
 * 获取配置
 */
export const apiDictGet = (key: string): Promise<Reply> => {
    return fetch.request({
        url: "/api/dict/get/" + key,
        method: "get",
    });
};
/**
 * 获取配置值
 */
export const apiDictVal = (key: string): Promise<Reply> => {
    return fetch.request({
        url: "/api/dict/" + key,
        method: "get",
    });
};



/**
 * 获取Basic配置值
 */
export const apiDictBasic = (): Promise<Reply> => {
    return fetch.request({
        url: "/api/dict/basic",
        method: "get",
    });
};
