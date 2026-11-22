import type { PaginationProps } from '../types/pagination';
import type { BasicTableProps } from '../types/table';
import { computed, unref, ref, ComputedRef, watch } from 'vue';

import { isBoolean } from '@/utils/is';
import { useDictStore } from '@/store/dict';

export function usePagination(refProps: ComputedRef<BasicTableProps>) {
  const configRef = ref<PaginationProps>({});
  const show = ref(true);
  // 分页条数可选值来自后端 page_size 字典（未就绪时用内置兜底）
  const dictStore = useDictStore();

  watch(
    () => unref(refProps).pagination,
    (pagination) => {
      if (!isBoolean(pagination) && pagination) {
        configRef.value = {
          ...unref(configRef),
          ...(pagination ?? {}),
        };
      }
    }
  );

  const getPaginationInfo = computed((): PaginationProps | boolean => {
    const { pagination } = unref(refProps);
    if (!unref(show) || (isBoolean(pagination) && !pagination)) {
      return false;
    }
    const info: PaginationProps = {
      page: 1, //当前页
      pageSize: dictStore.getPageSize, //分页大小
      pageSizes: dictStore.getPageSizes, // 每页条数
      showSizePicker: true,
      showQuickJumper: true,
      prefix: (pagingInfo) => `共 ${pagingInfo.itemCount} 条`, // 不需要可以通过 pagination 重置或者删除
      ...(isBoolean(pagination) ? {} : pagination),
      ...unref(configRef),
    };
    // 页面显式指定的条数不在字典里时补进可选项，避免下拉缺当前值
    const size = Number(info.pageSize);
    if (size > 0 && !info.pageSizes?.includes(size)) {
      info.pageSizes = [size, ...(info.pageSizes || [])].sort((a, b) => a - b);
    }
    return info;
  });

  function setPagination(info: Partial<PaginationProps>) {
    const paginationInfo = unref(getPaginationInfo);
    const next: PaginationProps = {
      ...(!isBoolean(paginationInfo) ? paginationInfo : {}),
      ...info,
    };
    // 可选条数以字典为准，不随局部更新被固化进 configRef，字典变了能实时生效
    delete next.pageSizes;
    configRef.value = next;
  }

  function getPagination() {
    return unref(getPaginationInfo);
  }

  function getShowPagination() {
    return unref(show);
  }

  async function setShowPagination(flag: boolean) {
    show.value = flag;
  }

  return { getPagination, getPaginationInfo, setShowPagination, getShowPagination, setPagination };
}
