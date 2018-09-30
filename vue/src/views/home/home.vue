<style lang="less">
@import "./home.less";
</style>
<template>
    <Card dis-hover>
        <Alert closable>An info prompt</Alert>
        <Row :gutter="10">
            <Col :md="24" :lg="8">
            <Card :bordered="false" dis-hover>
                <Row type="flex" class="user-infor">
                    <Col span="8">
                    <Row class-name="made-child-con-middle" type="flex" align="middle">
                        <i-circle :percent="100" stroke-color="#5cb85c">
                            <Icon type="ios-checkmark" size="60" style="color:#5cb85c"></Icon>
                        </i-circle>
                    </Row>
                    </Col>
                    <Col span="16" style="padding-left:6px;">
                    <Row class-name="made-child-con-middle" type="flex" align="middle">
                        <div>
                            <b class="card-user-infor-name">Admin</b>
                            <p>super admin</p>
                        </div>
                    </Row>
                    </Col>
                </Row>
                <div class="line-gray"></div>
                <Row class="margin-top-8">
                    <Col span="8">
                    <p class="notwrap">上次登录时间:</p>
                    </Col>
                    <Col span="16" class="padding-left-8">2017.09.12-13:32:20</Col>
                </Row>
                <Row class="margin-top-8">
                    <Col span="8">
                    <p class="notwrap">上次登录地点:</p>
                    </Col>
                    <Col span="16" class="padding-left-8">北京</Col>
                </Row>

            </Card>
            </Col>
        </Row>
    </Card>
</template>

<script>
export default {
  data() {
    return {
      count: {
        createUser: 496,
        visit: 3264,
        collection: 24389305,
        transfer: 39503498
      }
    };
  },
  computed: {},
  methods: {}
};
</script>
