<template>
    <n-card :bordered="false" title="添加用户" size="small" :segmented="{ content: 'hard' }">
        <n-grid cols="1 s:2 m:2 l:3 xl:3 xxl:3" responsive="screen">
            <n-grid-item>
                <n-form :label-width="100" :model="dataForm" label-placement="left" require-mark-placement="left" :rules="dataRules" ref="dataRef">
                    <n-form-item label="登录账号" path="num">
                        <n-input v-model:value="dataForm.num">
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
                    <n-form-item label="初始密码" path="passwd">
                        <n-input v-model:value="dataForm.passwd">
                            <template #suffix>默认1q2w3e4r</template>
                        </n-input>
                    </n-form-item>
                    <n-form-item label="用户角色" path="role_id">
                        <n-radio-group v-model:value="dataForm.role_id" name="radiogroup">
                            <n-space>
                                <n-radio v-for="role in roleAll" :key="role.id" :value="role.id">
                                    {{ role.name }}
                                </n-radio>
                            </n-space>
                        </n-radio-group>
                    </n-form-item>
                    <div>
                        <n-space justify="center">
                            <n-button type="warning" :loading="loading" @click="emitAdd">确认添加</n-button>
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
import { apiAdminAdd, apiAdminExist } from "@/api";
import { ref, onMounted } from "vue";
import { useMessage } from "naive-ui";
import { apiRoleAll } from "@/api";
import { useRouter } from "vue-router";
interface Role {
    id: number;
    name: string;
}
const router = useRouter();
const roleAll = ref<Role[]>([]);
const dataForm = ref({
    id: 0,
    num: "",
    name: "",
    passwd: "1q2w3e4r",
    passwd_enc: "",
    avatar: "",
    phone: "",
    email: "",
    role_id: 0,
});
const valideNum = (_, value, callback) => {
    apiAdminExist({ num: value }).then((resp) => {
        if (resp.code == 200) {
            callback(new Error("当前账号已经存在"));
        } else {
            callback();
        }
    });
};
const dataRules = {
    name: { required: true, message: "请输入用户姓名", trigger: "blur" },
    num: [
        { required: true, message: "请输入登录账号", trigger: "blur" },
        { min: 5, message: "请至少输入5个字符", trigger: "blur" },
        { max: 32, message: "最多输入32个字符", trigger: "blur" },
        { validator: valideNum, trigger: "blur" },
    ],
    // phone: { required: true, message: "请输入电话号码", trigger: "blur" },
    // email: { required: true, message: "请输入邮箱地址", trigger: "blur" },
    passwd: [
        { min: 6, message: "请至少输入6个字符", trigger: "blur" },
        { max: 32, message: "最多输入32个字符", trigger: "blur" },
    ],
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
};
onMounted(() => {
    init();
});
const emitReset = () => {
    dataForm.value = {
        id: 0,
        num: "",
        name: "",
        passwd: "1q2w3e4r",
        passwd_enc: "",
        avatar: "",
        phone: "",
        email: "",
        role_id: 0,
    };
    dataRef.value.restoreValidation();
};
const emitAdd = () => {
    console.log(dataForm.value);
    dataRef.value.validate((errors) => {
        if (!errors) {
            apiAdminAdd(dataForm.value).then((resp) => {
                if (resp.code == 200) {
                    message.success("添加成功", {
                        onAfterLeave() {
                            // init();
                            // emitReset();
                            router.push({ name: "admin-index" });
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
