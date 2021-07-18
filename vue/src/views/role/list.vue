	
<style lang="less" scoped>
.ivu-form-item {
	margin-bottom: 10px;
}
</style>
<template>
	<Card dis-hover>
		<Form inline>
			<!-- <FormItem>
				<Select v-model="cid" placeholder="请选择类别" style="width:200px">
					<Option v-for="item in columnAll" :value="item.id" :key="item.id">{{ item.name }}</Option>
				</Select>
			</FormItem> -->
			<FormItem>
				<Button v-auth="'role_show'" type="info" @click="init" icon="md-refresh" title="刷新数据">刷&nbsp;&nbsp;新</Button>
			</FormItem>
		</Form>
		<Table size="small" :columns="tabCol" :data="tabData"></Table>
		<Drawer title="角色授权" v-model="showGrant" width="520">
			<!-- <Alert closable type="warning">推荐角色序号越大,权限越少</Alert> -->
			<Form ref="grantForm" :model="grantForm" label-position="top" :show-message="false" style="padding-bottom:55px">
				<template v-for="(item,key) in grantTree">
					<FormItem :label="key">
						<CheckboxGroup v-model="grantForm.grant_ids">
							<Checkbox v-for="itm in item" :label="itm.id" :key="item.id" border>{{itm.name}}</Checkbox>
						</CheckboxGroup>
					</FormItem>
				</template>
			</Form>
			<div style="text-align: center;width: 100%;position:fixed;bottom: 0;right: 0;border-top: 1px solid #e8e8e8;padding: 10px 16px;background: #fff;width:520px">
				<ButtonGroup>
					<Button type="warning" :loading="loading" @click="emitGrant">提交保存</Button>
					<Button type="info" style="margin-left: 8px" @click="showGrant=false">取消关闭</Button>
				</ButtonGroup>
			</div>
		</Drawer>
	</Card>
</template>
<script>
import { admRoleAll, admRoleDrop, admGrantTree, adadmRoleGrantEdit, admRoleGrantAll } from "@/api/role";
export default {
	data() {
		return {
			showGrant: false,
			loading: false,
			grantTree: [],
			grantForm: { role_id: 0, grant_ids: [] },
			tabCol: [
				{ type: "index", width: 80, align: "center" },
				{ title: "角色名称", key: "name" },
				{ title: "角色说明", key: "intro", ellipsis: true, tooltip: true },
				{
					title: "录入日期",
					width: 200,
					render: (h, data) => {
						return h("div", data.row.ctime.replace(/T|\+08:00/g, " "));
					}
				},
				{
					title: "操作",
					width: 200,
					align: "center",
					render: (h, data) => {
						console.log(this.$auth("role_add"));
						let actios = [];
						if (this.$auth("role_grant")) {
							actios.push(
								h("Icon", {
									props: { type: "md-cog", size: "20", color: "#01AAED" },
									attrs: { title: "授权" },
									style: { marginRight: "15px" },
									on: {
										click: () => {
											this.grantForm.role_id = data.row.id;
											this.preinit();
											this.initGrant();
											this.showGrant = true;
										}
									}
								})
							);
						}
						if (this.$auth("role_edit")) {
							actios.push(
								h("Icon", {
									props: { type: "md-create", size: "20", color: "#FFB800" },
									attrs: { title: "修改" },
									style: { marginRight: "15px" },
									on: {
										click: () => {
											this.$router.push({
												name: "role-edit",
												params: { id: data.row.id }
											});
										}
									}
								})
							);
						}
						if (this.$auth("role_drop")) {
							actios.push(
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
							);
						}
						return h("a", actios);
					}
				}
			],
			tabData: []
		};
	},
	methods: {
		preinit() {
			admGrantTree().then((resp) => {
				if (resp.code == 200) {
					this.grantTree = resp.data;
				} else {
					this.grantTree = [];
					this.$Message.error({ content: resp.msg, duration: 3 });
				}
			});
		},
		initGrant() {
			admRoleGrantAll({ id: this.grantForm.role_id }).then((resp) => {
				if (resp.code == 200) {
					this.grantForm.grant_ids = resp.data.map((item) => {
						return item.id;
					});
				} else {
					this.grantForm.grant_ids = [];
					this.$Message.error({ content: resp.msg, duration: 3 });
				}
			});
		},
		init() {
			admRoleAll().then((resp) => {
				if (resp.code == 200) {
					this.tabData = resp.data;
				} else {
					this.tabData = [];
					this.$Message.error({ content: resp.msg, duration: 3 });
				}
			});
		},
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
					this.$Message.error({ content: resp.msg, duration: 3 });
				}
			});
		},
		emitGrant() {
			console.log(this.grantForm);
			this.loading = true;
			adadmRoleGrantEdit(this.grantForm).then((resp) => {
				this.loading = false;
				if (resp.code == 200) {
					this.$Message.success({
						content: "授权成功",
						onClose: () => {}
					});
					// 重新拉取授权
					this.$store.dispatch("fetchGrant");
				} else {
					this.$Message.error({ content: resp.msg, duration: 3 });
				}
			});
		}
	},
	mounted() {
		this.preinit();
		this.init();
	}
};
</script>
