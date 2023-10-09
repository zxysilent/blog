<template>
    <n-el class="note-card">
        <div class="note-card-body">
            <slot></slot>
        </div>
    </n-el>
</template>
<style lang="less">
.note-card {
    line-height: 1.6;
    border-radius: 3px;
    position: relative;
    transition: background-color 0.3s cubic-bezier(4, 0, 0.2, 1);
    background-color: #f4f6f7;
    text-align: start;
    word-break: break-word;
    padding: 8px;
    margin-top: 0.2rem;
    margin-left: 0.3rem;
    margin-right: 0.5rem;
    border-bottom: 3px solid #ffffff;
    cursor: pointer;
}
.note-card:hover {
    background-color: #d2d5d9;
    transform: translate(0, 0);
    transition: all 0.2s ease-in 0s;
}
.active {
    color: #232d47;
    background-color: #d2d5d9;
    border-color: #2d8cf0;
}
.note-card-body {
    color: #333639;
}
</style>
