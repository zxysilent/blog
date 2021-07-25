<template>
	<div>
		<textarea :id="'zxyslt-'+id">{{ content }}</textarea>
	</div>
</template>

<script>
// Import TinyMCE
import tinymce from "tinymce/tinymce";
import "tinymce/icons/default";
import "tinymce/themes/silver";
import "tinymce/plugins/advlist";
import "tinymce/plugins/codesample";
import "tinymce/plugins/hr"; //横线
import "tinymce/plugins/media"; //插入媒体文件
import "tinymce/plugins/save"; //保存
import "tinymce/plugins/toc"; //插入目录
import "tinymce/plugins/anchor"; //插入锚点
//import "tinymce/plugins/colorpicker"
import "tinymce/plugins/image"; //添加图片
import "tinymce/plugins/nonbreaking"; //插入一个nonbreaking
import "tinymce/plugins/searchreplace"; //打开搜索和替换对话框
import "tinymce/plugins/visualblocks"; //隐藏块级区域开关
import "tinymce/plugins/autolink";
//import "tinymce/plugins/contextmenu"
import "tinymce/plugins/imagetools"; //图像操作
import "tinymce/plugins/noneditable";
import "tinymce/plugins/spellchecker"; //拼写检查开关
import "tinymce/plugins/visualchars"; //隐藏字符串开关
// import "tinymce/plugins/autoresize"; //高度随内容自变化
import "tinymce/plugins/directionality"; //设置编写方向
import "tinymce/plugins/pagebreak"; //分页符
import "tinymce/plugins/tabfocus";
// import "tinymce/plugins/wordcount";
//import "tinymce/plugins/autosave"    //恢复到上次自动保存的数据
//import "tinymce/plugins/emoticons"       //表情
import "tinymce/plugins/insertdatetime"; //插入时间
import "tinymce/plugins/paste"; //粘贴
import "tinymce/plugins/table"; //表格
// import "tinymce/plugins/bbcode";//类markdown
// import "tinymce/plugins/fullpage"; //查看完整的文档属性
import "tinymce/plugins/legacyoutput";
// import "tinymce/plugins/preview"; //预览
// import "tinymce/plugins/template"; //插入模板
import "tinymce/plugins/charmap"; //开启字符映射表.
import "tinymce/plugins/fullscreen"; //全屏
import "tinymce/plugins/link"; //添加一个连接
// import "tinymce/plugins/print"; //打印
//import "tinymce/plugins/textcolor"//颜色（文本、背景色）
import "tinymce/plugins/code"; //源码预览
import "tinymce/plugins/help"; //帮助菜单
import "tinymce/plugins/lists"; //列表
// import "tinymce/plugins/quickbars";
import "tinymce/plugins/textpattern";
import "@/assets/zh_CN";
export default {
	name: "tiny-editor",
	props: {
		id: { type: String, required: true },
		value: { default: "" },
		height: {
			type: String,
			default: "calc(100vh - 110px)"
		}
	},
	data() {
		return { content: "", editor: null };
	},
	mounted() {
		this.content = this.value;
		this.init();
	},
	beforeDestroy: function() {
		//	this.editor.remove();
		this.editor.destroy();
	},

	watch: {
		value: function(newValue) {
			console.log(newValue);
			if (this.value == newValue) return;
			if (this.editor) {
				this.editor.setContent(newValue);
			} else {
				this.content = newValue;
			}
		}
	},
	methods: {
		emitResetor() {
			console.log("emitResetor");
			this.editor.setContent("");
		},
		setContent(ctt) {
			console.log(ctt);
			this.editor.setContent(ctt);
		},
		init() {
			tinymce.init({
				selector: "#zxyslt-" + this.id,
				skin: "oxide",
				skin_url: "/static/css/oxide",
				statusbar: false,
				image_advtab: true,
				menubar: "file edit insert view format table",
				height: this.height, // "calc(100vh - 110px)",
				// min_height: this.height,
				branding: false,
				language: "zh_CN",
				toolbar_items_size: "small",
				nonbreaking_force_tab: true,
				content_css: false,
				forced_root_block: "div",
				paste_retain_style_properties: "all",
				paste_word_valid_elements: "*[*]",
				paste_data_images: true,
				paste_convert_word_fake_lists: false,
				paste_webkit_styles: "all",
				paste_merge_formats: true,
				nonbreaking_force_tab: false,
				paste_auto_cleanup_on_paste: false,
				font_formats: `
                    微软雅黑=微软雅黑;
                    宋体=宋体;
                    黑体=黑体;
                    仿宋=仿宋;
                    楷体=楷体;
                    隶书=隶书;
                    幼圆=幼圆;
                    Andale Mono=andale mono,times;
                    Arial=arial, helvetica,
                    sans-serif;
                    Arial Black=arial black, avant garde;
                    Book Antiqua=book antiqua,palatino;
                    Comic Sans MS=comic sans ms,sans-serif;
                    Courier New=courier new,courier;
                    Georgia=georgia,palatino;
                    Helvetica=helvetica;
                    Impact=impact,chicago;
                    Symbol=symbol;
                    Tahoma=tahoma,arial,helvetica,sans-serif;
                    Terminal=terminal,monaco;
                    Times New Roman=times new roman,times;
                    Trebuchet MS=trebuchet ms,geneva;
                    Verdana=verdana,geneva;
                    Webdings=webdings;
                    Wingdings=wingdings,zapf dingbats`,
				plugins: [
					"advlist",
					"codesample",
					"hr",
					"media",
					"save",
					"toc",
					"anchor",
					// "colorpicker",
					"image",
					"nonbreaking",
					"searchreplace",
					"visualblocks",
					"autolink",
					// "contextmenu",
					"imagetools",
					"noneditable",
					"spellchecker",
					"visualchars",
					// "autoresize",
					"directionality",
					"pagebreak",
					"tabfocus",
					// "wordcount",
					// "autosave",
					// "emoticons",
					"insertdatetime",
					"paste",
					"table",
					// "bbcode",
					// "fullpage",
					"legacyoutput",
					// "preview",
					// "template",
					"charmap",
					"fullscreen",
					"link",
					// "print",
					// "textcolor",
					"code",
					"help",
					"lists",
					// "quickbars",
					"textpattern"
				],
				toolbar:
					"fontselect | fontsizeselect | bold italic underline strikethrough | forecolor backcolor | image  link | alignleft aligncenter alignright alignjustify | outdent indent | undo redo | fullscreen",
				images_upload_handler: function(blobInfo, success, failure) {
					var xhr, formData;
					xhr = new XMLHttpRequest();
					xhr.withCredentials = false;
					xhr.open("POST", process.env.VUE_APP_SRV + "/api/upload/image");
					xhr.onload = function() {
						var json;
						if (xhr.status < 200 || xhr.status >= 300) {
							failure("图片上传失败,HTTP Error: " + xhr.status);
							return;
						}
						json = JSON.parse(xhr.responseText);
						if (!json || json.code != 200) {
							failure("图片上传失败");
							return;
						}
						success(json.data);
					};
					formData = new FormData();
					formData.append("file", blobInfo.blob());
					formData.set("cut", true);
					formData.set("wd", 1100);
					xhr.send(formData);
				},
				setup: editor => {
					this.editor = editor;
					/* Helper functions */
					// var toDateHtml = function(date) {
					// 	return '<time datetime="' + date.toString() + '">' + date.toDateString() + "</time>";
					// };

					/* Basic button that just inserts the date */
					// editor.ui.registry.addButton("Pdf", {
					// 	text: "Pdf",
					// 	tooltip: "上传pdf",
					// 	onAction: function(_) {
					// 		editor.insertContent(toDateHtml(new Date()));
					// 	}
					// });
					editor.on("init", () => {
						editor.setContent(this.value);
					});
					editor.on("change keyup undo redo", () => {
						this.$emit("input", editor.getContent());
					});
					editor.on("NodeChange", () => {
						this.$emit("input", editor.getContent());
					});
				},
				file_picker_types: "file media",
				file_picker_callback: function(callback, value, meta) {
					// // Provide file and text for the link dialog
					let fileAccept = "*/*";
					// if (meta.filetype == "file") {
					// 	callback("mypage.html", { text: "My text" });
					// }

					// // Provide image and alt text for the image dialog
					if (meta.filetype == "image") {
						// callback("myimage.jpg", { alt: "My alt text" });
						fileAccept = "image/*";
					}

					// // Provide alternative source and posted for the media dialog
					if (meta.filetype == "media") {
						// callback("movie.mp4", { source2: "alt.ogg", poster: "image.jpg" });
						fileAccept = "video/*";
					}
					let input = document.createElement("input");
					input.setAttribute("type", "file");
					input.setAttribute("accept", fileAccept);
					input.onchange = function() {
						let file = this.files[0];
						let xhr, formData;
						xhr = new XMLHttpRequest();
						xhr.withCredentials = false;
						xhr.open("POST", process.env.VUE_APP_SRV + "/api/upload/file");
						xhr.onload = function() {
							let json;
							if (xhr.status < 200 || xhr.status >= 300) {
								//failure("上传失败,HTTP Error: " + xhr.status);
								return;
							}
							json = JSON.parse(xhr.responseText);
							if (!json || json.code != 200) {
								//failure("上传失败");
								return;
							}
							callback(json.data, {
								text: "附件 " + file.name,
								title: "点击下载"
							});
						};
						formData = new FormData();
						formData.append("file", file);
						xhr.send(formData);
					};
					input.click();
				}
			});
		}
	}
};
</script>
<style scoped>
</style>