<style lang="less">
.v-note-wrapper {
    z-index: 999 !important;
}
.v-note-wrapper .v-note-panel .v-note-edit.divarea-wrapper .content-input-wrapper {
    padding: 4px 10px 10px 10px !important;
}
.v-note-wrapper .v-note-op .v-right-item .op-icon {
    margin: 6px 1px 5px 0 !important;
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
    <div style="min-height: 600px; height: calc(100vh - 55px)">
        <mavon-editor
            ref="mdRef"
            @imgAdd="onImageAdd"
            @change="onChange"
            @save="onSave"
            :subfield="props.subfield"
            :editable="!props.readonly"
            :navigation="props.navigation"
            :defaultOpen="props.defaultOpen"
            :boxShadow="false"
            :toolbars="toolbars"
            v-model="markdown"
            style="height: 100%"
            :ishljs="true"
            :externalLink="externalLink">
        </mavon-editor>
    </div>
</template>
<script lang="ts" setup>
import { useMessage } from "naive-ui";
//@ts-ignore
import { mavonEditor } from "mavon-editor";
import "mavon-editor/dist/css/index.css";
import storage from "@/utils/storage";
import { ref, onMounted, watchEffect } from "vue";
const actionURL = "/api/upload/image";
const message = useMessage();
const emit = defineEmits(["update:value", "change", "save"]);
const props = defineProps({
    value: { type: String, required: true },
    readonly: { type: Boolean, default: false },
    subfield: { type: Boolean, default: true },
    navigation: { type: Boolean, default: false },
    defaultOpen: { type: String, default: "" },
});
const mdRef = ref();
const toolbars = {
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
    readmodel: false, // 沉浸式阅读
    htmlcode: false, // 展示html源码
    help: false, // 帮助
    /* 1.3.5 */
    undo: true, // 上一步
    redo: true, // 下一步
    trash: true, // 清空
    save: true, // 保存（触发events中的save事件）
    /* 1.4.2 */
    navigation: false, // 导航目录
    /* 2.1.8 */
    alignleft: true, // 左对齐
    aligncenter: true, // 居中
    alignright: true, // 右对齐
    /* 2.2.1 */
    subfield: true, // 单双栏模式
    preview: true, // 预览
};
const externalLink = {
    markdown_css: function () {
        // 这是你的markdown css文件路径
        return "/static/markdown/github-markdown.min.css";
    },
    hljs_js: function () {
        // 这是你的hljs文件路径
        return "/static/highlightjs/highlight.min.js";
    },
    hljs_css: function (css) {
        // 这是你的代码高亮配色文件路径
        return "/static/highlightjs/styles/" + css + ".min.css";
    },
    hljs_lang: function (lang) {
        // 这是你的代码高亮语言解析路径
        return "/static/highlightjs/languages/" + lang + ".min.js";
    },
    katex_css: function () {
        // 这是你的katex配色方案路径路径
        return "/static/katex/katex.min.css";
    },
    katex_js: function () {
        // 这是你的katex.js路径
        return "/static/katex/katex.min.js";
    },
};
const markdown = ref<string>(props.value);
watchEffect(() => {
    console.log("watchEffect");
    markdown.value = props.value;
});
const emptyHistory = () => {
    mdRef.value.$emptyHistory();
};
const onSave = (_, html) => {
    emit("save", html);
};
const onChange = (value, html) => {
    // console.log(value, html);
    emit("update:value", value);
    emit("change", html);
};

const onImageAdd = (pos, $file) => {
    const formData = new FormData();
    const xhr = new XMLHttpRequest();
    xhr.withCredentials = false;
    xhr.open("POST", actionURL);
    xhr.onload = () => {
        if (xhr.status < 200 || xhr.status >= 300) {
            message.warning("图片上传失败,HTTP Error: " + xhr.status);
            return;
        }
        const json = JSON.parse(xhr.responseText);
        if (!json || json.code != 200) {
            message.warning("图片上传失败");
            return;
        }
        mdRef.value.$img2Url(pos, json.data);
    };
    formData.append("token", storage.getToken());
    formData.append("file", $file);
    xhr.send(formData);
};
defineExpose({ emptyHistory });
onMounted(() => {});
</script>
