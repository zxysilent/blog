<template>
    <n-card :bordered="false" title="修改分类" class="mt-4 proCard" size="small" :segmented="{ content: 'hard' }">
        <n-grid cols="1 s:2 m:2 l:3 xl:3 xxl:3" responsive="screen">
            <n-grid-item>
                <n-form :label-width="100" :model="dataForm" label-placement="left" require-mark-placement="left" :rules="dataRules" ref="dataRef">
                    <n-form-item label="分类名称" path="name">
                        <n-input v-model:value="dataForm.name" />
                    </n-form-item>
                    <n-form-item label="分类颜色" path="color">
                        <n-color-picker v-model:value="dataForm.color" :modes="['hex']" :show-alpha="false" />
                    </n-form-item>
                    <n-form-item label="分类描述" path="intro">
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
                            <n-button type="warning" :loading="loading" @click="emitEdit">确认修改</n-button>
                            <n-button type="success" @click="emitReset">重置填写</n-button>
                            <n-button type="tertiary" @click="$router.push({ name: 'cate-index' })">取消返回</n-button>
                        </n-space>
                    </div>
                </n-form>
            </n-grid-item>
        </n-grid>
    </n-card>
</template>
<script lang="ts" setup>
import { apiCateGet, apiCateEdit } from "@/api";
import { ref, onMounted } from "vue";
import { useMessage } from "naive-ui";
import { useRoute } from "vue-router";
const route = useRoute();
const dataForm = ref({
    id: 0,
    name: "",
    color: "#333639",
    intro: "",
});

const dataRules = {
    name: { required: true, message: "请输入分类名称", trigger: "blur" },
    intro: { required: true, message: "请输入分类描述", trigger: "blur" },
};
const dataRef = ref();
const message = useMessage();
const loading = ref(false);
const init = () => {
    apiCateGet({ id: route.params.id }).then((resp) => {
        if (resp.code == 200) {
            dataForm.value = resp.data;
        } else {
            dataForm.value = {
                id: 0,
                name: "",
                color: "color",
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
            apiCateEdit(dataForm.value).then((resp) => {
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
