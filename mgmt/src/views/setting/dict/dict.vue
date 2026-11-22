<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="配置字典" />
    </div>

    <n-card :bordered="false" class="mt-4 proCard">
      <BasicTable
        :columns="columns"
        :request="loadDataTable"
        :row-key="(row: Dict) => row.id"
        ref="actionRef"
        :action-column="actionColumn"
      >
        <template #tableTitle>
          <n-space :size="10" align="center">
            <n-input
              v-model:value="search.q"
              placeholder="主键 / 名称 / 备注"
              clearable
              style="width: 200px"
              @keyup.enter="reload"
            />
            <n-select
              v-model:value="search.type"
              :options="typeOptions"
              placeholder="全部类型"
              clearable
              style="width: 140px"
              @update:value="reload"
            />
            <n-select
              v-model:value="search.status"
              :options="statusOptions"
              style="width: 130px"
              @update:value="reload"
            />
            <n-button type="primary" ghost @click="reload">查询</n-button>
            <n-button quaternary @click="resetSearch">重置</n-button>
            <n-button type="primary" @click="openAdd">
              <template #icon>
                <n-icon><PlusOutlined /></n-icon>
              </template>
              新增字典
            </n-button>
          </n-space>
        </template>
      </BasicTable>
    </n-card>

    <!-- 新增 / 编辑共用抽屉：选项较多，横向空间比弹窗更合适 -->
    <n-drawer v-model:show="showForm" :width="820" placement="right">
      <n-drawer-content :title="editId ? `${readonly ? '查看' : '编辑'}字典 - ${editId}` : '新增字典'" closable>
        <n-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-placement="left"
          label-width="90"
          :disabled="readonly"
        >
          <n-grid :cols="24" :x-gap="16">
            <n-form-item-gi :span="12" label="主键" path="id">
              <n-input
                v-model:value="form.id"
                placeholder="如 gender、page_size"
                :disabled="!!editId"
              />
            </n-form-item-gi>
            <n-form-item-gi :span="12" label="名称" path="name">
              <n-input v-model:value="form.name" placeholder="字典名称" />
            </n-form-item-gi>
            <n-form-item-gi :span="12" label="值类型" path="type">
              <n-select v-model:value="form.type" :options="dictTypes" />
            </n-form-item-gi>
            <n-form-item-gi :span="12" label="状态">
              <n-switch v-model:value="form.status">
                <template #checked>启用</template>
                <template #unchecked>禁用</template>
              </n-switch>
            </n-form-item-gi>
            <n-form-item-gi :span="24" label="备注">
              <n-input v-model:value="form.remark" placeholder="可选" />
            </n-form-item-gi>
          </n-grid>
        </n-form>

        <n-divider title-placement="left">字典选项</n-divider>
        <n-alert v-if="form.type === 'time'" type="info" :bordered="false" class="mb-2">
          {{ dictTimeHint }}
        </n-alert>

        <n-space v-if="!readonly" class="mb-2" :size="8" align="center">
          <n-button dashed size="small" @click="addRow">
            <template #icon>
              <n-icon><PlusOutlined /></n-icon>
            </template>
            添加一行
          </n-button>
          <n-button size="small" secondary type="primary" @click="resortRows">
            <template #icon>
              <n-icon><OrderedListOutlined /></n-icon>
            </template>
            按顺序重排序号
          </n-button>
          <n-text depth="3" class="text-sm">选项按序号从小到大在前端展示</n-text>
        </n-space>

        <n-table :bordered="false" size="small">
          <thead>
            <tr>
              <th style="width: 150px">字段名</th>
              <th style="width: 150px">展示名</th>
              <th>字典值</th>
              <th style="width: 100px">排序</th>
              <th style="width: 80px">启用</th>
              <th style="width: 130px">操作</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(row, idx) in rows" :key="row.rowKey">
              <td>
                <n-input v-model:value="row.field" size="small" placeholder="如 male" />
              </td>
              <td>
                <n-input v-model:value="row.label" size="small" placeholder="如 男" />
              </td>
              <td>
                <n-input v-model:value="row.value" size="small" placeholder="如 1" />
                <div v-if="rowError(row)" class="entry-error">{{ rowError(row) }}</div>
              </td>
              <td>
                <n-input-number v-model:value="row.sort" size="small" :min="0" style="width: 100%" />
              </td>
              <td>
                <n-switch v-model:value="row.status" size="small" />
              </td>
              <td>
                <n-space v-if="!readonly" :size="4">
                  <n-button
                    size="tiny"
                    secondary
                    circle
                    title="上移"
                    :disabled="idx === 0"
                    @click="moveRow(idx, -1)"
                  >
                    <template #icon>
                      <n-icon><ArrowUpOutlined /></n-icon>
                    </template>
                  </n-button>
                  <n-button
                    size="tiny"
                    secondary
                    circle
                    title="下移"
                    :disabled="idx === rows.length - 1"
                    @click="moveRow(idx, 1)"
                  >
                    <template #icon>
                      <n-icon><ArrowDownOutlined /></n-icon>
                    </template>
                  </n-button>
                  <n-button
                    size="tiny"
                    secondary
                    type="error"
                    circle
                    title="删除该选项"
                    @click="removeRow(idx)"
                  >
                    <template #icon>
                      <n-icon><DeleteOutlined /></n-icon>
                    </template>
                  </n-button>
                </n-space>
              </td>
            </tr>
            <tr v-if="!rows.length">
              <td colspan="6" class="entry-empty">还没有选项，点上方「添加一行」</td>
            </tr>
          </tbody>
        </n-table>

        <template #footer>
          <n-space justify="end">
            <n-button @click="showForm = false">取消</n-button>
            <n-button v-if="!readonly" type="primary" :loading="submitting" @click="submit">保存</n-button>
          </n-space>
        </template>
      </n-drawer-content>
    </n-drawer>
  </div>
