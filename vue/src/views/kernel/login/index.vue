<template>
    <div class="login" @keydown.enter="emitLogin">
        <div class="top">
            <div class="header">
                <a>
                    <img src="~@/assets/logo.svg" class="logo" alt="logo" />
                    <span class="title">{{ appForm.title }}</span>
                </a>
            </div>
            <div class="desc">{{ appForm.desc }}</div>
        </div>
        <div class="main">
            <n-form ref="loginRef" show-require-mark label-placement="left" require-mark-placement="left" size="large" label-width="auto" :model="loginForm" :rules="loginRules">
                <n-form-item path="num" label="账 号">
                    <n-input v-model:value="loginForm.num" placeholder="请输入账号">
                        <template #prefix>
                            <n-icon size="18" color="#808695">
                                <PersonOutline />
                            </n-icon>
                        </template>
                    </n-input>
                </n-form-item>
                <n-form-item path="passwd" label="密 码">
                    <n-input v-model:value="loginForm.passwd" type="password" show-password-on="click" placeholder="请输入密码">
                        <template #prefix>
                            <n-icon size="18" color="#808695">
                                <LockClosedOutline />
                            </n-icon>
                        </template>
                    </n-input>
                </n-form-item>
                <n-form-item path="vcode" label="验证码">
                    <n-row>
                        <n-col span="16">
                            <n-input size="large" type="text" v-model:value="loginForm.vcode" placeholder="请输入验证码"></n-input>
                        </n-col>
                        <n-col span="8">
                            <img
                                src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAGQAAAAoCAMAAAA/pq9xAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAyJpVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADw/eHBhY2tldCBiZWdpbj0i77u/IiBpZD0iVzVNME1wQ2VoaUh6cmVTek5UY3prYzlkIj8+IDx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IkFkb2JlIFhNUCBDb3JlIDUuMy1jMDExIDY2LjE0NTY2MSwgMjAxMi8wMi8wNi0xNDo1NjoyNyAgICAgICAgIj4gPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4gPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIgeG1sbnM6eG1wPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvIiB4bWxuczp4bXBNTT0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL21tLyIgeG1sbnM6c3RSZWY9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9zVHlwZS9SZXNvdXJjZVJlZiMiIHhtcDpDcmVhdG9yVG9vbD0iQWRvYmUgUGhvdG9zaG9wIENTNiAoV2luZG93cykiIHhtcE1NOkluc3RhbmNlSUQ9InhtcC5paWQ6NUQ2MzI1ODVGMEM5MTFFQTg4RjZCRTZDRDI5QTA0MDAiIHhtcE1NOkRvY3VtZW50SUQ9InhtcC5kaWQ6NUQ2MzI1ODZGMEM5MTFFQTg4RjZCRTZDRDI5QTA0MDAiPiA8eG1wTU06RGVyaXZlZEZyb20gc3RSZWY6aW5zdGFuY2VJRD0ieG1wLmlpZDo1RDYzMjU4M0YwQzkxMUVBODhGNkJFNkNEMjlBMDQwMCIgc3RSZWY6ZG9jdW1lbnRJRD0ieG1wLmRpZDo1RDYzMjU4NEYwQzkxMUVBODhGNkJFNkNEMjlBMDQwMCIvPiA8L3JkZjpEZXNjcmlwdGlvbj4gPC9yZGY6UkRGPiA8L3g6eG1wbWV0YT4gPD94cGFja2V0IGVuZD0iciI/PgyTenEAAAAGUExURf///wAAAFXC034AAAABdFJOUwBA5thmAAAAHUlEQVR42uzBAQ0AAADCoPdPbQ43oAAAAIAXE2AAD8gAAUGhJ40AAAAASUVORK5CYII="
                                ref="vcodeRef"
                                @click="reload"
                                title="点击刷新"
                                alt="点击刷新"
                                height="40" />
                        </n-col>
                    </n-row>
                </n-form-item>
                <n-form-item>
                    <n-button type="primary" @click="emitLogin" attr-type="button" :loading="loading" block>登&nbsp;&nbsp;&nbsp;&nbsp;录</n-button>
                </n-form-item>
                <!-- <n-form-item>
                    <n-button color="#01aaed" @click="casLogin" attr-type="button" block>CAS统一身份登录</n-button>
                </n-form-item> -->
            </n-form>
            <!-- <p class="other-login">
                <a href="/login/qq.html" title="QQ登陆">
                    <img src="~@/assets/login/qq.svg" alt="加载失败" class="ivu-icon" style="width: 30px" />
                </a>
                <a @click="wechatLogin" title="微信登陆">
                    <img src="~@/assets/login/wechat.svg" alt="加载失败" class="ivu-icon" style="width: 30px" />
                </a>
            </p>-->
            <div class="login-tips">
                <span>系统推荐浏览器</span>&nbsp;
                <br />
                <br />
                <a href="https://www.microsoft.com/zh-cn/edge" target="_blank" title="微软 Edge">
                    <n-icon size="20"> <LogoEdge /> </n-icon>&nbsp;Edge
                </a>
                <a href="https://www.google.cn/chrome/" target="_blank" title="谷歌 Chrome">
                    <n-icon size="20"> <LogoChrome /> </n-icon>&nbsp;Chrome
                </a>
                <a href="http://www.firefox.com.cn/" target="_blank" title="火狐 Firefox">
                    <n-icon size="20"> <LogoFirefox /> </n-icon>&nbsp;Firefox
                </a>
            </div>
        </div>
        <div class="footer">
            <div class="links">
                <a href="/">首页</a>
                <!-- <a>帮助</a> -->
                <!--<a>隐私</a>
                <a>条款</a>-->
                <!-- <a href="https://www.cwnu.edu.cn/" target="_blank">学校首页</a> -->
                <!-- <a href="https://nic.cwnu.edu.cn/" target="_blank">网络中心</a> -->
                <a href="https://github.com/zxysilent" target="_blank">技术支持</a>
            </div>
            <div class="copyright">
                Copyright &copy; {{ new Date().getFullYear() }}&nbsp; <a target="_blank" href="https://github.com/zxysilent">github.com/zxysilent</a>&nbsp;&nbsp;
                <a target="_blank" href="https://blog.zxysilent.com">blog.zxysilent.com</a>
            </div>
        </div>
    </div>
