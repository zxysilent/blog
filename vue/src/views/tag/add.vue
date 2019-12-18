
<template>
	<Card :bordered="false" dis-hover>
		<p slot="title">
			<Icon type="ios-add-circle-outline" /> 添加标签
		</p>
		<div style="max-width:520px">
			<Form ref="dataForm" :model="dataForm" :rules="dataRules" label-position="top">
				<FormItem label="标签名称" prop="name">
					<Input v-model="dataForm.name" placeholder="请填写标签名"></Input>
				</FormItem>
				<FormItem label="标签介绍" prop="intro">
					<Input v-model="dataForm.intro" placeholder="请填写标签介绍"></Input>
				</FormItem>
				<div>
					<Button type="warning" :loading="saveLoading" @click="cmtSave">保&nbsp;&nbsp;&nbsp;&nbsp;存</Button>
				</div>
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
				intro: [{ required: true, message: "请填写标签介绍", trigger: "blur", max: 64 }]
			},
			saveLoading: false
		};
	},
	methods: {
		cmtSave() {
			this.$refs["dataForm"].validate(valid => {
				if (valid) {
					this.saveLoading = true;
					admTagAdd(this.dataForm).then(resp => {
						this.saveLoading = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: "标签信息添加成功",
								onClose: () => {
									this.$router.push({
										name: "tag-list"
									});
								}
							});
						} else {
							this.$Message.error({ content: `标签信息添加失败,请重试`, duration: 3 });
						}
					});
				}
			});
		}
	}
};
</script>
