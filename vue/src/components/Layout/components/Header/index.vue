<template>
    <div class="layout-header">
        <!--顶部菜单-->
        <div class="layout-header-left" v-if="state.navMode === 'horizontal'">
            <div class="logo">
                <img src="~@/assets/logo.svg" alt="" />
                <h2 v-show="!collapsed" class="title">BLOG</h2>
            </div>
            <AsideMenu :collapsed="collapsed" :inverted="getInverted" mode="horizontal" />
        </div>
        <!--左侧菜单-->
        <div class="layout-header-left" v-else>
            <!-- 菜单收起 -->
            <div class="ml-1 layout-header-trigger layout-header-trigger-min" @click="() => $emit('update:collapsed', !collapsed)">
                <n-icon size="18" v-if="collapsed">
                    <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 1024 1024">
                        <path
                            d="M408 442h480c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8H408c-4.4 0-8 3.6-8 8v56c0 4.4 3.6 8 8 8zm-8 204c0 4.4 3.6 8 8 8h480c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8H408c-4.4 0-8 3.6-8 8v56zm504-486H120c-4.4 0-8 3.6-8 8v56c0 4.4 3.6 8 8 8h784c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8zm0 632H120c-4.4 0-8 3.6-8 8v56c0 4.4 3.6 8 8 8h784c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8zM142.4 642.1L298.7 519a8.84 8.84 0 0 0 0-13.9L142.4 381.9c-5.8-4.6-14.4-.5-14.4 6.9v246.3a8.9 8.9 0 0 0 14.4 7z"
                            fill="currentColor"></path>
                    </svg>
                </n-icon>
                <n-icon size="18" v-else>
                    <svg xmlns="http://www.w3.org/2000/svg" xmlns:xlink="http://www.w3.org/1999/xlink" viewBox="0 0 1024 1024">
                        <path
                            d="M408 442h480c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8H408c-4.4 0-8 3.6-8 8v56c0 4.4 3.6 8 8 8zm-8 204c0 4.4 3.6 8 8 8h480c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8H408c-4.4 0-8 3.6-8 8v56zm504-486H120c-4.4 0-8 3.6-8 8v56c0 4.4 3.6 8 8 8h784c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8zm0 632H120c-4.4 0-8 3.6-8 8v56c0 4.4 3.6 8 8 8h784c4.4 0 8-3.6 8-8v-56c0-4.4-3.6-8-8-8zM115.4 518.9L271.7 642c5.8 4.6 14.4.5 14.4-6.9V388.9c0-7.4-8.5-11.5-14.4-6.9L115.4 505.1a8.74 8.74 0 0 0 0 13.8z"
                            fill="currentColor"></path>
                    </svg>
                </n-icon>
            </div>
            <div class="mr-1 layout-header-trigger layout-header-trigger-min" @click="notePage">
                <n-icon size="18">
                    <svg viewBox="0 0 16 16">
                        <g fill="none">
                            <path
                                d="M5 1a.5.5 0 0 1 .5.5V2h2v-.5a.5.5 0 0 1 1 0V2h2v-.5a.5.5 0 0 1 1 0V2A1.5 1.5 0 0 1 13 3.5v2.536a2.547 2.547 0 0 0-1 .406V3.5a.5.5 0 0 0-.5-.5h-7a.5.5 0 0 0-.5.5v10a.5.5 0 0 0 .5.5h1.547v.002a1.59 1.59 0 0 0 .068.998H4.5A1.5 1.5 0 0 1 3 13.5v-10A1.5 1.5 0 0 1 4.5 2v-.5A.5.5 0 0 1 5 1zm5 7c.107 0 .206.034.288.091L9.378 9H6a.5.5 0 0 1 0-1h4zm-3.004 3.435A.5.5 0 0 0 6.5 11H6a.5.5 0 0 0 0 1h.5a.498.498 0 0 0 .157-.025c.097-.189.21-.37.339-.54zM6 5a.5.5 0 0 0 0 1h4a.5.5 0 0 0 0-1H6zm6.338 2.455a1.56 1.56 0 0 1 2.207 2.207l-4.289 4.288a2.777 2.777 0 0 1-1.29.731l-1.211.303a.61.61 0 0 1-.74-.74l.304-1.21c.122-.489.374-.935.73-1.29l4.289-4.289z"
                                fill="currentColor"></path>
                        </g>
                    </svg>
                </n-icon>
            </div>
            <!-- 刷新 -->
            <div class="mr-1 layout-header-trigger layout-header-trigger-min" @click="reloadPage">
                <n-icon size="18">
                    <ReloadOutline />
                </n-icon>
            </div>
        </div>
        <div class="layout-header-right">
            <!-- <div class="layout-header-trigger layout-header-trigger-min" v-for="item in iconList" :key="item.icon.name">
                <n-tooltip placement="bottom">
                    <template #trigger>
                        <n-icon size="18">
                            <component :is="item.icon" v-on="item.eventObject || {}" />
                        </n-icon>
                    </template>
                    <span>{{ item.tips }}</span>
                </n-tooltip>
            </div> -->
            <!--切换全屏-->
            <!-- <div class="layout-header-trigger layout-header-trigger-min">
                <n-tooltip placement="bottom">
                    <template #trigger>
                        <n-icon size="18" @click="toggleFullScreen">
                            <FullscreenExitOutlined v-if="fullScreen"></FullscreenExitOutlined>
                            <FullscreenOutlined v-else></FullscreenOutlined>
                        </n-icon>
                    </template>
                    <span>全屏</span>
                </n-tooltip>
            </div> -->
            <!-- 个人中心 -->
            <div class="layout-header-trigger layout-header-trigger-min">
                <n-dropdown trigger="hover" @select="avatarSelect" :options="avatarOptions">
                    <!-- <div class="avatar">
                        <n-icon size="30" color="#2d8cf0" style="font-weight: bold">
                            <PersonSharp />
                        </n-icon>
                    </div> -->
                    <div class="avatar">
                        <n-avatar round :src="dataForm.avatar"> </n-avatar>
                    </div>
                </n-dropdown>
            </div>
        </div>
    </div>
