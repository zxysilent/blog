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
								<Icon type="md-paper-plane" size="20" />
							</a>
						</Tooltip>
						<Divider type="vertical" />
						<Tooltip content="个人中心" placement="bottom">
							<router-link to="/profile" tag="span">
								<Icon type="ios-person" size="20" />
							</router-link>
						</Tooltip>
						<Divider type="vertical" />
						<Tooltip content="退出登陆" placement="bottom">
							<Icon type="md-log-out" size="20" @click="logout" />
						</Tooltip>
					</div>
				</div>
				<Menu ref="menu" :active-name="$route.name.split('-')[0]" :open-names="[$route.meta.module,$route.name.split('-')]" theme="dark" width="auto" accordion>
					<MenuItem name="home" to="/home">
					<Icon type="ios-home-outline" />管理主页</MenuItem>
					<Submenu name="sys">
						<template slot="title">
							<Icon type="ios-apps-outline" />系统模块
						</template>
						<MenuItem name="global" to="/global">
						<Icon type="ios-cog-outline" />配置管理</MenuItem>
                        <MenuItem name="dict" to="/dict/list">
						<Icon type="ios-key-outline" />字典管理</MenuItem>
						<MenuItem name="profile" to="/profile">
						<Icon type="ios-person-outline" />个人信息</MenuItem>
					</Submenu>
					<Submenu name="app">
						<template slot="title">
							<Icon type="ios-apps-outline" />博客模块
						</template>
						<MenuItem name="post" to="/post/list">
						<Icon type="ios-megaphone-outline" />文章管理</MenuItem>
						<MenuItem name="page" to="/page/list">
						<Icon type="ios-list-box-outline" />页面管理</MenuItem>
						<MenuItem name="cate" to="/cate/list">
						<Icon type="ios-school-outline" />分类管理</MenuItem>
						<MenuItem name="tag" to="/tag/list">
						<Icon type="ios-pricetags-outline" />标签管理</MenuItem>
					</Submenu>
					<!-- <Submenu name="setting">
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
					</Submenu> -->
					<!-- <MenuItem name="home" to="/home">
					<Icon type="ios-home-outline" /> 管理主页</MenuItem>
					<template v-for="item in getMenus">
						<Submenu v-if="item.children && item.show" :name="item.name" :key="item.name">
							<template slot="title">
								<Icon :type="item.icon"></Icon>
								{{item.title}}
							</template>
							<MenuItem :name="itm.name" v-for="itm in item.children" v-if="itm.show" :to="itm.path" :key="itm.name">
							<Icon :type="itm.icon"></Icon>{{itm.title}}
							</MenuItem>
						</Submenu>
						<MenuItem v-else-if="item.show" :name="item.name" :to="item.path" :key="item.name">
						<Icon :type="item.icon"></Icon>
						{{item.title}}
						</MenuItem>
					</template> -->
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
// import { mapGetters, mapMutations } from "vuex";
import Storage from "@/utils/storage";
import { admAuthGet } from "@/api/auth";
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
		// 使用对象展开运算符将 getter 混入 computed 对象中
		// ...mapGetters(["getMenus"])
	},
	methods: {
		init() {
			admAuthGet().then((resp) => {
				if (resp.code == 200) {
					this.user = resp.data;
				}
			});
		},
		logout() {
			// this.$store.commit("logout", this);
			// 移除 vuex
			Storage.clear();
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
<style lang="less" scoped>
.layout {
	height: 100%;
	&-content {
		padding: 4px;
		height: ~"calc(100% - 10px)";
		overflow: auto;
		min-width: 500px;
	}
	&-sider {
		.ivu-layout-sider-children {
			overflow-y: scroll;
			margin-right: -18px;
		}
		.sider-container {
			height: 100%;
			width: 200px;
			transition: width 0.3s;
			.sider-header {
				padding: 8px 8px 0px 8px;
				text-align: center;
				color: #f5f5f5;
				a {
					color: #f5f5f5;
				}
				a:hover {
					color: #2d8cf0;
				}
				img {
					width: 90%;
					max-width: 120px;
					border-radius: 50%;
				}
				.header-user {
					font-size: 20px;
					padding: 2px;
				}
				.header-menu {
					padding: 2px;
					.ivu-tooltip {
						cursor: pointer;
					}
					.ivu-tooltip:hover {
						color: #2d8cf0;
					}
				}
			}
		}
	}
}
.ivu-menu-item > i {
	margin-right: 8px !important;
}
.ivu-menu-submenu > .ivu-menu > .ivu-menu-item > i {
	margin-right: 8px !important;
}
.collased-menu-dropdown {
	width: 100%;
	margin: 0;
	line-height: normal;
	padding: 7px 0 6px 16px;
	clear: both;
	font-size: 12px !important;
	white-space: nowrap;
	list-style: none;
	cursor: pointer;
	transition: background 0.2s ease-in-out;
	&:hover {
		background: rgba(100, 100, 100, 0.1);
	}
	& * {
		color: #515a6e;
	}
	.ivu-menu-item > i {
		margin-right: 12px !important;
	}
	.ivu-menu-submenu > .ivu-menu > .ivu-menu-item > i {
		margin-right: 8px !important;
	}
}

.ivu-select-dropdown.ivu-dropdown-transfer {
	max-height: 400px;
}
.ivu-layout-sider {
	background: #515a6e;
}
.ivu-menu-vertical.ivu-menu-light:after {
	width: 0;
}
</style>