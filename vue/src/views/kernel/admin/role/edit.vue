<template>
    <n-card :bordered="false" title="修改角色" size="small" :segmented="{ content: 'hard' }">
        <n-grid cols="1 s:2 m:2 l:3 xl:3 xxl:3" responsive="screen">
            <n-grid-item>
                <n-form :label-width="100" :model="dataForm" label-placement="left" require-mark-placement="left" :rules="dataRules" ref="dataRef">
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
                            }"
                        />
                    </n-form-item>
                    <div>
                        <n-space justify="center">
                            <n-button type="warning" :loading="loading" @click="emitEdit">确认修改</n-button>
                            <n-button type="success" @click="emitReset">重置填写</n-button>
                            <n-button type="tertiary" @click="$router.push({ name: 'role-index' })">取消返回</n-button>
                        </n-space>
                    </div>
                </n-form>
            </n-grid-item>
        </n-grid>
    </n-card>
</template>
<script lang="ts" setup>
import { apiRoleGet, apiRoleEdit } from "@/api";
import { ref, onMounted } from "vue";
import { useMessage } from "naive-ui";
import { useRoute } from "vue-router";
const route = useRoute();
const dataForm = ref({
    id: 0,
    name: "",
    intro: "",
});

const dataRules = {
    name: { required: true, message: "请输入角色名称", trigger: "blur" },
    intro: { required: true, message: "请输入角色描述", trigger: "blur" },
};
const dataRef = ref();
const message = useMessage();
const loading = ref(false);
const init = () => {
    apiRoleGet({ id: route.params.id }).then((resp) => {
        if (resp.code == 200) {
            dataForm.value = resp.data;
        } else {
            dataForm.value = {
                id: 0,
                name: "",
                intro: "",
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
            apiRoleEdit(dataForm.value).then((resp) => {
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
