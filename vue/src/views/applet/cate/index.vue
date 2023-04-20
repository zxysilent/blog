<template>
    <div>
        <n-card :bordered="false" title="分类列表" size="small" :segmented="{ content: 'hard' }">
            <n-form inline label-placement="left" require-mark-placement="left" label-width="auto" ref="formRef">
                <n-form-item v-auth="'role_add'">
                    <n-button ghost attr-type="button" @click="$router.push({ name: 'cate-add' })">
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
import { NButton, NPopconfirm, useMessage, NTime, NPopover, NSpace, NTag } from "naive-ui";
import { Refresh, AddCircleOutline, FilterSharp } from "@vicons/ionicons5";
import { useRouter } from "vue-router";
import { useAdminStore } from "@/store/admin";
import { apiCatePage, apiCateDrop } from "@/api";
const adminStore = useAdminStore();
const router = useRouter();
const message = useMessage();
const filter = reactive({ mult: "", pi: 1, ps: 12 });
const tabPage = reactive({
    page: filter.pi,
    itemCount: 0,
    pageSize: filter.ps,
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
    apiCatePage(filter).then((resp) => {
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
    { title: "分类名称", key: "name", width: 200, ellipsis: { tooltip: { width: "trigger" } } },
    { title: "分类说明", key: "intro", width: 200, ellipsis: { tooltip: { width: "trigger" } } },
    {
        title: "分类颜色",
        width: 100,
        render(row) {
            // return h(NColorPicker, { value: row.color, modes: ["hex"], size: "small", disabled: true });
            // return h(NTag, { color: { color: row.color, textColor: "#dcf3fd", borderColor: "#fafafc" }, size: "small" }, { default: () => row.color });
            return h(NButton, { color: row.color, size: "tiny" }, { default: () => row.color });
        },
    },
    {
        title: "类型",
        width: 80,
        render(row) {
            if (row.kind == 1) {
                return h(NTag, { size: "small", type: "info" }, { default: () => "文章" });
            } else if (row.kind == 3) {
                return h(NTag, { size: "small", type: "success" }, { default: () => "笔记" });
            } else {
                return h(NTag, { size: "small" }, { default: () => "页面" });
            }
        },
    },

    {
        title: "修改时间",
        key: "updated",
        width: 150,
        render(row) {
            return h(NTime, { type: "datetime", time: row.updated });
        },
    },
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
        width: 100,
        align: "center",
        render(row) {
            let actios = [
                h(
                    NButton,
                    {
                        size: "tiny",
                        block: true,
                        quaternary: true,
                        onClick: () => {
                            router.push({ name: "cate-detail", params: { id: row.id } });
                        },
                    },
                    { default: () => "查看" }
                ),
            ];
            if (adminStore.Authed("role_edit")) {
                actios.push(h(NButton, { block: true, quaternary: true, size: "tiny", onClick: () => router.push({ name: "cate-edit", params: { id: row.id } }) }, { default: () => "编辑" }));
            }
            if (adminStore.Authed("role_drop")) {
                actios.push(
                    h(
                        NPopconfirm,
                        { size: "20", onPositiveClick: () => emitDrop(row), trigger: "click" },
                        { trigger: () => h(NButton, { size: "tiny", block: true, quaternary: true, type: "error" }, { default: () => "删除" }), default: () => "确定删除?" }
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
    apiCateDrop({ id: row.id }).then((resp) => {
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
