<template>
    <div>
        <n-card :bordered="false" title="字典列表" class="mt-4 proCard" size="small" :segmented="{ content: 'hard' }">
            <n-form inline label-placement="left" require-mark-placement="left" label-width="auto" ref="formRef">
                <n-form-item v-auth="'user_add'">
                    <n-button ghost attr-type="button" @click="$router.push({ name: 'dict-add' })">
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
import type { VNode } from "vue";
import { h, ref, reactive, onMounted } from "vue";
import { Refresh, AddCircleOutline, FilterSharp } from "@vicons/ionicons5";
import { NButton, NTime, NTag, useMessage, NPopconfirm, NPopover, NSpace } from "naive-ui";
import { useRouter } from "vue-router";
import { useAdminStore } from "@/store/admin";
import { apiDictPage, apiDictDrop } from "@/api";
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
    apiDictPage(filter).then((resp) => {
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
    { title: "配置名", key: "key", width: 150, ellipsis: { tooltip: { width: "trigger" } } },
    { title: "配置值", key: "value", minWidth: 150, ellipsis: { tooltip: { width: "trigger" } } },
    { title: "配置说明", key: "intro", width: 150, ellipsis: { tooltip: { width: "trigger" } } },
    {
        title: "类型",
        width: 80,
        render(row) {
            if (row.inner == 0) {
                return h(NTag, { size: "small" }, { default: () => "普通" });
            } else if (row.inner == 1) {
                return h(NTag, { size: "small", type: "success" }, { default: () => "内置" });
            }
            return "";
        },
    },
    {
        title: "修改时间",
        width: 150,
        render(row) {
            return h(NTime, { type: "datetime", time: row.updated });
        },
    },
    {
        title: "创建时间",
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
            let actios: VNode[] = [];
            actios.push(h(NButton, { block: true, size: "tiny", quaternary: true, onClick: () => router.push({ name: "dict-detail", params: { key: row.key } }) }, { default: () => "查看" }));
            actios.push(h(NButton, { block: true, size: "tiny", quaternary: true, target: "_blank", tag: "a", href: "/api/dict/" + row.key, style: { marginRight: "6px" } }, { default: () => "API" }));
            if (adminStore.Authed("dict_edit") && row.inner == 0) {
                actios.push(h(NButton, { block: true, size: "tiny", quaternary: true, onClick: () => router.push({ name: "dict-edit", params: { key: row.key } }) }, { default: () => "编辑" }));
            }
            if (adminStore.Authed("dict_drop") && row.inner == 0) {
                actios.push(
                    h(
                        NPopconfirm,
                        { size: "20", onPositiveClick: () => emitDrop(row), trigger: "click" },
                        { trigger: () => h(NButton, { block: true, size: "tiny", quaternary: true, type: "error" }, { default: () => "删除" }), default: () => "确定删除？" }
                    )
                );
            }
            // return h("div", actios);
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
const emitDrop = (row) => {
    apiDictDrop({ key: row.key }).then((resp) => {
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
const scrollWidth = () => {
    let width = 0;
    tabCol.forEach((item) => {
        width += item.width ? item.width : item.minWidth ? item.minWidth : 150;
    });
    return width;
};
const rowKey = (row) => {
    return row.id;
};
</script>
<style lang="less" scoped></style>
