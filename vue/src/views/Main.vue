<style lang="less">
@import "./main.less";
</style>
<template>
    <Layout style="height: 100%" class="main">
        <Sider hide-trigger :width="200" class="left-sider" :style="{overflow: 'hidden'}">
            <div class="sidebar-menu-con">
                <div class="logo-con">
                    <img src="../images/logo.png" alt="" srcset="">
                    <div>{{ user.name }}</div>
                    <div class="header-avator-con">
                        <Tooltip content="主页" placement="bottom">
                            <router-link to="/home" tag="span">
                                <Icon type="ios-home-outline" size="20" />
                            </router-link>
                        </Tooltip>
                        <Divider type="vertical" />
                        <Tooltip content="个人中心" placement="bottom">
                            <router-link to="/self" tag="span">
                                <Icon type="ios-person-outline" size="20" />
                            </router-link>
                        </Tooltip>
                        <Divider type="vertical" />
                        <message-tip v-model="mesCount"></message-tip>
                    </div>
                </div>
                <Menu :active-name="$route.name" :open-names="$route.name.split('-')" width="auto" accordion>
                    <!-- <MenuItem name="home"  to="/home">
                    <Icon type="ios-home-outline" /> 主页</MenuItem> -->
                    <Submenu name="post">
                        <template slot="title">
                            <Icon type="ios-megaphone-outline" />
                            博客管理
                        </template>
                        <MenuItem name="post-list" to="/post/list">
                        <Icon type="ios-list-box-outline" />文章列表</MenuItem>
                        <MenuItem name="post-add"  to="/post/add">
                        <Icon type="ios-add-circle-outline" />添加文章</MenuItem>
                    </Submenu>
                    <Submenu name="page">
                        <template slot="title">
                            <Icon type="ios-map-outline" />
                            页面管理
                        </template>
                        <MenuItem name="page-list"  to="/page/list">
                        <Icon type="ios-list-box-outline" />页面列表</MenuItem>
                        <MenuItem name="page-add">
                        <Icon type="ios-add-circle-outline" />添加页面</MenuItem>
                    </Submenu>
                    <Submenu name="cate">
                        <template slot="title">
                            <Icon type="ios-school-outline" />
                            分类管理
                        </template>
                        <MenuItem name="cate-list" to="/cate/list">
                        <Icon type="ios-list-box-outline" />分类列表</MenuItem>
                        <MenuItem name="cate-add" to="/cate/add">
                        <Icon type="ios-add-circle-outline" />添加分类</MenuItem>
                    </Submenu>
                    <Submenu name="tag">
                        <template slot="title">
                            <Icon type="ios-pricetags-outline" />
                            标签管理
                        </template>
                        <MenuItem name="tag-list" to="/tag/list">
                        <Icon type="ios-list-box-outline" />标签列表</MenuItem>
                        <MenuItem name="tag-add"  to="/tag/add">
                        <Icon type="ios-add-circle-outline" />添加标签</MenuItem>
                    </Submenu>
                    <Submenu name="setting">
                        <template slot="title">
                            <Icon type="ios-cog-outline" />
                            系统设置
                        </template>
                        <MenuItem name="setting-base" to="/setting/base">
                        <Icon type="ios-cog-outline" />基本设置</MenuItem>
                        <MenuItem name="setting-comment" to="/setting/comment">
                        <Icon type="ios-text-outline" />评论设置</MenuItem>
                        <MenuItem name="setting-analytic" to="/setting/analytic">
                        <Icon type="ios-pulse" /> 统计设置</MenuItem>
                        <MenuItem name="setting-custom" to="/setting/custom">
                        <Icon type="ios-code-working" /> 自 定 义</MenuItem>
                    </Submenu>
                </Menu>
            </div>
        </Sider>
        <Layout>
            <Header class="header-con">
                <div class="header-bar">
                    <div class="custom-content-con">
                        <Dropdown transfer trigger="click" @on-click="handleClickUserDropdown">
                            <a href="javascript:void(0)">
                                <span class="main-user-name">{{ user.name }}</span>
                                <Icon type="md-arrow-dropdown"></Icon>
                            </a>
                            <DropdownMenu slot="list">
                                <DropdownItem name="self">个人中心</DropdownItem>
                                <DropdownItem name="loginout" divided>退出登录</DropdownItem>
                            </DropdownMenu>
                        </Dropdown>
                    </div>
                </div>
            </Header>
            <Content class="main-content-con">
                <Layout class="main-layout-con">
                    <Content class="content-wrapper">
                        <keep-alive>
                            <router-view />
                        </keep-alive>
                    </Content>
                </Layout>
            </Content>
        </Layout>
    </Layout>
</template>
<script>
import messageTip from "./main-components/message-tip.vue";
import util from "@/libs/util.js";
import { auth } from "@/api/auth";
export default {
  components: {
    messageTip
  },
  data() {
    return {
      collapsed: false,
      user: {
        name: "--",
        num: "--",
      },
      shrink: false,
      userName: "",
      isFullScreen: false,
      openedSubmenuArr: this.$store.state.app.openedSubmenuArr
    };
  },
  computed: {
    menuitemClasses: function() {
      return ["menu-item", this.isCollapsed ? "collapsed-menu" : ""];
    },
    menuList() {
      return this.$store.state.app.menuList;
    },
    avatorPath() {
      return localStorage.avatorImgPath;
    },
    cachePage() {
      return this.$store.state.app.cachePage;
    },
    menuTheme() {
      return this.$store.state.app.menuTheme;
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
      this.$store.commit("updateMenulist");
      let messageCount = 3;
      this.messageCount = messageCount.toString();
      this.$store.commit("setMessageCount", 3);
    },
    toggleClick() {
      this.shrink = !this.shrink;
    },
    handleClickUserDropdown(name) {
      if (name === "self") {
        this.$router.push({
          name: "self"
        });
      } else if (name === "loginout") {
        // 退出登录
        this.$store.commit("logout", this);
        this.$store.commit("clearOpenedSubmenu");
        this.$router.push({
          name: "login"
        });
      }
    },
    handleSubmenuChange(val) {
      // console.log(val)
    },
    beforePush(name) {
      // if (name === 'accesstest_index') {
      //     return false;
      // } else {
      //     return true;
      // }
      return true;
    },
    fullscreenChange(isFullScreen) {
      // console.log(isFullScreen);
    }
  },
  watch: {
    $route(to) {
      localStorage.currentPageName = to.name;
    }
  },
  created() {
    this.init();
    // 显示打开的页面的列表
    this.$store.commit("setOpenedList");
  }
};
</script>
