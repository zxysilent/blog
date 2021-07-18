<template>
	<Upload :before-upload="upBefore" :max-size="maxSize" :show-upload-list="false" :data="upData" :on-success="upAfter" :on-error="upError" :on-exceeded-size="errorSize" :on-format-error="errorFormat" :accept="accept" :action="upAction">
		<Button icon="ios-cloud-upload-outline" type="dashed">{{ctx}}</Button>
		<Poptip :content="msg" v-model="show" placement="right">&nbsp;</Poptip>
	</Upload>
</template>
<script>
export default {
	name: "upimage",
	props: {
		//url
		value: {
			type: String
		},
		//文件限制大小
		size: {
			type: Number,
			default: 6
		},
		accept: {
			type: String,
			default: "image/*"
		},
		content: {
			type: String,
			default: "请选择图片"
		}
	},
	data() {
		return {
			msg: "",
			ctx: this.content,
			show: false
		};
	},
	computed: {
		tip() {
			return `大小不能超过${this.size}M`;
		},
		maxSize() {
			return this.size * 1024;
		},
		upData() {
			return {
				ver: 1009
			};
		},
		upAction() {
			return process.env.VUE_APP_SRV + "/api/upload/image?cut=true";
		}
	},
	watch: {
		show(o, n) {
			setTimeout(() => {
				this.show = false;
			}, 3000);
			console.log(o, n);
		}
	},
	methods: {
		upBefore(file) {
			this.msg = "上传中...";
			this.show = true;
			return true;
		},
		upError() {
			this.msg = "上传失败,请重试";
			this.show = true;
		},
		errorSize() {
			this.msg = `请选择文件小于${this.size}M的文件`;
			this.show = true;
		},
		errorFormat() {
			this.msg = "请选择正确的文件格式";
			this.show = true;
		},
		upAfter(resp) {
			this.showSpin = false;
			if (resp.code == 200) {
				this.$emit("input", resp.data);
				this.msg = "上传成功";
				this.ctx = "上传成功";
				this.show = true;
			} else {
				this.msg = "上传失败,请重试";
				this.ctx = "上传失败,请重试";
				this.show = true;
			}
		}
	},
	mounted() {}
};
</script>