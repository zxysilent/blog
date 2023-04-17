<template>
    <div class="login">
        <div class="top">
            <div class="header">
                <a>
                    <img src="~@/assets/logo.svg" class="logo" alt="logo" />
                    <span class="title">欢迎登陆</span>
                </a>
            </div>
            <div class="desc">学生学校邮箱自助申请系统</div>
        </div>
        <div class="main">
            <n-form ref="loginRef" show-require-mark label-placement="left" require-mark-placement="left" size="large" label-width="auto">
                <n-form-item>
                    <n-button color="#01aaed" @click="casLogin" attr-type="button" block> CAS统一身份登录 </n-button>
                </n-form-item>
            </n-form>
        </div>
        <div class="footer">
            <div class="links">
                <!-- <a href="/">首页</a> -->
                <!-- <a>帮助</a>
                <a>隐私</a>
                <a>条款</a> -->
                <a href="https://www.cwnu.edu.cn/" target="_blank"> 学校首页 </a>
                <a href="http://ehall.cwnu.edu.cn/" target="_blank"> 信息门户 </a>
                <a href="https://nic.cwnu.edu.cn/" target="_blank"> 网络中心 </a>
                <a href="https://github.com/zxysilent" target="_blank"> 技术支持 </a>
            </div>
            <div class="copyright">
                Copyright &copy; {{ new Date().getFullYear() }}&nbsp;<a target="_blank" href="https://github.com/zxysilent">github.com/zxysilent</a>&nbsp;&nbsp;<a target="_blank" href="https://blog.zxysilent.com"
                    >blog.zxysilent.com</a
                >
            </div>
        </div>
    </div>
</template>
<script lang="ts" setup>
// import { reactive, ref, onMounted } from "vue";
// import { useRouter } from "vue-router";
// import { useAdminStore } from "@/store/user";
// import { useMessage } from "naive-ui";
// import { PersonOutline, LockClosedOutline, LogoEdge, LogoChrome, LogoFirefox } from "@vicons/ionicons5";
const casLogin = () => {
    location.href = "/auth/cas";
};
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
        margin-top: 100px;
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
