<template>
    <n-card :bordered="false" title="权限详情" class="mt-4 proCard" size="small" :segmented="{ content: 'hard' }">
        <n-grid cols="1 s:2 m:2 l:3 xl:3 xxl:3" responsive="screen">
            <n-grid-item>
                <n-form :label-width="100" :model="dataForm" disabled label-placement="left" require-mark-placement="left" ref="dataRef">
                    <n-form-item label="权限名称" path="name">
                        <n-input v-model:value="dataForm.name" />
                    </n-form-item>
                    <n-form-item label="权限标识" path="guid">
                        <n-input v-model:value="dataForm.guid" />
                    </n-form-item>
                    <n-form-item label="权限分组" path="group">
                        <n-input v-model:value="dataForm.group" />
                    </n-form-item>
                    <n-form-item label="权限排序" path="sort">
                        <n-input-number v-model:value="dataForm.sort" :min="0" :max="2000" />
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
                            <n-button type="tertiary" @click="$router.push({ name: 'grant-index' })">取消返回</n-button>
                        </n-space>
                    </div>
                </n-form>
            </n-grid-item>
        </n-grid>
    </n-card>
</template>
<script lang="ts" setup>
import { apiGrantGet } from "@/api";
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
const route = useRoute();
const dataForm = ref({
    id: 0,
    name: "",
    guid: "",
    group: "",
    sort: 1000,
    updated: 0,
    created: 0,
});

const dataRef = ref();
const init = () => {
    apiGrantGet({ id: route.params.id }).then((resp) => {
        if (resp.code == 200) {
            dataForm.value = resp.data;
        } else {
            dataForm.value = {
                id: 0,
                name: "",
                guid: "",
                group: "",
                sort: 1000,
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
