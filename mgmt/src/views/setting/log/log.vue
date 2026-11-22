<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="操作日志" />
    </div>

    <n-card :bordered="false" class="mt-4 proCard">
      <n-space class="mb-3" :size="12" align="center">
        <n-input
          v-model:value="filter.q"
          placeholder="模块 / 动作 / 路径"
          clearable
          style="width: 200px"
          @keyup.enter="reload"
        />
        <n-select
          v-model:value="filter.method"
          :options="methodOptions"
          placeholder="全部方法"
          clearable
          style="width: 130px"
          @update:value="reload"
        />
        <n-select
          v-model:value="filter.status"
          :options="statusOptions"
          placeholder="全部结果"
          clearable
          style="width: 130px"
          @update:value="reload"
        />
        <n-date-picker
          v-model:value="range"
          type="datetimerange"
          clearable
          start-placeholder="开始时间"
          end-placeholder="结束时间"
          style="width: 360px"
          @update:value="reload"
        />
        <n-button type="primary" ghost @click="reload">查询</n-button>
        <n-button quaternary @click="resetFilter">重置</n-button>
        <n-button
          v-if="hasPermission(['log:delete'])"
          type="error"
          secondary
          :loading="clearing"
          @click="confirmClear"
        >
          <template #icon>
            <n-icon><DeleteOutlined /></n-icon>
          </template>
          清空日志
        </n-button>
      </n-space>

      <BasicTable
        :columns="columns"
        :request="loadDataTable"
        :row-key="(row: Log) => row.id"
        ref="actionRef"
        :action-column="actionColumn"
        :scroll-x="1500"
      />
    </n-card>

    <!-- 详情：请求参数已脱敏，这里做一次格式化只为了让 JSON 更好读 -->
    <n-drawer v-model:show="showDetail" :width="560" placement="right">
      <n-drawer-content :title="`日志详情 - ${detail?.action || ''}`" closable>
        <n-descriptions v-if="detail" :column="1" label-placement="left" size="small" bordered>
          <n-descriptions-item label="操作人">
            {{ detail.user?.name || '-' }}
            <n-text v-if="detail.user?.num" depth="3" class="ml-1">{{ detail.user.num }}</n-text>
          </n-descriptions-item>
          <n-descriptions-item label="模块">{{ detail.module || '-' }}</n-descriptions-item>
          <n-descriptions-item label="动作">{{ detail.action || '-' }}</n-descriptions-item>
          <n-descriptions-item label="权限标识">
            <n-text code>{{ detail.perm || '-' }}</n-text>
          </n-descriptions-item>
          <n-descriptions-item label="请求">
            <n-tag size="tiny" :bordered="false">{{ detail.method }}</n-tag>
            <n-text code class="ml-1">{{ detail.path }}</n-text>
          </n-descriptions-item>
          <n-descriptions-item label="来源 ip">{{ detail.ip || '-' }}</n-descriptions-item>
          <n-descriptions-item label="结果">
            <n-tag size="tiny" :bordered="false" :type="detail.status === 200 ? 'success' : 'error'">
              {{ statusText(detail.status) }}
            </n-tag>
          </n-descriptions-item>
          <n-descriptions-item label="时间">{{ fmtTime(detail.created) }}</n-descriptions-item>
        </n-descriptions>

        <n-divider title-placement="left">请求参数</n-divider>
        <n-empty v-if="!detail?.data" description="无参数" class="py-4" />
        <pre v-else class="log-data">{{ prettyData }}</pre>
      </n-drawer-content>
    </n-drawer>
  </div>
</template>

