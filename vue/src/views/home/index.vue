<template>
    <div>
        <n-card :bordered="false" title="快速入口" class="mt-4 proCard" size="small" :segmented="{ content: 'hard' }">
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
        <n-card v-auth="'home_basic'" :bordered="false" title="系统信息" class="mt-4 proCard" size="small" :segmented="{ content: 'hard' }">
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
                <!-- <template #title>
                    <a href="https://anyway.fm/" style="text-decoration: none">Anyway.FM</a>
                </template> -->
                <!-- <template #header>
                    <n-space>
                        <n-button text> 播客 </n-button>/<n-button text> 精选 </n-button>/<n-button text> 超级精选 </n-button> 
                    </n-space>
                </template> -->
                <!-- <template #avatar>
                    <n-avatar src="https://cdnimg103.lizhi.fm/user/2017/02/04/2583325032200238082_160x160.jpg" />
                </template> -->
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
        <!-- 
        <n-card :bordered="false" title="快捷入口" class="mt-4 proCard" size="small" :segmented="{ content: 'hard' }">
            <n-grid x-gap="12" y-gap="6" cols="4 xs:1 s:2 m:3 l:4 xl:5 xxl:6" responsive="screen">
                <n-gi>
                    <n-button type="info" attr-type="button" secondary @click="$router.push({name:'query-logon'})" strong block>
                        <template #icon>
                            <n-icon> </n-icon>
                        </template>
                        注册查询
                    </n-button>
                </n-gi>
                <n-gi>
                    <n-button type="success" attr-type="button" secondary @click="$router.push({name:'query-auth'})" strong block>
                        <template #icon>
                            <n-icon> </n-icon>
                        </template>
                        认证查询
                    </n-button>
                </n-gi>
            </n-grid>
        </n-card> -->
        <!-- <n-grid class="mt-4" cols="2 s:1 m:1 l:2 xl:2 xxl:2" responsive="screen" :x-gap="12" :y-gap="9">
            <n-gi>
                <n-card :segmented="{ content: 'hard' }" content-style="padding: 0;" :bordered="false" size="small" title="快捷操作">
                    <div class="flex flex-wrap project-card">
                        <n-card size="small" class="cursor-pointer project-card-item" hoverable>
                            <div class="flex flex-col justify-center text-gray-500">
                                <span class="text-center">
                                    <n-icon size="30" color="#68c755">
                                        <HomeOutline />
                                    </n-icon>
                                </span>
                                <span class="text-lx text-center">管理主页</span>
                            </div>
                        </n-card>
                        <n-card size="small" class="cursor-pointer project-card-item" hoverable>
                            <div class="flex flex-col justify-center text-gray-500">
                                <span class="text-center">
                                    <n-icon size="30" color="#fab251">
                                        <SettingsOutline />
                                    </n-icon>
                                </span>
                                <span class="text-lx text-center">个人设置</span>
                            </div>
                        </n-card>
                        <n-card size="small" class="cursor-pointer project-card-item" hoverable>
                            <div class="flex flex-col justify-center text-gray-500">
                                <span class="text-center">
                                    <n-icon size="30" color="#7238d1">
                                        <InformationCircleOutline />
                                    </n-icon>
                                </span>
                                <span class="text-lx text-center">关于</span>
                            </div>
                        </n-card>
                    </div>
                </n-card>
                <n-card :segmented="{ content: 'hard' }" content-style="padding: 0;" :bordered="false" size="small" title="项目" class="mt-2">
                    <div class="flex flex-wrap project-card">
                        <n-card size="small" class="cursor-pointer project-card-item ms:w-1/2 md:w-1/3" hoverable>
                            <div class="flex">
                                <span>
                                    <n-icon size="30">
                                        <LogoGithub />
                                    </n-icon>
                                </span>
                                <span class="text-lg ml-4">Github</span>
                            </div>
                            <div class="flex mt-2 h-10 text-gray-400">是一个面向开源及私有软件项目的托管平台。</div>
                            <div class="flex mt-2 h-10 text-gray-400">开源君，2021-07-04</div>
                        </n-card>
                        <n-card size="small" class="cursor-pointer project-card-item ms:w-1/2 md:w-1/3" hoverable>
                            <div class="flex">
                                <span>
                                    <n-icon size="30" color="#42b983">
                                        <LogoVue />
                                    </n-icon>
                                </span>
                                <span class="text-lg ml-4">Vue</span>
                            </div>
                            <div class="flex mt-2 h-10 text-gray-400">渐进式 JavaScript 框架</div>
                            <div class="flex mt-2 h-10 text-gray-400">学不动也要学，2021-07-04</div>
                        </n-card>
                        <n-card size="small" class="cursor-pointer project-card-item ms:w-1/2 md:w-1/3" hoverable>
                            <div class="flex">
                                <span>
                                    <n-icon size="30" color="#e44c27">
                                        <LogoHtml5 />
                                    </n-icon>
                                </span>
                                <span class="text-lg ml-4">Html5</span>
                            </div>
                            <div class="flex mt-2 h-10 text-gray-400">HTML5是互联网的下一代标准。</div>
                            <div class="flex mt-2 h-10 text-gray-400">撸码也是一种艺术 2021-04-01</div>
                        </n-card>
                    </div>
                </n-card>
            </n-gi>
             <n-gi>
                <n-card :segmented="{ content: 'hard' }" content-style="padding-top: 0;padding-bottom: 0;" :bordered="false" size="small" title="动态">
                    <template #header-extra><a href="javascript:;">更多</a></template>
                    <n-list>
                        <n-list-item>
                            <template #prefix>
                                <n-avatar circle :size="40" :src="schoolboy" />
                            </template>
                            <n-thing title="Ah Jung 刚才把工作台页面随便写了一些，凑合能看了！">
                                <template #description><p class="text-xs text-gray-500">2021-07-04 22:37:16</p></template>
                            </n-thing>
                        </n-list-item>
                         <n-list-item>
                            <template #prefix>
                                <n-avatar circle :size="40" :src="schoolboy" />
                            </template>
                            <n-thing title="页面切换其实也支持缓存，只是加了过度效果，看起来像是重新渲染了">
                                <template #description>
                                    <p class="text-gray-400">
                                        <n-input type="text" placeholder="不信，输点文字试试" />
                                    </p>
                                </template>
                            </n-thing>
                        </n-list-item> 
                    </n-list>
                </n-card>
            </n-gi> 
        </n-grid> -->
    </div>
</template>

<script lang="ts" setup>
// import schoolboy from "@/assets/images/schoolboy.png";
import { ref, onMounted } from "vue";
// import { InformationCircleOutline, SettingsOutline, HomeOutline, LogoGithub, LogoVue, ReloadOutline, LogoHtml5 } from "@vicons/ionicons5";
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
