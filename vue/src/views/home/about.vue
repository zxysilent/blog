<template>
    <div>
        <n-card :bordered="false" title="关于"> <strong>管理后台</strong></n-card>
        <!-- <n-card :bordered="false" title="应用信息" class="mt-4 proCard" size="small" :segmented="{ content: 'hard' }">
            <n-descriptions bordered label-placement="top" :column="3" class="py-2">
                <n-descriptions-item label="系统架构">
                    <n-tag type="info"> {{ StateGo.arch }} </n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="操作系统">
                    <n-tag type="info"> {{ StateGo.os }} </n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="运行环境">
                    <n-tag type="info"> {{ StateGo.version }} </n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="逻辑处理器">
                    <n-tag type="info"> {{ StateGo.num_cpu }} </n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="当前协程数">
                    <n-tag type="info"> {{ StateGo.num_goroutine }} </n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="作者主页">
                    <div class="flex items-center">
                        <a href="https://github.com/zxysilent" class="py-2" style="color: #2d8cf0" target="_blank">点击前往</a>
                    </div>
                </n-descriptions-item>
            </n-descriptions>
        </n-card> -->
        <!-- <n-card :bordered="false" title="测试区域">
            <Images v-model:value="data"></Images>
            <br />
            {{ data }}
        </n-card> -->
        <!-- <n-card :bordered="false" title="测试区域">
            <Mce :value="'xxxx'" height="400px"></Mce>
        </n-card> -->
        <n-card :bordered="false" title="测试区域"> </n-card>
    </div>
</template>
<script lang="ts" setup></script>
<style lang="less" scoped></style>
