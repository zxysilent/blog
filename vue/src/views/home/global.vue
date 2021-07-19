
<template>
	<Card :bordered="false" dis-hover>
		<p slot="title">
			<Icon type="ios-cog-outline" /> 配置中心
		</p>
		<div style="max-width:520px">
			<Form ref="dataForm" :model="dataForm" :label-width="120" label-colon label-position="right" :rules="dataRules">
				<FormItem label="小程序图" prop="mpg_cover">
					<Poptip v-if="denyEdit" trigger="hover" placement="right-start">
						<div slot="content">
							<img :src="dataForm.mpg_cover" style="max-height:300px" alt="加载失败">
						</div>
						<img :src="dataForm.mpg_cover" width="102" height="102" style="border-radius: 4px;" alt="加载失败">
					</Poptip>
					<uimage v-else v-model="dataForm.mpg_cover"></uimage>
				</FormItem>
				<FormItem label="小程序标题" prop="mpg_title">
					<Input :disabled="denyEdit" v-model="dataForm.mpg_title" maxlength="20" show-word-limit placeholder="请填写小程序标题"></Input>
				</FormItem>
				<FormItem label="报名限制" prop="app_allow">
					<InputNumber :disabled="denyEdit" v-model="dataForm.app_allow" style="width:120px" :max="1000" :min="1" :step="1" placeholder="请填写报名限制"></InputNumber>&nbsp;通过年龄限制(有效信息填写人数)
				</FormItem>
				<FormItem label="录取限制" prop="app_limit">
					<InputNumber :disabled="denyEdit" v-model="dataForm.app_limit" style="width:120px" :max="1000" :min="1" :step="1" placeholder="请填写录取限制"></InputNumber>&nbsp;填报完后筛选(最终需要的人数)
				</FormItem>
				<FormItem label="一批次">
					<FormItem label="开始日期" prop="one_min">
						<DatePicker :disabled="denyEdit" type="date" v-model="dataForm.one_min" :editable="false" :clearable="false" placeholder="开始日期"></DatePicker>&nbsp;0点开始
					</FormItem>
					<FormItem label="结束日期" prop="one_max">
						<DatePicker :disabled="denyEdit" type="date" v-model="dataForm.one_max" :editable="false" :clearable="false" placeholder="结束日期"></DatePicker>&nbsp;0点之前
					</FormItem>
				</FormItem>
				<FormItem label="二批次">
					<FormItem label="开始日期" prop="two_min">
						<DatePicker :disabled="denyEdit" type="date" v-model="dataForm.two_min" :editable="false" :clearable="false" placeholder="开始日期"></DatePicker>&nbsp;0点开始
					</FormItem>
					<FormItem label="结束日期" prop="two_max">
						<DatePicker :disabled="denyEdit" type="date" v-model="dataForm.two_max" :editable="false" :clearable="false" placeholder="结束日期"></DatePicker>&nbsp;0点之前
					</FormItem>
				</FormItem>
				<FormItem label="报名通过人数" prop="app_count">
					<span>{{ dataForm.app_count}}</span>
				</FormItem>
				<FormItem label="系统状态" prop="show">
					<i-Switch :disabled="denyEdit" size="large" v-model="dataForm.state"><span slot="open">开启</span>
						<span slot="close">关闭</span>
					</i-Switch>
				</FormItem>
				<FormItem label="关闭提示信息" prop="close_msg">
					<Input :disabled="denyEdit" v-model="dataForm.close_msg" maxlength="20" show-word-limit placeholder="请填写关闭提示信息"></Input>
				</FormItem>
				<FormItem>
					<Button type="success" @click="emitReset">取消刷新</Button>
					<Button v-auth="'global_edit'" type="info" @click="denyEdit=false" style="margin-left: 5px">开启编辑</Button>
					<Poptip v-auth="'global_edit'" confirm title="确定要修改吗？" @on-ok="emitEdit" @on-cancel="emitReset">
						<Button type="warning" :loading="loading" :disabled="denyEdit" style="margin-left: 5px" @click="">提交保存</Button>
					</Poptip>
				</FormItem>
			</Form>
		</div>
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
				mpg_cover: "",
				mpg_title: "",
				app_allow: 0,
				app_limit: 0,
				app_count: 0,
				one_min: new Date(),
				one_max: new Date(),
				two_min: new Date(),
				two_max: new Date(),
				state: true,
				close_msg: ""
			},
			dataRules: {
				mpg_cover: [{ required: true, message: "请填写小程序图", trigger: "blur", max: 200 }],
				mpg_title: [{ required: true, message: "请填写小程序标题", trigger: "blur", max: 200 }],
				app_allow: [{ type: "number", required: true, message: "请填写报名限制", trigger: "change", min: 1 }],
				app_limit: [{ type: "number", required: true, message: "请填写录取限制", trigger: "change", min: 1 }],
				close_msg: [{ required: true, message: "请填写小程序提示信息", trigger: "blur", max: 200 }]
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
