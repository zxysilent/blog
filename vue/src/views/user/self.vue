<template>
	<div>
		<Card dis-hover>
			<p slot="title">
				<Icon type="ios-person-outline" /> 个人信息
			</p>
			<div style="max-width:640px">
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
						<Button type="text" size="small" @click="modalPasswd = true">修改密码</Button>
					</FormItem>
					<FormItem label="创建时间">
						<span>{{ dataForm.ctime.replace(/T|\+08:00/g, " ")}}</span>
					</FormItem>
					<FormItem>
						<Button type="warning" :loading="loading" @click="emitEdit">提交保存</Button>
						<Button type="success" @click="init()" style="margin-left: 10px">重置填写</Button>
					</FormItem>
				</Form>
			</div>
		</Card>
		<Modal v-model="modalPasswd" :closable='false' :mask-closable=false :width="500">
			<h3 slot="header" style="color:#2D8CF0">
				<Icon type="ios-eye-off-outline" /> 修改密码
			</h3>
			<Form ref="passwdForm" :model="passwdForm" :label-width="100" label-position="right" :rules="passwdRules">
				<FormItem label="原密码" prop="oldPasswd">
					<Input v-model="passwdForm.oldPasswd" type="password" placeholder="请输入现在使用的密码"></Input>
				</FormItem>
				<FormItem label="新密码" prop="newPasswd">
					<Input v-model="passwdForm.newPasswd" type="password" password placeholder="请输入新密码，至少6位字符"></Input>
				</FormItem>
				<FormItem label="确认新密码" prop="rePass">
					<Input v-model="passwdForm.rePass" type="password" password placeholder="请再次输入新密码"></Input>
				</FormItem>
			</Form>
			<div slot="footer">
				<Button type="warning" :loading="loadingPasswd" @click="emitPasswd">提交保存</Button>
				<Button type="success" @click="resetPasswd" style="margin-left: 10px">重置填写</Button>
				<Button type="info" @click="cancelPasswd" style="margin-left: 10px">取消关闭</Button>
			</div>
		</Modal>
	</div>
</template>
<script>
import md5 from "js-md5";
import { admAuthGet, admAuthEdit, admAuthPasswd } from "@/api/auth";
export default {
	data() {
		return {
			dataForm: { id: 1, name: "", num: "", phone: "", ctime: "" },
			loading: false,
			modalPasswd: false, // 修改密码模态框显示
			loadingPasswd: false,
			dataRules: {
				num: [
					{ required: true, message: "请输入账号", trigger: "blur" },
					{ min: 5, message: "请至少输入5个字符", trigger: "blur" },
					{ max: 32, message: "最多输入32个字符", trigger: "blur" }
				],
				name: [{ required: true, message: "请输入姓名", trigger: "blur" }]
			},
			passwdForm: { oldPasswd: "", newPasswd: "", rePass: "" },
			passwdRules: {
				oldPasswd: [{ required: true, message: "请输入原密码", trigger: "blur" }],
				newPasswd: [
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
			if (value !== this.passwdForm.newPasswd) {
				callback(new Error("两次输入密码不一致"));
			} else {
				callback();
			}
		},
		init() {
			admAuthGet().then((resp) => {
				if (resp.code == 200) {
					this.dataForm = resp.data;
				} else {
					this.$Message.error({ content: resp.msg, duration: 3 });
				}
			});
		},
		resetPasswd() {
			this.$refs.passwdForm.resetFields();
		},
		emitEdit() {
			this.$refs.dataForm.validate((valid) => {
				if (valid) {
					this.loading = true;
					admAuthEdit(this.dataForm).then((resp) => {
						this.loading = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: "信息修改成功",
								onClose: () => {
									this.loadingPasswd = false;
								}
							});
						} else {
							this.$Message.error({ content: resp.msg, duration: 3 });
						}
					});
				}
			});
		},
		cancelPasswd() {
			this.modalPasswd = false;
			this.$refs.passwdForm.resetFields();
		},
		emitPasswd() {
			this.$refs.passwdForm.validate((valid) => {
				if (valid) {
					this.loadingPasswd = true;
					let passObj = {
						id: this.dataForm.id,
						opasswd: md5(this.passwdForm.oldPasswd).substr(1, 30),
						npasswd: md5(this.passwdForm.newPasswd).substr(1, 30)
					};
					admAuthPasswd(passObj).then((resp) => {
						this.loadingPasswd = false;
						if (resp.code == 200) {
							this.$Message.success({
								content: "密码修改成功,请重新登陆",
								onClose: () => {
									localStorage.clear();
									sessionStorage.clear();
									this.$router.push({ name: "login" });
								}
							});
						} else {
							this.$Message.error({ content: resp.msg, duration: 2 });
						}
					});
				}
			});
		}
	},
	created() {
		this.init();
	}
};
</script>
