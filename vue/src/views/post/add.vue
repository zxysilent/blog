<style lang="less">
@import "./add.less";
</style>

<template>
    <div>
        <Card dis-hover>
            <Row :gutter="6">
                <Col span="18">
                <Form>
                    <FormItem style="margin-bottom:15px">
                        <Input v-model="articleTitle" placeholder="请输入标题" ></Input>
                    </FormItem>
                    <FormItem style="margin-bottom:15px" label="https://blog.zxysilent.com/post/">
                        <Row>
                            <Col span="10">
                            <Input type="text"  placeholder="请输入访问路径"></Input>
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
                <div>
                    <mavon-editor :boxShadow="false" :toolbars="toolbars" @fullScreen="fullScreen" v-model="dataForm.markdown_content" style="min-height: 600px;height: calc(100vh - 200px);"></mavon-editor>
                </div>
                </Col>
                <Col span="6">
                <Card dis-hover :bordered="false">
                    <p slot="title">
                        <Icon type="ios-settings-outline" />
                        设置
                    </p>
                    <Form ref="dataForm" :model="dataForm" :rules="dataRules" :label-width="80">
                        <FormItem label="分类">
                            <RadioGroup v-model="dataForm.cateId">
                                <Radio label="cate" v-for="item in cateAll" :label="item.id" :key="item.id">{{ item.name }} [{{ item.intro }}]</Radio>
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
      },
      articleTitle: "",
      articleError: "",
      showLink: false,
      fixedLink: "",
      articlePath: "",
      articlePathHasEdited: false,
      editLink: false,
      editPathButtonType: "ghost",
      editPathButtonText: "编辑",
      articleStateList: [{ value: "草稿" }, { value: "等待复审" }],
      editOpenness: false,
      Openness: "公开",
      currentOpenness: "公开",
      topArticle: false,
      publishTime: "",
      publishTimeType: "immediately",
      editPublishTime: false, // 是否正在编辑发布时间
      articleTagSelected: [], // 文章选中的标签
      articleTagList: [], // 所有标签列表
      classificationList: [],
      classificationSelected: [], // 在所有分类目录中选中的目录数组
      offenUsedClass: [],
      offenUsedClassSelected: [], // 常用目录选中的目录
      classificationFinalSelected: [], // 最后实际选择的目录
      publishLoading: false,
      addingNewTag: false, // 添加新标签
      newTagName: "" // 新建标签名
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
    fullScreen(status,value){
        console.log(status,value)
    },
    handleSubmit(name) {
      this.$refs[name].validate(valid => {
        if (valid) {
          this.$Message.success("Success!");
        } else {
          this.$Message.error("Fail!");
        }
      });
    },
    handleReset(name) {
      this.$refs[name].resetFields();
    },
    handleArticletitleBlur() {
      if (this.articleTitle.length !== 0) {
        // this.articleError = '';
        localStorage.articleTitle = this.articleTitle; // 本地存储文章标题
        if (!this.articlePathHasEdited) {
          let date = new Date();
          let year = date.getFullYear();
          let month = date.getMonth() + 1;
          let day = date.getDate();
          this.fixedLink =
            window.location.host + "/" + year + "/" + month + "/" + day + "/";
          this.articlePath = this.articleTitle;
          this.articlePathHasEdited = true;
          this.showLink = true;
        }
      } else {
        // this.articleError = '文章标题不可为空哦';
        this.$Message.error("文章标题不可为空哦");
      }
    },
    editArticlePath() {
      this.editLink = !this.editLink;
      this.editPathButtonType =
        this.editPathButtonType === "ghost" ? "success" : "ghost";
      this.editPathButtonText =
        this.editPathButtonText === "编辑" ? "完成" : "编辑";
    },
    handleEditOpenness() {
      this.editOpenness = !this.editOpenness;
    },
    handleSaveOpenness() {
      this.Openness = this.currentOpenness;
      this.editOpenness = false;
    },
    cancelEditOpenness() {
      this.currentOpenness = this.Openness;
      this.editOpenness = false;
    },
    handleEditPublishTime() {
      this.editPublishTime = !this.editPublishTime;
    },
    handleSavePublishTime() {
      this.publishTimeType = "timing";
      this.editPublishTime = false;
    },
    cancelEditPublishTime() {
      this.publishTimeType = "immediately";
      this.editPublishTime = false;
    },
    setPublishTime(datetime) {
      this.publishTime = datetime;
    },
    setClassificationInAll(selectedArray) {
      this.classificationFinalSelected = selectedArray.map(item => {
        return item.title;
      });
      localStorage.classificationSelected = JSON.stringify(
        this.classificationFinalSelected
      ); // 本地存储所选目录列表
    },
    setClassificationInOffen(selectedArray) {
      this.classificationFinalSelected = selectedArray;
    },
    handleAddNewTag() {
      this.addingNewTag = !this.addingNewTag;
    },
    createNewTag() {
      if (this.newTagName.length !== 0) {
        this.articleTagList.push({ value: this.newTagName });
        this.addingNewTag = false;
        setTimeout(() => {
          this.newTagName = "";
        }, 200);
      } else {
        this.$Message.error("请输入标签名");
      }
    },
    cancelCreateNewTag() {
      this.newTagName = "";
      this.addingNewTag = false;
    },
    canPublish() {
      if (this.articleTitle.length === 0) {
        this.$Message.error("请输入文章标题");
        return false;
      } else {
        return true;
      }
    },
    handlePreview() {
      if (this.canPublish()) {
        if (this.publishTimeType === "immediately") {
          let date = new Date();
          let year = date.getFullYear();
          let month = date.getMonth() + 1;
          let day = date.getDate();
          let hour = date.getHours();
          let minute = date.getMinutes();
          let second = date.getSeconds();
          localStorage.publishTime =
            year +
            " 年 " +
            month +
            " 月 " +
            day +
            " 日 -- " +
            hour +
            " : " +
            minute +
            " : " +
            second;
        } else {
          localStorage.publishTime = this.publishTime; // 本地存储发布时间
        }
        localStorage.content = tinymce.activeEditor.getContent();
        this.$router.push({
          name: "preview"
        });
      }
    },
    handleSaveDraft() {
      if (!this.canPublish()) {
        //
      }
    },
    handlePublish() {
      if (this.canPublish()) {
        this.publishLoading = true;
        setTimeout(() => {
          this.publishLoading = false;
          this.$Notice.success({
            title: "保存成功",
            desc: "文章《" + this.articleTitle + "》保存成功"
          });
        }, 1000);
      }
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
  mounted() {
    this.articleTagList = [
      { value: "vue" },
      { value: "iview" },
      { value: "ES6" },
      { value: "webpack" },
      { value: "babel" },
      { value: "eslint" }
    ];
    this.classificationList = [
      {
        title: "Vue实例",
        expand: true,
        children: [
          {
            title: "数据与方法",
            expand: true
          },
          {
            title: "生命周期",
            expand: true
          }
        ]
      },
      {
        title: "Class与Style绑定",
        expand: true,
        children: [
          {
            title: "绑定HTML class",
            expand: true,
            children: [
              {
                title: "对象语法",
                expand: true
              },
              {
                title: "数组语法",
                expand: true
              },
              {
                title: "用在组件上",
                expand: true
              }
            ]
          },
          {
            title: "生命周期",
            expand: true
          }
        ]
      },
      {
        title: "模板语法",
        expand: true,
        children: [
          {
            title: "插值",
            expand: true
          },
          {
            title: "指令",
            expand: true
          },
          {
            title: "缩写",
            expand: true
          }
        ]
      }
    ];
    this.offenUsedClass = [
      {
        title: "vue实例"
      },
      {
        title: "生命周期"
      },
      {
        title: "模板语法"
      },
      {
        title: "插值"
      },
      {
        title: "缩写"
      }
    ];
  },
  created() {
    this.dataForm.create_time = new Date();
    this.init();
  },
  destroyed() {}
};
</script>
