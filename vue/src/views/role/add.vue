
<template>
	<Card :bordered="false" dis-hover>
		<p slot="title">
			<Icon type="ios-add-circle-outline" /> 添加角色
		</p>
		<div style="max-width:520px">
			<Form ref="dataForm" :model="dataForm" :rules="dataRules" label-position="top">
				<FormItem label="角色名称" prop="name">
					<Input v-model="dataForm.name" placeholder="请填写角色名称"></Input>
				</FormItem>
				<FormItem label="角色介绍" prop="intro">
					<Input v-model="dataForm.intro" placeholder="请填写角色介绍"></Input>
				</FormItem>
				<FormItem>
					<Button type="warning" :loading="loading" @click="emitAdd">提交保存</Button>
					<Button type="success" @click="emitReset()" style="margin-left: 8px">重置填写</Button>
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
			dataForm: { name: "", intro: "" },
			dataRules: {
				name: [{ required: true, message: "请填写角色名称", trigger: "blur", max: 128 }],
				intro: [{ required: true, message: "请填写角色介绍", trigger: "blur", max: 128 }]
			},
			loading: false
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
									this.$router.push({ name: "role-list" });
								}
							});
						} else {
							this.$Message.error({ content: resp.msg, duration: 3 });
						}
					});
				}
			});
		}
	}
};
</script>
