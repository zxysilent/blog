<style lang="less" scoped>
@import "./error.less";
</style>
<template>
	<div class="error">
		<div class="content">
			<img src="../../assets/jwt.svg" alt="登陆检测">
			<div class="text">
				<h4>jwt</h4>
				<h5>Oh~~嘿呀要重新登陆啦~</h5>
			</div>
			<div class="btn-group">
				<Button size="large" type="info" @click="goLogin">去登陆</Button>
			</div>
		</div>
	</div>
</template>
<script>
export default {
	name: "errjwt",
	data() {
		return {};
	},
	methods: {
		goLogin() {
			this.$router.replace({
				name: "login"
			});
		}
	}
};
</script>