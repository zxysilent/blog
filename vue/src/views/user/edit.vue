<template>
	<div>
		<Card dis-hover>
			<p slot="title">
				<Icon type="ios-create-outline" /> 修改用户
			</p>
			<div style="max-width:640px">
				<Form ref="dataForm" :model="dataForm" :label-width="120" label-colon label-position="right" :rules="dataRules">
					<FormItem label="登录账号" prop="num">
						<Input readonly disabled v-model="dataForm.num"> </Input>
					</FormItem>
					<FormItem label="用户姓名" prop="name">
						<Input v-model="dataForm.name"></Input>
					</FormItem>
					<FormItem label="用户电话" prop="phone">
						<Input v-model="dataForm.phone"></Input>
					</FormItem>
					<FormItem label="邮箱地址" prop="email">
						<Input v-model="dataForm.email"></Input>
					</FormItem>
					<FormItem>
						<Button type="warning" :loading="loading" @click="emitEdit">提交保存</Button>
						<Button type="success" @click="emitReset()" style="margin-left: 10px">重置填写</Button>
						<Button :to="{name:'user-list'}" style="margin-left: 10px">返回列表</Button>
					</FormItem>
				</Form>
			</div>
		</Card>
	</div>
</template>
<script>
import { admUserGet, admUserEdit } from "@/api/user";
export default {
	data() {
		return {
			dataForm: { name: "", num: "", phone: "", id: 0 },
			roles: [],
			loading: false,
			dataRules: {
				name: [{ required: true, message: "请输入姓名", trigger: "blur" }]
			}
		};
	},
	methods: {
		init() {
			admUserGet({ id: this.dataForm.id }).then((resp) => {
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
					admUserEdit(this.dataForm).then((resp) => {
						this.loading = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: "修改成功",
								onClose: () => {
									this.init();
								}
							});
						} else {
							this.$Message.error({
								content: resp.msg,
								duration: 3
							});
						}
					});
				}
			});
		}
	},
	created() {
		this.dataForm.id = parseInt(this.$route.params.id, 10);
		this.init();
	}
};
</script>
