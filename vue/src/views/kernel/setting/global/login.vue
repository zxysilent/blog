<template>
    <n-grid cols="1 s:1 m:1 l:2 xl:3 xxl:3" responsive="screen">
        <n-grid-item>
            <n-form :label-width="120" :model="dataValue" label-placement="left" require-mark-placement="left" :rules="dataRules" ref="dataRef" :disabled="disable">
                <n-form-item label="标题" path="title">
                    <n-input v-model:value="dataValue.title" />
                </n-form-item>
                <n-form-item label="描述" path="desc">
                    <n-input v-model:value="dataValue.desc" />
                </n-form-item>
                <n-form-item label="API">
                    <n-space>
                        <n-button text tag="a" :href="'/api/dict/get/' + KEY" target="_blank" type="primary"> 标准查看 </n-button>
                        <n-button text tag="a" :href="'/api/dict/' + KEY" target="_blank" type="primary"> 应用查看 </n-button>
                    </n-space>
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
const KEY = "login";
interface Value {
    title: string;
    desc: string;
}
interface Applet {
    key: string;
    updated: number;
    created: number;
}
const dataValue = ref<Value>({ title: "欢迎登录", desc: "便捷、高效、可靠" });
const dataForm = ref<Applet>({ key: KEY, updated: 0, created: 0 });
const disable = ref<boolean>(true);
const dataRules = {
    title: { required: true, message: "请输入标题", trigger: "blur" },
    desc: { required: true, message: "请输入提示", trigger: "blur" },
};
const dataRef = ref();
const message = useMessage();
const loading = ref(false);
const init = () => {
    apiDictGet(KEY).then((resp) => {
        if (resp.code == 200) {
            dataForm.value.updated = resp.data.updated;
            dataForm.value.created = resp.data.created;
            dataValue.value = JSON.parse(resp.data.value);
        } else {
            dataValue.value = { title: "欢迎登录", desc: "便捷、高效、可靠" };
            dataForm.value.updated = 0;
            dataForm.value.created = 0;
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
