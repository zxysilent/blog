<style lang="less">
@import "./article.less";
</style>
<template>
	<div>
		<Card dis-hover>
			<Row :gutter="2">
				<Col span="20">
				<Form ref="dataForm" :rules="dataRules" :model="dataForm">
					<FormItem style="margin-bottom:15px" prop="title">
						<Input v-model="dataForm.title" placeholder="请输入标题"></Input>
					</FormItem>
					<FormItem style="margin-bottom:15px" prop="path">
						<Row>
							<Col span="3" style="min-width:195px">
							{{prefix}}
							</Col>
							<Col span="8"><Input type="text" v-bind:disabled="isEdit" v-model="dataForm.path" placeholder="请输入访问路径"></Input>
							</Col>
							<Col span="10" style="min-width:250px"> .html&nbsp;<a title="预览" @click="clkPreview">
								<Icon type="ios-eye" size="22" /></a></Button>
							<Button type="info" @click="cmtDraft" :loading="draftLoading">
								<Icon type="ios-trash" size="20" />存草稿</Button>&nbsp;
							<Button type="warning" @click="cmtPublish" :loading="publishLoading">
								<Icon type="ios-send" size="20" />发 布</Button>
							</Col>
						</Row>
					</FormItem>
				</Form>
				<div style="min-height: 600px;height: calc(100vh - 145px);">
					<button ref="diy" type="button" @click="diyMore" class="op-icon ivu-icon ivu-icon-md-code" aria-hidden="true" title="更多"></button>
					<mavon-editor ref="md" @imgAdd="imgAdd" :boxShadow="false" @change="change" :toolbars="toolbars" v-model="dataForm.markdown_content" style="height:100%">
					</mavon-editor>
				</div>
				</Col>
				<Col span="4">
				<Card dis-hover :bordered="false">
					<p slot="title">
						<Icon type="ios-settings-outline" />
						设置
					</p>
					<Form ref="optsForm" label-position="top">
						<FormItem v-if="isPost" label="分类">
							<RadioGroup v-model="dataForm.cate_id">
								<Radio label="cate" v-for="item in cateAll" :label="item.id" :key="item.id">{{item.name}}[{{item.intro}}]</Radio>
							</RadioGroup>
						</FormItem>
						<FormItem label="权限" prop="is_public">
							<i-switch v-model="dataForm.is_public"><span slot="open">公开</span>
								<span slot="close">私密</span></i-switch>
						</FormItem>
						<FormItem label="评论" prop="allow_comment">
							<Checkbox v-model="dataForm.allow_comment">允许评论</Checkbox>
						</FormItem>
						<FormItem label="发布日期">
							<DatePicker v-model="dataForm.create_time" type="datetime" placeholder="选择发布日期和时间" :clearable="false" :editable="false"></DatePicker>
						</FormItem>
						<FormItem v-if="isPost" label="标签">
							<Select v-model="tags" multiple placeholder="请选择文章标签">
								<Option v-for="item in tagAll" :value="item.id" :key="item.id">{{ item.name }}</Option>
							</Select>
						</FormItem>
						<br><br>
					</Form>
				</Card>
				</Col>
			</Row>
		</Card>
	</div>
