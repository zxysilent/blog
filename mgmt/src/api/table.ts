import { Alova } from '@/utils/alova';

//获取table
export function getTableList(params) {
  return Alova.Get('/table/list', { params });
}
