<template>
	<Card dis-hover>
		<p slot="title">
			<Icon type="md-add" /> 添加角色
		</p>
		<div style="max-width:520px">
			<Form ref="dataForm" :model="dataForm" :label-width="120" label-colon label-position="right" :rules="dataRules">
				<FormItem label="角色名称" prop="name">
					<Input v-model="dataForm.name" placeholder="请填写角色名称"></Input>
				</FormItem>
				<FormItem label="角色介绍" prop="intro">
					<Input v-model="dataForm.intro" type="textarea" :autosize="{minRows: 2,maxRows: 4}" maxlength="128" show-word-limit placeholder="请填写角色介绍"></Input>
				</FormItem>
				<FormItem>
					<Button type="warning" :loading="loading" @click="emitAdd">提交保存</Button>
					<Button type="success" @click="emitReset" style="margin-left: 8px">重置填写</Button>
				</FormItem>
			</Form>
		</div>
	</Card>
</template>
<script>
import { admRoleAdd } from "@/api/role";
export default {
	data() {
		return {
			loading: false,
			dataForm: { id: 0, name: "", intro: "" },
			dataRules: {
				name: [{ required: true, message: "请填写角色名", trigger: "blur", max: 128 }],
				intro: [{ required: true, message: "请填写角色详细", trigger: "blur", max: 128 }]
			}
		};
	},
	methods: {
		emitReset() {
			this.$refs.dataForm.resetFields();
		},
		emitAdd() {
			this.$refs["dataForm"].validate((valid) => {
				if (valid) {
					this.loading = true;
					admRoleAdd(this.dataForm).then((resp) => {
						this.loading = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: "添加成功",
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
		//
	}
};
</script>
