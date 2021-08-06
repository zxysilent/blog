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
							<Col span="8"><Input type="text" v-bind:disabled="true" v-model="dataForm.path" placeholder="请输入访问路径"></Input>
							</Col>
							<Col span="10" style="min-width:250px"> .html&nbsp;<a title="预览" @click="emitView">
								<Icon type="ios-eye" size="22" />
							</a>
							<Button type="info" @click="emitDraft" :loading="draftLoading">
								<Icon type="ios-trash" size="20" />存草稿
							</Button>&nbsp;
							<Button type="warning" @click="emitFinish" :loading="finishLoading">
								<Icon type="ios-send" size="20" />发 布
							</Button>
							<Button :to="{name:'post-list'}" style="margin-left: 10px">返回列表</Button>
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
						<FormItem label="分类">
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
						<FormItem label="标签">
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
import { apiPostGet, admPostEdit } from "@/api/post";
export default {
	components: {
		Mavon
	},
	data() {
		return {
			cateAll: [],
			tagAll: [],
			tags: [],
			draftLoading: false,
			finishLoading: false,
			// 访问前缀
			prefix: process.env.VUE_APP_SRV + "/post/",
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
			dataRules: {
				title: [{ required: true, message: " ", trigger: "blur" }],
				path: [{ required: true, message: " ", trigger: "blur" }]
			}
		};
	},
	methods: {
		preinit() {
			apiCateAll().then((resp) => {
				if (resp.code == 200) {
					this.cateAll = resp.data;
				} else {
					this.$Message.error({ content: resp.msg, duration: 3 });
				}
			});
			apiTagAll().then((resp) => {
				if (resp.code == 200) {
					this.tagAll = resp.data;
				} else {
					this.$Message.error({ content: resp.msg, duration: 3 });
				}
			});
		},
		init() {
			apiPostGet({ id: this.dataForm.id }).then((resp) => {
				if (resp.code == 200) {
                    // 临时保存
					this.tags = resp.data.tags.map((item) => {
						return item.id;
					});
					this.dataForm = resp.data;
				} else {
					this.$Message.error({ content: resp.msg, duration: 3 });
				}
			});
		},
		onChange(html) {
			console.log(html);
			// 保存 html
			this.dataForm.richtext = html;
		},
		emitView() {
			window.open(this.prefix + this.dataForm.path + ".html");
		},
		// 存草稿
		emitDraft() {
			this.$refs.dataForm.validate((valid) => {
				if (valid) {
					if (this.dataForm.richtext == "") {
						this.$Message.warning("请填写内容");
						return;
					}
					console.log(this.dataForm);
					this.dataForm.status = 1; //草稿
					this.draftLoading = true;
                    // 还原
					this.dataForm.tags = this.tags.map((item) => {
						return { id: item };
					});
					admPostEdit(this.dataForm).then((resp) => {
						this.draftLoading = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: "保存成功",
								onClose: () => {
									this.init();
								}
							});
						} else {
							this.$Message.error({ content: resp.msg, duration: 3 });
						}
					});
				}
			});
		},
		// 发布
		emitFinish() {
			this.$refs.dataForm.validate((valid) => {
				if (valid) {
					if (this.dataForm.richtext == "") {
						this.$Message.warning("请填写内容");
						return;
					}
					console.log(this.dataForm);
					this.dataForm.status = 2; //发布
					this.finishLoading = true;
                    // 还原
					this.dataForm.tags = this.tags.map((item) => {
						return { id: item };
					});
					admPostEdit(this.dataForm).then((resp) => {
						this.finishLoading = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: "保存成功",
								onClose: () => {
									this.init();
								}
							});
						} else {
							this.$Message.error({ content: resp.msg, duration: 3 });
						}
					});
				}
			});
		}
	},

	mounted() {
		this.dataForm.id = parseInt(this.$route.params.id, 10);
		this.preinit();
		this.init();
	}
};
</script>
