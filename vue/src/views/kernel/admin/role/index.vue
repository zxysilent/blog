<template>
    <div>
        <n-card :bordered="false" title="角色列表" size="small" :segmented="{ content: 'hard' }">
            <n-form inline label-placement="left" require-mark-placement="left" label-width="auto" ref="formRef">
                <n-form-item v-auth="'role_add'">
                    <n-button ghost attr-type="button" @click="$router.push({ name: 'role-add' })">
                        <template #icon>
                            <n-icon>
                                <AddCircleOutline />
                            </n-icon>
                        </template>
                        添加
                    </n-button>
                </n-form-item>
                <n-form-item>
                    <n-button type="info" attr-type="button" @click="init">
                        <template #icon>
                            <n-icon>
                                <Refresh />
                            </n-icon>
                        </template>
                        刷新
                    </n-button>
                </n-form-item>
            </n-form>
            <n-data-table size="small" max-height="70vh" remote :columns="tabCol" :data="tabData" :row-key="rowKey" :scroll-x="scrollWidth()" :pagination="tabPage" />
        </n-card>
        <n-drawer v-model:show="drawer" :width="520">
            <n-drawer-content title="角色授权">
                <n-form label-placement="top" :show-feedback="false" :model="roleForm" style="padding-bottom: 55px">
                    <template v-for="(items, key) in grantTree">
                        <n-form-item :label="key">
                            <n-checkbox-group v-model:value="roleForm.grant_ids">
                                <n-space item-style="display: flex;">
                                    <n-checkbox v-for="item in items" :label="item['name']" :value="item['id']" />
                                </n-space>
                            </n-checkbox-group>
                        </n-form-item>
                    </template>
                </n-form>
                <div style="text-align: center; width: 100%; position: fixed; bottom: 0; right: 0; border-top: 1px solid #e8e8e8; padding: 10px 16px; background: #fff; width: 520px">
                    <n-space justify="center">
                        <n-button type="warning" :loading="loading" @click="emitSave">提交保存</n-button>
                        <n-button type="success" @click="drawer = false">取消关闭</n-button>
                    </n-space>
                </div>
            </n-drawer-content>
        </n-drawer>
    </div>
</template>

