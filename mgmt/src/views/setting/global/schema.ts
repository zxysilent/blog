/**
 * 内置配置的字段描述
 *
 * 「全局配置」页按这份 schema 渲染表单：后端把配置存成 JSON 文本（config.value），
 * 并且对 auth / upload / sql 有区间校验，逐字段渲染既避免把「一段裸 JSON」
 * 丢给使用者，也让必填/范围/单位与后端保持一致（不一致时以后端报错为准）。
 *
 * 未列入此处的配置（用户自建、或 layout 这类由专门页面维护的）不在全局配置页出现。
 */

/** 表单项类型 */
export type FieldType = 'text' | 'textarea' | 'number' | 'switch' | 'password' | 'radio';

/** 单选项 */
export interface FieldOption {
  label: string;
  value: string;
  /** 选项补充说明，显示在单选下方 */
  tip?: string;
}

export interface Field {
  key: string;
  label: string;
  type: FieldType;
  placeholder?: string;
  /** 数字型：最小值/最大值/单位 */
  min?: number;
  max?: number;
  step?: number;
  suffix?: string;
  /** 单选型：可选项 */
  options?: FieldOption[];
  /** 补充说明，显示在控件下方 */
  tip?: string;
}

export interface ConfigSchema {
  /** 配置主键，与后端 config.id 一致 */
  id: string;
  /** 页签名称 */
  tab: string;
  fields: Field[];
}

export const configSchemas: ConfigSchema[] = [
  {
    id: 'basic',
    tab: '基础配置',
    fields: [
      { key: 'title', label: '网站标题', type: 'text', placeholder: '如：我的博客' },
      { key: 'site_url', label: '网站地址', type: 'text', placeholder: '如：https://blog.example.com' },
      { key: 'logo_url', label: 'Logo 地址', type: 'text', placeholder: '图片 URL，可留空' },
      { key: 'author', label: '网站作者', type: 'text' },
      { key: 'keywords', label: 'SEO 关键词', type: 'text', placeholder: '英文逗号分隔' },
      { key: 'description', label: 'SEO 描述', type: 'textarea' },
      { key: 'beian_miit', label: 'ICP 备案', type: 'text' },
      { key: 'beian_nism', label: '公安备案', type: 'text' },
      { key: 'copyright', label: '版权信息', type: 'textarea' },
    ],
  },
  {
    id: 'login',
    tab: '登录页面',
    fields: [
      { key: 'title', label: '页面标题', type: 'text', placeholder: '如：欢迎登录' },
      { key: 'desc', label: '页面描述', type: 'text', placeholder: '如：便捷、高效、可靠' },
    ],
  },
  {
    id: 'blog',
    tab: '博文配置',
    fields: [
      {
        key: 'page_size',
        label: '分页条数',
        type: 'number',
        min: 1,
        max: 100,
        tip: '前台列表每页展示的条数',
      },
      { key: 'analytic', label: '统计脚本', type: 'textarea', placeholder: '如：百度统计代码' },
      { key: 'comment', label: '评论脚本', type: 'textarea', placeholder: '第三方评论系统代码' },
      { key: 'github_url', label: 'GitHub 地址', type: 'text' },
      { key: 'weibo_url', label: '微博地址', type: 'text' },
    ],
  },
  {
    id: 'auth',
    tab: '认证配置',
    fields: [
      {
        key: 'token_exp',
        label: '令牌有效期',
        type: 'number',
        min: 1,
        max: 720,
        suffix: '小时',
        tip: '只影响新签发的令牌，已签发的不变',
      },
      {
        key: 'token_secret',
        label: '加密私钥',
        type: 'password',
        tip: '16~64 字符；修改后所有已登录用户的令牌立即失效',
      },
      {
        key: 'audit',
        label: '操作审计',
        type: 'switch',
        tip: '关闭后不再记录写操作与登录日志',
      },
      {
        key: 'docs',
        label: '接口文档',
        type: 'switch',
        tip: '开启后可通过 /swagger/index.html 查看 OpenAPI 文档，生产环境建议关闭',
      },
    ],
  },
  {
    id: 'upload',
    tab: '上传配置',
    fields: [
      {
        key: 'file_max_size',
        label: '文件上限',
        type: 'number',
        min: 1,
        max: 512,
        suffix: 'MB',
      },
      {
        key: 'image_max_size',
        label: '图片上限',
        type: 'number',
        min: 1,
        max: 100,
        suffix: 'MB',
        tip: '不能超过文件上限',
      },
      {
        key: 'image_max_width',
        label: '图片宽度上限',
        type: 'number',
        min: 100,
        max: 10000,
        suffix: 'px',
      },
      {
        key: 'image_max_height',
        label: '图片高度上限',
        type: 'number',
        min: 100,
        max: 10000,
        suffix: 'px',
      },
      {
        key: 'crop',
        label: '自动缩小',
        type: 'switch',
        tip: '开启时超出边长上限的图片会被等比缩小；关闭则直接拒绝',
      },
    ],
  },
  {
    id: 'sql',
    tab: 'SQL 日志',
    fields: [
      { key: 'show', label: '打印全部 SQL', type: 'switch', tip: '调试用，生产环境建议关闭' },
      {
        key: 'slow_ms',
        label: '慢 SQL 阈值',
        type: 'number',
        min: 0,
        max: 60000,
        suffix: 'ms',
        tip: '超过该耗时按 warn 记录；0 表示关闭慢日志',
      },
    ],
  },
];

/** 按主键取 schema */
export function findSchema(id: string): ConfigSchema | undefined {
  return configSchemas.find((it) => it.id === id);
}
