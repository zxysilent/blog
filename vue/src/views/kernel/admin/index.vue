<template>
    <div>
        <n-card :bordered="false" title="用户列表" size="small" :segmented="{ content: 'hard' }">
            <n-form inline label-placement="left" require-mark-placement="left" label-width="auto" ref="formRef">
                <n-form-item v-auth="'admin_add'">
                    <n-button ghost attr-type="button" @click="$router.push({ name: 'admin-add' })">
                        <template #icon>
                            <n-icon>
                                <AddCircleOutline />
                            </n-icon>
                        </template>
                        添加
                    </n-button>
                </n-form-item>
                <n-form-item label="角色">
                    <n-select style="width: 119px" laceholder="角色" remote v-model:value="filter.role_id" :options="roleAll" />
                </n-form-item>
                <n-form-item label="账号">
                    <n-input placeholder="账号" v-model:value="filter.mult" clearable />
                </n-form-item>
                <n-form-item>
                    <n-button type="primary" attr-type="button" @click="init">
                        <template #icon>
                            <n-icon>
                                <Search />
                            </n-icon>
                        </template>
                        搜索
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
    </div>
</template>

<script lang="ts" setup>
import { h, ref, reactive, onMounted } from "vue";
import { NButton, NPopconfirm, useMessage, NTime } from "naive-ui";
import { Search, Refresh, AddCircleOutline } from "@vicons/ionicons5";
import { useRouter } from "vue-router";
import { useAdminStore } from "@/store/admin";
import { apiAdminPage, apiAdminDrop } from "@/api";
import { apiRoleAll } from "@/api";
const adminStore = useAdminStore();
const router = useRouter();
const message = useMessage();
const roleAll = ref([]);
const filter = reactive({ role_id: 0, mult: "", pi: 1, ps: 15 });
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
const preInit = () => {
    apiRoleAll({}).then((resp) => {
        if (resp.code == 200) {
            resp.data.unshift({ id: 0, name: "所有" });
            roleAll.value = resp.data.map((item) => {
                return {
                    label: item.name,
                    value: item.id,
                };
            });
        } else {
            roleAll.value = [];
        }
    });
};
const init = () => {
    apiAdminPage(filter).then((resp) => {
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
    preInit();
    init();
});
const tabCol = [
    { title: "登录账号", key: "num", width: 100, ellipsis: { tooltip: { width: "trigger" } } },
    { title: "用户姓名", key: "name", width: 150, ellipsis: { tooltip: { width: "trigger" } } },
    {
        title: "用户角色",
        key: "role_id",
        width: 100,
        render(row) {
            return row.role.name;
        },
    },
    { title: "电话号码", key: "phone", width: 150, ellipsis: { tooltip: { width: "trigger" } } },
    { title: "联系邮箱", key: "email", width: 150, ellipsis: { tooltip: { width: "trigger" } } },
    {
        title: "创建时间",
        key: "created",
        width: 150,
        render(row) {
            return h(NTime, { type: "datetime", time: row.created });
        },
    },
    {
        title: "操作",
        key: "actions",
        width: 150,
        render(row) {
            let actios = [h(NButton, { size: "tiny", type: "success", style: { marginRight: "6px" }, onClick: () => router.push({ name: "admin-detail", params: { id: row.id } }) }, { default: () => "查看" })];
            if (adminStore.Authed("admin_edit")) {
                actios.push(h(NButton, { size: "tiny", type: "warning", style: { marginRight: "6px" }, onClick: () => router.push({ name: "admin-edit", params: { id: row.id } }) }, { default: () => "编辑" }));
            }
            if (adminStore.Authed("admin_drop")) {
                actios.push(
                    h(
                        NPopconfirm,
                        { size: "20", onPositiveClick: () => emitDrop(row), trigger: "click" },
                        { trigger: () => h(NButton, { size: "tiny", type: "error" }, { default: () => "删除" }), default: () => "确定删除？" }
                    )
                );
            }
            return h("div", actios);
            // return h("div", [
            //     // h(NButton, { size: "tiny", type: "warning", style: { marginRight: "6px" }, onClick: () => sendAdmin(row) }, { default: () => "修改" }),

            // ]);
        },
    },
];
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
    apiAdminDrop({ id: row.id }).then((resp) => {
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
