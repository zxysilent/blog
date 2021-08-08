<template>
	<Card dis-hover>
		<p slot="title">
			<Icon type="ios-create-outline" /> 修改字典
		</p>
		<div style="max-width:520px">
			<Form ref="dataForm" :model="dataForm" :label-width="120" label-colon label-position="right" :rules="dataRules">
				<FormItem label="字典名" prop="key">
					<Input v-model="dataForm.key" readonly disabled placeholder="请填写字典名"><span slot="append">不可修改</span></Input>
				</FormItem>
				<FormItem label="字典值" prop="value">
					<Input v-model="dataForm.value" type="textarea" :autosize="{minRows: 4,maxRows: 6}" maxlength="255" show-word-limit placeholder="请填写字典值"></Input>
				</FormItem>
				<FormItem label="字典说明" prop="intro">
					<Input v-model="dataForm.intro" type="textarea" :autosize="{minRows: 2,maxRows: 4}" maxlength="255" show-word-limit placeholder="请填写字典说明"></Input>
				</FormItem>
				<FormItem>
					<Button type="warning" :loading="loading" @click="emitEdit">提交保存</Button>
					<Button type="success" @click="emitReset" style="margin-left: 8px">重置填写</Button>
					<Button :to="{name:'dict-list'}" style="margin-left: 8px">返回列表</Button>
				</FormItem>
			</Form>
		</div>
	</Card>
</template>
<script>
import { apiDictGet, admDictEdit } from "@/api/dict";
export default {
	data() {
		return {
			loading: false,
			dataForm: { key: "", value: "", intro: "" },
			dataRules: {
				key: [{ required: true, message: "请填写字典名", trigger: "blur", max: 64 }],
				value: [{ required: true, message: "请填写字典值", trigger: "blur", max: 255 }]
			}
		};
	},
	methods: {
		init() {
			apiDictGet({ key: this.dataForm.key }).then((resp) => {
				if (resp.code == 200) {
					this.dataForm = resp.data;
				}
			});
		},
		emitReset() {
			this.init();
		},
		emitEdit() {
			this.$refs["dataForm"].validate((valid) => {
				if (valid) {
					this.loading = true;
					admDictEdit(this.dataForm).then((resp) => {
						this.loading = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: "修改成功",
								onClose: () => {
									this.emitReset();
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
		this.dataForm.key = this.$route.params.key;
		this.init();
	}
};
</script>
