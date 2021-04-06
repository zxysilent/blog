<template>
	<div>
		<Card dis-hover>
			<p slot="title">
				<Icon type="md-add" /> 添加用户
			</p>
			<div style="max-width:520px">
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
					<FormItem label="用户角色" prop="role">
						<RadioGroup v-model="dataForm.role">
							<Radio v-for="item in roleArr" :label="item.role" :key="item.role">{{item.tip}}</Radio>
						</RadioGroup>
					</FormItem>
					<FormItem>
						<Button type="warning" :loading="loading" @click="emit">提交保存</Button>
						<Button type="success" @click="reset()" style="margin-left: 8px">重置填写</Button>
					</FormItem>
				</Form>
			</div>
		</Card>
	</div>
</template>
<script>
import md5 from "js-md5";
import { apiUserExist, admUserAdd } from "@/api/user";
export default {
	data() {
		return {
			auth: {},
			dataForm: { name: "", passRaw: "1q2w3e", num: "", phone: "", role: 0 },
			roleArr: [],
			loading: false,
			dataRules: {
				num: [
					{ required: true, message: "请输入账号", trigger: "blur" },
					{ min: 5, message: "请至少输入5个字符", trigger: "blur" },
					{ max: 32, message: "最多输入32个字符", trigger: "blur" },
					{ validator: this.valideNum, trigger: "blur" }
				],
				name: [{ required: true, message: "请输入姓名", trigger: "blur" }],
				pass: [
					{ min: 6, message: "请至少输入6个字符", trigger: "blur" },
					{ max: 32, message: "最多输入32个字符", trigger: "blur" }
				]
			}
		};
	},
	methods: {
		valideNum(rule, value, callback) {
			apiUserExist(value).then(resp => {
				if (resp.code == 200) {
					callback(new Error("当前账号已经存在"));
				} else {
					callback();
				}
			});
		},
		reset() {
			this.$refs.dataForm.resetFields();
			if (this.auth.role == 8192) {
				this.dataForm.role = 4096;
			} else if (this.auth.role == 4096) {
				this.dataForm.role = 2048;
			}
		},
		emit() {
			this.$refs.dataForm.validate(valid => {
				if (valid) {
					this.loading = true;
					this.dataForm.pass = md5(this.dataForm.passRaw).substr(1, 30);
					admUserAdd(this.dataForm).then(resp => {
						this.loading = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: "添加成功",
								onClose: () => {
									this.reset();
									this.showEnt = false;
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
		},
		genRoleArr(rl) {
			let arr = [];
			if (rl == 8192) {
				arr.push({ role: 4096, tip: "系统管理员" });
				arr.push({ role: 2048, tip: "内容管理员" });
				this.dataForm.role = 4096;
			} else if (rl == 4096) {
				arr.push({ role: 2048, tip: "内容管理员" });
				this.dataForm.role = 2048;
			}
			return arr;
		}
	},
	mounted() {
		this.auth = this.$auth();
		this.roleArr = this.genRoleArr(this.auth .role);
	}
};
</script>
