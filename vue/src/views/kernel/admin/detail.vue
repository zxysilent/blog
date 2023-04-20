<template>
    <n-card :bordered="false" title="用户详情" size="small" :segmented="{ content: 'hard' }">
        <n-grid cols="1 s:2 m:2 l:3 xl:3 xxl:3" responsive="screen">
            <n-grid-item>
                <n-form :label-width="80" disabled :model="dataForm" label-placement="left" require-mark-placement="left" ref="dataRef">
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
                    <n-form-item label="用户角色" path="role_id">
                        <n-radio-group v-model:value="dataForm.role_id" name="radiogroup">
                            <n-space>
                                <n-radio v-for="role in roleAll" :key="role.id" :value="role.id">
                                    {{ role.name }}
                                </n-radio>
                            </n-space>
                        </n-radio-group>
                    </n-form-item>
                    <n-form-item label="更新时间" :show-feedback="false">
                    <n-time :time="dataForm.updated" />
                </n-form-item>
                <n-form-item label="创建时间" :show-feedback="false">
                    <n-time :time="dataForm.created" />
                </n-form-item>
                    <br />
                    <div>
                        <n-space justify="center">
                            <n-button type="success" @click="init">信息刷新</n-button>
                            <n-button type="tertiary" @click="$router.push({ name: 'admin-index' })">取消返回</n-button>
                        </n-space>
                    </div>
                </n-form>
            </n-grid-item>
        </n-grid>
    </n-card>
</template>
<script lang="ts" setup>
import { apiAdminGet } from "@/api";
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
import { apiRoleAll } from "@/api";
interface Role {
    id: number;
    name: string;
}

const roleAll = ref<Role[]>([]);
const route = useRoute();
let dataForm = ref({
    id: 0,
    num: "",
    name: "",
    avatar: "",
    phone: "",
    email: "",
    role_id: 0,
    updated: 0,
    created: 0,
});
const dataRef = ref();
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
                phone: "",
                email: "",
                role_id: 0,
                updated: 0,
                created: 0,
            };
        }
    });
};
onMounted(() => {
    init();
});
</script>
