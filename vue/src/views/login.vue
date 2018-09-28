<style lang="less">
@import "./login.less";
</style>

<template>
    <div class="login" @keydown.enter="handleSubmit">
        <div class="login-con">
            <Card :bordered="false">
                <p slot="title">
                    <Icon type="log-in"></Icon>
                    欢迎登录
                </p>
                <div class="form-con">
                    <Form ref="loginForm" :model="user" :rules="rules">
                        <FormItem prop="num">
                            <Input type="text" v-model="user.num" placeholder="请输入用户名">
                            <Icon type="ios-person-outline" slot="prepend"></Icon>
                            </Input>
                        </FormItem>
                        <FormItem prop="pass">
                            <Input type="password" v-model="user.pass" placeholder="请输入密码">
                            <Icon type="ios-lock-outline" slot="prepend"></Icon>
                            </Input>
                        </FormItem>
                        <FormItem>
                            <Button @click="handleSubmit" type="primary" long>登录</Button>
                        </FormItem>
                    </Form>
                    <p class="login-tip">官网
                        <a href="/" title="author">
                            <Icon type="ios-send-outline" />
                        </a>
                    </p>
                </div>
            </Card>
        </div>
    </div>
</template>

<script>
import md5 from "js-md5";
import { login } from "@/api/auth";
import util from "@/libs/util";
export default {
  data() {
    return {
      user: {
        num: "zxysilent",
        pass: "zxyslt"
      },
      rules: {
        num: [{ required: true, message: "账号不能为空", trigger: "blur" }],
        pass: [{ required: true, message: "密码不能为空", trigger: "blur" }]
      }
    };
  },
  methods: {
    handleSubmit() {
      let that = this;
      that.$refs.loginForm.validate(valid => {
        if (valid) {
          let data = {
            num: that.user.num,
            pass: md5(that.user.pass).substr(1, 30)
          };
          login(data).then(res => {
            if (res.code == 200) {
              this.$Message.success({
                content: "登陆成功",
                onClose: () => {
                  util.setToken(res.data);
                  this.$store.commit("updateMenulist");
                  that.$router.push({ name: "home_index" });
                }
              });
            } else {
              this.$Message.error(res.msg);
            }
          });
        }
      });
    }
  }
};
</script>