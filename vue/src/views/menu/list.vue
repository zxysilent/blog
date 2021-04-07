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
	</Card>
</template>
<script>
import { admMenuTree, admMenuDrop, admMenuEditShow } from "@/api/menu";
export default {
	data() {
		return {
			tabCol: [
				{ key: "id", width: 60, align: "center" },
				{ title: "标题", key: "title", width: 150, tree: true },
				{ title: "名称", key: "name", width: 150 },
				{ title: "访问路径", key: "path", width: 150 },
				{
					title: "导航显示",
					width: 100,
					key: "show",
					render: (h, data) => {
						let that = this;
						return h(
							"i-switch",
							{
								props: { value: data.row.show, size: "large" },
								on: {
									"on-change": function (val) {
										admMenuEditShow({ id: data.row.id, show: val }).then((resp) => {
											if (resp.code != 200) {
												that.$Message.error({
													content: resp.msg,
													onClose: () => {
														that.init();
													}
												});
											}
										});
									}
								}
							},
							[
								h("span", { slot: "open", domProps: { innerHTML: "显示" } }),
								h("span", { slot: "close", domProps: { innerHTML: "隐藏" } })
							]
						);
					}
				},
				{
					title: "图标",
					width: 80,
					key: "icon",
					align: "center",
					render: (h, data) => {
						return h("Icon", {
							props: { type: data.row.icon, size: "24" },
							attrs: { title: data.row.icon }
						});
					}
				},
				{ title: "组件", width: 200, key: "comp" },
				{ title: "序号", width: 80, key: "sort" },
				{
					title: "操作",
					key: "action",
					width: 100,
					align: "center",
					render: (h, data) => {
						// let allow = this.getRole(49);
						//if (allow) {
						return h("a", [
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
		init() {
			this.loading = true;
			admMenuTree({}).then((resp) => {
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
			admMenuDrop({ id: data.row.id }).then((resp) => {
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
		}
	},
	created() {
		this.init();
	}
};
</script>