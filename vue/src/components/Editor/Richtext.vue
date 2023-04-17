<template>
    <div class="tinymce-box">
        <textarea :id="'mce' + props.id">{{ content }}</textarea>
    </div>
</template>
<script lang="ts" setup>
// Import TinyMCE
import tinymce from "tinymce/tinymce";
import "tinymce/themes/silver";
import "tinymce/models/dom";
import "tinymce/icons/default";
// plugins
import "tinymce/plugins/advlist";
import "tinymce/plugins/anchor";
import "tinymce/plugins/autolink";
// import "tinymce/plugins/autoresize";
// import "tinymce/plugins/autosave";
import "tinymce/plugins/charmap";
import "tinymce/plugins/code";
import "tinymce/plugins/codesample";
import "tinymce/plugins/directionality";
// import "tinymce/plugins/emoticons";
import "tinymce/plugins/fullscreen";
// import "tinymce/plugins/help";
import "tinymce/plugins/image";
import "tinymce/plugins/importcss";
import "tinymce/plugins/insertdatetime";
import "tinymce/plugins/link";
import "tinymce/plugins/lists";
import "tinymce/plugins/media";
import "tinymce/plugins/nonbreaking";
import "tinymce/plugins/pagebreak";
import "tinymce/plugins/preview";
// import "tinymce/plugins/quickbars";
import "tinymce/plugins/save";
import "tinymce/plugins/searchreplace";
import "tinymce/plugins/table";
// import "tinymce/plugins/template";
import "tinymce/plugins/visualblocks";
import "tinymce/plugins/visualchars";
import "tinymce/plugins/wordcount";
import "@/assets/zh-Hans";
import { ref, onMounted, onUnmounted } from "vue";
import storage from "@/utils/storage";
const imageActionURL = "/api/upload/image";
const fileActionURL = "/api/upload/file";
const emit = defineEmits(["update:value"]);
const props = defineProps({
    id: { type: String, default: new Date().getTime() + "" },
    value: { type: String, required: true },
    size: { type: Number, default: 8 },
    onChange: { type: Function, default: undefined },
    fileListStyle: { type: Object, default: undefined },
    height: {
        type: String,
        default: "calc(100vh - 110px)",
    },
});
const content = ref<string>(props.value);
const editorRef = ref();
const images_upload_handler = (blobInfo, progress) =>
    new Promise<string>((resolve, reject) => {
        const xhr = new XMLHttpRequest();
        xhr.withCredentials = false;
        xhr.open("POST", imageActionURL);
        xhr.upload.onprogress = (e) => {
            progress((e.loaded / e.total) * 100);
        };
        xhr.onload = () => {
            if (xhr.status === 403) {
                reject({ message: "图片上传失败 HTTP Error: " + xhr.status, remove: true });
                return;
            }
            if (xhr.status < 200 || xhr.status >= 300) {
                reject({ message: "图片上传失败 HTTP Error: " + xhr.status, remove: true });
                return;
            }
            const json = JSON.parse(xhr.responseText);
            if (!json || json.code != 200) {
                reject({ message: "图片上传失败 HTTP Error: " + xhr.status, remove: true });
                return;
            }
            resolve(json.data);
        };

        xhr.onerror = () => {
            reject("Image upload failed due to a XHR Transport error. Code: " + xhr.status);
        };
        const formData = new FormData();
        formData.append("file", blobInfo.blob());
        formData.append("token", storage.getToken());
        xhr.send(formData);
    });
