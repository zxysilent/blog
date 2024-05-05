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
        <!-- <Sidebar></Sidebar> -->
        <n-layout-sider collapsed position="static" collapse-mode="width" :collapsed-width="64" inverted class="layout-sider">
            <n-menu inverted :collapsed-width="64" :collapsed-icon-size="22" :options="menuOptions" @update:value="onUpdate" :value="activeMenu" />
        </n-layout-sider>
        <n-layout-content>
            <n-layout has-sider>
                <n-layout-sider
                    collapse-mode="transform"
                    :collapsed-width="0"
                    :width="220"
                    show-trigger="bar"
                    :inverted="inverted"
                    :collapsed="collapsed"
                    @collapse="collapsed = true"
                    @expand="collapsed = false"
                    style="height: 100vh">
                    <n-scrollbar style="border-right: 1px solid #f2f6fc">
                        <Card label-field="name" v-if="memCate.id > 0" key-field="id" @click="onFolderBack">
                            <n-icon size="20">
                                <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 24 24">
                                    <path d="M9 11l-4 4l4 4m-4-4h11a4 4 0 0 0 0-8h-1" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path>
                                </svg>
                            </n-icon>
                            返回 {{ memCate?.name }}
                        </Card>
                        <Card v-for="item in noted?.cates" @click="onFolder(item)" label-field="name" key-field="id">
                            <n-icon size="20" color="#FFBA49">
                                <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 16 16">
                                    <g fill="none">
                                        <path d="M7.207 4h-.032l-1.113-.89A.5.5 0 0 0 5.75 3H4a2 2 0 0 0-2 2v.5h3.557L7.207 4z" fill="currentColor"></path>
                                        <path d="M8.693 4L6.086 6.37a.5.5 0 0 1-.336.13H2V11a2 2 0 0 0 2 2h8a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2H8.693z" fill="currentColor"></path>
                                    </g>
                                </svg>
                            </n-icon>
                            {{ item.name }}
                        </Card>
                        <Card v-for="item in noted?.posts" @click="onActive(item.id)" :class="{ active: noteActive == item.id }">
                            <n-icon size="20" color="#38B7FC">
                                <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 1024 1024">
                                    <path
                                        d="M854.6 288.6L639.4 73.4c-6-6-14.1-9.4-22.6-9.4H192c-17.7 0-32 14.3-32 32v832c0 17.7 14.3 32 32 32h640c17.7 0 32-14.3 32-32V311.3c0-8.5-3.4-16.7-9.4-22.7zM790.2 326H602V137.8L790.2 326zm1.8 562H232V136h302v216a42 42 0 0 0 42 42h216v494zM429 481.2c-1.9-4.4-6.2-7.2-11-7.2h-35c-6.6 0-12 5.4-12 12v272c0 6.6 5.4 12 12 12h27.1c6.6 0 12-5.4 12-12V582.1l66.8 150.2a12 12 0 0 0 11 7.1H524c4.7 0 9-2.8 11-7.1l66.8-150.6V758c0 6.6 5.4 12 12 12H641c6.6 0 12-5.4 12-12V486c0-6.6-5.4-12-12-12h-34.7c-4.8 0-9.1 2.8-11 7.2l-83.1 191l-83.2-191z"
                                        fill="currentColor"></path>
                                </svg>
                            </n-icon>
                            {{ item.title }}
                            <br />
                            <n-time type="date" style="font-size: small" :time="item.created" />
                        </Card>

                        <n-tree
                        :data="noted.posts"
    block-line
    :render-label="renderLabel"
    :render-suffix="renderSuffix"
    :selectable="false"
  />
                    </n-scrollbar>
                </n-layout-sider>
                <n-layout-content content-style="padding: 4px 0 0 24px;">
                    <n-form :label-width="100" :model="dataForm" label-placement="left" require-mark-placement="left" :rules="dataRules" ref="dataRef">
                        <n-grid responsive="screen" cols="14 s:14 m:14 l:24 xl:24 xxl:24" :x-gap="6" :y-gap="12">
                            <n-form-item-gi :show-feedback="false" :span="4">
                                <n-input type="text" placeholder="搜索" clearable v-model:value="mult">
                                    <template #suffix>
                                        <n-icon id="search" @click="init(memCate.id)">
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
        </n-layout-content>
    </n-layout>
