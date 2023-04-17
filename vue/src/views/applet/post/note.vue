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
.active {
    color: #2d8cf0;
    background-color: #eaf3fd;
}
</style>
<template>
    <n-layout has-sider>
        <n-layout-sider
            collapse-mode="transform"
            :collapsed-width="0"
            :width="160"
            show-trigger="bar"
            :inverted="inverted"
            content-style="padding: 8px;"
            :collapsed="collapsed"
            @collapse="collapsed = true"
            @expand="collapsed = false"
            style="height: 100vh"
            bordered>
            <!-- <n-menu :inverted="inverted" :collapsed-width="0" :icon-size="0" :indent="0" :options="cateAll" :collapsed="collapsed" /> -->
            <n-space vertical>
                <n-space>
                    <n-button size="small" type="info" dashed block @click="create">
                        <template #icon>
                            <n-icon>
                                <AddSharp></AddSharp>
                            </n-icon> </template
                        >新建
                    </n-button>
                    <n-button size="small" dashed block @click="$router.push({ name: 'post-index' })"> 返回 </n-button>
                </n-space>
                <n-input type="text" placeholder="笔记" clearable v-model:value="mult">
                    <template #suffix>
                        <n-icon id="search" @click="init">
                            <SearchOutline></SearchOutline>
                        </n-icon>
                    </template>
                </n-input>
                <!-- <n-button v-for="item in noteAll" quaternary block @click="init(item.id)" style="text-align: left">
                    <n-ellipsis>
                        {{ item.title }}
                    </n-ellipsis>
                </n-button> -->
            </n-space>
            <n-list hoverable clickable :show-divider="false" style="margin-top: 0.5rem">
                <n-list-item v-for="item in noteAll" @click="onActive(item.id)" :class="{ active: active == item.id }">
                    <n-ellipsis style="width: 135px">
                        {{ item.title }}
                    </n-ellipsis>
                </n-list-item>
            </n-list>
        </n-layout-sider>
        <n-layout-content content-style="padding: 4px 0 0 24px;">
            <n-form :label-width="100" :model="dataForm" label-placement="left" require-mark-placement="left" :rules="dataRules" ref="dataRef">
                <n-grid responsive="screen" cols="12 s:12 m:16 l:24 xl:24 xxl:24" :x-gap="6" :y-gap="12">
                    <!-- <n-grid-item span="12"> </n-grid-item> -->
                    <n-form-item-gi :show-feedback="false" :span="12" path="title"> <n-input v-model:value="dataForm.title" placeholder="请输入标题" /> </n-form-item-gi>
                    <n-form-item-gi :show-feedback="false" :span="5">
                        <n-select placeholder="标签" max-tag-count="responsive" style="width: 100%" multiple remote v-model:value="dataForm.tags" :options="tagAll" />
                    </n-form-item-gi>
                    <n-grid-item :show-feedback="false" span="24"> <Note ref="mdRef" v-model:value="dataForm.markdown" @change="onChange" @save="onSave" /> </n-grid-item>
                    <!-- <n-grid-item span="24"> </n-grid-item> -->
                </n-grid>
            </n-form>
        </n-layout-content>
    </n-layout>
    <!-- <n-card :bordered="false" size="small" :segmented="{ content: 'hard' }">
       
    </n-card> -->
</template>
<script lang="ts" setup>
import { apiPostSave, apiPostList, apiTagList, apiPostGet } from "@/api";
import { Note } from "@/components/Editor";
import { AddSharp, SearchOutline } from "@vicons/ionicons5";
import { ref, onMounted, nextTick } from "vue";
import { useMessage } from "naive-ui";
interface Note {
    id: number;
    title: string;
    markdown: string;
}
const active = ref(0);
const mult = ref("");
const mdRef = ref();
const dataRef = ref();
const message = useMessage();
const dataForm = ref({
    id: 0,
    title: "",
    path: "",
    kind: 3,
    status: 2,
    summary: "",
    cate_id: 0,
    allow: false,
    richtext: "",
    markdown: "",
    tags: [],
    updated: 0,
    created: 0,
});
const inverted = ref(false);
const collapsed = ref(false);
const noteAll = ref<Note[]>([]);
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
        status: 2,
        summary: "",
        cate_id: 0,
        allow: false,
        richtext: "",
        markdown: "",
        tags: [],
        updated: 0,
        created: 0,
    };
    active.value = 0;
};

const onChange = (val) => {
    console.log(val);
    dataForm.value.richtext = val;
};
const onActive = (id) => {
    if (active.value === id) {
        return;
    }
    active.value = id;
    detail();
};
const detail = () => {
    console.log("active", active.value);
    apiPostGet({ id: active.value }).then((resp) => {
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
                status: 1,
                summary: "",
                cate_id: 0,
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
    apiPostList({ kind: 3, mult: mult.value }).then((resp) => {
        if (resp.code == 200) {
            noteAll.value = resp.data;
            if (!active.value) {
                active.value = noteAll.value[0].id;
            }
            detail();
        } else {
            noteAll.value = [];
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
            apiPostSave(dataForm.value).then((resp) => {
                if (resp.code == 200) {
                    init();
                    message.success("保存成功");
                } else {
                    message.error(resp.msg);
                }
            });
        }
    });
};
</script>
