
<template>
    <div>
        <Card>
            <p slot="title">
                <Icon type="ios-settings-outline" /> 基本设置
            </p>
            <div>
                <Form ref="userForm" :model="userForm" label-position="top" :rules="inforValidate">
                    <FormItem label="站点名称" prop="name">
                        <div style="display:inline-block;width:600px;">
                            <Input v-model="userForm.name"></Input>
                        </div>
                    </FormItem>
                    <FormItem label="Logo" prop="cellphone">
                        <div style="display:inline-block;width:600px;">
                            <Input v-model="userForm.cellphone" @on-keydown="hasChangePhone"></Input>
                        </div>
                    </FormItem>
                    <FormItem label="站点描述" prop="name">
                        <div style="display:inline-block;width:600px;">
                            <Input v-model="userForm.name"></Input>
                        </div>
                    </FormItem>
                    <FormItem label="站点地址" prop="name">
                        <div style="display:inline-block;width:600px;">
                            <Input v-model="userForm.name"></Input>
                        </div>
                    </FormItem>
                    <FormItem label="Favicon" prop="name">
                        <div style="display:inline-block;width:600px;">
                            <Input v-model="userForm.name"></Input>
                        </div>
                    </FormItem>
                    <FormItem label="关键词" prop="name">
                        <div style="display:inline-block;width:600px;">
                            <Input v-model="userForm.name"></Input>
                        </div>
                    </FormItem>
                    <FormItem label="Github" prop="name">
                        <div style="display:inline-block;width:600px;">
                            <Input v-model="userForm.name"></Input>
                        </div>
                    </FormItem>
                    <FormItem label="Weibo" prop="name">
                        <div style="display:inline-block;width:600px;">
                            <Input v-model="userForm.name"></Input>
                        </div>
                    </FormItem>
                    <FormItem label="备案号" prop="name">
                        <div style="display:inline-block;width:600px;">
                            <Input v-model="userForm.name"></Input>
                        </div>
                    </FormItem>
                    <div>
                        <Button type="warning" :loading="saveLoading" @click="cmtSave">保&nbsp;&nbsp;&nbsp;&nbsp;存</Button>
                    </div>
                </Form>
            </div>
        </Card>
    </div>
</template>

<script>
export default {
	name: "ownspace_index",
	data() {
		return {
			userForm: {
				name: "",
				cellphone: "",
				company: "",
				department: ""
			},
			saveLoading: false,
			identifyError: "", // 验证码错误
			editPasswordModal: false, // 修改密码模态框显示
			savePassLoading: false,
			oldPassError: "",
			identifyCodeRight: false, // 验证码是否正确
			checkIdentifyCodeLoading: false,
			editPasswordForm: {
				oldPass: "",
				newPass: "",
				rePass: ""
			},
			inputCodeVisible: false, // 显示填写验证码box
			initPhone: "",
			gettingIdentifyCodeBtnContent: "获取验证码" // “获取验证码”按钮的文字
		};
	},
	methods: {
		cancelEditUserInfor() {
			this.$store.commit("removeTag", "ownspace_index");
			localStorage.pageOpenedList = JSON.stringify(
				this.$store.state.app.pageOpenedList
			);
			let lastPageName = "";
			if (this.$store.state.app.pageOpenedList.length > 1) {
				lastPageName = this.$store.state.app.pageOpenedList[1].name;
			} else {
				lastPageName = this.$store.state.app.pageOpenedList[0].name;
			}
			this.$router.push({
				name: lastPageName
			});
		},
		cmtSave() {
			this.$refs["userForm"].validate(valid => {
				if (valid) {
					if (
						this.phoneHasChanged &&
						this.userForm.cellphone !== this.initPhone
					) {
						// 手机号码修改过了而且修改之后的手机号和原来的不一样
					} else {
						this.saveInfoAjax();
					}
				}
			});
		},
		init() {
			this.userForm.name = "";
			this.userForm.cellphone = "17712345678";
			this.initPhone = "17712345678";
			this.userForm.company = "TalkingData";
			this.userForm.department = "可视化部门";
		},
		saveInfoAjax() {
			this.save_loading = true;
			setTimeout(() => {
				this.$Message.success("保存成功");
				this.save_loading = false;
			}, 1000);
		}
	},
	mounted() {
		this.init();
	}
};
</script>