const file_picker_callback = (callback, value, meta) => {
    console.log("meta", value, meta);
    let fileAccept = "*/*";
    if (meta.filetype == "image") {
        fileAccept = "image/*";
    }
    if (meta.filetype == "media") {
        fileAccept = "video/*"; // callback("movie.mp4", { source2: "alt.ogg", poster: "image.jpg" });
    }
    const input = document.createElement("input");
    input.setAttribute("type", "file");
    input.setAttribute("accept", fileAccept);
    input.onchange = () => {
        //@ts-ignore
        const file = input.files[0];
        const xhr = new XMLHttpRequest();
        xhr.withCredentials = false;
        xhr.open("POST", fileActionURL);
        xhr.onload = function () {
            if (xhr.status < 200 || xhr.status >= 300) {
                return;
            }
            const json = JSON.parse(xhr.responseText);
            if (!json || json.code != 200) {
                return;
            }
            callback(json.data, {
                text: "附件 " + file.name,
                title: "点击下载",
            });
        };
        const formData = new FormData();
        formData.append("file", file);
        formData.append("token", storage.getToken());
        xhr.send(formData);
    };
    input.click();
};
const init = () => {
    tinymce
        .init({
            selector: "#mce" + props.id,
            promotion: false, //升级插件提示
            skin: "oxide",
            skin_url: "/static/css/oxide",
            statusbar: false, //输入状态
            image_advtab: true,
            menubar: "edit insert view format table",
            height: props.height, // "calc(100vh - 110px)",
            branding: false,
            language: "zh-Hans",
            toolbar_items_size: "small",
            nonbreaking_force_tab: true,
            content_css: false,
            forced_root_block: "div",
            paste_data_images: true,
            paste_webkit_styles: "all",
            paste_merge_formats: true,
            elementpath: false, //元素路径是否显示
            custom_undo_redo_levels: 10, //撤销和重做的次数
            draggable_modal: true, //对话框允许拖拽
            relative_urls: false,
            // remove_script_host: false,
            // nonbreaking_force_tab: false,
            paste_auto_cleanup_on_paste: false,
            font_family_formats: `
                    微软雅黑=微软雅黑;
                    宋体=宋体;
                    黑体=黑体;
                    仿宋=仿宋;
                    楷体=楷体;
                    隶书=隶书;
                    幼圆=幼圆;
                    Andale Mono=andale mono,times;
                    Arial=arial, helvetica, sans-serif;
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
            image_title: true,
            /* enable automatic uploads of images represented by blob or data URIs*/
            automatic_uploads: true,
            plugins: [
                "advlist",
                "anchor",
                "autolink",
                "charmap",
                "code",
                "codesample",
                "directionality",
                "fullscreen",
                "image",
                "importcss",
                "insertdatetime",
                "link",
                "lists",
                "media",
                "nonbreaking",
                "pagebreak",
                "preview",
                "save",
                "searchreplace",
                "table",
                // "template",
                "visualblocks",
                "visualchars",
                "wordcount",
                // "autoresize",
                // "autosave",
                // "emoticons",
                // "help",
                // "quickbars",
            ],
            //fontfamily blocks fontsize
            toolbar: "bold italic underline strikethrough | forecolor backcolor | image  link | alignleft aligncenter alignright alignjustify  | outdent indent undo redo | preview fullscreen",
            images_upload_handler: images_upload_handler,
            setup: (editor) => {
                editorRef.value = editor;
                /* Helper functions */
                // var toDateHtml = function(date) {
                // 	return '<time datetime="' + date.toString() + '">' + date.toDateString() + "</time>";
                // };
                // editor.ui.registry.addButton("Pdf", {
                // 	text: "Pdf",
                // 	tooltip: "上传pdf",
                // 	onAction: function(_) {
                // 		editor.insertContent(toDateHtml(new Date()));
                // 	}
                // });
                editor.on("init", () => {
                    editor.setContent(props.value);
                });
                editor.on("change keyup undo redo", () => {
                    emit("update:value", editor.getContent());
                });
                editor.on("NodeChange", () => {
                    emit("update:value", editor.getContent());
                });
            },
            file_picker_types: "file media",
            file_picker_callback: file_picker_callback,
            init_instance_callback: (editor) => {
                console.log(`Editor: ${editor.id} is now initialized.`);
            },
        })
        .then((editor) => {
            console.log("then", editor);
        })
        .catch((err) => {
            console.log("err", err);
        });
};
onMounted(() => {
    console.log("onMounted");
    init();
});
onUnmounted(() => {
    tinymce.remove(); //销毁
});
</script>
<style scoped>
.tinymce-box {
    width: 100%;
    display: block;
}
</style>