</template>
<script>
import { mavonEditor } from "mavon-editor";
import "mavon-editor/dist/css/index.css";
import toolbars from "./toolbars";
import { apiCateAll } from "@/api/cate";
import { apiTagAll } from "@/api/tag";
import util from "@/utils.js";
import { apiPostGet, admPostOpts, apiPostTagGet } from "@/api/post";
// 通用 文章/页面 + 添加/修改
// 减少js体积
export default {
	components: {
		mavonEditor
	},
	data() {
		return {
			cateAll: [],
			tagAll: [],
			draftLoading: false,
			publishLoading: false,
			dataForm: {
				title: "",
				path: "",
				status: 0,
				summary: "",
				cate_id: 0,
				is_public: true,
				allow_comment: false,
				create_time: "",
				content: "",
				markdown_content: ""
			},
			tags: [],
			toolbars: toolbars,
			dataRules: {
				title: [{ required: true, message: " ", trigger: "blur" }],
				path: [{ required: true, message: " ", trigger: "blur" }]
			}
		};
	},
	methods: {
		init() {
			if (this.isPost) {
				this.getCate();
				this.getTag();
			}
			if (this.isEdit) {
				this.getOne();
			}
			if (this.isPost && this.isEdit) {
				apiPostTagGet(this.dataId).then(resp => {
					if (resp.code == 200) {
						this.tags = resp.data;
					}
				});
			}
		},
		getCate() {
			apiCateAll().then(resp => {
				if (resp.code == 200) {
					this.cateAll = resp.data;
					if (this.isAdd) {
						this.dataForm.cate_id = resp.data[0].id;
					}
				} else {
					this.cateAll = [];
					this.$Message.warning("未查询到分类信息,请重试！");
				}
			});
		},
		getTag() {
			apiTagAll().then(resp => {
				if (resp.code == 200) {
					this.tagAll = resp.data;
				} else {
					this.tagAll = [];
					this.$Message.warning("未查询到标签信息,请重试！");
				}
			});
		},
		getOne() {
			apiPostGet(this.dataId).then(resp => {
				if (resp.code == 200) {
					this.dataForm = resp.data;
				} else {
					this.dataForm = {
						title: "",
						path: "",
						summary: "",
						cate_id: 0,
						status: 0,
						is_public: true,
						allow_comment: false,
						create_time: "",
						content: "",
						markdown_content: ""
					};
					this.$Message.warning("未查询到信息，请重试！");
				}
			});
		},
		change(value, html) {
			// 保存 html
			this.dataForm.content = html;
		},
		clkPreview() {
			if (!this.dataForm.id) {
				this.$Message.warning("添加状态无法预览");
				return;
			}
			let ext = "";
			if (this.isPost) {
				ext = ".html";
			}
			window.open(this.prefix + this.dataForm.path + ext);
		},
		diyMore(val) {
			this.$refs.md.insertText(this.$refs.md.getTextareaDom(), {
				prefix: "",
				subfix: "",
				str: "<!--more-->"
			});
		},
		imgAdd(pos, $file) {
			let formData = new FormData();
			let xhr = new XMLHttpRequest();
			xhr.withCredentials = false;
			xhr.open("POST", process.env.VUE_APP_SRV + "/adm/upload");
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

			formData.append("token", util.getToken());
			formData.append("file", $file);
			xhr.send(formData);
		},
		// 存草稿
		cmtDraft() {
			this.$refs.dataForm.validate(valid => {
				if (valid) {
					if (this.dataForm.content == "") {
						this.$Message.warning("请填写内容");
						return;
					}
					console.log(this.dataForm);
					this.dataForm.status = 0; //草稿
					this.draftLoading = true;
					admPostOpts({
						post: this.dataForm,
						edit: this.isEdit,
						type: this.isPost ? 0 : 1,
						tags: this.tags
					}).then(resp => {
						this.draftLoading = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: resp.msg
							});
						} else if (resp.code == 300) {
							this.$Message.warning({
								content: resp.msg
							});
						} else {
							this.$Message.error({
								content: resp.msg
							});
						}
					});
				}
			});
		},
		// 发布
		cmtPublish() {
			this.$refs.dataForm.validate(valid => {
				if (valid) {
					if (this.dataForm.content == "") {
						this.$Message.warning("请填写内容");
						return;
					}
					this.dataForm.status = 3; //发布
					this.publishLoading = true;
					admPostOpts({
						post: this.dataForm,
						edit: this.isEdit,
						type: this.isPost ? 0 : 1,
						tags: this.tags
					}).then(resp => {
						this.publishLoading = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: resp.msg,
								onClose: () => {
									if (!this.isEdit) {
										if (this.isPost) {
											this.$router.push({
												name: "post-list"
											});
										} else {
											this.$router.push({
												name: "page-list"
											});
										}
									}
								}
							});
						} else if (resp.code == 300) {
							this.$Message.warning({
								content: resp.msg
							});
						} else {
							this.$Message.error({
								content: resp.msg
							});
						}
					});
				}
			});
		}
	},
	computed: {
		// 是否 post
		isPost() {
			return this.$route.name.indexOf("post") > -1;
		},
		// 是否编辑
		isEdit() {
			return this.$route.name.indexOf("edit") > -1;
		},
		// 是否添加
		isAdd() {
			return this.$route.name.indexOf("add") > -1;
		},
		// 数据id
		dataId() {
			if (this.isEdit) {
				return parseInt(this.$route.params.id, 10);
			}
			return 0;
		},
		// 访问前缀
		prefix() {
			if (this.isPost) {
				return process.env.VUE_APP_SRV + "/post/";
			}
			return process.env.VUE_APP_SRV + "/page/";
		}
	},
	mounted() {
		var md = this.$refs.md;
		var toolbar_left = md.$refs.toolbar_left;
		var diy = this.$refs.diy;
		toolbar_left.$el.append(diy);
		//   beforeRouteUpdate(to, from, next) {
		//     //参数不同才触发
		//     console.log("beforeRouteUpdate");
		//     // don't forget to call
		//     next();
		//   },
		//   watch: {
		//     $route(to, from) {
		//       //刷新参数放到这里里面去触发就可以刷新相同界面了
		//       console.log("watch");
		//     }
		//   },
		if (this.isAdd) {
			this.dataForm.create_time = new Date();
		}
		this.init();
	}
};
</script>
