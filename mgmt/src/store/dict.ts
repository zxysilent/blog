import { defineStore } from 'pinia';
import { store } from '@/store';
import { getDictEntries } from '@/api/dict';
import componentSetting from '@/settings/componentSetting';

/**
 * 分页条数字典
 *
 * 后端内置了 page_size 字典（配置字典里可维护选项），列表页的「x 条/页」
 * 下拉直接读它，避免前后端两处各写一份可选值。
 * 字典读取失败时退回 componentSetting 里的内置兜底，保证列表还能用。
 */
export const PAGE_SIZE_DICT = 'page_size';

const { table } = componentSetting;
const FALLBACK_SIZES = [...table.pageSizes];
const FALLBACK_SIZE = table.defaultPageSize;

export interface IDictState {
  /** 可选分页条数，升序 */
  pageSizes: number[];
  /** 默认分页条数 */
  pageSize: number;
  /** 字典是否已成功加载 */
  loaded: boolean;
}

/** 并发调用只发一次请求 */
let pending: Promise<void> | null = null;

export const useDictStore = defineStore({
  id: 'app-dict',
  state: (): IDictState => ({
    pageSizes: FALLBACK_SIZES,
    pageSize: FALLBACK_SIZE,
    loaded: false,
  }),
  getters: {
    getPageSizes(state): number[] {
      return state.pageSizes;
    },
    getPageSize(state): number {
      return state.pageSize;
    },
  },
  actions: {
    /** 加载分页条数字典，已加载过则直接返回 */
    async loadPageSize(force = false): Promise<void> {
      if (!force) {
        if (this.loaded) return;
        if (pending) return pending;
      }
      pending = (async () => {
        try {
          const rows = await getDictEntries(PAGE_SIZE_DICT);
          const sizes = [
            ...new Set(
              (rows || [])
                .map((it) => Number(it.value))
                .filter((n) => Number.isInteger(n) && n > 0)
            ),
          ].sort((a, b) => a - b);
          if (sizes.length) {
            this.pageSizes = sizes;
            // 当前默认值不在字典里时，用字典的第一项
            if (!sizes.includes(this.pageSize)) this.pageSize = sizes[0];
            this.loaded = true;
          }
        } catch {
          // 接口异常保留内置兜底，不阻塞列表
        } finally {
          pending = null;
        }
      })();
      return pending;
    },
  },
});

// Need to be used outside the setup
export function useDict() {
  return useDictStore(store);
}
