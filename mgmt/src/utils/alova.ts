import { createAlova } from 'alova';
import VueHook from 'alova/vue';
import adapterFetch from 'alova/fetch';
import { isString } from 'lodash-es';
import { useUser } from '@/store/user';
import { storage } from '@/utils/Storage';
import { useGlobSetting } from '@/hooks/useSetting';
import { PageEnum } from '@/enums/pageEnum';
import { ResultEnum } from '@/enums/httpEnum';
import { isUrl } from '@/utils';

const { apiUrl, urlPrefix } = useGlobSetting();

const requestAdapter = adapterFetch();

export const Alova = createAlova({
  baseURL: apiUrl,
  statesHook: VueHook,
  // 关闭全局请求缓存
  // alova v3 默认给 GET 开 5 分钟内存缓存，管理端列表在增删改后必须立刻看到最新数据，
  // 这里显式置为 null，避免「新增成功但列表还是旧数据」。
  cacheFor: null,
  requestAdapter,
  beforeRequest(method) {
    const userStore = useUser();
    const token = userStore.getToken;
    // 添加 token 到请求头
    if (!method.meta?.ignoreToken && token) {
      method.config.headers['token'] = token;
    }
    // 处理 api 请求前缀
    const isUrlStr = isUrl(method.url as string);
    if (!isUrlStr && urlPrefix) {
      method.url = `${urlPrefix}${method.url}`;
    }
    if (!isUrlStr && apiUrl && isString(apiUrl)) {
      method.url = `${apiUrl}${method.url}`;
    }
  },
  responded: {
    onSuccess: async (response, method) => {
      const res = (response.json && (await response.json())) || response.body;

      // 是否返回原生响应头 比如：需要获取响应头时使用该属性
      if (method.meta?.isReturnNativeResponse) {
        return res;
      }

      // 后端统一信封：{ code, msg, data, count }
      // count 为分页总数，无分页时不返回
      const { code, msg, data, count } = res || {};

      // 不进行任何处理，直接返回完整信封
      // 用于需要拿到 code/msg 自己判断的场景（如登录页）
      if (method.meta?.isTransformResponse === false) {
        return res;
      }

      // @ts-expect-error 全局挂载的 naive-ui discrete api
      const Message = window.$message;
      // @ts-expect-error 全局挂载的 naive-ui discrete api
      const Modal = window.$dialog;

      const LoginPath = PageEnum.BASE_LOGIN;

      if (ResultEnum.SUCCESS === code) {
        // 分页接口把 count 一并带回，方便调用方组装表格数据
        if (typeof count === 'number') {
          return { data, count };
        }
        return data;
      }

      // 令牌失效/非法：清理本地登录态后重新登录
      if (ResultEnum.TOKEN === code) {
        Modal?.warning({
          title: '提示',
          content: msg || '登录身份已失效，请重新登录!',
          okText: '确定',
          closable: false,
          maskClosable: false,
          onOk: async () => {
            storage.clear();
            // hash 模式下登录页是 <当前目录>/#/login：
            // 保留 pathname（生产环境是 /mgmt/）只替换 hash，避免跳出 SPA
            const { pathname, search } = window.location;
            window.location.replace(`${pathname}${search}#${LoginPath}`);
          },
        });
        throw new Error(msg || '登录已失效');
      }

      // 无权限与业务失败：统一提示并中止后续流程
      Message?.error(msg || '请求失败');
      throw new Error(msg || '请求失败');
    },
  },
});

// 项目，多个不同 api 地址，可导出多个实例
// export const AlovaTwo = createAlova({
//   baseURL: 'http://localhost:9001',
// });
