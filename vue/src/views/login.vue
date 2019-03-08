<style lang="less" scoped>
@import "./login.less";
</style>

<template>
	<div class="login" @keydown.enter="handleSubmit">
		<div class="login-ctn">
			<Card shadow>
				<p slot="title">
					<Icon type="md-log-in" />
					欢迎登录
				</p>
				<Form ref="loginForm" :model="user" :rules="rules">
					<FormItem prop="num" label="用户名">
						<Input prefix="ios-person-outline" type="text" v-model="user.num" placeholder="请输入用户名">
						</Input>
					</FormItem>
					<FormItem prop="pass" label="密码">
						<Input prefix="ios-lock-outline" type="password" v-model="user.pass" placeholder="请输入密码">
						</Input>
					</FormItem>
					<FormItem>
						<Button @click="handleSubmit" type="primary" long>登 录</Button>
					</FormItem>
				</Form>
				<p class="login-ctn-tip">首页
					<a href="/" title="首页">
						<Icon type="ios-send-outline" size="20" />
					</a>
				</p>
			</Card>
		</div>
	</div>
</template>
<script>
import md5 from "js-md5";
import { login } from "@/api/auth";
import util from "@/util/util";
export default {
	data() {
		return {
			user: {
				num: "zxysilent",
				pass: "zxyslt"
			},
			rules: {
				num: [{ required: true, message: "账号不能为空", trigger: "blur" }],
				pass: [{ required: true, message: "密码不能为空", trigger: "blur" }]
			}
		};
	},
	methods: {
		handleSubmit() {
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
