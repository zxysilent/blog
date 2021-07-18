<style lang="less" scoped>
.control {
	cursor: pointer;
}
.control:hover {
	color: black;
}
</style>
<template>
	<div class="audiox">
		<Icon type="ios-play" title="播放" class="control" v-if="state" @click.native="onPlay" size="24" />
		<Icon type="ios-pause" title="暂停" class="control" v-else @click.native="onPause" size="24" />
		<audio :src="url" controls="controls" :id="'adx'+id" style="display: none;visibility: hidden;"></audio>
	</div>
</template>
<script>
export default {
	name: "audiox",
	data() {
		return { audios: null, state: true };
	},
	props: {
		id: Number,
		url: String
	},

	methods: {
		onPlay() {
			this.state = false;
			this.audios.play();
		},
		onPause() {
			this.state = true;
			this.audios.pause();
		}
	},
	mounted() {
		this.audios = document.getElementById("adx" + this.id);
	}
};
</script>