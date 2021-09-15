<template>
	<div class="login" @keydown.enter="emitLogin">
		<div class="top">
			<div class="header">
				<a>
					<img src="../assets/logo.png" class="logo" alt="logo">
					<span class="title">欢迎登陆</span>
				</a>
			</div>
			<div class="desc">
				github.com/zxysilent/blog
			</div>
		</div>
		<div class="main">
			<Form ref="loginForm" label-position="right" :label-width="65" :model="dataForm" :rules="rules">
				<FormItem prop="num" label="账 号">
					<Input size="large" prefix="ios-person-outline" type="text" v-model="dataForm.num" placeholder="请输入账号">
					</Input>
				</FormItem>
				<FormItem prop="passwd" label="密 码">
					<Input size="large" prefix="ios-lock-outline" type="password" password v-model="dataForm.passwd" placeholder="请输入密码">
					</Input>
				</FormItem>
				<FormItem prop="vcode" label="验证码">
					<Row>
						<Col span="16">
						<Input size="large" type="text" v-model="dataForm.vcode" placeholder="请输入验证码"> </Input>
						</Col>
						<Col span="8">
						<img src="data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAGQAAAAoCAMAAAA/pq9xAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAAAyJpVFh0WE1MOmNvbS5hZG9iZS54bXAAAAAAADw/eHBhY2tldCBiZWdpbj0i77u/IiBpZD0iVzVNME1wQ2VoaUh6cmVTek5UY3prYzlkIj8+IDx4OnhtcG1ldGEgeG1sbnM6eD0iYWRvYmU6bnM6bWV0YS8iIHg6eG1wdGs9IkFkb2JlIFhNUCBDb3JlIDUuMy1jMDExIDY2LjE0NTY2MSwgMjAxMi8wMi8wNi0xNDo1NjoyNyAgICAgICAgIj4gPHJkZjpSREYgeG1sbnM6cmRmPSJodHRwOi8vd3d3LnczLm9yZy8xOTk5LzAyLzIyLXJkZi1zeW50YXgtbnMjIj4gPHJkZjpEZXNjcmlwdGlvbiByZGY6YWJvdXQ9IiIgeG1sbnM6eG1wPSJodHRwOi8vbnMuYWRvYmUuY29tL3hhcC8xLjAvIiB4bWxuczp4bXBNTT0iaHR0cDovL25zLmFkb2JlLmNvbS94YXAvMS4wL21tLyIgeG1sbnM6c3RSZWY9Imh0dHA6Ly9ucy5hZG9iZS5jb20veGFwLzEuMC9zVHlwZS9SZXNvdXJjZVJlZiMiIHhtcDpDcmVhdG9yVG9vbD0iQWRvYmUgUGhvdG9zaG9wIENTNiAoV2luZG93cykiIHhtcE1NOkluc3RhbmNlSUQ9InhtcC5paWQ6NUQ2MzI1ODVGMEM5MTFFQTg4RjZCRTZDRDI5QTA0MDAiIHhtcE1NOkRvY3VtZW50SUQ9InhtcC5kaWQ6NUQ2MzI1ODZGMEM5MTFFQTg4RjZCRTZDRDI5QTA0MDAiPiA8eG1wTU06RGVyaXZlZEZyb20gc3RSZWY6aW5zdGFuY2VJRD0ieG1wLmlpZDo1RDYzMjU4M0YwQzkxMUVBODhGNkJFNkNEMjlBMDQwMCIgc3RSZWY6ZG9jdW1lbnRJRD0ieG1wLmRpZDo1RDYzMjU4NEYwQzkxMUVBODhGNkJFNkNEMjlBMDQwMCIvPiA8L3JkZjpEZXNjcmlwdGlvbj4gPC9yZGY6UkRGPiA8L3g6eG1wbWV0YT4gPD94cGFja2V0IGVuZD0iciI/PgyTenEAAAAGUExURf///wAAAFXC034AAAABdFJOUwBA5thmAAAAHUlEQVR42uzBAQ0AAADCoPdPbQ43oAAAAIAXE2AAD8gAAUGhJ40AAAAASUVORK5CYII=" ref="vcode" @click="reload" title="点击刷新" alt="点击刷新" height="40" />
						</Col>
					</Row>
				</FormItem>
				<FormItem :label-width="5">
					<Button size="large" @click="emitLogin" type="primary" long>登 录</Button>
				</FormItem>
			</Form>
			<p class="other-login">
				<a href="/login/qq.html" title="QQ登陆">
					<img :src="logins.qq" alt="加载失败" class="ivu-icon" style="width:30px">
				</a>
				<a @click="wechatLogin" title="微信登陆">
					<img :src="logins.wechat" alt="加载失败" class="ivu-icon" style="width:30px">
				</a>
			</p>
			<div class="login-tips">
				<span>系统推荐浏览器</span>&nbsp;<br><br>
				<a href="https://www.microsoft.com/zh-cn/edge" target="_blank" title="微软 Edge">
					<Icon type="logo-codepen" size="20" />&nbsp;Edge
				</a>
				<a href="https://www.google.cn/chrome/" target="_blank" title="谷歌 Chrome">
					<Icon type="logo-chrome" size="20" />&nbsp;Chrome
				</a>
				<a href="http://www.firefox.com.cn/" target="_blank" title="火狐 Firefox">
					<Icon type="ios-globe" size="20" />&nbsp;Firefox
				</a>
			</div>
		</div>
		<div class="footer">
			<div class="links">
				<a href="/">首页</a>
				<a>帮助</a>
				<a>隐私</a>
				<a>条款</a>
			</div>
			<div class="copyright">
				Copyright &copy; {{new Date().getFullYear()}}&nbsp;<a target="_blank" href="https://github.com/zxysilent">github.com/zxysilent</a>&nbsp;&nbsp;<a target="_blank" href="https://blog.zxysilent.com">blog.zxysilent.com</a>
			</div>
		</div>
	</div>
