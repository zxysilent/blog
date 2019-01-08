<template>
	<div id="main" class="app-main">
		<router-view></router-view>
	</div>
</template>

<script>
export default {
	mounted() {},
	beforeDestroy() {},
	methods: {}
};
</script>

<style>
html,
body {
	width: 100%;
	height: 100%;
	background: #f0f0f0;
	overflow: auto;
}
.app-main {
	width: 100%;
	height: 100%;
	min-width: 1000px;
}
</style>
