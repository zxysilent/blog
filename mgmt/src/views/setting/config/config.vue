<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="配置管理" />
    </div>

    <n-card :bordered="false" class="mt-4 proCard">
      <n-space class="mb-3" :size="12" align="center">
        <n-select
          v-model:value="groupFilter"
          :options="groupOptions"
          placeholder="全部分组"
          clearable
          style="width: 170px"
        />
        <n-button quaternary @click="groupFilter = null">重置</n-button>
        <n-input
          v-model:value="keyword"
          placeholder="主键 / 名称 / 描述"
          clearable
          style="width: 220px"
        />
        <n-button type="primary" @click="openCreate">
          <template #icon>
            <n-icon><PlusOutlined /></n-icon>
          </template>
          新增配置
        </n-button>
      </n-space>

      <n-spin :show="loading">
        <n-space class="mb-3" :size="8" align="center">
          <n-tag v-for="it in groupStats" :key="it.group" size="small" :bordered="false">
            {{ it.group }} {{ it.count }}
          </n-tag>
        </n-space>
        <n-data-table
          :columns="columns"
          :data="filteredList"
          :row-key="(row: Config) => row.id"
          :pagination="{ pageSize: 20 }"
          :scroll-x="1200"
          size="small"
        />
      </n-spin>
    </n-card>

    <!-- ============ 新增 / 查看 / 编辑（同一个抽屉，三种模式） ============ -->
    <n-drawer v-model:show="showForm" :width="620" placement="right">
      <n-drawer-content :title="formTitle" closable>
        <n-alert v-if="isEdit && cfg.editable" type="info" :showIcon="false" class="mb-3">
          内置配置的值建议在「全局配置」里按表单修改，这里改的是原始值
        </n-alert>

        <n-form
          ref="cfgFormRef"
          :model="cfg"
          :rules="cfgRules"
          label-placement="left"
          label-width="100"
          :disabled="isView"
        >
          <n-form-item label="主键" path="id">
            <n-input
              v-model:value="cfg.id"
              placeholder="如 feature_flag，留空自动生成"
              :disabled="!isCreate"
            />
          </n-form-item>
          <n-form-item label="名称" path="name">
            <n-input v-model:value="cfg.name" placeholder="显示名，如 功能开关" />
          </n-form-item>
          <n-form-item label="分组">
            <n-auto-complete
              v-model:value="cfg.group"
              :options="groupInputOptions"
              placeholder="如 站点 / 系统 / 其它"
              clearable
            />
          </n-form-item>
          <n-form-item label="值类型" path="type">
            <n-select v-model:value="cfg.type" :options="typeOptions" :disabled="!isCreate" />
          </n-form-item>
          <n-grid :cols="3" :x-gap="12">
            <n-form-item-gi label="最小值">
              <n-input-number v-model:value="cfg.min" :precision="0" style="width: 100%" />
            </n-form-item-gi>
            <n-form-item-gi label="最大值">
              <n-input-number v-model:value="cfg.max" :precision="0" style="width: 100%" />
            </n-form-item-gi>
            <n-form-item-gi label="排序">
              <n-input-number v-model:value="cfg.sort" :precision="0" style="width: 100%" />
            </n-form-item-gi>
          </n-grid>
          <div class="field-tip mb-3">
            区间仅对数字型生效（字符串型表示长度），0 表示不限制
          </div>
          <n-form-item label="描述">
            <n-input v-model:value="cfg.describe" type="textarea" :rows="2" placeholder="可选" />
          </n-form-item>
          <n-form-item label="值" path="value">
            <n-input
              v-model:value="cfg.value"
              type="textarea"
              :rows="8"
              placeholder='按类型填写，如 true、12、1.5、{"k":"v"}'
            />
          </n-form-item>
        </n-form>

        <div class="config-meta">
          <span>修改时间：{{ fmtTime(cfg.updated) }}</span>
          <span>创建时间：{{ fmtTime(cfg.created) }}</span>
        </div>

        <template #footer>
          <n-space justify="end">
            <n-button @click="showForm = false">{{ isView ? '关闭' : '取消' }}</n-button>
            <n-button v-if="!isView" type="primary" :loading="submitting" @click="submitForm">
              {{ isCreate ? '确认新增' : '保存' }}
            </n-button>
          </n-space>
        </template>
      </n-drawer-content>
    </n-drawer>
  </div>
</template>

