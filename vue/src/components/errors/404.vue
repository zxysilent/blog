<template>
	<error code="404" desc="Oh~~您的页面好像飞走了~" :url="url" />
</template>
<script>
import url404 from "@/assets/404.svg";
import error from "./error.vue";
export default {
	name: "E404",
	components: {
		error
	},
	data() {
		return {
			url: url404
		};
	}
};
</script>
