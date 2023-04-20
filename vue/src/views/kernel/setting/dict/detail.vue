<template>
    <n-card :bordered="false" title="修改配置" size="small" :segmented="{ content: 'hard' }">
        <n-grid cols="1 s:2 m:2 l:3 xl:3 xxl:3" responsive="screen">
            <n-grid-item>
                <n-form :label-width="100" :model="dataForm" disabled label-placement="left" require-mark-placement="left" ref="dataRef">
                    <n-form-item label="配置名" path="key">
                        <n-input v-model:value="dataForm.key" maxlength="64" show-count readonly disabled>
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
                            <n-button type="tertiary" @click="$router.push({ name: 'dict-index' })">取消返回</n-button>
                        </n-space>
                    </div>
                </n-form>
            </n-grid-item>
        </n-grid>
    </n-card>
</template>
<script lang="ts" setup>
import { apiDictGet } from "@/api/ext";
import { ref, onMounted } from "vue";
import { useRoute } from "vue-router";
const route = useRoute();
const kinds = ref<string[]>(["string", "time", "bool", "int", "json"]);
const dataForm = ref({
    key: "",
    value: "",
    kind: "string",
    intro: "",
    updated: 0,
    created: 0,
});
const dataRef = ref();
const init = () => {
    apiDictGet(route.params.key as string).then((resp) => {
        if (resp.code == 200) {
            dataForm.value = resp.data;
        } else {
            dataForm.value = {
                key: "",
                value: "",
                kind: "string",
                intro: "",
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
