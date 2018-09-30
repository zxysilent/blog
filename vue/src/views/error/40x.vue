<style lang="less">
@import "./error.less";
</style>

<template>
    <Card :bordered="false" style="width:600px;margin: 0 auto;">
        <div class="error-title">4
            <span class="error-0">
                <Icon type="android-lock"></Icon>
            </span>
            X
        </div>
        <p class="error-message">登陆失效,请重新登陆</p>
        <div class="error-btn">
            <Button @click="goLogin" size="large" style="width: 200px;" type="text">去登陆</Button>
        </div>
    </Card>
</template>

<script>
export default {
  name: "Error403",
  methods: {
    goLogin() {
      this.$store.commit("logout", this);
      this.$router.push({
        name: "login"
      });
    }
  }
};
</script>
