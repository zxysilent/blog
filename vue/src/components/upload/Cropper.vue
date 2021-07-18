<style lang="less" scoped>
.ivu-btn-icon-only {
	width: 45px;
}
.bg {
	background-image: url("data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQAQMAAAAlPW0iAAAAA3NCSVQICAjb4U/gAAAABlBMVEXMzMz////TjRV2AAAACXBIWXMAAArrAAAK6wGCiw1aAAAAHHRFWHRTb2Z0d2FyZQBBZG9iZSBGaXJld29ya3MgQ1M26LyyjAAAABFJREFUCJlj+M/AgBVhF/0PAH6/D/HkDxOGAAAAAElFTkSuQmCC");
}
.cropper-wrapper {
	width: 430px;
	height: 250px;
	.img-box {
		height: 250px;
		width: 310px;
		border: 1px solid #ebebeb;
		display: inline-block;
		.bg;
		img {
			max-width: 100%;
			display: block;
		}
	}
	.right-con {
		display: inline-block;
		width: 120px;
		vertical-align: top;
		box-sizing: border-box;
		padding: 0 5px;
	}
}
</style>
<template>
	<div class="cropper-wrapper">
		<div class="img-box">
			<img class="cropper-image" :id="imgId">
		</div>
		<div class="right-con">
			<Row>
				<Col span="12">
				<Upload action="image/upload" accept="image/*" :show-upload-list="false" :before-upload="beforeUpload">
					<Button title="选择图片" icon="md-cloud-upload"></Button>
				</Upload>
				</Col>
				<Col span="12">
				<Button title="删除图片" icon="md-close" @click="resetCrop"></Button></Col>
			</Row>
			<div v-show="insideSrc">
				<Input style="width:100px;margin-top:1px;" size="small" readonly v-model="wd">
				<span slot="prepend">宽</span>
				</Input>
				<Input style="width:100px;margin-top:1px;" size="small" readonly v-model="hd">
				<span slot="prepend">高</span>
				</Input>
				<Row style="margin-top:1px;">
					<Col span="12">
					<Button title="旋转" icon="md-refresh" @click="rotate"></Button>
					</Col>
					<Col span="12">
					<Button title="缩小" icon="md-remove" @click="shrink"></Button>
					</Col>
				</Row>
				<Row style="margin-top:1px;">
					<Col span="12">
					<Button title="放大"icon="md-add" @click="magnify"></Button>
					</Col>
					<Col span="12">
					<Button title="还原" icon="md-sync" @click="reset"></Button>
					</Col>
				</Row>
				<Row style="margin-top:1px;">
					<Col span="12">
					<Button title="上移" icon="md-arrow-round-up" @click="move(0, -10)"></Button>
					</Col>
					<Col span="12">
					<Button title="下移" icon="md-arrow-round-down" @click="move(0, 10)"></Button>
					</Col>
				</Row>
				<Row style="margin-top:1px;">
					<Col span="12">
					<Button title="左移" icon="md-arrow-round-back" @click="move(-10, 0)"></Button>
					</Col>
					<Col span="12">
					<Button title="右移" icon="md-arrow-round-forward" @click="move(10, 0)"></Button>
					</Col>
				</Row>

				<Row style="margin-top:1px;">
					<Col span="24"><Button style="width: 100px;margin-top: 1px;" @click="crop" type="info">确认上传</Button></Col>
				</Row>
			</div>
		</div>
	</div>
</template>
<script>
import Cropper from "cropperjs";
import "cropperjs/dist/cropper.min.css";
export default {
	name: "cropper",
	props: {
		src: {
			type: String,
			default: ""
		}
	},
	data() {
		return {
			wd: 0,
			hd: 0,
			cropper: null,
			insideSrc: ""
		};
	},
	computed: {
		imgId() {
			return `cropper${this._uid}`;
		}
	},
	watch: {
		src(src) {
			this.replace(src);
		},
		insideSrc(src) {
			this.replace(src);
		}
	},
	methods: {
		beforeUpload(file) {
			const reader = new FileReader();
			reader.readAsDataURL(file);
			reader.onload = event => {
				this.insideSrc = event.srcElement.result;
			};
			return false;
		},
		resetCrop() {
			this.replace("");
			this.cropper.destroy();
			this.$emit("on-reset");
			this.init();
		},
		replace(src) {
			this.cropper.replace(src);
			this.insideSrc = src;
		},
		rotate() {
			this.cropper.rotate(90);
		},
		shrink() {
			this.cropper.zoom(-0.1);
		},
		magnify() {
			this.cropper.zoom(0.1);
		},
		reset() {
			this.cropper.reset();
		},
		move(...argu) {
			this.cropper.move(...argu);
		},
		crop() {
			this.cropper.getCroppedCanvas().toBlob(blob => {
				this.$emit("on-crop", blob);
			});
		},
		init() {
			let that = this;
			let dom = document.getElementById(this.imgId);
			this.cropper = new Cropper(dom, {
				checkCrossOrigin: true,
				autoCropArea: 1,
				crop(event) {
					that.wd = event.detail.width.toFixed(1);
					that.hd = event.detail.height.toFixed(1);
				}
			});
		}
	},
	beforeDestroy: function() {
		this.cropper.destroy();
	},
	mounted() {
		this.$nextTick(() => {
			this.init();
		});
	}
};
</script>