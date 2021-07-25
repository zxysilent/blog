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
								<Icon type="ios-eye" size="22" />
							</a></Button>
							<Button type="info" @click="cmtDraft" :loading="draftLoading">
								<Icon type="ios-trash" size="20" />存草稿
							</Button>&nbsp;
							<Button type="warning" @click="cmtPublish" :loading="publishLoading">
								<Icon type="ios-send" size="20" />发 布
							</Button>
							</Col>
						</Row>
					</FormItem>
				</Form>
				<!-- markdown -->
				<mavon v-model="dataForm.markdown" @on-change="onChange"></mavon>
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
						<FormItem label="评论" prop="allow">
							<Checkbox v-model="dataForm.allow">允许评论</Checkbox>
						</FormItem>
						<FormItem label="发布日期">
							<DatePicker v-model="dataForm.created" type="datetime" placeholder="选择发布日期和时间" :clearable="false" :editable="false"></DatePicker>
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
import Mavon from "@/components/markdown/Mavon.vue";
import { apiCateAll } from "@/api/cate";
import { apiTagAll } from "@/api/tag";
import Storage from "@/utils/storage";
import { apiPostGet, admPostOpts, apiPostTagGet } from "@/api/post";
// 通用 文章/页面 + 添加/修改
// 减少js体积
export default {
	components: {
		Mavon
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
				allow: false,
				created: "",
				richtext: "",
				markdown: ""
			},
			tags: [],
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
				apiPostTagGet(this.dataId).then((resp) => {
					if (resp.code == 200) {
						this.tags = resp.data;
					}
				});
			}
		},
		getCate() {
			apiCateAll().then((resp) => {
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
			apiTagAll().then((resp) => {
				if (resp.code == 200) {
					this.tagAll = resp.data;
				} else {
					this.tagAll = [];
					this.$Message.warning("未查询到标签信息,请重试！");
				}
			});
		},
		getOne() {
			apiPostGet(this.dataId).then((resp) => {
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
						allow: false,
						created: "",
						richtext: "",
						markdown_richtext: ""
					};
					this.$Message.warning("未查询到信息，请重试！");
				}
			});
		},
		onChange(html) {
            console.log(html)
			// 保存 html
			this.dataForm.richtext = html;
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

			formData.append("token", Storage.getToken());
			formData.append("file", $file);
			xhr.send(formData);
		},
		// 存草稿
		cmtDraft() {
			this.$refs.dataForm.validate((valid) => {
				if (valid) {
					if (this.dataForm.richtext == "") {
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
					}).then((resp) => {
						this.draftLoading = false;
						if (resp.code == 200) {
							this.$Message.success({
								richtext: resp.msg
							});
						} else if (resp.code == 300) {
							this.$Message.warning({
								richtext: resp.msg
							});
						} else {
							this.$Message.error({
								richtext: resp.msg
							});
						}
					});
				}
			});
		},
		// 发布
		cmtPublish() {
			this.$refs.dataForm.validate((valid) => {
				if (valid) {
					if (this.dataForm.richtext == "") {
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
					}).then((resp) => {
						this.publishLoading = false;
						if (resp.code == 200) {
							this.$Message.success({
								richtext: resp.msg,
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
								richtext: resp.msg
							});
						} else {
							this.$Message.error({
								richtext: resp.msg
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
		// var md = this.$refs.md;
		// var toolbar_left = md.$refs.toolbar_left;
		// var diy = this.$refs.diy;
		// toolbar_left.$el.append(diy);
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
			this.dataForm.created = new Date();
		}
		this.init();
	}
};
</script>
