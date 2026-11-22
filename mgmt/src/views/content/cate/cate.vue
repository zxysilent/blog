<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="分类管理" />
    </div>

    <n-card :bordered="false" class="mt-4 proCard">
      <n-space class="mb-3" :size="12" align="center">
        <n-select
          v-model:value="filter.kind"
          :options="kindOptions"
          style="width: 130px"
          @update:value="loadTree"
        />
        <n-input
          v-model:value="filter.q"
          placeholder="栏目名称 / 路径"
          clearable
          style="width: 200px"
          @keyup.enter="loadTree"
        />
        <n-button type="primary" ghost @click="loadTree">查询</n-button>
        <n-button type="primary" @click="openAdd()">
          <template #icon>
            <n-icon><PlusOutlined /></n-icon>
          </template>
          新增分类
        </n-button>
        <n-button quaternary @click="toggleExpand">{{ expanded ? '收起全部' : '展开全部' }}</n-button>
      </n-space>

      <n-spin :show="loading">
        <n-data-table
          :columns="columns"
          :data="rows"
          :row-key="(row: Cate) => row.id"
          :pagination="false"
          :expanded-row-keys="expandedKeys"
          @update:expanded-row-keys="(keys: string[]) => (expandedKeys = keys)"
          size="small"
        />
      </n-spin>
    </n-card>

    <!-- 新增/编辑/查看弹窗 -->
    <n-modal
      v-model:show="showForm"
      preset="card"
      :title="readonly ? '查看分类' : editId ? '编辑分类' : '新增分类'"
      style="width: 640px"
    >
      <n-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-placement="left"
        label-width="90"
        :disabled="readonly"
      >
        <n-form-item label="上级分类">
          <n-tree-select
            v-model:value="form.pid"
            :options="parentOptions"
            placeholder="留空为顶级分类"
            clearable
            key-field="id"
            label-field="name"
            children-field="children"
          />
        </n-form-item>
        <n-form-item label="类型" path="kind">
          <n-select v-model:value="form.kind" :options="kindAllOptions" :disabled="!!editId" />
        </n-form-item>
        <n-form-item label="名称" path="name">
          <n-input v-model:value="form.name" placeholder="栏目名称，同级唯一" />
        </n-form-item>
        <n-form-item label="路径">
          <n-input v-model:value="form.path" placeholder="前台访问路径，留空自动生成" />
        </n-form-item>
        <n-form-item label="描述">
          <n-input v-model:value="form.intro" type="textarea" :rows="2" placeholder="可选" />
        </n-form-item>
        <n-grid :cols="2" :x-gap="12">
          <n-form-item-gi label="排序">
            <n-input-number v-model:value="form.sort" :min="0" style="width: 100%">
              <template #suffix>越小越前</template>
            </n-input-number>
          </n-form-item-gi>
          <n-form-item-gi label="颜色">
            <n-color-picker v-model:value="form.color" :show-alpha="false" />
          </n-form-item-gi>
        </n-grid>
        <n-form-item label="选项">
          <n-space>
            <n-checkbox v-model:checked="form.use">可被内容选中</n-checkbox>
            <n-checkbox v-model:checked="form.show">前台导航显示</n-checkbox>
            <n-checkbox v-model:checked="form.auth">需登录访问</n-checkbox>
          </n-space>
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showForm = false">{{ readonly ? '关闭' : '取消' }}</n-button>
          <n-button v-if="!readonly" type="primary" :loading="submitting" @click="submit">确定</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
  import { h, onMounted, reactive, ref } from 'vue';
  import { NButton, NTag, useDialog, useMessage } from 'naive-ui';
  import type { DataTableColumns } from 'naive-ui';
  import {
    PlusOutlined,
    EyeOutlined,
    EditOutlined,
    DeleteOutlined,
    CopyOutlined,
    ArrowUpOutlined,
    ArrowDownOutlined,
  } from '@vicons/antd';
  import { BasicTable } from '@/components/Table';
  import { TableAction } from '@/components/Table';
  import { addCate, deleteCate, editCate, getCateAll, getCateTree } from '@/api/cate';
  import type { Cate } from '@/api/types';

  const message = useMessage();
  const loading = ref(false);
  const submitting = ref(false);
  const moving = ref(false);
  const rows = ref<Cate[]>([]);
  /** 未过滤的完整树，同级排序调整要基于完整树计算 */
  const treeAll = ref<Cate[]>([]);
  const expandedKeys = ref<string[]>([]);
  const expanded = ref(false);
  const formRef = ref();
  const showForm = ref(false);
  const editId = ref('');
  /** 查看模式：表单与选项只读 */
  const readonly = ref(false);
  const dialog = useDialog();
  const parentOptions = ref<Cate[]>([]);

  /** 内容类型：与后端 model.Kind* 一致 */
  const kindMap: Record<number, string> = { 1: '文章', 2: '页面', 3: '笔记' };
  const kindAllOptions = [
    { label: '文章', value: 1 },
    { label: '页面', value: 2 },
    { label: '笔记', value: 3 },
  ];
  const kindOptions = [{ label: '全部类型', value: 0 }, ...kindAllOptions];

  const filter = reactive({ kind: 0, q: '' });

  const form = reactive({
    pid: null as string | null,
    kind: 1,
    name: '',
    path: '',
    intro: '',
    sort: 1000,
    color: '',
    use: true,
    show: true,
    auth: false,
  });

  const rules = {
    name: { required: true, message: '请输入栏目名称', trigger: 'blur' },
  };

  async function loadTree() {
    loading.value = true;
    try {
      // 后端树接口只支持 kind/use/show/pid 过滤，关键字在前端做，
      // 这样命中节点的祖先链能保留下来，层级关系不会丢
      const params: Recordable = {};
      if (filter.kind > 0) params.kind = filter.kind;
      // 树接口不带内容计数，列表接口带，两者合并后才能展示「内容数」
      const [tree, all] = await Promise.all([getCateTree(params), getCateAll()]);
      const counts = new Map(all.map((it) => [it.id, it.count ?? 0]));
      applyCount(tree, counts);
      treeAll.value = tree;
      parentOptions.value = tree;
      const kw = filter.q.trim().toLowerCase();
      rows.value = kw ? filterTree(tree, kw) : tree;
      // 默认展开顶层，方便一眼看完整结构
      if (!expandedKeys.value.length) {
        expandedKeys.value = tree.map((it) => it.id);
      }
    } finally {
      loading.value = false;
    }
  }

  /** 展开 / 收起全部：把所有含子级的节点主键写入/清空 */
  function toggleExpand() {
    if (expanded.value) {
      expandedKeys.value = [];
    } else {
      const keys: string[] = [];
      const walk = (nodes: Cate[]) => {
        nodes.forEach((it) => {
          if (it.children?.length) {
            keys.push(it.id);
            walk(it.children);
          }
        });
      };
      walk(rows.value);
      expandedKeys.value = keys;
    }
    expanded.value = !expanded.value;
  }

  /** 取某节点的全部同级（pid 为空表示顶层） */
  function siblingsOf(pid: string): Cate[] {
    if (!pid) return treeAll.value;
    const hit = findNode(treeAll.value, pid);
    return hit?.children || [];
  }

  function findNode(nodes: Cate[], id: string): Cate | undefined {
    for (const it of nodes) {
      if (it.id === id) return it;
      const sub = findNode(it.children || [], id);
      if (sub) return sub;
    }
    return undefined;
  }

  /**
   * 上移 / 下移一位
   * sort 字段默认都是 1000，直接交换不会生效，所以交换后按新顺序重写整组 sort
   */
  async function move(row: Cate, dir: -1 | 1) {
    const siblings = [...siblingsOf(row.pid || '')];
    const idx = siblings.findIndex((it) => it.id === row.id);
    const target = idx + dir;
    if (idx < 0 || target < 0 || target >= siblings.length) return;
    [siblings[idx], siblings[target]] = [siblings[target], siblings[idx]];
    moving.value = true;
    try {
      await Promise.all(siblings.map((it, i) => editCate(it.id, { sort: (i + 1) * 10 })));
      message.success('顺序已调整');
      await loadTree();
    } finally {
      moving.value = false;
    }
  }

  async function copyPath(row: Cate) {
    try {
      await navigator.clipboard.writeText(row.path || '');
      message.success('路径已复制');
    } catch {
      message.warning('浏览器不允许复制，请手动选择');
    }
  }

  /** 把内容计数写回到树节点上 */
  function applyCount(nodes: Cate[], counts: Map<string, number>) {
    nodes.forEach((it) => {
      it.count = counts.get(it.id) ?? 0;
      if (it.children?.length) applyCount(it.children, counts);
    });
  }

  /** 关键字过滤：自身或后代命中的节点保留，其余剪掉 */
  function filterTree(nodes: Cate[], kw: string): Cate[] {
    const out: Cate[] = [];
    nodes.forEach((it) => {
      const children = it.children?.length ? filterTree(it.children, kw) : [];
      const hit =
        it.name.toLowerCase().includes(kw) || (it.path || '').toLowerCase().includes(kw);
      if (hit || children.length) {
        out.push({ ...it, children: children.length ? children : undefined });
      }
    });
    return out;
  }

  const columns: DataTableColumns<Cate> = [
    { title: '名称', key: 'name', width: 240 },
    {
      title: '路径',
      key: 'path',
      width: 200,
      render: (row) =>
        row.path
          ? h('div', { class: 'flex items-center gap-1' }, [
              h('span', row.path),
              h(
                NButton,
                { size: 'tiny', quaternary: true, onClick: () => copyPath(row) },
                { default: () => '复制' }
              ),
            ])
          : '-',
    },
    {
      title: '类型',
      key: 'kind',
      width: 90,
      render: (row) => h(NTag, { size: 'small', type: 'info' }, { default: () => kindMap[row.kind] || '-' }),
    },
    { title: '排序', key: 'sort', width: 80 },
    {
      title: '内容数',
      key: 'count',
      width: 90,
      render: (row) => row.count ?? 0,
    },
    {
      title: '前台',
      key: 'show',
      width: 90,
      render: (row) =>
        h(
          NTag,
          { size: 'small', type: row.show ? 'success' : 'default' },
          { default: () => (row.show ? '显示' : '隐藏') }
        ),
    },
    {
      title: '可选中',
      key: 'use',
      width: 90,
      render: (row) =>
        h(
          NTag,
          { size: 'small', type: row.use ? 'success' : 'default' },
          { default: () => (row.use ? '是' : '否') }
        ),
    },
    {
      title: '操作',
      key: 'actions',
      width: 140,
      fixed: 'right',
      // 默认三个：查看 / 编辑 / 更多，其余操作统一收进下拉
      render: (row) =>
        h(TableAction, {
          style: 'button',
          actions: [
            { label: '查看', icon: EyeOutlined, onClick: () => openView(row) },
            { label: '编辑', icon: EditOutlined, onClick: () => openEdit(row) },
          ],
          dropDownActions: [
            { label: '加子级', key: 'add', icon: PlusOutlined },
            { label: '上移', key: 'up', icon: ArrowUpOutlined, ifShow: () => !moving.value },
            { label: '下移', key: 'down', icon: ArrowDownOutlined, ifShow: () => !moving.value },
            { label: '复制路径', key: 'copy', icon: CopyOutlined },
            // 内置栏目不可删除，避免拆掉站点结构
            { label: '删除', key: 'delete', icon: DeleteOutlined, type: 'error', ifShow: () => !row.inner },
          ],
          select: (key: string) => {
            if (key === 'add') openAdd(row.id);
            if (key === 'up') move(row, -1);
            if (key === 'down') move(row, 1);
            if (key === 'copy') copyPath(row);
            if (key === 'delete') confirmRemove(row);
          },
        }),
    },
  ];

  function resetForm() {
    Object.assign(form, {
      pid: null,
      kind: filter.kind > 0 ? filter.kind : 1,
      name: '',
      path: '',
      intro: '',
      sort: 1000,
      color: '',
      use: true,
      show: true,
      auth: false,
    });
  }

  function openAdd(pid?: string) {
    editId.value = '';
    readonly.value = false;
    resetForm();
    form.pid = pid || null;
    showForm.value = true;
  }

  function fillForm(row: Cate) {
    Object.assign(form, {
      pid: row.pid || null,
      kind: row.kind,
      name: row.name,
      path: row.path,
      intro: row.intro,
      sort: row.sort,
      color: row.color,
      use: row.use,
      show: row.show,
      auth: row.auth,
    });
  }

  /** 只读查看：复用同一个弹窗，表单禁用，底部不出现确定按钮 */
  function openView(row: Cate) {
    editId.value = row.id;
    readonly.value = true;
    fillForm(row);
    showForm.value = true;
  }

  function openEdit(row: Cate) {
    editId.value = row.id;
    readonly.value = false;
    fillForm(row);
    showForm.value = true;
  }

  async function submit() {
    try {
      await formRef.value?.validate();
    } catch {
      return;
    }
    const payload = {
      ...form,
      pid: form.pid || '',
    };
    submitting.value = true;
    try {
      if (editId.value) {
        await editCate(editId.value, payload);
        message.success('修改成功');
      } else {
        await addCate(payload);
        message.success('新增成功');
      }
      showForm.value = false;
      await loadTree();
    } finally {
      submitting.value = false;
    }
  }

  async function remove(row: Cate) {
    await deleteCate(row.id);
    message.success('删除成功');
    await loadTree();
  }

  /** 删除确认：子栏目会提升为顶级，属于结构性改动 */
  function confirmRemove(row: Cate) {
    dialog.warning({
      title: '删除确认',
      content: `确定删除栏目「${row.name}」？子栏目会提升为顶级，删除后不可恢复。`,
      positiveText: '删除',
      negativeText: '取消',
      onPositiveClick: () => remove(row),
    });
  }

  onMounted(loadTree);
</script>
