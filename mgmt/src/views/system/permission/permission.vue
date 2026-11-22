<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="权限标识" />
    </div>

    <n-card :bordered="false" class="mt-4 proCard">
      <n-space class="mb-3" :size="12" align="center">
        <n-radio-group v-model:value="viewMode">
          <n-radio-button value="table">列表</n-radio-button>
          <n-radio-button value="tree">分组树</n-radio-button>
        </n-radio-group>
        <n-select
          v-model:value="filter.group"
          :options="groupOptions"
          placeholder="全部分组"
          clearable
          style="width: 170px"
          @update:value="reload"
        />
        <n-input
          v-model:value="filter.q"
          placeholder="权限标识 / 名称"
          clearable
          style="width: 220px"
          @keyup.enter="reload"
        />
        <n-button type="primary" ghost @click="reload">查询</n-button>
        <n-button quaternary @click="resetFilter">重置</n-button>
      </n-space>

      <!-- 列表视图 -->
      <BasicTable
        v-if="viewMode === 'table'"
        :columns="columns"
        :request="loadDataTable"
        :row-key="(row: Permission) => row.id"
        ref="actionRef"
      />

      <!-- 分组树视图 -->
      <n-spin v-else :show="treeLoading">
        <n-empty v-if="!treeRows.length" description="没有匹配的权限标识" class="py-8" />
        <n-data-table
          v-else
          :columns="treeColumns"
          :data="treeRows"
          :row-key="(row: any) => row.id"
          :pagination="false"
          :default-expand-all="true"
          size="small"
        />
      </n-spin>
    </n-card>

    <!-- 详情 -->
    <n-modal v-model:show="showDetail" preset="card" title="权限详情" style="width: 520px">
      <n-descriptions v-if="detail" :column="1" label-placement="left" bordered>
        <n-descriptions-item label="分组">{{ detail.group }}</n-descriptions-item>
        <n-descriptions-item label="分组键">{{ detail.group_key || '-' }}</n-descriptions-item>
        <n-descriptions-item label="权限标识">
          <n-text code>{{ detail.perm }}</n-text>
        </n-descriptions-item>
        <n-descriptions-item label="名称">{{ detail.name }}</n-descriptions-item>
        <n-descriptions-item label="排序">{{ detail.sort }}</n-descriptions-item>
        <n-descriptions-item label="主键">
          <n-text depth="3">{{ detail.id }}</n-text>
        </n-descriptions-item>
        <n-descriptions-item label="创建时间">{{ fmtTime(detail.created) }}</n-descriptions-item>
      </n-descriptions>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showDetail = false">关闭</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
  import { h, onMounted, reactive, ref, watch } from 'vue';
  import { NTag, useMessage } from 'naive-ui';
  import { BasicTable, TableAction } from '@/components/Table';
  import { EyeOutlined, CopyOutlined } from '@vicons/antd';
  import { getPermissionList } from '@/api/permission';
  import { getPermissionTree } from '@/api/auth';
  import { fmtTime } from '@/api/common';
  import type { Permission, PermissionNode } from '@/api/types';

  const actionRef = ref();
  const viewMode = ref<'table' | 'tree'>('table');
  const filter = reactive({ group: null as string | null, q: '' });

  /** 分组概览与下拉选项都来自权限树，避免额外接口 */
  const groups = ref<{ key: string; name: string; total: number }[]>([]);
  const groupOptions = ref<{ label: string; value: string }[]>([]);

  const columns = [
    {
      title: '分组',
      key: 'group',
      width: 140,
      render: (row: Permission) =>
        h(NTag, { size: 'small', type: 'info' }, { default: () => row.group }),
    },
    {
      title: '权限标识',
      key: 'perm',
      width: 200,
      render: (row: Permission) => h('code', row.perm),
    },
    { title: '名称', key: 'name', minWidth: 180 },
    { title: '排序', key: 'sort', width: 90 },
    {
      title: '创建时间',
      key: 'created',
      width: 170,
      render: (row: Permission) => fmtTime(row.created),
    },
    {
      title: '操作',
      key: 'action',
      width: 100,
      fixed: 'right',
      render: (row: Permission) =>
        h(TableAction, {
          style: 'button',
          actions: [{ label: '查看', icon: EyeOutlined, onClick: () => openDetail(row) }],
          dropDownActions: [{ label: '复制标识', key: 'copy', icon: CopyOutlined }],
          select: (key: string) => {
            if (key === 'copy') copyPerm(row.perm);
          },
        }),
    },
  ];

  const treeColumns = [
    { title: '分组 / 权限', key: 'name', minWidth: 260 },
    {
      title: '权限标识',
      key: 'perm',
      width: 200,
      render: (row: any) => (row.perm ? h('code', row.perm) : '-'),
    },
    {
      title: '分组键',
      key: 'group_key',
      width: 130,
      render: (row: any) => row.group_key || '-',
    },
    {
      title: '排序',
      key: 'sort',
      width: 90,
      render: (row: any) => row.sort ?? '-',
    },
    {
      title: '类型',
      key: 'leaf',
      width: 110,
      render: (row: any) =>
        row.leaf
          ? h(NTag, { size: 'small', type: 'default' }, { default: () => '权限' })
          : h(NTag, { size: 'small', type: 'info' }, { default: () => '分组' }),
    },
  ];

  const treeLoading = ref(false);
  const treeRows = ref<any[]>([]);

  const loadDataTable = async (res: any) => {
    const params: Recordable = { ...res };
    if (filter.group) params.group = filter.group;
    if (filter.q) params.q = filter.q;
    return await getPermissionList(params);
  };

  function reload() {
    actionRef.value?.reload();
    if (viewMode.value === 'tree') loadTree();
  }

  function resetFilter() {
    filter.group = null;
    filter.q = '';
    reload();
  }

  /** 把权限树摊平成可直接渲染的两级表格 */
  function toRows(nodes: PermissionNode[]): any[] {
    return (nodes || []).map((grp) => ({
      id: `group:${grp.group_key || grp.group}`,
      name: grp.name,
      group_key: grp.group_key,
      perm: '',
      sort: grp.sort,
      leaf: false,
      children: (grp.children || []).map((it) => ({
        id: it.id,
        name: it.name,
        group_key: it.group_key,
        perm: it.perm,
        sort: it.sort,
        leaf: true,
      })),
    }));
  }

  async function loadTree() {
    treeLoading.value = true;
    try {
      const tree = await getPermissionTree();
      // 分组维度也支持上面的筛选条件
      const kw = filter.q.trim().toLowerCase();
      let rows = toRows(tree);
      if (filter.group) rows = rows.filter((it) => it.group_key === filter.group);
      if (kw) {
        rows = rows
          .map((it) => ({
            ...it,
            children: it.children.filter(
              (c: any) =>
                String(c.name).toLowerCase().includes(kw) ||
                String(c.perm).toLowerCase().includes(kw)
            ),
          }))
          .filter((it) => it.children.length || String(it.name).toLowerCase().includes(kw));
      }
      treeRows.value = rows;
    } finally {
      treeLoading.value = false;
    }
  }

  const showDetail = ref(false);
  const detail = ref<Permission | null>(null);
  const message = useMessage();

  function openDetail(row: Permission) {
    detail.value = row;
    showDetail.value = true;
  }

  /** 权限标识常被后端代码引用，提供一键复制 */
  async function copyPerm(perm: string) {
    try {
      await navigator.clipboard.writeText(perm);
      message.success('已复制：' + perm);
    } catch {
      message.warning('浏览器不支持剪贴板，请手动复制');
    }
  }

  onMounted(async () => {
    const tree = await getPermissionTree().catch(() => []);
    groups.value = (tree || []).map((grp) => ({
      key: grp.group_key || grp.group,
      name: grp.name,
      total: grp.children?.length || 0,
    }));
    groupOptions.value = groups.value.map((it) => ({ label: `${it.name}（${it.total}）`, value: it.key }));
    if (viewMode.value === 'tree') loadTree();
  });

  // 切到树视图时才拉数据，避免列表页多一次请求
  watch(viewMode, (mode) => {
    if (mode === 'tree') loadTree();
  });
</script>

<style lang="less" scoped>
  /* 表格横向滚动时操作列不被内容挤压 */
  :deep(.n-data-table) {
    .n-data-table-td,
    .n-data-table-th {
      white-space: nowrap;
    }
  }
</style>
