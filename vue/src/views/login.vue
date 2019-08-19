<style lang="less" scoped>
@import "./login.less";
</style>
<template>
	<div class="login" @keydown.enter="submit">
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
			<Form ref="loginForm" :model="user" :rules="rules">
				<FormItem prop="num" label="用户名">
					<Input size="large" prefix="ios-person-outline" type="text" v-model="user.num" placeholder="请输入用户名">
					</Input>
				</FormItem>
				<FormItem prop="pass" label="密码">
					<Input size="large" prefix="ios-lock-outline" type="password" v-model="user.pass" placeholder="请输入密码">
					</Input>
				</FormItem>
				<FormItem>
					<Button size="large" @click="submit" type="primary" long>登 录</Button>
				</FormItem>
			</Form>
			<p class="login-tip">首页
				<a href="/" title="首页">
					<Icon type="ios-send-outline" size="20" />
				</a>
			</p>
			<div class="other-login">
				<span>其他登录方式</span>&nbsp;
				<a>
					<Icon type="logo-github" />
				</a>
			</div>
		</div>
		<div class="footer">
			<div class="links">
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
import { login } from "@/api/auth";
import util from "@/init/util";
export default {
	data() {
		return {
			user: { num: "zxysilent", pass: "zxyslt" },
			rules: {
				num: [{ required: true, message: "账号不能为空", trigger: "blur" }],
				pass: [{ required: true, message: "密码不能为空", trigger: "blur" }]
			}
		};
	},
	methods: {
		submit() {
			let that = this;
			that.$refs.loginForm.validate(valid => {
				if (valid) {
					let data = {
						num: that.user.num,
						pass: md5(that.user.pass).substr(1, 30)
					};
					login(data).then(res => {
						if (res.code == 200) {
							this.$Message.success({
								content: "登陆成功",
								onClose: () => {
									util.setToken(res.data);
									that.$router.push({ name: "home" });
								}
							});
						} else {
							this.$Message.error(res.msg);
						}
					});
				}
			});
		}
	}
};
</script>