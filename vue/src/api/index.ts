import fetch, { Reply } from '@/utils/fetch';

/**
 * 验证码
 */
export const apiAuthVcode = (params): Promise<Reply> => {
    return fetch.request({
        url: `/api/auth/vcode`,
        method: "get",
		params: params
    });
};

/**
 * 通过id获取单条字典
 */
export const apiDictGet = (key, params = null): Promise<Reply> => {
    return fetch.request({
        url: `/api/dict/get/${key}`,
        method: "get",
		params: params
    });
};

/**
 * 通过id获取单条字典
 */
export const apiDictVal = (key, params = null): Promise<Reply> => {
    return fetch.request({
        url: `/api/dict/${key}`,
        method: "get",
		params: params
    });
};

/**
 * 登陆
 */
export const apiAuthLogin = (data): Promise<Reply> => {
    return fetch.request({
        url: `/api/auth/login`,
        method: "post",
        data: data
    });
};

/**
 * 获取当前用户的授权
 */
export const apiAuthGrant = (params): Promise<Reply> => {
    return fetch.request({
        url: `/api/auth/grant`,
        method: "get",
		params: params
    });
};

/**
 * 获取服务器go信息
 */
export const apiStatusGo = (params): Promise<Reply> => {
    return fetch.request({
        url: `/api/status/app`,
        method: "get",
		params: params
    });
};

/**
 * 获取服务器go信息
 */
export const apiStatusGo = (params): Promise<Reply> => {
    return fetch.request({
        url: `/api/status/go`,
        method: "get",
		params: params
    });
};

/**
 * 上传文件
 */
export const apiUploadFile = (data): Promise<Reply> => {
    return fetch.request({
        url: `/api/upload/file`,
        method: "post",
        data: data
    });
};

/**
 * 上传图片并裁剪
 */
export const apiUploadImage = (data): Promise<Reply> => {
    return fetch.request({
        url: `/api/upload/image`,
        method: "post",
        data: data
    });
};

/**
 * 通过key删除单条字典
 */
export const apiDictDrop = (data): Promise<Reply> => {
    return fetch.request({
        url: `/api/dict/drop`,
        method: "delete",
        data: data
    });
};

/**
 * 通过id删除单条授权信息
 */
export const apiGrantDrop = (id, data = null): Promise<Reply> => {
    return fetch.request({
        url: `/api/grant/${id}`,
        method: "delete",
        data: data
    });
};

/**
 * 获取admin分页数据
 */
export const apiAdminAll = (params): Promise<Reply> => {
    return fetch.request({
        url: `/api/admin/all`,
        method: "get",
		params: params
    });
};

/**
 * 获取admin分页数据
 */
export const apiAdminPage = (params): Promise<Reply> => {
    return fetch.request({
        url: `/api/admin/page`,
        method: "get",
		params: params
    });
};

/**
 * 通过id获取admin信息
 */
export const apiAdminGet = (id, params = null): Promise<Reply> => {
    return fetch.request({
        url: `/api/admin/${id}`,
        method: "get",
		params: params
    });
};

/**
 * 获取某个用户信息
 */
export const apiAdminExist = (num, params = null): Promise<Reply> => {
    return fetch.request({
        url: `/api/admin/${num}/exist`,
        method: "get",
		params: params
    });
};

/**
 * 当前登录信息
 */
export const apiAuthGet = (params): Promise<Reply> => {
    return fetch.request({
        url: `/api/auth/get`,
        method: "get",
		params: params
    });
};

/**
 * 获取字典分页
 */
export const apiDictPage = (params): Promise<Reply> => {
    return fetch.request({
        url: `/api/dict/page`,
        method: "get",
		params: params
    });
};

/**
 * 获取授权分页信息
 */
export const apiGrantPage = (params): Promise<Reply> => {
    return fetch.request({
        url: `/api/grant/page`,
        method: "get",
		params: params
    });
};

/**
 * 通过id获取单条授权信息
 */
