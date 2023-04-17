<template>
    <n-card :bordered="false" title="修改用户" class="mt-4 proCard" size="small" :segmented="{ content: 'hard' }">
        <n-grid cols="1 s:2 m:2 l:3 xl:3 xxl:3" responsive="screen">
            <n-grid-item>
                <n-form :label-width="100" :model="dataForm" label-placement="left" require-mark-placement="left" :rules="dataRules" ref="dataRef">
                    <n-form-item label="登录账号" path="num">
                        <n-input v-model:value="dataForm.num" readonly disabled>
                            <template #suffix>不可修改</template>
                        </n-input>
                    </n-form-item>
                    <n-form-item label="用户姓名" path="name">
                        <n-input v-model:value="dataForm.name" />
                    </n-form-item>
                    <n-form-item label="电话号码" path="phone">
                        <n-input v-model:value="dataForm.phone" />
                    </n-form-item>
                    <n-form-item label="邮箱地址" path="email">
                        <n-input v-model:value="dataForm.email" />
                    </n-form-item>
                    <n-form-item label="用户密码" path="passwd">
                        <n-input v-model:value="dataForm.passwd">
                            <template #suffix>为空不修改</template>
                        </n-input>
                    </n-form-item>
                    <n-form-item label="用户角色" path="role_id">
                        <n-radio-group v-model:value="dataForm.role_id">
                            <n-space>
                                <n-radio v-for="role in roleAll" :key="role.id" :value="role.id">
                                    {{ role.name }}
                                </n-radio>
                            </n-space>
                        </n-radio-group>
                    </n-form-item>
                    <div>
                        <n-space justify="center">
                            <n-button type="warning" :loading="loading" @click="emitEdit">确认修改</n-button>
                            <n-button type="success" @click="emitReset">重置填写</n-button>
                            <n-button type="tertiary" @click="$router.push({ name: 'admin-index' })">取消返回</n-button>
                        </n-space>
                    </div>
                </n-form>
            </n-grid-item>
        </n-grid>
    </n-card>
</template>
<script lang="ts" setup>
import { apiAdminGet, apiAdminEdit } from "@/api";
import { ref, onMounted } from "vue";
import { useMessage } from "naive-ui";
import { apiRoleAll } from "@/api";
import { useRoute } from "vue-router";
interface Role {
    id: number;
    name: string;
}
const route = useRoute();
const roleAll = ref<Role[]>([]);
const dataForm = ref({
    id: 0,
    num: "",
    name: "",
    passwd: "",
    avatar: "",
    phone: "",
    email: "",
    role_id: 0,
});

const dataRules = {
    name: { required: true, message: "请输入用户姓名", trigger: "blur" },
    num: { required: true, message: "请输入登录账号", trigger: "blur" },
    // passwd: { required: true, message: "请输入初始密码", trigger: "blur" },
    role_id: { required: true, type: "number", min: 1, message: "请选择用户角色", trigger: ["blur", "change"] },
};
const dataRef = ref();
const message = useMessage();
const loading = ref(false);
const init = () => {
    apiRoleAll({}).then((resp) => {
        if (resp.code == 200) {
            roleAll.value = resp.data;
        } else {
            roleAll.value = [];
        }
    });
    apiAdminGet({ id: route.params.id }).then((resp) => {
        if (resp.code == 200) {
            dataForm.value = resp.data;
        } else {
            dataForm.value = {
                id: 0,
                num: "",
                name: "",
                avatar: "",
                passwd: "",
                phone: "",
                email: "",
                role_id: 0,
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
    if (dataForm.value.passwd && dataForm.value.passwd != "") {
        console.log(dataForm.value.passwd);
    }
    dataRef.value.validate((errors) => {
        if (!errors) {
            apiAdminEdit(dataForm.value).then((resp) => {
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
