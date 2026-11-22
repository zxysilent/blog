<template>
  <div>
    <div class="n-layout-page-header">
      <n-card :bordered="false" title="角色权限" />
    </div>

    <n-card :bordered="false" class="mt-4 proCard">
      <n-space class="mb-3" :size="12" align="center">
        <n-input
          v-model:value="keyword"
          placeholder="角色名称 / 说明"
          clearable
          style="width: 220px"
          @keyup.enter="reload"
        />
        <n-button type="primary" ghost @click="reload">查询</n-button>
        <n-button quaternary @click="resetFilter">重置</n-button>
        <n-button type="primary" @click="openAdd">
          <template #icon>
            <n-icon><PlusOutlined /></n-icon>
          </template>
          新增角色
        </n-button>
      </n-space>

      <BasicTable
        :columns="columns"
        :request="loadDataTable"
        :row-key="(row: Role) => row.id"
        ref="actionRef"
        :action-column="actionColumn"
        :scroll-x="1200"
      />
    </n-card>

    <!-- 新增/编辑 -->
    <n-modal
      v-model:show="showForm"
      preset="card"
      :title="editId ? '编辑角色' : '新增角色'"
      style="width: 560px"
    >
      <n-form ref="formRef" :model="form" :rules="rules" label-placement="left" label-width="90">
        <n-form-item label="角色名称" path="name">
          <n-input v-model:value="form.name" placeholder="请输入角色名称" />
        </n-form-item>
        <n-form-item label="说明">
          <n-input v-model:value="form.intro" type="textarea" :rows="2" placeholder="可选" />
        </n-form-item>
      </n-form>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showForm = false">取消</n-button>
          <n-button type="primary" :loading="submitting" @click="submit">
            {{ editId ? '保存修改' : '确认新增' }}
          </n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 权限分配 -->
    <n-drawer v-model:show="showAuth" :width="720" placement="right">
      <n-drawer-content
        closable
        :native-scrollbar="false"
        body-style="overflow: hidden;"
        body-content-style="display: flex; flex-direction: column; height: 100%; min-height: 0;"
      >
        <template #header>
          <div class="flex items-center gap-2">
            <span>分配权限</span>
            <n-tag v-if="authRole" size="small" :bordered="false" type="info">
              {{ authRole.name }}
            </n-tag>
          </div>
        </template>

        <n-spin :show="authLoading" class="auth-spin">
          <n-alert v-if="isRoot" type="warning" :bordered="false" class="mb-3">
            根角色天然拥有全部权限，无需分配，也无法取消。
          </n-alert>

          <div class="auth-bar">
            <n-input
              v-model:value="permKeyword"
              placeholder="搜索权限名称 / 标识"
              clearable
              style="width: 220px"
            >
              <template #prefix>
                <n-icon><SearchOutlined /></n-icon>
              </template>
            </n-input>
            <n-space :size="6" align="center">
              <n-button size="small" type="primary" secondary :disabled="isRoot" @click="checkAll">
                <template #icon>
                  <n-icon><CheckSquareOutlined /></n-icon>
                </template>
                全选
              </n-button>
              <n-button size="small" secondary :disabled="isRoot" @click="checkClear">
                <template #icon>
                  <n-icon><CloseCircleOutlined /></n-icon>
                </template>
                清空
              </n-button>
              <n-button size="small" secondary :disabled="isRoot" @click="checkInvert">
                <template #icon>
                  <n-icon><SwapOutlined /></n-icon>
                </template>
                反选
              </n-button>
              <n-button size="small" dashed @click="toggleExpand">
                <template #icon>
                  <n-icon>
                    <UpOutlined v-if="expanded" />
                    <DownOutlined v-else />
                  </n-icon>
                </template>
                {{ expanded ? '收起' : '展开' }}
              </n-button>
            </n-space>
          </div>

          <div class="auth-stats">
            <n-tag size="small" :bordered="false" type="success">
              已选 {{ checkedKeys.length }} / {{ leafTotal }} 项
            </n-tag>
            <n-tag
              v-for="it in groupStats"
              :key="it.key"
              size="small"
              :bordered="false"
              :type="it.checked ? 'success' : 'default'"
            >
              {{ it.group }} {{ it.checked }}/{{ it.total }}
            </n-tag>
          </div>

          <n-empty v-if="!treeView.length" description="没有匹配的权限标识" class="py-6" />
          <n-tree
            v-else
            block-line
            cascade
            checkable
            check-on-click
            class="auth-tree"
            :data="treeView"
            :checked-keys="checkedKeys"
            :expanded-keys="expandedKeys"
            :pattern="permKeyword"
            :show-irrelevant-nodes="false"
            @update:checked-keys="onChecked"
            @update:expanded-keys="(keys: string[]) => (expandedKeys = keys)"
          />
        </n-spin>

        <template #footer>
          <n-space justify="end">
            <n-button @click="showAuth = false">关闭</n-button>
            <n-button type="primary" :loading="submitting" :disabled="isRoot" @click="submitAuth">
              保存授权
            </n-button>
          </n-space>
        </template>
      </n-drawer-content>
    </n-drawer>

    <!-- 查看已授权 -->
    <n-modal
      v-model:show="showView"
      preset="card"
      :title="`已授权权限 - ${viewRole?.name || ''}`"
      style="width: 640px"
    >
      <n-spin :show="viewLoading">
        <n-empty v-if="!viewList.length" description="尚未授权任何权限" class="py-6" />
        <n-space v-else vertical :size="14">
          <div v-for="grp in viewGroups" :key="grp.group">
            <n-space :size="8" align="center" class="mb-1">
              <n-text strong>{{ grp.group }}</n-text>
              <n-tag size="tiny" :bordered="false">{{ grp.items.length }}</n-tag>
            </n-space>
            <n-space :size="6">
              <n-tag v-for="it in grp.items" :key="it.id" size="small" :bordered="false">
                {{ it.name }}
                <n-text depth="3" class="ml-1">{{ it.perm }}</n-text>
              </n-tag>
            </n-space>
          </div>
        </n-space>
      </n-spin>
      <template #footer>
        <n-space justify="end">
          <n-button @click="showView = false">关闭</n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script lang="ts" setup>
  import { computed, h, onMounted, reactive, ref } from 'vue';
  import { NTag, useDialog, useMessage } from 'naive-ui';
  import { BasicTable, TableAction } from '@/components/Table';
  import {
    PlusOutlined,
    EyeOutlined,
    EditOutlined,
    DeleteOutlined,
    SafetyCertificateOutlined,
    SearchOutlined,
    CheckSquareOutlined,
    CloseCircleOutlined,
    SwapOutlined,
    DownOutlined,
    UpOutlined,
  } from '@vicons/antd';
  import {
    addRole,
    deleteRole,
    editRole,
    getRoleList,
    getRolePermissions,
    setRolePermissions,
  } from '@/api/role';
  import { getPermissionTree } from '@/api/auth';
  import { fmtTime } from '@/api/common';
  import type { Permission, PermissionNode, Role } from '@/api/types';

  const ROOT_ID = '00000000000000000000000000';

  const message = useMessage();
  const dialog = useDialog();
  const actionRef = ref();
  const formRef = ref();
  const submitting = ref(false);
  const keyword = ref('');

  const showForm = ref(false);
  const editId = ref('');
  const form = reactive({ name: '', intro: '' });
  const rules = {
    name: { required: true, message: '请输入角色名称', trigger: 'blur' },
  };

  // ---------------- 列表 ----------------
  const columns = [
    { title: '角色名称', key: 'name', width: 170 },
    { title: '说明', key: 'intro', minWidth: 200, render: (row: Role) => row.intro || '-' },
    {
      title: '用户数',
      key: 'user_count',
      width: 100,
      render: (row: Role) => row.user_count ?? 0,
    },
    {
      title: '权限数',
      key: 'perm_count',
      width: 100,
      render: (row: Role) =>
        row.id === ROOT_ID
          ? h(NTag, { size: 'small', type: 'error', bordered: false }, { default: () => '全部' })
          : (row.perm_count ?? 0),
    },
    {
      title: '来源',
      key: 'inner',
      width: 100,
      render: (row: Role) =>
        h(
          NTag,
          { type: row.inner ? 'warning' : 'default', size: 'small' },
          { default: () => (row.inner ? '内置' : '自定义') }
        ),
    },
    {
      title: '创建时间',
      key: 'created',
      width: 170,
      render: (row: Role) => fmtTime(row.created),
    },
  ];

  const actionColumn = reactive({
    width: 140,
    title: '操作',
    key: 'action',
    fixed: 'right',
    render(record: Role) {
      // 默认三个：查看 / 编辑 / 更多，其余操作统一收进下拉
      const actions: any[] = [
        { label: '查看', icon: EyeOutlined, onClick: () => openView(record) },
        {
          label: '编辑',
          icon: EditOutlined,
          ifShow: () => !record.inner,
          onClick: () => openEdit(record),
        },
      ];
      const dropDownActions: any[] = [
        {
          label: '分配权限',
          key: 'auth',
          icon: SafetyCertificateOutlined,
          // 根角色天然拥有全部权限，无需分配
          ifShow: () => record.id !== ROOT_ID,
        },
        { label: '删除', key: 'delete', icon: DeleteOutlined, ifShow: () => !record.inner },
      ];
      const select = (key: string) => {
        if (key === 'auth') openAuth(record);
        if (key === 'delete') confirmDelete(record);
      };
      return h(TableAction, { style: 'button', actions, dropDownActions, select });
    },
  });

  const loadDataTable = async (res: any) => {
    const params: Recordable = { ...res };
    if (keyword.value) params.q = keyword.value;
    return await getRoleList(params);
  };

  function reload() {
    actionRef.value?.reload();
  }

  function resetFilter() {
    keyword.value = '';
    reload();
  }

  // ---------------- 新增/编辑 ----------------
  // 层级属于后端权限模型的内部细节，前端不展示也不提交：
  // 新增时由后端自动取「当前登录角色层级 + 1」，且始终校验不得高于操作者
  function openAdd() {
    editId.value = '';
    Object.assign(form, { name: '', intro: '' });
    showForm.value = true;
  }

  function openEdit(row: Role) {
    editId.value = row.id;
    Object.assign(form, { name: row.name, intro: row.intro });
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
        await editRole(editId.value, { name: form.name, intro: form.intro });
        message.success('修改成功');
      } else {
        await addRole({ name: form.name, intro: form.intro });
        message.success('新增成功');
      }
      showForm.value = false;
      reload();
    } finally {
      submitting.value = false;
    }
  }

  function confirmDelete(row: Role) {
    if ((row.user_count ?? 0) > 0) {
      message.warning(`该角色下还有 ${row.user_count} 个用户，请先调整用户的角色`);
      return;
    }
    dialog.warning({
      title: '删除确认',
      content: `确定删除角色「${row.name}」？删除后其授权关系一并清理，不可恢复。`,
      positiveText: '删除',
      negativeText: '取消',
      onPositiveClick: async () => {
        await deleteRole(row.id);
        message.success('删除成功');
        reload();
      },
    });
  }

  // ---------------- 授权 ----------------
  const showAuth = ref(false);
  const authLoading = ref(false);
  const authRole = ref<Role | null>(null);
  const permNodes = ref<PermissionNode[]>([]);
  const permKeyword = ref('');
  const checkedKeys = ref<string[]>([]);
  const expandedKeys = ref<string[]>([]);
  const expanded = ref(true);

  const isRoot = computed(() => authRole.value?.id === ROOT_ID);

  /** 权限树转 naive-ui 结构：分组节点用 group: 前缀，叶子节点用权限 id
   * 根角色无需分配，整棵树设为不可勾选
   * label 里带上权限标识，这样 n-tree 的 pattern 搜索能同时命中名称与标识
   */
  function toTreeData(nodes: PermissionNode[], readonly = false): any[] {
    return (nodes || []).map((node) => {
      const group = !node.id;
      const key = group ? `group:${node.group_key || node.group}` : node.id;
      return {
        key,
        group,
        label: node.perm ? `${node.name}（${node.perm}）` : node.name,
        checkable: !readonly,
        children: node.children?.length ? toTreeData(node.children, readonly) : undefined,
      };
    });
  }

  /** 关键字过滤交给 n-tree 的 pattern：树数据保持不变，避免勾选状态被重算 */
  const treeView = computed(() => safeTree.value);

  const safeTree = computed(() => toTreeData(permNodes.value, isRoot.value));

  /** 叶子权限总数 */
  const leafTotal = computed(() => countLeaf(permNodes.value));
  function countLeaf(nodes: PermissionNode[]): number {
    let total = 0;
    (nodes || []).forEach((it) => {
      if (it.children?.length) total += countLeaf(it.children);
      else total += 1;
    });
    return total;
  }

  /** 各分组已选/总数，用于顶部统计 */
  const groupStats = computed(() =>
    (permNodes.value || []).map((grp) => {
      const children = grp.children || [];
      const checked = children.filter((it) => it.id && checkedKeys.value.includes(it.id)).length;
      return { key: grp.group_key || grp.group, group: grp.name, checked, total: children.length };
    })
  );

  function onChecked(keys: string[]) {
    // 只保留真实权限 id，过滤掉分组占位 key
    checkedKeys.value = keys.filter((it) => !it.startsWith('group:'));
  }

  function allLeafIds(): string[] {
    const out: string[] = [];
    const walk = (nodes: PermissionNode[]) => {
      (nodes || []).forEach((it) => {
        if (it.children?.length) walk(it.children);
        else if (it.id) out.push(it.id);
      });
    };
    walk(permNodes.value);
    return out;
  }

  function checkAll() {
    checkedKeys.value = allLeafIds();
  }

  function checkClear() {
    checkedKeys.value = [];
  }

  function checkInvert() {
    const all = allLeafIds();
    const set = new Set(checkedKeys.value);
    checkedKeys.value = all.filter((it) => !set.has(it));
  }

  function toggleExpand() {
    if (expanded.value) {
      expandedKeys.value = [];
    } else {
      expandedKeys.value = safeTree.value.map((it: any) => String(it.key));
    }
    expanded.value = !expanded.value;
  }

  async function openAuth(row: Role) {
    authRole.value = row;
    showAuth.value = true;
    permKeyword.value = '';
    authLoading.value = true;
    checkedKeys.value = [];
    try {
      const [tree, granted] = await Promise.all([getPermissionTree(), getRolePermissions(row.id)]);
      permNodes.value = tree;
      checkedKeys.value = granted.map((it) => it.id).filter(Boolean);
      expandedKeys.value = toTreeData(tree, row.id === ROOT_ID).map((it: any) => String(it.key));
      expanded.value = true;
    } finally {
      authLoading.value = false;
    }
  }

  async function submitAuth() {
    if (!authRole.value) return;
    submitting.value = true;
    try {
      await setRolePermissions(authRole.value.id, checkedKeys.value);
      message.success('授权成功');
      showAuth.value = false;
      reload();
    } finally {
      submitting.value = false;
    }
  }

  // ---------------- 查看已授权 ----------------
  const showView = ref(false);
  const viewLoading = ref(false);
  const viewRole = ref<Role | null>(null);
  const viewList = ref<Permission[]>([]);

  const viewGroups = computed(() => {
    const map = new Map<string, Permission[]>();
    viewList.value.forEach((it) => {
      const key = it.group || '未分组';
      if (!map.has(key)) map.set(key, []);
      map.get(key)!.push(it);
    });
    return Array.from(map.entries()).map(([group, items]) => ({ group, items }));
  });

  async function openView(row: Role) {
    viewRole.value = row;
    showView.value = true;
    viewLoading.value = true;
    viewList.value = [];
    try {
      viewList.value = await getRolePermissions(row.id);
    } finally {
      viewLoading.value = false;
    }
  }

  onMounted(() => {
    // 预取权限树，首次点开授权抽屉时更快
    getPermissionTree()
      .then((tree) => (permNodes.value = tree))
      .catch(() => undefined);
  });
</script>

<style lang="less" scoped>
  /* 授权抽屉：工具条与统计固定在顶部，权限树占满剩余高度并自己滚动 */
  .auth-bar {
    display: flex;
    flex-wrap: wrap;
    gap: 8px;
    align-items: center;
    padding-bottom: 8px;
    margin-bottom: 8px;
    border-bottom: 1px solid var(--divider-color);
  }

  .auth-stats {
    display: flex;
    flex-wrap: wrap;
    gap: 6px;
    align-items: center;
    padding-bottom: 8px;
    margin-bottom: 8px;
  }

  .auth-spin {
    display: flex;
    flex: 1;
    flex-direction: column;
    min-height: 0;

    :deep(.n-spin-container),
    :deep(.n-spin-content) {
      display: flex;
      flex: 1;
      flex-direction: column;
      min-height: 0;
    }
  }

  .auth-tree {
    flex: 1;
    min-height: 0;
    overflow: auto;
  }
</style>