<script lang="ts" setup>
  import { computed, h, reactive, ref } from 'vue';
  import { NTag, useDialog, useMessage } from 'naive-ui';
  import { BasicTable, TableAction } from '@/components/Table';
  import { DeleteOutlined, EyeOutlined } from '@vicons/antd';
  import { clearLogs, deleteLog, getLogList } from '@/api/log';
  import { fmtTime } from '@/api/common';
  import { usePermission } from '@/hooks/usePermission';
  import { ResultEnum } from '@/enums/httpEnum';
  import type { Log } from '@/api/types';

  /**
   * 操作日志
   *
   * 日志由后端审计中间件自动写入（只记写接口与登录），这里只做查询、查看与清理。
   * 请求参数里的密码类字段在后端已替换为 ***，前端不再做二次处理。
   */
  const message = useMessage();
  const dialog = useDialog();
  const { hasPermission } = usePermission();
  const actionRef = ref();
  const clearing = ref(false);
  const showDetail = ref(false);
  const detail = ref<Log | null>(null);
  const range = ref<[number, number] | null>(null);

  const filter = reactive({
    q: '',
    method: null as string | null,
    status: null as number | null,
  });

  const methodOptions = ['POST', 'PUT', 'DELETE', 'GET'].map((it) => ({ label: it, value: it }));
  const statusOptions = [
    { label: '成功', value: ResultEnum.SUCCESS },
    { label: '失败', value: ResultEnum.FAIL },
    { label: '无权限', value: ResultEnum.DENY },
    { label: '令牌失效', value: ResultEnum.TOKEN },
  ];

  function statusText(code: number) {
    const hit = statusOptions.find((it) => it.value === code);
    return hit ? hit.label : String(code);
  }

  const columns = [
    { title: '时间', key: 'created', width: 180, render: (row: Log) => fmtTime(row.created) },
    {
      title: '操作人',
      key: 'user',
      width: 130,
      render: (row: Log) => row.user?.name || row.user?.num || '-',
    },
    { title: '模块', key: 'module', width: 120, ellipsis: { tooltip: true } },
    // key 不能叫 action：操作列也用 action 作 key，重名会顶掉右侧操作列
    { title: '动作', key: 'act', width: 140, ellipsis: { tooltip: true }, render: (row: Log) => row.action },
    {
      title: '权限标识',
      key: 'perm',
      width: 170,
      ellipsis: { tooltip: true },
      render: (row: Log) =>
        row.perm ? h('code', { class: 'log-perm' }, row.perm) : '-',
    },
    {
      title: '请求',
      key: 'method',
      width: 100,
      render: (row: Log) => h(NTag, { size: 'tiny', bordered: false }, { default: () => row.method }),
    },
    { title: '路径', key: 'path', width: 220, ellipsis: { tooltip: true } },
    { title: '来源 ip', key: 'ip', width: 130 },
    {
      title: '结果',
      key: 'status',
      width: 100,
      render: (row: Log) =>
        h(
          NTag,
          { size: 'tiny', bordered: false, type: row.status === 200 ? 'success' : 'error' },
          { default: () => statusText(row.status) }
        ),
    },
  ];

  const actionColumn = reactive({
    width: 100,
    title: '操作',
    key: 'action',
    fixed: 'right',
    render(record: Log) {
      // 主操作只留「查看」，删除放进更多，操作列保持窄
      const actions: any[] = [{ label: '查看', icon: EyeOutlined, onClick: () => openDetail(record) }];
      const dropDownActions: any[] = [
        {
          label: '删除',
          key: 'delete',
          icon: DeleteOutlined,
          auth: ['log:delete'],
          type: 'error',
        },
      ];
      const select = (key: string) => {
        if (key === 'delete') confirmDelete(record);
      };
      return h(TableAction, { style: 'button', actions, dropDownActions, select });
    },
  });

  const loadDataTable = async (res: any) => {
    const params: Recordable = { ...res };
    if (filter.q) params.q = filter.q;
    if (filter.method) params.method = filter.method;
    if (filter.status) params.status = filter.status;
    if (range.value?.length === 2) {
      params.from = range.value[0];
      params.to = range.value[1];
    }
    return await getLogList(params);
  };

  function reload() {
    actionRef.value?.reload();
  }

  function resetFilter() {
    filter.q = '';
    filter.method = null;
    filter.status = null;
    range.value = null;
    reload();
  }

  function openDetail(row: Log) {
    detail.value = row;
    showDetail.value = true;
  }

  /**
   * 日志里的 data 是后端序列化后的字符串，尽量格式化展示，失败就原样输出。
   * 配置类接口会把整段 JSON 当成字符串字段提交，这里再解一层，
   * 避免出现「{\"value\":\"{\\\"nav_mode\\\"...}\"}」这种双层转义。
   */
  function deepParse(val: any, depth = 0): any {
    if (depth > 3) return val;
    if (typeof val === 'string') {
      const s = val.trim();
      if (!s || !/^[[{]/.test(s)) return val;
      try {
        return deepParse(JSON.parse(s), depth + 1);
      } catch {
        return val;
      }
    }
    if (Array.isArray(val)) return val.map((it) => deepParse(it, depth + 1));
    if (val && typeof val === 'object') {
      const out: Recordable = {};
      Object.keys(val).forEach((key) => (out[key] = deepParse(val[key], depth + 1)));
      return out;
    }
    return val;
  }

  const prettyData = computed(() => {
    const raw = detail.value?.data || '';
    if (!raw) return '';
    try {
      return JSON.stringify(deepParse(JSON.parse(raw)), null, 2);
    } catch {
      return raw;
    }
  });

  function confirmDelete(row: Log) {
    dialog.warning({
      title: '删除日志',
      content: `确定删除该条日志（${row.module} / ${row.action}）吗？`,
      positiveText: '删除',
      negativeText: '取消',
      onPositiveClick: async () => {
        await deleteLog(row.id);
        message.success('删除成功');
        reload();
      },
    });
  }

  function confirmClear() {
    dialog.warning({
      title: '清空日志',
      content: '将删除全部操作日志，且不可恢复，确定继续吗？',
      positiveText: '清空',
      negativeText: '取消',
      onPositiveClick: async () => {
        clearing.value = true;
        try {
          await clearLogs();
          message.success('清空成功');
          reload();
        } finally {
          clearing.value = false;
        }
      },
    });
  }
</script>

<style lang="less" scoped>
  .log-perm {
    font-family: ui-monospace, sfmono-regular, menlo, consolas, monospace;
  }

  .log-data {
    max-height: 420px;
    padding: 12px;
    overflow: auto;
    font-family: ui-monospace, sfmono-regular, menlo, consolas, monospace;
    font-size: 12px;
    line-height: 1.7;
    white-space: pre-wrap;
    word-break: break-all;
    background: var(--action-color);
    border-radius: 4px;
  }
</style>
