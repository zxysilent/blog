<template>
    <n-upload
        ref="uploadRef"
        :headers="{ Authorization: storage.getToken() }"
        list-type="image-card"
        :action="actionURL"
        :max="1"
        v-model:file-list="files"
        @before-upload="beforeUpload"
        @remove="onRemove"
        @finish="onFinish"
        :file-list-style="props.fileListStyle"
        response-type="json"
        accept="image/png,image/jpeg,image/bmp">
        <n-icon size="25">
            <CloudUploadOutline />
        </n-icon>
    </n-upload>
</template>
<script lang="ts" setup>
import { useMessage } from "naive-ui";
import storage from "@/utils/storage";
import { ref, watchEffect } from "vue";
import type { UploadFileInfo } from "naive-ui";
import { CloudUploadOutline } from "@vicons/ionicons5";
const emit = defineEmits(["update:value"]);
const mb = 1024 * 1024; //1mb
const props = defineProps({
    value: { type: String, required: true },
    size: { type: Number, default: 8 },
    onChange: { type: Function, default: undefined },
    fileListStyle: { type: Object, default: undefined },
});
const files = ref<UploadFileInfo[]>([]);
const uploadRef = ref();
const url = ref<string>(props.value);
const message = useMessage();
const actionURL = "/api/upload/image";
watchEffect(() => {
    console.log("watchEffect", props.value);
    if (props.value != "") {
        files.value = [{ id: new Date().getTime() + "", name: "image", status: "finished", url: props.value, type: "image/jpeg" }];
    } else {
        files.value = [];
    }
});
const beforeUpload = async (data: { file: UploadFileInfo; fileList: UploadFileInfo[] }) => {
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
const onFinish = ({ file, event }: { file: UploadFileInfo; event?: ProgressEvent }) => {
    console.log(file, event);
    const resp = (event?.target as XMLHttpRequest).response;
    console.log(resp);
    file.url = resp.data;
    if (resp.code == 200) {
        url.value = resp.data;
        emit("update:value", url.value);
        if (props.onChange) {
            props.onChange(url.value);
        }
    } else {
        file.status = "error";
        emit("update:value", "");
        message.error(resp.msg);
        if (props.onChange) {
            props.onChange("");
        }
    }
    file.url = resp.data;
    return file;
};
// const reset = () => {
//     console.log("reset");
//     uploadRef.value.clear();
// };
const onRemove = (options: { file: UploadFileInfo; fileList: Array<UploadFileInfo> }) => {
    console.log(options);
    emit("update:value", "");
    if (props.onChange) {
        props.onChange("");
    }
    return true;
};
// defineExpose({ reset });
</script>
