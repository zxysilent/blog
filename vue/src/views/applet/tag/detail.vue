<template>
    <n-card :bordered="false" title="标签详情" class="mt-4 proCard" size="small" :segmented="{ content: 'hard' }">
        <n-grid cols="1 s:2 m:2 l:3 xl:3 xxl:3" responsive="screen">
            <n-grid-item>
                <n-form :label-width="100" :model="dataForm" disabled label-placement="left" require-mark-placement="left" :rules="dataRules" ref="dataRef">
                    <n-form-item label="标签名称" path="name">
                        <n-input v-model:value="dataForm.name" />
                    </n-form-item>
                    <n-form-item label="标签颜色" path="color">
                        <n-color-picker v-model:value="dataForm.color" :modes="['hex']" :show-alpha="false" />
                    </n-form-item>
                    <n-form-item label="标签描述" path="intro">
                        <n-input
                            v-model:value="dataForm.intro"
                            type="textarea"
                            :autosize="{
                                minRows: 3,
                                maxRows: 5,
                            }" />
                    </n-form-item>
                    <n-form-item label="更新时间" :show-feedback="false">
                        <n-time :time="dataForm.updated" />
                    </n-form-item>
                    <n-form-item label="创建时间" :show-feedback="false">
                        <n-time :time="dataForm.created" />
                    </n-form-item>
                    <div>
                        <n-space justify="center">
                            <n-button type="success" @click="init">信息刷新</n-button>
                            <n-button type="tertiary" @click="$router.push({ name: 'tag-index' })">取消返回</n-button>
                        </n-space>
                    </div>
                </n-form>
            </n-grid-item>
        </n-grid>
    </n-card>
</template>
<script lang="ts" setup>
import { apiTagGet } from "@/api";
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
const route = useRoute();
const dataForm = ref({
    id: 0,
    name: "",
    color: "#333639",
    intro: "",
    updated: 0,
    created: 0,
});

const dataRules = {
    name: { required: true, message: "请输入标签名称", trigger: "blur" },
    intro: { required: true, message: "请输入标签描述", trigger: "blur" },
};
const dataRef = ref();
const init = () => {
    apiTagGet({ id: route.params.id }).then((resp) => {
        if (resp.code == 200) {
            dataForm.value = resp.data;
        } else {
            dataForm.value = {
                id: 0,
                name: "",
                color: "",
                intro: "",
                updated: 0,
                created: 0,
            };
        }
    });
};
onMounted(() => {
    init();
});
</script>
