<template>
	<error code="50x" desc="Oh~~服务器休息了~" :src="src" />
</template>
<script>
import error50x from "@/assets/50x.svg";
import error from "./error.vue";
export default {
	name: "err50x",
	components: {
		error
	},
	data() {
		return {
			src: error50x
		};
	}
};
</script>
