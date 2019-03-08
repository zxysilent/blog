<template>
	<error code="40x" desc="Oh~~嘿呀没有权限啦~" :src="src" />
</template>

<script>
import error40x from "@/assets/40x.svg";
import error from "./error.vue";
export default {
	name: "err40x",
	components: {
		error
	},
	data() {
		return {
			src: error40x
		};
	}
};
</script>