</template>
<script>
import md5 from "js-md5";
import loginQQ from "@/assets/login/qq.svg";
import loginWechat from "@/assets/login/wechat.svg";
import { apiAuthLogin, apiAuthVcode } from "@/api/auth";
import Storage from "@/utils/storage";
export default {
	data() {
		return {
			dataForm: { num: "", passwd: "", vcode: "", vreal: "" },
			logins: { qq: loginQQ, wechat: loginWechat },
			rules: {
				num: [{ required: true, message: "账号不能为空", trigger: "blur" }],
				passwd: [{ required: true, message: "密码不能为空", trigger: "blur" }],
				vcode: [{ required: true, message: "验证码不能为空", trigger: "blur" }]
			}
		};
	},
	methods: {
		wechatLogin() {
			this.$Message.warning("暂不支持");
		},
		emitLogin() {
			this.$refs.loginForm.validate((valid) => {
				if (valid) {
					let data = {
						num: this.dataForm.num,
						vcode: this.dataForm.vcode,
						vreal: this.dataForm.vreal,
						passwd: md5(this.dataForm.passwd).substr(1, 30)
					};
					apiAuthLogin(data).then((resp) => {
						if (resp.code == 200) {
							this.$Message.success({
								content: "登陆成功",
								onClose: () => {
									Storage.setToken(resp.data);
									this.$router.push({ name: "home" });
								}
							});
						} else {
							this.reload();
							this.$Message.error(resp.msg);
						}
					});
				}
			});
		},
		reload() {
			apiAuthVcode().then((resp) => {
				if (resp.code == 200) {
					this.$refs.vcode.src = "data:image/png;base64," + resp.data.vcode;
					this.dataForm.vcode = "";
					this.dataForm.vreal = resp.data.vreal;
				} else {
					this.dataForm.vcode = "";
					this.dataForm.vreal = "";
				}
			});
		}
	},
	mounted() {
		this.reload();
		Storage.clear();
	}
};
</script>
<style lang="less" scoped>
.login {
	width: 100%;
	min-height: 100%;
	background: #f0f2f5 url(../assets/bg.svg) no-repeat 50%;
	background-size: 100%;
	padding: 110px 0 144px;
	position: relative;
	a {
		text-decoration: none;
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
				font-family: "Helvetica Neue", Helvetica, "PingFang SC", "Hiragino Sans GB", "Microsoft YaHei",
					"微软雅黑", Arial, sans-serif;
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