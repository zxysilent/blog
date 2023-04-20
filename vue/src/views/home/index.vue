<template>
    <div>
        <n-card :bordered="false" title="快速入口" size="small" :segmented="{ content: 'hard' }">
            <n-space>
                <n-button strong secondary @click="gotoDoc"> 文档 </n-button>
                <n-button strong secondary type="tertiary" @click="gotoHome"> 首页 </n-button>
                <n-button strong secondary type="primary"> Primary </n-button>
                <n-button strong secondary type="info"> Info </n-button>
                <n-button strong secondary type="success"> Success </n-button>
                <n-button strong secondary type="warning"> Warning </n-button>
                <n-button strong secondary type="error"> Error </n-button>
                <n-button strong secondary round> Default </n-button>
                <n-button strong secondary round type="primary"> Primary </n-button>
                <n-button strong secondary round type="info"> Info </n-button>
                <n-button strong secondary round type="success"> Success </n-button>
                <n-button strong secondary round type="warning"> Warning </n-button>
                <n-button strong secondary round type="error"> Error </n-button>
            </n-space>
        </n-card>
        <n-card v-auth="'home_basic'" :bordered="false" title="系统信息" size="small" :segmented="{ content: 'hard' }">
            <n-descriptions bordered label-placement="top" :column="3" class="py-2">
                <n-descriptions-item label="系统架构">
                    <n-tag type="info"> {{ stateGo.arch }} </n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="操作系统">
                    <n-tag type="info"> {{ stateGo.os }} </n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="运行环境">
                    <n-tag type="info"> {{ stateGo.version }} </n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="逻辑处理器">
                    <n-tag type="info"> {{ stateGo.num_cpu }} </n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="当前协程数">
                    <n-tag type="info"> {{ stateGo.num_goroutine }} </n-tag>
                </n-descriptions-item>
                <n-descriptions-item label="作者主页">
                    <div class="flex items-center">
                        <a href="https://github.com/zxysilent" class="py-2" style="color: #2d8cf0" target="_blank">点击前往</a>
                    </div>
                </n-descriptions-item>
            </n-descriptions>
        </n-card>
        <n-card v-if="0" :bordered="false" title="应用信息" v-auth="'home_mgmt'">
            <n-page-header subtitle="应用统计信息">
                <n-grid :cols="4">
                    <n-gi>
                        <n-statistic label="日志" :value="stateApp.log" />
                    </n-gi>
                    <n-gi>
                        <n-statistic label="用户" :value="stateApp.user" />
                    </n-gi>
                    <n-gi>
                        <n-statistic label="资料" :value="stateApp.pdf" />
                    </n-gi>
                    <n-gi>
                        <n-statistic label="题库" :value="stateApp.lib" />
                    </n-gi>
                </n-grid>
                <template #extra>
                    <n-space>
                        <n-button tertiary @click="init">
                            <template #icon>
                                <n-icon>
                                    <ReloadOutline></ReloadOutline>
                                </n-icon>
                            </template>
                            刷新
                        </n-button>
                        <!-- <n-dropdown :options="options" placement="bottom-start">
                            <n-button :bordered="false" style="padding: 0 4px">···</n-button>
                        </n-dropdown> -->
                    </n-space>
                </template>
                <template #footer>截止到 <n-time type="datetime" :time="stateApp.updated"></n-time></template>
            </n-page-header>
        </n-card>
    </div>
</template>

<script lang="ts" setup>
import { ref, onMounted } from "vue";
import { ReloadOutline } from "@vicons/ionicons5";
import { apiStatusApp, apiStatusGo } from "@/api";
import { useAdminStore } from "@/store/admin";

const adminStore = useAdminStore();
const stateGo = ref({ arch: "--", num_cpu: 0, os: "--", version: "--", num_goroutine: 0 });
const stateApp = ref({ user: 0, pdf: 0, book: 0, lib: 0, log: 0, updated: 0 });
const init = () => {
    apiStatusGo({}).then((resp) => {
        if (resp.code == 200) {
            stateGo.value = resp.data;
        }
    });
    if (adminStore.Authed("grant_admin")) {
        apiStatusApp({}).then((resp) => {
            if (resp.code == 200) {
                stateApp.value = resp.data;
            }
        });
    }
};
const gotoDoc = () => {
    window.open(import.meta.env.VITE_APP_SRV + "/swagger/index.html", "_blank");
};
const gotoHome = () => {
    window.open(import.meta.env.VITE_APP_SRV + "/", "_blank");
};
onMounted(() => {
    init();
});
</script>
<style lang="less" scoped>
.project-card {
    margin-right: -6px;
    &-item {
        margin: -1px;
        width: 33.333333%;
    }
}
</style>
