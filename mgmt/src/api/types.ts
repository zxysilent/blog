/**
 * 后端数据模型与请求/响应类型
 * 与 go 后端 internal/model 保持一致，字段名全部为后端下发的 json 名
 */

/** 统一响应信封，HTTP 恒为 200，业务结果看 code */
export interface Reply<T = any> {
  code: number;
  msg: string;
  data?: T;
  /** 分页接口返回的数据总数 */
  count?: number;
}

/** 分页请求参数（后端为 page/size/q） */
export interface PageArgs {
  page?: number;
  size?: number;
  /** 模糊查询关键字 */
  q?: string;
  [key: string]: any;
}

/** 表格组件需要的分页结果 */
export interface PageResult<T> {
  list: T[];
  page: number;
  pageCount: number;
  itemCount: number;
}

/** 登录参数 */
export interface LoginArgs {
  num: string;
  passwd: string;
  /** 验证码明文 */
  vcode: string;
  /** 验证码签名，由 /api/captcha 下发 */
  vreal: string;
}

/** 图形验证码 */
export interface Captcha {
  /** base64 图片 */
  image: string;
  /** 签名，登录时回传 */
  real: string;
}

/** 角色 */
export interface Role {
  id: string;
  name: string;
  intro: string;
  /** 层级-数值越小权限越高，0 为根角色 */
  level: number;
  /** 内置角色，禁止删除 */
  inner: boolean;
  /** 关联用户数，列表接口填充 */
  user_count?: number;
  /** 已授权权限数，列表接口填充 */
  perm_count?: number;
  updated: number;
  created: number;
}

/** 权限标识，由路由注册自动同步，只读 */
export interface Permission {
  id: string;
  perm: string;
  name: string;
  group: string;
  group_key?: string;
  sort: number;
  created: number;
}

/** 权限树节点 */
export interface PermissionNode {
  id?: string;
  name: string;
  /** 分组键，如 user、role，用于派生分组授权权限 */
  group_key?: string;
  group: string;
  perm?: string;
  sort: number;
  children?: PermissionNode[];
}

/** 用户 */
export interface User {
  id: string;
  num: string;
  name: string;
  avatar: string;
  phone: string;
  email: string;
  role_id: string;
  role?: Role;
  /** 1 管理员 2 普通用户 */
  type: number;
  /** 1 正常 2 禁用 */
  status: number;
  updated: number;
  created: number;
}

/** 字典选项 */
export interface DictEntry {
  id: string;
  dict_id?: string;
  pid?: string;
  field: string;
  label: string;
  value: string;
  sort: number;
  status: boolean;
  type?: string;
  children?: DictEntry[];
  created?: number;
  updated?: number;
}

/** 字典 */
export interface Dict {
  id: string;
  name: string;
  type: string;
  status: boolean;
  inner: boolean;
  remark: string;
  entries: DictEntry[];
  created: number;
  updated: number;
}

/** 配置 */
export interface Config {
  id: string;
  group: string;
  name: string;
  value: string;
  min: number;
  max: number;
  type: string;
  describe: string;
  editable: boolean;
  sort: number;
  created: number;
  updated: number;
}

/** 配置分组统计 */
export interface ConfigGroup {
  group: string;
  count: number;
}

/** 栏目 */
export interface Cate {
  id: string;
  pid: string;
  kind: number;
  name: string;
  path: string;
  url: string;
  intro: string;
  icon: string;
  color: string;
  sort: number;
  level: number;
  use: boolean;
  show: boolean;
  inner: boolean;
  auth: boolean;
  itpl: string;
  ltpl: string;
  dtpl: string;
  count?: number;
  children?: Cate[];
  updated: number;
  created: number;
}

/** 标签 */
export interface Tag {
  id: string;
  name: string;
  intro: string;
  color: string;
  updated: number;
  created: number;
}

/** 内容 */
export interface Post {
  id: string;
  pid: string;
  cate_id: string;
  uid: string;
  /** 1 文章 2 页面 3 笔记 */
  kind: number;
  /** 1 草稿 2 已发布 */
  status: number;
  show: boolean;
  istop: boolean;
  hits: number;
  title: string;
  path: string;
  summary: string;
  author: string;
  origin: string;
  icon: string;
  markdown: string;
  richtext: string;
  allow: boolean;
  inner: boolean;
  tags: Tag[];
  cate?: Cate;
  updated: number;
  created: number;
}

/** 操作日志 */
export interface Log {
  id: string;
  user_id: string;
  user?: User;
  module: string;
  action: string;
  perm: string;
  method: string;
  path: string;
  ip: string;
  status: number;
  data: string;
  created: number;
}

/** 站点公开配置 */
export interface Site {
  basic: Record<string, any>;
  blog: Record<string, any>;
  /** 登录页配置：title / desc / style(mgmt|vue) */
  login?: Record<string, any>;
  /** 是否开放接口文档（认证配置里的开关） */
  docs?: boolean;
}
