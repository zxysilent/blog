<template>
    <n-upload
        ref="uploadRef"
        :headers="{ Authorization: storage.getToken() }"
        list-type="image-card"
        :action="actionURL"
        :max="props.max"
        v-model:file-list="files"
        @before-upload="beforeUpload"
        @remove="onRemove"
        @change="onChange"
        @finish="onFinish"
        response-type="json"
        accept="image/png,image/jpeg,image/bmp">
        <n-icon size="25">
            <CloudUploadOutline />
        </n-icon>
    </n-upload>
</template>
<script lang="ts" setup>
import { ref } from "vue";
import { useMessage } from "naive-ui";
import storage from "@/utils/storage";
import type { UploadFileInfo } from "naive-ui";
import { CloudUploadOutline } from "@vicons/ionicons5";
const emit = defineEmits(["update:value"]);
const mb = 1024 * 1024; //1mb
const props = defineProps({
    value: { type: Array<string>, required: true },
    size: { type: Number, default: 8 },
    onChange: { type: Function },
    max: { type: Number, default: 3 },
});
const uploadRef = ref();
const message = useMessage();
const actionURL = "/api/upload/image";
const files = ref<UploadFileInfo[]>([]);
files.value = props.value.map((item) => {
    return { id: new Date().getTime() + "", name: "image", status: "finished", url: item };
});
console.log(files.value);
const beforeUpload = async (data: { file: UploadFileInfo; fileList: UploadFileInfo[] }) => {
    console.log(data.file.file?.type);
    console.log(data.file.file?.size);
    if (data.file.file?.type.indexOf("image") == -1) {
        message.error("只支持上传图片");
        return false;
    }
    if ((data.file.file?.size as Number) > mb * props.size) {
        message.error(`文件大小不能超过${props.size}MB`);
        return false;
    }
    return true;
};
const onChange = (data: { fileList: UploadFileInfo[] }) => {
    let arr = data.fileList.map((item) => {
        return item.url;
    });
    console.log(data.fileList);
    console.log("arr", arr);
    emit("update:value", arr);
    if (props.onChange) {
        props.onChange(arr);
    }
};
const onFinish = ({ file, event }: { file: UploadFileInfo; event?: ProgressEvent }) => {
    console.log(file, event);
    const resp = (event?.target as XMLHttpRequest).response;
    console.log(resp);
    file.url = resp.data;
    if (resp.code == 200) {
        file.url = resp.data;
    } else {
        file.status = "error";
        message.error(resp.msg);
    }
    return file;
};
const onRemove = (options: { file: UploadFileInfo; fileList: Array<UploadFileInfo> }) => {
    console.log(options);
    return true;
};
</script>
