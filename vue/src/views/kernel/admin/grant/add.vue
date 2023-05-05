<template>
    <n-card :bordered="false" title="添加权限" size="small" :segmented="{ content: 'hard' }">
        <n-grid cols="1 s:2 m:2 l:3 xl:3 xxl:3" responsive="screen">
            <n-grid-item>
                <n-form :label-width="100" :model="dataForm" label-placement="left" require-mark-placement="left" :rules="dataRules" ref="dataRef">
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
                    <div>
                        <n-space justify="center">
                            <n-button type="warning" :loading="loading" @click="emitAdd">确认添加</n-button>
                            <n-button type="success" @click="emitReset">重置填写</n-button>
                            <n-button type="tertiary" @click="$router.push({ name: 'grant-index' })">取消返回</n-button>
                        </n-space>
                    </div>
                </n-form>
            </n-grid-item>
        </n-grid>
    </n-card>
</template>
<script lang="ts" setup>
import { apiGrantAdd } from "@/api";
import { ref, onMounted } from "vue";
import { useMessage } from "naive-ui";
const dataForm = ref({
    id: 0,
    name: "",
    guid: "",
    group: "",
    sort: 1000,
});
const dataRules = {
    name: { required: true, message: "请输入权限名称", trigger: "blur" },
    guid: { required: true, message: "请输入权限标识", trigger: "blur" },
    group: { required: true, message: "请输入权限分组", trigger: "blur" },
};
const dataRef = ref();
const message = useMessage();
const loading = ref(false);

onMounted(() => {});
const emitReset = () => {
    dataForm.value = {
        id: 0,
        name: "",
        guid: "",
        group: "",
        sort: 1000,
    };
    dataRef.value.restoreValidation();
};
const emitAdd = () => {
    console.log(dataForm.value);
    dataRef.value.validate((errors) => {
        if (!errors) {
            apiGrantAdd(dataForm).then((resp) => {
                if (resp.code == 200) {
                    message.success("添加成功", {
                        onAfterLeave() {
                            emitReset();
                        },
                    });
                } else {
                    message.error(resp.msg);
                }
            });
        }
    });
};
</script>
