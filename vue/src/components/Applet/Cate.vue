<template>
    <n-tree-select :value="props.value" @update:value="onUpdate" :options="valall" key-field="id" label-field="name" :loading="loading" remote :disabled="props.disabled" />
</template>
<script lang="ts" setup>
import { apiCateList } from "@/api";
import { ref, onMounted } from "vue";
const emit = defineEmits(["update:value"]);
const props = defineProps({
    value: {
        type: [Number, String],
        required: true,
    },
    kind: {
        type: Number,
        default: 0,
    },
    search: {
        type: Boolean,
        default: false,
    },
    disabled: {
        type: Boolean,
        default: false,
    },
    onChange: {
        type: Function,
        default: undefined,
    },
});
const loading = ref(false);
const onUpdate = (value: number, raw: any) => {
    if (props.onChange) {
        props.onChange(raw);
    }
    emit("update:value", value);
};
const valall = ref([]);
const init = () => {
    loading.value = true;
    apiCateList({ kind: props.kind }).then((resp) => {
        loading.value = false;
        if (resp.code == 200) {
            resp.data.unshift({ id: 0, name: props.search ? "所有" : "请选择", kind: 0 });
            valall.value = resp.data.map((item) => {
                if (item.kind == 1) {
                    item.name += "(文章)";
                } else if (item.kind == 3) {
                    item.name += "(笔记)";
                } else if (item.kind == 2) {
                    item.name += "(页面)";
                }
                return item;
            });
        } else {
            valall.value = [];
        }
    });
};
onMounted(() => {
    init();
});
</script>
