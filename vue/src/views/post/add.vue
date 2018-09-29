<style lang="less">
@import "./add.less";
</style>

<template>
    <div>
        <Card dis-hover>
            <Row :gutter="2">
                <Col span="20">
                <Form>
                    <FormItem style="margin-bottom:15px" prop="titel">
                        <Input v-model="dataForm.title" placeholder="请输入标题" ></Input>
                    </FormItem>
                    <FormItem style="margin-bottom:15px" label="https://blog.zxysilent.com/post/" prop="path">
                        <Row>
                            <Col span="10">
                            <Input type="text"  v-model="dataForm.path" placeholder="请输入访问路径"></Input>
                            </Col>
                            <Col span="8">.html
                            <Button type="dashed" @click="clkPreview">
                                <Icon type="ios-eye" size="20" /> 预 览</Button>
                            <Button type="info">
                                <Icon type="ios-trash" size="20" /> 存草稿</Button>
                            <Button type="warning">
                                <Icon type="ios-send" size="20" /> 发 布</Button>
                            </Col>
                        </Row>
                    </FormItem>
                </Form>
                <div style="min-height: 600px;height: calc(100vh - 200px);">
                    <mavon-editor :boxShadow="false" :toolbars="toolbars" @fullScreen="fullScreen" v-model="dataForm.markdown_content" style="height:100%"></mavon-editor>
                </div>
                </Col>
                <Col span="4">
                <Card dis-hover :bordered="false">
                    <p slot="title">
                        <Icon type="ios-settings-outline" />
                        设置
                    </p>
                    <Form ref="dataForm" :model="dataForm" :rules="dataRules"  label-position="top">
                        <FormItem label="分类">
                            <RadioGroup v-model="dataForm.cateId">
                                <Radio label="cate" v-for="item in cateAll" :label="item.id" :key="item.id">{{item.name}}[{{item.intro}}]</Radio>
                            </RadioGroup>
                        </FormItem>
                        <FormItem label="权限" prop="is_public">
                            <i-switch :value="dataForm.is_public"><span slot="open">公开</span>
                                <span slot="close">私密</span></i-switch>
                        </FormItem>
                        <FormItem label="评论" prop="all_comment">
                            <Checkbox v-model="dataForm.all_comment">允许评论</Checkbox>
                        </FormItem>
                        <FormItem label="发布日期">
                            <DatePicker v-model="dataForm.create_time" type="datetime" placeholder="选择发布日期和时间" :clearable="false" :editable="false"></DatePicker>
                        </FormItem>
                        <FormItem label="标签">
                            <Select v-model="dataForm.tags" multiple @on-change="handleSelectTag" placeholder="请选择文章标签">
                                <Option v-for="item in tagAll" :value="item.id" :key="item.id">{{ item.name }}</Option>
                            </Select>
                        </FormItem>
                        <br><br>
                    </Form>
                </Card>
                </Col>
            </Row>
        </Card>
        <Modal v-model="showPreview" fullscreen title="预览文章">
            <div v-html="dataForm.content"></div>
        </Modal>
    </div>
</template>

<script>
import { mavonEditor } from "mavon-editor";
import "mavon-editor/dist/css/index.css";
import { cateAll } from "@/api/cate";
import { tagAll } from "@/api/tag";

export default {
  components: {
    mavonEditor
  },
  data() {
    return {
      cateAll: [],
      tagAll: [],
      showPreview: false,
      showFull: false,
      dataForm: {
        title: "",
        path: "",
        summary: "",
        cateId: 0,
        is_public: true,
        all_comment: false,
        create_time: "",
        tags: [],
        content: "",
        markdown_content: ""
      },
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
      },
      dataRules: {
        title: [
          {
            required: true,
            message: "The name cannot be empty",
            trigger: "blur"
          }
        ],
        path: [
          {
            required: true,
            message: "The name cannot be empty",
            trigger: "blur"
          }
        ]
      }
    };
  },
  methods: {
    init() {
      cateAll().then(resp => {
        if (resp.code == 200) {
          this.cateAll = resp.data;
          this.dataForm.cateId = resp.data[0].id;
        } else {
          this.cateAll = [];
          this.$Message.warning("未查询到分类信息,请重试！");
        }
      });
      tagAll().then(resp => {
        if (resp.code == 200) {
          this.tagAll = resp.data;
        } else {
          this.tagAll = [];
          this.$Message.warning("未查询到标签信息,请重试！");
        }
      });
    },
    clkPreview() {
      this.showPreview = true;
    },
    fullScreen(status, value) {
      this.showFull = status;
      console.log(status, value);
    },
    handleSelectTag() {
      localStorage.tagsList = JSON.stringify(this.articleTagSelected); // 本地存储文章标签列表
    }
  },
  computed: {
    completeUrl() {
      let finalUrl = this.fixedLink + this.articlePath;
      localStorage.finalUrl = finalUrl; // 本地存储完整文章路径
      return finalUrl;
    }
  },
  mounted() {},
  created() {
    this.dataForm.create_time = new Date();
    this.init();
  },
  destroyed() {}
};
</script>
