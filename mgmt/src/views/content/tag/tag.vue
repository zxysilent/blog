<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="标签管理" />
    </div>

    <n-card :bordered="false" class="mt-4 proCard">
      <n-space class="mb-3" :size="12" align="center">
        <n-input
          v-model:value="keyword"
          placeholder="标签名 / 描述"
          clearable
          style="width: 220px"
          @keyup.enter="reload"
        />
        <n-button type="primary" ghost @click="reload">查询</n-button>
        <n-button type="primary" @click="openAdd">
          <template #icon>
            <n-icon><PlusOutlined /></n-icon>
          </template>
          新增标签
        </n-button>
      </n-space>

      <n-space v-if="checkedKeys.length" class="mb-3" :size="8" align="center">
        <n-tag type="info" :bordered="false">已选 {{ checkedKeys.length }} 项</n-tag>
        <n-button size="small" type="error" ghost :loading="batching" @click="confirmBatchDelete()">
          批量删除
        </n-button>
        <n-button size="small" quaternary @click="checkedKeys = []">清空选择</n-button>
      </n-space>

      <BasicTable
        v-model:checked-row-keys="checkedKeys"
        :columns="columns"
        :request="loadDataTable"
        :row-key="(row: Tag) => row.id"
        ref="actionRef"
        :action-column="actionColumn"
      />
    </n-card>

    <n-modal
      v-model:show="showForm"
      preset="card"
      :title="readonly ? '查看标签' : editId ? '编辑标签' : '新增标签'"
      style="width: 480px"
    >
      <n-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-placement="left"
        label-width="80"
        :disabled="readonly"
      >
        <n-form-item label="名称" path="name">
          <n-input v-model:value="form.name" placeholder="标签名，全局唯一" />
        </n-form-item>
        <n-form-item label="描述">
          <n-input v-model:value="form.intro" type="textarea" :rows="2" placeholder="可选" />
        </n-form-item>
        <n-form-item label="颜色">
          <n-color-picker v-model:value="form.color" :show-alpha="false" />
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
  import { h, reactive, ref } from 'vue';
  import { NTag, useDialog, useMessage } from 'naive-ui';
  import { BasicTable, TableAction } from '@/components/Table';
  import { PlusOutlined, EyeOutlined, EditOutlined, DeleteOutlined } from '@vicons/antd';
  import { addTag, deleteTag, editTag, getTagList } from '@/api/tag';
  import { fmtTime } from '@/api/common';
  import type { Tag } from '@/api/types';

  const message = useMessage();
  const dialog = useDialog();
  const actionRef = ref();
  const formRef = ref();
  const submitting = ref(false);
  const batching = ref(false);
  const showForm = ref(false);
  const editId = ref('');
  /** 查看模式：表单只读 */
  const readonly = ref(false);
  const keyword = ref('');
  const checkedKeys = ref<string[]>([]);

  const form = reactive({ name: '', intro: '', color: '' });
  const rules = { name: { required: true, message: '请输入标签名', trigger: 'blur' } };

  const columns = [
    { type: 'selection' },
    {
      title: '名称',
      key: 'name',
      width: 200,
      render: (row: Tag) =>
        row.color
          ? h(NTag, { size: 'small', color: { color: row.color, textColor: '#fff' } }, { default: () => row.name })
          : h(NTag, { size: 'small' }, { default: () => row.name }),
    },
    { title: '描述', key: 'intro', ellipsis: { tooltip: true } },
    {
      title: '创建时间',
      key: 'created',
      width: 180,
      render: (row: Tag) => fmtTime(row.created),
    },
  ];

  const actionColumn = reactive({
    width: 140,
    title: '操作',
    key: 'action',
    fixed: 'right',
    render(record: Tag) {
      // 默认三个：查看 / 编辑 / 更多，其余操作统一收进下拉
      const actions: any[] = [
        { label: '查看', icon: EyeOutlined, onClick: () => openView(record) },
        { label: '编辑', icon: EditOutlined, onClick: () => openEdit(record) },
      ];
      const dropDownActions: any[] = [
        { label: '删除', key: 'delete', icon: DeleteOutlined, type: 'error' },
      ];
      const select = (key: string) => {
        if (key === 'delete') confirmBatchDelete([record.id]);
      };
      return h(TableAction, { style: 'button', actions, dropDownActions, select });
    },
  });

  const loadDataTable = async (res: any) => {
    const params: Recordable = { ...res };
    if (keyword.value) params.q = keyword.value;
    return await getTagList(params);
  };

  function reload() {
    checkedKeys.value = [];
    actionRef.value?.reload();
  }

  function openAdd() {
    editId.value = '';
    readonly.value = false;
    Object.assign(form, { name: '', intro: '', color: '' });
    showForm.value = true;
  }

  /** 只读查看：复用同一个弹窗，表单禁用，底部不出现确定按钮 */
  function openView(row: Tag) {
    editId.value = row.id;
    readonly.value = true;
    Object.assign(form, { name: row.name, intro: row.intro, color: row.color });
    showForm.value = true;
  }

  function openEdit(row: Tag) {
    editId.value = row.id;
    readonly.value = false;
    Object.assign(form, { name: row.name, intro: row.intro, color: row.color });
    showForm.value = true;
  }

  async function submit() {
    try {
      await formRef.value?.validate();
    } catch {
      return;
    }
    submitting.value = true;
    try {
      if (editId.value) {
        await editTag(editId.value, { ...form });
        message.success('修改成功');
      } else {
        await addTag({ ...form });
        message.success('新增成功');
      }
      showForm.value = false;
      actionRef.value?.reload();
    } finally {
      submitting.value = false;
    }
  }

  async function remove(ids: string[]) {
    batching.value = true;
    try {
      const res = await Promise.allSettled(ids.map((id) => deleteTag(id)));
      const failed = res.filter((it) => it.status === 'rejected').length;
      if (failed) message.warning(`${failed} 条删除失败`);
      else message.success('删除成功');
      checkedKeys.value = [];
      actionRef.value?.reload();
    } finally {
      batching.value = false;
    }
  }

  /** 删除前确认；单条与批量走同一个入口 */
  function confirmBatchDelete(ids: string[] = [...checkedKeys.value]) {
    if (!ids.length) return;
    const tip = ids.length > 1 ? `确定删除选中的 ${ids.length} 个标签？` : '确定删除该标签？';
    dialog.warning({
      title: '删除确认',
      content: `${tip}删除后内容与标签的关联会一并解除。`,
      positiveText: '删除',
      negativeText: '取消',
      onPositiveClick: () => remove([...ids]),
    });
  }
</script>
