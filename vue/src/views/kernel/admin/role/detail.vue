<template>
    <n-card :bordered="false" title="角色详情" size="small" :segmented="{ content: 'hard' }">
        <n-grid cols="1 s:2 m:2 l:3 xl:3 xxl:3" responsive="screen">
            <n-grid-item>
                <n-form :label-width="100" :model="dataForm" disabled label-placement="left" require-mark-placement="left" ref="dataRef">
                    <n-form-item label="角色名称" path="name">
                        <n-input v-model:value="dataForm.name" />
                    </n-form-item>
                    <n-form-item label="角色描述" path="intro">
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
                            <n-button type="tertiary" @click="$router.push({ name: 'role-index' })">取消返回</n-button>
                        </n-space>
                    </div>
                </n-form>
            </n-grid-item>
        </n-grid>
    </n-card>
</template>
<script lang="ts" setup>
import { apiRoleGet } from "@/api";
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
const route = useRoute();
const dataForm = ref({
    id: 0,
    name: "",
    intro: "",
    updated: 0,
    created: 0,
});
const dataRef = ref();
const init = () => {
    apiRoleGet({ id: route.params.id }).then((resp) => {
        if (resp.code == 200) {
            dataForm.value = resp.data;
        } else {
            dataForm.value = {
                id: 0,
                name: "",
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
