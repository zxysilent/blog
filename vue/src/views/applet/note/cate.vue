<template>
    <n-el class="note-cate" :style="{ top: `${props.index * 45 + 200}px`, backgroundColor: props.color }">
        <slot></slot>
    </n-el>
</template>
<script lang="ts" setup>
const props = defineProps({
    index: {
        type: Number,
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
</script>

<style lang="less">
.note-cate {
    line-height: 1.5;
    box-sizing: border-box;
    border-width: 0;
    border-style: solid;
    border-color: #e5e7eb;
    position: fixed;
    right: 0px;
    // top: 50%;
    display: flex;
    align-items: center;
    width: 60px;
    height: 40px;
    color: #fff;
    border-radius: 20px 0 0 20px;
    cursor: pointer;
    z-index: 1000;
    padding-left: 10px;
    transition: width 0.5s;
    opacity: 0.6;
    text-overflow: ellipsis;
    white-space: nowrap;
    // background-color: #2d8cf0;
}
.note-cate:hover {
    width: 80px;
    opacity: 1;
    transition: all 0.5s;
}
.cate-active {
    opacity: 1;
    width: 80px;
}
.note-cate-body {
    color: #333639;
}
</style>
