<style lang="less">
@import "./error.less";
</style>

<template>
    <Card :bordered="false" style="width:600px;margin: 0 auto;">
        <div class="error-title">4
            <span class="error-0">
                <Icon type="android-lock"></Icon>
            </span>
            4
        </div>
        <p class="error-message">YOU&nbsp;&nbsp;LOOK&nbsp;&nbsp;LOST</p>
        <div class="error-btn">
            <Button @click="goHome" size="large" style="width: 200px;" type="text">返回首页</Button>
            <Button @click="backPage" size="large" style="width: 200px;margin-left: 40px;" type="primary">返回上一页</Button>
        </div>
    </Card>
</template>

<script>
export default {
  name: "Error404",
  methods: {
    backPage() {
      this.$router.go(-1);
    },
    goHome() {
      this.$router.push({
        name: "home"
      });
    }
  }
};
</script>