<script lang="ts" setup>
import { h, ref, reactive, onMounted } from "vue";
import { NButton, NPopconfirm, useMessage, NTime, NTag, NPopover, NSpace } from "naive-ui";
import { Refresh, AddCircleOutline, FilterSharp } from "@vicons/ionicons5";
import { useRouter } from "vue-router";
import { useAdminStore } from "@/store/admin";
import { apiRolePage, apiRoleDrop, apiGrantTree, apiRoleGrantAll, apiRoleGrantEdit } from "@/api";
const adminStore = useAdminStore();
const router = useRouter();
const drawer = ref(false);
const loading = ref(false);
const message = useMessage();
const filter = reactive({ mult: "", pi: 1, ps: 15 });
const tabPage = reactive({
    page: 1,
    itemCount: 0,
    pageSize: 12,
    showSizePicker: true,
    pageSizes: [8, 10, 12, 15, 20, 30],
    prefix({ itemCount }) {
        return `Total is ${itemCount}.`;
    },
    onChange: (page) => {
        filter.pi = page;
        tabPage.page = page;
        init();
    },
    onUpdatePageSize: (pageSize) => {
        filter.pi = 1;
        tabPage.page = 1;
        filter.ps = pageSize;
        tabPage.pageSize = pageSize;
        console.log(pageSize);
        init();
    },
});
const tabData = ref([]);
const init = () => {
    apiRolePage(filter).then((resp) => {
        if (resp.code == 200) {
            tabData.value = resp.data.items;
            tabPage.itemCount = resp.data.count;
        } else {
            tabData.value = [];
            tabPage.page = 1;
            tabPage.itemCount = 0;
        }
    });
};
onMounted(() => {
    init();
});
const tabCol = [
    { title: "角色名称", key: "name", width: 150, ellipsis: { tooltip: { width: "trigger" } } },
    {
        title: "角色类型",
        width: 80,
        render(row) {
            if (row.inner) {
                return h(NTag, { size: "small", type: "success" }, { default: () => "内置" });
            }
            return h(NTag, { size: "small" }, { default: () => "普通" });
        },
    },
    { title: "角色说明", key: "intro", width: 200, ellipsis: { tooltip: { width: "trigger" } } },
    {
        title: "修改时间",
        key: "updated",
        width: 200,
        render(row) {
            return h(NTime, { type: "datetime", time: row.updated });
        },
    },
    {
        title: "创建时间",
        key: "created",
        width: 200,
        render(row) {
            return h(NTime, { type: "datetime", time: row.created });
        },
    },
    {
        title: "操作",
        key: "actions",
        width: 120,
        align: "center",
        // fixed:"right",
        render(row) {
            let actios = [
                h(
                    NButton,
                    {
                        size: "tiny",
                        block: true,
                        quaternary: true,
                        onClick: () => {
                            router.push({ name: "role-detail", params: { id: row.id } });
                        },
                    },
                    { default: () => "查看" }
                ),
            ];
            if (adminStore.Authed("role_grant")) {
                actios.push(
                    h(
                        NButton,
                        {
                            size: "tiny",
                            block: true,
                            quaternary: true,
                            onClick: () => {
                                roleForm.role_id = row.id;
                                roleGrantAll();
                                roleGrantTree();
                                drawer.value = true;
                            },
                        },
                        { default: () => "授权" }
                    )
                );
            }
            if (adminStore.Authed("role_edit")) {
                actios.push(h(NButton, { block: true, quaternary: true, size: "tiny", onClick: () => router.push({ name: "role-edit", params: { id: row.id } }) }, { default: () => "编辑" }));
            }
            if (adminStore.Authed("role_drop")) {
                actios.push(
                    h(
                        NPopconfirm,
                        { size: "20", onPositiveClick: () => emitDrop(row), trigger: "click" },
                        { trigger: () => h(NButton, { size: "tiny", block: true, quaternary: true, type: "error" }, { default: () => "删除" }), default: () => "确定删除？" }
                    )
                );
            }
            return h(
                NPopover,
                { size: "20", showArrow: false, trigger: "hover", placement: "bottom-start" },
                {
                    trigger: () => h(NButton, { size: "tiny", type: "primary", ghost: true }, { default: () => "更多", icon: () => h(FilterSharp) }),
                    default: () => h(NSpace, { vertical: true, size: "small", itemStyle: { flex: 1, color: "red" } }, { default: () => actios }),
                }
            );
        },
    },
];
const grantTree = ref([]);
const roleGrantTree = () => {
    apiGrantTree({}).then((resp) => {
        if (resp.code == 200) {
            grantTree.value = resp.data;
        } else {
            tabData.value = [];
        }
    });
};

const roleForm = reactive({
    role_id: 0,
    grant_ids: [],
});

const roleGrantAll = () => {
    roleForm.grant_ids = [];
    apiRoleGrantAll({ role_id: roleForm.role_id }).then((resp) => {
        if (resp.code == 200) {
            roleForm.grant_ids = resp.data.map((item) => {
                return item.id;
            });
        } else {
            roleForm.grant_ids = [];
        }
    });
};
const emitSave = () => {
    loading.value = true;
    apiRoleGrantEdit(roleForm).then((resp) => {
        loading.value = false;
        if (resp.code == 200) {
            message.success("保存成功", {
                onAfterLeave() {},
            });
        } else {
            message.error(resp.msg);
        }
    });
};
const scrollWidth = () => {
    let width = 0;
    tabCol.forEach((item) => {
        width += item.width;
    });
    return width;
};
const rowKey = (row) => {
    return row.id;
};
const emitDrop = (row) => {
    apiRoleDrop({ id: row.id }).then((resp) => {
        if (resp.code == 200) {
            message.success("删除成功", {
                onAfterLeave() {
                    init();
                },
            });
        } else {
            message.error(resp.msg);
        }
    });
};
</script>

<style lang="less" scoped></style>
