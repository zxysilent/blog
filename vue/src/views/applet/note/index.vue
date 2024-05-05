<style lang="less" scoped>
#search {
    cursor: pointer;
}

#search:hover {
    color: #57a3f3;
}

.n-list.n-list--hoverable .n-list-item {
    padding: 6px 4px;
}

.layout-sider {
    min-height: 100vh;
    box-shadow: 2px 0 8px 0 rgb(29 35 41 / 5%);
    position: relative;
    z-index: 13;
    transition: all 0.2s ease-in-out;
}
</style>
<template>
    <n-layout position="absolute" has-sider>
        <n-layout-sider
            collapse-mode="transform"
            :collapsed-width="0"
            :width="260"
            show-trigger="bar"
            :inverted="inverted"
            :collapsed="collapsed"
            @collapse="collapsed = true"
            @expand="collapsed = false"
            style="height: 100vh">
            <n-scrollbar style="border-right: 1px solid #f2f6fc">
                <n-space vertical :size="12">
                    <n-input type="text" placeholder="搜索" clearable v-model:value="mult">
                        <template #suffix>
                            <n-icon id="search" @click="init">
                                <SearchOutline></SearchOutline>
                            </n-icon>
                        </template>
                    </n-input>
                    <n-tree :data="notes" block-node :render-label="renderLabel"  :render-suffix="renderSuffix"  :selected-keys="selectedKeys" :cancelable="false" :node-props="nodeProps" key-field="id" />
                </n-space>
            </n-scrollbar>
        </n-layout-sider>
        <n-layout-content content-style="padding: 4px 0 0 24px;">
            <n-form :label-width="100" :model="dataForm" label-placement="left" require-mark-placement="left" :rules="dataRules" ref="dataRef">
                <n-grid responsive="screen" cols="14 s:14 m:14 l:24 xl:16 xxl:16" :x-gap="6" :y-gap="12">
                    <!-- <n-grid-item span="12"> </n-grid-item> -->
                    <n-form-item-gi :show-feedback="false" >
                        <n-button ghost attr-type="button">
                            <template #icon>
                                <n-icon>
                                    <AddCircleOutline />
                                </n-icon>
                            </template>
                            添加
                        </n-button>
                    </n-form-item-gi>
                    <n-form-item-gi :show-feedback="false" >
                        <n-button type="info" attr-type="button" @click="init">
                            <template #icon>
                                <n-icon>
                                    <Refresh />
                                </n-icon>
                            </template>
                            刷新
                        </n-button>
                    </n-form-item-gi>
                    <n-form-item-gi :show-feedback="false" :span="10" path="title"> <n-input v-model:value="dataForm.title" placeholder="请输入标题" /> </n-form-item-gi>
                    <n-grid-item :show-feedback="false" span="24">
                        <Note ref="mdRef" v-model:value="dataForm.markdown" @change="onChange" @save="onSave" />
                    </n-grid-item>
                    <!-- <n-grid-item span="24"> </n-grid-item> -->
                </n-grid>
            </n-form>
        </n-layout-content>
    </n-layout>
</template>
<script lang="ts" setup>
import { apiNoteSave, apiNoteList, apiTagList, apiPostGet } from "@/api";
// import Card from "./card.vue";
import { Refresh, AddCircleOutline, AddOutline, SearchOutline } from "@vicons/ionicons5";
import { useMessage, NIcon, NButton, NEllipsis } from "naive-ui";
import { Note } from "@/components/Editor";
import { ref, onMounted, nextTick, h, withModifiers } from "vue";
// function renderIcon(icon: Component) { return () => h(NIcon, null, { default: () => h(icon) }); }
// interface Cate {
//     id: number;
//     pid: number;
//     name: string;
//     color: string;
//     intro: string;
//     created: number;
// }
interface Post {
    id: number;
    title: string;
    markdown: string;
    level: number;
    created: number;
}
const selectedKeys = ref<any[]>([])
function renderLabel({ option }: { option: Post }) {
    return h(NEllipsis, { style: { maxWidth: `${220 - option.level * 20}px` }, tooltip: { placement: "top-start" } }, ()=>`${option.title}`);
}

const notes = ref<Post[]>();
const noteActive = ref(0);
const mult = ref("");
const mdRef = ref();
const dataRef = ref();
const message = useMessage();
const dataForm = ref({
    id: 0,
    title: "",
    path: "",
    kind: 3,
    pid: 0,
    status: 2,
    summary: "",
    allow: false,
    richtext: "",
    markdown: "",
    tags: [],
    updated: 0,
    created: 0,
});
const inverted = ref(false);
const collapsed = ref(false);
const tagAll = ref([]);
const dataRules = {
    title: { required: true, message: "请输入标题", trigger: "blur" },
    path: { required: true, message: "请输入访问路径", trigger: "blur" },
};

const create = () => {
    dataForm.value = {
        id: 0,
        title: "",
        path: "",
        kind: 3,
        pid: noteActive.value,
        status: 2,
        summary: "",
        allow: false,
        richtext: "",
        markdown: "",
        tags: [],
        updated: 0,
        created: 0,
    };
    noteActive.value = 0;
};
const onChange = (val) => {
    console.log(val);
    dataForm.value.richtext = val;
};
const renderSuffix = ({ option }: { option: Post }) => {
    console.log(option)
    if (option.id != noteActive.value || option.level > 4) {
        return ""
    }
    return h(NButton, { text: true, type: "primary", size: 'tiny', onClick: withModifiers(() => create(), ['stop']), }, { icon: () => h(AddOutline) })
}
const nodeProps = ({ option }: { option: Post }) => {
    return {
        onClick() {
            console.log(option)
            if (noteActive.value === option.id) {
                return;
            }
            noteActive.value = option.id;
            selectedKeys.value = [option.id]
            noteGet();
        }
    };
};
const noteGet = () => {
    console.log("active", noteActive.value);
    apiPostGet({ id: noteActive.value }).then((resp) => {
        if (resp.code == 200) {
            dataForm.value = resp.data;
            nextTick(() => {
                mdRef.value.emptyHistory();
            });
        } else {
            dataForm.value = {
                id: 0,
                title: "",
                path: "",
                kind: 1,
                pid: 0,
                status: 1,
                summary: "",
                allow: false,
                richtext: "",
                markdown: "",
                tags: [],
                updated: 0,
                created: 0,
            };
        }
    });
};

const init = () => {
    apiNoteList({ mult: mult.value }).then((resp) => {
        if (resp.code == 200) {
            notes.value = resp.data;
            if (notes.value?.length) {
                noteActive.value = notes.value[0].id;
                selectedKeys.value = [notes.value[0].id]
                noteGet();
            }
        } else {
            notes.value = [];
        }
    });
};
const preInit = () => {
    apiTagList({}).then((resp) => {
        if (resp.code == 200) {
            tagAll.value = resp.data.map((item) => {
                return {
                    label: item.name,
                    value: item.id,
                };
            });
        } else {
            tagAll.value = [];
        }
    });
};

const watchWidth = () => {
    const Width = document.body.clientWidth;
    if (Width <= 950) {
        collapsed.value = true;
    } else collapsed.value = false;
};

onMounted(() => {
    window.addEventListener("resize", watchWidth);
    preInit();
    init();
});
const onSave = () => {
    dataRef.value.validate((errors) => {
        if (!errors) {
            apiNoteSave(dataForm.value).then((resp) => {
                if (resp.code == 200) {
                    // message.success("保存成功");
                } else {
                    message.error(resp.msg);
                }
            });
        }
    });
};
</script>
