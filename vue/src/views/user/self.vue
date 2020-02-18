<template>
	<div>
		<Card dis-hover>
			<p slot="title">
				<Icon type="ios-person-outline" /> 个人信息
			</p>
			<div style="max-width:520px">
				<Alert closable type="error">部分信息「用户名」修改后刷新后有效</Alert>
				<Form ref="userForm" :model="userForm" :label-width="100" label-position="right" :rules="userRules">
					<FormItem label="用户账号：" prop="num">
						<Input readonly disabled v-model="userForm.num">
						</Input>
					</FormItem>
					<FormItem label="用户姓名：" prop="name">
						<Input v-model="userForm.name"></Input>
					</FormItem>
					<FormItem label="电话：" prop="phone">
						<Input v-model="userForm.phone"></Input>
					</FormItem>
					<FormItem label="邮箱地址：" prop="email">
						<Input v-model="userForm.email"></Input>
					</FormItem>
					<FormItem label="登录密码：">
						<Button type="text" size="small" @click="showModel">修改密码</Button>
					</FormItem>
					<FormItem label="创建时间：">
						<span>{{ userForm.ctime.replace(/T|\+08:00|Z/g, " ") }}</span>
					</FormItem>
					<FormItem label="备注信息：" prop="remark">
						<span>{{userForm.remark}}</span>
					</FormItem>
					<FormItem>
						<Button type="warning" :loading="loadingSaveInfo" @click="saveUser">提交保存</Button>
						<Button type="success" @click="resetEdit()" style="margin-left: 8px">重置填写</Button>
					</FormItem>
				</Form>
			</div>
		</Card>
		<Modal v-model="showPasswordModal" :closable='false' :mask-closable=false :width="500">
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
				<Button type="warning" :loading="loadingSavePass" @click="savePass">提交保存</Button>
				<Button type="success" @click="resetForm('passForm')" style="margin-left: 8px">重置填写</Button>
				<Button type="info" @click="cancelPass()" style="margin-left: 8px">取消关闭</Button>
			</div>
		</Modal>
	</div>
</template>
<script>
import md5 from "js-md5";
import util from "@/utils.js";
import { admAuth } from "@/api/auth";
import { admUserPass, admUserEditSelf } from "@/api/user";
export default {
	data() {
		const valideRePassword = (rule, value, callback) => {
			if (value !== this.passForm.newPass) {
				callback(new Error("两次输入密码不一致"));
			} else {
				callback();
			}
		};
		return {
			userForm: { id: 0, num: "", name: "", email: "", phone: "", remark: "", ctime: "" },
			loadingSaveInfo: false,
			showPasswordModal: false, // 修改密码模态框显示
			loadingSavePass: false,
			oldPassError: "",
			userRules: {
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
					{ validator: valideRePassword, trigger: "blur" }
				]
			}
		};
	},
	methods: {
		showModel() {
			this.showPasswordModal = true;
		},
		resetEdit() {
			this.init();
		},
		resetForm(form) {
			this.$refs[form].resetFields();
		},
		saveUser() {
			this.$refs["userForm"].validate(valid => {
				if (valid) {
					this.loadingSaveInfo = true;
					admUserEditSelf(this.userForm).then(res => {
						this.loadingSaveInfo = false;
						if (res.code == 200) {
							this.$Message.success({
								content: "信息修改成功",
								onClose: () => {
									this.loadingSavePass = false;
								}
							});
						} else {
							this.$Message.error({ content: `信息修改失败,请重试`, duration: 3 });
						}
					});
				}
			});
		},
		cancelPass() {
			this.showPasswordModal = false;
		},
		savePass() {
			this.$refs["passForm"].validate(valid => {
				if (valid) {
					this.loadingSavePass = true;
					let passObj = {
						id: this.userForm.id,
						opass: md5(this.passForm.oldPass).substr(1, 30),
						npass: md5(this.passForm.newPass).substr(1, 30)
					};
					admUserPass(passObj).then(res => {
						this.loadingSavePass = false;
						if (res.code == 200) {
							this.$Message.success({
								content: "密码修改成功,请重新登陆",
								onClose: () => {
									// 退出登录
									// this.$store.commit("logout", this);
									this.$router.push({
										name: "login"
									});
								}
							});
						} else {
							this.$Message.error({ content: `密码修改失败`, duration: 2 });
						}
					});
				}
			});
		},
		init() {
			admAuth().then(res => {
				if (res.code == 200) {
					this.userForm = res.data;
				} else {
					this.$Message.error({ content: `未查询到数据，请重试`, duration: 3 });
				}
			});
		}
	},
	mounted() {
		this.init();
	}
};
</script>
