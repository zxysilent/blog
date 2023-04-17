<template>
    <n-grid cols="1 s:1 m:1 l:2 xl:3 xxl:3" responsive="screen">
        <n-grid-item>
            <n-form :label-width="120" :model="dataValue" label-placement="left" require-mark-placement="left" :rules="dataRules" ref="dataRef" :disabled="disable">
                <n-form-item label="githu地址" path="github_url">
                    <n-input v-model:value="dataValue.github_url" maxlength="64" show-count placeholder="请填写githu地址"></n-input>
                </n-form-item>
                <n-form-item label="微博地址" path="weibo_url">
                    <n-input v-model:value="dataValue.weibo_url" maxlength="64" show-count placeholder="请填写微博地址"></n-input>
                </n-form-item>
                <n-form-item label="分页大小" path="page_size">
                    <n-slider v-model:value="dataValue.page_size" :marks="marks" :step="1" :min="1" :max="10" show-stops></n-slider>
                </n-form-item>
                <n-form-item label="统计配置" path="analytic">
                    <n-input type="textarea" :autosize="{ minRows: 5, maxRows: 8 }" v-model:value="dataValue.analytic" maxlength="512" show-count placeholder="请填写自定义js"></n-input>
                </n-form-item>
                <n-form-item label="评论配置" path="comment">
                    <n-input type="textarea" :autosize="{ minRows: 5, maxRows: 8 }" v-model:value="dataValue.comment" maxlength="512" show-count placeholder="请填写自定义js"></n-input>
                </n-form-item>
                <n-form-item label="API" :show-feedback="false">
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
const KEY = "blog";
const marks = { 1: "1", 2: "2", 3: "3", 4: "4", 5: "5", 6: "6", 7: "7", 8: "8", 9: "9", 10: "10" };
interface Value {
    github_url: string;
    page_size: number;
    weibo_url: string;
    analytic: string;
    comment: string;
}
interface Applet {
    key: string;
    updated: number;
    created: number;
}
const dataValue = ref<Value>({ github_url: "", page_size: 1, weibo_url: "", analytic: "", comment: "" });
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
            dataValue.value = { github_url: "", page_size: 1, weibo_url: "", analytic: "", comment: "" };
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
