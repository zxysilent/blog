<template>
    <n-grid cols="1" responsive="screen" class="-mt-5">
        <n-grid-item>
            <n-list>
                <n-list-item>
                    <template #suffix>
                        <n-button type="primary" text @click="modalPasswd = true">修改</n-button>
                    </template>
                    <n-thing title="账户密码">
                        <template #description><span class="text-gray-400">绑定手机和邮箱，并设置密码，帐号更安全</span></template>
                    </n-thing>
                </n-list-item>
                <!-- <n-list-item>
          <template #suffix>
            <n-button type="primary" text>修改</n-button>
          </template>
          <n-thing title="绑定手机">
            <template #description
              ><span class="text-gray-400">已绑定手机号：+86189****4877</span></template
            >
          </n-thing>
        </n-list-item> -->
                <!-- <n-list-item>
          <template #suffix>
            <n-button type="primary" text>设置</n-button>
          </template>
          <n-thing title="密保问题">
            <template #description
              ><span class="text-gray-400"
                >未设置密保问题，密保问题可有效保护账户安全</span
              ></template
            >
          </n-thing>
        </n-list-item> -->
                <!-- <n-list-item>
          <template #suffix>
            <n-button type="primary" text>修改</n-button>
          </template>
          <n-thing title="个性域名">
            <template #description
              ><span class="text-gray-400">已绑定域名：https://www.naiveui.com</span></template
            >
          </n-thing>
        </n-list-item> -->
            </n-list>
        </n-grid-item>
    </n-grid>
    <n-modal v-model:show="modalPasswd" :mask-closable="false">
        <n-card style="width: 512px" title="修改密码" :bordered="false" size="huge" role="dialog" aria-modal="true">
            <n-form :model="dataForm" label-placement="left" require-mark-placement="left" label-width="auto" ref="dataRef" :rules="dataRules">
                <n-form-item path="opasswd" label="原密码">
                    <n-input type="password" v-model:value="dataForm.opasswd" placeholder="请输入现在使用的密码" />
                </n-form-item>
                <n-form-item path="npasswd" label="新密码">
                    <n-input type="password" show-password-on="click" v-model:value="dataForm.npasswd" placeholder="请输入新密码，至少6位字符" />
                </n-form-item>
                <n-form-item first path="rpasswd" label="重复密码">
                    <n-input type="password" show-password-on="click" :disabled="!dataForm.npasswd" placeholder="请再次输入新密码" v-model:value="dataForm.rpasswd" />
                </n-form-item>
                <!-- <n-row :gutter="[0, 24]">
                    <n-col :span="24">
                        <div style="display: flex; justify-content: flex-end">
                            <n-button @click="handleValidateButtonClick" :disabled="model.age === null" round type="primary"> 验证 </n-button>
                        </div>
                    </n-col>
                </n-row> -->
            </n-form>
            <template #footer>
                <n-space justify="end">
                    <n-button type="warning" :loading="loading" @click="emitEdit">提交保存</n-button>
                    <n-button type="success" @click="emitReset">重置填写</n-button>
                    <n-button type="tertiary" @click="emitCancel">取消关闭</n-button>
                </n-space>
            </template>
        </n-card>
    </n-modal>
</template>

<script lang="ts" setup>
import { apiAuthEditPasswd } from "@/api";
import { reactive, ref } from "vue";
import storage from "@/utils/storage";
import { useRouter } from "vue-router";
import { useMessage } from "naive-ui";
const loading = ref(false);
const message = useMessage();
const router = useRouter();
const modalPasswd = ref(false);
const dataRef = ref();
const dataForm = reactive({ id: 0, opasswd: "", npasswd: "", rpasswd: "" });
const dataRules = {
    opasswd: { required: true, message: "请输入原密码", trigger: "blur" },
    npasswd: [
        { required: true, message: "请输入新密码", trigger: "blur" },
        { min: 6, message: "请至少输入6个字符", trigger: "blur" },
        { max: 32, message: "最多输入32个字符", trigger: "blur" },
    ],
    rpasswd: [
        { required: true, message: "请再次输入密码", trigger: ["input", "blur"] },
        {
            validator: (_, value) => {
                return value === dataForm.npasswd;
            },
            message: "两次密码输入不一致",
            trigger: ["blur"],
        },
    ],
};
const emitEdit = () => {
    dataRef.value.validate((errors) => {
        if (!errors) {
            loading.value = true;
            apiAuthEditPasswd(dataForm).then((resp) => {
                loading.value = false;
                if (resp.code == 200) {
                    message.success("密码修改成功,请重新登陆", {
                        onAfterLeave() {
                            storage.clear();
                            router.push({ name: "login" });
                        },
                    });
                } else {
                    message.error(resp.msg);
                }
            });
        }
    });
};
const emitReset = () => {
    dataForm.opasswd = "";
    dataForm.npasswd = "";
    dataForm.rpasswd = "";
    dataRef.value.restoreValidation();
};
const emitCancel = () => {
    dataForm.opasswd = "";
    dataForm.npasswd = "";
    dataForm.rpasswd = "";
    modalPasswd.value = false;
};
</script>
