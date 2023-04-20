<template>
    <div>
        <n-card :bordered="false" title="日志列表" size="small" :segmented="{ content: 'hard' }">
            <n-form inline label-placement="left" require-mark-placement="left" label-width="auto" ref="formRef">
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
            <n-drawer-content title="日志详情"> {{ detail }} </n-drawer-content>
        </n-drawer>
    </div>
</template>

<script lang="ts" setup>
import { h, ref, reactive, onMounted } from "vue";
import { NButton, NTime } from "naive-ui";
import { Refresh } from "@vicons/ionicons5";
import { apiLogPage } from "@/api";
const drawer = ref(false);
const detail = ref("");
const tabPage = reactive({
    page: 1,
    mult: "",
    itemCount: 0,
    pageSize: 12,
    showSizePicker: true,
    pageSizes: [8, 10, 12, 15, 20, 30],
    prefix({ itemCount }) {
        return `Total is ${itemCount}.`;
    },
    onChange: (page) => {
        tabPage.page = page;
        init();
    },
    onUpdatePageSize: (pageSize) => {
        tabPage.pageSize = pageSize;
        tabPage.page = 1;
        console.log(pageSize);
        init();
    },
});
const tabData = ref([]);
const init = () => {
    apiLogPage({ pi: tabPage.page, ps: tabPage.pageSize, mult: tabPage.mult }).then((resp) => {
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
    {
        title: "用户",
        width: 150,
        tooltip: { width: "trigger" },
        render(row) {
            return h("div", row.admin ? row.admin.name + "(" + row.admin.num + ")" : "");
        },
    },
    { title: "模块", key: "module", width: 120, tooltip: { width: "trigger" } },
    { title: "事件", key: "action", width: 200, tooltip: { width: "trigger" } },
    {
        title: "时间",
        key: "created",
        width: 200,
        render(row) {
            return h(NTime, { type: "datetime", time: row.created });
        },
    },
    { title: "数据", key: "data", width: 250, ellipsis: { tooltip: { width: "trigger" } } },
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
                        type: "success",
                        tertiary: true,
                        style: { marginRight: "6px" },
                        onClick: () => {
                            drawer.value = true;
                            detail.value = row.data.replaceAll("\\/", "/").replaceAll('\\"', '"');
                            console.log(detail.value);
                        },
                    },
                    { default: () => "查看" }
                ),
            ];
            return h("div", actios);
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
</script>
<style lang="less" scoped></style>
