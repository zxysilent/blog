<template>
    <n-card :bordered="false" title="添加配置" size="small" :segmented="{ content: 'hard' }">
        <n-grid cols="1 s:2 m:2 l:3 xl:3 xxl:3" responsive="screen">
            <n-grid-item>
                <n-form :label-width="100" :model="dataForm" label-placement="left" require-mark-placement="left" :rules="dataRules" ref="dataRef">
                    <n-form-item label="配置名" path="key">
                        <n-input v-model:value="dataForm.key" maxlength="64" show-count>
                            <template #suffix>不可修改</template>
                        </n-input>
                    </n-form-item>
                    <n-form-item label="数据类型" path="kind">
                        <n-radio-group v-model:value="dataForm.kind" name="radiogroup">
                            <n-space>
                                <n-radio v-for="kind in kinds" :key="kind" :value="kind">
                                    {{ kind }}
                                </n-radio>
                            </n-space>
                        </n-radio-group>
                    </n-form-item>
                    <n-form-item label="配置值" path="value">
                        <n-input type="textarea" :autosize="{ minRows: 5, maxRows: 8 }" v-model:value="dataForm.value" maxlength="255" clearable show-count />
                    </n-form-item>
                    <n-form-item label="说明" path="intro">
                        <n-input v-model:value="dataForm.intro" />
                    </n-form-item>
                    <div>
                        <n-space justify="center">
                            <n-button type="warning" :loading="loading" @click="emitAdd">确认添加</n-button>
                            <n-button type="success" @click="emitReset">重置填写</n-button>
                            <n-button type="tertiary" @click="$router.push({ name: 'dict-index' })">取消返回</n-button>
                        </n-space>
                    </div>
                </n-form>
            </n-grid-item>
        </n-grid>
    </n-card>
</template>
<script lang="ts" setup>
import { apiDictAdd } from "@/api";
import { apiDictGet } from "@/api/ext";
import { reactive, ref, onMounted } from "vue";
import { useMessage } from "naive-ui";

const kinds = ref<string[]>(["string", "time", "bool", "int", "json"]);
const dataForm = reactive({
    key: "",
    value: "",
    kind: "string",
    intro: "",
});
const valideKey = (_, value, callback) => {
    if (value == "") {
        return;
    }
    apiDictGet(value).then((resp) => {
        if (resp.code == 200) {
            callback(new Error("当前配置名已经存在"));
        } else {
            callback();
        }
    });
};
const dataRules = {
    key: [
        { required: true, message: "请输入配置名", trigger: "blur" },
        { min: 1, message: "请至少输入1个字符", trigger: "blur" },
        { max: 64, message: "最多输入64个字符", trigger: "blur" },
        { validator: valideKey, trigger: "blur" },
    ],
    value: { required: true, message: "请输入配置值", trigger: "blur" },
    // intro: { required: true, message: "请输入配置说明", trigger: "blur" },
    kind: { required: true, type: "string", message: "请选择数据类型", trigger: ["blur", "change"] },
};
const dataRef = ref();
const message = useMessage();
const loading = ref(false);

onMounted(() => {});
const emitReset = () => {
    dataForm.key = "";
    dataForm.value = "";
    dataForm.kind = "string";
    dataForm.intro = "";
    dataRef.value.restoreValidation();
};
const emitAdd = () => {
    console.log(dataForm);
    loading.value = true;
    dataRef.value.validate((errors) => {
        loading.value = false;
        if (!errors) {
            apiDictAdd(dataForm).then((resp) => {
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
