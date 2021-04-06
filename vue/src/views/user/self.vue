<template>
	<div>
		<Card dis-hover>
			<p slot="title">
				<Icon type="ios-person-outline" /> 个人信息
			</p>
			<div style="max-width:520px">
				<Alert closable type="error">部分信息「用户名」修改后刷新后有效</Alert>
				<Form ref="dataForm" :model="dataForm" :label-width="120" label-colon label-position="right" :rules="dataRules">
					<FormItem label="登录账号" prop="num">
						<Input readonly disabled v-model="dataForm.num"> </Input>
					</FormItem>
					<FormItem label="用户姓名" prop="name">
						<Input v-model="dataForm.name"></Input>
					</FormItem>
					<FormItem label="电话" prop="phone">
						<Input v-model="dataForm.phone"></Input>
					</FormItem>
					<FormItem label="邮箱地址" prop="email">
						<Input v-model="dataForm.email"></Input>
					</FormItem>
					<FormItem label="登录密码">
						<Button type="text" size="small" @click="modalPass = true">修改密码</Button>
					</FormItem>
					<FormItem label="创建时间">
						<span>{{ dataForm.ctime.replace(/T|\+08:00/g, " ")}}</span>
					</FormItem>
					<FormItem>
						<Button type="warning" :loading="loading" @click="emitPasswd">提交保存</Button>
						<Button type="success" @click="init()" style="margin-left: 8px">重置填写</Button>
					</FormItem>
				</Form>
			</div>
		</Card>
		<Modal v-model="modalPass" :closable='false' :mask-closable=false :width="500">
			<h3 slot="header" style="color:#2D8CF0">
				<Icon type="ios-eye-off-outline" /> 修改密码</h3>
			<Form ref="passForm" :model="passForm" :label-width="100" label-position="right" :rules="passRules">
				<FormItem label="原密码" prop="oldPass" :error="oldPassError">
					<Input v-model="passForm.oldPass" type="password" placeholder="请输入现在使用的密码"></Input>
				</FormItem>
				<FormItem label="新密码" prop="newPass">
					<Input v-model="passForm.newPass" type="password" placeholder="请输入新密码，至少6位字符"></Input>
				</FormItem>
				<FormItem label="确认新密码" prop="rePass">
					<Input v-model="passForm.rePass" type="password" placeholder="请再次输入新密码"></Input>
				</FormItem>
			</Form>
			<div slot="footer">
				<Button type="warning" :loading="loadingPass" @click="emit">提交保存</Button>
				<Button type="success" @click="resetPasswd('passForm')" style="margin-left: 8px">重置填写</Button>
				<Button type="info" @click="cancelPass()" style="margin-left: 8px">取消关闭</Button>
			</div>
		</Modal>
	</div>
</template>
<script>
import md5 from "js-md5";
import { apiUserExist } from "@/api/user";
import { admAuthGet, admAuthEditInfo, admAuthEditPasswd } from "@/api/auth";
export default {
	data() {
		return {
			dataForm: { id: 1, name: "", num: "", phone: "", ctime: "" },
			loading: false,
			modalPass: false, // 修改密码模态框显示
			loadingPass: false,
			oldPassError: "",
			dataRules: {
				num: [
					{ required: true, message: "请输入账号", trigger: "blur" },
					{ min: 5, message: "请至少输入5个字符", trigger: "blur" },
					{ max: 32, message: "最多输入32个字符", trigger: "blur" },
					{ validator: this.valideNum, trigger: "blur" }
				],
				name: [{ required: true, message: "请输入姓名", trigger: "blur" }]
			},
			passForm: { oldPass: "", newPass: "", rePass: "" },
			passRules: {
				oldPass: [{ required: true, message: "请输入原密码", trigger: "blur" }],
				newPass: [
					{ required: true, message: "请输入新密码", trigger: "blur" },
					{ min: 6, message: "请至少输入6个字符", trigger: "blur" },
					{ max: 32, message: "最多输入32个字符", trigger: "blur" }
				],
				rePass: [
					{ required: true, message: "请再次输入新密码", trigger: "blur" },
					{ validator: this.valideRePassword, trigger: "blur" }
				]
			}
		};
	},
	methods: {
		valideRePassword(rule, value, callback) {
			if (value !== this.passForm.newPass) {
				callback(new Error("两次输入密码不一致"));
			} else {
				callback();
			}
		},
		init() {
			admAuthGet().then(resp => {
				if (resp.code == 200) {
					this.dataForm = resp.data;
				} else {
					this.$Message.error({
						content: resp.msg,
						duration: 3
					});
				}
			});
		},
		resetPasswd(form) {
			this.$refs.passForm.resetFields();
		},
		emitPasswd() {
			this.$refs.dataForm.validate(valid => {
				if (valid) {
					this.loading = true;
					admAuthEditInfo(this.dataForm).then(resp => {
						this.loading = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: "信息修改成功",
								onClose: () => {
									this.loadingPass = false;
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
		cancelPass() {
			this.modalPass = false;
			this.$refs.passForm.resetFields();
		},
		emit() {
			this.$refs.passForm.validate(valid => {
				if (valid) {
					this.loadingPass = true;
					let passObj = {
						id: this.dataForm.id,
						opass: md5(this.passForm.oldPass).substr(1, 30),
						npass: md5(this.passForm.newPass).substr(1, 30)
					};
					admAuthEditPasswd(passObj).then(resp => {
						this.loadingPass = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: "密码修改成功,请重新登陆",
								onClose: () => {
									localStorage.clear();
									sessionStorage.clear();
									this.$router.push({
										name: "login"
									});
								}
							});
						} else {
							this.$Message.error({
								content: resp.msg,
								duration: 2
							});
						}
					});
				}
			});
		}
	},
	mounted() {
		this.init();
	}
};
</script>
