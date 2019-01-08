<style lang="less">
@import "./main.less";
</style>
<template>
    <Layout style="height: 100%" class="main">
        <Sider hide-trigger :width="200" class="left-sider" :style="{overflow: 'hidden'}">
            <div class="sidebar-menu-con">
                <div class="logo-con">
                    <img src="../images/logo.png" alt="" srcset="">
                    <div class="header-user">{{ user.name }}</div>
                    <div class="header-top-con">
                        <Tooltip content="前台主页" placement="bottom">
                           <a href="/" target="_blank"><Icon type="md-paper-plane" size="20" /></a>
                        </Tooltip>
                        <Divider type="vertical" />
                        <Tooltip content="个人中心" placement="bottom">
                            <router-link to="/core/self" tag="span">
                                <Icon type="ios-person" size="20" />
                            </router-link>
                        </Tooltip>
                        <Divider type="vertical" />
                        <Tooltip content="退出登陆" placement="bottom">
                            <Icon type="md-log-out" size="20" @click="logout" />
                        </Tooltip>
                    </div>
                </div>
                <Menu :active-name="$route.name" :open-names="$route.name.split('-')" theme="dark" width="auto" accordion>
                    <MenuItem name="home" to="/core/home">
                    <Icon type="ios-home-outline" /> 管理主页</MenuItem>
                    <Submenu name="post">
                        <template slot="title">
                            <Icon type="ios-megaphone-outline" />
                            博客管理
                        </template>
                        <MenuItem name="post-list" to="/core/post/list">
                        <Icon type="ios-list-box-outline" />文章列表</MenuItem>
                        <MenuItem name="post-add" to="/core/post/add">
                        <Icon type="ios-add-circle-outline" />添加文章</MenuItem>
                    </Submenu>
                    <Submenu name="page">
                        <template slot="title">
                            <Icon type="ios-map-outline" />
                            页面管理
                        </template>
                        <MenuItem name="page-list" to="/core/page/list">
                        <Icon type="ios-list-box-outline" />页面列表</MenuItem>
                        <MenuItem name="page-add" to="/core/page/add">
                        <Icon type="ios-add-circle-outline" />添加页面</MenuItem>
                    </Submenu>
                    <Submenu name="cate">
                        <template slot="title">
                            <Icon type="ios-school-outline" />
                            分类管理
                        </template>
                        <MenuItem name="cate-list" to="/core/cate/list">
                        <Icon type="ios-list-box-outline" />分类列表</MenuItem>
                        <MenuItem name="cate-add" to="/core/cate/add">
                        <Icon type="ios-add-circle-outline" />添加分类</MenuItem>
                    </Submenu>
                    <Submenu name="tag">
                        <template slot="title">
                            <Icon type="ios-pricetags-outline" />
                            标签管理
                        </template>
                        <MenuItem name="tag-list" to="/core/tag/list">
                        <Icon type="ios-list-box-outline" />标签列表</MenuItem>
                        <MenuItem name="tag-add" to="/core/tag/add">
                        <Icon type="ios-add-circle-outline" />添加标签</MenuItem>
                    </Submenu>
                    <Submenu name="setting">
                        <template slot="title">
                            <Icon type="ios-cog-outline" />
                            系统设置
                        </template>
                        <MenuItem name="setting-base" to="/core/setting/base">
                        <Icon type="ios-cog-outline" />基本设置</MenuItem>
                        <MenuItem name="setting-comment" to="/core/setting/comment">
                        <Icon type="ios-text-outline" />评论设置</MenuItem>
                        <MenuItem name="setting-analytic" to="/core/setting/analytic">
                        <Icon type="ios-pulse" /> 统计设置</MenuItem>
                        <MenuItem name="setting-custom" to="/core/setting/custom">
                        <Icon type="ios-code-working" /> 自 定 义</MenuItem>
                    </Submenu>
                </Menu>
            </div>
        </Sider>
        <Layout>
            <Content class="main-content-con">
                <Layout class="main-layout-con">
                    <Content class="content-wrapper">
                        <!-- <keep-alive> -->
                        <router-view :key="rkey" />
                        <!-- </keep-alive> -->
                    </Content>
                </Layout>
            </Content>
        </Layout>
    </Layout>
</template>
<script>
import messageTip from "./main-components/message-tip.vue";
import util from "@/util/util.js";
import { auth } from "@/api/auth";
export default {
  components: {
    messageTip
  },
  data() {
    return {
      rkey: 6655,
      user: {
        name: "--",
        num: "--"
      }
    };
  },
  computed: {
    cachePage() {
      return this.$store.state.app.cachePage;
    },
    mesCount() {
      return this.$store.state.app.messageCount;
    }
  },
  methods: {
    init() {
      auth().then(resp => {
        if (resp.code == 200) {
          this.user = resp.data;
        }
      });
    },
    logout() {
      this.$store.commit("logout", this);
      this.$router.push({
        name: "login"
      });
    }
  },
  watch: {
    $route(to, from) {
      //   console.log("to.name=>", to.name);
      // 强制 触发生命周期
      if (
        to.name == "post-edit" ||
        to.name == "post-add" ||
        to.name == "page-edit" ||
        to.name == "page-add"
      ) {
        this.rkey = new Date().getTime();
      } else {
        this.rkey = 6655;
      }
    }
  },
  created() {
    this.init();
  }
};
</script>
