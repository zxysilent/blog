
<template>
	<Card :bordered="false" dis-hover>
		<p slot="title">
			<Icon type="ios-hammer-outline" /> 信息配置
		</p>
		<div style="max-width:768px">
			<Form ref="dataForm" :model="dataForm" :label-width="120" label-colon label-position="right" :rules="dataRules">
				<FormItem label="学校名" prop="school_name">
					<Input :disabled="denyEdit" v-model="dataForm.school_name" maxlength="32" show-word-limit placeholder="请填写学校名"></Input>
				</FormItem>
				<Row>
					<Col span="12">
					<FormItem label="学校电话" prop="school_phone">
						<Input :disabled="denyEdit" v-model="dataForm.school_phone" maxlength="32" show-word-limit placeholder="请填写学校电话"></Input>
					</FormItem>
					</Col>
					<Col span="12">
					<FormItem label="学校地址" prop="school_addr">
						<Input :disabled="denyEdit" v-model="dataForm.school_addr" maxlength="32" show-word-limit placeholder="请填写学校地址"></Input>
					</FormItem>
					</Col>
				</Row>
				<Row>
					<Col span="12">
					<FormItem label="说明标题" prop="school_aim">
						<Input type="textarea" :autosize="{minRows: 3,maxRows: 6}" :disabled="denyEdit" v-model="dataForm.school_aim" maxlength="128" show-word-limit placeholder="请填写说明标题"></Input>
					</FormItem>
					</Col>
					<Col span="12">

					<FormItem label="筛选提示" prop="check_msg">
						<Input type="textarea" :autosize="{minRows: 3,maxRows: 6}" :disabled="denyEdit" v-model="dataForm.check_msg" maxlength="128" show-word-limit placeholder="请填写筛选提示"></Input>
					</FormItem>
					</Col>
				</Row>
				<Row>
					<Col span="12">
					<FormItem label="筛选成功" prop="check_succ">
						<Input type="textarea" :autosize="{minRows: 3,maxRows: 6}" :disabled="denyEdit" v-model="dataForm.check_succ" maxlength="128" show-word-limit placeholder="请填写筛选成功"></Input>
					</FormItem>
					</Col>
					<Col span="12">
					<FormItem label="筛选失败" prop="check_fail">
						<Input type="textarea" :autosize="{minRows: 3,maxRows: 6}" :disabled="denyEdit" v-model="dataForm.check_fail" maxlength="128" show-word-limit placeholder="请填写筛选失败"></Input>
					</FormItem>
					</Col>
				</Row>
				<FormItem label="学校简介" prop="school_intro">
					<Input type="textarea" :autosize="{minRows: 4,maxRows: 12}" :disabled="denyEdit" v-model="dataForm.school_intro" maxlength="512" show-word-limit placeholder="请填写学校简介"></Input>
				</FormItem>
				<FormItem label="使用说明" prop="school_usage">
					<Input type="textarea" :autosize="{minRows: 4,maxRows: 8}" :disabled="denyEdit" v-model="dataForm.school_usage" maxlength="512" show-word-limit placeholder="请填写使用说明"></Input>
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
import { apiGlobalGet, admGlobalEditUsage } from "@/api/home";
import Uimage from "@/components/upload/Uimage.vue";
export default {
	components: { Uimage },
	data() {
		return {
			loading: false,
			denyEdit: true,
			dataForm: {
				id: 0,
				school_name: "",
				school_intro: "",
				school_phone: "",
				school_addr: "",
				school_aim: "",
				school_usage: "",
				check_msg: "",
				check_succ: "",
				check_fail: ""
			},
			dataRules: {
				school_name: [{ required: true, message: "请填写学校名", trigger: "blur", max: 128 }],
				school_intro: [{ required: true, message: "请填写学校简介", trigger: "blur", max: 1024 }],
				school_phone: [{ required: true, message: "请填写学校电话", trigger: "blur", max: 128 }],
				school_addr: [{ required: true, message: "请填写学校地址", trigger: "blur", max: 128 }],
				school_aim: [{ required: true, message: "请填写使用标题", trigger: "blur", max: 128 }],
				school_usage: [{ required: true, message: "请填写使用说明", trigger: "blur", max: 1024 }],
				check_msg: [{ required: true, message: "请填写筛选提示", trigger: "blur", max: 128 }],
				check_succ: [{ required: true, message: "请填写筛选成功", trigger: "blur", max: 128 }],
				check_fail: [{ required: true, message: "请填写筛选失败", trigger: "blur", max: 128 }]
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
					admGlobalEditUsage(this.dataForm).then((resp) => {
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
