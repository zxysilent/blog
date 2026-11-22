<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="后台首页" />
    </div>

    <!-- 快速入口：按权限出，缺权限的入口不渲染，点了也不会被拦 -->
    <n-card v-if="entries.length" :bordered="false" title="快速入口" size="small" class="mt-4 proCard">
      <n-space :size="10" :wrap="true">
        <n-button
          v-for="it in entries"
          :key="it.label"
          size="small"
          strong
          secondary
          type="primary"
          @click="go(it)"
        >
          <template #icon>
            <n-icon><component :is="it.icon" /></n-icon>
          </template>
          {{ it.label }}
        </n-button>
      </n-space>
    </n-card>

    <!-- 系统信息：运行时信息，需要 status:read -->
    <n-card v-if="hasStatus" :bordered="false" title="系统信息" size="small" class="mt-4 proCard">
      <n-descriptions bordered label-placement="top" :column="3">
        <n-descriptions-item label="系统架构">
          <n-tag type="info" size="small">{{ status.arch || '-' }}</n-tag>
        </n-descriptions-item>
        <n-descriptions-item label="操作系统">
          <n-tag type="info" size="small">{{ status.os || '-' }}</n-tag>
        </n-descriptions-item>
        <n-descriptions-item label="运行环境">
          <n-tag type="info" size="small">{{ status.version || '-' }}</n-tag>
        </n-descriptions-item>
        <n-descriptions-item label="逻辑处理器">
          <n-tag type="info" size="small">{{ status.num_cpu ?? '-' }}</n-tag>
        </n-descriptions-item>
        <n-descriptions-item label="当前协程数">
          <n-tag type="info" size="small">{{ status.num_goroutine ?? '-' }}</n-tag>
        </n-descriptions-item>
        <n-descriptions-item label="作者主页">
          <n-a :href="AUTHOR_URL" target="_blank" class="text-sm">点击前往</n-a>
        </n-descriptions-item>
      </n-descriptions>
    </n-card>
  </div>
</template>

<script lang="ts" setup>
  import { computed, onMounted, ref } from 'vue';
  import { useRouter } from 'vue-router';
  import {
    AppstoreOutlined,
    FileTextOutlined,
    HistoryOutlined,
    HomeOutlined,
    PlusOutlined,
    ProfileOutlined,
    SettingOutlined,
    TagsOutlined,
    UserOutlined,
  } from '@vicons/antd';
  import { usePermission } from '@/hooks/usePermission';
  import { getSite, getStatus } from '@/api/site';

  /** 作者主页：与旧版后台保持一致 */
  const AUTHOR_URL = 'https://github.com/zxysilent';

  /**
   * 后台首页
   *
   * 只放「去哪儿」和「运行环境」两件事：
   * - 快速入口：每个入口带自己的权限标识，没有权限就不渲染（系统日志 = log:list）；
   * - 系统信息：运行时数据，需要 status:read，取不到就整块不显示。
   */
  const router = useRouter();
  const { hasPermission } = usePermission();

  const hasStatus = computed(() => hasPermission(['status:read']));
  const status = ref<Record<string, any>>({});
  /** 接口文档开关（认证配置），关闭后不展示入口 */
  const docsOn = ref(true);

  /** 站内跳转与站外打开区分开：external 的用新窗口 */
  type Entry = {
    label: string;
    icon: any;
    perm?: string;
    path?: string;
    external?: boolean;
    /** 额外条件，为 false 时不渲染 */
    show?: () => boolean;
  };

  const allEntries: Entry[] = [
    // 站外：接口文档与博客前台，不需要权限
    {
      label: '接口文档',
      icon: FileTextOutlined,
      path: '/swagger/index.html',
      external: true,
      show: () => docsOn.value,
    },
    { label: '前台首页', icon: HomeOutlined, path: 'site', external: true },
    { label: '写博文', icon: PlusOutlined, perm: 'post:create', path: '/content/post/add' },
    { label: '栏目管理', icon: AppstoreOutlined, perm: 'cate:list', path: '/content/cate' },
    { label: '标签管理', icon: TagsOutlined, perm: 'tag:list', path: '/content/tag' },
    { label: '用户管理', icon: UserOutlined, perm: 'user:list', path: '/system/user' },
    { label: '全局配置', icon: SettingOutlined, perm: 'config:list', path: '/setting/global' },
    { label: '个人设置', icon: ProfileOutlined, path: '/setting/profile' },
    // 审计日志：属于敏感数据，必须显式授权 log:list 才出现
    { label: '操作日志', icon: HistoryOutlined, perm: 'log:list', path: '/setting/log' },
  ];

  const entries = computed(() =>
    allEntries.filter(
      (it) => (it.perm ? hasPermission([it.perm]) : true) && (it.show ? it.show() : true)
    )
  );

  function go(it: Entry) {
    if (!it.external) {
      router.push(it.path as string);
      return;
    }
    // 接口文档与前台首页都由 Go 后端提供：开发环境在 8085（见 .env.development 的 VITE_PROXY），
    // 生产由后端同源托管 mgmt，直接用当前 Origin 即可。
    const origin = import.meta.env.DEV ? 'http://127.0.0.1:8085' : location.origin;
    window.open(it.path === 'site' ? origin : origin + it.path, '_blank');
  }

  onMounted(() => {
    // 接口文档开关在认证配置里，/site 已一并下发；取不到就按开启处理，
    // 反正点进去也是 404，不会比直接隐藏更危险。
    getSite()
      .then((res) => (docsOn.value = res?.docs !== false))
      .catch(() => undefined);
    if (hasStatus.value) {
      getStatus()
        .then((res) => (status.value = res || {}))
        .catch(() => undefined);
    }
  });
</script>