</template>
<script lang="ts" setup>
import { reactive, computed, onMounted } from "vue";
import { useRouter } from "vue-router";
import { apiAuthGet } from "@/api";
import { ReloadOutline } from "@vicons/ionicons5";
import { useDialog, useMessage } from "naive-ui";
import { useAdminStore } from "@/store/admin";
import { AsideMenu } from "@/components/Layout/components/Menu";
import setting from "@/utils/global";
import storage from "@/utils/storage";
const dataForm = reactive({ id: 1, name: "", num: "", phone: "", email: "", created: "", avatar: "" });
const init = () => {
    apiAuthGet({}).then((resp) => {
        if (resp.code == 200) {
            dataForm.id = resp.data.id;
            dataForm.num = resp.data.num;
            dataForm.name = resp.data.name;
            dataForm.phone = resp.data.phone;
            dataForm.avatar = resp.data.avatar;
            dataForm.created = resp.data.created;
        }
    });
};
onMounted(() => {
    init();
});
const router = useRouter();
const props = defineProps({
    collapsed: { type: Boolean },
    inverted: { type: Boolean },
});
const adminStore = useAdminStore();
const message = useMessage();
const dialog = useDialog();

const state = reactive({
    username: adminStore.name || "",
    navMode: setting.navMode,
    navTheme: setting.navTheme,
});

const getInverted = computed(() => {
    return ["light", "header-dark"].includes(setting.navTheme) ? props.inverted : !props.inverted;
});

