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
// .active {
//     color: #2d8cf0;
//     background-color: #eaf3fd;
// }

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
        <!-- <Menu :options="cateAll"></Menu> -->
        <!-- <n-layout-sider collapsed position="static" collapse-mode="width" :collapsed-width="64" inverted class="layout-sider">
            <n-menu inverted key-field="id" label-field="name" :collapsed-width="64" :collapsed-icon-size="22" :options="cateAll"  />
        </n-layout-sider> -->
        <n-layout-content>
            <n-layout has-sider>
                <n-layout-sider
                    collapse-mode="transform"
                    :collapsed-width="0"
                    :width="160"
                    show-trigger="bar"
                    :inverted="inverted"
                    :collapsed="collapsed"
                    @collapse="collapsed = true"
                    @expand="collapsed = false"
                    style="height: 100vh">
                    <n-scrollbar style="border-right: 1px solid #f2f6fc">
                        <Card v-for="item in noteAll" @click="onActive(item.id)" :class="{ active: noteActive == item.id }"> {{ item.title }}</Card>
                    </n-scrollbar>
                </n-layout-sider>
                <n-layout-content content-style="padding: 4px 0 0 24px;">
                    <n-form :label-width="100" :model="dataForm" label-placement="left" require-mark-placement="left" :rules="dataRules" ref="dataRef">
                        <n-grid responsive="screen" cols="14 s:14 m:14 l:24 xl:24 xxl:24" :x-gap="6" :y-gap="12">
                            <n-form-item-gi :show-feedback="false" :span="4">
                                <n-input type="text" placeholder="搜索" clearable v-model:value="mult">
                                    <template #suffix>
                                        <n-icon id="search" @click="init(cateActive)">
                                            <SearchOutline></SearchOutline>
                                        </n-icon>
                                    </template>
                                </n-input>
                            </n-form-item-gi>
                            <!-- <n-grid-item span="12"> </n-grid-item> -->
                            <n-form-item-gi :show-feedback="false" :span="10" path="title"> <n-input v-model:value="dataForm.title" placeholder="请输入标题" /> </n-form-item-gi>
                            <n-form-item-gi :show-feedback="false" :span="5">
                                <n-select placeholder="标签" max-tag-count="responsive" style="width: 100%" multiple remote v-model:value="dataForm.tags" :options="tagAll" />
                            </n-form-item-gi>
                            <n-grid-item :show-feedback="false" span="24">
                                <Note ref="mdRef" v-model:value="dataForm.markdown" @change="onChange" @save="onSave" />
                            </n-grid-item>
                            <!-- <n-grid-item span="24"> </n-grid-item> -->
                        </n-grid>
                    </n-form>
                </n-layout-content>
            </n-layout>
            <n-dropdown :options="dropdowns" placement="top" @select="onSelect">
                <Floating />
            </n-dropdown>
            <!-- <Cates :class="{ 'cate-active': cateActive == 0 }" :index="0" @click="onCateActive(0)">所有</Cates> -->
            <Cates v-for="(item, index) in cateAll" :class="{ 'cate-active': cateActive == item.id }" :index="index + 1" :color="item.color" @click="onCateActive(item.id)" :tip="item.name"> {{ item.name }}</Cates>
        </n-layout-content>
    </n-layout>
</template>
<script lang="ts" setup>
import { apiPostSave, apiPostList, apiTagList, apiPostGet, apiCateList } from "@/api";
import { Note } from "@/components/Editor";
import { SearchOutline } from "@vicons/ionicons5";
import Floating from "./floating.vue";
import { useRouter } from "vue-router";
import Card from "./card.vue";
import Cates from "./cate.vue";
import Menu from "./menu.vue";
import { ref, onMounted, nextTick } from "vue";
import { useMessage, NIcon } from "naive-ui";
const router = useRouter();
interface Note {
    id: number;
    title: string;
    markdown: string;
}
interface Cate {
    id: number;
    name: string;
    color: string;
    intro: string;
}
const dropdowns = [
    { label: "新建", key: "new" },
    { label: "返回", key: "back" },
    // { label: "删除", key: "drop" },
];
const noteActive = ref(0);
const cateActive = ref(0);
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
const cateAll = ref<Cate[]>([]);
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
        cate_id: cateActive.value,
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
const onCateActive = (id) => {
    if (cateActive.value === id) {
        return;
    }
    cateActive.value = id;
    init(id);
};
const onActive = (id) => {
    if (noteActive.value === id) {
        return;
    }
    noteActive.value = id;
    detail();
};
const onSelect = (key: string) => {
    console.log(key);
    switch (key) {
        case "new":
            create();
            break;
        case "back":
            router.push({ name: "post-index" });
            break;
    }
};
const detail = () => {
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
const init = (cateId = 0) => {
    apiPostList({ kind: 3, mult: mult.value, cate_id: cateId }).then((resp) => {
        if (resp.code == 200) {
            noteAll.value = resp.data;
            noteActive.value = noteAll.value[0].id;
            detail();
        } else {
            noteAll.value = [];
        }
    });
};
const preInit = () => {
    apiCateList({ kind: 3 }).then((resp) => {
        if (resp.code == 200) {
            cateAll.value = resp.data;
        } else {
            cateAll.value = [];
        }
    });
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
                    init(cateActive.value);
                    message.success("保存成功");
                } else {
                    message.error(resp.msg);
                }
            });
        }
    });
};
</script>