</template>
<script lang="ts" setup>
import { reactive, ref, onMounted } from "vue";
import { useRouter } from "vue-router";
import { useAdminStore } from "@/store/admin";
import { useMessage } from "naive-ui";
import { apiDictVal } from "@/api/ext";
import { apiAuthLogin, apiAuthVcode } from "@/api";
import { PersonOutline, LockClosedOutline, LogoEdge, LogoChrome, LogoFirefox } from "@vicons/ionicons5";
import storage from "@/utils/storage";
const adminStore = useAdminStore();
const loginRef = ref();
const vcodeRef = ref();
const loginRules = {
    num: { required: true, message: "账号不能为空", trigger: "blur" },
    passwd: { required: true, message: "密码不能为空", trigger: "blur" },
    vcode: { required: true, message: "验证码不能为空", trigger: "blur" },
};
const loginForm = reactive({
    num: "",
    passwd: "",
    vcode: "",
    vreal: "",
});
const appForm = ref({
    title: "欢迎登录",
    desc: "便捷、高效、可靠",
});
const loading = ref(false);
const message = useMessage();
const router = useRouter();
const reload = () => {
    apiAuthVcode({}).then((resp) => {
        if (resp.code == 200) {
            vcodeRef.value.src = "data:image/png;base64," + resp.data.image;
            loginForm.vcode = "";
            loginForm.vreal = resp.data.real;
        } else {
            loginForm.vcode = "";
            loginForm.vreal = "";
        }
    });
};
const emitLogin = () => {
    loginRef.value.validate((errors) => {
        if (!errors) {
            loading.value = true;
            apiAuthLogin(loginForm).then((resp) => {
                loading.value = false;
                if (resp.code == 200) {
                    storage.setToken(resp.data);
                    adminStore.emitGranted(false);
                    message.success("登录成功", {
                        duration: 1500,
                        onAfterLeave: () => {
                            router.push({ name: "root" });
                        },
                    });
                } else {
                    message.error(resp.msg);
                    reload();
                }
            });
        }
    });
};
// const wechatLogin = () => {
//     message.warning("暂不支持");
// };
// const casLogin = () => {
//     location.href = "/auth/cas";
// };
const init = () => {
    apiDictVal("login").then((resp) => {
        if (resp.code == 200) {
            appForm.value = resp.data;
        } else {
            appForm.value.title = "欢迎登录";
            appForm.value.desc = "便捷、高效、可靠";
        }
    });
};
onMounted(() => {
    // storage.clear();
    init();
    reload();
});
</script>
<style lang="less" scoped>
.login {
    width: 100%;
    min-height: 100%;
    background: #f0f2f5 url(../../../assets/bg.svg) no-repeat 50%;
    background-size: 100%;
    padding: 110px 0 144px;
    position: relative;
    a {
        text-decoration: none;
    }
    img {
        display: inline-block;
    }
    .top {
        text-align: center;
        .header {
            height: 44px;
            line-height: 44px;
            .logo {
                height: 44px;
                vertical-align: top;
                margin-right: 16px;
                border-style: none;
                border-radius: 50%;
            }
            .title {
                font-size: 30px;
                color: rgba(0, 0, 0, 0.85);
                font-family: "Helvetica Neue", Helvetica, "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei", "微软雅黑", Arial, sans-serif;
                font-weight: 600;
                position: relative;
                top: 2px;
            }
        }
        .desc {
            font-size: 16px;
            color: rgba(0, 0, 0, 0.45);
            margin-top: 12px;
            margin-bottom: 30px;
        }
    }
    .main {
        min-width: 260px;
        width: 368px;
        margin: 0 auto;
        .login-tips {
            text-align: left;
            margin-top: 10px;
            line-height: 22px;
            font-size: 16px;
            a {
                color: rgba(0, 0, 0, 0.45);
                transition: all 0.3s;
                &:not(:last-child) {
                    margin-right: 40px;
                }
            }
        }
        .other-login {
            text-align: center;
            color: #c3c3c3;
            a {
                margin: 10px;
            }
        }
    }
    .footer {
        position: absolute;
        width: 100%;
        bottom: 0;
        padding: 0 16px;
        margin: 48px 0 24px;
        text-align: center;
        .links {
            margin-bottom: 8px;
            font-size: 14px;
            a {
                color: rgba(0, 0, 0, 0.45);
                transition: all 0.3s;
                &:not(:last-child) {
                    margin-right: 40px;
                }
            }
        }
        .copyright {
            color: rgba(0, 0, 0, 0.45);
            font-size: 14px;
        }
    }
}
</style>
