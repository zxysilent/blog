<style lang="less" scoped>
.ivu-form-inline .ivu-form-item {
	margin-bottom: 10px;
}
</style>
<template>
	<Card dis-hover>
		<Form inline>
			<!-- <FormItem>
				<Select v-model="cid" placeholder="请选择文章类别" style="width:200px">
					<Option v-for="item in columnAll" :value="item.id" :key="item.id">{{ item.name }}</Option>
				</Select>
			</FormItem> -->
			<FormItem>
				<Button type="info" @click="init" icon="md-refresh" title="刷新数据">刷&nbsp;&nbsp;新</Button>
			</FormItem>
			<!-- <FormItem>
				<Button type="warning" @click="reCache" icon="ios-alert-outline" :title="'重新加载数据库缓存'+'\n'+'适用于直接修改数据库'">重载缓存</Button>
			</FormItem> -->
		</Form>
		<Table row-key="id" size="small" :loading="loading" border :columns="tabCol" :data="tabData"></Table>
		<Drawer title="授权" v-model="drawer" width="520" :mask-closable="false">
			<Tree :data="menuTree" show-checkbox ></Tree>
		</Drawer>
	</Card>
</template>
<script>
import { admMenuTree } from "@/api/menu";
import { admRoleAll, admRoleDrop } from "@/api/role";
export default {
	data() {
		return {
			drawer: true,
			formData: {
				name: "",
				url: "",
				owner: "",
				type: "",
				approver: "",
				date: "",
				desc: ""
			},
            menuTree:[],
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
										this.$router.push({
											name: "menu-edit",
											params: { id: data.row.id }
										});
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
			tabData: []
		};
	},
	methods: {
		preinit() {
			admMenuTree({}).then((resp) => {
				if (resp.code == 200) {
					this.menuTree = resp.data;
				} else {
					this.$Message.error({
						content: resp.msg,
						duration: 3
					});
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
		}
	},
	created() {
		this.preinit();
		this.init();
	}
};
</script>