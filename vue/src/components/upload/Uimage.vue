
<style lang="less" scoped>
.uimage {
	display: inline-block;
	width: 102px;
	height: 102px;
	text-align: center;
	line-height: 102px;
	border-radius: 4px;
	overflow: hidden;
	background: #fff;
	position: relative;
	margin-right: 4px;
}
.uimage img {
	width: 100%;
	height: 100%;
}
.uimage-cover {
	display: none;
	position: absolute;
	top: 0;
	bottom: 0;
	left: 0;
	right: 0;
	background: rgba(0, 0, 0, 0.6);
}
.uimage:hover .uimage-cover {
	display: block;
}
.uimage-cover i {
	color: #f3f3f3;
	font-size: 25px;
	cursor: pointer;
	margin: 0 5px;
}
</style>
<template>
	<div>
		<div class="uimage">
			<template v-if="fullurl != ''">
				<img :src="fullurl" alt="加载失败">
				<div class="uimage-cover">
					<Icon type="ios-eye-outline" title="预览" color="#5FB878" @click.native="emitPreview(fullurl)"></Icon>
					<Icon type="ios-trash-outline" title="删除" color="#FF5722" @click.native="emitDrop(value)"></Icon>
				</div>
			</template>
			<template v-else>
				<Upload :show-upload-list="false" :on-success="upAfter" :accept="accept" :max-size="maxSize" :on-format-error="errorFormat" :on-exceeded-size="errorSize" type="drag" :action="upAction" style="display: inline-block;width:100px;">
					<div style="width: 100px;height:100px;line-height: 100px;">
						<Icon type="ios-camera" size="30" title="上传"></Icon>
					</div>
				</Upload>
			</template>
		</div>
		<Modal v-model="preview">
			<p slot="header" style="color:#515a6e;text-align:center">
				<Icon type="ios-eye-outline" />
				<span>预览查看</span>
			</p>
			<div style="text-align:center">
				<img :src="fullurl" style="width: 100%;max-height: 600px">
			</div>
			<div slot="footer">
				<Button type="info" @click="preview=false">确&nbsp;定</Button>
			</div>
		</Modal>
	</div>
</template>
<script>
export default {
	name: "uimage",
	props: {
		//url
		value: {
			type: String,
			default: ""
		},
		//文件限制大小
		size: {
			type: Number,
			default: 4
		},
		accept: {
			type: String,
			default: "image/png,image/bmp,image/jpeg"
		}
	},
	data() {
		return {
			preview: false
		};
	},
	computed: {
		maxSize() {
			return this.size * 1024;
		},
		upData() {
			return {
				ver: 1009
			};
		},
		fullurl() {
			console.log(this.value);
			if (this.value.indexOf("/") == 0) {
				return process.env.VUE_APP_SRV + this.value;
			}
			return this.value;
		},
		upAction() {
			return process.env.VUE_APP_SRV + "/api/upload/image";
		}
	},
	methods: {
		emitPreview(name) {
			this.preview = true;
		},
		emitDrop(val) {
			this.$emit("input", "");
			this.url = "";
		},
		errorFormat(file) {
			this.$Message.warning("请选择图片格式：.png、.jpg、.bmp");
		},
		errorSize(file) {
			this.$Message.warning(`请选择小于${this.size}M的图片`);
		},
		upAfter(resp) {
			if (resp.code == 200) {
				this.$emit("input", resp.data);
				// this.url = resp.data;
			} else {
				this.$emit("input", "");
				// this.url = "";
				this.$Message.error(`上传失败：${resp.msg}`);
			}
		}
	},
	mounted() {
		console.log(this.url);
	}
};
</script>