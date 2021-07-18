<template>
	<error code="500" desc="Oh~~服务器休息了~" :url="url" />
</template>
<script>
import url500 from "@/assets/500.svg";
import error from "./error.vue";
export default {
	name: "E500",
	components: {
		error
	},
	data() {
		return {
			url: url500
		};
	}
};
</script>
