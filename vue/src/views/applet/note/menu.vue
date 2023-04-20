<template>
    <n-el class="note-menu">
        <div class="note-menu-item" v-for="item in options" :style="{ backgroundColor: item.color, '--tooltip-color': item.color }">
            {{ item.name }}
        </div>
    </n-el>
</template>
<script lang="ts" setup>
import { computed } from "vue";
interface Cate {
    id: number;
    name: string;
    color: string;
    intro: string;
}
const props = defineProps({
    options: {
        type: Array<Cate>,
        required: true,
    },
    tip: {
        type: String,
        default: "",
    },
    color: {
        type: String,
        default: "#2d8cf0",
    },
});
const options = computed(() => {
    return props.options;
});
</script>
<style lang="less">
.note-menu {
    width: 80px;
    flex-shrink: 0;
    box-sizing: border-box;
    position: relative;
    z-index: 1;
    background-color: #f1f2f3;
    display: flex;
    flex-direction: column;
    // justify-content: flex-end;
}

.active {
    color: #333639;
    background-color: #eaf3fd;
    border-color: #2d8cf0;
}
.note-menu-item {
    cursor: pointer;
    color: #e5e7eb;
    margin-top: 0.2rem;
    padding: 6px;
    width: 80px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    opacity: 0.6;
    background-color: aquamarine;
}
.note-menu-item:hover {
    opacity: 1;
    color: var(--tooltip-color);
    background-color: #ffffff !important;
}
</style>