</template>

<script lang="ts" setup>
  import { h, reactive, ref } from 'vue';
  import { NTag, useDialog, useMessage } from 'naive-ui';
  import type { DataTableColumns } from 'naive-ui';
  import { BasicTable, TableAction } from '@/components/Table';
  import {
    PlusOutlined,
    EyeOutlined,
    EditOutlined,
    CopyOutlined,
    DeleteOutlined,
    OrderedListOutlined,
    ArrowUpOutlined,
    ArrowDownOutlined,
  } from '@vicons/antd';
  import { addDict, deleteDict, editDict, getDict, getDictList } from '@/api/dict';
  import { fmtTime } from '@/api/common';
  import type { Dict } from '@/api/types';
  import { PAGE_SIZE_DICT, useDictStore } from '@/store/dict';
  import {
    checkDictValue,
    dictTypeMap,
    dictTypes,
    dictTimeHint,
    emptyRow,
    type EntryRow,
  } from './entry';

  const message = useMessage();
  const dialog = useDialog();
  const actionRef = ref();
  const formRef = ref();
  const submitting = ref(false);
  const showForm = ref(false);
  const editId = ref('');
  /** 查看模式：表单与选项只读 */
  const readonly = ref(false);
  const search = reactive({ q: '', type: null as string | null, status: 0 });

  /** 类型下拉：与 entry.ts 的 dictTypes 同源，前端过滤用 */
  const typeOptions = dictTypes.map((it) => ({ label: it.label, value: it.value }));
  const statusOptions = [
    { label: '全部状态', value: 0 },
    { label: '启用', value: 1 },
    { label: '禁用', value: 2 },
  ];

  const form = reactive({
    id: '',
    name: '',
    type: 'string',
    status: true,
    remark: '',
  });
  const rows = ref<EntryRow[]>([]);

  const rules = {
    id: { required: true, message: '请输入主键', trigger: 'blur' },
    name: { required: true, message: '请输入名称', trigger: 'blur' },
    type: { required: true, message: '请选择值类型', trigger: ['blur', 'change'] },
  };

  const columns: DataTableColumns<Dict> = [
    {
      title: '主键',
      key: 'id',
      width: 160,
      render: (row) =>
        h('span', { class: 'font-mono' }, row.id),
    },
    { title: '名称', key: 'name', width: 160 },
    {
      title: '值类型',
      key: 'type',
      width: 110,
      render: (row) => h(NTag, { size: 'small', type: 'info' }, { default: () => dictTypeMap[row.type] || row.type }),
    },
    {
      title: '选项数',
      key: 'entries',
      width: 90,
      render: (row) => (row.entries?.length ?? 0),
    },
    {
      title: '状态',
      key: 'status',
      width: 90,
      render: (row) =>
        h(
          NTag,
          { size: 'small', type: row.status ? 'success' : 'default' },
          { default: () => (row.status ? '启用' : '禁用') }
        ),
    },
    { title: '备注', key: 'remark', ellipsis: { tooltip: true } },
    {
      title: '更新时间',
      key: 'updated',
      width: 180,
      render: (row) => fmtTime(row.updated),
    },
  ];

  const actionColumn = reactive({
    width: 140,
    title: '操作',
    key: 'action',
    fixed: 'right',
    render(record: Dict) {
      // 默认三个：查看 / 编辑 / 更多，其余操作统一收进下拉
      const actions: any[] = [
        { label: '查看', icon: EyeOutlined, onClick: () => openEdit(record, true) },
        { label: '编辑', icon: EditOutlined, onClick: () => openEdit(record, false) },
      ];
      const dropDownActions: any[] = [
        { label: '复制', key: 'copy', icon: CopyOutlined },
        {
          // 内置字典禁止删除，避免拆掉代码里引用到的字典
          label: '删除',
          key: 'delete',
          icon: DeleteOutlined,
          ifShow: () => !record.inner,
        },
      ];
      const select = (key: string) => {
        if (key === 'copy') openCopy(record);
        if (key === 'delete') confirmDelete(record);
      };
      return h(TableAction, { style: 'button', actions, dropDownActions, select });
    },
  });

  const loadDataTable = async (res: any) => {
    const params: Recordable = { ...res };
    if (search.q) params.q = search.q;
    // 字典总量通常很小（十几条），且后端没有类型/状态过滤参数，
    // 这里在需要时一次拉全量，再按类型与状态在内存里筛，保证计数准确
    if (search.type || search.status > 0) params.size = 1000;
    const page = await getDictList(params);
    let list = page.list;
    if (search.type) list = list.filter((it) => it.type === search.type);
    if (search.status === 1) list = list.filter((it) => it.status);
    if (search.status === 2) list = list.filter((it) => !it.status);
    return { ...page, list, itemCount: list.length, pageCount: 1 };
  };

  function reload() {
    actionRef.value?.reload();
  }

  function resetSearch() {
    search.q = '';
    search.type = null;
    search.status = 0;
    reload();
  }

  function resetForm() {
    Object.assign(form, { id: '', name: '', type: 'string', status: true, remark: '' });
    rows.value = [];
  }

  function openAdd() {
    editId.value = '';
    readonly.value = false;
    resetForm();
    rows.value = [emptyRow(1)];
    showForm.value = true;
  }

  async function openEdit(row: Dict, viewOnly = false) {
    editId.value = row.id;
    readonly.value = viewOnly;
    const mod = await getDict(row.id);
    Object.assign(form, {
      id: mod.id,
      name: mod.name,
      type: mod.type,
      status: mod.status,
      remark: mod.remark,
    });
    rows.value = flattenEntries(mod.entries || []);
    showForm.value = true;
  }

  /** 后端返回的是树，字典当前实现为单层，这里拉平并保留 id */
  function flattenEntries(list: any[]): EntryRow[] {
    const out: EntryRow[] = [];
    const walk = (items: any[]) => {
      items.forEach((it) => {
        out.push({
          rowKey: it.id || `${it.field}-${out.length}`,
          id: it.id,
          field: it.field,
          label: it.label,
          value: it.value,
          sort: it.sort ?? out.length + 1,
          status: it.status ?? true,
        });
        if (it.children?.length) walk(it.children);
      });
    };
    walk(list);
    return out;
  }

  function addRow() {
    rows.value.push(emptyRow(rows.value.length + 1));
  }

  function removeRow(idx: number) {
    rows.value.splice(idx, 1);
  }

  /** 上移 / 下移一行 */
  function moveRow(idx: number, dir: -1 | 1) {
    const target = idx + dir;
    if (target < 0 || target >= rows.value.length) return;
    const list = [...rows.value];
    [list[idx], list[target]] = [list[target], list[idx]];
    rows.value = list;
  }

  /** 按当前行序把 sort 重写为 1..N，避免出现同序号导致展示顺序不确定 */
  function resortRows() {
    rows.value = rows.value.map((it, i) => ({ ...it, sort: i + 1 }));
    message.success('序号已按当前顺序重排');
  }

  function rowError(row: EntryRow): string {
    return checkDictValue(form.type, row.value);
  }

  async function submit() {
    try {
      await formRef.value?.validate();
    } catch {
      return;
    }
    // 前端先做一轮校验，减少来回请求；后端仍会二次校验
    const seen = new Set<string>();
    for (const [idx, row] of rows.value.entries()) {
      if (!row.field.trim()) {
        message.error(`第 ${idx + 1} 行字段名不能为空`);
        return;
      }
      if (seen.has(row.field.trim())) {
        message.error(`字段名 ${row.field} 重复`);
        return;
      }
      seen.add(row.field.trim());
      const err = checkDictValue(form.type, row.value);
      if (err) {
        message.error(`选项 ${row.field}：${err}`);
        return;
      }
    }

    const entries = rows.value.map((it) => ({
      id: it.id,
      field: it.field.trim(),
      label: it.label.trim() || it.value.trim(),
      value: it.value.trim(),
      sort: it.sort,
      status: it.status,
    }));

    submitting.value = true;
    try {
      if (editId.value) {
        await editDict(editId.value, {
          name: form.name,
          type: form.type,
          status: form.status,
          remark: form.remark,
          entries,
        });
        message.success('修改成功');
      } else {
        await addDict({
          id: form.id.trim(),
          name: form.name,
          type: form.type,
          status: form.status,
          remark: form.remark,
          entries,
        });
        message.success('新增成功');
      }
      // 分页条数字典被改过时立即刷新缓存，其它列表页的分页下拉随之更新
      if ((editId.value || form.id.trim()) === PAGE_SIZE_DICT) {
        await useDictStore().loadPageSize(true);
      }
      showForm.value = false;
      reload();
    } finally {
      submitting.value = false;
    }
  }

  async function remove(row: Dict) {
    await deleteDict(row.id);
    message.success('删除成功');
    reload();
  }

  /** 删除前确认；字典被代码或表单引用时误删代价较高 */
  function confirmDelete(row: Dict) {
    dialog.warning({
      title: '删除确认',
      content: `确定删除字典「${row.name}（${row.id}）」？其选项会一并删除，引用该字典的表单将取不到选项。`,
      positiveText: '删除',
      negativeText: '取消',
      onPositiveClick: () => remove(row),
    });
  }

  /**
   * 复制字典：把当前字典的名称/类型/选项带过来，主键留空由后端生成，
   * 名称自动加「副本」后缀，方便快速派生一个相似字典
   */
  async function openCopy(row: Dict) {
    const mod = await getDict(row.id);
    editId.value = '';
    readonly.value = false;
    Object.assign(form, {
      id: '',
      name: `${mod.name}副本`,
      type: mod.type,
      status: mod.status,
      remark: mod.remark,
    });
    rows.value = flattenEntries(mod.entries || []).map((it) => ({ ...it, id: undefined }));
    showForm.value = true;
  }
</script>

<style lang="less" scoped>
  .entry-error {
    margin-top: 2px;
    font-size: 12px;
    color: var(--error-color);
  }

  .entry-empty {
    padding: 16px;
    text-align: center;
    color: var(--text-color-3);
  }
</style>