// 刷新页面
const reloadPage = () => {
    location.reload();
};
// 笔记模式
const notePage = () => {
    router.push({ name: "note" });
};
// // 切换全屏
// const fullScreen = ref(false);
// const toggleFullscreenIcon = () => {
//     fullScreen.value = document.fullscreenElement!=null;
// };
// // 监听全屏切换事件
// document.addEventListener("fullscreenchange", toggleFullscreenIcon);
// // 全屏切换
// const toggleFullScreen = () => {
//     if (!document.fullscreenElement) {
//         document.documentElement.requestFullscreen();
//     } else {
//         if (document.exitFullscreen) {
//             document.exitFullscreen();
//         }
//     }
// };

// 图标列表
// const iconList = [
//     {
//         icon: "SearchOutline",
//         tips: "搜索",
//     },
//     {
//         icon: "GithubOutlined",
//         tips: "github",
//         eventObject: {
//             click: () => window.open("https://github.com/jekip/naive-ui-admin"),
//         },
//     },
// ];
const avatarOptions = [
    { label: "个人设置", key: 1 },
    { label: "退出登录", key: 2 },
];

//头像下拉菜单
const avatarSelect = (key) => {
    switch (key) {
        case 1:
            router.push({ name: "account" });
            break;
        case 2:
            dialog.info({
                title: "提示",
                content: "您确定要退出登录吗",
                positiveText: "确定",
                negativeText: "取消",
                onPositiveClick: () => {
                    storage.clear();
                    message.success("成功退出登录", {
                        duration: 1500,
                        onAfterLeave: () => {
                            router.replace({
                                name: "login",
                                // query: { redirect: route.fullPath },
                            });
                        },
                    });
                },
                onNegativeClick: () => {},
            });
            break;
    }
};
</script>
<style lang="less" scoped>
.layout-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0;
    height: @header-height;
    box-shadow: 0 1px 4px rgb(0 21 41 / 8%);
    transition: all 0.2s ease-in-out;
    width: 100%;
    z-index: 11;
    &-left {
        display: flex;
        align-items: center;
        .logo {
            display: flex;
            align-items: center;
            justify-content: center;
            height: 64px;
            line-height: 64px;
            overflow: hidden;
            white-space: nowrap;
            padding-left: 10px;
            img {
                width: auto;
                height: 32px;
                margin-right: 10px;
            }
            .title {
                margin: 0;
                font-size: inherit;
                font-weight: inherit;
            }
        }
        ::v-deep(.ant-breadcrumb span:last-child .link-text) {
            color: #515a6e;
        }
        .n-breadcrumb {
            display: inline-block;
        }
        &-menu {
            color: var(--text-color);
        }
    }
    &-right {
        display: flex;
        align-items: center;
        margin-right: 20px;
        .avatar {
            display: flex;
            align-items: center;
            height: 64px;
        }
        > * {
            cursor: pointer;
        }
    }
    &-trigger {
        display: inline-block;
        width: 64px;
        height: 64px;
        text-align: center;
        cursor: pointer;
        transition: all 0.2s ease-in-out;
        .n-icon {
            display: flex;
            align-items: center;
            height: 64px;
            line-height: 64px;
        }
        &:hover {
            background: hsla(0, 0%, 100%, 0.08);
        }
        .anticon {
            font-size: 16px;
            color: #515a6e;
        }
    }
    &-trigger-min {
        width: auto;
        padding: 0 12px;
    }
}
.layout-header-light {
    background: #fff;
    color: #515a6e;
    .n-icon {
        color: #515a6e;
    }
    .layout-header-left {
        ::v-deep(.n-breadcrumb .n-breadcrumb-item:last-child .n-breadcrumb-item__link) {
            color: #515a6e;
        }
    }
    .layout-header-trigger {
        &:hover {
            background: #f8f8f9;
        }
    }
}
.layout-header-fix {
    position: fixed;
    top: 0;
    right: 0;
    left: 200px;
    z-index: 11;
}

//::v-deep(.menu-router-link) {
//  color: #515a6e;
//
//  &:hover {
//    color: #1890ff;
//  }
//}
</style>
