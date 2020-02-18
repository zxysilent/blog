<style lang="less" scoped>
@import "./layout.less";
</style>
<template>
	<Layout class="layout">
		<Sider hide-trigger :width="200" class="layout-sider" :style="{overflow: 'hidden'}">
			<div class="sider-container">
				<div class="sider-header ">
					<img src="../assets/logo.png" alt="" srcset="">
					<div class="header-user">{{ user.name }}</div>
					<div class="header-menu">
						<Tooltip content="前台主页" placement="bottom">
							<a href="/" target="_blank">
								<Icon type="md-paper-plane" size="20" /></a>
						</Tooltip>
						<Divider type="vertical" />
						<Tooltip content="个人中心" placement="bottom">
							<router-link to="/self" tag="span">
								<Icon type="ios-person" size="20" />
							</router-link>
						</Tooltip>
						<Divider type="vertical" />
						<Tooltip content="退出登陆" placement="bottom">
							<Icon type="md-log-out" size="20" @click="logout" />
						</Tooltip>
					</div>
				</div>
				<Menu ref="menu" :active-name="$route.name" :open-names="$route.name.split('-')" theme="dark" width="auto" accordion>
					<MenuItem name="home" to="/home">
					<Icon type="ios-home-outline" /> 管理主页</MenuItem>
					<Submenu name="post">
						<template slot="title">
							<Icon type="ios-megaphone-outline" />
							博客管理
						</template>
						<MenuItem name="post-list" to="/post/list">
						<Icon type="ios-list-box-outline" />文章列表</MenuItem>
						<MenuItem name="post-add" to="/post/add">
						<Icon type="ios-add-circle-outline" />添加文章</MenuItem>
					</Submenu>
					<Submenu name="page">
						<template slot="title">
							<Icon type="ios-map-outline" />
							页面管理
						</template>
						<MenuItem name="page-list" to="/page/list">
						<Icon type="ios-list-box-outline" />页面列表</MenuItem>
						<MenuItem name="page-add" to="/page/add">
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
						<MenuItem name="tag-add" to="/tag/add">
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
		<Content class="layout-content">
			<!-- <keep-alive> -->
			<router-view :key="rkey" />
			<!-- </keep-alive> -->
		</Content>
	</Layout>
</template>
<script>
import Utils from "@/utils.js";
import { admAuth } from "@/api/auth";
export default {
	data() {
		return {
			rkey: 6655,
			user: { name: "--", num: "--" }
		};
	},
	computed: {
		// cachePage() {
		// 	return this.$store.state.app.cachePage;
		// }
	},
	methods: {
		init() {
			admAuth().then(resp => {
				if (resp.code == 200) {
					this.user = resp.data;
				}
			});
		},
		logout() {
			// this.$store.commit("logout", this);
			// 移除 vuex
			localStorage.clear();
			sessionStorage.clear();
			this.$router.push({ name: "login" });
		}
	},
	updated() {
		// 手动更新展开的子目录
		this.$nextTick(() => {
			console.log("nextTick");
			if (this.$refs.menu) {
				this.$refs.menu.updateOpened();
			}
		});
	},
	watch: {
		$route(to, from) {
			// console.log("to.name=>", to.name);
			// 强制触发生命周期
			if (to.name == "post-edit" || to.name == "post-add" || to.name == "page-edit" || to.name == "page-add") {
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
