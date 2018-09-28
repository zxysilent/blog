<style lang="less">
@import "./add.less";
</style>

<template>
    <div>
     <Card>
            <Row>
            <Col span="18">
                <Form :label-width="80">
                    <FormItem label="文章标题" :error="articleError">
                        <Input v-model="articleTitle" @on-blur="handleArticletitleBlur" icon="android-list" />
                    </FormItem>
                    <FormItem>
                        <Input v-model="value11">
                        <span slot="prepend">https://blog.zxysilent.com/post/</span>
                        <span slot="append">.html
                            <Icon type="ios-eye-outline" />
                        </span>
                        </Input>
                    </FormItem>
                </Form>
                <div class="margin-top-20">
                    <mavon-editor style="height: 100%"></mavon-editor>
                </div>
            </Col>
            <Col span="6" class="padding-left-10">
            <Card>
                <p slot="title">
                    <Icon type="paper-airplane"></Icon>
                    发布
                </p>
                <p class="margin-top-10">
                    <Icon type="android-time"></Icon>权限
                    <i-switch  size="large" />
                </p>
                <p class="margin-top-10">
                    <Icon type="ios-eye"></Icon>&nbsp;&nbsp;公开度：&nbsp;
                    <i-switch  />
                </p>
                <p class="margin-top-10">
                    <Icon type="ios-calendar-outline"></Icon>&nbsp;&nbsp;
                    <DatePicker @on-change="setPublishTime" type="datetime" style="width:200px;" placeholder="选择日期和时间"></DatePicker>
                </p>
                <Row class="margin-top-20 publish-button-con">
                    <span class="publish-button">
                        <Button @click="handlePreview">预览</Button>
                    </span>
                    <span class="publish-button">
                        <Button @click="handleSaveDraft">保存草稿</Button>
                    </span>
                    <span class="publish-button">
                        <Button @click="handlePublish" :loading="publishLoading" icon="ios-checkmark" style="width:90px;" type="primary">发布</Button>
                    </span>
                </Row>
            </Card>
            <div class="margin-top-10">
                <Card>
                    <p slot="title">
                        <Icon type="ios-pricetags-outline"></Icon>
                        标签
                    </p>
                    <Row>
                        <Col span="24">
                        <Select v-model="articleTagSelected" multiple @on-change="handleSelectTag" placeholder="请选择文章标签">
                            <Option v-for="item in articleTagList" :value="item.value" :key="item.value">{{ item.value }}</Option>
                        </Select>
                        </Col>
                    </Row>
                    <transition name="add-new-tag">
                        <div v-show="addingNewTag" class="add-new-tag-con">
                            <Col span="14">
                            <Input v-model="newTagName" placeholder="请输入标签名" />
                            </Col>
                            <Col span="5" class="padding-left-10">
                            <Button @click="createNewTag" long type="primary">确定</Button>
                            </Col>
                            <Col span="5" class="padding-left-10">
                            <Button @click="cancelCreateNewTag" long type="ghost">取消</Button>
                            </Col>
                        </div>
                    </transition>
                </Card>
            </div>
            </Col>
        </Row>
     </Card>
    </div>
</template>

<script>
import { mavonEditor } from "mavon-editor";
import "mavon-editor/dist/css/index.css";
export default {
	name: "artical-publish",
	components: {
		mavonEditor
		// or 'mavon-editor': mavonEditor
	},
	data() {
		return {
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
						window.location.host +
						"/" +
						year +
						"/" +
						month +
						"/" +
						day +
						"/";
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
	destroyed() {}
};
</script>