export const apiGrantGet = (id, params = null): Promise<Reply> => {
    return fetch.request({
        url: `/api/grant/${id}`,
        method: "get",
		params: params
    });
};

/**
 * 获取所有角色信息
 */
export const apiRoleAll = (params): Promise<Reply> => {
    return fetch.request({
        url: `/api/role/all`,
        method: "get",
		params: params
    });
};

/**
 * 通过角色id获取所有角色授权信息
 */
export const apiRoleGrantAll = (params): Promise<Reply> => {
    return fetch.request({
        url: `/api/role/grant/all`,
        method: "get",
		params: params
    });
};

/**
 * 获取角色分页信息
 */
export const apiRolePage = (params): Promise<Reply> => {
    return fetch.request({
        url: `/api/role/page`,
        method: "get",
		params: params
    });
};

/**
 * 通过id获取单条角色信息
 */
export const apiRoleGet = (id, params = null): Promise<Reply> => {
    return fetch.request({
        url: `/api/role/${id}`,
        method: "get",
		params: params
    });
};

/**
 * 添加admin信息
 */
export const apiAdminAdd = (data): Promise<Reply> => {
    return fetch.request({
        url: `/api/admin/add`,
        method: "post",
        data: data
    });
};

/**
 * 删除admin信息
 */
export const apiAdminDrop = (data): Promise<Reply> => {
    return fetch.request({
        url: `/api/admin/drop`,
        method: "post",
        data: data
    });
};

/**
 * 修改admin信息
 */
export const apiAdminEdit = (data): Promise<Reply> => {
    return fetch.request({
        url: `/api/admin/edit`,
        method: "post",
        data: data
    });
};

/**
 * 修改用户锁定状态
 */
export const apiAdminEditLock = (data): Promise<Reply> => {
    return fetch.request({
        url: `/api/admin/edit/lock`,
        method: "post",
        data: data
    });
};

/**
 * 修改自己的信息
 */
export const apiAuthEdit = (data): Promise<Reply> => {
    return fetch.request({
        url: `/api/auth/edit`,
        method: "post",
        data: data
    });
};

/**
 * 修改自己的密码
 */
export const apiAuthEditPasswd = (data): Promise<Reply> => {
    return fetch.request({
        url: `/api/auth/edit/passwd`,
        method: "post",
        data: data
    });
};

/**
 * 添加字典
 */
export const apiDictAdd = (data): Promise<Reply> => {
    return fetch.request({
        url: `/api/dict/add`,
        method: "post",
        data: data
    });
};

/**
 * 添加授权信息
 */
export const apiGrantAdd = (data): Promise<Reply> => {
    return fetch.request({
        url: `/api/grant/add`,
        method: "post",
        data: data
    });
};

/**
 * 修改授权信息
 */
export const apiGrantEdit = (data): Promise<Reply> => {
    return fetch.request({
        url: `/api/grant/edit`,
        method: "post",
        data: data
    });
};

/**
 * 添加角色信息
 */
export const apiRoleAdd = (data): Promise<Reply> => {
    return fetch.request({
        url: `/api/role/add`,
        method: "post",
        data: data
    });
};

/**
 * 通过id删除单条角色信息
 */
export const apiRoleDrop = (data): Promise<Reply> => {
    return fetch.request({
        url: `/api/role/drop`,
        method: "post",
        data: data
    });
};

/**
 * 修改角色信息
 */
export const apiRoleEdit = (data): Promise<Reply> => {
    return fetch.request({
        url: `/api/role/edit`,
        method: "post",
        data: data
    });
};

/**
 * 修改角色授权信息
 */
export const apiRoleGrantEdit = (data): Promise<Reply> => {
    return fetch.request({
        url: `/api/role/grant/edit`,
        method: "post",
        data: data
    });
};

/**
 * 修改字典
 */
export const apiDictEdit = (data): Promise<Reply> => {
    return fetch.request({
        url: `/api/dict/edit`,
        method: "put",
        data: data
    });
};
