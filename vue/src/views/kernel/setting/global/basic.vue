<template>
    <n-grid cols="1 s:1 m:1 l:2 xl:3 xxl:3" responsive="screen">
        <n-grid-item>
            <n-form :label-width="120" :model="dataValue" label-placement="left" require-mark-placement="left" :rules="dataRules" ref="dataRef" :disabled="disable">
                <n-form-item label="网站地址" path="site_url">
                    <n-input v-model:value="dataValue.site_url" maxlength="64" show-count placeholder="请填写网站地址"></n-input>
                </n-form-item>
                <n-form-item label="网站标题" path="title">
                    <n-input v-model:value="dataValue.title" maxlength="64" show-count placeholder="请填写网站标题"></n-input>
                </n-form-item>
                <n-form-item label="Logo地址" path="logo_url">
                    <Image v-model:value="dataValue.logo_url" />
                </n-form-item>
                <n-form-item label="网站作者" path="author">
                    <n-input v-model:value="dataValue.author" maxlength="64" show-count placeholder="请填写网站作者"></n-input>
                </n-form-item>
                <n-form-item label="SEO关键词" path="keywords">
                    <n-input v-model:value="dataValue.keywords" maxlength="64" show-count placeholder="请填写SEO关键词"></n-input>
                </n-form-item>
                <n-form-item label="SEO描述" path="description">
                    <n-input v-model:value="dataValue.description" type="textarea" :autosize="{ minRows: 2, maxRows: 4 }" maxlength="128" show-count placeholder="请填写SEO描述词"></n-input>
                </n-form-item>
                <n-form-item label="ICP备案" path="beian_miit">
                    <n-input v-model:value="dataValue.beian_miit" maxlength="64" show-count placeholder="请填写ICP备案"></n-input>
                </n-form-item>
                <n-form-item label="公安备案" path="beian_nism">
                    <n-input v-model:value="dataValue.beian_nism" maxlength="64" show-count placeholder="请填写公安备案"></n-input>
                </n-form-item>
                <n-form-item label="版权" path="copyright">
                    <n-input v-model:value="dataValue.copyright" type="textarea" :autosize="{ minRows: 2, maxRows: 4 }" maxlength="128" show-count placeholder="请填写版权"></n-input>
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
import { Image } from "@/components/Uploader";
const KEY = "basic";
interface Value {
    title: string;
    site_url: string;
    logo_url: string;
    author: string;
    copyright: string;
    description: string;
    keywords: string;
    beian_miit: string;
    beian_nism: string;
}
interface Applet {
    key: string;
    updated: number;
    created: number;
}
const dataValue = ref<Value>({ title: "", site_url: "", logo_url: "", author: "", copyright: "", description: "", keywords: "", beian_miit: "", beian_nism: "" });
const dataForm = ref<Applet>({ key: KEY, updated: 0, created: 0 });
const disable = ref<boolean>(true);
const dataRules = {
    title: { required: true, message: "请输入网站标题", trigger: "blur" },
    site_url: { required: true, message: "请输入网站地址", trigger: "blur" },
    logo_url: { required: true, message: "请上传Logo", trigger: "blur" },
    copyright: { required: true, message: "请输入版权", trigger: "blur" },
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
            dataValue.value = { title: "", site_url: "", logo_url: "", author: "", copyright: "", description: "", keywords: "", beian_miit: "", beian_nism: "" };
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
    dataRef.value.restoreValidation();
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