</template>
<script lang="ts" setup>
import { apiNoteSave, apiNoteList, apiTagList, apiPostGet, apiCateGet } from "@/api";
import Card from "./card.vue";
import Sidebar from "./sidebar.vue";
import { defineComponent, } from 'vue'
import { NButton, TreeOption } from 'naive-ui'
import { repeat } from 'seemly'
import { Note } from "@/components/Editor";
import { SearchOutline } from "@vicons/ionicons5";
import { RouterLink } from "vue-router";
import { ref, onMounted, nextTick, h, Component } from "vue";
import { useMessage, NIcon } from "naive-ui";
import type { MenuOption } from "naive-ui";
import { FolderOutline, AddOutline, PricetagsOutline, GridOutline, ArrowBackCircleOutline } from "@vicons/ionicons5";
function renderIcon(icon: Component) {
    return () => h(NIcon, null, { default: () => h(icon) });
}
const activeMenu = ref("newest");
interface Cate {
    id: number;
    pid: number;
    name: string;
    color: string;
    intro: string;
    created: number;
}
interface Post {
    id: number;
    title: string;
    markdown: string;
    created: number;
}
interface Note {
    empty: boolean;
    cates: Array<Cate>;
    posts: Array<Post>;
}

function renderLabel({ option }: { option: TreeOption }) {
    return `${option.title}`
}

function renderSuffix({ option }: { option: TreeOption }) {
    return h(
        NButton,
        { text: true, type: 'primary' },
        { default: () => `S` }
    )
}
const memCate = ref<Cate>({ id: 0, pid: 0, name: "", color: "", intro: "", created: 0 });
const onFolder = (cate: Cate) => {
    memCate.value = cate;
    init(cate.id, false);
};
const onFolderBack = () => {
    init(memCate.value?.pid, false);
    if (memCate.value.pid == 0) {
        memCate.value = { id: 0, pid: 0, name: "", color: "", intro: "", created: 0 };
        return;
    }
    cateGet(memCate.value.pid);
};
const noted = ref<Note>({
    empty: true,
    posts: [],
    cates: [],
});
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
        cate_id: memCate.value.id,
        allow: false,
        richtext: "",
        markdown: "",
        tags: [],
        updated: 0,
        created: 0,
    };
    noteActive.value = 0;
};
const menuOptions: MenuOption[] = [
    { label: "新建", key: "add", icon: renderIcon(AddOutline) },
    { label: "最新", key: "newest", icon: renderIcon(GridOutline) },
    { label: "目录", key: "folder", icon: renderIcon(FolderOutline) },
    // { label: "标签", key: "tags", icon: renderIcon(PricetagsOutline) },
    {
        label: () => h(RouterLink, { to: { name: "post-index" } }, { default: () => "返回" }),
        key: "back",
        icon: renderIcon(ArrowBackCircleOutline),
    },
];
const onUpdate = (val) => {
    console.log(val);
    switch (val) {
        case "add": {
            activeMenu.value = "none";
            create();
            break;
        }
        case "newest": {
            activeMenu.value = "newest";
            memCate.value = { id: 0, pid: 0, name: "", color: "", intro: "", created: 0 };
            init(0, true);
            break;
        }
        case "folder": {
            activeMenu.value = "folder";
            preInit();
            memCate.value = { id: 0, pid: 0, name: "", color: "", intro: "", created: 0 };
            init(0, false);
            break;
        }
        case "tags": {
            activeMenu.value = "tags";
            break;
        }
    }
};
const onChange = (val) => {
    console.log(val);
    dataForm.value.richtext = val;
};
const onActive = (id) => {
    if (noteActive.value === id) {
        return;
    }
    noteActive.value = id;
    noteGet();
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
const cateGet = (cateId: number) => {
    apiCateGet({ id: cateId }).then((resp) => {
        if (resp.code == 200) {
            memCate.value = resp.data;
        } else {
            memCate.value = { id: 0, pid: 0, name: "", color: "", intro: "", created: 0 };
        }
    });
};
const init = (cateId = 0, newest = true) => {
    apiNoteList({ mult: mult.value, cate_id: cateId, newest: newest }).then((resp) => {
        if (resp.code == 200) {
            if (noted.value?.posts.length) {
                //@ts-ignore
                noteActive.value = noted.value?.posts[0].id;
                noteGet();
            }
            noted.value = resp.data;
        } else {
            noted.value = { empty: true, posts: [], cates: [] };
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
    init(0, true);
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
