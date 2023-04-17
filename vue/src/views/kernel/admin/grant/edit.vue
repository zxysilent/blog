<template>
    <n-card :bordered="false" title="修改权限" class="mt-4 proCard" size="small" :segmented="{ content: 'hard' }">
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
                            <n-button type="warning" :loading="loading" @click="emitEdit">确认修改</n-button>
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
import { apiGrantGet, apiGrantEdit } from "@/api";
import { ref, onMounted } from "vue";
import { useMessage } from "naive-ui";
import { useRoute } from "vue-router";
const route = useRoute();
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
            };
        }
    });
};
onMounted(() => {
    init();
});
const emitReset = () => {
    init();
};
const emitEdit = () => {
    console.log(dataForm);
    dataRef.value.validate((errors) => {
        if (!errors) {
            apiGrantEdit(dataForm.value).then((resp) => {
                if (resp.code == 200) {
                    message.success("修改成功", {
                        onAfterLeave() {
                            init();
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
