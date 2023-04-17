<template>
    <n-grid cols="1 s:1 m:1 l:2 xl:3 xxl:3" responsive="screen">
        <n-grid-item>
            <n-form :label-width="120" :model="dataValue" label-placement="left" require-mark-placement="left" :rules="dataRules" ref="dataRef" :disabled="disable">
                <n-form-item label="API">
                    <n-button text tag="a" href="/api/global/get" target="_blank" type="primary"> 查看 </n-button>
                </n-form-item>
                <n-form-item label="更新时间" :show-feedback="false">
                    <n-time :time="dataForm.updated" />
                </n-form-item>
                <n-form-item label="创建时间" :show-feedback="false">
                    <n-time :time="dataForm.created" />
                </n-form-item>
                <div>
                    <n-space justify="center">
                        <n-button type="warning" :disabled="disable" :loading="loading" @click="emitEdit">确认修改</n-button>
                        <n-button type="success" @click="emitReset">取消刷新</n-button>
                        <n-button type="tertiary" :disabled="!disable" @click="openEdit">开启编辑</n-button>
                    </n-space>
                </div>
            </n-form>
        </n-grid-item>
    </n-grid>
</template>
<script lang="ts" setup>
import { apiDictEdit } from "@/api";
import { apiDictGet } from "@/api/ext";
import { ref, onMounted } from "vue";
import { useMessage } from "naive-ui";
const KEY = "applet";
interface Value {
    data: string;
}
interface Applet {
    key: string;
    value: Value;
    updated: number;
    created: number;
}

const dataValue = ref<Value>({ data: "" });
const dataForm = ref<Applet>({ key: KEY, value: dataValue.value, updated: 0, created: 0 });
const disable = ref<boolean>(true);
const dataRules = {
    data: { type: "number", required: true, message: "请输入轮播广告高度", trigger: "blur", min: 1 },
    number: { type: "number", required: true, message: "请输广告小标宽度", trigger: "blur", min: 1 },
};
const dataRef = ref();
const message = useMessage();
const loading = ref(false);
const init = () => {
    apiDictGet(KEY).then((resp) => {
        if (resp.code == 200) {
            dataForm.value.key = KEY;
            dataForm.value.updated = resp.data.updated;
            dataForm.value.created = resp.data.created;
            dataForm.value.value = JSON.parse(resp.data.value);
        } else {
            dataValue.value = { data: "" };
            dataForm.value = {
                key: KEY,
                value: dataValue.value,
                updated: 0,
                created: 0,
            };
        }
    });
};
onMounted(() => {
    init();
    disable.value = true;
});
const emitReset = () => {
    init();
    disable.value = true;
};
const openEdit = () => {
    disable.value = false;
};
const emitEdit = () => {
    console.log(dataForm);
    dataRef.value.validate((errors) => {
        if (!errors) {
            const data = {
                key: KEY,
                value: JSON.stringify(dataValue.value),
            };
            apiDictEdit(data).then((resp) => {
                disable.value = true;
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
