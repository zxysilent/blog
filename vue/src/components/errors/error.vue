<style lang="less" scoped>
@import "./error.less";
</style>
<template>
	<div class="error">
		<div class="imgBlock">
			<div class="imgEle" :style="{backgroundImage: `url(${url})`}">
			</div>
		</div>
		<div class="content">
			<h1>{{code }}</h1>
			<div class="desc">{{ desc }}</div>
			<div class="actions">
				<Button type="success" @click="backHome">首 页</Button>&nbsp;&nbsp;
				<Button type="info" @click="backLogin">&nbsp;&nbsp;登 陆&nbsp;&nbsp;</Button>&nbsp;&nbsp;
				<Button type="warning" @click="backPrev">上一页</Button>
			</div>
		</div>
	</div>
</template>
<script>
export default {
	name: "error",
	data() {
		return {};
	},
	props: {
		code: String,
		desc: String,
		url: String
	},
	methods: {
		backHome() {
			this.$router.replace({ name: "home" });
		},
		backLogin() {
			localStorage.clear();
			sessionStorage.clear();
			this.$router.replace({ name: "login" });
		},
		backPrev() {
			this.$router.go(-1);
		}
	}
};
</script>