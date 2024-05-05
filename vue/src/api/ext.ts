import fetch, { Reply } from '@/utils/fetch';
/**
 * 获取Basic配置值
 */
export const apiDictBasic = (): Promise<Reply> => {
    return fetch.request({
        url: "/api/dict/basic",
        method: "get",
    });
};
