
<template>
	<Card :bordered="false" dis-hover>
		<p slot="title">
			<Icon type="ios-cog-outline" /> 配置中心
		</p>
		<Tabs :animated="false">
			<TabPane label="基础配置">
				<div style="max-width:640px">
					<Form ref="dataForm" :model="dataForm" :label-width="120" label-colon label-position="right" :rules="dataRules">
						<!-- <FormItem label="小程序图" prop="mpg_cover">
							<Poptip v-if="denyEdit" trigger="hover" placement="right-start">
								<div slot="content">
									<img :src="dataForm.mpg_cover" style="max-height:300px" alt="加载失败">
								</div>
								<img :src="dataForm.mpg_cover" width="102" height="102" style="border-radius: 4px;" alt="加载失败">
							</Poptip>
							<uimage v-else v-model="dataForm.mpg_cover"></uimage>
						</FormItem> -->
						<FormItem label="网站地址" prop="site_url">
							<Input v-model="dataForm.site_url" maxlength="64" show-word-limit placeholder="请填写网站地址"></Input>
						</FormItem>
						<FormItem label="Logo地址" prop="logo_url">
							<Input v-model="dataForm.logo_url" maxlength="64" show-word-limit placeholder="请填写Logo地址"></Input>
						</FormItem>
						<FormItem label="网站标题" prop="title">
							<Input v-model="dataForm.title" maxlength="64" show-word-limit placeholder="请填写网站标题"></Input>
						</FormItem>
                        <FormItem label="网站作者" prop="author">
							<Input v-model="dataForm.author" maxlength="64" show-word-limit placeholder="请填写网站作者"></Input>
						</FormItem>
						<FormItem label="SEO关键词" prop="keywords">
							<Input v-model="dataForm.keywords" maxlength="64" show-word-limit placeholder="请填写SEO关键词"></Input>
						</FormItem>
						<FormItem label="SEO描述" prop="description">
							<Input v-model="dataForm.description" type="textarea" :autosize="{minRows: 2,maxRows: 4}" maxlength="128" show-word-limit placeholder="请填写SEO描述词"></Input>
						</FormItem>
						<FormItem label="ICO地址" prop="favicon_url">
							<Input v-model="dataForm.favicon_url" maxlength="64" show-word-limit placeholder="请填写ICO地址"></Input>
						</FormItem>
						<FormItem label="ICP备案" prop="beian_miit">
							<Input v-model="dataForm.beian_miit" maxlength="64" show-word-limit placeholder="请填写ICP备案"></Input>
						</FormItem>
						<FormItem label="公安备案" prop="beian_nism">
							<Input v-model="dataForm.beian_nism" maxlength="64" show-word-limit placeholder="请填写公安备案"></Input>
						</FormItem>
						<FormItem label="版权" prop="copyright">
							<Input v-model="dataForm.copyright" type="textarea" :autosize="{minRows: 2,maxRows: 4}" maxlength="128" show-word-limit placeholder="请填写版权"></Input>
						</FormItem>
						<FormItem>
							<Poptip confirm title="确定要修改吗？" @on-ok="emitEdit" @on-cancel="emitReset">
								<Button type="warning" :loading="loading">提交保存</Button>
							</Poptip>
							<Button type="success" @click="emitReset" style="margin-left: 10px">重置填写</Button>
						</FormItem>
					</Form>
				</div>
			</TabPane>
			<TabPane label="统计评论">
				<div style="max-width:640px">
					<Form :model="dataForm" :label-width="120" label-colon label-position="right">
						<FormItem label="统计配置" prop="analytic">
							<Input type="textarea" :autosize="{minRows: 8,maxRows: 12}" v-model="dataForm.analytic" maxlength="512" show-word-limit placeholder="请填写自定义css"></Input>
						</FormItem>
						<FormItem label="评论配置" prop="comment">
							<Input type="textarea" :autosize="{minRows: 8,maxRows: 12}" v-model="dataForm.comment" maxlength="512" show-word-limit placeholder="请填写自定义js"></Input>
						</FormItem>
						<FormItem>
							<Poptip confirm title="确定要修改吗？" @on-ok="emitEdit" @on-cancel="emitReset">
								<Button type="warning" :loading="loading">提交保存</Button>
							</Poptip>
							<Button type="success" @click="emitReset" style="margin-left: 10px">重置填写</Button>
						</FormItem>
					</Form>
				</div>
			</TabPane>
			<TabPane label="个性化配置">
				<div style="max-width:640px">
					<Form :model="dataForm" :label-width="120" label-colon label-position="right">
						<FormItem label="githu地址" prop="github_url">
							<Input v-model="dataForm.github_url" maxlength="64" show-word-limit placeholder="请填写githu地址"></Input>
						</FormItem>
						<FormItem label="微博地址" prop="weibo_url">
							<Input v-model="dataForm.weibo_url" maxlength="64" show-word-limit placeholder="请填写微博地址"></Input>
						</FormItem>
						<FormItem label="分页大小" prop="page_size">
							<Slider v-model="dataForm.page_size" :step="1" :min="1" :max="10" show-stops></Slider>
						</FormItem>
						<FormItem label="CSS" prop="site_css">
							<Input type="textarea" :autosize="{minRows: 8,maxRows: 12}" v-model="dataForm.site_css" maxlength="512" show-word-limit placeholder="请填写自定义css"></Input>
						</FormItem>
						<FormItem label="JS" prop="site_js">
							<Input type="textarea" :autosize="{minRows: 8,maxRows: 12}" v-model="dataForm.site_js" maxlength="512" show-word-limit placeholder="请填写自定义js"></Input>
						</FormItem>
						<FormItem>
							<Poptip confirm title="确定要修改吗？" @on-ok="emitEdit" @on-cancel="emitReset">
								<Button type="warning" :loading="loading">提交保存</Button>
							</Poptip>
							<Button type="success" @click="emitReset" style="margin-left: 10px">重置填写</Button>
						</FormItem>
					</Form>
				</div>
			</TabPane>
		</Tabs>
	</Card>
