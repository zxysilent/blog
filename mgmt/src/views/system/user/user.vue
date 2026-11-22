<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="用户管理" />
    </div>

    <!-- 概览：点卡片切筛选 -->
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
        <n-input
          v-model:value="filter.q"
          placeholder="账号 / 名称 / 手机 / 邮箱"
          clearable
          style="width: 220px"
          @keyup.enter="reload"
        />
        <n-select
          v-model:value="filter.role_id"
          :options="roleFilterOptions"
          placeholder="全部角色"
          clearable
          style="width: 180px"
          @update:value="reload"
        />
        <n-select
          v-model:value="filter.status"
          :options="statusOptions"
          style="width: 130px"
          @update:value="reload"
        />
        <n-button type="primary" ghost @click="reload">查询</n-button>
        <n-button quaternary @click="resetFilter">重置</n-button>
        <n-button type="primary" @click="openAdd">
          <template #icon>
            <n-icon><PlusOutlined /></n-icon>
          </template>
          新增用户
        </n-button>
      </n-space>

      <!-- 批量操作 -->
      <n-space v-if="checkedKeys.length" class="mb-3" :size="8" align="center">
        <n-tag type="info" :bordered="false">已选 {{ checkedKeys.length }} 项</n-tag>
        <n-button size="small" type="success" ghost :loading="batching" @click="batchStatus(1)">
          批量启用
        </n-button>
        <n-button size="small" ghost :loading="batching" @click="batchStatus(2)">批量禁用</n-button>
        <n-button size="small" type="error" ghost :loading="batching" @click="confirmDelete(checkedKeys)">
          批量删除
        </n-button>
        <n-button size="small" quaternary @click="checkedKeys = []">清空选择</n-button>
      </n-space>

      <BasicTable
        v-model:checked-row-keys="checkedKeys"
        :columns="columns"
        :request="loadDataTable"
        :row-key="(row: User) => row.id"
        ref="actionRef"
        :action-column="actionColumn"
        :scroll-x="1400"
      />
    </n-card>

    <!-- 新增/编辑/查看 -->
    <n-drawer v-model:show="showForm" :width="620" placement="right">
      <n-drawer-content :title="drawerTitle" closable>
        <n-form
          ref="formRef"
          :model="form"
          :rules="rules"
          label-placement="left"
          label-width="90"
          :disabled="readonly"
        >
          <n-form-item label="账号" path="num">
            <n-input v-model:value="form.num" placeholder="登录账号，创建后不可修改" :disabled="!!editId" />
          </n-form-item>
          <n-form-item label="名称" path="name">
            <n-input v-model:value="form.name" placeholder="显示名称" />
          </n-form-item>
          <n-form-item label="头像">
            <n-space vertical class="w-full">
              <n-input v-model:value="form.avatar" placeholder="图片地址，可留空" clearable />
              <n-space align="center" :size="12">
                <n-upload :show-file-list="false" accept="image/*" :custom-request="onUploadAvatar">
                  <n-button size="small" :loading="uploading">
                    <template #icon>
                      <n-icon><CloudUploadOutlined /></n-icon>
                    </template>
                    上传头像
                  </n-button>
                </n-upload>
                <n-avatar v-if="form.avatar" round :size="40" :src="form.avatar" />
                <span class="text-sm text-gray-400">留空时列表按名称首字展示</span>
              </n-space>
            </n-space>
          </n-form-item>
          <n-form-item :label="editId ? '新密码' : '密码'" path="passwd">
            <n-input
              v-model:value="form.passwd"
              type="password"
              show-password-on="click"
              :placeholder="editId ? '留空表示不修改' : '至少 6 位'"
            />
          </n-form-item>
          <n-form-item label="角色" path="role_id">
            <n-select
              v-model:value="form.role_id"
              :options="roleOptions"
              placeholder="请选择角色（只能分配权限低于自己的角色）"
            />
          </n-form-item>
          <n-form-item label="手机">
            <n-input v-model:value="form.phone" placeholder="可选" />
          </n-form-item>
          <n-form-item label="邮箱">
            <n-input v-model:value="form.email" placeholder="可选" />
          </n-form-item>
          <n-form-item label="状态">
            <n-radio-group v-model:value="form.status">
              <n-radio-button :value="1">正常</n-radio-button>
              <n-radio-button :value="2">禁用</n-radio-button>
            </n-radio-group>
          </n-form-item>
        </n-form>
        <template #footer>
          <n-space justify="end">
            <n-button @click="showForm = false">{{ readonly ? '关闭' : '取消' }}</n-button>
            <n-button v-if="!readonly" type="primary" :loading="submitting" @click="submit">
              {{ editId ? '保存修改' : '确认新增' }}
            </n-button>
          </n-space>
        </template>
      </n-drawer-content>
    </n-drawer>

    <!-- 重置密码 -->
    <n-modal v-model:show="showPwd" preset="card" title="重置密码" style="width: 480px">
      <n-alert type="info" :bordered="false" class="mb-3">
        将为「{{ pwdUser?.name || pwdUser?.num }}」设置新密码，旧密码立即失效。
      </n-alert>
      <n-space vertical :size="12">
        <n-input-group>
          <n-input v-model:value="pwdValue" placeholder="至少 6 位" />
          <n-button quaternary @click="genPwd">随机生成</n-button>
        </n-input-group>
        <n-space align="center" :size="8">
          <n-tag v-if="pwdValue" :bordered="false">{{ pwdValue }}</n-tag>
          <n-button size="tiny" quaternary :disabled="!pwdValue" @click="copyPwd">复制</n-button>
        </n-space>
      </n-space>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showPwd = false">取消</n-button>
          <n-button type="primary" :loading="submitting" :disabled="!pwdValue" @click="submitPwd">
            确认重置
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
  import { computed, h, onMounted, reactive, ref } from 'vue';
  import { NAvatar, NTag, useDialog, useMessage } from 'naive-ui';
  import { BasicTable, TableAction } from '@/components/Table';
  import {
    PlusOutlined,
    CloudUploadOutlined,
    EyeOutlined,
    EditOutlined,
    KeyOutlined,
    CheckCircleOutlined,
    StopOutlined,
    DeleteOutlined,
  } from '@vicons/antd';
  import { addUser, deleteUser, editUser, getUserList } from '@/api/user';
  import { getRoleList } from '@/api/role';
  import { uploadImage } from '@/api/site';
  import { fmtTime } from '@/api/common';
  import type { Role, User } from '@/api/types';

  /** 后端常量：用户状态 1 正常 2 禁用 */
  const STATUS_NORMAL = 1;
  const STATUS_DISABLE = 2;
  /** 根角色主键，仅用于给角色标签上色，不参与任何展示逻辑 */
  const ROOT_ROLE_ID = '00000000000000000000000000';

  const message = useMessage();
  const dialog = useDialog();
  const actionRef = ref();
  const formRef = ref();
  const submitting = ref(false);
  const batching = ref(false);
  const uploading = ref(false);

  const filter = reactive({ q: '', role_id: null as string | null, status: 0 });
  const statusOptions = [
    { label: '全部状态', value: 0 },
    { label: '正常', value: 1 },
    { label: '禁用', value: 2 },
  ];

  const checkedKeys = ref<string[]>([]);
  const showForm = ref(false);
  const editId = ref('');
  /** 查看模式：表单与选项只读 */
  const readonly = ref(false);
  const drawerTitle = computed(() => {
    if (readonly.value) return '查看用户';
    return editId.value ? '编辑用户' : '新增用户';
  });
  const allRoles = ref<Role[]>([]);

  const form = reactive({
    num: '',
    name: '',
    avatar: '',
    passwd: '',
    role_id: null as string | null,
    phone: '',
    email: '',
    status: STATUS_NORMAL,
  });

  const rules = {
    num: { required: true, message: '请输入账号', trigger: 'blur' },
    name: { required: true, message: '请输入名称', trigger: 'blur' },
    role_id: { required: true, message: '请选择角色', trigger: ['blur', 'change'] },
    passwd: {
      trigger: 'blur',
      validator: (_rule: any, value: string) => {
        // 新增必填；编辑留空表示不修改
        if (!editId.value && !value) return new Error('请输入密码');
        if (value && value.length < 6) return new Error('密码至少 6 位');
        return true;
      },
    },
  };

  /** 角色下拉：只列权限低于自己（层级数值更大）的角色，后端已按此过滤 */
  const roleOptions = computed(() =>
    allRoles.value.map((it) => ({ label: it.name, value: it.id }))
  );
  const roleFilterOptions = computed(() =>
    allRoles.value.map((it) => ({ label: it.name, value: it.id }))
  );

  const stats = reactive({ total: 0, normal: 0, disabled: 0, roles: 0 });
  const statCards = computed(() => [
    {
      key: 'total',
      label: '用户总数',
      count: stats.total,
      active: filter.status === 0 && !filter.role_id,
      pick: () => {
        filter.status = 0;
        filter.role_id = null;
        reload();
      },
    },
    {
      key: 'normal',
      label: '正常',
      count: stats.normal,
      active: filter.status === 1,
      pick: () => {
        filter.status = 1;
        reload();
      },
    },
    {
      key: 'disabled',
      label: '禁用',
      count: stats.disabled,
      active: filter.status === 2,
      pick: () => {
        filter.status = 2;
        reload();
      },
    },
    {
      key: 'roles',
      label: '可分配角色',
      count: stats.roles,
      active: !!filter.role_id,
      pick: () => {
        filter.role_id = null;
        reload();
      },
    },
  ]);

  const columns = [
    { type: 'selection' },
    {
      title: '账号',
      key: 'num',
      width: 200,
      render: (row: User) =>
        h('div', { class: 'flex items-center gap-2' }, [
          h(
            NAvatar,
            { round: true, size: 28, src: row.avatar || undefined },
            { default: () => (row.name || row.num || '?').slice(0, 1) }
          ),
          h('span', row.num),
        ]),
    },
    { title: '名称', key: 'name', width: 130 },
    {
      title: '角色',
      key: 'role',
      width: 150,
      render: (row: User) =>
        row.role
          ? h(
              NTag,
              {
                size: 'small',
                type: row.role.id === ROOT_ROLE_ID ? 'error' : 'info',
                bordered: false,
              },
              { default: () => row.role?.name }
            )
          : '-',
    },
    { title: '手机', key: 'phone', width: 130, render: (row: User) => row.phone || '-' },
    { title: '邮箱', key: 'email', width: 200, render: (row: User) => row.email || '-' },
    {
      title: '状态',
      key: 'status',
      width: 100,
      render: (row: User) =>
        h(
          NTag,
          { size: 'small', type: row.status === STATUS_NORMAL ? 'success' : 'error' },
          { default: () => (row.status === STATUS_NORMAL ? '正常' : '禁用') }
        ),
    },
    {
      title: '创建时间',
      key: 'created',
      width: 170,
      render: (row: User) => fmtTime(row.created),
    },
  ];

  const actionColumn = reactive({
    width: 140,
    title: '操作',
    key: 'action',
    fixed: 'right',
    render(record: User) {
      return h(TableAction, {
        style: 'button',
        // 默认三个：查看 / 编辑 / 更多，其余操作统一收进下拉
        actions: [
          { label: '查看', icon: EyeOutlined, onClick: () => openView(record) },
          { label: '编辑', icon: EditOutlined, onClick: () => openEdit(record) },
        ],
        dropDownActions: [
          { label: '重置密码', key: 'pwd', icon: KeyOutlined },
          {
            label: record.status === STATUS_NORMAL ? '禁用' : '启用',
            key: 'status',
            icon: record.status === STATUS_NORMAL ? StopOutlined : CheckCircleOutlined,
          },
          { label: '删除', key: 'delete', icon: DeleteOutlined, type: 'error' },
        ],
        select: (key: string) => {
          if (key === 'pwd') openPwd(record);
          if (key === 'status') toggleStatus(record);
          if (key === 'delete') confirmDelete([record.id]);
        },
      });
    },
  });

  /** 界面筛选翻译成后端参数 */
  function buildParams(extra: Recordable = {}): Recordable {
    const params: Recordable = { ...extra };
    if (filter.q) params.q = filter.q;
    if (filter.role_id) params.role_id = filter.role_id;
    if (filter.status > 0) params.status = filter.status;
    return params;
  }

  const loadDataTable = async (res: any) => {
    const page = await getUserList(buildParams({ ...res }));
    refreshStats();
    return page;
  };

  /** 概览只按角色维度统计，不受状态/关键字影响，避免点卡片时数字跳动 */
  async function refreshStats() {
    const base: Recordable = { page: 1, size: 1 };
    if (filter.q) base.q = filter.q;
    if (filter.role_id) base.role_id = filter.role_id;
    const [all, normal, disabled] = await Promise.all([
      getUserList({ ...base }),
      getUserList({ ...base, status: 1 }),
      getUserList({ ...base, status: 2 }),
    ]);
    stats.total = all.itemCount;
    stats.normal = normal.itemCount;
    stats.disabled = disabled.itemCount;
    stats.roles = allRoles.value.length;
  }

  function reload() {
    checkedKeys.value = [];
    actionRef.value?.reload();
  }

  function resetFilter() {
    Object.assign(filter, { q: '', role_id: null, status: 0 });
    reload();
  }

  async function loadRoles() {
    const page = await getRoleList({ page: 1, size: 1000 });
    allRoles.value = page.list;
  }

  function openAdd() {
    editId.value = '';
    readonly.value = false;
    Object.assign(form, {
      num: '',
      name: '',
      avatar: '',
      passwd: '',
      // 默认选层级最低的角色，避免误建高权限账号
      role_id: allRoles.value.length ? allRoles.value[allRoles.value.length - 1].id : null,
      phone: '',
      email: '',
      status: STATUS_NORMAL,
    });
    showForm.value = true;
  }

  function fillForm(row: User) {
    Object.assign(form, {
      num: row.num,
      name: row.name,
      avatar: row.avatar || '',
      passwd: '',
      role_id: row.role_id,
      phone: row.phone || '',
      email: row.email || '',
      status: row.status,
    });
  }

  /** 只读查看：复用同一个抽屉，表单禁用，底部不出现保存按钮 */
  function openView(row: User) {
    editId.value = row.id;
    readonly.value = true;
    fillForm(row);
    showForm.value = true;
  }

  function openEdit(row: User) {
    editId.value = row.id;
    readonly.value = false;
    fillForm(row);
    showForm.value = true;
  }

  async function onUploadAvatar({ file, onFinish, onError }: any) {
    uploading.value = true;
    try {
      const res: any = await uploadImage(file.file as File);
      const url = typeof res === 'string' ? res : res?.data || '';
      if (!url) throw new Error('上传失败');
      form.avatar = url;
      message.success('上传成功');
      onFinish();
    } catch (err: any) {
      message.error(err?.message || '上传失败');
      onError();
    } finally {
      uploading.value = false;
    }
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
        // 账号不可改，密码留空表示不修改
        const payload: Recordable = {
          name: form.name,
          avatar: form.avatar,
          role_id: form.role_id,
          phone: form.phone,
          email: form.email,
          status: form.status,
        };
        if (form.passwd) payload.passwd = form.passwd;
        await editUser(editId.value, payload);
        message.success('修改成功');
      } else {
        await addUser({
          num: form.num,
          name: form.name,
          avatar: form.avatar,
          passwd: form.passwd,
          role_id: form.role_id,
          phone: form.phone,
          email: form.email,
          status: form.status,
        });
        message.success('新增成功');
      }
      showForm.value = false;
      reload();
    } finally {
      submitting.value = false;
    }
  }

  // ---------------- 重置密码 ----------------
  const showPwd = ref(false);
  const pwdUser = ref<User | null>(null);
  const pwdValue = ref('');

  function openPwd(row: User) {
    pwdUser.value = row;
    pwdValue.value = '';
    showPwd.value = true;
  }

  function genPwd() {
    // 去掉易混淆字符，方便口头/手抄转达
    const chars = 'abcdefghijkmnpqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ23456789';
    let out = '';
    for (let i = 0; i < 12; i++) out += chars[Math.floor(Math.random() * chars.length)];
    pwdValue.value = out;
  }

  async function copyPwd() {
    try {
      await navigator.clipboard.writeText(pwdValue.value);
      message.success('已复制');
    } catch {
      message.warning('浏览器不允许复制，请手动选择');
    }
  }

  async function submitPwd() {
    if (!pwdUser.value || !pwdValue.value) return;
    submitting.value = true;
    try {
      await editUser(pwdUser.value.id, { passwd: pwdValue.value });
      message.success('密码已重置');
      showPwd.value = false;
    } finally {
      submitting.value = false;
    }
  }

  // ---------------- 状态与删除 ----------------
  async function toggleStatus(row: User) {
    const next = row.status === STATUS_NORMAL ? STATUS_DISABLE : STATUS_NORMAL;
    await editUser(row.id, { status: next });
    message.success(next === STATUS_NORMAL ? '已启用' : '已禁用');
    reload();
  }

  function confirmDelete(ids: string[]) {
    const tip = ids.length > 1 ? `确定删除选中的 ${ids.length} 个用户？` : '确定删除该用户？';
    dialog.warning({
      title: '删除确认',
      content: `${tip}删除后该账号立即失去登录能力，且不可恢复。`,
      positiveText: '删除',
      negativeText: '取消',
      onPositiveClick: async () => {
        await runBatch([...ids], (id) => deleteUser(id));
        message.success('删除成功');
        reload();
      },
    });
  }

  async function batchStatus(status: number) {
    await runBatch([...checkedKeys.value], (id) => editUser(id, { status }));
    message.success(status === STATUS_NORMAL ? '批量启用完成' : '批量禁用完成');
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
    await loadRoles().catch(() => undefined);
    refreshStats();
  });
</script>

<style lang="less" scoped>
  .stat-active {
    box-shadow: 0 0 0 1px var(--primary-color);
  }
</style>
