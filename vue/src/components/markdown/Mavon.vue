<style lang="less">
.v-note-wrapper {
	z-index: 999 !important;
}
.v-note-wrapper .v-note-panel .v-note-edit.divarea-wrapper .content-input-wrapper {
	padding: 4px 10px 10px 10px !important;
}
.markdown-body {
	font-size: 14px !important;
}

.v-note-wrapper .v-note-panel .v-note-show .v-show-content,
.v-note-wrapper .v-note-panel .v-note-show .v-show-content-html {
	padding: 4px 10px 10px 10px !important;
}

.fullscreen {
	z-index: 1001 !important;
}
.add-image-link-wrapper .add-image-link {
	padding: 20px !important;
}
[type="reset"],
[type="submit"],
button,
html [type="button"] {
	-webkit-appearance: none !important;
}
</style>
<template>
	<div>
		<div style="min-height: 600px;height: calc(100vh - 145px);">
			<button ref="diy" type="button" @click="diyMore" class="op-icon ivu-icon ivu-icon-md-code" aria-hidden="true" title="更多"></button>
			<mavon-editor ref="md" @imgAdd="onImageAdd" :boxShadow="false" @change="onChange" :toolbars="toolbars" v-model="markdown" style="height:100%">
			</mavon-editor>
		</div>
	</div>
</template>
<script>
import { mavonEditor } from "mavon-editor";
import "mavon-editor/dist/css/index.css";
import Storage from "@/utils/storage";
// 减少js体积
export default {
	name: "mavon",
	components: {
		mavonEditor
	},
	props: {
		value: { type: String, required: true }
	},
	data() {
		return {
			richtext: "",
			markdown: "",
			toolbars: {
				bold: true, // 粗体
				italic: true, // 斜体
				header: true, // 标题
				underline: true, // 下划线
				strikethrough: true, // 中划线
				mark: true, // 标记
				superscript: true, // 上角标
				subscript: true, // 下角标
				quote: true, // 引用
				ol: true, // 有序列表
				ul: true, // 无序列表
				link: true, // 链接
				imagelink: true, // 图片链接
				code: true, // code
				table: true, // 表格
				fullscreen: true, // 全屏编辑
				//readmodel: true, // 沉浸式阅读
				htmlcode: true, // 展示html源码
				// help: true, // 帮助
				/* 1.3.5 */
				undo: true, // 上一步
				redo: true, // 下一步
				trash: true, // 清空
				// save: true, // 保存（触发events中的save事件）
				/* 1.4.2 */
				// navigation: true, // 导航目录
				/* 2.1.8 */
				// alignleft: true, // 左对齐
				// aligncenter: true, // 居中
				// alignright: true, // 右对齐
				/* 2.2.1 */
				subfield: true, // 单双栏模式
				preview: true // 预览
			}
		};
	},
	watch: {
		value: function (newValue) {
			// console.log(newValue);
			this.markdown = newValue;
		}
	},
	methods: {
		onChange(value, html) {
			// console.log(value, html);
			this.richtext = html;
			this.$emit("input", value);
			this.$emit("on-change", html);
		},
		diyMore(val) {
			this.$refs.md.insertText(this.$refs.md.getTextareaDom(), {
				prefix: "",
				subfix: "",
				str: "<!--more-->"
			});
		},
		onImageAdd(pos, $file) {
			let formData = new FormData();
			let xhr = new XMLHttpRequest();
			xhr.withCredentials = false;
			xhr.open("POST", process.env.VUE_APP_SRV + "/adm/upload/image");
			xhr.onload = () => {
				var json;
				if (xhr.status < 200 || xhr.status >= 300) {
					this.$Message.warning("图片上传失败,HTTP Error: " + xhr.status);
					return;
				}
				json = JSON.parse(xhr.responseText);
				if (!json || json.code != 200) {
					this.$Message.warning("图片上传失败");
					return;
				}
				this.$refs.md.$img2Url(pos, process.env.VUE_APP_SRV + json.data);
			};
			formData.append("token", Storage.getToken());
			formData.append("file", $file);
			xhr.send(formData);
		}
	},
	mounted() {
		var md = this.$refs.md;
		var toolbarLeft = md.$refs.toolbar_left;
		var diy = this.$refs.diy;
		toolbarLeft.$el.append(diy);
		this.markdown = this.value;
	}
};
</script>
