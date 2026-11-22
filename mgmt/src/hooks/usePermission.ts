import { useUserStore } from '@/store/user';

/**
 * 权限判断
 *
 * 后端 /api/auth/permissions 返回的是权限标识字符串数组（如 ['user:list']），
 * 但模板类项目历史上也见过 [{ label, value }] 的形态，这里统一取标量值，
 * 两种数据形状都能判定，避免出现「明明有权限却判成没有」。
 */
function permValue(item: unknown): string {
  if (typeof item === 'string') return item;
  const value = (item as Recordable)?.value;
  return typeof value === 'string' ? value : '';
}

export function usePermission() {
  const userStore = useUserStore();

  /**
   * 检查权限
   * @param accesses
   */
  function _somePermissions(accesses: string[]) {
    return userStore.getPermissions.some((item) => accesses.includes(permValue(item)));
  }

  /**
   * 判断是否存在权限
   * 可用于 v-if 显示逻辑
   * */
  function hasPermission(accesses: string[]): boolean {
    if (!accesses || !accesses.length) return true;
    return _somePermissions(accesses);
  }

  /**
   * 是否包含指定的所有权限
   * @param accesses
   */
  function hasEveryPermission(accesses: string[]): boolean {
    if (!Array.isArray(accesses)) {
      throw new Error(`[hasEveryPermission]: ${accesses} should be a array !`);
    }
    if (!accesses.length) return true;
    const owned = userStore.getPermissions.map(permValue);
    return accesses.every((item) => owned.includes(item));
  }

  /**
   * 是否包含其中某个权限
   * @param accesses
   */
  function hasSomePermission(accesses: string[]): boolean {
    if (!Array.isArray(accesses)) {
      throw new Error(`[hasSomePermission]: ${accesses} should be a array !`);
    }
    return hasPermission(accesses);
  }

  return { hasPermission, hasEveryPermission, hasSomePermission };
}
