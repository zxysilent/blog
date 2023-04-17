<template>
    <n-grid cols="2 s:2 m:2 l:3 xl:3 xxl:3" responsive="screen">
        <n-grid-item>
            <n-form :label-width="80" label-placement="left" :disabled="disable" require-mark-placement="left" :model="dataForm" :rules="dataRules" ref="dataRef">
                <n-form-item label="登录账号" path="num">
                    <n-input v-model:value="dataForm.num" readonly disabled />
                </n-form-item>
                <n-form-item label="用户姓名" path="name">
                    <n-input v-model:value="dataForm.name" />
                </n-form-item>
                <n-form-item label="用户头像">
                    <Image ref="uploadRef" v-model:value="dataForm.avatar"></Image>
                </n-form-item>
                <n-form-item label="电话" path="phone">
                    <n-input v-model:value="dataForm.phone" />
                </n-form-item>
                <n-form-item label="邮箱" path="email">
                    <n-input v-model:value="dataForm.email" />
                </n-form-item>
                <div>
                    <n-space justify="center">
                        <n-button type="warning" :disabled="disable" :loading="loading" @click="emitEdit">确认修改</n-button>
                        <n-button type="success" @click="emitReset">取消刷新</n-button>
                        <n-button type="tertiary"  :disabled="!disable" @click="openEdit">开启编辑</n-button>
                    </n-space>
                </div>
            </n-form>
        </n-grid-item>
    </n-grid>
</template>

<script lang="ts" setup>
import { apiAuthGet, apiAuthEdit } from "@/api";
import { reactive, ref, onMounted } from "vue";
import { useMessage } from "naive-ui";
import { Image } from "@/components/Uploader";
const dataForm = reactive({ id: 1, name: "", num: "", phone: "", email: "", created: "", avatar: "" });
const dataRules = {
    name: { required: true, message: "请输入用户姓名", trigger: "blur" },
};
const loading = ref(false);
const disable = ref<boolean>(true);
const dataRef = ref();
const message = useMessage();
const init = () => {
    apiAuthGet({}).then((resp) => {
        if (resp.code == 200) {
            dataForm.id = resp.data.id;
            dataForm.num = resp.data.num;
            dataForm.name = resp.data.name;
            dataForm.phone = resp.data.phone;
            dataForm.avatar = resp.data.avatar;
            dataForm.created = resp.data.created;
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
    dataRef.value.validate((errors) => {
        if (!errors) {
            loading.value = true;
            apiAuthEdit(dataForm).then((resp) => {
                loading.value = false;
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
