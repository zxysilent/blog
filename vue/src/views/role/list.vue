<style lang="less" scoped>
.drawer-content {
	overflow: auto;
}

.drawer-content::-webkit-scrollbar {
	display: none;
}
.drawer-footer {
	z-index: 10;
	width: 100%;
	position: absolute;
	bottom: 0;
	left: 0;
	border-top: 1px solid #e8e8e8;
	padding: 10px 16px;
	text-align: left;
	background: #fff;
}
</style>
<template>
	<Card dis-hover>
		<Form inline>
			<FormItem>
				<Button type="info" @click="init" icon="md-refresh" title="刷新数据">刷&nbsp;&nbsp;新</Button>
			</FormItem>
			<!-- <FormItem>
				<Button type="warning" @click="reCache" icon="ios-alert-outline" :title="'重新加载数据库缓存'+'\n'+'适用于直接修改数据库'">重载缓存</Button>
			</FormItem> -->
		</Form>
		<Table row-key="id" size="small" :loading="loading" border :columns="tabCol" :data="tabData"></Table>
		<Drawer title="授权" v-model="drawer" width="520">
			<!-- <Tree :data="menuTree" show-checkbox></Tree> -->
			<div :style="{ maxHeight: maxHeight }" class="drawer-content">
				<div style="position: relative">
					<Tree ref="tree" :data="menuTree" show-checkbox :check-directly="true" :render="renderContent"></Tree>
					<Spin size="large" fix v-if="treeLoading"></Spin>
				</div>
			</div>
			<div class="drawer-footer">
				<Button type="primary" :loading="roleEditLoading" @click="emitRoleEdit">提交保存</Button>
				<!-- <Button @click="selectTreeAll">全选/反选</Button> -->
				<!-- @on-change="changeOpen" -->
				<!-- <Select v-model="openLevel" style="width: 110px" transfer>
					<Option value="0">展开所有</Option>
					<Option value="1">收合所有</Option>
				</Select> -->
				<Button type="text" @click="drawer=false">取消</Button>
			</div>
		</Drawer>
	</Card>
</template>
<script>
import { admMenuTree } from "@/api/menu";
import { admRoleAll, admRoleDrop, admRoleMenuAll, adadmRoleMenuEdit } from "@/api/role";
export default {
	data() {
		return {
			drawer: false,
			maxHeight: "100px",
			treeLoading: false,
			roleEditLoading: false,
			menuTree: [],
			tabCol: [
				{ key: "id", width: 60, align: "center" },
				{ title: "角色名称", key: "name", width: 150 },
				{ title: "角色说明", width: 250, key: "intro", ellipsis: true, tooltip: true },
				{
					title: "操作",
					key: "action",
					width: 150,
					align: "center",
					render: (h, data) => {
						// let allow = this.getRole(49);
						//if (allow) {
						return h("a", [
							h("Icon", {
								props: { type: "md-cog", size: "20", color: "#01AAED" },
								attrs: { title: "授权" },
								style: { marginRight: "10px" },
								on: {
									click: () => {
										// this.$router.push({
										// 	name: "menu-edit",
										// 	params: { id: data.row.id }
										// });
										this.role_id = data.row.id;
										this.menu_ids = [];
										this.roleMenuAll(this.role_id);
										this.drawer = true;
									}
								}
							}),
							h("Icon", {
								props: { type: "md-create", size: "20", color: "#FFB800" },
								attrs: { title: "修改" },
								style: { marginRight: "10px" },
								on: {
									click: () => {
										this.$router.push({ name: "menu-edit", params: { id: data.row.id } });
									}
								}
							}),
							h(
								"Poptip",
								{
									props: { confirm: true, title: "确定要删除吗？" },
									on: {
										"on-ok": () => {
											this.emitDrop(data);
										}
									}
								},
								[
									h("Icon", {
										props: { type: "ios-trash", size: "20", color: "#FF5722" },
										attrs: { title: "删除" }
									})
								]
							)
						]);
						// }
					}
				}
			],
			loading: false,
			tabData: [],
			//角色菜单
			role_id: 0,
			menu_ids: []
		};
	},
	methods: {
		preinit() {
			admMenuTree({}).then((resp) => {
				if (resp.code == 200) {
					this.menuTree = resp.data;
				} else {
					this.$Message.error({ content: resp.msg, duration: 3 });
				}
			});
		},
		init() {
			this.loading = true;
			admRoleAll({ table: true }).then((resp) => {
				if (resp.code == 200) {
					this.tabData = resp.data;
				} else {
					this.tabData = [];
				}
				this.loading = false;
			});
		},
		renderContent(h, { root, node, data }) {
			return h("span", { style: { display: "inline-block", width: "100%" } }, [
				h("span", [
					h("Icon", {
						props: { type: data.icon, size: "16" },
						style: { marginRight: "8px" }
					}),
					h("span", data.title)
				])
			]);
		},
		//删除
		emitDrop(data) {
			admRoleDrop({ id: data.row.id }).then((resp) => {
				if (resp.code == 200) {
					this.$Message.success({
						content: "删除成功",
						onClose: () => {
							this.tabData.splice(data.index, 1);
						}
					});
				} else {
					this.$Message.error(resp.msg);
				}
			});
		},
		// 提交权限
		emitRoleEdit() {
			this.roleEditLoading = true;
			let menuIds = [];
			let selectedNodes = this.$refs.tree.getCheckedAndIndeterminateNodes();
			selectedNodes.forEach(function (item) {
				menuIds.push(item.id);
			});
			this.menu_ids = menuIds;
			adadmRoleMenuEdit({
				role_id: this.role_id,
				menu_ids: this.menu_ids
			}).then((resp) => {
				this.roleEditLoading = false;
				if (resp.code == 200) {
					this.$Message.success({
						content: "授权成功",
						onClose: () => {
						}
					});
				} else {
					this.$Message.error({ content: resp.msg, duration: 3 });
				}
				// if (res.success) {
				// 	this.$Message.success("操作成功");
				// 	// 标记重新获取菜单数据
				// 	this.$store.commit("setAdded", false);
				// 	util.initRouter(this);
				// 	this.getDataList();
				// 	this.permModalVisible = false;
				// }
			});
		},
		//
		roleMenuAll(role_id) {
			admRoleMenuAll({ id: role_id }).then((resp) => {
				if (resp.code == 200) {
					console.log(resp);
				} else {
					this.$Message.error({ content: resp.msg, duration: 3 });
				}
			});
		}
	},
	created() {
		this.maxHeight = Number(document.documentElement.clientHeight - 121) + "px";
		this.preinit();
		this.init();
	}
};
</script>