<style lang="less">
@import "./image-editor.less";
</style>

<template>
    <div class="image-editor">
        <Row :gutter="10">
            <Col span="12">
            <Card>
                <p slot="title">
                    <Icon type="qr-scanner"></Icon>
                    基础实例
                </p>
                <Row :gutter="10">
                    <Col span="14" class="image-editor-con1">
                    <div class="cropper">
                        <img id="cropimg1" alt="">
                            </div>
                        </Col>
                        <Col span="10" class="image-editor-con1">
                        <Row type="flex" justify="center" align="middle" class="image-editor-con1-preview-con">
                            <div id="preview1"></div>
                        </Row>
                        <div class="image-editor-con1-btn-con margin-top-10">
                            <input type="file" accept="image/png, image/jpeg, image/gif, image/jpg" @change="handleChange1" id="fileinput1" class="fileinput" />
                            <label class="filelabel" for="fileinput1">
                                <Icon type="image"></Icon>&nbsp;选择图片
                            </label>
                            <span><Button @click="handlecrop1" type="primary" icon="crop">裁剪</Button></span>
                        </div>
                        <Modal v-model="option1.showCropedImage">
                            <p slot="header">预览裁剪之后的图片</p>
                            <img :src="option1.cropedImg" alt="" v-if="option1.showCropedImage" style="width: 100%;">
                            </Modal>
                            </Col>
                </Row>
            </Card>
            </Col>
            <Col span="12">
            <Card>
                <p slot="title">
                    <Icon type="android-options"></Icon>
                    获取图片数据
                </p>
                <Row :gutter="10">
                    <Col span="14" class="image-editor-con2">
                    <div class="cropper">
                        <img id="cropimg2" src="../../../images/bg.jpg" alt="">
                            </div>
                        </Col>
                        <Col span="10" class="image-editor-con2">
                        <p><b>x:</b>{{ cropdata2.x }}</p>
                        <p><b>y:</b>{{ cropdata2.y }}</p>
                        <p><b>width:</b>{{ cropdata2.w }}</p>
                        <p><b>heigh:</b>{{ cropdata2.h }}</p>
                        <p><b>deg:</b>{{ cropdata2.deg }}</p>
                        <p><b>scaleX:</b>{{ cropdata2.scaleX }}</p>
                        <p><b>scaleY:</b>{{ cropdata2.scaleY }}</p>
                        <div class="margin-top-10" style="text-align: center;">
                            <ButtonGroup>
                                <Button @click="cropper2.rotate(-90)" type="primary">
                                    <Icon :size="14" type="arrow-return-left"></Icon>
                                </Button>
                                <Button @click="cropper2.rotate(90)" type="primary">
                                    <Icon :size="14" type="arrow-return-right"></Icon>
                                </Button>
                                <Button @click="cropper2.zoom(0.1)" type="primary">
                                    <Icon :size="14" type="plus-round"></Icon>
                                </Button>
                                <Button @click="cropper2.zoom(-0.1)" type="primary">
                                    <Icon :size="14" type="minus-round"></Icon>
                                </Button>
                                <Button @click="cropper2.scaleX(-cropper2.getData().scaleX)" type="primary">
                                    <Icon :size="14" type="android-more-horizontal"></Icon>
                                </Button>
                                <Button @click="cropper2.scaleY(-cropper2.getData().scaleY)" type="primary">
                                    <Icon :size="14" type="android-more-vertical"></Icon>
                                </Button>
                            </ButtonGroup>
                        </div>
                        </Col>
                </Row>
            </Card>
            </Col>
        </Row>
    </div>
</template>

<script>
import Cropper from "cropperjs";
import "./cropper.min.css";
export default {
  name: "image-editor",
  data() {
    return {
      cropper1: {},
      option1: {
        showCropedImage: false,
        cropedImg: ""
      },
      cropper2: {},
      cropdata2: {
        x: "",
        y: "",
        w: "",
        h: "",
        deg: "",
        scaleX: "",
        scaleY: ""
      }
    };
  },
  methods: {
    handleChange1(e) {
      let file = e.target.files[0];
      let reader = new FileReader();
      reader.onload = () => {
        this.cropper2.replace(reader.result);
        reader.onload = null;
      };
      reader.readAsDataURL(file);
    },
    handlecrop1() {
      let file = this.cropper.getCroppedCanvas().toDataURL();
      this.option1.cropedImg = file;
      this.option1.showCropedImage = true;
    }
  },
  mounted() {
    let img1 = document.getElementById("cropimg1");
    this.cropper1 = new Cropper(img1, {
      dragMode: "move",
      preview: "#preview1",
      restore: false,
      center: false,
      highlight: false,
      cropBoxMovable: false,
      toggleDragModeOnDblclick: false
    });

    let img2 = document.getElementById("cropimg2");
    this.cropper2 = new Cropper(img2, {
      dragMode: "move",
      preview: "#preview1",
      restore: false,
      center: false,
      highlight: false,
      cropBoxMovable: false,
      toggleDragModeOnDblclick: false
    });
    img2.addEventListener("crop", e => {
      this.cropdata2.x = parseInt(e.detail.x);
      this.cropdata2.y = parseInt(e.detail.y);
      this.cropdata2.w = parseInt(e.detail.width);
      this.cropdata2.h = parseInt(e.detail.height);
      this.cropdata2.deg = parseInt(e.detail.rotate);
      this.cropdata2.scaleX = parseInt(e.detail.scaleX);
      this.cropdata2.scaleY = parseInt(e.detail.scaleY);
    });
  }
};
</script>

<style>
</style>
