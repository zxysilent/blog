<template>
	<div>
		<Card dis-hover>
			<p slot="title">
				<Icon type="md-add" /> 添加用户
			</p>
			<div style="max-width:640px">
				<Form ref="dataForm" :model="dataForm" :label-width="120" label-colon label-position="right" :rules="dataRules">
					<FormItem label="登录账号" prop="num">
						<Input v-model="dataForm.num"><span slot="append">不可修改</span></Input>
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
					<FormItem label="登陆密码" prop="passRaw">
						<Input v-model="dataForm.passRaw"><span slot="append">默认1q2w3e</span></Input>
					</FormItem>
					<FormItem label="用户角色" prop="role_id">
						<RadioGroup v-model="dataForm.role_id">
							<Radio v-for="item in roleAll" :label="item.id" :key="item.role">{{item.name}}</Radio>
						</RadioGroup>
					</FormItem>
					<FormItem>
						<Button type="warning" :loading="loading" @click="emitAdd">提交保存</Button>
						<Button type="success" @click="emitReset()" style="margin-left: 10px">重置填写</Button>
                        <Button :to="{name:'user-list'}" style="margin-left: 10px">返回列表</Button>
					</FormItem>
				</Form>
			</div>
		</Card>
	</div>
</template>
<script>
import md5 from "js-md5";
import { apiUserExist, admUserAdd } from "@/api/user";
import { admRoleAll } from "@/api/role";
export default {
	data() {
		return {
			auth: {},
			roleAll: [],
			dataForm: { name: "", passRaw: "1q2w3e", num: "", phone: "", role_id: 0 },
			loading: false,
			dataRules: {
				num: [
					{ required: true, message: "请输入账号", trigger: "blur" },
					{ min: 5, message: "请至少输入5个字符", trigger: "blur" },
					{ max: 32, message: "最多输入32个字符", trigger: "blur" },
					{ validator: this.valideNum, trigger: "blur" }
				],
				name: [{ required: true, message: "请输入姓名", trigger: "blur" }],
				role_id: [{ type: "number", required: true, message: "请选择用户角色", trigger: "change", min: 1 }],
				passRaw: [
					{ min: 6, message: "请至少输入6个字符", trigger: "blur" },
					{ max: 32, message: "最多输入32个字符", trigger: "blur" }
				]
			}
		};
	},
	methods: {
		preinit() {
			admRoleAll().then((resp) => {
				if (resp.code == 200) {
					this.roleAll = resp.data;
				} else {
					this.$Message.error({ content: resp.msg, duration: 3 });
				}
			});
		},
		valideNum(rule, value, callback) {
			apiUserExist({ num: value }).then((resp) => {
				if (resp.code == 200) {
					callback(new Error("当前账号已经存在"));
				} else {
					callback();
				}
			});
		},
		emitReset() {
			this.$refs.dataForm.resetFields();
		},
		emitAdd() {
			this.$refs.dataForm.validate((valid) => {
				if (valid) {
					this.loading = true;
					this.dataForm.pass = md5(this.dataForm.passRaw).substr(1, 30);
					admUserAdd(this.dataForm).then((resp) => {
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
	created() {
		this.preinit();
	}
};
</script>