<script lang="ts" setup>
  /**
   * 配置管理
   *
   * 只负责「配置项」本身：清单浏览、按分组筛选、自建配置的增删改。
   * 内置配置的表单化编辑在「全局配置」页，布局配置在右上角抽屉，
   * 三处职责分开，避免同一个值有多个入口各改一半。
   *
   * 查看/编辑/新增共用一个抽屉：三者字段完全一样，拆成两个入口只会让人
   * 疑惑「改名称该点哪个」。
   */
  import { computed, h, onMounted, reactive, ref } from 'vue';
  import { NTag, useDialog, useMessage } from 'naive-ui';
  import type { DataTableColumns } from 'naive-ui';
  import { PlusOutlined, EyeOutlined, EditOutlined, DeleteOutlined } from '@vicons/antd';
  import { TableAction } from '@/components/Table';
  import { addConfig, deleteConfig, editConfig, getConfigList } from '@/api/config';
  import { fmtTime } from '@/api/common';
  import type { Config, ConfigGroup } from '@/api/types';

  const message = useMessage();
  const dialog = useDialog();

  const loading = ref(false);
  const list = ref<Config[]>([]);
  const groupFilter = ref<string | null>(null);
  const keyword = ref('');

  /** 按分组统计数量，供筛选下拉与标签展示 */
  const groupStats = computed<ConfigGroup[]>(() => {
    const mp = new Map<string, number>();
    list.value.forEach((it) => {
      const key = it.group || '其它';
      mp.set(key, (mp.get(key) || 0) + 1);
    });
    return [...mp.entries()].map(([group, count]) => ({ group, count }));
  });

  const groupOptions = computed(() =>
    groupStats.value.map((it) => ({ label: `${it.group}（${it.count}）`, value: it.group }))
  );
  const groupInputOptions = computed(() =>
    groupStats.value.map((it) => ({ label: it.group, value: it.group }))
  );

  const filteredList = computed(() => {
    const kw = keyword.value.trim().toLowerCase();
    return list.value.filter((it) => {
      if (groupFilter.value && it.group !== groupFilter.value) return false;
      if (!kw) return true;
      return (
        it.id.toLowerCase().includes(kw) ||
        (it.name || '').toLowerCase().includes(kw) ||
        (it.describe || '').toLowerCase().includes(kw)
      );
    });
  });

  const typeOptions = [
    { label: '字符串 string', value: 'string' },
    { label: '布尔 bool', value: 'bool' },
    { label: '整数 int', value: 'int' },
    { label: '浮点 float', value: 'float' },
    { label: 'JSON json', value: 'json' },
  ];

  /** 值摘要，太长的只取前 40 字符 */
  function valueBrief(mod: Config): string {
    const v = mod.value || '';
    return v.length > 40 ? `${v.slice(0, 40)}…` : v || '（空）';
  }

  const columns: DataTableColumns<Config> = [
    { title: '主键', key: 'id', width: 160, render: (row: Config) => h('code', row.id) },
    { title: '名称', key: 'name', width: 150 },
    {
      title: '分组',
      key: 'group',
      width: 110,
      render: (row: Config) =>
        h(NTag, { size: 'small', type: 'info', bordered: false }, { default: () => row.group }),
    },
    { title: '类型', key: 'type', width: 90 },
    { title: '值', key: 'value', minWidth: 220, render: (row: Config) => h('code', valueBrief(row)) },
    {
      title: '来源',
      key: 'editable',
      width: 100,
      render: (row: Config) =>
        row.editable
          ? h(NTag, { size: 'small', type: 'warning', bordered: false }, { default: () => '内置' })
          : h(NTag, { size: 'small', bordered: false }, { default: () => '自建' }),
    },
    { title: '修改时间', key: 'updated', width: 170, render: (row: Config) => fmtTime(row.updated) },
    {
      title: '操作',
      key: 'action',
      width: 140,
      fixed: 'right',
      render: (row: Config) =>
        h(TableAction, {
          style: 'button',
          actions: [
            { label: '查看', icon: EyeOutlined, onClick: () => openView(row) },
            { label: '编辑', icon: EditOutlined, onClick: () => openEdit(row) },
          ],
          dropDownActions: [
            {
              label: '删除',
              key: 'delete',
              icon: DeleteOutlined,
              type: 'error',
              ifShow: () => !row.editable,
            },
          ],
          select: (key: string) => {
            if (key === 'delete') confirmDelete(row);
          },
        }),
    },
  ];

  // ---------------- 抽屉：新增 / 查看 / 编辑 ----------------
  type FormMode = 'create' | 'view' | 'edit';

  const showForm = ref(false);
  const formMode = ref<FormMode>('create');
  const submitting = ref(false);
  const cfgFormRef = ref();

  const isCreate = computed(() => formMode.value === 'create');
  const isView = computed(() => formMode.value === 'view');
  const isEdit = computed(() => formMode.value === 'edit');

  const formTitle = computed(() => {
    if (isCreate.value) return '新增配置';
    return isView.value ? '查看配置' : '编辑配置';
  });

  const cfg = reactive({
    id: '',
    group: '',
    name: '',
    type: 'string',
    min: 0,
    max: 0,
    sort: 1000,
    describe: '',
    value: '',
    editable: false,
    updated: 0,
    created: 0,
  });

  const cfgRules = {
    id: {
      trigger: 'blur',
      validator: (_r: any, v: string) => {
        if (!isCreate.value) return true;
        if (!v) return true; // 留空由后端生成
        return /^[A-Za-z0-9_-]{1,64}$/.test(v) ? true : new Error('只能包含字母、数字、下划线和短横线');
      },
    },
    name: { required: true, message: '请输入名称', trigger: 'blur' },
    value: {
      trigger: 'blur',
      validator: (_r: any, v: string) => {
        const t = cfg.type;
        if (t === 'bool' && !['true', 'false'].includes(String(v))) {
          return new Error('布尔值只能是 true 或 false');
        }
        if (t === 'int' && !/^-?\d+$/.test(String(v))) return new Error('请输入整数');
        if (t === 'float' && Number.isNaN(Number(v))) return new Error('请输入数字');
        if (t === 'json') {
          try {
            JSON.parse(v);
          } catch {
            return new Error('不是合法的 JSON');
          }
        }
        return true;
      },
    },
  };

  function fillForm(row: Config) {
    Object.assign(cfg, {
      id: row.id,
      group: row.group,
      name: row.name,
      type: row.type,
      min: row.min ?? 0,
      max: row.max ?? 0,
      sort: row.sort ?? 1000,
      describe: row.describe || '',
      value: row.value || '',
      editable: row.editable,
      updated: row.updated || 0,
      created: row.created || 0,
    });
  }

  function openCreate() {
    formMode.value = 'create';
    Object.assign(cfg, {
      id: '',
      group: groupFilter.value || '其它',
      name: '',
      type: 'string',
      min: 0,
      max: 0,
      sort: 1000,
      describe: '',
      value: '',
      editable: false,
      updated: 0,
      created: 0,
    });
    showForm.value = true;
  }

  function openView(row: Config) {
    formMode.value = 'view';
    fillForm(row);
    showForm.value = true;
  }

  function openEdit(row: Config) {
    formMode.value = 'edit';
    fillForm(row);
    showForm.value = true;
  }

  async function submitForm() {
    try {
      await cfgFormRef.value?.validate();
    } catch {
      return;
    }
    if (cfg.max !== 0 && cfg.min > cfg.max) {
      message.warning('最小值不能大于最大值');
      return;
    }
    submitting.value = true;
    try {
      if (isCreate.value) {
        await addConfig({
          id: cfg.id,
          group: cfg.group,
          name: cfg.name,
          type: cfg.type,
          min: cfg.min,
          max: cfg.max,
          sort: cfg.sort,
          describe: cfg.describe,
          value: cfg.value,
        });
        message.success('新增成功');
      } else {
        // 主键与值类型创建后不可改，其余元信息与值一并提交
        await editConfig(cfg.id, {
          name: cfg.name,
          group: cfg.group,
          min: cfg.min,
          max: cfg.max,
          sort: cfg.sort,
          describe: cfg.describe,
          value: cfg.value,
        });
        message.success('修改成功');
      }
      showForm.value = false;
      await load();
    } finally {
      submitting.value = false;
    }
  }

  // ---------------- 删除 ----------------
  async function removeConfig(id: string) {
    await deleteConfig(id);
    message.success('删除成功');
    await load();
  }

  function confirmDelete(row: Config) {
    if (row.editable) {
      message.warning('内置配置不可删除');
      return;
    }
    dialog.warning({
      title: '删除确认',
      content: `确定删除配置「${row.name}」？站点可能依赖该配置`,
      positiveText: '删除',
      negativeText: '取消',
      onPositiveClick: () => removeConfig(row.id),
    });
  }

  async function load() {
    loading.value = true;
    try {
      const res = await getConfigList({ page: 1, size: 1000 });
      list.value = res.list;
    } finally {
      loading.value = false;
    }
  }

  onMounted(load);
</script>

<style lang="less" scoped>
  .field-tip {
    margin-top: -8px;
    font-size: 12px;
    line-height: 1.5;
    color: var(--text-color-3);
  }

  .config-meta {
    display: flex;
    gap: 16px;
    margin-top: 16px;
    font-size: 12px;
    color: var(--text-color-3);
  }
</style>
