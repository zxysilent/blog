
<style lang="less" scoped>
.uimage {
	display: inline-block;
	width: 60px;
	height: 60px;
	text-align: center;
	line-height: 60px;
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
                 <Tag color="primary">已上传</Tag>
				<div class="uimage-cover">
					<Icon type="ios-trash-outline" color="#FF5722" title="删除" @click.native="emitDrop(value)"></Icon>
				</div>
			</template>
			<template v-else>
				<Upload :show-upload-list="false" :on-success="upAfter" :accept="accept" :max-size="maxSize" :on-format-error="errorFormat" :on-exceeded-size="errorSize" type="drag" :action="upAction" style="display: inline-block;width:58px;">
					<div style="width: 58px;height:58px;line-height: 58px;">
						<Icon type="ios-cloud-upload" size="20" title="上传"></Icon>
					</div>
				</Upload>
			</template>
		</div>
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
			default: 8
		},
		accept: {
			type: String,
			default: "*/*"
		}
	},
	data() {
		return {};
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
			return process.env.VUE_APP_SRV + "/api/upload/file";
		}
	},
	methods: {
		emitDrop(val) {
			this.$emit("input", "");
			this.url = "";
		},
		errorFormat(file) {
			this.$Message.warning("请选择正确的文件格式");
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