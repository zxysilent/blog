<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="博文管理" />
    </div>

    <!-- 概览：点卡片直接切筛选 -->
    <n-grid :cols="4" :x-gap="12" :y-gap="12" class="mt-4">
      <n-grid-item v-for="it in statCards" :key="it.key">
        <n-card
          :bordered="false"
          size="small"
          hoverable
          class="cursor-pointer"
          :class="{ 'stat-active': it.active }"
          @click="it.pick"
        >
          <n-statistic :label="it.label" :value="it.count" />
        </n-card>
      </n-grid-item>
    </n-grid>

    <n-card :bordered="false" class="mt-4 proCard">
      <n-space class="mb-3" :size="12" align="center">
        <n-select
          v-model:value="filter.kind"
          :options="kindOptions"
          style="width: 120px"
          @update:value="reload"
        />
        <n-tree-select
          v-model:value="filter.cate_id"
          :options="cateOptions"
          placeholder="全部栏目"
          clearable
          key-field="id"
          label-field="name"
          children-field="children"
          style="width: 180px"
          @update:value="reload"
        />
        <n-select
          v-model:value="filter.status"
          :options="statusOptions"
          style="width: 120px"
          @update:value="reload"
        />
        <n-select
          v-model:value="filter.istop"
          :options="istopOptions"
          style="width: 120px"
          @update:value="reload"
        />
        <n-select
          v-model:value="filter.show"
          :options="showOptions"
          style="width: 130px"
          @update:value="reload"
        />
        <n-input
          v-model:value="filter.q"
          placeholder="标题 / 路径 / 摘要"
          clearable
          style="width: 200px"
          @keyup.enter="reload"
        />
        <n-button type="primary" ghost @click="reload">查询</n-button>
        <n-button quaternary @click="resetFilter">重置</n-button>
        <n-button type="primary" @click="goAdd">
          <template #icon>
            <n-icon><PlusOutlined /></n-icon>
          </template>
          新增博文
        </n-button>
      </n-space>

      <!-- 批量操作：勾选后出现 -->
      <n-space v-if="checkedKeys.length" class="mb-3" :size="8" align="center">
        <n-tag type="info" :bordered="false">已选 {{ checkedKeys.length }} 项</n-tag>
        <n-button size="small" type="success" ghost :loading="batching" @click="batchStatus(2)">
          批量发布
        </n-button>
        <n-button size="small" ghost :loading="batching" @click="batchStatus(1)">批量撤回</n-button>
        <n-button size="small" ghost :loading="batching" @click="batchTop(true)">批量置顶</n-button>
        <n-button size="small" ghost :loading="batching" @click="batchTop(false)">
          取消置顶
        </n-button>
        <n-button size="small" type="error" ghost :loading="batching" @click="confirmDelete(checkedKeys)">
          批量删除
        </n-button>
        <n-button size="small" quaternary @click="checkedKeys = []">清空选择</n-button>
      </n-space>

      <BasicTable
        v-model:checked-row-keys="checkedKeys"
        :columns="columns"
        :request="loadDataTable"
        :row-key="(row: Post) => row.id"
        ref="actionRef"
        :action-column="actionColumn"
        :scroll-x="1420"
      />
    </n-card>
  </div>
</template>

