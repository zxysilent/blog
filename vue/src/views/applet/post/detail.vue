<template>
    <n-card :bordered="false" class="mt-4 proCard" size="small" :segmented="{ content: 'hard' }">
        <n-form disabled :label-width="100" :model="dataForm" label-placement="left" require-mark-placement="left" ref="dataRef">
            <n-grid responsive="screen" cols="4 s:8 m:12 l:16 xl:20 xxl:24" :x-gap="4" :y-gap="12">
                <!-- <n-grid-item span="12"> </n-grid-item> -->
                <n-form-item-gi :show-feedback="false" :span="12" path="title"> <n-input v-model:value="dataForm.title" placeholder="请输入标题" /> </n-form-item-gi>
                <n-form-item-gi :show-feedback="false" :span="12" path="path">
                    <n-input v-model:value="dataForm.path" placeholder="请输入访问路径">
                        <template #prefix> https://www. </template>
                        <template #suffix> .html </template>
                    </n-input>
                </n-form-item-gi>
                <n-form-item-gi :show-feedback="false" :span="1"> <n-checkbox v-model:checked="dataForm.allow"> 评论 </n-checkbox> </n-form-item-gi>
                <n-form-item-gi :show-feedback="false" :span="1"> <n-checkbox v-model:checked="dataForm.status" :checked-value="1" :unchecked-value="2"> 草稿 </n-checkbox> </n-form-item-gi>
                <n-form-item-gi :show-feedback="false" :span="3"> <n-date-picker style="width: 100%" v-model:value="dataForm.created" type="datetime" /> </n-form-item-gi>
                <n-form-item-gi :show-feedback="false" :span="1">
                    <n-switch size="large" :round="false" v-model:value="dataForm.kind" :checked-value="1" :unchecked-value="2"> <template #checked>文章 </template> <template #unchecked> 页面 </template> </n-switch>
                </n-form-item-gi>
                <n-form-item-gi :show-feedback="false" :span="2">
                    <n-select placeholder="分类" max-tag-count="responsive" style="width: 100%" remote v-model:value="dataForm.cate_id" :options="cateAll" />
                </n-form-item-gi>
                <n-form-item-gi :show-feedback="false" :span="5">
                    <n-select placeholder="标签" max-tag-count="responsive" style="width: 100%" multiple remote v-model:value="dataForm.tags" :options="tagAll" />
                </n-form-item-gi>
                <n-grid-item span="6">
                    <n-space>
                        <n-button type="success" @click="init">信息刷新</n-button>
                        <n-button type="tertiary" @click="$router.push({ name: 'post-index' })">取消返回</n-button>
                    </n-space>
                </n-grid-item>
                <n-grid-item :show-feedback="false" span="24"> <Markdown v-model:value="dataForm.markdown" readonly :subfield="false" navigation :defaultOpen="'preview'" /> </n-grid-item>
                <!-- <n-grid-item span="24"> </n-grid-item> -->
            </n-grid>
        </n-form>
    </n-card>
</template>
<script lang="ts" setup>
import { apiPostGet, apiCateList, apiTagList } from "@/api";
import { ref, onMounted } from "vue";
import { Markdown } from "@/components/Editor";
import { useRoute } from "vue-router";
const route = useRoute();
const cateAll = ref([]);
const tagAll = ref([]);
const dataForm = ref({
    id: 0,
    title: "",
    path: "",
    kind: 1,
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

const dataRef = ref();
const init = () => {
    apiPostGet({ id: route.params.id }).then((resp) => {
        if (resp.code == 200) {
            dataForm.value = resp.data;
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
const preInit = () => {
    apiCateList({}).then((resp) => {
        if (resp.code == 200) {
            cateAll.value = resp.data.map((item) => {
                return {
                    label: item.name,
                    value: item.id,
                };
            });
            dataForm.value.cate_id = resp.data[0].id;
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

onMounted(() => {
    preInit();
    init();
});
</script>
