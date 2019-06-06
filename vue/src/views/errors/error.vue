<style lang="less" scoped>
@import "./error.less";
</style>
<template>
	<div class="error">
		<div class="content">
			<img :src="src" :alt="code">
			<div class="text">
				<h4>{{ code }}</h4>
				<h5>{{ desc }}</h5>
			</div>
			<div class="btn-group">
				<Button size="large" type="success" @click="backHome">返回首页</Button>&nbsp;&nbsp;
				<Button size="large" type="info" @click="backLogin">&nbsp;&nbsp;去登陆&nbsp;&nbsp;</Button>&nbsp;&nbsp;
				<Button size="large" type="warning" @click="backPrev">返回上一页</Button>
			</div>
		</div>
	</div>
</template>
<script>
export default {
	name: "error",
	data() {
		return {
			second: 5,
			timer: null
		};
	},
	props: {
		code: String,
		desc: String,
		src: String
	},
	methods: {
		backHome() {
			this.$router.replace({
				name: "home"
			});
		},
		backLogin() {
			localStorage.clear();
			sessionStorage.clear();
			this.$router.replace({
				name: "login"
			});
		},
		backPrev() {
			this.$router.go(-1);
		}
	},
	mounted() {
		// this.timer = setInterval(() => {
		// 	if (this.second === 0) this.backPrev();
		// 	else this.second--;
		// }, 1000);
	},
	beforeDestroy() {
		// clearInterval(this.timer);
	}
};
</script>