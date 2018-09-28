<style lang="less">
@import "./main.less";
</style>
<template>
    <div class="main" :class="{'main-hide-text': shrink}">
        <div class="sidebar-menu-con">
            <div :style="{background: '#fff'}" class="ivu-shrinkable-menu">
                <div class="logo-con">
                    <img src="../images/logo.png" alt="" srcset="">
                    <p class="main-user-name">{{ user.name }}</p>
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
                    <Submenu name="博客管理">
                        <template slot="title">
                            <Icon type="ios-ionic-outline" />
                            博客管理
                        </template>
                        <MenuGroup title="post">
                            <MenuItem name="post-list" to="/post/list">文章列表</MenuItem>
                            <MenuItem name="post-add">添加文章</MenuItem>
                        </MenuGroup>
                        <MenuGroup title="page">
                            <MenuItem name="page-list">页面列表</MenuItem>
                            <MenuItem name="page-add">添加页面</MenuItem>
                        </MenuGroup>
                    </Submenu>
                    <Submenu name="类型管理">
                        <template slot="title">
                            <Icon type="ios-school-outline" />
                            类型管理
                        </template>
                        <MenuGroup title="cate">
                            <MenuItem name="cate-list">分类列表</MenuItem>
                            <MenuItem name="cate-add">添加分类</MenuItem>
                        </MenuGroup>
                        <MenuGroup title="tag">
                            <MenuItem name="tag-list">标签列表</MenuItem>
                            <MenuItem name="tag-add">添加标签</MenuItem>
                        </MenuGroup>
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
        </div>
        <div class="main-header-con">
            <div class="main-header">
                <div class="navicon-con">
                    <Button icon="ios-home-outline" type="text" shape="circle" @click="toggleClick">
                    </Button>
                </div>
                <div class="header-middle-con">
                    <div class="main-breadcrumb">
                        sss
                    </div>
                </div>
                <div class="header-avator-con">
                    <div class="user-dropdown-menu-con">
                        <Row type="flex" justify="end" align="middle" class="user-dropdown-innercon">
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
                            <!-- <Avatar size="large" :src="avatorPath" style="background: #619fe7;margin-left: 10px;"></Avatar> -->
                        </Row>
                    </div>
                </div>
            </div>
        </div>
        <div class="single-page-con">
            <div class="single-page">
                <keep-alive :include="cachePage">
                    <router-view></router-view>
                </keep-alive>
            </div>
        </div>
    </div>
</template>
<script>
import lockScreen from "./main-components/lockscreen/lockscreen.vue";
import messageTip from "./main-components/message-tip.vue";
import util from "@/libs/util.js";
import { auth } from "@/api/auth";
export default {
  components: {
    lockScreen,
    messageTip
  },
  data() {
    return {
      user: {
        name: "--",
        num: "--",
        avatar: ""
      },
      shrink: false,
      userName: "",
      isFullScreen: false,
      openedSubmenuArr: this.$store.state.app.openedSubmenuArr
    };
  },
  computed: {
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