</template>
<script>
import { apiGlobalGet, admGlobalEdit } from "@/api/home";
import Uimage from "@/components/upload/Uimage.vue";
export default {
	components: { Uimage },
	data() {
		return {
			loading: false,
			denyEdit: true,
			dataForm: {
				id: 0,
				site_url: "",
				logo_url: "",
				title: "",
                author:"",
				keywords: "",
				description: "",
				favicon_url: "",
				beian_miit: "",
				beian_nism: "",
				copyright: "",
				site_js: "",
				site_css: "",
				page_size: 0,
				analytic: "",
				comment: "",
				github_url: "",
				weibo_url: ""
			},
			dataRules: {
				site_url: [{ required: true, message: "请填写网站地址", trigger: "blur", max: 128 }],
				logo_url: [{ required: true, message: "请填写Logo地址", trigger: "blur", max: 128 }],
				title: [{ required: true, message: "请填写网站标题", trigger: "blur", max: 128 }],
                author: [{ required: true, message: "请填写网站作者", trigger: "blur", max: 128 }],
				favicon_url: [{ required: true, message: "请填写ICO地址", trigger: "blur", max: 128 }]
			}
		};
	},
	methods: {
		init() {
			apiGlobalGet().then((resp) => {
				this.denyEdit = true;
				if (resp.code == 200) {
					this.dataForm = resp.data;
				} else {
					this.$Message.error({ content: resp.msg, duration: 3 });
				}
			});
		},
		emitReset() {
			this.init();
		},
		emitEdit() {
			this.$refs.dataForm.validate((valid) => {
				if (valid) {
					this.loading = true;
					admGlobalEdit(this.dataForm).then((resp) => {
						this.loading = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: "修改成功",
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
		this.init();
	}
};
</script>
