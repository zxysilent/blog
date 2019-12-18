
<template>
	<Card :bordered="false" dis-hover>
		<p slot="title">
			<Icon type="ios-add-circle-outline" /> 添加分类
		</p>
		<div style="max-width:520px">
			<Form ref="dataForm" :model="dataForm" :rules="dataRules" label-position="top">
				<FormItem label="分类名称" prop="name">
					<Input v-model="dataForm.name" placeholder="请填写分类名"></Input>
				</FormItem>
				<FormItem label="分类介绍" prop="intro">
					<Input v-model="dataForm.intro" placeholder="请填写分类介绍"></Input>
				</FormItem>
				<div>
					<Button type="warning" :loading="saveLoading" @click="cmtSave">保&nbsp;&nbsp;&nbsp;&nbsp;存</Button>
				</div>
			</Form>
		</div>
	</Card>
</template>
<script>
import { admCateAdd } from "@/api/cate";
export default {
	data() {
		return {
			dataForm: { name: "", intro: "" },
			dataRules: {
				name: [{ required: true, message: "请填写分类名", trigger: "blur", max: 64 }],
				intro: [{ required: true, message: "请填写分类介绍", trigger: "blur", max: 64 }]
			},
			saveLoading: false
		};
	},
	methods: {
		cmtSave() {
			this.$refs["dataForm"].validate(valid => {
				if (valid) {
					this.saveLoading = true;
					admCateAdd(this.dataForm).then(resp => {
						this.saveLoading = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: "分类信息添加成功",
								onClose: () => {
									this.$router.push({
										name: "cate-list"
									});
								}
							});
						} else {
							this.$Message.error({ content: `分类信息添加失败,请重试`, duration: 3 });
						}
					});
				}
			});
		}
	}
};
</script>
