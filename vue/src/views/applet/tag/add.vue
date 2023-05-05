<template>
    <n-card :bordered="false" title="添加标签" size="small" :segmented="{ content: 'hard' }">
        <n-grid cols="1 s:2 m:2 l:3 xl:3 xxl:3" responsive="screen">
            <n-grid-item>
                <n-form :label-width="100" :model="dataForm" label-placement="left" require-mark-placement="left" :rules="dataRules" ref="dataRef">
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
                    <div>
                        <n-space justify="center">
                            <n-button type="warning" :loading="loading" @click="emitAdd">确认添加</n-button>
                            <n-button type="success" @click="emitReset">重置填写</n-button>
                            <n-button type="tertiary" @click="$router.push({ name: 'tag-index' })">取消返回</n-button>
                        </n-space>
                    </div>
                </n-form>
            </n-grid-item>
        </n-grid>
    </n-card>
</template>
<script lang="ts" setup>
import { apiTagAdd } from "@/api";
import { ref, onMounted } from "vue";
import { useMessage } from "naive-ui";
const dataForm = ref({
    id: 0,
    name: "",
    color: "#333639",
    intro: "",
});
const dataRules = {
    name: { required: true, message: "请输入标签名称", trigger: "blur" },
    intro: { required: true, message: "请输入标签描述", trigger: "blur" },
};
const dataRef = ref();
const message = useMessage();
const loading = ref(false);

onMounted(() => {});
const emitReset = () => {
    dataForm.value = {
        id: 0,
        name: "",
        color: "#333639",
        intro: "",
    };
    dataRef.value.restoreValidation();
};
const emitAdd = () => {
    console.log(dataForm.value);
    dataRef.value.validate((errors) => {
        if (!errors) {
            apiTagAdd(dataForm.value).then((resp) => {
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
