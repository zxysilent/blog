<template>
    <div>
        <n-card :bordered="false" title="权限列表" size="small" :segmented="{ content: 'hard' }">
            <n-form inline label-placement="left" require-mark-placement="left" label-width="auto" ref="formRef">
                <n-form-item v-auth="'grant_add'">
                    <n-button ghost attr-type="button" @click="$router.push({ name: 'grant-add' })">
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
    </div>
</template>

<script lang="ts" setup>
import { h, ref, reactive, onMounted } from "vue";
import { NButton, NPopconfirm, useMessage, NTime, NTag, NPopover, NSpace } from "naive-ui";
import { Refresh, AddCircleOutline, FilterSharp } from "@vicons/ionicons5";
import { useRouter } from "vue-router";
import { useAdminStore } from "@/store/admin";
import { apiGrantPage, apiGrantDrop } from "@/api";
const adminStore = useAdminStore();
const router = useRouter();
const message = useMessage();
const filter = reactive({ mult: "", pi: 1, ps: 12 });
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
    apiGrantPage(filter).then((resp) => {
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
    { title: "权限名称", key: "name", width: 150, ellipsis: { tooltip: { width: "trigger" } } },
    { title: "权限标识", key: "guid", width: 150, ellipsis: { tooltip: { width: "trigger" } } },
    { title: "权限分组", key: "group", width: 150, ellipsis: { tooltip: { width: "trigger" } } },
    {
        title: "权限类型",
        width: 80,
        render(row) {
            if (row.inner) {
                return h(NTag, { type: "success", size: "small" }, { default: () => "内置" });
            }
            return h(NTag, { size: "small" }, { default: () => "普通" });
        },
    },
    { title: "排序", key: "sort", width: 100, ellipsis: { tooltip: { width: "trigger" } } },
    {
        title: "路由",
        width: 200,
        ellipsis: { tooltip: { width: "trigger" } },
        render(row) {
            return row.routes ? `[${row.routes.join(",")}]` : "";
        },
    },
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
        width: 100,
        align: "center",
        render(row) {
            let actios = [
                h(
                    NButton,
                    {
                        block: true,
                        size: "tiny",
                        quaternary: true,
                        onClick: () => {
                            router.push({ name: "grant-detail", params: { id: row.id } });
                        },
                    },
                    { default: () => "查看" }
                ),
            ];
            if (adminStore.Authed("grant_edit")) {
                actios.push(h(NButton, { block: true, size: "tiny", quaternary: true, onClick: () => router.push({ name: "grant-edit", params: { id: row.id } }) }, { default: () => "编辑" }));
            }
            if (adminStore.Authed("grant_drop")) {
                actios.push(
                    h(
                        NPopconfirm,
                        { size: "20", onPositiveClick: () => emitDrop(row), trigger: "click" },
                        { trigger: () => h(NButton, { block: true, quaternary: true, size: "tiny", type: "error" }, { default: () => "删除" }), default: () => "确定删除？" }
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
    apiGrantDrop({ id: row.id }).then((resp) => {
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