<script lang="ts" setup>
  import { computed, h, onMounted, reactive, ref } from 'vue';
  import { useRouter } from 'vue-router';
  import { NTag, useDialog, useMessage } from 'naive-ui';
  import { BasicTable, TableAction } from '@/components/Table';
  import {
    PlusOutlined,
    EyeOutlined,
    EditOutlined,
    DeleteOutlined,
    CheckCircleOutlined,
    StopOutlined,
    VerticalAlignTopOutlined,
  } from '@vicons/antd';
  import { deletePost, editPost, getPostList, setPostTop } from '@/api/post';
  import { getCateTree } from '@/api/cate';
  import { fmtTime } from '@/api/common';
  import type { Cate, Post } from '@/api/types';

  const router = useRouter();
  const message = useMessage();
  const dialog = useDialog();
  const actionRef = ref();
  const cateOptions = ref<Cate[]>([]);
  const checkedKeys = ref<string[]>([]);
  const batching = ref(false);

  /** 与后端 model.Kind* / Status* 一致 */
  const kindMap: Record<number, { text: string; type: 'info' | 'success' | 'default' }> = {
    1: { text: '文章', type: 'info' },
    2: { text: '页面', type: 'default' },
    3: { text: '笔记', type: 'success' },
  };
  const kindOptions = [
    { label: '全部类型', value: 0 },
    { label: '文章', value: 1 },
    { label: '页面', value: 2 },
    { label: '笔记', value: 3 },
  ];
  const statusOptions = [
    { label: '全部状态', value: 0 },
    { label: '已发布', value: 2 },
    { label: '草稿', value: 1 },
  ];
  // 布尔筛选在后端是指针，前端用 0/1/2 三态表达「不限 / 是 / 否」
  const istopOptions = [
    { label: '置顶不限', value: 0 },
    { label: '已置顶', value: 1 },
    { label: '未置顶', value: 2 },
  ];
  const showOptions = [
    { label: '前台不限', value: 0 },
    { label: '前台可见', value: 1 },
    { label: '前台隐藏', value: 2 },
  ];

  const filter = reactive({
    kind: 0,
    cate_id: null as string | null,
    status: 0,
    istop: 0,
    show: 0,
    q: '',
  });

  /** 概览计数，按当前「类型 + 栏目」统计，点卡片切筛选 */
  const stats = reactive({ all: 0, published: 0, draft: 0, top: 0 });

  const statCards = computed(() => [
    {
      key: 'all',
      label: '全部内容',
      count: stats.all,
      active: filter.status === 0 && filter.istop === 0,
      pick: () => {
        filter.status = 0;
        filter.istop = 0;
        reload();
      },
    },
    {
      key: 'published',
      label: '已发布',
      count: stats.published,
      active: filter.status === 2,
      pick: () => {
        filter.status = 2;
        reload();
      },
    },
    {
      key: 'draft',
      label: '草稿',
      count: stats.draft,
      active: filter.status === 1,
      pick: () => {
        filter.status = 1;
        reload();
      },
    },
    {
      key: 'top',
      label: '置顶',
      count: stats.top,
      active: filter.istop === 1,
      pick: () => {
        filter.istop = 1;
        reload();
      },
    },
  ]);

  const columns = [
    { type: 'selection' },
    {
      title: '封面',
      key: 'icon',
      width: 74,
      render: (row: Post) =>
        row.icon
          ? h('img', {
              src: row.icon,
              style: 'width:44px;height:30px;object-fit:cover;border-radius:4px',
            })
          : h('span', { class: 'text-gray-400' }, '-'),
    },
    {
      title: '标题',
      key: 'title',
      minWidth: 260,
      render: (row: Post) =>
        h('span', { class: 'flex items-center gap-1' }, [
          row.istop
            ? h(NTag, { size: 'tiny', type: 'error', bordered: false }, { default: () => '顶' })
            : null,
          row.title,
          row.inner
            ? h(NTag, { size: 'tiny', type: 'warning', bordered: false }, { default: () => '内置' })
            : null,
        ]),
    },
    {
      title: '类型',
      key: 'kind',
      width: 84,
      render: (row: Post) => {
        const it = kindMap[row.kind] || { text: '-', type: 'default' as const };
        return h(NTag, { size: 'small', type: it.type }, { default: () => it.text });
      },
    },
    {
      title: '状态',
      key: 'status',
      width: 90,
      render: (row: Post) =>
        h(
          NTag,
          { size: 'small', type: row.status === 2 ? 'success' : 'default' },
          { default: () => (row.status === 2 ? '已发布' : '草稿') }
        ),
    },
    {
      title: '前台',
      key: 'show',
      width: 88,
      render: (row: Post) =>
        h(
          NTag,
          { size: 'small', type: row.show ? 'success' : 'default', bordered: false },
          { default: () => (row.show ? '可见' : '隐藏') }
        ),
    },
    {
      title: '栏目',
      key: 'cate',
      width: 130,
      // 列表接口只返回 cate_id，这里用已加载的栏目树做映射
      render: (row: Post) => cateName(row.cate_id),
    },
    { title: '路径', key: 'path', width: 170 },
    { title: '浏览', key: 'hits', width: 80 },
    {
      title: '更新时间',
      key: 'updated',
      width: 170,
      render: (row: Post) => fmtTime(row.updated),
    },
  ];

  const actionColumn = reactive({
    width: 140,
    title: '操作',
    key: 'action',
    fixed: 'right',
    render(record: Post) {
      return h(TableAction, {
        style: 'button',
        // 默认三个：查看 / 编辑 / 更多，其余操作统一收进下拉
        actions: [
          {
            label: '查看',
            icon: EyeOutlined,
            onClick: () => router.push({ name: 'content_post_detail', params: { id: record.id } }),
          },
          {
            label: '编辑',
            icon: EditOutlined,
            onClick: () => router.push({ name: 'content_post_edit', params: { id: record.id } }),
          },
        ],
        dropDownActions: [
          {
            label: record.status === 2 ? '撤回为草稿' : '发布',
            key: 'status',
            icon: record.status === 2 ? StopOutlined : CheckCircleOutlined,
          },
          {
            label: record.istop ? '取消置顶' : '置顶',
            key: 'top',
            icon: VerticalAlignTopOutlined,
          },
          { label: '删除', key: 'delete', icon: DeleteOutlined, type: 'error' },
        ],
        // 下拉项统一在这里分发，避免给每个 option 单独绑 onClick
        select: (key: string) => {
          if (key === 'status') toggleStatus(record);
          if (key === 'top') toggleTop(record);
          if (key === 'delete') confirmDelete([record.id]);
        },
      });
    },
  });

  /** 栏目 id 转名称 */
  function cateName(id: string): string {
    if (!id) return '-';
    const hit = findCate(cateOptions.value, id);
    return hit?.name || '-';
  }

  function findCate(list: Cate[], id: string): Cate | undefined {
    for (const it of list) {
      if (it.id === id) return it;
      const sub = findCate(it.children || [], id);
      if (sub) return sub;
    }
    return undefined;
  }

  /** 把界面上的筛选条件翻译成后端查询参数 */
  function buildParams(extra: Recordable = {}): Recordable {
    const params: Recordable = { ...extra };
    if (filter.kind > 0) params.kind = filter.kind;
    if (filter.status > 0) params.status = filter.status;
    if (filter.cate_id) params.cate_id = filter.cate_id;
    if (filter.istop === 1) params.istop = true;
    if (filter.istop === 2) params.istop = false;
    if (filter.show === 1) params.show = true;
    if (filter.show === 2) params.show = false;
    if (filter.q) params.q = filter.q;
    return params;
  }

  const loadDataTable = async (res: any) => {
    const page = await getPostList(buildParams({ ...res }));
    refreshStats();
    return page;
  };

  /** 概览只按类型与栏目统计，其他筛选不参与，避免点卡片时数字乱跳 */
  async function refreshStats() {
    const base = buildParams();
    delete base.status;
    delete base.istop;
    delete base.show;
    delete base.q;
    const [all, published, draft, top] = await Promise.all([
      getPostList({ ...base, page: 1, size: 1 }),
      getPostList({ ...base, page: 1, size: 1, status: 2 }),
      getPostList({ ...base, page: 1, size: 1, status: 1 }),
      getPostList({ ...base, page: 1, size: 1, istop: true }),
    ]);
    stats.all = all.itemCount;
    stats.published = published.itemCount;
    stats.draft = draft.itemCount;
    stats.top = top.itemCount;
  }

  function reload() {
    checkedKeys.value = [];
    actionRef.value?.reload();
  }

  function resetFilter() {
    Object.assign(filter, { kind: 0, cate_id: null, status: 0, istop: 0, show: 0, q: '' });
    reload();
  }

  function goAdd() {
    router.push({ name: 'content_post_add' });
  }

  async function toggleTop(row: Post) {
    await setPostTop(row.id, !row.istop);
    message.success(row.istop ? '已取消置顶' : '已置顶');
    reload();
  }

  async function toggleStatus(row: Post) {
    const next = row.status === 2 ? 1 : 2;
    await editPost(row.id, { status: next });
    message.success(next === 2 ? '已发布' : '已撤回为草稿');
    reload();
  }

  function confirmDelete(ids: string[]) {
    const tip = ids.length > 1 ? `确定删除选中的 ${ids.length} 条内容？` : '确定删除该内容？';
    dialog.warning({
      title: '删除确认',
      content: `${tip}标签关联会一并清理，删除后不可恢复。`,
      positiveText: '删除',
      negativeText: '取消',
      onPositiveClick: async () => {
        await runBatch([...ids], (id) => deletePost(id));
        message.success('删除成功');
        reload();
      },
    });
  }

  async function batchStatus(status: number) {
    await runBatch([...checkedKeys.value], (id) => editPost(id, { status }));
    message.success(status === 2 ? '批量发布完成' : '批量撤回完成');
    reload();
  }

  async function batchTop(istop: boolean) {
    await runBatch([...checkedKeys.value], (id) => setPostTop(id, istop));
    message.success(istop ? '批量置顶完成' : '批量取消置顶完成');
    reload();
  }

  /** 批量执行：逐条请求，汇总失败数，单条失败不中断整体 */
  async function runBatch(ids: string[], fn: (id: string) => Promise<any>) {
    batching.value = true;
    try {
      const res = await Promise.allSettled(ids.map((id) => fn(id)));
      const failed = res.filter((it) => it.status === 'rejected').length;
      if (failed) message.warning(`${failed} 条处理失败`);
    } finally {
      batching.value = false;
    }
  }

  onMounted(async () => {
    cateOptions.value = await getCateTree({});
  });
</script>

<style lang="less" scoped>
  .stat-active {
    box-shadow: 0 0 0 1px var(--primary-color);
  }
</style>
