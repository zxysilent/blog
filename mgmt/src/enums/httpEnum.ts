/**
 * @description: 请求结果集
 * 与后端 pkg/reply 保持一致：HTTP 恒为 200，业务结果看 code
 */
export enum ResultEnum {
  /** 成功 */
  SUCCESS = 200,
  /** 失败-业务校验未通过 */
  FAIL = 300,
  /** 没有权限 */
  DENY = 340,
  /** 令牌无效或已过期，需要重新登录 */
  TOKEN = 360,
  ERROR = -1,
  TIMEOUT = 10042,
  TYPE = 'success',
}

/**
 * @description: 请求方法
 */
export enum RequestEnum {
  GET = 'GET',
  POST = 'POST',
  PATCH = 'PATCH',
  PUT = 'PUT',
  DELETE = 'DELETE',
}

/**
 * @description:  常用的contentTyp类型
 */
export enum ContentTypeEnum {
  // json
  JSON = 'application/json;charset=UTF-8',
  // json
  TEXT = 'text/plain;charset=UTF-8',
  // form-data 一般配合qs
  FORM_URLENCODED = 'application/x-www-form-urlencoded;charset=UTF-8',
  // form-data  上传
  FORM_DATA = 'multipart/form-data;charset=UTF-8',
}
