
<template>
	<Card :bordered="false" dis-hover>
		<p slot="title">
			<Icon type="ios-add-circle-outline" /> 添加标签
		</p>
		<div style="max-width:640px">
			<Form ref="dataForm" :model="dataForm" :rules="dataRules" label-position="top">
				<FormItem label="标签名称" prop="name">
					<Input v-model="dataForm.name" maxlength="64" show-word-limit placeholder="请填写标签名称"></Input>
				</FormItem>
				<FormItem label="标签介绍" prop="intro">
					<Input v-model="dataForm.intro" type="textarea" :autosize="{minRows:2,maxRows:4}" maxlength="128" show-word-limit show-word-limit placeholder="请填写标签介绍"></Input>
				</FormItem>
				<FormItem>
					<Button type="warning" :loading="loading" @click="emitAdd">提交保存</Button>
					<Button type="success" @click="emitReset()" style="margin-left: 10px">重置填写</Button>
					<Button :to="{name:'tag-list'}" style="margin-left: 10px">返回列表</Button>
				</FormItem>
			</Form>
		</div>
	</Card>
</template>
<script>
import { admTagAdd } from "@/api/tag";
export default {
	data() {
		return {
			dataForm: { name: "", intro: "" },
			dataRules: {
				name: [{ required: true, message: "请填写标签名", trigger: "blur", max: 64 }],
				intro: [{ required: true, message: "请填写标签介绍", trigger: "blur", max: 128 }]
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
					admTagAdd(this.dataForm).then((resp) => {
						this.loading = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: "添加成功",
								onClose: () => {
									this.$router.push({ name: "tag-list" });
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
