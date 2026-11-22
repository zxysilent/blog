/**
 * 字典类型与选项的基础约定
 * 与 go 后端 internal/model/dict.go 保持一致
 */

/** 值类型 */
export const dictTypes = [
  { label: '字符串 string', value: 'string' },
  { label: '布尔 bool', value: 'bool' },
  { label: '整数 int', value: 'int' },
  { label: '时间 time', value: 'time' },
];

export const dictTypeMap: Record<string, string> = {
  string: '字符串',
  bool: '布尔',
  int: '整数',
  time: '时间',
};

/** 时间类型允许的格式，后端按顺序尝试解析 */
export const dictTimeHint = '支持 2006-01-02 15:04:05 / 2006-01-02 / 15:04:05 / RFC3339';

/** 按类型校验字典值，返回错误提示，合法时返回空串 */
export function checkDictValue(type: string, value: string): string {
  const v = (value || '').trim();
  if (type === 'bool' && v !== 'true' && v !== 'false') {
    return '布尔值只能是 true 或 false';
  }
  if (type === 'int' && !/^-?\d+$/.test(v)) {
    return '必须是整数';
  }
  if (type === 'time' && v === '') {
    return '时间不能为空';
  }
  return '';
}

/** 选项行：与后端 DictEntry 对齐，前端额外维护一个稳定的 key */
export interface EntryRow {
  /** 前端行标识，不提交给后端 */
  rowKey: string;
  id?: string;
  /** 字段名，同层唯一 */
  field: string;
  label: string;
  value: string;
  sort: number;
  status: boolean;
}

let seq = 0;
/** 生成前端行 key */
export function nextRowKey(): string {
  seq += 1;
  return `r${Date.now().toString(36)}${seq}`;
}

/** 新建空白行 */
export function emptyRow(index = 1): EntryRow {
  return {
    rowKey: nextRowKey(),
    field: '',
    label: '',
    value: '',
    sort: index,
    status: true,
  };
}
