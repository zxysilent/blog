<template>
    <div>
        <n-card :bordered="false" title="博文列表" class="mt-4 proCard" size="small" :segmented="{ content: 'hard' }">
            <n-form inline label-placement="left" require-mark-placement="left" label-width="auto" ref="formRef">
                <n-form-item v-auth="'role_add'">
                    <n-button ghost attr-type="button" @click="$router.push({ name: 'post-add' })">
                        <template #icon>
                            <n-icon>
                                <AddCircleOutline />
                            </n-icon>
                        </template>
                        添加
                    </n-button>
                </n-form-item>
                <n-form-item label="类型">
                    <n-select style="width: 80px" v-model:value="filter.kind" :options="kindAll" />
                </n-form-item>
                <n-form-item label="分类">
                    <n-select style="width: 120px" v-model:value="filter.cate_id" :options="cateAll" />
                </n-form-item>
                <n-form-item label="标题">
                    <n-input placeholder="标题" v-model:value="filter.mult" clearable />
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
            <n-modal v-model:show="noteShare.show" preset="dialog" title="笔记分享">
                <n-space vertical>
                    <n-p>有效期</n-p>
                    <n-radio-group style="display: block" v-model:value="noteShare.day" @update:value="noteShare.url = ''">
                        <n-space justify="space-around">
                            <n-radio :value="1">1天 </n-radio>
                            <n-radio :value="3">3天 </n-radio>
                            <n-radio :value="7">7天 </n-radio>
                            <n-radio :value="-1">永久 </n-radio>
                            <n-radio :value="0">自定义 </n-radio>
                        </n-space>
                    </n-radio-group>
                    <n-space justify="space-around" v-if="showDiy" style="display: block">
                        <n-input-number min="1" max="30" v-model:value="noteShare.diy_day" :precision="0" @update:value="noteShare.url = ''">
                            <template #suffix> 天 </template>
                        </n-input-number>
                    </n-space>
                    <n-space v-if="noteShare.url" style="display: block">
                        <n-input readonly placeholder="" :value="noteShare.url" title="复制">
                            <template #suffix>
                                <n-button @click="copyURL" text>
                                    <template #icon>
                                        <n-icon>
                                            <CopyOutline />
                                        </n-icon>
                                    </template>
                                </n-button>
                            </template>
                        </n-input>
                    </n-space>
                </n-space>
                <template #action>
                    <n-button type="info" @click="emitShare"> 创建链接 </n-button>
                </template>
            </n-modal>
        </n-card>
    </div>
</template>

<script lang="ts" setup>
import { h, ref, reactive, onMounted, computed } from "vue";
import type { VNode } from "vue";
import { NButton, NPopconfirm, useMessage, NTime, NPopover, NSpace, NTag } from "naive-ui";
import { Refresh, AddCircleOutline, FilterSharp, Search, CopyOutline } from "@vicons/ionicons5";
import { useRouter } from "vue-router";
import { useAdminStore } from "@/store/admin";
import { apiPostPage, apiPostDrop, apiCateList, apiPostShare } from "@/api";
const kindAll = [
    { label: "所有", value: 0 },
    { label: "文章", value: 1 },
    { label: "页面", value: 2 },
    { label: "笔记", value: 3 },
];
const cateAll = ref([]);
const noteShare = reactive({
    show: false,
    diy_day: 1,
    day: 1,
    url: "",
    path: "",
    post_id: 0,
});
const showDiy = computed(() => {
    return noteShare.day === 0;
});
const adminStore = useAdminStore();
const router = useRouter();
const message = useMessage();
const filter = reactive({ kind: 0, cate_id: 0, mult: "", pi: 1, ps: 12 });
const tabPage = reactive({
    page: 1,
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
const copyURL = () => {
    navigator.clipboard
        .writeText(noteShare.url)
        .then(() => {
            message.success("已复制到剪贴板");
        })
        .catch(() => {
            message.error("复制失败,请手动复制");
        });
};
const preInit = () => {
    apiCateList({}).then((resp) => {
        if (resp.code == 200) {
            resp.data.unshift({ id: 0, name: "所有" });
            cateAll.value = resp.data.map((item) => {
                return {
                    label: item.name,
                    value: item.id,
                };
            });
        } else {
            cateAll.value = [];
        }
    });
};
const init = () => {
    apiPostPage(filter).then((resp) => {
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
const emitShare = () => {
    let data = {
        post_id: noteShare.post_id,
        day: noteShare.day,
    };
    if (noteShare.day == 0) {
        data.day = noteShare.diy_day;
    }
    apiPostShare(data).then((resp) => {
        if (resp.code == 200) {
            noteShare.url = `${import.meta.env.VITE_APP_SRV}/notes/${noteShare.path}.html?s=${resp.data}`;
        } else {
            noteShare.url = "";
            message.warning("分享失败,请重试");
        }
    });
};
onMounted(() => {
    preInit();
    init();
});
const tabCol = [
    {
        title: "#",
        key: "key",
        width: 60,
        render: (_, index) => {
            return `${index + 1}`;
        },
    },
    { title: "路径", key: "path", minWidth: 200, ellipsis: { tooltip: { width: "trigger" } } },
    { title: "标题", key: "title", minWidth: 200, ellipsis: { tooltip: { width: "trigger" } } },
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
        title: "状态",
        width: 80,
        render(row) {
            if (row.status) {
                return h(NTag, { size: "small", type: "success" }, { default: () => "发布" });
            }
            return h(NTag, { size: "small" }, { default: () => "草稿" });
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
        width: 120,
        align: "center",
        render(row) {
            let actios: VNode[] = [];
            let preview = "";
            switch (row.kind) {
                case 1:
                    preview = "posts";
                    break;
                case 2:
                    preview = "pages";
                    break;
                case 3:
                    preview = "notes";
                    break;
            }
            if (preview) {
                actios.push(h(NButton, { block: true, size: "tiny", quaternary: true, target: "_blank", tag: "a", href: `${import.meta.env.VITE_APP_SRV}/${preview}/${row.path}.html` }, { default: () => "预览" }));
            }
            if (preview && row.kind == 3) {
                actios.push(
                    h(
                        NButton,
                        {
                            block: true,
                            size: "tiny",
                            quaternary: true,
                            onClick: () => {
                                noteShare.post_id = row.id;
                                noteShare.day = 1;
                                noteShare.path = row.path;
                                noteShare.diy_day = 1;
                                noteShare.url = "";
                                noteShare.show = true;
                            },
                        },
                        { default: () => "分享" }
                    )
                );
            }
            actios.push(
                h(
                    NButton,
                    {
                        size: "tiny",
                        block: true,
                        quaternary: true,
                        onClick: () => {
                            router.push({ name: "post-detail", params: { id: row.id } });
                        },
                    },
                    { default: () => "查看" }
                )
            );
            if (adminStore.Authed("role_edit")) {
                actios.push(h(NButton, { block: true, quaternary: true, size: "tiny", onClick: () => router.push({ name: "post-edit", params: { id: row.id } }) }, { default: () => "编辑" }));
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
        //@ts-ignore
        width += item.minWidth ? item.minWidth : item.width;
    });
    return width;
};
const rowKey = (row) => {
    return row.id;
};
const emitDrop = (row) => {
    apiPostDrop({ id: row.id }).then((resp) => {
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
