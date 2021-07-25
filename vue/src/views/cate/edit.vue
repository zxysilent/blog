
<template>
	<Card :bordered="false" dis-hover>
		<p slot="title">
			<Icon type="ios-create-outline" /> 修改分类
		</p>
		<div style="max-width:640px">
			<Form ref="dataForm" :model="dataForm" :rules="dataRules" label-position="top">
				<FormItem label="分类名称" prop="name">
					<Input v-model="dataForm.name" maxlength="64" show-word-limit placeholder="请填写分类名称"></Input>
				</FormItem>
				<FormItem label="分类介绍" prop="intro">
					<Input v-model="dataForm.intro" type="textarea" :autosize="{minRows:2,maxRows:4}" maxlength="128" show-word-limit placeholder="请填写分类介绍"></Input>
				</FormItem>
				<FormItem>
					<Button type="warning" :loading="loading" @click="emitEdit">提交保存</Button>
					<Button type="success" @click="emitReset" style="margin-left: 10px">重置填写</Button>
					<Button :to="{name:'cate-list'}" style="margin-left: 10px">返回列表</Button>
				</FormItem>
			</Form>
		</div>
	</Card>
</template>
<script>
import { apiCateGet, admCateEdit } from "@/api/cate";
export default {
	data() {
		return {
			dataForm: { id: 0, name: "", intro: "" },
			dataRules: {
				name: [{ required: true, message: "请填写分名称", trigger: "blur", max: 64 }],
				intro: [{ required: true, message: "请填写分类介绍", trigger: "blur", max: 128 }]
			},
			loading: false
		};
	},
	methods: {
		init() {
			apiCateGet({ id: this.dataForm.id }).then((resp) => {
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
					admCateEdit(this.dataForm).then((resp) => {
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
